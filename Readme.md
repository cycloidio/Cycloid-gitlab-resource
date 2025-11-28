# Concourse gitlab resource

This a [concourse resource](https://concourse-ci.org/resources.html) meant to interact
with Gitlab CI.

Before reading this readme, please read the links below for context:
- https://concourse-ci.org/resources.html
- https://concourse-ci.org/implementing-resource-types.html

## Quickstart

The resource is made in go, all actions like build/test/etc have been added to a
[justfile](https://github.com/casey/just).

Below is a list of the main commands

### Build

```bash
go build
```

### Build docker image

```bash
docker build . --tag <tag>
```

### Tests

```bash
go test ./...
```

### Lint

```bash
just lint
```
## Implementation

### Entrypoint

A concourse resource must contain three scripts:
```
/opt/resource/check
/opt/resource/in
/opt/resource/out
```

To integrate in the docker image, the binary will be copied in `/opt/resource/gitlab-resource`,
to integrate with concourse, we created three scripts [here](./scripts/) named
`check`, `in`, `out` that contains a shebang like `#! /opt/resource/gitlab-resource`.

This will call gitlab resource like that:
```
gitlab-resource /opt/resource/<check,in,out> <args>
```

The [main](./main.go) contains the logic to handle that an trigger the action.

This trick allow us to make a very small docker image with only the binary and
some utility files, making the attack surface minimal.

### Features

Each feature will have its own sets of params and settings, they must be added
in the [models folder](./models/).

Then you will have to create a handler that implements the [Feature interface](./features/interface.go)
and add its creation in the NewFeatureHandler function.

For the implementation details, look up how the [deployments](./features/deployments/)
or [environments](./features/environments/) are made.

## Tests

There are no end to end tests for now (as it would require way more time to implement
right now).

Some units tests have been itegrated using the defaults go test system.
