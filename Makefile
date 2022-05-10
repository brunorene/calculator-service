SHELL := /bin/bash
.SHELLFLAGS = -euo pipefail -c
.DELETE_ON_ERROR:
.ONESHELL:

GOARCH ?= amd64

os := $(shell uname)
ifeq ("$(os)", "Linux")
	GOOS = linux
else ifeq ("$(os)", "Darwin")
	GOOS = darwin
endif

gciParams := --NoInlineComments --NoPrefixComments --Section Standard --Section Default --Section "Prefix(github.com/sky-uk/$(team)/vault/$(projectName))" .

.PHONY: setup
setup:
	@echo "== setup"
	@go install github.com/daixiang0/gci@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/kyoh86/richgo@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go mod download

.PHONY: clean
clean:
	@rm -rf build

.PHONY: build-image
build-image:
	@docker build -t local/service/calculator .

.PHONY: run-service
run-service: build-image
	@docker stop calculator && docker rm calculator || true
	@docker run -d -p 8090:8090 --name calculator local/service/calculator

.PHONY: check
check: check-system-dependencies vet lint check-format test

.PHONY: check-system-dependencies
check-system-dependencies:
	@echo "== checking system dependencies"
ifeq (, $(shell which go))
	$(error "golang not found in PATH")
endif

.PHONY: vet
vet:
	@echo "== vet"
	@go vet ./...

.PHONY: lint
lint:
	@echo "== lint"
	@GOGC=45 golangci-lint --timeout 5m run

.PHONY: check-format
check-format:
	@echo "== check formatting"
	@test -z "$(shell GOGC=45 gci diff $(gciParams))"
	@test -z "$(shell GOGC=45 gofumpt -l main.go operator/)"

.PHONY: format
format:
	@echo "== format"
	@gci write $(gciParams)
	@gofumpt -w main.go operator/
	@sync

build/calculator-service: check
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o build/calculator-service


.PHONY: test
test:
	@richgo test -race -cover ./...

.PHONY: build
build: build/calculator-service

.PHONY: nft
nft: run-service
	@richgo test -v nft_test.go 