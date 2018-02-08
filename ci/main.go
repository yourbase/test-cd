// CIBot is a git repo webhook handler that runs CI jobs given an event.
//
// See docs/design/CIbot.md for a discussion.
//
// For now, it relies on Bazel's sandboxing for security. This is probably not
// enough.
//
// It listens for a webhook event from github, then runs the ci.sh script.
//
// A local directory is used to cache build artifacts between container
// invocations, which speeds things up considerably.
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	// TODO:Use go-github instead. See for example:
	// https://github.com/mlarraz/threshold/blob/master/threshold.go
	"github.com/google/go-github/github"
	joonix "github.com/joonix/log"
	"github.com/phayes/hookserve/hookserve"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/yourbase/yourbase/bazel"
	"github.com/yourbase/yourbase/ci/logview"
)

type gitHubCredentials struct {
	// TODO: Consolidate into a token? Are we able to clone repos with the
	// token?
	GithubUsername string `split_words:"true"`
	GithubPassword string `split_words:"true" required:true`
	GithubToken    string `split_words:"true"`
}

const (
	gitBase   = "ci/github.com/"
	cacheBase = "cache/github.com/"
	ciShell   = "ci.sh"
	cdShell   = "cd.sh"
)

var (
	homeDir = "/"
)

var (
	port    = flag.Int("port", 8080, "port to serve requests")
	runOnce = flag.Bool("runOnce", false, "Exit after one execution, for testing")
	tmpBase = flag.Bool("tmpBase", false, "Use the Bazel test tmp directory to store files (bazel cache, git sources).")
)

func init() {
	log.SetFormatter(&joonix.FluentdFormatter{})
	usr, err := user.Current()
	if err != nil {
		return
	}
	homeDir = usr.HomeDir
}

// commandWithLog takes a cmd object and a logrus logger, executes the command, and logs
// the cmd's stderr and stdout to the logrus logger.
func commandWithLog(cmd *exec.Cmd, logger *log.Entry) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	multi := io.MultiReader(stdout, stderr)

	if err = cmd.Start(); err != nil {
		logger.WithField("status", err.Error()).Error("Start Error")
		return err
	}

	std := bufio.NewScanner(multi)

	lastLine := time.Now()
	for std.Scan() {
		logger.WithField("duration", time.Since(lastLine).Seconds()).Info(std.Text())
		lastLine = time.Now()
	}
	if err = cmd.Wait(); err != nil {
		logger.WithField("status", err.Error()).Error("Command Failed")
		return err
	}
	if err = std.Err(); err != nil {
		logger.WithField("status", err.Error()).Error("Error")
		return err
	}

	logger.WithField("status", cmd.ProcessState.String()).Info("Completed")
	return nil
}

type ciEvent struct {
	repo    string
	branch  string
	repoURL string
	commit  string
	owner   string
	build   string

	// Previous commit for this branch before this push. It's set to
	// 0000000000000000000000000000000000000000 if the branch is being created.
	beforeCommit string
}

type ciRunner struct {
	githubClient      *github.Client
	githubCredentials *gitHubCredentials
}

func (c *ciRunner) gitSetup(logger *log.Entry, event *ciEvent) error {
	owner, repo, commit := event.owner, event.repo, event.commit
	creds := c.githubCredentials

	auth := ""
	if creds.GithubUsername != "" || creds.GithubPassword != "" {
		auth = fmt.Sprintf("%s:%s@", creds.GithubUsername, creds.GithubPassword)
	}
	repoAddress := fmt.Sprintf("https://%sgithub.com/%s/%s.git", auth, owner, repo)
	base := homeDir
	if *tmpBase {
		base = os.Getenv("TEST_TMPDIR")
	}
	localGitDir := path.Join(base, gitBase, owner, repo)
	localCacheDir := path.Join(base, cacheBase, owner, repo)

	os.MkdirAll(localCacheDir, 0755)

	logger.Infoln("Opening git directory", localGitDir)
	r, err := git.PlainOpen(localGitDir)
	if err != nil {
		logger.Infof("Cloning https://github.com/%s/%s.git", owner, repo)
		r, err = git.PlainClone(localGitDir, false, &git.CloneOptions{
			URL:               repoAddress,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
	}
	if err != nil {
		// Nuke the directory because it might have the wrong auth bits and
		// the next attempt should have a fresh start.
		os.RemoveAll(localGitDir)
		return fmt.Errorf("PlainClone of %q failed: %v", repoAddress, err)
	}

	logger.Infoln("Fetching git remote")
	err = r.Fetch(&git.FetchOptions{RemoteName: "origin"})
	if err != git.NoErrAlreadyUpToDate && err != nil {
		return fmt.Errorf("Git Fetch failed: %v", err)
	}
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("Could not obtain worktree: %v", err)
	}

	logger.Infoln("Checking out commit", commit)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commit),
	})
	if err != nil {
		return fmt.Errorf("Checkout failed: %v", err)
	}
	if err := os.Chdir(localGitDir); err != nil {
		return fmt.Errorf("chdir %v: %v", localGitDir, err)
	}
	return nil
}

func (c *ciRunner) setCIStatus(ctx context.Context, event *ciEvent, status string, description string) error {
	return c.setRepoStatus(ctx, "https://yourbase.io/ci", event, status, description)
}

func (c *ciRunner) setCDStatus(ctx context.Context, event *ciEvent, status string, description string) error {
	return c.setRepoStatus(ctx, "https://yourbase.io/cd", event, status, description)
}

func (c *ciRunner) setRepoStatus(ctx context.Context, checkContext string, event *ciEvent, status string, description string) error {

	targetURL := logview.LogURL(event.owner, event.repo, event.branch, event.build)
	repoStatus := &github.RepoStatus{
		State:       &status,
		Description: &description,
		Context:     &checkContext,
		TargetURL:   &targetURL,
	}
	_, _, err := c.githubClient.Repositories.CreateStatus(ctx, event.owner, event.repo, event.commit, repoStatus)
	return err
}

