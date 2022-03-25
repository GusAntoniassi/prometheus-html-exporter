ARG ARCH="amd64"
ARG OS="linux"

FROM golang:1.17.8-alpine3.15 AS builder
LABEL maintainer="GusAntoniassi"

COPY . /src/prometheus-html-exporter
WORKDIR /src/prometheus-html-exporter

RUN \
    go mod download && \
    env GOOS="$OS" GOARCH="$ARCH" CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o /bin/prometheus-html-exporter

# -----------
FROM alpine:3.15.0

WORKDIR /app

COPY --from=builder /bin/prometheus-html-exporter ./
COPY --from=builder /src/prometheus-html-exporter/LICENSE ./
COPY --from=builder /src/prometheus-html-exporter/examples/* ./

ENTRYPOINT ["/app/prometheus-html-exporter"]
