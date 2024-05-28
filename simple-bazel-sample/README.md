# Simple Bazel Example

This is a simple example project using Bazel.

You can build it using [bazelisk]
bazelisk build //...  --experimental_convenience_symlinks=ignore
vs
bazelisk build //... --experimental_convenience_symlinks=ignore --extra_execution_platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 --host_platform=@io_bazel_rules_go//go/toolchain:linux_amd64
