package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "google.golang.org/api",
)

external_go_package(
    name = "cloudtrace/v1",
    base_pkg = "google.golang.org/api",
    deps = [
        "@go_google_api//:option",
        "@go_x_net//:context",
        "@go_x_net//:context/ctxhttp",
        "@go_google_api//:gensupport",
        "@go_google_api//:googleapi",
        "@go_google_api//:transport/http",
    ],
)

external_go_package(
    name = "gensupport",
    base_pkg = "google.golang.org/api",
    deps = [
        "@go_google_api//:googleapi",
        "@go_x_net//:context",
        "@go_x_net//:context/ctxhttp",
    ],
)

external_go_package(
    name = "googleapi",
    base_pkg = "google.golang.org/api",
    deps = [
        "@go_google_api//:googleapi/internal/uritemplates",
    ],
)

external_go_package(
    name = "support/bundler",
    base_pkg = "google.golang.org/api",
    deps = [
        "@go_x_sync//:semaphore",
    ]
)

external_go_package(
    name = "googleapi/transport",
    base_pkg = "google.golang.org/api",
)

external_go_package(
    name = "googleapi/internal/uritemplates",
    base_pkg = "google.golang.org/api",
)

external_go_package(
    name = "option",
    base_pkg = "google.golang.org/api",
    deps = [
        "@go_grpc//:grpc",
        "@go_google_api//:internal",
        "@go_x_oauth2//:oauth2",
        "@go_x_oauth2//:google",
    ],
    exclude_srcs = [
        "credentials_notgo19.go",
    ],
)

external_go_package(
    name = "internal",
    base_pkg = "google.golang.org/api",
    deps = [
        "@go_grpc//:grpc",
        "@go_grpc//:naming",
        "@go_x_oauth2//:oauth2",
        "@go_x_oauth2//:google",
    ],
)

external_go_package(
    name = "transport",
    base_pkg = "google.golang.org/api",
    deps = [
        "@go_x_net//:context",
        "@go_x_oauth2//:oauth2",
        "@go_x_oauth2//:google",
        "@go_grpc//:grpc",
        "@go_grpc//:credentials",
        "@go_grpc//:credentials/oauth",
        "@go_google_api//:transport/http",
        "@go_google_api//:internal",
        "@go_google_api//:option",
        "@go_google_api//:transport/grpc",
        "@go_appengine//:socket",
    ],
    exclude_srcs = [
        "not_go19.go",
    ],
)

external_go_package(
    name = "transport/http",
    base_pkg = "google.golang.org/api",
    deps = [
        "@go_x_oauth2//:oauth2",
        "@go_appengine//:urlfetch",
        "@go_opencensus//:plugin/ochttp",
        "@go_google_api//:option",
        "@go_google_api//:internal",
        "@go_google_api//:googleapi/transport",
        "@go_google_api//:transport/http/internal/propagation",
    ],
)

external_go_package(
    name = "transport/http/internal/propagation",
    base_pkg = "google.golang.org/api",
    deps = [
        "@go_opencensus//:trace",
        "@go_opencensus//:trace/propagation",
    ]
)

external_go_package(
    name = "transport/grpc",
    base_pkg = "google.golang.org/api",
    deps = [
        "@go_grpc//:grpc",
        "@go_grpc//:balancer/grpclb",
        "@go_grpc//:credentials",
        "@go_grpc//:credentials/google",
        "@go_grpc//:credentials/oauth",
        "@go_appengine//:appengine",
        "@go_appengine//:socket",
        "@go_x_oauth2//:oauth2",
        "@go_google_api//:option",
        "@go_google_api//:internal",
        "@go_opencensus//:plugin/ochttp",
        "@go_opencensus//:plugin/ocgrpc",
    ],
    exclude_srcs = [
        "dial_socketopt.go",
    ],
)