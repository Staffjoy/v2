package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "google.golang.org/api",
)

external_go_package(
    name = "cloudtrace/v1",
    base_pkg = "google.golang.org/api",
    deps = [
        "@go_x_net//:context",
        "@go_x_net//:context/ctxhttp",
        "@go_google_api//:gensupport",
        "@go_google_api//:googleapi",
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
    name = "googleapi/internal/uritemplates",
    base_pkg = "google.golang.org/api",
)

external_go_package(
    name = "option",
    base_pkg = "google.golang.org/api",
    deps = [
        "@go_grpc//:grpc",
        "@go_x_oauth2//:oauth2",
        "@go_google_api//:internal",
    ],
)

external_go_package(
    name = "internal",
    base_pkg = "google.golang.org/api",
    deps = [
        "@go_grpc//:grpc",
        "@go_grpc//:naming",
        "@go_x_oauth2//:oauth2",
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
        "@go_google_api//:internal",
        "@go_google_api//:option",
        "@go_appengine//:socket",
    ],
)

