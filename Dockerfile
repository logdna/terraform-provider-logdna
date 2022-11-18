FROM goreleaser/goreleaser:v1.12.3 as goreleaser
FROM golang:1.18-bullseye as build

COPY --from=goreleaser /usr/bin/goreleaser /usr/local/bin/goreleaser

COPY gpgkey.asc /opt/secrets/gpgkey.asc
RUN gpg --import /opt/secrets/gpgkey.asc

RUN go install github.com/mattn/goveralls@latest

WORKDIR /opt/build
