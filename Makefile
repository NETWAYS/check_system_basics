.PHONY: test coverage lint vet build build-all

TARGET_BASENAME := check_system_basics
TARGET_386 := $(TARGET_BASENAME)_386
TARGET_amd64 := $(TARGET_BASENAME)_amd64
TARGET_arm64 := $(TARGET_BASENAME)_arm64
TARGET_arm6 := $(TARGET_BASENAME)_arm6
TARGET_riscv64 := $(TARGET_BASENAME)_riscv64

GIT_COMMIT := $(shell git rev-list -1 HEAD)
VERSION = $(GIT_COMMIT)
GIT_LAST_TAG_COMMIT := $(shell git rev-list --tags -1)

ifeq ($(GIT_COMMIT), $(GIT_LAST_TAG_COMMIT))
	VERSION = $(shell git tag -l --contains $(GIT_COMMIT))
endif

GO_LINKERFLAGS := "-X github.com/NETWAYS/check_system_basics/cmd.version=$(VERSION)"

GO_LINKEROPTS := -ldflags $(GO_LINKERFLAGS)

GO_FILES = $(shell find . -iname '*.go')

build:
	go build $(GO_LINKEROPTS)

build-all: $(TARGET_arm6) $(TARGET_amd64) $(TARGET_386) $(TARGET_arm64) $(TARGET_riscv64)

lint:
	go fmt $(go list ./... | grep -v /vendor/)
vet:
	go vet $(go list ./... | grep -v /vendor/)
test:
	go test -v -cover ./...
coverage:
	go test -v -cover -coverprofile=coverage.out ./... &&\
	go tool cover -html=coverage.out -o coverage.html

$(TARGET_amd64): $(GO_FILES)
	CGO_ENABLED=0 go build $(GO_LINKEROPTS) -o $(TARGET_amd64)

$(TARGET_386): $(GO_FILES)
	CGO_ENABLED=0 GOARCH=386 go build $(GO_LINKEROPTS) -o $(TARGET_386)

$(TARGET_arm64): $(GO_FILES)
	CGO_ENABLED=0 GOARCH=arm64 go build $(GO_LINKEROPTS) -o $(TARGET_arm64)

$(TARGET_arm6): $(GO_FILES)
	CGO_ENABLED=0 GOARCH=arm GOARM=6 go build $(GO_LINKEROPTS) -o $(TARGET_arm6)

$(TARGET_riscv64): $(GO_FILES)
	CGO_ENABLED=0 GOARCH=riscv64 go build $(GO_LINKEROPTS) -o $(TARGET_riscv64)
