build_dir := "./build"

playground: lint

lint:
    go vet
    golangci-lint run

watch *cmd:
    watchexec -w . -w .justfile -e go -c -- just {{ cmd }}

build:
    go build

build-release:
    go version
    CGO_ENABLED=0 go build -ldflags "-s -w" -o gitlab-resource
    upx gitlab-resource

build-docker:
    docker build . --tag cycloid/gitlab-resource:latest

test-docker: build-docker
    docker inspect docker.io/cycloid/gitlab-resource:latest | jq -r .[0].Size
