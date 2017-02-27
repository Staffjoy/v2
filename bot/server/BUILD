package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_binary")
load("@bazel_tools//tools/build_defs/docker:docker.bzl", "docker_build")

go_binary(
    name = "botserver",
    srcs = [
        "alerts.go",
        "dispatch.go",
        "helpers.go",
        "main.go",
        "onboarding.go",
        "quirks.go",
    ],
    deps = [
        "//account",
        "//auth",
        "//bot",
        "//company",
        "//email",
        "//environments",
        "//healthcheck",
        "//sms",
        "@go_gorilla_mux//:mux",
        "@go_grpc//:codes",
        "@go_grpc//:grpc",
        "@go_grpc//:metadata",
        "@go_logrus//:logrus",
        "@go_protobuf//:ptypes",
        "@go_protobuf//:ptypes/empty",
        "@go_x_net//:context",
    ],
)

docker_build(
    name = "docker",
    base = "//docker:ubuntu-trusty",
    cmd = "./botserver",
    files = [
        ":botserver",
    ],
    visibility = ["//visibility:public"],
)
