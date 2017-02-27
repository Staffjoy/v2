package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_library")

go_library(
    name = "environments",
    srcs = [
        "environments.go",
        "helpers.go",
        "sentry.go",
    ],
    deps = [
        "@go_intercom//:intercom-go.v2",
        "@go_logrus//:logrus",
        "@go_raven//:raven-go",
    ],
)
