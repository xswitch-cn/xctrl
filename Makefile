GOPATH:=$(shell go env GOPATH)
LOCAL_YML=

ifeq ($(VERSION),)
VERSION := latest
endif

.PHONY: proto
proto:
	# protoc --proto_path=${GOPATH}/src:.  --stack_out=../ core/proto/xctrl/*.proto
	# protoc --proto_path=. --go_out=. core/proto/xctrl/*.proto --stack_out=../ core/proto/xctrl/*.proto
	protoc --proto_path=. --go_out=. core/proto/xctrl/*.proto --stack_out=. core/proto/xctrl/*.proto
java:
	protoc --proto_path=${GOPATH}/src:. --java_out=../ core/proto/xctrl/*.proto

# go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

doc-html:
	protoc --doc_out=core/proto/xctrl/doc --doc_opt=template/defaut.html,xctrl.html core/proto/xctrl/xctrl.proto

doc-md:
	protoc --doc_out=core/proto/xctrl/doc --doc_opt=markdown,xctrl.md core/proto/xctrl/xctrl.proto
