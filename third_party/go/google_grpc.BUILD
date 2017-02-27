package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "google.golang.org/grpc",
    deps = [
        "@go_x_net//:context",
        "@go_x_net//:trace",
        "@go_x_net//:http2",
        "@go_grpc//:grpclog",
        "@go_grpc//:naming",
        "@go_grpc//:codes",
        "@go_grpc//:transport",
        "@go_grpc//:credentials",
        "@go_grpc//:metadata",
        "@go_grpc//:internal",
        "@go_grpc//:stats",
        "@go_grpc//:tap",
        "@go_protobuf//:proto",
    ],
)

external_go_package(
    name = "grpclog",
    base_pkg = "google.golang.org/grpc",
)

external_go_package(
    name = "codes",
    base_pkg = "google.golang.org/grpc",
)

external_go_package(
    name = "credentials",
    base_pkg = "google.golang.org/grpc",
    deps = [
        "@go_x_net//:context",
    ],
    exclude_srcs = [
        "credentials_util_pre_go17.go",
    ],
)

external_go_package(
    name = "credentials/oauth",
    base_pkg = "google.golang.org/grpc",
    deps = [
        "@go_x_net//:context",
        "@go_x_oauth2//:oauth2",
        "@go_x_oauth2//:google",
        "@go_x_oauth2//:jwt",
        "@go_grpc//:credentials",
    ],
)

external_go_package(
    name = "transport",
    base_pkg = "google.golang.org/grpc",
    deps = [
        "@go_x_net//:http2",
        "@go_x_net//:context",
        "@go_x_net//:trace",
        "@go_x_net//:http2/hpack",
        "@go_grpc//:codes",
        "@go_grpc//:credentials",
        "@go_grpc//:metadata",
        "@go_grpc//:peer",
        "@go_grpc//:grpclog",
        "@go_grpc//:stats",
        "@go_grpc//:tap",
    ],
    exclude_srcs = [
        "go16.go",
        "pre_go16.go",
    ],
)

external_go_package(
    name = "naming",
    base_pkg = "google.golang.org/grpc",
)

external_go_package(
    name = "metadata",
    base_pkg = "google.golang.org/grpc",
    deps = [
      "@go_x_net//:context",
    ],
)

external_go_package(
    name = "peer",
    base_pkg = "google.golang.org/grpc",
    deps = [
      "@go_x_net//:context",
      "@go_grpc//:credentials",
    ],
)

external_go_package(
    name = "internal",
    base_pkg = "google.golang.org/grpc",
)

external_go_package(
    name = "stats",
    base_pkg = "google.golang.org/grpc",
    deps = [
      "@go_x_net//:context",
      "@go_grpc//:grpclog",
    ],
)

external_go_package(
    name = "tap",
    base_pkg = "google.golang.org/grpc",
    deps = [
      "@go_x_net//:context",
    ],
)
