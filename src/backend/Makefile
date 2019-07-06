.PHONY: all cclean

GO_FILES := $(shell find . -type f -name "*.go")

all: backend

backend: signature.pb.go $(GO_FILES) go.mod
	go build -o backend ./cmd

signature.pb.go: signature.proto
	protoc --go_out=paths=source_relative:. signature.proto

clean:
	rm -f backend
