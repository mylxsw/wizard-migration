Version := $(shell date "+%Y%m%d%H%M")
GitCommit := $(shell git rev-parse HEAD)
DIR := $(shell pwd)
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"

run: build
	./build/debug/wizard-migration

build:
	go build -race -ldflags $(LDFLAGS) -o build/debug/wizard-migration main.go

build-release:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/release/wizard-migration-darwin main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/release/wizard-migration.exe main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/release/wizard-migration-linux main.go

clean:
	rm -fr build/debug/* build/release/*

.PHONY: run build build-release clean build-dashboard run-dashboard static-gen doc-gen
