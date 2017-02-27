package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_library")

go_library(
    name = "suite",
    srcs = [
        "data.go",
        "suite.go",
        "users.go",
    ],
    deps = [
        "//environments",
        "@go_logrus//:logrus",
    ],
)
