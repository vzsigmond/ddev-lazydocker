#ddev-generated
ARG BASE_IMAGE_BUILDER=golang
ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.21

# --- Stage 1: Build docker-proxy, launcher, and lazydocker
FROM ${BASE_IMAGE_BUILDER}:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

ARG GOARCH=amd64
ARG GOARM
ARG VERSION=dev
ARG VCS_REF=local-build

RUN apk add --no-cache git bash coreutils gcc musl-dev

# Build LazyDocker from source (no need to keep source locally)
WORKDIR /lazydocker-src
RUN git clone --depth 1 https://github.com/jesseduffield/lazydocker.git . && \
    CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} GOARM=${GOARM} \
    go build -a -mod=vendor \
    -ldflags="-s -w \
      -X main.commit=${VCS_REF} \
      -X main.version=${VERSION} \
      -X main.buildSource=Docker" \
    -o /tmp/gobuild/lazydocker

# Build your custom proxy + launcher from subdirectory
WORKDIR /docker-proxy
COPY ./docker-proxy /docker-proxy

RUN go build -o /tmp/gobuild/docker-proxy docker-proxy.go
RUN go build -o /tmp/gobuild/lazydocker-launcher main.go
RUN chmod +x /tmp/gobuild/docker-proxy /tmp/gobuild/lazydocker-launcher

# --- Final image with everything included
FROM alpine

RUN apk add --no-cache libc6-compat curl

COPY --from=builder /tmp/gobuild/lazydocker /bin/lazydocker
COPY --from=builder /tmp/gobuild/docker-proxy /bin/docker-proxy
COPY --from=builder /tmp/gobuild/lazydocker-launcher /bin/lazydocker-launcher

ENTRYPOINT ["/bin/lazydocker-launcher"]
