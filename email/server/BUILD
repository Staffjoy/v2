package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_binary")
load("@io_bazel_rules_docker//container:container.bzl", "container_image")

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
        "@go_mailgun//:v3",
        "@go_protobuf//:ptypes",
        "@go_protobuf//:ptypes/empty",
        "@go_x_net//:context",
    ],
)

container_image(
    name = "docker",
    base = "@ubuntu//image",
    cmd = "./emailserver",
    files = [
        ":emailserver",
    ],
    visibility = ["//visibility:public"],
)
