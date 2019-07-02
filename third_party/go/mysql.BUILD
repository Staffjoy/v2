package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "github.com/go-sql-driver/mysql",
    exclude_srcs = [
        "appengine.go",
        "conncheck_dummy.go",
    ],
)
