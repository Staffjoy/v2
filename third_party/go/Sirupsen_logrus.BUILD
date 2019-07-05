package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "github.com/sirupsen/logrus",
    exclude_srcs = [
        "terminal_check_bsd.go",
        "terminal_check_unix.go",
        "terminal_check_solaris.go",
        "terminal_check_windows.go",
        "terminal_check_no_terminal.go",
        "terminal_check_notappengine.go",
    ],
)

external_go_package(
    name = "hooks/syslog",
    base_pkg = "github.com/sirupsen/logrus",
)

external_go_package(
    name = "formatters/logstash",
    base_pkg = "github.com/sirupsen/logrus",
)
