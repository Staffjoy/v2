package(default_visibility = ["//visibility:public"])

load("@bazel_tools//tools/build_defs/docker:docker.bzl", "docker_build")
load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar", "pkg_deb")

docker_build(
    name = "docker",
    base = "//docker:nginx",
    debs = [
        ":app-deb",
    ],
    files = [
        ":app-tar-src",
    ],
    ports = [
        "80",
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
