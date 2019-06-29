package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_binary")
load("@io_bazel_rules_docker//container:container.bzl", "container_image")

go_binary(
    name = "accountserver",
    srcs = [
        "helpers.go",
        "main.go",
        "notifications.go",
        "server.go",
    ],
    deps = [
        "//account",
        "//auditlog",
        "//auth",
        "//company",
        "//crypto",
        "//email",
        "//environments",
        "//healthcheck",
        "//helpers",
        "//sms",
        "@go_cloud//:trace",
        "@go_gorilla_mux//:mux",
        "@go_gorp//:gorp",
        "@go_grpc//:codes",
        "@go_grpc//:grpc",
        "@go_grpc//:metadata",
        "@go_grpc_gateway//:runtime",
        "@go_intercom//:intercom-go.v2",
        "@go_libphonenumber//:libphonenumber",
        "@go_logrus//:logrus",
        "@go_mysql//:mysql",
        "@go_protobuf//:ptypes",
        "@go_protobuf//:ptypes/empty",
        "@go_x_net//:context",
    ],
)

container_image(
    name = "docker",
    base = "@ubuntu//image",
    cmd = "./accountserver",
    files = [
        ":accountserver",
    ],
    visibility = ["//visibility:public"],
)
