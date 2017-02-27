package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_binary")
load("@bazel_tools//tools/build_defs/docker:docker.bzl", "docker_build")

go_binary(
    name = "envserver",
    srcs = [
        "main.go",
    ],
    deps = [
        "@go_gorilla_context//:context",
    ],
)

docker_build(
    name = "docker",
    base = "@docker_ubuntu//:trusty",
    files = [
        ":envserver",
    ],
)
