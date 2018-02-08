// TODO: switch to beta.3
local k = import "external/com_github_ksonnet_lib/ksonnet.beta.2/k.libsonnet";

// Specify the import objects that we need
local container = k.extensions.v1beta1.deployment.mixin.spec.template.spec.containersType;
local deployment = k.extensions.v1beta1.deployment;
local ingress = k.extensions.v1beta1.ingress;
local service = k.core.v1.service;

local containerPort = container.portsType;
local servicePort = k.core.v1.service.mixin.spec.portsType;

local appPort = 8080;
local svcPort = 80;

local podLabels = {app: std.extVar("mc_app")};

local upperCode = std.codepoint("A") - std.codepoint("a");

local toUpper(s) =
  local arr = std.stringChars(s);
  local upArr = std.map(function(c)
    if "a" <= c && c <= "z"then
      std.char(std.codepoint(c) + upperCode)
    else
      c
  , arr);
  std.join("", upArr);

local inputSecrets = std.extVar("mc_secrets");

// Secrets are encoded in mc_secrets in the format:
// "k8s_secret_name=key_name;k8s_secret_name_2=key_name_2"
local secrets = container.env(std.map(function(c)
  local parts = std.split(c, "/");
  if std.length(parts) != 2 then
    {}
  else
    local upper = std.map(toUpper, parts);
    local secretName = std.format("SECRET_%s_%s", upper);
    {
      "name": secretName,
      "valueFrom": {
        "secretKeyRef": {
          "name" : parts[0],
          "key": parts[1],
        },
      },
    }, std.split(inputSecrets, ";")));

# All YourBase apps have the ON_KUBE environment set when
#	running on kubernetes.
local baseEnv = container.env(
  {
      "name": "ON_KUBE",
      "value": "true"
  }
);

local cont =
  container.new(std.extVar("mc_app"), std.extVar("mc_image"))
  + container.ports(containerPort.containerPort(appPort))
  + container.imagePullPolicy("Always"); # TODO: use a pinned digest instead.

local appContainer = (
  if std.length(inputSecrets) > 0 then
    cont + baseEnv + secrets
  else
    cont);

# TODO: Move back to 2+ replicas once we can read CI bot output more reliably.
local appDeployment =
  deployment.new(std.extVar("mc_app") + "-deployment", 1, appContainer, podLabels);

local appService = service
  .new(
    std.extVar("mc_svc"),
    podLabels,
    null) +
  {"spec": {
    "selector": podLabels,
    "type": "LoadBalancer",
    "ports" : [{
      "protocol": "TCP",
      "port": svcPort,
      "targetPort": appPort,
    }]
  },
};

// FIXME: we don't need an ingress for each service.
// REMOVED
local appIngress = ingress.new() + {
  "metadata": {
    "name": std.extVar("mc_ingress"),
    "annotations": {
      "kubernetes.io/ingress.global-static-ip-name": "mc-ingress"
    }
  },
  "spec": {
    "backend": {
      "serviceName": std.extVar("mc_svc"),
      "servicePort": svcPort
    }
  }
};

k.core.v1.list.new([appDeployment, appService])