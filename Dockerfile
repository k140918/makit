# ビルドステージ
FROM golang:1-bullseye AS builder
WORKDIR /work
ARG CGO_ENABLED=0
ENV CGO_ENABLED=${CGO_ENABLED}
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o makit ./cmd/main

# 配布用イメージステージ
FROM scratch
ARG VERSION=0.1.3
LABEL org.opencontainers.image.source=https://github.com/k140918/makit \
      org.opencontainers.image.version=${VERSION} \
      org.opencontainers.image.title=makit \
      org.opencontainers.image.description="A CLI tool to create files and directories with a single command."

COPY --from=builder /work/makit /opt/makit/makit

WORKDIR /workdir
ENTRYPOINT [ "/opt/makit/makit" ]
