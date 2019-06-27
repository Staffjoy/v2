package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_binary")
load("@bazel_tools//tools/build_defs/docker:docker.bzl", "docker_build")

go_binary(
    name = "companyserver",
    srcs = [
        "admin.go",
        "companies.go",
        "directory.go",
        "helpers.go",
        "internal.go",
        "jobs.go",
        "main.go",
        "permissions.go",
        "shifts.go",
        "teams.go",
        "timezones.go",
        "worker.go",
    ],
    deps = [
        "//account",
        "//auditlog",
        "//auth",
        "//bot",
        "//company",
        "//crypto",
        "//email",
        "//environments",
        "//healthcheck",
        "//helpers",
        "//suite",
        "@go_cloud//:trace",
        "@go_gorilla_mux//:mux",
        "@go_gorp//:gorp",
        "@go_grpc//:codes",
        "@go_grpc//:grpc",
        "@go_grpc//:metadata",
        "@go_grpc_gateway//:runtime",
        "@go_libphonenumber//:libphonenumber",
        "@go_logrus//:logrus",
        "@go_mysql//:mysql",
        "@go_protobuf//:ptypes",
        "@go_protobuf//:ptypes/empty",
        "@go_x_net//:context",
    ],
)

docker_build(
    name = "docker",
    base = "//docker:ubuntu-xenial",
    cmd = "./companyserver",
    files = [
        ":companyserver",
    ],
    visibility = ["//visibility:public"],
)
