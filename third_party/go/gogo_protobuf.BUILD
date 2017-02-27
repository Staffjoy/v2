package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "github.com/gogo/protobuf",
)

external_go_package(
    name = "proto",
    base_pkg = "github.com/gogo/protobuf",
    exclude_srcs = [
        "pointer_reflect.go",
    ]
)

external_go_package(
    name = "sortkeys",
    base_pkg = "github.com/gogo/protobuf",
    deps = [
    ],
)

external_go_package(
    name = "types",
    base_pkg = "github.com/gogo/protobuf",
    deps = [
        "@go_gogo_protobuf//:proto",
        "@go_gogo_protobuf//:sortkeys",
    ],
)

external_go_package(
    name = "gogoproto",
    base_pkg = "github.com/gogo/protobuf",
    deps = [
        "@go_gogo_protobuf//:proto",
        "@go_gogo_protobuf//:protoc-gen-gogo/descriptor",
    ],
)

external_go_package(
    name = "protoc-gen-gogo/descriptor",
    base_pkg = "github.com/gogo/protobuf",
        deps = [
        "@go_gogo_protobuf//:proto",

    ],
)
