package(default_visibility = ["//visibility:public"])

load("//tools:build_rules/go.bzl", "go_library")

go_library(
    name = "crypto",
    srcs = [
        "hash.go",
        "salt.go",
        "sign.go",
        "uuid.go",
    ],
    deps = [
        "@go_jwt//:jwt-go",
        "@go_x_crypto//:bcrypt",
    ],
)
