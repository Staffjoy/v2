package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_binary")
load("@io_bazel_rules_docker//container:container.bzl", "container_image")

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

container_image(
    name = "docker",
    base = "@ubuntu//image",
    cmd = "./botserver",
    files = [
        ":botserver",
    ],
    visibility = ["//visibility:public"],
)
