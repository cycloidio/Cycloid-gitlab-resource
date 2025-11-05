playground: lint

lint:
    go vet
    golangci-lint run

watch *cmd:
    watchexec -w . -w .justfile -e go -c -- just {{ cmd }}
