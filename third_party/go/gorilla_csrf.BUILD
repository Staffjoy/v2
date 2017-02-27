package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "github.com/gorilla/csrf",
    deps = [
        "@go_pkg_errors//:errors",
        "@go_gorilla_context//:context",
        "@go_gorilla_securecookie//:securecookie",
    ],
    exclude_srcs = [
        "context_legacy.go",
    ],
)
