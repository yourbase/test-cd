#!/bin/bash

kubectl -n $USER logs -l app=ci-server-app
