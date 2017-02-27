package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_library")

go_library(
    name = "helpers",
    srcs = [
        "helpers.go",
    ],
    deps = [
        "//account",
        "//auth",
        "@go_grpc//:metadata",
        "@go_x_net//:context",
    ],
)
