#DOCKER STUFF

new_http_archive(
  name = "docker_ubuntu",
  build_file = "BUILD.ubuntu",
  url = "https://codeload.github.com/tianon/docker-brew-ubuntu-core/zip/52c8214ecac89d45592d16ce7c14ef82ac7b0822",
  sha256 = "a7386a64ad61298ee518885b414f70f9dba86eda61aebc1bca99bd91b07dd32c",
  type = "zip",
  strip_prefix = "docker-brew-ubuntu-core-52c8214ecac89d45592d16ce7c14ef82ac7b0822"
)

# Docker base images(s)
load("//docker:docker_pull.bzl", "docker_pull")

[docker_pull(
    name = name,
    dockerfile = "//docker:Dockerfile." + name,
    tag = "local:" + name,
) for name in [
    "ubuntu-trusty",
    "ubuntu-xenial",
    "docker-nginx",
]]

# NGINX

http_file(
    name='nginx',
    url='http://nginx.org/packages/ubuntu/pool/nginx/n/nginx/nginx_1.10.1-1~trusty_amd64.deb',
    sha256='06b589dc9b3e064faa7fbc6b6c6de629a3ec59254ac8b54770fa3dc8dd1718f1',
)

# NODEJS

http_file(
    name='nodejs',
    url="https://deb.nodesource.com/node_6.x/pool/main/n/nodejs/nodejs-dbg_6.4.0-1nodesource1~trusty1_amd64.deb",
    sha256="6a481ab1ec13849ca0465f2a97255ef3291760c7dd327a424a715c015aef1543",
)

# GOLANG INIT
load("//tools/go:go_configure.bzl", "go_configure")

go_configure()

bind(
    name = "go_package_prefix",
    actual = "//:go_package_prefix",
)

# GOLANG DEPS

new_git_repository(
    name = "go_recaptcha",
    build_file = "third_party/go/dpapathanasiou_recaptcha.BUILD",
    commit = "962f1d77fed91285eb86c988e3ae8e7948e554c8",
    remote = "https://github.com/dpapathanasiou/go-recaptcha.git",
)

new_git_repository(
    name = "go_libphonenumber",
    build_file = "third_party/go/ttacon_libphonenumber.BUILD",
    commit = "5cb77679a4c77d45f2496c9ed8e60b5eec03bb47",
    remote = "https://github.com/ttacon/libphonenumber.git",
)

new_git_repository(
    name = "go_builder",
    build_file = "third_party/go/ttacon_builder.BUILD",
    commit = "7f152c3cf4714fd6318739f8f3dbcd14c2a18b39",
    remote = "https://github.com/ttacon/builder.git",
)

new_git_repository(
    name = "go_jwt",
    build_file = "third_party/go/dgrijalva_jwt.BUILD",
    commit = "d2709f9f1f31ebcda9651b03077758c1f3a0018c",
    remote = "https://github.com/dgrijalva/jwt-go.git",
)

new_git_repository(
    name = "go_gorilla_sessions",
    build_file = "third_party/go/gorilla_sessions.BUILD",
    commit = "ca9ada44574153444b00d3fd9c8559e4cc95f896",
    remote = "https://github.com/gorilla/sessions.git",
)

new_git_repository(
    name = "go_gorilla_securecookie",
    build_file = "third_party/go/gorilla_securecookie.BUILD",
    commit = "667fe4e3466a040b780561fe9b51a83a3753eefc",
    remote = "https://github.com/gorilla/securecookie.git",
)

new_git_repository(
    name = "go_gorilla_context",
    build_file = "third_party/go/gorilla_context.BUILD",
    commit = "aed02d124ae4a0e94fea4541c8effd05bf0c8296",
    remote = "https://github.com/gorilla/context.git",
)

new_git_repository(
    name = "go_gorilla_csrf",
    build_file = "third_party/go/gorilla_csrf.BUILD",
    commit = "fdae182b1882857ae6a246467084c30af79be824",
    remote = "https://github.com/gorilla/csrf.git",
)

new_git_repository(
    name = "go_gorilla_mux",
    build_file = "third_party/go/gorilla_mux.BUILD",
    commit = "0eeaf8392f5b04950925b8a69fe70f110fa7cbfc",
    remote = "https://github.com/gorilla/mux.git",
)


new_git_repository(
    name = "go_pkg_errors",
    build_file = "third_party/go/pkg_errors.BUILD",
    commit = "645ef00459ed84a119197bfb8d8205042c6df63d",
    remote = "https://github.com/pkg/errors.git",
)

new_git_repository(
    name = "go_mandrill",
    build_file = "third_party/go/keighl_mandrill.BUILD",
    commit = "6a59523fcf7d27e9230141f0e2563ba976a92b8f",
    remote = "https://github.com/keighl/mandrill.git",
)

new_git_repository(
    name = "go_logrus",
    build_file = "third_party/go/Sirupsen_logrus.BUILD",
    commit = "4b6ea7319e214d98c938f12692336f7ca9348d6b",
    remote = "https://github.com/Sirupsen/logrus.git",
)

new_git_repository(
    name = "go_testify",
    build_file = "third_party/go/stretchr_testify.BUILD",
    commit = "f390dcf405f7b83c997eac1b06768bb9f44dec18",
    remote = "https://github.com/stretchr/testify.git",
)

new_git_repository(
    name = "go_negroni",
    build_file = "third_party/go/urfave_negroni.BUILD",
    commit = "fde5e16d32adc7ad637e9cd9ad21d4ebc6192535",
    remote = "https://github.com/urfave/negroni.git",
)

