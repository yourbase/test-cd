load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")
load("@io_bazel_rules_docker//container:container.bzl", "container_image")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load(
    "@io_bazel_rules_jsonnet//jsonnet:jsonnet.bzl",
    "jsonnet_library",
    "jsonnet_to_json",
)
load(":secrets.bzl", "expand_secrets")
load("@k8s_deploy//:defaults.bzl", "k8s_deploy")
load("//:k8s.bzl", "image_chroot", "cluster")

def go_http_server(name, embed=None, environment_access=None, app_config=None,
  args=None, files=None, base=None, enable_uniformity_testing=True,
  env=None, secrets=[], importRepo=""):
  """Create a deployable Go server with bells and whistles.

  Arguments:
    - name: it must be a *globally unique name* in the cloud namespace.
    - files: additional files to add to the server
    - base: alternative base container image. Useful if you need a richer
      system including a shell.
    - enable_uniformity_testing: if enabled, run uniformity tests again this
      server.
    - secrets: list of string of secrets we import from a local file. Requires
      a "secrets" repository defined in the WORKSPACE. Each secret should be
      formatted as "file/key" where file is a `file.jsonnet` in the secrets
      repository.
    - env: dictionary of environment variables and their value. Unfortunate we
      can't use make variables in the values.
    - args: extra arguments to pass to the binary. BUG: for now, this only works
      with bazel run on the _image.binary. The arguments do not get added
      to container images and to pods.

  Output:
    ..._image: container image
    ..._image_binary: the raw executable go_binary
    ..._runtime_params.json: runtime config in JSON format

  TODO: use environment_access and app_config.

  """

  go_image(
    name = "%s_image" % name,
    importpath = "unused-for-now",
    embed = embed,
    visibility = ["//visibility:public"],
    args = args,
    data = files,
    base = base,
  )

  binpath = "%s_image.binary" % name

  if enable_uniformity_testing:
    go_test(
      name = "%s_uniformity_test" % name,
      embed = [ importRepo + "//testing:http_uniformity_lib"],
      data = [
          "%s_image.binary" % name,
      ],
      # Doesn't seem to work.
      args = [ "$(location :%s_image.binary)" % name ],
    )

  # Proof-of-concept of how to provide JSON configs to apps.
  native.genrule(
    name = "%s_runtime_params" % name,
    outs = ["%s_runtime_params.json" % name ],
    # TODO: Find a nicer way to dict-to-struct this. Maybe write it out
    # explicitly.
    cmd = "echo '" + struct(prod=environment_access["production"]).to_json() + "' > $@"
  )
  # ["github/password", "github/username"] => ["github"]
  secretFiles = [key for key in {s.split("/")[0]:1 for s in secrets}]
  expand_secrets(name, secretFiles)

  repo = "deft-cove-184100/" + name + "_image"
  dnsName = name.replace("_", "-")

  # TODO: add the BUILD_USER to the image path, so people don't
  # publish images on other people's directories.
  # img = image_chroot.replace("{BUILD_USER}", "$(BUILD_USER)")
  # I think we need go_http_server to be a rule and not a macro, so we can do
  # something like ctx.action.expand_template.
  img = "%s/%s_image:latest" % (image_chroot, name)

  jsonnet_to_json(
      name = name + "_kube_deployment_json",
      # importRepo is a workaround to allow us to import go_http_server
      # from external packages.
      # https://github.com/bazelbuild/rules_jsonnet/issues/36
      src = importRepo + "//bazel/templates:deployment.jsonnet",
      deps = [importRepo + "//bazel:ksonnet-lib"],
      outs = [
          name + "_kube_deployment.json",
      ],
      vars = {
        "mc_svc": dnsName + "-svc",
        "mc_app": dnsName + "-app",
        "mc_image": img,
        "mc_secrets": ";".join(secrets),
      }
  )

  k8s_deploy(
    name = name + "_deploy",
    template = name + "_kube_deployment.json",
    # The image_chroot is applied here.
    images = {
       img : name + "_image",
    }
  )
