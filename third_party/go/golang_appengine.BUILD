package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "google.golang.org/appengine",
    deps = [
        "@go_protobuf//:proto",
        "@go_x_net//:context",
        "@go_appengine//:internal",
        "@go_appengine//:internal/datastore",
        "@go_appengine//:internal/app_identity",
        "@go_appengine//:internal/modules",
    ],
)

external_go_package(
    base_pkg = "google.golang.org/appengine",
    name = "urlfetch",
    deps = [
        "@go_protobuf//:proto",
        "@go_x_net//:context",
        "@go_appengine//:internal",
        "@go_appengine//:internal/urlfetch",
    ],
)

external_go_package(
    base_pkg = "google.golang.org/appengine",
    name = "internal/urlfetch",
    deps = [
        "@go_protobuf//:proto",
    ],
)

external_go_package(
    base_pkg = "google.golang.org/appengine",
    name = "internal/datastore",
    deps = [
        "@go_protobuf//:proto",
    ],
)

external_go_package(
    base_pkg = "google.golang.org/appengine",
    name = "internal",
    deps = [
        "@go_x_net//:context",
        "@go_appengine//:internal/base",
        "@go_protobuf//:proto",
        "@go_appengine//:internal/log",
        "@go_appengine//:internal/remote_api",
        "@go_appengine//:internal/datastore",
    ],
    exclude_srcs = [
        "api_classic.go",
        "identity_classic.go",
        "main.go",
    ],
)

external_go_package(
    base_pkg = "google.golang.org/appengine",
    name = "internal/base",
    deps = [
        "@go_x_net//:context",
        "@go_protobuf//:proto",
        "@go_appengine//:internal/log",
    ],
)

external_go_package(
    base_pkg = "google.golang.org/appengine",
    name = "internal/log",
    deps = [
        "@go_protobuf//:proto",
    ],
)

external_go_package(
    base_pkg = "google.golang.org/appengine",
    name = "internal/remote_api",
    deps = [
        "@go_protobuf//:proto",
    ],
)

external_go_package(
    base_pkg = "google.golang.org/appengine",
    name = "internal/app_identity",
    deps = [
        "@go_protobuf//:proto",
    ],
)

external_go_package(
    base_pkg = "google.golang.org/appengine",
    name = "internal/modules",
    deps = [
        "@go_protobuf//:proto",
    ],
)

external_go_package(
    base_pkg = "google.golang.org/appengine",
    name = "socket",
    deps = [
        "@go_protobuf//:proto",
        "@go_x_net//:context",
        "@go_appengine//:internal",
        "@go_appengine//:internal/socket",
    ],
    exclude_srcs = [
        "socket_classic.go",
    ]
)

external_go_package(
    base_pkg = "google.golang.org/appengine",
    name = "internal/socket",
    deps = [
        "@go_protobuf//:proto",
    ],
)
