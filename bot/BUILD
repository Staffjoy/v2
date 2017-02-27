package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_library")

go_library(
    name = "bot",
    srcs = [
        "bot.pb.go",
        "client.go",
        "constants.go",
    ],
    deps = [
        "//company",
        "@go_gogo_protobuf//:gogoproto",
        "@go_gogo_protobuf//:proto",
        "@go_grpc//:grpc",
        "@go_protobuf//:proto",
        "@go_protobuf//:ptypes/empty",
        "@go_x_net//:context",
    ],
)
