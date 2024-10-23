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


def _encourage_impl(ctx):
    args = ctx.actions.args()

    outputs = []
    for src in ctx.files.srcs:
        file_name_parts = src.basename.rsplit(".", 1)
        file_name_without_ext = file_name_parts[0]
        output_file = ctx.actions.declare_file("{name}_encouraged.{extension}".format(name = file_name_without_ext, extension = src.extension))
        outputs.append(output_file)
        args.add("{src_path}:{out_path}".format(src_path=src.path ,out_path=output_file.path))

    inputs = depset(
        direct = ctx.files.srcs,
    )

    ctx.actions.run(
        outputs = outputs,
        inputs = inputs,
        executable = ctx.executable._encourager,
        arguments = [args],
        use_default_shell_env = True,
    )

    return [DefaultInfo(
        files = depset(outputs),
    )]

encourage = rule(
    implementation = _encourage_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = [".go"],
            doc = "Source files to add a message to",
        ),
        "_encourager": attr.label(
            default = "//internal:encourager",
            executable = True,
            cfg = "exec",
        ),
    },
    doc = "Adds an encouraging message to the input files.",
    executable = False,
)
