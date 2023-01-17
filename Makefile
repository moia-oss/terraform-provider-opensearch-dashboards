GO := CGO_ENABLED=0 go

LDFLAGS += -X main.Version=$(shell git describe --tags --abbrev=0)
LDFLAGS += -X main.Revision=$(shell git rev-parse --short=7 HEAD)
LDFLAGS += -X "main.BuildDate=$(DATE)"
LDFLAGS += -extldflags '-static'

PACKAGES = $(shell go list ./...)

OPENSEARCH_VERSION = 1.3.6
CONTAINER_RUNTIME = podman #set this variable to 'docker' if that is the container runtime on your machine

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

start_opensearch_container: create_opensearch_container wait_for_opensearch

create_opensearch_container:
	${CONTAINER_RUNTIME} network create opensearch_network
	${CONTAINER_RUNTIME} run -d --name=opensearch -p 9200:9200 -p 9600:9600 -e "discovery.type=single-node" -e "plugins.security.disabled=true" opensearchproject/opensearch:${OPENSEARCH_VERSION}
	${CONTAINER_RUNTIME} run -d --name=opensearch_dashboards -p 5601:5601 -e "OPENSEARCH_HOSTS=[\"http://opensearch:9200\"]" -e "DISABLE_SECURITY_DASHBOARDS_PLUGIN=true" opensearchproject/opensearch-dashboards:${OPENSEARCH_VERSION}
	${CONTAINER_RUNTIME} network connect opensearch_network opensearch
	${CONTAINER_RUNTIME} network connect opensearch_network opensearch_dashboards

remove_opensearch_container:
	${CONTAINER_RUNTIME} stop opensearch
	${CONTAINER_RUNTIME} rm opensearch
	${CONTAINER_RUNTIME} stop opensearch_dashboards
	${CONTAINER_RUNTIME} rm opensearch_dashboards
	${CONTAINER_RUNTIME} network rm opensearch_network

restart_opensearch_container: remove_opensearch_container start_opensearch_container

# init_smoke_test always re-initializes the .terraform folder even if it is already present
init_smoke_test:
	cd smoketest && rm -rf .terraform && rm .terraform.lock.hcl && terraform init
smoke_test: init_smoke_test smoke_test_fast

# this only creates the .terraform folder if it is not already present
smoketest/.terraform:
	cd smoketest && rm -rf .terraform && rm .terraform.lock.hcl && terraform init

# smoke_test_fast runs faster than smoke_test because we skip the initial cleanup / terraform init step.
# But it can cause errors in some cases (for example when a new version was released between runs) because if
# ls returns multiple results the error message doesn't clearly indicate that you need to delete and re-initialize the .terraform-folder.
# so to avoid errors for the casual user (and since the pipeline needs to initialize on every run anyways) the usage of smoke_test
# is recommended in most cases
# But when you are currently developing a new feature, this make-command may save you some annoying waits ;)
smoke_test_fast: smoketest/.terraform
	sed -i -E "s/hashes = \[//g" smoketest/.terraform.lock.hcl
	sed -i -E "s/\]//g" smoketest/.terraform.lock.hcl
	sed -i -E "s/\".*:.*\",//g" smoketest/.terraform.lock.hcl
	# for some reason on mac the sed command creates this weird file which we don't want
	rm -f smoketest/.terraform.lock.hcl-E
	set -e ;\
	BASE_PATH="smoketest/.terraform/providers/registry.terraform.io/moia-oss/opensearch-dashboards"; \
    VERSION=$$(cd $$BASE_PATH && ls); \
    ENVIRONMENT=$$(cd $$BASE_PATH/$$VERSION && ls); \
    export BUILD_DEST_PATH="$$BASE_PATH/$$VERSION/$$ENVIRONMENT/terraform-provider-opensearch-dashboards_v$$VERSION"; \
    echo Building to $$BUILD_DEST_PATH; \
    CGO_ENABLED=0 go build -o $$BUILD_DEST_PATH . ;
	cd smoketest && terraform apply -auto-approve

wait_for_opensearch:
	set -e ;\
	max_attempts=120; \
	attempt_counter=0; \
	echo "------------------------------------"; \
	until $$(curl --output /dev/null --silent --head --fail http://localhost:9200); do \
	    if [ $${attempt_counter} -eq $${max_attempts} ];then \
	      echo "After $${max_attempts} seconds Opensearch has still not started. Check the logs for errors:"; \
	      ${CONTAINER_RUNTIME} logs opensearch; \
	      exit 1; \
	    fi; \
	    printf "Waiting for Opensearch to start ... %3d/%d\r " $${attempt_counter} $${max_attempts}; \
	    attempt_counter=$$(($$attempt_counter+1)); \
	    sleep 1; \
	done; \
	echo; \
	echo "Opensearch has started on http://localhost:9200"; \
	max_attempts=60; \
	attempt_counter=0; \
	until $$(curl --output /dev/null --silent --head --fail http://localhost:5601); do \
	    if [ $${attempt_counter} -eq $${max_attempts} ];then \
	      echo "After $${max_attempts} seconds Opensearch dashboards has still not started. Check the logs for errors:"; \
	      ${CONTAINER_RUNTIME} logs opensearch_dashboards; \
	      exit 1; \
	    fi; \
	    printf "Waiting for Opensearch Dashboards to start ... %3d/%d\r " $${attempt_counter} $${max_attempts}; \
	    attempt_counter=$$(($$attempt_counter+1)); \
	    sleep 1; \
	done; \
	echo; \
	echo "Opensearch Dashboards has started on http://localhost:5601";