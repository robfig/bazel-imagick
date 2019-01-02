load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "b7a62250a3a73277ade0ce306d22f122365b513f5402222403e507f2f997d421",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.16.3/rules_go-0.16.3.tar.gz",
)

git_repository(
    name = "bazel_gazelle",
    commit = "3c4f16ae4a2117f0908f58107d1b55e533c9a431",
    remote = "https://github.com/bazelbuild/bazel-gazelle",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

# NOTE: for local testing
# local_repository(
#     name = "com_github_robfig_imagick",
#     path = "/Users/robfig/alpha/gocode/src/github.com/robfig/imagick",
# )

go_repository(
    name = "com_github_beevik_etree",
    commit = "09746331a38f3da8c023a85a65c8c152070de725",
    importpath = "github.com/beevik/etree",
    build_file_proto_mode = "disable",
    build_file_name = "BUILD.bazel",
)

go_repository(
    name = "com_github_robfig_imagick",
    commit = "6656e0a5523d77ec1484d83c551ba881470e035d",
    importpath = "github.com/robfig/imagick",
    build_file_proto_mode = "disable",
    build_file_name = "BUILD.bazel",
)

go_repository(
    name = "com_github_stretchr_testify",
    commit = "b89eecf5ca5db6d3ba60b237ffe3df7bafb7662f",
    importpath = "github.com/stretchr/testify",
    build_file_proto_mode = "disable",
    build_file_name = "BUILD.bazel",
)
