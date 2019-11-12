PREFIX?=$(shell pwd)
BUILDDIR := ${PREFIX}/dist
BINDIR := ${PREFIX}/bin

PKG := github.com/retailnext/sentry-sendmail

VERSION := $(shell cat VERSION)
GITCOMMIT := $(shell git rev-parse HEAD)

BTVARS=\
	-X $(PKG).GitCommit=$(GITCOMMIT)\
	-X $(PKG).Version=$(VERSION)

GO_LDFLAGS=-ldflags "-s -w $(BTVARS)"

BUILDTAGS :=

.PHONY: all
all: clean build

.PHONY: build
build:
	@echo "==> $@ <=="
	@echo "Building linux binary dist..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -trimpath -tags "$(BUILDTAGS)" $(GO_LDFLAGS) -o $(BUILDDIR)/sendmail ./cmd/sendmail
	@echo "Building deb..."
	@GO111MODULE=on go run -ldflags "-X main.Version=$(VERSION)" ./cmd/build-deb
	@echo "Building binary..."
	@CGO_ENABLED=0 GO111MODULE=on go build -trimpath -tags "$(BUILDTAGS)" $(GO_LDFLAGS) -o $(BINDIR)/sendmail ./cmd/sendmail

.PHONY: clean
clean:
	@echo "==> $@ <=="
	@echo "Cleaning..."
	@$(RM) -r $(BUILDDIR)
	@$(RM) -r $(BINDIR)
