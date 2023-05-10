GOPATH:=$(shell go env GOPATH)
LOCAL_YML=

ifeq ($(VERSION),)
VERSION := latest
endif

.PHONY: proto
proto:
	# protoc --proto_path=${GOPATH}/src:.  --stack_out=../ core/proto/xctrl/*.proto
	# protoc --proto_path=. --go_out=. core/proto/xctrl/*.proto --stack_out=../ core/proto/xctrl/*.proto
	protoc --proto_path=. --go_out=. core/proto/xctrl/*.proto --xctrl_out=. core/proto/xctrl/*.proto
java:
	protoc --proto_path=${GOPATH}/src:. --java_out=../ core/proto/xctrl/*.proto

# go get -u github.com/chuanlinzhang/protoc-gen-doc/cmd/protoc-gen-doc

doc-html:
	protoc --doc_out=core/proto/xctrl/doc --doc_opt=template/default.html,xctrl.html core/proto/xctrl/xctrl.proto

doc-md:
	protoc --doc_out=core/proto/doc --doc_opt=template/default.md,base.md core/proto/base/base.proto
	protoc --doc_out=core/proto/doc --doc_opt=template/default.md,xctrl.md core/proto/xctrl/xctrl.proto
	sed -i -e 's/#map<string, string>/#map-string-string/' core/proto/doc/xctrl.md
