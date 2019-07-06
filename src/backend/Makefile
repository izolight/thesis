.PHONY: all clean

SOURCE_FILES := $(shell find . -type f -name "*.go")

all: backend

backend: $(SOURCE_FILES) go.mod
	go build -o backend ./cmd

clean:
	rm -f backend
