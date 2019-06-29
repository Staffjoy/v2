package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_binary")
load("@io_bazel_rules_docker//container:container.bzl", "container_image")

go_binary(
    name = "ical",
    srcs = [
        "cal.go",
        "main.go",
    ],
    deps = [
        "//account",
        "//auth",
        "//company",
        "//environments",
        "//errorpages",
        "//healthcheck",
        "//middlewares",
        "@go_gorilla_mux//:mux",
        "@go_grpc//:grpc",
        "@go_grpc//:metadata",
        "@go_logrus//:logrus",
        "@go_negroni//:negroni",
        "@go_x_net//:context",
    ],
)

container_image(
    name = "docker",
    base = "@ubuntu//image",
    cmd = "./ical",
    files = [
        ":ical",
    ],
    visibility = ["//visibility:public"],
)
