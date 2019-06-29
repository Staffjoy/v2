package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_binary")
load("@io_bazel_rules_docker//container:container.bzl", "container_image")

go_binary(
    name = "faraday",
    srcs = [
        "main.go",
        "mobileconfig.go",
        "robotstxt.go",
        "security_mw.go",
        "service_mw.go",
        "trace_mw.go",
    ],
    deps = [
        "//auth",
        "//environments",
        "//errorpages",
        "//faraday/services",
        "//healthcheck",
        "//middlewares",
        "@go_cloud//:trace",
        "@go_google_api//:option",
        "@go_gorilla_context//:context",
        "@go_gorilla_mux//:mux",
        "@go_logrus//:logrus",
        "@go_negroni//:negroni",
        "@go_x_net//:context",
    ],
)

container_image(
    name = "docker",
    base = "@ubuntu//image",
    cmd = "./faraday",
    files = [
        ":faraday",
    ],
    visibility = ["//visibility:public"],
)
