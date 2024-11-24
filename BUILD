load("@gazelle//:def.bzl", "gazelle")
# gazelle:prefix github.com/olarozenfeld/camel-up

load("@rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

gazelle(name = "gazelle")

go_binary(
    name = "camelup",
    embed = [":camelup_lib"],
)

go_test(
    name = "camelup_test",
    srcs = glob(["*_test.go"]),
    embed = [":camelup_lib"],
    deps = ["@org_gonum_v1_gonum//stat:stat"],
)

go_library(
    name = "camelup_lib",
    srcs = glob(["*.go"], exclude = ["*_test.go"]),
    importpath = "github.com/olarozenfeld/camelup",
    deps = ["@com_github_fatih_color//:color"],
)