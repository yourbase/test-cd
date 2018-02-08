package main

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func newGithubClient(creds *gitHubCredentials) *github.Client {
	oc := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: creds.GithubToken},
	))
	return github.NewClient(oc)
}
