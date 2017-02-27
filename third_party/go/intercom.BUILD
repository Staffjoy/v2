package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "gopkg.in/intercom/intercom-go.v2",
    deps = [
        "@go_intercom//:interfaces",
    ],
)

external_go_package(
    name = "interfaces",
    base_pkg = "gopkg.in/intercom/intercom-go.v2",
    deps = [
        "@go_querystring//:query",
    ],
)


