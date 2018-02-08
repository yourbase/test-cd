#!/bin/bash

set -eu

$(dirname $0)/logs.sh | \
	egrep -e '("Running CI command"|Completed|Failed)' --color=always | \
	GREP_COLOR='01;36' egrep -e "commit[^ ]+|$" --color=always
