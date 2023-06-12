GOPATH:=$(shell go env GOPATH)
LOCAL_YML=

ifeq ($(VERSION),)
VERSION := latest
endif

.PHONY: setup
setup:
	go mod tidy

.PHONY: test
test:
	go test ./... -cover

