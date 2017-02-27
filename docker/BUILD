#
# copied/prepped https://github.com/bazelbuild/continuous-integration
#

package(default_visibility = ["//visibility:public"])

load("@bazel_tools//tools/build_defs/docker:docker.bzl", "docker_build")

# TRUSTY (14.04) BASE
docker_build(
    name = "ubuntu-trusty",
    base = "@ubuntu-trusty//:image",
    directory = "/etc/ssl/certs/",
    files = [
        "//third_party:ubuntu/ca-certificates.crt",
    ],
)

# XENIAL (16.04) BASE
docker_build(
    name = "ubuntu-xenial",
    base = "@ubuntu-xenial//:image",
    directory = "/etc/ssl/certs/",
    files = [
        "//third_party:ubuntu/ca-certificates.crt",
    ],
)

# https://hub.docker.com/_/nginx/
docker_build(
    name = "nginx",
    base = "@docker-nginx//:image",
    directory = "/etc/ssl/certs/",
    files = [
        "//third_party:ubuntu/ca-certificates.crt",
    ],
)
