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
