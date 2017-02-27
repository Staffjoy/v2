package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_library")

go_library(
    name = "auditlog",
    srcs = [
        "auditlog.go",
    ],
    deps = [
        "@go_logrus//:logrus",
        "@go_structs//:structs",
    ],
)
