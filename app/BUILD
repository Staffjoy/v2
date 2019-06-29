package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_docker//container:container.bzl", "container_image")
load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_deb", "pkg_tar")

container_image(
    name = "docker",
    base = "@nginx//image",
    #debs = [
    #    ":app-deb",
    #],
    files = [
        ":app-tar-src",
    ],
    ports = [
        "80",
    ],
    tars = [
        ":app-tar",
    ],
)

pkg_tar(
    name = "app-tar-src",
    files = glob([
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
    name = "app-tar",
    files = glob([
        "dist/bundle*.js",
        "dist/assets/*",
    ]) + ["dist/index.html"],
    mode = "0644",
    package_dir = "/usr/share/nginx/html",
    strip_prefix = "/app/dist",
)

pkg_deb(
    name = "app-deb",
    data = ":app-tar",
    depends = [
        "nginx",
    ],
    description = "staffjoy app static deb",
    homepage = "https://www.staffjoy.com",
    maintainer = "Staffjoy <hello@staffjoy.com>",
    package = "staffjoyapp",
    version = "0.0.1",
)
