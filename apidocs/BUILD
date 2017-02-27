package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_library")

go_library(
    name = "apidocs",
    srcs = [
        "bindata.go",
        "serve.go",
    ],
    deps = [
        "@go_assetfs//:go-bindata-assetfs",
        "@go_logrus//:logrus",
    ],
)
