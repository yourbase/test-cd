#/bin/bash

set -eu


target_host="$($(dirname $0)/external_ip.sh)"
# TODO: This restarts the CI/CD itself with an old version. Replace with a
# contained repo.
payload="$(dirname $0)/../testdata/cd_payload.json"

curl -f \
	-H "content-type: application/json" \
	-H "User-Agent: GitHub-Hookshot/f2b2366" \
	-H "X-GitHub-Delivery: 418f1f00-fcf7-11e7-9fba-db41916318a8" \
	-H "X-Hub-Signature: sha1=44c8e2f4ee2cb1a7667d82f741c8bbcdc4f8e6a8" \
	-H "X-GitHub-Event: push" --data "@${payload}" http://${target_host}/postreceive
