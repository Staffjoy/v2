package(default_visibility = ["@//visibility:public"])

load("@//third_party:go/build.bzl", "external_go_package")

external_go_package(
    base_pkg = "go.opencensus.io",
)

external_go_package(
    name = "plugin/ocgrpc",
    base_pkg = "go.opencensus.io",
    deps = [
        "@go_x_net//:context",
        "@go_grpc//:stats",
        "@go_grpc//:codes",
        "@go_grpc//:status",
        "@go_grpc//:grpclog",
        "@go_grpc//:metadata",
        "@go_opencensus//:tag",
        "@go_opencensus//:trace",
        "@go_opencensus//:stats",
        "@go_opencensus//:stats/view",
        #"@go_opencensus//:metric/metricdata",
        "@go_opencensus//:trace/propagation",
    ],
)

external_go_package(
    name = "plugin/ochttp",
    base_pkg = "go.opencensus.io",
    deps = [
        "@go_opencensus//:tag",
        "@go_opencensus//:stats",
        "@go_opencensus//:stats/view",
        "@go_opencensus//:trace",
        "@go_opencensus//:trace/propagation",
        "@go_opencensus//:plugin/ochttp/propagation/b3",
    ],
)

external_go_package(
    name = "plugin/ochttp/propagation/b3",
    base_pkg = "go.opencensus.io",
    deps = [
        "@go_opencensus//:trace",
        "@go_opencensus//:trace/propagation",
    ],
)

external_go_package(
    name = "exporterutil",
    base_pkg = "go.opencensus.io",
)

external_go_package(
    name = "trace",
    base_pkg = "go.opencensus.io",
    deps = [
        "@go_opencensus//:internal",
        "@go_opencensus//:trace/internal",
        #"@go_opencensus//:trace/tracestate",
        "@golang_lru//:simplelru",
    ],
    exclude_srcs = [
        "trace_nongo11.go",
    ],
)

#external_go_package(
#    name = "trace/tracestate",
#    base_pkg = "go.opencensus.io",
#)

external_go_package(
    name = "trace/internal",
    base_pkg = "go.opencensus.io",
)

external_go_package(
    name = "trace/propagation",
    base_pkg = "go.opencensus.io",
    deps = [
        "@go_opencensus//:trace",
    ],
)

external_go_package(
    name = "internal",
    base_pkg = "go.opencensus.io",
    deps = [
        "@go_opencensus//:exporterutil",
    ],
)

external_go_package(
    name = "internal/tagencoding",
    base_pkg = "go.opencensus.io",
)

external_go_package(
    name = "stats",
    base_pkg = "go.opencensus.io",
    deps = [
        #"@go_opencensus//:metric/metricdata",
        "@go_opencensus//:stats/internal",
        "@go_opencensus//:tag",
    ],
)
external_go_package(
    name = "stats/internal",
    base_pkg = "go.opencensus.io",
    deps = [
        "@go_opencensus//:tag",
    ],
)

external_go_package(
    name = "stats/view",
    base_pkg = "go.opencensus.io",
    deps = [
        "@go_opencensus//:tag",
        "@go_opencensus//:stats",
        "@go_opencensus//:stats/internal",
        #"@go_opencensus//:metric/metricdata",
        #"@go_opencensus//:metric/metricproducer",
        "@go_opencensus//:internal/tagencoding",
    ],
)

#external_go_package(
#    name = "resource",
#    base_pkg = "go.opencensus.io",
#)

external_go_package(
    name = "tag",
    base_pkg = "go.opencensus.io",
    exclude_srcs = [
        "profile_not19.go",
    ],
)

#external_go_package(
#    name = "metric/metricdata",
#    base_pkg = "go.opencensus.io",
#    deps = [
#        "@go_opencensus//:resource",
#    ],
#)

#external_go_package(
#    name = "metric/metricproducer",
#    base_pkg = "go.opencensus.io",
#    deps = [
#        "@go_opencensus//:metric/metricdata",
#    ],
#)