GO := CGO_ENABLED=0 go

LDFLAGS += -X main.Version=$(shell git describe --tags --abbrev=0)
LDFLAGS += -X main.Revision=$(shell git rev-parse --short=7 HEAD)
LDFLAGS += -X "main.BuildDate=$(DATE)"
LDFLAGS += -extldflags '-static'

PACKAGES = $(shell go list ./...)

default: release

.PHONY: clean
clean:
	$(GO) clean -i ./...
	rm -rf ./bin/

.PHONY: go/fmt
go/fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: go/lint
go/lint:
	golangci-lint run

.PHONY: go/test
go/test:
	@for PKG in $(PACKAGES); do $(GO) test -cover $$PKG || exit 1; done;

.PHONY: release
release: go/test
	goreleaser build
