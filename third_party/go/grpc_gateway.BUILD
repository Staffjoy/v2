package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "github.com/grpc-ecosystem/grpc-gateway",
    deps = [
        "@go_x_net//:context",
    ],
)

external_go_package(
    name = "utilities",
    base_pkg = "github.com/grpc-ecosystem/grpc-gateway",
)

external_go_package(
    name = "internal",
    base_pkg = "github.com/grpc-ecosystem/grpc-gateway",
    deps = [
      "@go_protobuf//:proto",
      "@go_protobuf//:ptypes/any",
    ],
)

external_go_package(
    name = "runtime",
    base_pkg = "github.com/grpc-ecosystem/grpc-gateway",
    deps = [
        "@go_grpc_gateway//:utilities",
        "@go_grpc_gateway//:internal",
        "@go_x_net//:context",
        "@go_grpc//:grpc",
        "@go_grpc//:grpclog",
        "@go_grpc//:codes",
        "@go_grpc//:metadata",
        "@go_grpc//:status",
        "@go_genproto//:protobuf/field_mask",
        "@go_genproto//:googleapis/api/httpbody",
        "@go_protobuf//:proto",
        "@go_protobuf//:jsonpb",
        "@go_protobuf//:ptypes/any",
        "@go_protobuf//:ptypes/timestamp",
        "@go_protobuf//:ptypes/duration",
        "@go_protobuf//:ptypes/wrappers",
        "@go_protobuf//:protoc-gen-go/generator",
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