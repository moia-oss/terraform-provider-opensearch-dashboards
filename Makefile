GO := CGO_ENABLED=0 go

LDFLAGS += -X main.Version=$(shell git describe --tags --abbrev=0)
LDFLAGS += -X main.Revision=$(shell git rev-parse --short=7 HEAD)
LDFLAGS += -X "main.BuildDate=$(DATE)"
LDFLAGS += -extldflags '-static'

PACKAGES = $(shell go list ./...)

OPENSEARCH_VERSION = 1.3.6
CONTAINER_RUNTIME = docker #replace this variable with 'podman' if that is the container runtime on your machine

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

start_opensearch_container:
	${CONTAINER_RUNTIME} run -d --name=opensearch -p 9200:9200 -p 9600:9600 -e "discovery.type=single-node" opensearchproject/opensearch:${OPENSEARCH_VERSION}

remove_opensearch_container:
	${CONTAINER_RUNTIME} stop opensearch
	${CONTAINER_RUNTIME} rm opensearch

smoke_test:
	cd smoketest && rm -rf .terraform && rm .terraform.lock.hcl && terraform init
	sed -i -E "s/hashes = \[//g" smoketest/.terraform.lock.hcl
	sed -i -E "s/\]//g" smoketest/.terraform.lock.hcl
	sed -i -E "s/\".*:.*\",//g" smoketest/.terraform.lock.hcl
	# for some reason on mac the sed command creates this weird file which we don't want
	rm -f smoketest/.terraform.lock.hcl-E
	set -e ;\
	TMP="hello world" ;\
	echo hi $$TMP ;\
	BASE_PATH="smoketest/.terraform/providers/registry.terraform.io/moia-oss/opensearch-dashboards"; \
    VERSION=$$(cd $$BASE_PATH && ls); \
    ENVIRONMENT=$$(cd $$BASE_PATH/$$VERSION && ls); \
    export BUILD_DEST_PATH="$$BASE_PATH/$$VERSION/$$ENVIRONMENT/terraform-provider-opensearch-dashboards_v$$VERSION"; \
    echo Building to $$BUILD_DEST_PATH; \
    CGO_ENABLED=0 go build -o $$BUILD_DEST_PATH . ; \
	# TODO: apply (with -y)
