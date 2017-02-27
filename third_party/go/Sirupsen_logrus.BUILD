package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "github.com/Sirupsen/logrus",
    exclude_srcs = [
        "terminal_bsd.go",
        "terminal_solaris.go",
        "terminal_windows.go",
    ],
)

external_go_package(
    name = "hooks/syslog",
    base_pkg = "github.com/Sirupsen/logrus",
)

external_go_package(
    name = "formatters/logstash",
    base_pkg = "github.com/Sirupsen/logrus",
)
