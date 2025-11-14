build_dir := "./build"

playground: lint

lint:
    go vet
    golangci-lint run

watch *cmd:
    watchexec -w . -w .justfile -e go -c -- just {{ cmd }}

test:
    go test ./...

build: test
    go build

build-release:
    go version
    CGO_ENABLED=0 go build -ldflags "-s -w" -o gitlab-resource
    upx gitlab-resource

build-docker:
    docker build . --tag docker.io/cycloid/gitlab-resource:latest

push-docker: build-docker
    docker push cycloid/gitlab-resource:latest

test-docker-check: build-docker
    cat check_delta.json | docker run -i -a STDIN -a STDERR -a STDOUT --rm -v "$(pwd):/code" -w /code cycloid/gitlab-resource:latest /opt/resource/check

test-docker-create: build-docker
    cat create.json | docker run -i -a STDIN -a STDERR -a STDOUT --rm -v "$(pwd):/code" -w /code cycloid/gitlab-resource:latest /opt/resource/out .
