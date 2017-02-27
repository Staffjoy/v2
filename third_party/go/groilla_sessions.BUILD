package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "github.com/gorilla/sessions",
    deps = [
        "@go_gorilla_context//:context",
        "@go_gorilla_securecookie//:securecookie",
    ],
)
