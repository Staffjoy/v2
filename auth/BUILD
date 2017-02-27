package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_library")

go_library(
    name = "auth",
    srcs = [
        "auth.go",
        "sessions.go",
    ],
    deps = [
        "//crypto",
        "//environments",
        "@go_grpc//:metadata",
    ],
)
