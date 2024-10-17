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

def _my_rule_impl(ctx):
    output = ctx.actions.declare_file("{name}_/main.a".format(name = ctx.label.name))
    
    args = ctx.actions.args()
    args.add("-o", output)
    args.add_all(ctx.files.srcs)

    inputs = depset(
        direct = ctx.files.srcs,
    )

    ctx.actions.run(
        outputs = [output],
        inputs = inputs,
        executable = ctx.executable._encourager,
        arguments = [args],
        use_default_shell_env = True,
    )

    return [DefaultInfo(
        files = depset([output]),
    )]

my_rule = rule(
    implementation = _my_rule_impl,
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
