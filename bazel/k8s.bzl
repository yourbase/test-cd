load("@io_bazel_rules_k8s//k8s:object.bzl", "k8s_object")
load("@io_bazel_rules_k8s//k8s:k8s.bzl", "k8s_defaults")

_NAMESPACE = "prod"

def k8s_cluster(name, image_chroot):
  """k8s_cluster installs rules for operating a Kubernetes cluster.

  This macro should normally be called from a WORKSPACE file.
  """

  k8s_defaults(
      name = "k8s_deploy",
      cluster = name,
      image_chroot = image_chroot,
      kind = "deployment",
      namespace = _NAMESPACE,
  )

  [k8s_defaults(
      name = "k8s_" + kind,
      cluster = name,
      kind = kind,
      namespace = _NAMESPACE,
  ) for kind in [
      "service",
      "secret",
  ]]
