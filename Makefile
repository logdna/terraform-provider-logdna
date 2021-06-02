TEST?=$$(go list ./... | grep -v 'vendor')
COVERAGE_DIR?=coverage
COVERAGE_FILENAME?=coverprofile.out
COVERAGE_FILE=$(COVERAGE_DIR)/$(COVERAGE_FILENAME)
HOSTNAME=logdna.com
NAMESPACE=logdna
NAME=logdna
BINARY=terraform-provider-${NAME}
VERSION=1.1.0
OS_ARCH=darwin_amd64

default: install

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

# Test just the provider directory. Since there aren't multiple dirs, results will show without buffering
test:
	TF_ACC=1 go test -v $(TEST_ARGS) ./logdna

testcov:
	mkdir -p $(COVERAGE_DIR)
	TF_ACC=1 go test $(TEST) -v $(TEST_ARGS) -coverprofile $(COVERAGE_FILE)
	go tool cover -html $(COVERAGE_FILE) -o $(COVERAGE_FILE).html

.PHONY: build release install test testacc testcov
