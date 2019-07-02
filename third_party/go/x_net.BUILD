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
        "pre_go19.go",
    ],
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
    ],
    exclude_srcs = [
        "trace_go16.go",
    ],
)

external_go_package(
    name = "http2",
    base_pkg = "golang.org/x/net",
    deps = [
        "@go_x_net//:idna",
        "@go_x_net//:http2/hpack",
        "@go_x_net//:http/httpguts",
        "@go_x_net//:context",
    ],
    exclude_srcs = [
        "not_go111.go",
    ],
)

external_go_package(
    name = "http2/hpack",
    base_pkg = "golang.org/x/net",
)

external_go_package(
    name = "http/httpguts",
    base_pkg = "golang.org/x/net",
    deps = [
        "@go_x_net//:idna",
    ],
)

external_go_package(
    name = "idna",
    base_pkg = "golang.org/x/net",
    deps = [
        "@go_x_text//:unicode/bidi",
        "@go_x_text//:unicode/norm",
        "@go_x_text//:secure/bidirule",
    ],
    exclude_srcs = [
        "tables9.0.0.go",
        "tables10.0.0.go",
        "idna9.0.0.go",
    ],
)

external_go_package(
    name = "html",
    base_pkg = "golang.org/x/net",
)

external_go_package(
    name = "html/atom",
    base_pkg = "golang.org/x/net",
)
