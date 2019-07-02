package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "github.com/googleapis/gax-go",
)

external_go_package(
    name = "v2",
    base_pkg = "github.com/googleapis/gax-go",
    deps = [
        "@go_grpc//:grpc",
        "@go_grpc//:codes",
        "@go_grpc//:status",
    ]
)
