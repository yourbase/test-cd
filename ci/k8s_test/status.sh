#!/bin/bash

set -eu

# TODO: aggregate status by commit or so.
kubectl -n $USER logs -l app=ci-server-app | \
	egrep -e '("Running CI command"|Completed)' --color=always | \
	GREP_COLOR='01;36' egrep -e "commit[^ ]+|$" --color=always
