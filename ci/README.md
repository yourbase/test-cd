# CI

# Setup

## Secrets

TODO: automatically deploy secrets before deploying the app.

Place secrets in a `github.jsonnet` file in /secrets/ on your local filesystem. That path can be changed what's in the WORKSPACE file. Example:

```
$ cat /secrets/github.jsonnet 
{
	"username": std.base64("nictuku"),
	"password" : std.base64("password here"),
	"token" : std.base64("github token here")
}
```

Deploy the secrets to your kubernetes namespace:

```
$ bazel run //ci:ci_server_deploy_secret_github.apply
# ...
(00:00:51) INFO: Running command line: bazel-bin/ci/ci_server_deploy_secret_github.apply
secret "github" created
```

## Deploy the CI

```
$ bazel run --experimental_platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //ci:ci_server_deploy.apply
```

Run a live test and check the log output. The CI build should generate a _failure_ (that's on purpose).

```
$ bash ci/tools/live_test.sh
$ bash ci/tools/logs.sh
# or
$ bash ci/tools/check_status.sh
```

## Check the kibana logs

Do the k8s port forward first (see the LoggingInfrastructure doc). Then head to the kibana dashboard at http://localhost:5601/ and search for "commit".

TODO: make this nicer.

Note: The CI logs aren't generating a logstash-friendly output. We gotta fix that.

# Autodeploys

All services are deployed to the `root` namespace whenever files that affect them are updated on git.

## Emergency procedures

Use sparingly

### How to attach to the ci server pod

i.e: how run a shell inside the ci server.

Assuming we want the CI server running on the $USER namespace:

```
POD=$(kubectl -n $USER -l app=ci-server-app get pods -o name|cut -d/ -f2)
kubectl -n $USER exec -it "$POD" -- /bin/bash
```

### How to trigger cd.sh manually inside the ci server pod

```export TRAVIS_COMMIT_RANGE=614470e59245c56863e55b53f6687f0b0253afd7..
cd /ci/github.com/yourbase/yourbase
bash -x /app/ci/ci_server_image.binary.runfiles/__main__/ci/cd.sh
```

### How to deploy manually to `root`

Normally this shouldn't be needed, until our autodeploy is robust, we may occasionally need to redeploy things.

Attach to the CI server pod and run:
```
cd /ci/github.com/yourbase/yourbase

export PATH="/app/ci/ci_server_image.binary.runfiles/io_k8s_kubernetes/cmd/kubectl/linux_amd64_pure_stripped/:$PATH"

# Find and deploy all deployables
bazel \
  --output_user_root=/root/bazel-cache/ run \
  $(bazel --output_user_root=/root/bazel-cache/ query "filter('_deploy.apply$', attr(generator_function, go_http_server, ...))")
```