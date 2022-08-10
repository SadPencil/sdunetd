SHELL := /bin/bash

BUILD_DIR = build
GO_BUILD_ARGS = -trimpath -ldflags=all="-s -w"
export CGO_ENABLED = 0

build:
	go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd

run:
	go run .

build-all:
	GOOS=linux GOARCH=arm go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-linux-arm
	GOOS=linux GOARCH=arm64 go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-linux-arm64
	GOOS=linux GOARCH=386 go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-linux-386
	GOOS=linux GOARCH=amd64 go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-linux-amd64
	GOMIPS=softfloat GOOS=linux GOARCH=mips go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-linux-mips-softfloat
	GOOS=linux GOARCH=mips64 go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-linux-mips64
	GOMIPS=softfloat GOOS=linux GOARCH=mipsle go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-linux-mipsle-softfloat
	GOOS=linux GOARCH=mips64le go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-linux-mips64le
	GOOS=linux GOARCH=arm go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-linux-arm
	GOOS=linux GOARCH=arm64 go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-linux-arm64
	GOOS=windows GOARCH=amd64 go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-windows-amd64.exe
	GOOS=windows GOARCH=386 go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-windows-386.exe
	GOOS=windows GOARCH=arm64 go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-windows-arm64.exe
	GOOS=darwin GOARCH=amd64 go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-darwin-amd64
	GOOS=freebsd GOARCH=amd64 go build $(GO_BUILD_ARGS) -o $(BUILD_DIR)/sdunetd-freebsd-amd64

upx: build-all
	upx --best --ultra-brute $(BUILD_DIR)/*

clean:
	rm -r $(BUILD_DIR)

all: build-all upx

.PHONY: all build clean upx build-all run