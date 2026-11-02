FROM golang:1.18.1-alpine3.15 AS builder

RUN go version

COPY . /github.com/sku4/alice-checklist/
WORKDIR /github.com/sku4/alice-checklist/

RUN go mod download
# Version is injected at build time (git short SHA in CI, "dev" locally) and
# logged on startup so a freshly deployed version is identifiable in the logs.
ARG VERSION=dev
RUN GOOS=linux go build -ldflags "-X github.com/sku4/alice-checklist/pkg/version.Version=${VERSION}" -o ./.bin/app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/sku4/alice-checklist/.bin/app .
COPY --from=0 /github.com/sku4/alice-checklist/configs/config.yml configs/config.yml
COPY --from=0 /github.com/sku4/alice-checklist/lang/json/*.json lang/json/

EXPOSE 8000

CMD ["./app"]