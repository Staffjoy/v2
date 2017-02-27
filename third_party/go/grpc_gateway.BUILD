package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "github.com/grpc-ecosystem/grpc-gateway",
    deps = [
        "@go_x_net//:context",
    ],
)

external_go_package(
    name = "runtime",
    base_pkg = "github.com/grpc-ecosystem/grpc-gateway",
    deps = [
        "@go_x_net//:context",
        "@go_grpc//:grpc",
        "@go_grpc//:grpclog",
        "@go_grpc//:codes",
        "@go_grpc//:metadata",
        "@go_protobuf//:proto",
        "@go_grpc_gateway//:runtime/internal",
        "@go_grpc_gateway//:utilities",
        "@go_protobuf//:jsonpb",
    ],
)

external_go_package(
    name = "runtime/internal",
    base_pkg = "github.com/grpc-ecosystem/grpc-gateway",
    deps = [
        "@go_protobuf//:proto",
    ],
)

external_go_package(
    name = "utilities",
    base_pkg = "github.com/grpc-ecosystem/grpc-gateway",
)

external_go_package(
    name = "protoc-gen-grpc-gateway/httprule",
    base_pkg = "github.com/grpc-ecosystem/grpc-gateway",
    deps = [
      "@go_glog//:glog",
      "@go_grpc_gateway//:utilities",
    ],
)
external_go_package(
    name = "third_party/googleapis/google/api",
    base_pkg = "github.com/grpc-ecosystem/grpc-gateway",
    deps = [
            "@go_protobuf//:proto",
      "@go_protobuf//:protoc-gen-go/descriptor",

            ],
)