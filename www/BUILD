package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_binary")
load("@bazel_tools//tools/build_defs/docker:docker.bzl", "docker_build")

go_binary(
    name = "www",
    srcs = [
        "activate.go",
        "bindata.go",
        "breaktime_episodes.go",
        "confirmreset.go",
        "login.go",
        "logout.go",
        "main.go",
        "new_company.go",
        "page.go",
        "reset.go",
        "signup.go",
    ],
    deps = [
        "//account",
        "//auth",
        "//company",
        "//crypto",
        "//email",
        "//environments",
        "//errorpages",
        "//faraday/services",
        "//healthcheck",
        "//helpers",
        "//middlewares",
        "//suite",
        "@go_assetfs//:go-bindata-assetfs",
        "@go_blackfriday//:blackfriday",
        "@go_gorilla_csrf//:csrf",
        "@go_gorilla_mux//:mux",
        "@go_grpc//:codes",
        "@go_grpc//:grpc",
        "@go_grpc//:metadata",
        "@go_logrus//:logrus",
        "@go_negroni//:negroni",
        "@go_recaptcha//:go-recaptcha",
        "@go_x_net//:context",
    ],
)

docker_build(
    name = "docker",
    base = "//docker:ubuntu-trusty",
    cmd = "./www",
    files = [
        ":www",
    ],
    visibility = ["//visibility:public"],
)
