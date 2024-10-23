# Copyright 2024 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@gazelle//:def.bzl", "gazelle")
load("@rules_go//go:def.bzl", "go_binary", "go_library")
load("//internal:rules.bzl", "encourage")


# gazelle:prefix github.com/cgrotz/blog-examples/custom-bazel-rule
gazelle(name = "custom-bazel-rule")

go_library(
    name = "example-lib",
    srcs = [":encouraged"],
    importpath = "github.com/cgrotz/blog-examples/custom-bazel-rule",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "example",
    embed = [":example-lib"],
    visibility = ["//visibility:public"],
)

encourage(
    name = "encouraged",
    srcs = ["main.go"],
)