GOPATH:=$(shell go env GOPATH)
LOCAL_YML=

ifeq ($(VERSION),)
VERSION := latest
endif

.PHONY: setup
setup:
	go mod tidy
	cd xctrl/cmd/protoc-gen-xctrl && go install && cd -

.PHONY: proto
proto:
	protoc --proto_path=. --go_out=. proto/xctrl/*.proto --xctrl_out=. proto/xctrl/*.proto
	protoc --proto_path=. --go_out=. proto/cman/*.proto --xctrl_out=. proto/cman/*.proto

java:
	protoc --proto_path=${GOPATH}/src:. --java_out=../ proto/xctrl/*.proto

doc-html:
	protoc --doc_out=proto/xctrl/doc --doc_opt=template/default.html,xctrl.html proto/xctrl/xctrl.proto

doc-md:
	protoc --doc_out=proto/doc --doc_opt=template/default.md,base.md proto/base/base.proto
	protoc --doc_out=proto/doc --doc_opt=template/default.md,xctrl.md proto/xctrl/xctrl.proto
	sed -i -e 's/#map<string, string>/#map-string-string/' proto/doc/xctrl.md
