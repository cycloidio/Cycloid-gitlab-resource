build_dir := "./build"

playground: lint build
    ./gitlab-resource out . < test.json

lint:
    go vet
    golangci-lint run

watch *cmd:
    watchexec -w . -w .justfile -e go -c -- just {{ cmd }}

test: lint
    go test ./...

build: test
    go build

build-release:
    go version
    CGO_ENABLED=0 go build -ldflags "-s -w" -o gitlab-resource
    upx gitlab-resource

docker-login:
    cy get cy://org/cycloid/credentials/dockerhub-machine?key=.raw.raw.password \
        | docker login --password-stdin \
          -u "$(cy get cy://org/cycloid/credentials/dockerhub-machine?key=.raw.raw.login)"

docker-build:
    docker build . --tag docker.io/cycloid/gitlab-resource:latest

docker-push: docker-build
    docker push cycloid/gitlab-resource:latest

test-docker-check: docker-build
    cat check_delta.json | docker run -i -a STDIN -a STDERR -a STDOUT --rm -v "$(pwd):/code" -w /code cycloid/gitlab-resource:latest /opt/resource/check

test-docker-create: docker-build
    cat create.json | docker run -i -a STDIN -a STDERR -a STDOUT --rm -v "$(pwd):/code" -w /code cycloid/gitlab-resource:latest /opt/resource/out .
