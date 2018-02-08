local k = import "external/com_github_ksonnet_lib/ksonnet.beta.3/k.libsonnet";
local g = import "REPO/IMPORT.jsonnet";
local secret = k.core.v1.secret;

secret.new(std.extVar("secret_name"), g)
