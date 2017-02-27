package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_library")

go_library(
    name = "middlewares",
    srcs = [
        "recovery.go",
    ],
    deps = [
        "//environments",
        "//errorpages",
        "@go_logrus//:logrus",
        "@go_raven//:raven-go",
    ],
)
