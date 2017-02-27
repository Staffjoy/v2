package(default_visibility = ["//visibility:public"])

load("@bazel_tools//tools/build_defs/docker:docker.bzl", "docker_build")
load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar", "pkg_deb")

docker_build(
    name = "docker",
    base = "//docker:nginx",
    debs = [
        ":deb",
    ],
    files = [
        ":tar-src",
    ],
    ports = [
        "80",
    ],
)

pkg_tar(
    name = "tar-src",
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
    name = "tar",
    files = glob([
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
