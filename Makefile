SHELL := /bin/bash

GO111MODULE = on

protobuf:
	protoc --go_out=. ./protobuf/*.proto
.PHONY: protobuf