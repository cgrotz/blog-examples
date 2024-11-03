load("@gazelle//:def.bzl", "gazelle")
load("@rules_go//go:def.bzl", "go_binary", "go_library")
load("//internal:rules.bzl", "encourager_rule")


# gazelle:prefix github.com/cgrotz/blog-examples/custom-bazel-rule
gazelle(name = "custom-bazel-rule")

go_library(
    name = "example-lib",
    srcs = ["main.go"],
    importpath = "github.com/cgrotz/blog-examples/custom-bazel-rule",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "example",
    embed = [":example-lib"],
    visibility = ["//visibility:public"],
)

encourager_rule(
    name = "encouraged",
    srcs = ["main.go"],
)