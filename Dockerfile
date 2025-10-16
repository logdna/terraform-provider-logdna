FROM goreleaser/goreleaser:v2.12.5 AS goreleaser
FROM golang:1.18-bullseye AS build

COPY --from=goreleaser /usr/bin/goreleaser /usr/local/bin/goreleaser

COPY gpgkey.asc /opt/secrets/gpgkey.asc
RUN gpg --import /opt/secrets/gpgkey.asc

RUN go install github.com/mattn/goveralls@latest

WORKDIR /opt/build
