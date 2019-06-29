package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_docker//container:container.bzl", "container_image")
load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_deb", "pkg_tar")

container_image(
    name = "docker",
    base = "@nginx//image",
    #debs = [
    #    ":deb",
    #],
    files = [
        ":tar-src",
    ],
    ports = [
        "80",
    ],
    tars = [
        ":tar",
    ],
)

pkg_tar(
    name = "tar-src",
    srcs = glob([
        "src/*",
    ]) + [
        "index.template.ejs",
        "package.json",
        "README.md",
        "webpack.config.js",
    ],
    mode = "0440",
)

pkg_tar(
    name = "tar",
    srcs = glob([
        "dist/bundle*.js",
        "dist/assets/*",
    ]) + ["dist/index.html"],
    mode = "0644",
    package_dir = "/usr/share/nginx/html",
    strip_prefix = "/myaccount/dist",
)

pkg_deb(
    name = "deb",
    data = ":tar",
    depends = [
        "nginx",
    ],
    description = "staffjoy myaccount static deb",
    homepage = "https://www.staffjoy.com",
    maintainer = "Staffjoy <hello@staffjoy.com>",
    package = "staffjoy",
    version = "0.0.1",
)
