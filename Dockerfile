ARG GO_VERSION

FROM golang:${GO_VERSION}-buster

ARG GO_VERSION

WORKDIR /opt/app

COPY ./go.mod ./go.sum ./main.go ./Makefile /opt/app/
COPY ./logdna /opt/app/logdna

RUN go get golang.org/x/lint/golint \
  && go get github.com/mattn/goveralls \
  && go get

ENV COVERAGE_FILENAME=coverprofile-${GO_VERSION}

CMD golint -set_exit_status **/*.go \
  && make testcov \
  && goveralls -coverprofile=coverage/${COVERAGE_FILENAME} -service jenkins
