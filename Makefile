Version := $(shell date "+%Y%m%d%H%M")
GitCommit := $(shell git rev-parse HEAD)
DIR := $(shell pwd)
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"

run: build
	./build/debug/wizard-migration

build:
	go build -race -ldflags $(LDFLAGS) -o build/debug/wizard-migration main.go

build-release:
	mkdir -p build/release/ && cd build/release/
	# https://github.com/karalabe/xgo
	xgo -ldflags="$(LDFLAGS)"-targets=linux/amd64,windows/amd64,darwin/amd64 github.com/mylxsw/wizard-migration

clean:
	rm -fr build/debug/* build/release/*

.PHONY: run build build-release clean build-dashboard run-dashboard static-gen doc-gen
