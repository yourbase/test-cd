#! /usr/bin/env bash

# This is YourBase's modified ci.sh that runs inside a CI container with access
# to a persistent caching volume.

# Copyright 2018 The YourBase Authors. All rights reserved.
# Copyright 2015 The Bazel Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# This script looks at the files changed in git against origin/master
# (actually a common ancestor of origin/master and the current commit) and
# queries for all deployable targets associated with those files.
#
# Running this script on a CI server should allow you to only test the targets
# that have changed since the last time your merged or fast forwarded.

# TODO: This script definitely needs a test.

set -eu

echo "Args: $@"

function kubeauth() {
 # TODO: Use $(realpath $0) to determine the runfiles dir.
 # TODO: Pass the kubectl path via $(location) in bazel.

 # Make kubectl available to other commands.
 export PATH="/app/ci/ci_server_image.binary.runfiles/io_k8s_kubernetes/cmd/kubectl/linux_amd64_pure_stripped/:$PATH"
 mkdir -p ~/.kube/
 umask 0600
 TOKEN="$(</var/run/secrets/kubernetes.io/serviceaccount/token)"

 CAFILE=/run/secrets/kubernetes.io/serviceaccount/ca.crt

 # This name must match what's in k8s.bzl. rules_k8s will issue commands
 # for that name and we need auth setup for that cluster.
 C=gke_deft-cove-184100_us-central1-a_cluster-2

 kubectl config set-cluster "${C}" --server=https://$KUBERNETES_SERVICE_HOST \
   --certificate-authority="${CAFILE}"
 kubectl config set-credentials default-admin --token="${TOKEN}"
 kubectl config set-context default-system --cluster="${C}" --user=default-admin
 kubectl config use-context default-system

# This is more or less the equivalent of:
# docker login -u _json_key -p "${SECRET_KUBE_REGISTRY_JSON_KEY}" https://gcr.io
#
# Reference:
# https://cloud.google.com/container-registry/docs/advanced-authentication
# TODO: Support different private registries.
mkdir -p ~/.docker
umask 0600

JSON_KEY_AUTH="$(echo -n "_json_key:$SECRET_KUBE_REGISTRY_JSON_KEY" |base64 -w 0)"

 cat <<EOF > ~/.docker/config.json
{
    "auths": {
        "https://gcr.io": {
                "auth": "$JSON_KEY_AUTH"
        }
    }
}
EOF
}

IS_TEST="${1:-}"
TEST_MODE="false"
if [[ "$IS_TEST" == "test" ]]; then
	TEST_MODE=true
fi

if [[ "$TEST_MODE" == "true" ]]; then
	COMMIT_RANGE="--cached"
else
	COMMIT_RANGE=${TRAVIS_COMMIT_RANGE:-$(git merge-base origin/master HEAD)".."}
fi

if [[ -z ${COMMIT_RANGE} ]]; then
        echo "Could not find commit range." 2>&1 > /dev/null
	exit 1
fi
echo "Commit range: $COMMIT_RANGE"

CACHE_DIR=${CACHE_DIR:-"$HOME/bazel-cache"}

function brun() {
    C="--output_user_root=$CACHE_DIR"

    if [[ "$TEST_MODE" == "true" ]]; then
	C=""
    fi
    bazel $C "$@"
}

# TODO: Switch to a neutral namespace.
if [[ "$TEST_MODE" != "true" ]]; then
cat <<EOF > ~/.bazelrc
startup --host_jvm_args=-Dbazel.DigestFunction=sha256
build --experimental_remote_spawn_cache
build --remote_rest_cache=http://130.211.129.0:80
EOF
fi

# Go to the root of the repo
cd "$(git rev-parse --show-toplevel)"

echo "Bazel info:"
brun info | sed -e 's/^/    /'

kubeauth

# Get a list of the current files in package form by querying Bazel.
files=()
for file in $(git diff --name-only ${COMMIT_RANGE} ); do
  # Ignore deleted files for now. See note on ci.sh.
  if [[ ! -f "${file}" ]]; then
      continue
  fi
  IFS=$'\n' read -r -a files <<< "$(brun query --noshow_progress $file)"
  brun query --noshow_progress $file
done

echo "Updated files: ${files[*]}"

# Query for the associated deployables
deployables=$(brun query \
    --noshow_progress \
    "filter('_deploy.apply$', attr(generator_function, go_http_server, rdeps(//..., set(${files[*]}))))")

if [[ ! -z $deployables ]]; then
  echo "Deploying"
  brun run \
  $deployables
fi