GOPATH:=$(shell go env GOPATH)
LOCAL_YML=

ifeq ($(VERSION),)
VERSION := latest
endif

.PHONY: proto
proto:
	# protoc --proto_path=${GOPATH}/src:.  --stack_out=../ core/proto/xctrl/*.proto
	protoc --proto_path=. --go_out=. core/proto/xctrl/*.proto --stack_out=../ core/proto/xctrl/*.proto

java:
	protoc --proto_path=${GOPATH}/src:. --java_out=../ core/proto/xctrl/*.proto
