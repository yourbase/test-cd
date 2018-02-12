# This file contains YourBase's internal WORKSPACE to be used when developing
# YB itself.

# TODO: Refactor into smaller .bzl files that can be imported by users.
# https://github.com/yourbase/yourbase/issues/4

http_archive(
    name = "bazel_gazelle",
    sha256 = "0103991d994db55b3b5d7b06336f8ae355739635e0c2379dea16b8213ea5a223",
    url = "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.9/bazel-gazelle-0.9.tar.gz",
)

#### For Go
http_archive(
    name = "io_bazel_rules_go",
    #sha256 = "91fca9cf860a1476abdc185a5f675b641b60d3acf0596679a27b580af60bf19c",
    #url = "https://github.com/bazelbuild/rules_go/releases/download/0.7.0/rules_go-0.7.0.tar.gz",

    # contains gopkg.in fix.
    #473ed9b2bf3279bfb52cc85dd3a9d22a2f99d9a6
    #strip_prefix = "rules_go-1baa6e5acbcbf4eed0e2a93420d1e02fb373eaf9",
    #urls = ["https://github.com/bazelbuild/rules_go/archive/1baa6e5acbcbf4eed0e2a93420d1e02fb373eaf9.tar.gz"],

    # from kubernetes version
    sha256 = "0efdc3cca8ac1c29e1c837bee260dab537dfd373eb4c43c7d50246a142a7c098",
    strip_prefix = "rules_go-74d8ad8f9f59a1d9a7cf066d0980f9e394acccd7",
    urls = ["https://github.com/bazelbuild/rules_go/archive/74d8ad8f9f59a1d9a7cf066d0980f9e394acccd7.tar.gz"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")

go_rules_dependencies()

go_register_toolchains()

#### For go_image
git_repository(
    name = "io_bazel_rules_docker",
    # Does not work anymore because of changes in rules_go
    # https://github.com/bazelbuild/rules_docker/issues/262
    # tag = "v0.3.0"
    # There hasn't been a release yet, so using a commit.
    commit = "8aeab63328a82fdb8e8eb12f677a4e5ce6b183b1",
    remote = "https://github.com/bazelbuild/rules_docker.git",
)

load(
    "@io_bazel_rules_docker//go:image.bzl",
    "container_pull",
    _go_image_repos = "repositories",
    container_repositories = "repositories",
)

_go_image_repos()

#### For rules_k8s
container_repositories()

# This requires rules_docker to be fully instantiated before
# it is pulled in.
git_repository(
    name = "io_bazel_rules_k8s",
    commit = "8240d175e08b3e4c2a1f3d6038d33800fb1cd692",
    remote = "https://github.com/bazelbuild/rules_k8s.git",
)

load("@io_bazel_rules_k8s//k8s:k8s.bzl", "k8s_repositories", "k8s_defaults")

k8s_repositories()

# We can't use master because we can get stuck fetching old versions because of
# caches.
# This is self-referential so kinda hard to wrap our heads around. It might make
# sense to move the cli to a separate repo.
http_archive(
    name = "com_github_yourbase_yourbase",
    strip_prefix = "yourbase-855521282bb3c6d93eefe306c55d353a4e76faec",
    urls = ["https://github.com/yourbase/yourbase/archive/855521282bb3c6d93eefe306c55d353a4e76faec.tar.gz"],
)

#local_repository(
#    name = "com_github_yourbase_yourbase",
#    path = "/Users/yves/src/github.com/yourbase/yourbase",
#)

load("@com_github_yourbase_yourbase//:k8s.bzl", "cluster", "image_chroot", "namespace")
load("@com_github_yourbase_yourbase//bazel:k8s.bzl", "k8s_cluster")

k8s_cluster(
    # This is our testing cluster. Talk to Yves to get access.
    # From `kubectl config current-context`
    name = cluster,
    image_chroot = image_chroot,
)

#### For Google APIs

go_repository(
    name = "org_golang_x_oauth2",
    commit = "f95fa95eaa936d9d87489b15d1d18b97c1ba9c28",
    importpath = "golang.org/x/oauth2",
)

go_repository(
    name = "org_golang_google_api",
    commit = "92db9b55d2aa90e54a33cfaa8caf354afae68157",
    importpath = "google.golang.org/api",
)

go_repository(
    name = "com_github_gordonklaus_portaudio",
    commit = "e66c30a9c4ca11f93538cf8c004831bfb76f3838",
    importpath = "github.com/gordonklaus/portaudio",
)

go_repository(
    name = "com_github_google_go_genproto",
    commit = "7f0da29060c682909f650ad8ed4e515bd74fa12a",
    importpath = "github.com/google/go-genproto",
)

go_repository(
    name = "com_google_cloud_go",
    commit = "290422ce3a4bcdb0c7998c2a994e23bc950b0bdb",
    importpath = "cloud.google.com/go",
)

# BSD 3-clause
#go_repository(
go_repository(
    name = "com_github_phayes_hookserve",
    # Support the `before` field for Push events.
    # TODO: Send PR upstream once proven stable.
    commit = "075a31f8bf7fcdf6218283eb54fec92807e0bfd5",
    importpath = "github.com/phayes/hookserve",
    remote = "https://github.com/nictuku/hookserve.git",
    vcs = "git",
)

# BSD-style though LICENSE is missing.
go_repository(
    name = "com_github_bmatsuo_go_jsontree",
    commit = "8a1cc1e88d44390691f69163f6bc597aac267cc3",
    importpath = "github.com/bmatsuo/go-jsontree",
)

### For the Kubernetes client

go_repository(
    name = "io_k8s_apimachinery",
    importpath = "k8s.io/apimachinery",
    strip_prefix = "apimachinery-d2536e3e90fd52c9ebc06c63c8ee44c2d2ef16fd",
    urls = ["https://github.com/nictuku/apimachinery/archive/d2536e3e90fd52c9ebc06c63c8ee44c2d2ef16fd.tar.gz"],
)

go_repository(
    name = "io_k8s_client_go",
    importpath = "k8s.io/client-go",
    tag = "v4.0.0",
)

go_repository(
    name = "io_k8s_kube_openapi",
    commit = "39a7bf85c140f972372c2a0d1ee40adbf0c8bfe1",
    importpath = "k8s.io/kube-openapi",
)

go_repository(
    name = "com_github_google_gofuzz",
    commit = "24818f796faf91cd76ec7bddd72458fbced7a6c1",
    importpath = "github.com/google/gofuzz",
)

go_repository(
    name = "com_github_go_openapi_spec",
    commit = "a4fa9574c7aa73b2fc54e251eb9524d0482bb592",
    importpath = "github.com/go-openapi/spec",
)

go_repository(
    name = "io_k8s_api",
    commit = "218912509d74a117d05a718bb926d0948e531c20",
    importpath = "k8s.io/api",
)

go_repository(
    name = "com_github_juju_ratelimit",
    commit = "59fac5042749a5afb9af70e813da1dd5474f0167",
    importpath = "github.com/juju/ratelimit",
)

go_repository(
    name = "com_github_emicklei_go_restful",
    commit = "2dd44038f0b95ae693b266c5f87593b5d2fdd78d",
    importpath = "github.com/emicklei/go-restful",
)

go_repository(
    name = "com_github_go_openapi_loads",
    commit = "c3e1ca4c0b6160cac10aeef7e8b425cc95b9c820",
    importpath = "github.com/go-openapi/loads",
)

go_repository(
    name = "com_github_ugorji_go",
    commit = "84cb69a8af8316eed8cf4a3c9368a56977850062",
    importpath = "github.com/ugorji/go",
)

go_repository(
    name = "com_github_emicklei_go_restful_swagger12",
    commit = "7524189396c68dc4b04d53852f9edc00f816b123",
    importpath = "github.com/emicklei/go-restful-swagger12",
)

http_archive(
    name = "io_kubernetes_build",
    sha256 = "f4946917d95c54aaa98d1092757256e491f8f48fd550179134f00f902bc0b4ce",
    strip_prefix = "repo-infra-c75960142a50de16ac6225b0843b1ff3476ab0b4",
    urls = ["https://github.com/kubernetes/repo-infra/archive/c75960142a50de16ac6225b0843b1ff3476ab0b4.tar.gz"],
)

go_repository(
    name = "com_github_go_openapi_analysis",
    commit = "8ed83f2ea9f00f945516462951a288eaa68bf0d6",
    importpath = "github.com/go-openapi/analysis",
)

go_repository(
    name = "com_github_go_openapi_jsonpointer",
    commit = "779f45308c19820f1a69e9a4cd965f496e0da10f",
    importpath = "github.com/go-openapi/jsonpointer",
)

go_repository(
    name = "com_github_fsouza_go_dockerclient",
    importpath = "github.com/fsouza/go-dockerclient",
    strip_prefix = "go-dockerclient-5a3fcc6464bd83ccd04f4953488a52934be72348",
    urls = ["https://github.com/nictuku/go-dockerclient/archive/5a3fcc6464bd83ccd04f4953488a52934be72348.tar.gz"],
)

go_repository(
    name = "com_github_docker_go_units",
    commit = "0dadbb0345b35ec7ef35e228dabb8de89a65bf52",
    importpath = "github.com/docker/go-units",
)

go_repository(
    name = "com_github_docker_docker",
    build_file_generation = "on",
    build_file_name = "BUILD.bazel",

    # runtime.PluginSpec not defined
    #commit = "dfe2c023a34de3e1731e789f4344ef4d85070bc6",

    # Temporary fix:
    #strip_prefix = "docker-fcd145b091709a720c130a52196904dadafc5dda",
    #urls = ["https://github.com/nictuku/docker/archive/fcd145b091709a720c130a52196904dadafc5dda.tar.gz"],
    # Attempt fix (my PR was merged):
    commit = "1cea9d3bdb427bbdd7f14c6b11a3f6cef332bd34",
    importpath = "github.com/docker/docker",
)

# For kubectl:
# @io_k8s_kubernetes//cmd/kubectl:kubectl
go_repository(
    name = "io_k8s_kubernetes",
    importpath = "k8s.io/kubernetes",
    strip_prefix = "kubernetes-1.10.0-alpha.2",
    urls = ["https://github.com/kubernetes/kubernetes/archive/v1.10.0-alpha.2.tar.gz"],
)

go_repository(
    name = "com_github_opencontainers_runc",
    commit = "fb871d9cd069aad4c8beacd261bbe36ebbf7ffd9",
    importpath = "github.com/opencontainers/runc",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    commit = "95cd2b9c79aa5e72ab0bc69b7ccc2be15bf850f6",
    importpath = "github.com/sirupsen/logrus",
)

go_repository(
    name = "com_github_pkg_errors",
    commit = "f15c970de5b76fac0b59abb32d62c17cc7bed265",
    importpath = "github.com/pkg/errors",
)

go_repository(
    name = "com_github_Nvveen_Gotty",
    commit = "cd527374f1e5bff4938207604a14f2e38a9cf512",
    importpath = "github.com/Nvveen/Gotty",
)

go_repository(
    name = "org_golang_x_sys",
    commit = "4ff8c001ce4cc464e644b922325097228fce14d8",
    importpath = "golang.org/x/sys",
)

go_repository(
    name = "com_github_opencontainers_image_spec",
    commit = "577479e4dc273d3779f00c223c7e0dba4cd6b8b0",
    importpath = "github.com/opencontainers/image-spec",
)

go_repository(
    name = "com_github_docker_go_connections",
    commit = "3ede32e2033de7505e6500d6c868c2b9ed9f169d",
    importpath = "github.com/docker/go-connections",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "b080dc9a8c480b08e698fb1219160d598526310f",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "com_github_Microsoft_go_winio",
    commit = "78439966b38d69bf38227fbf57ac8a6fee70f69a",
    importpath = "github.com/Microsoft/go-winio",
)

go_repository(
    name = "com_github_Azure_go_ansiterm",
    commit = "d6e3b3328b783f23731bc4d058875b0371ff8109",
    importpath = "github.com/Azure/go-ansiterm",
)

go_repository(
    name = "com_github_opencontainers_go_digest",
    commit = "279bed98673dd5bef374d3b6e4b09e2af76183bf",
    importpath = "github.com/opencontainers/go-digest",
)

go_repository(
    name = "com_github_containerd_continuity",
    commit = "0cf103d319cc2d7efe085224094f466d1f8b9640",
    importpath = "github.com/containerd/continuity",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    commit = "95cd2b9c79aa5e72ab0bc69b7ccc2be15bf850f6",
    importpath = "github.com/sirupsen/logrus",
)

go_repository(
    name = "com_github_bshuster_repo_logruzio",
    build_file_generation = "on",
    build_file_name = "BUILD.bazel",
    commit = "3561e3108e18f992618aed3daa5b69f9349a74d8",
    importpath = "github.com/bshuster-repo/logruzio",
    # TODO: Send a PR upstream once this is stable.
    remote = "https://github.com/nictuku/logruzio.git",
    vcs = "git",
)

go_repository(
    name = "in_gopkg_src_d_go_git_v4",
    commit = "f9879dd043f84936a1f8acb8a53b74332a7ae135",
    importpath = "gopkg.in/src-d/go-git.v4",
)

go_repository(
    name = "in_gopkg_src_d_go_billy_v3",
    commit = "c329b7bc7b9d24905d2bc1b85bfa29f7ae266314",
    importpath = "gopkg.in/src-d/go-billy.v3",
)

go_repository(
    name = "com_github_sergi_go_diff",
    commit = "1744e2970ca51c86172c8190fadad617561ed6e7",
    importpath = "github.com/sergi/go-diff",
)

go_repository(
    name = "in_gopkg_src_d_go_git_v4",
    commit = "f9879dd043f84936a1f8acb8a53b74332a7ae135",
    importpath = "gopkg.in/src-d/go-git.v4",
)

go_repository(
    name = "com_github_jbenet_go_context",
    commit = "d14ea06fba99483203c19d92cfcd13ebe73135f4",
    importpath = "github.com/jbenet/go-context",
)

go_repository(
    name = "com_github_src_d_gcfg",
    commit = "f187355171c936ac84a82793659ebb4936bc1c23",
    importpath = "github.com/src-d/gcfg",
)

go_repository(
    name = "in_gopkg_warnings_v0",
    commit = "ec4a0fea49c7b46c2aeb0b51aac55779c607e52b",
    importpath = "gopkg.in/warnings.v0",
)

go_repository(
    name = "com_github_xanzy_ssh_agent",
    commit = "ba9c9e33906f58169366275e3450db66139a31a9",
    importpath = "github.com/xanzy/ssh-agent",
)

go_repository(
    name = "com_github_mitchellh_go_homedir",
    commit = "b8bc1bf767474819792c23f32d8286a45736f1c6",
    importpath = "github.com/mitchellh/go-homedir",
)

git_repository(
    name = "io_bazel_rules_jsonnet",
    commit = "745566196e26107042a5f787f34ded8eaeb908fb",
    remote = "https://github.com/bazelbuild/rules_jsonnet.git",
)

load("@io_bazel_rules_jsonnet//jsonnet:jsonnet.bzl", "jsonnet_repositories")

jsonnet_repositories()

ksonnet_build_file = """
filegroup(
    name = "ksonnet_files",
    srcs = glob([
        "ksonnet.beta.2/*.libsonnet",
        "ksonnet.beta.3/*.libsonnet",
    ]),
    visibility = ["//visibility:public"],
)
"""

new_git_repository(
    name = "com_github_ksonnet_lib",
    build_file_content = ksonnet_build_file,
    commit = "46d8bb9e605dc3d3977e2e2054e921ed32dd699f",
    remote = "https://github.com/ksonnet/ksonnet-lib.git",
)

container_pull(
    name = "docker_bazel",
    registry = "index.docker.io",
    repository = "insready/bazel",
    tag = "latest",  # Not ideal but no other option.
)

# For gpc_cli. This is a large repo, so use http_archive to speed up the
# download.
http_archive(
    name = "com_github_grpc_grpc",
    sha256 = "01e411f6e9b299a68cd859b301d0065f1349f586f5f9c0f6d47e8f7490ebe81d",
    strip_prefix = "grpc-0ea629c61ec70a35075e800bc3f85651f00e746f",
    urls = ["https://github.com/grpc/grpc/archive/0ea629c61ec70a35075e800bc3f85651f00e746f.tar.gz"],
)

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

# Secrets used by our testing infrastructure.
load("@com_github_yourbase_yourbase//bazel:secrets.bzl", "secret_repo")

secret_repo(
    name = "secrets",
    path = "secrets",
)

go_repository(
    name = "com_github_kelseyhightower_envconfig",
    commit = "462fda1f11d8cad3660e52737b8beefd27acfb3f",
    importpath = "github.com/kelseyhightower/envconfig",
)

go_repository(
    name = "com_github_google_go_github",
    commit = "922ceac0585d40f97d283d921f872fc50480e06e",
    importpath = "github.com/google/go-github",
)

go_repository(
    name = "com_github_google_go_querystring",
    commit = "53e6ce116135b80d037921a7fdd5138cf32d7a8a",
    importpath = "github.com/google/go-querystring",
)

go_repository(
    name = "org_golang_x_oauth2",
    commit = "b28fcf2b08a19742b43084fb40ab78ac6c3d8067",
    importpath = "golang.org/x/oauth2",
)

go_repository(
    name = "com_github_joonix_log",
    commit = "9f489441df72b5a985b0ee0423850c7999a3a09b",
    importpath = "github.com/joonix/log",
)

go_repository(
    name = "in_gopkg_olivere_elastic_v5",
    commit = "1094ee281ca61a783c9ba22bf86e7e0a8aa2d112",
    importpath = "gopkg.in/olivere/elastic.v5",
)

go_repository(
    name = "com_github_mailru_easyjson",
    commit = "32fa128f234d041f196a9f3e0fea5ac9772c08e1",
    importpath = "github.com/mailru/easyjson",
)

go_repository(
    name = "com_github_satori_go_uuid",
    commit = "36e9d2ebbde5e3f13ab2e25625fd453271d6522e",
    importpath = "github.com/satori/go.uuid",
)

go_repository(
    name = "com_github_spf13_cobra",
    commit = "f91529fc609202eededff4de2dc0ba2f662240a3",
    importpath = "github.com/spf13/cobra",
)

go_repository(
    name = "com_github_spf13_viper",
    commit = "aafc9e6bc7b7bb53ddaa75a5ef49a17d6e654be5",
    importpath = "github.com/spf13/viper",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    commit = "d670f9405373e636a5a2765eea47fac0c9bc91a4",
    importpath = "gopkg.in/yaml.v2",
)

go_repository(
    name = "com_github_spf13_jwalterweatherman",
    commit = "7c0cea34c8ece3fbeb2b27ab9b59511d360fb394",
    importpath = "github.com/spf13/jwalterweatherman",
)

go_repository(
    name = "com_github_spf13_pflag",
    commit = "4c012f6dcd9546820e378d0bdda4d8fc772cdfea",
    importpath = "github.com/spf13/pflag",
)

go_repository(
    name = "com_github_pelletier_go_toml",
    commit = "acdc4509485b587f5e675510c4f2c63e90ff68a8",
    importpath = "github.com/pelletier/go-toml",
)

go_repository(
    name = "com_github_spf13_cast",
    commit = "acbeb36b902d72a7a4c18e8f3241075e7ab763e4",
    importpath = "github.com/spf13/cast",
)

go_repository(
    name = "com_github_spf13_afero",
    commit = "bb8f1927f2a9d3ab41c9340aa034f6b803f4359c",
    importpath = "github.com/spf13/afero",
)

go_repository(
    name = "com_github_mitchellh_mapstructure",
    commit = "b4575eea38cca1123ec2dc90c26529b5c5acfcff",
    importpath = "github.com/mitchellh/mapstructure",
)

go_repository(
    name = "com_github_magiconair_properties",
    commit = "49d762b9817ba1c2e9d0c69183c2b4a8b8f1d934",
    importpath = "github.com/magiconair/properties",
)

go_repository(
    name = "com_github_hashicorp_hcl_printer",
    commit = "23c074d0eceb2b8a5bfdbb271ab780cde70f05a8",
    importpath = "github.com/hashicorp/hcl/printer",
)

go_repository(
    name = "com_github_hashicorp_hcl",
    commit = "23c074d0eceb2b8a5bfdbb271ab780cde70f05a8",
    importpath = "github.com/hashicorp/hcl",
)

go_repository(
    name = "com_github_fsnotify_fsnotify",
    commit = "c2828203cd70a50dcccfb2761f8b1f8ceef9a8e9",
    importpath = "github.com/fsnotify/fsnotify",
)

go_repository(
    name = "com_github_inconshreveable_mousetrap",
    commit = "76626ae9c91c4f2a10f34cad8ce83ea42c93bb75",
    importpath = "github.com/inconshreveable/mousetrap",
)

# For bazel-remote:
# @com_github_buchgr_bazel_remote//:bazel-remote
go_repository(
    name = "com_github_buchgr_bazel_remote",
    commit = "3ef26ca18171e0a9f1a1aabde6686036e63ced8a",
    importpath = "github.com/buchgr/bazel-remote",
)

go_repository(
    name = "com_github_djherbis_atime",
    commit = "89517e96e10b93292169a79fd4523807bdc5d5fa",
    importpath = "github.com/djherbis/atime",
)
