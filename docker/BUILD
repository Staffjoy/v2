#
# copied/prepped https://github.com/bazelbuild/continuous-integration
#

package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_docker//container:container.bzl", "container_image")

# BIONIC (18.04) BASE
container_image(
    name = "ubuntu",
    base = "@ubuntu//:image",
    directory = "/etc/ssl/certs/",
    files = [
        "//third_party:ubuntu/ca-certificates.crt",
    ],
)

# https://hub.docker.com/_/nginx/
container_image(
    name = "nginx",
    base = "@docker-nginx//:image",
    directory = "/etc/ssl/certs/",
    files = [
        "//third_party:ubuntu/ca-certificates.crt",
    ],
)
