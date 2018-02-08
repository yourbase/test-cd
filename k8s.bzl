# This is our testing cluster. Talk to Yves to get access to it.

# From `kubectl config current-context`
cluster = "gke_deft-cove-184100_us-central1-a_cluster-2"

# TODO: Add {BUILD_USER} to the image chroot when we support template expansion
# in the jsonnet templates.
image_chroot = "gcr.io/deft-cove-184100/yourbase"

namespace = "prod"
