package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "golang.org/x/crypto",
)

external_go_package(
    base_pkg = "golang.org/x/crypto",
    name = "blowfish",
)



external_go_package(
    base_pkg = "golang.org/x/crypto",
    name = "bcrypt",
    deps = [
        "@go_x_crypto//:blowfish",
    ],
)

