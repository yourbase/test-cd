# Secrets

This directory contains example secrets. It's also used to make sure the
secrets bazel rules are valid, even if the content isn't expected. It's
important that all Bazel rules in the repo are valid.

As long as the .gitignore file is configured correctly, this directory will be
ignored by git.

To copy secrets between namespaces:

```
kubectl get secrets -o json --namespace old | jq '.items[].metadata.namespace = "new"' |kubectl create -f -
```
