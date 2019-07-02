package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "cloud.google.com/go",
)

external_go_package(
    name = "trace",
    base_pkg = "cloud.google.com/go",
    deps = [
        "@go_grpc//:grpc",
        "@go_grpc//:metadata",
        "@go_x_time//:rate",
        "@go_x_net//:context",
        "@go_google_api//:cloudtrace/v1",
        "@go_google_api//:gensupport",
        "@go_google_api//:option",
        "@go_google_api//:transport",
        "@go_google_api//:transport/http",
        "@go_cloud//:internal/tracecontext",
        "@go_google_api//:support/bundler",
    ],
)

external_go_package(
    name = "compute/metadata",
    base_pkg = "cloud.google.com/go",
    deps = [
        "@go_x_net//:context",
        "@go_x_net//:context/ctxhttp",
        "@go_cloud//:internal",
    ],
)

external_go_package(
    name = "internal",
    base_pkg = "cloud.google.com/go",
    deps = [
        "@go_grpc//:status",
        "@go_google_api//:googleapi",
        "@googleapis_gax//:v2",
    ],
)

external_go_package(
    name = "internal/tracecontext",
    base_pkg = "cloud.google.com/go",
)