new_git_repository(
    name = "go_assetfs",
    build_file = "third_party/go/elazarl_assetfs.BUILD",
    commit = "e1a2a7ec64b07d04ac9ebb072404fe8b7b60de1b",
    remote = "https://github.com/elazarl/go-bindata-assetfs.git",
)

new_git_repository(
    name = "go_raven",
    build_file = "third_party/go/getsentry_raven.BUILD",
    commit = "379f8d0a68ca237cf8893a1cdfd4f574125e2c51",
    remote = "https://github.com/getsentry/raven-go.git",
)

new_git_repository(
    name = "go_grpc",
    build_file = "third_party/go/google_grpc.BUILD",
    commit = "e59af7a0a8bf571556b40c3f871dbc4298f77693",
    remote = "https://github.com/grpc/grpc-go.git",
)

new_git_repository(
    name = "go_grpc_gateway",
    build_file = "third_party/go/grpc_gateway.BUILD",
    commit = "84398b94e188ee336f307779b57b3aa91af7063c",
    remote = "https://github.com/grpc-ecosystem/grpc-gateway.git",
)

new_git_repository(
    name = "go_gogo_protobuf",
    build_file = "third_party/go/gogo_protobuf.BUILD",
    commit = "a9cd0c35b97daf74d0ebf3514c5254814b2703b4",
    remote = "https://github.com/gogo/protobuf.git",
)

new_git_repository(
    name = "go_glog",
    build_file = "third_party/go/glog.BUILD",
    commit = "23def4e6c14b4da8ac2ed8007337bc5eb5007998",
    remote = "https://github.com/golang/glog.git",
)

new_git_repository(
    name = "go_protobuf",
    build_file = "third_party/go/protobuf.BUILD",
    commit = "1f49d83d9aa00e6ce4fc8258c71cc7786aec968a",
    remote = "https://github.com/golang/protobuf.git",
)

new_git_repository(
    name = "go_certifi",
    build_file = "third_party/go/certifi_gocertifi.BUILD",
    commit = "ec89d50f00d39494f5b3ec5cf2fe75c53467a937",
    remote = "https://github.com/certifi/gocertifi.git",
)

new_git_repository(
    name = "go_cloud",
    build_file = "third_party/go/google_cloud.BUILD",
    commit = "c033d081db673449a5095963f987693c186fcf34",
    remote = "https://github.com/GoogleCloudPlatform/google-cloud-go.git",
)

new_git_repository(
    name = "go_intercom",
    build_file = "third_party/go/intercom.BUILD",
    commit = "2f809a5bfee1c01cbef2dd76453ef0f9123e289e",
    remote = "https://github.com/intercom/intercom-go.git"
)

new_git_repository(
    name = "go_querystring",
    build_file = "third_party/go/google_querystring.BUILD",
    commit = "9235644dd9e52eeae6fa48efd539fdc351a0af53",
    remote = "https://github.com/google/go-querystring",
)

new_git_repository(
    name = "go_google_api",
    build_file = "third_party/go/google_api.BUILD",
    commit = "a69f0f19d246419bb931b0ac8f4f8d3f3e6d4feb",
    remote = "https://github.com/google/google-api-go-client.git",
)

new_git_repository(
    name = "go_appengine",
    build_file = "third_party/go/golang_appengine.BUILD",
    commit = "4f7eeb5305a4ba1966344836ba4af9996b7b4e05",
    remote = "https://github.com/golang/appengine.git",
)

new_git_repository(
    name = "go_gorp",
    build_file = "third_party/go/gorp.BUILD",
    commit = "c87af80f3cc5036b55b83d77171e156791085e2e",
    remote = "https://github.com/go-gorp/gorp.git",
)

new_git_repository(
    name = "go_blackfriday",
    build_file = "third_party/go/russross_blackfriday.BUILD",
    commit = "5f33e7b7878355cd2b7e6b8eefc48a5472c69f70",
    remote = "https://github.com/russross/blackfriday.git",
)

new_git_repository(
    name = "go_sanitized_anchor_name",
    build_file = "third_party/go/shurcool_sanitized_anchor_name.BUILD",
    commit = "1dba4b3954bc059efc3991ec364f9f9a35f597d2",
    remote = "https://github.com/shurcool/sanitized_anchor_name.git",
)

new_git_repository(
    name = "go_structs",
    build_file = "third_party/go/fatih_structs.BUILD",
    commit = "dc3312cb1a4513a366c4c9e622ad55c32df12ed3",
    remote = "https://github.com/fatih/structs.git",
)

new_git_repository(
    name = "go_mysql",
    build_file = "third_party/go/mysql.BUILD",
    commit = "0b58b37b664c21f3010e836f1b931e1d0b0b0685",
    remote = "https://github.com/go-sql-driver/mysql.git",
)


new_git_repository(
    name = "go_x_net",
    build_file = "third_party/go/x_net.BUILD",
    commit = "6250b412798208e6c90b03b7c4f226de5aa299e2",
    remote = "https://github.com/golang/net.git"
)

new_git_repository(
    name = "go_x_oauth2",
    build_file = "third_party/go/x_oauth2.BUILD",
    commit = "3c3a985cb79f52a3190fbc056984415ca6763d01",
    remote = "https://github.com/golang/oauth2.git"
)

new_git_repository(
    name = "go_x_time",
    build_file = "third_party/go/x_time.BUILD",
    commit = "a4bde12657593d5e90d0533a3e4fd95e635124cb",
    remote = "https://github.com/golang/time.git"
)

new_git_repository(
    name = "go_x_crypto",
    build_file = "third_party/go/x_crypto.BUILD",
    commit = "6ab629be5e31660579425a738ba8870beb5b7404",
    remote = "https://github.com/golang/crypto.git"
)
