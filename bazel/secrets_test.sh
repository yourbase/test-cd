#!/bin/bash

# Compare the two target files and return an error if they are different.

rundir="$(dirname $0)"

ls -la $rundir

cmp --silent "$rundir/$1" "$rundir/$2"

if [[ $? -eq 0 ]]; then
  echo '### SUCCESS: Files Are Identical! ###'
  exit 0
fi
echo '### WARNING: Files Are Different! ###'

diff -Naur "$rundir/$1" "$rundir/$2"
exit 1

