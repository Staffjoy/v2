package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "github.com/mailgun/mailgun-go/v3",
    deps = [
        "@go_pkg_errors//:errors",
        "@go_mailgun//:events",
        "@easyjson//:easyjson",
    ],
    exclude_srcs = [
        "mock*.go",
    ],
)

external_go_package(
    base_pkg = "github.com/mailgun/mailgun-go/v3",
    name = "events",
    deps = [
        "@easyjson//:easyjson",
        "@easyjson//:jlexer",
        "@easyjson//:jwriter",
    ],
)