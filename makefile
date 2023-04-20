GOCMD			:=$(shell which go)
NPMCMD			:=$(shell which npm)
GOBUILD			:=$(GOCMD) build

IMPORT_PATH		:=IpProxyPool/common
BUILD_TIME		:=$(shell date "+%F %T")
COMMIT_ID       :=$(shell git rev-parse HEAD)
GO_VERSION      :=$(shell $(GOCMD) version)
VERSION			:=v1.0.0
BUILD_USER		:=$(shell whoami)

FLAG			:="-w -s -X '${IMPORT_PATH}.BuildTime=${BUILD_TIME}' -X '${IMPORT_PATH}.CommitID=${COMMIT_ID}' -X '${IMPORT_PATH}.GoVersion=${GO_VERSION}'  -X '${IMPORT_PATH}.Version=${VERSION}' -X '${IMPORT_PATH}.BuildUser=${BUILD_USER}'"

BINARY_DIR=bin
BINARY_NAME:=IpProxyPool

# linux
build-linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)-$(VERSION)-linux
# linux
build:
	@CGO_ENABLED=0 $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)

#mac
build-darwin:
	CGO_ENABLED=0 GOOS=darwin $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)-$(VERSION)-darwin

# windows
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)-$(VERSION)-win.exe
# 全平台
build-all:
	make build-linux
	make build-darwin
	make build-win
	./upx $(BINARY_DIR)/$(BINARY_NAME)-*
clean:
	rm -rf bin
