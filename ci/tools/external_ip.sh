#!/bin/bash

echo "$(kubectl get svc -n $USER ci-server-svc -o jsonpath="{.status.loadBalancer.ingress[*].ip}")"

