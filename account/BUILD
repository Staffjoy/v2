package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_library")

go_library(
    name = "account",
    srcs = [
        "account.pb.go",
        "client.go",
        "constants.go",
    ],
    deps = [
        "//auth",
        "@go_gogo_protobuf//:gogoproto",
        "@go_gogo_protobuf//:proto",
        "@go_gogo_protobuf//:types",
        "@go_grpc//:codes",
        "@go_grpc//:grpc",
        "@go_grpc//:grpclog",
        "@go_grpc//:metadata",
        "@go_grpc_gateway//:runtime",
        "@go_grpc_gateway//:third_party/googleapis/google/api",
        "@go_grpc_gateway//:utilities",
        "@go_protobuf//:proto",
        "@go_protobuf//:ptypes/empty",
        "@go_protobuf//:ptypes/timestamp",
        "@go_x_net//:context",
    ],
)
