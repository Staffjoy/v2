package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_binary")
load("@bazel_tools//tools/build_defs/docker:docker.bzl", "docker_build")

go_binary(
    name = "emailserver",
    srcs = [
        "main.go",
    ],
    deps = [
        "//email",
        "//environments",
        "//healthcheck",
        "@go_gorilla_mux//:mux",
        "@go_grpc//:codes",
        "@go_grpc//:grpc",
        "@go_logrus//:logrus",
        "@go_mandrill//:mandrill",
        "@go_protobuf//:ptypes",
        "@go_protobuf//:ptypes/empty",
        "@go_x_net//:context",
    ],
)

docker_build(
    name = "docker",
    base = "//docker:ubuntu-xenial",
    cmd = "./emailserver",
    files = [
        ":emailserver",
    ],
    visibility = ["//visibility:public"],
)
