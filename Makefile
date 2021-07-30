TEST?=$$(go list ./... | grep -v 'vendor')
COVERAGE_DIR?=coverage
COVERAGE_FILENAME?=coverprofile.out
COVERAGE_FILE:=$(COVERAGE_DIR)/$(COVERAGE_FILENAME)
HOSTNAME:=logdna.com
NAMESPACE:=logdna
NAME:=logdna
PROJECT:=terraform-provider-$(NAME)
VCS_REF:=$(shell git rev-parse --short HEAD)
BUILD_IMAGE_NAME?=$(PROJECT):$(VCS_REF)
GOOS:=$(shell go env GOOS)
GOARCH:=$(shell go env GOARCH)

GOLANG_LINT_VERSION=1.41.1
DOCKER_RUN=docker run --rm -i$(shell [ -t 0 ] && echo t)
BUILD_ENV=$(DOCKER_RUN) -v $(PWD):/opt/build:Z $(BUILD_FLAGS) $(BUILD_IMAGE_NAME)
LINT_CMD=$(DOCKER_RUN) -v $(PWD):/app -w /app golangci/golangci-lint:v$(GOLANG_LINT_VERSION) golangci-lint run -v
VERSION_CMD=$(DOCKER_RUN) -v $(PWD):/app -w /app -- ghcr.io/caarlos0/svu

split-bin-filename = $(word $2,$(subst _v, ,$1))

default: install-local

.env-%:
	@if [ -z '${${*}}' ]; then echo 'Environment variable $* not set' && exit 1; fi

build-image:
	@test -f gpgkey.asc || (echo GPG key missing: ./gpgkey.asc; exit 1;)
	docker build . --rm -t $(BUILD_IMAGE_NAME) 

build: build-image
	$(BUILD_ENV) goreleaser build --rm-dist --snapshot --single-target

build-local:
	@goreleaser --version >/dev/null 2>&1 || (echo "ERROR: goreleaser is required."; exit 1)
	goreleaser build --rm-dist --snapshot --single-target

install-local: BIN_DIR=./dist/$(PROJECT)_$(GOOS)_$(GOARCH)
install-local: BIN=$(shell basename ./dist/*/*)
install-local: VERSION=$(call split-bin-filename,$(shell basename ./dist/*/*),2)
install-local: TARGET_DIR=$(HOME)/.terraform.d/plugins/$(HOSTNAME)/$(NAMESPACE)/$(NAME)/$(VERSION)/${GOOS}_${GOARCH}
install-local: build-local
	mkdir -p ${TARGET_DIR}
	cp ${BIN_DIR}/${BIN} ${TARGET_DIR}/${PROJECT}

lint:
	$(LINT_CMD)

test-local: .env-SERVICE_KEY lint
	TF_ACC=1 go test -v $(TEST_ARGS) ./logdna -coverprofile $(COVERAGE_FILE)
	go tool cover -html $(COVERAGE_FILE) -o $(COVERAGE_FILE).html

test: BUILD_FLAGS:=--env SERVICE_KEY --env TF_ACC=1
test: .env-SERVICE_KEY build-image lint
	$(BUILD_ENV) go test -v $(TEST_ARGS) ./logdna

testcov: BUILD_FLAGS:=--env SERVICE_KEY --env TF_ACC=1 
testcov: .env-SERVICE_KEY build-image lint
	mkdir -p $(COVERAGE_DIR)
	$(BUILD_ENV) go test $(TEST) -v $(TEST_ARGS) -coverprofile $(COVERAGE_FILE)
	$(BUILD_ENV) go tool cover -html $(COVERAGE_FILE) -o $(COVERAGE_FILE).html

postcov: BUILD_FLAGS:=--env COVERALLS_TOKEN --env GIT_BRANCH
postcov: testcov .env-COVERALLS_TOKEN
	$(BUILD_ENV) goveralls -coverprofile=${COVERAGE_FILE} -service jenkins

test-release: build-image
	$(BUILD_ENV) goreleaser release --rm-dist --snapshot

release: BUILD_FLAGS:=--env GITHUB_TOKEN
release: build-image .env-GITHUB_TOKEN
	$(BUILD_ENV) goreleaser release --rm-dist

version-%:
	@$(VERSION_CMD) $*

.PHONY: build build-image build-local install-local lint test-local test testcov postcov test-release release version-%
