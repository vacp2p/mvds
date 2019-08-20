SHELL := /bin/bash

GO111MODULE = on

build:
	go build
.PHONY: build

test:
	go test -v ./...
.PHONY: test

protobuf:
	protoc --go_out=. ./protobuf/*.proto
.PHONY: protobuf

lint:
	golangci-lint run -v
.PHONY: lint

install-linter:
	# install linter
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.17.1
.PHONY: install-linter

vendor:
	go mod tidy
	go mod vendor
	modvendor -copy="**/*.c **/*.h" -v
.PHONY: vendor

generate:
	go generate ./...
.PHONY: generate

create-migration:
	@if [ -z "$$DIR" ]; then \
		echo 'missing DIR var'; \
		exit 1; \
	fi

	@if [ -z "$$NAME" ]; then \
		echo 'missing NAME var'; \
		exit 1; \
	fi

	mkdir -p $(DIR)
	touch $(DIR)/`date +"%s"`_$(NAME).down.sql ./$(DIR)/`date +"%s"`_$(NAME).up.sql
.PHONY: create-migration
