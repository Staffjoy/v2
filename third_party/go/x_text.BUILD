package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "golang.org/x/text",
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "secure/bidirule",
    deps = [
        "@go_x_text//:transform",
        "@go_x_text//:unicode/bidi",
    ],
    exclude_srcs = [
        "bidirule9.0.0.go",
    ],
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "unicode/norm",
    deps = [
        "@go_x_text//:transform",
        "@go_x_text//:internal/ucd",
        "@go_x_text//:internal/gen",
        "@go_x_text//:internal/triegen",
    ],
    exclude_srcs = [
        "triegen.go",
        "maketables.go",
        "data9.0.0.go",
        "data10.0.0.go",
        "tables9.0.0.go",
        "tables10.0.0.go",
    ],
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "unicode/cldr",
    exclude_srcs = [
        "makexml.go",
    ],
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "internal/gen",
    deps = [
        "@go_x_text//:unicode/cldr",
    ],
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "internal/triegen",
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "internal/ucd",
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "internal/tag",
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "transform",
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "internal/colltab",
    deps = [
        "@go_x_text//:language",
        "@go_x_text//:unicode/norm",
    ],
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "internal/language",
    deps = [
        "@go_x_text//:internal/gen",
        "@go_x_text//:internal/tag",
        "@go_x_text//:unicode/cldr",
    ],
    exclude_srcs = [
        "gen.go",
        "gen_common.go",
    ],
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "internal/language/compact",
    deps = [
        "@go_x_text//:internal/language",
        "@go_x_text//:internal/gen",
        "@go_x_text//:unicode/cldr",
    ],
    exclude_srcs = [
        "gen.go",
        "gen_index.go",
        "gen_parents.go",
    ],
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "language",
    deps = [
        "@go_x_text//:internal/language",
        "@go_x_text//:internal/language/compact",
        "@go_x_text//:internal/tag",
        "@go_x_text//:internal/gen",
        "@go_x_text//:unicode/cldr",
    ],
    exclude_srcs = [
        "gen.go",
        "go1_1.go",
    ],
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "unicode/rangetable",
    deps = [
        "@go_x_text//:collate",
    ],
    exclude_srcs = [
        "gen.go",
        "tables9.0.0.go",
        "tables10.0.0.go",
    ],
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "collate",
    deps = [
        "@go_x_text//:language",
        "@go_x_text//:internal/colltab",
        "@go_x_text//:unicode/norm",
    ],
    exclude_srcs = [
        "maketables.go",
    ],
)

external_go_package(
    base_pkg = "golang.org/x/text",
    name = "unicode/bidi",
    deps = [
        "@go_x_text//:unicode/rangetable",
        "@go_x_text//:internal/ucd",
        "@go_x_text//:internal/gen",
        "@go_x_text//:internal/triegen",
    ],
    exclude_srcs = [
        "gen.go",
        "gen_ranges.go",
        "gen_trieval.go",
        "tables9.0.0.go",
        "tables10.0.0.go",
    ],
)
