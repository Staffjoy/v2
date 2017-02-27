package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "golang.org/x/net",
)

external_go_package(
    base_pkg = "golang.org/x/net",
    name = "context/ctxhttp",
    deps = [
        "@go_x_net//:context",
    ],
    exclude_srcs = [
        "ctxhttp_pre17.go",
      ],
)


external_go_package(
    base_pkg = "golang.org/x/net",
    name = "context",
    exclude_srcs = [
        "pre_go17.go",
    ]
)

external_go_package(
    name = "internal/timeseries",
    base_pkg = "golang.org/x/net",
)

external_go_package(
    name = "trace",
    base_pkg = "golang.org/x/net",
    deps = [
        "@go_x_net//:internal/timeseries",
        "@go_x_net//:context",
    ]
)

external_go_package(
    name = "http2",
    base_pkg = "golang.org/x/net",
    deps = [
        "@go_x_net//:http2/hpack",
        "@go_x_net//:lex/httplex",
        "@go_x_net//:context",
    ],
    exclude_srcs = [
        "not_go17.go",
        "not_go16.go",
    ],
)

external_go_package(
    name = "http2/hpack",
    base_pkg = "golang.org/x/net",
)

external_go_package(
    name = "lex/httplex",
    base_pkg = "golang.org/x/net",
)

external_go_package(
    name = "html",
    base_pkg = "golang.org/x/net",
)

external_go_package(
    name = "html/atom",
    base_pkg = "golang.org/x/net",
)
