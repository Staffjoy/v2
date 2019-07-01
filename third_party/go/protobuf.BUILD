package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "github.com/golang/protobuf",
)

external_go_package(
    name = "proto",
    base_pkg = "github.com/golang/protobuf",
    exclude_srcs = [
        "pointer_unsafe.go",
    ],
)

external_go_package(
    name = "protoc-gen-go/descriptor",
    base_pkg = "github.com/golang/protobuf",
    deps = [
        "@go_protobuf//:proto",
    ],
)

external_go_package(
    name = "protoc-gen-go/generator",
    base_pkg = "github.com/golang/protobuf",
    deps = [
      "@go_protobuf//:proto",
      "@go_protobuf//:protoc-gen-go/descriptor",
      "@go_protobuf//:protoc-gen-go/plugin",
      "@go_protobuf//:protoc-gen-go/generator/internal/remap",
    ],
)

external_go_package(
    name = "protoc-gen-go/generator/internal/remap",
    base_pkg = "github.com/golang/protobuf",
)

external_go_package(
    name = "protoc-gen-go/plugin",
    base_pkg = "github.com/golang/protobuf",
    deps = [
      "@go_protobuf//:proto",
      "@go_protobuf//:protoc-gen-go/descriptor",
    ],
)

external_go_package(
    name = "jsonpb",
    base_pkg = "github.com/golang/protobuf",
    deps = [
      "@go_protobuf//:proto",
      "@go_protobuf//:ptypes/struct",
    ],
)

external_go_package(
    name = "ptypes",
    base_pkg = "github.com/golang/protobuf",
    deps = [
        "@go_protobuf//:proto",
        "@go_protobuf//:ptypes/any",
        "@go_protobuf//:ptypes/timestamp",
        "@go_protobuf//:ptypes/duration",
    ],
)

external_go_package(
    name = "ptypes/any",
    base_pkg = "github.com/golang/protobuf",
    deps = [
          "@go_protobuf//:proto",
    ],
)

external_go_package(
    name = "ptypes/struct",
    base_pkg = "github.com/golang/protobuf",
    deps = [
          "@go_protobuf//:proto",
    ],
)

external_go_package(
    name = "ptypes/empty",
    base_pkg = "github.com/golang/protobuf",
    deps = [
          "@go_protobuf//:proto",
    ],
)

external_go_package(
    name = "ptypes/duration",
    base_pkg = "github.com/golang/protobuf",
    deps = [
          "@go_protobuf//:proto",
    ],
)

external_go_package(
    name = "ptypes/wrappers",
    base_pkg = "github.com/golang/protobuf",
    deps = [
          "@go_protobuf//:proto",
    ],
)

external_go_package(
    name = "ptypes/timestamp",
    base_pkg = "github.com/golang/protobuf",
    deps = [
          "@go_protobuf//:proto",
    ],
)
