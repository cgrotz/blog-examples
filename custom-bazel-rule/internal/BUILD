load("@gazelle//:def.bzl", "gazelle")
load("@rules_go//go:def.bzl", "go_binary", "go_library")

# gazelle:prefix github.com/cgrotz/blog-examples/custom-bazel-rule
gazelle(name = "custom-bazel-rule")

go_library(
    name = "encourager-lib",
    srcs = ["encourager.go"],
    importpath = "github.com/cgrotz/blog-examples/custom-bazel-rule/internal",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "encourager",
    embed = [":encourager-lib"],
    visibility = ["//visibility:public"],
)
