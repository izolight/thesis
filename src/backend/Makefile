.PHONY: all clean fmt vet

GO_FILES := $(shell find . -type f -name "*.go")

all: fmt vet backend

backend: pkg/signature.pb.go $(GO_FILES) go.mod
	go build -o backend ./cmd

pkg/signature.pb.go: pkg/signature.proto
	protoc --go_out=paths=source_relative:. signature.proto

clean:
	rm -f backend

fmt:
	go fmt $$(go list ./...)

vet:
	go vet ./...