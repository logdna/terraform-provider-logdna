FROM goreleaser/goreleaser:v0.171.0 as goreleaser
FROM golang:1.18.8-buster as build

COPY --from=goreleaser /usr/local/bin/goreleaser /usr/local/bin/goreleaser

COPY gpgkey.asc /opt/secrets/gpgkey.asc
RUN gpg --import /opt/secrets/gpgkey.asc

RUN go get github.com/mattn/goveralls

WORKDIR /opt/build
