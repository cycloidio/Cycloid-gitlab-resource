FROM golang:1.25 AS build

# Install required packages
RUN apt update && apt install -y gcc upx

# Create a non-root user for the runtime
RUN groupadd -f -g 1000 gitlab && useradd -u 1000 -g gitlab gitlab

RUN mkdir -p /build
WORKDIR /build

COPY . /build

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o gitlab-resource && \
      upx gitlab-resource

FROM scratch

WORKDIR /opt/resource
USER gitlab:gitlab
COPY --from=build /etc/passwd /etc/shadow /etc/group /etc/
COPY --from=build --chown=gitlab:gitlab /build/gitlab-resource /build/scripts/* /opt/resource/
