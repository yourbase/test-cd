load(
    "@io_bazel_rules_jsonnet//jsonnet:jsonnet.bzl",
    "jsonnet_to_json",
)
load("@k8s_secret//:defaults.bzl", "k8s_secret")

secrets_build_file = """
load(
    "@io_bazel_rules_jsonnet//jsonnet:jsonnet.bzl",
        "jsonnet_library",
)
jsonnet_library(
    name = "secrets",
    srcs = glob([
        "*.jsonnet",
    ]),
    visibility = ["//visibility:public"],
)
"""

# Files in this directory must be named "secret_name.jsonnet" and the content
# of the file will become the `data` section of the kubernetes secret.
def secret_repo(name, path):
  native.new_local_repository(
    name = name,
    path = path,
    build_file_content = secrets_build_file,
  )

def secret_template(name, repo):
  # Massage the jsonnet import paths. It might be better to use something like $(location).
  if repo.startswith("@"):
    repo = repo.replace("@", "external/")
    repo = repo.replace("//", "/")
  else:
    repo = repo.replace("//", "")
  native.genrule(
    name = name + "-tmpl",
    outs = [name + "_secret.jsonnet"],
    tools = ["//bazel/templates:secrets.jsonnet"],
    cmd = "sed -e 's|REPO|%s|' -e 's|IMPORT|%s|' $(location //bazel/templates:secrets.jsonnet) > $@" % (repo, name),
  )

def expand_secrets(name, secrets, repo="@secrets//"):
  # repo is a jsonnet path to the directory containing the secrets.jsonnet file.
  for secret in secrets:
    secret_template(secret, repo)

    jsonnet_to_json(
        name = "%s_secret_%s" % (name, secret.lower()),
        src = "%s_secret.jsonnet" % secret.lower(),
	deps = ["//bazel:ksonnet-lib", "%s:secrets" % repo],
        outs = [
            "%s_secret_%s.json" % (name, secret.lower()),
        ],
        vars = {
	        "secret_name": secret,
        }
    )

  for secret in secrets:
    k8s_secret(
      name = "%s_deploy_secret_%s" % (name, secret),
      template = ":%s_secret_%s.json" % (name, secret.lower()),
    )