// runBazelCommand executes a script for Bazel repositories. This is not safe
// for concurrent use.
func (c *ciRunner) runBazelCommand(event *ciEvent, logger *log.Entry, command string) error {
	if err := c.gitSetup(logger, event); err != nil {
		logger.Errorf("Failed to setup git: %v", err)
		return err
	}
	runfile := path.Join(bazel.Workspace(), "ci", command)

	logger.Infoln("Running command", runfile)
	cmd := exec.Command(runfile)
	env := os.Environ()
	if *tmpBase {
		env = append(env, fmt.Sprintf("CACHE_DIR=%s", os.Getenv("TEST_TMPDIR")))
	} else {
		env = append(env, fmt.Sprintf("CACHE_DIR=%s", filepath.Join(os.Getenv("HOME"), "bazel-cache")))
	}
	if event.branch == "master" {
		// Normally, we find tests to run by looking for files that changed
		// between the origin/master and the new commit.
		// But if we are updating master itself, we rely on GitHub to
		// tell us what commit was in origin/master before this push.
		// This should be fast and it should be good enough for a while.
		//
		// But it's not perfect. Suppose we push a broken CL that breaks
		// test Foo - all future pushes may turn out green even though there
		// is a broken test in the repository. That's good for test isolation
		// and to prevent developers of unrelated code projects from hurting
		// each other. But we still need a way to call-out project owners
		// and let them know that one of their tests has been persistently
		// broken.
		env = append(env, fmt.Sprintf("TRAVIS_COMMIT_RANGE=%s..", event.beforeCommit))
	}

	cmd.Env = env
	return commandWithLog(cmd, logger)
}

func (c *ciRunner) handleEvent(event hookserve.Event) {

	ctx := context.Background()

	ev := &ciEvent{
		commit:       event.Commit,
		branch:       event.Branch,
		owner:        event.Owner,
		repo:         event.Repo,
		repoURL:      fmt.Sprintf("github.com/%v/%v", event.Owner, event.Repo),
		beforeCommit: event.Before,
	}
	u2, err := uuid.NewV4()
	if err != nil {
		log.Printf("bazel CI UUID failed: %s", err)
		if err := c.setCIStatus(ctx, ev, gitHubStatusFailure, "CI internal error"); err != nil {
			log.Printf("setCIStatus error: %v", err)
		}
		return
	}
	ev.build = strings.Replace(u2.String(), "-", "", -1)
	logger := eventLogger(ev)

	logger.Infof("Handling GitHub Push event: %s", event)

	if err := c.setCIStatus(ctx, ev, gitHubStatusPending, "CI test running"); err != nil {
		logger.Errorf("setCIStatus: %v", err)
	}
	status := gitHubStatusSuccess
	logger.Infof("Starting CI")

	if err = c.runBazelCommand(ev, logger, ciShell); err != nil {
		logger.Errorf("bazel CI failed: %s", err)
		status = gitHubStatusFailure
	}
	if err := c.setCIStatus(ctx, ev, status, "CI test finished"); err != nil {
		logger.Errorf("setCIStatus: %v", err)
	}
	if *runOnce {
		if err != nil {
			logger.Fatal(err)
		}
		os.Exit(0)
	}
	logger.Infof("CI Finished")

	if ev.branch == "master" {
		// Do deployments. We'll eventually move this somewhere else but it
		// makes sense to keep it simple.
		if status == gitHubStatusSuccess {
			logger.Infof("Starting CD")
			if err := c.setCDStatus(ctx, ev, gitHubStatusPending, "Deployment starting"); err != nil {
				logger.Errorf("setCDStatus failed: %v", err)
			}
			if err = c.runBazelCommand(ev, logger, cdShell); err != nil {
				logger.Errorf("bazel CD failed: %s", err)
				// TODO: Don't block stuff for now because I know it won't work.
				status = gitHubStatusSuccess
			}
			if err := c.setCDStatus(ctx, ev, status, "Deployment finished"); err != nil {
				logger.Errorf("setCDStatus failed: %v", err)
			}
			logger.Infof("CD Finished")
		}

	}

}

const (
	gitHubStatusSuccess = "success"
	gitHubStatusPending = "pending"
	gitHubStatusFailure = "failure"
)

func eventLogger(event *ciEvent) *log.Entry {
	return log.WithFields(log.Fields{
		"owner":   event.owner,
		"repoURL": event.repoURL,
		"commit":  event.commit,
		"branch":  event.branch,
		"build":   event.build,
	})
}

func main() {
	flag.Parse()

	creds := new(gitHubCredentials)
	err := envconfig.Process("secret", creds)
	if err != nil {
		log.Fatal(err.Error())
	}

	server := hookserve.NewServer()
	// TODO: Use a secret.
	// server.Secret = hookSecret
	server.Path = "/postreceive"
	http.Handle("/postreceive", server)

	http.HandleFunc("/builds/", logview.Handler)

	go func() {
		log.Println("Starting server on port", *port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
	}()

	runner := &ciRunner{
		githubClient:      newGithubClient(creds),
		githubCredentials: creds,
	}
	for event := range server.Events {
		// This runs one Bazel command at a time. If GitHub sends us another
		// event, it will be sitting in the server.Events channel until we
		// have time to run it. That's OK as long as we don't restart while
		// there are items in the queue. In the short term, we should add
		// shutdown lameducking to give us time to clear the queue. That will
		// make restarts a little annoying. Longer term, we should use a
		// distributed worker pool that keeps track of pending jobs.
		runner.handleEvent(event)
	}
}
