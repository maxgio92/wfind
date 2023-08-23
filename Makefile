APP := wfind
VERSION := 0.1.0

user := maxgio92
oci_image := quay.io/$(user)/$(APP)

bins := docker git go gofumpt golangci-lint

PACKAGE_NAME          := github.com/$(user)/$(APP)
GOLANG_CROSS_VERSION  ?= v$(shell sed -nE 's/go[[:space:]]+([[:digit:]]\.[[:digit:]]+)/\1/p' go.mod)

GIT_HEAD_COMMIT ?= $$($(git) rev-parse --short HEAD)
GIT_TAG_COMMIT  ?= $$($(git) rev-parse --short $(VERSION))
GIT_MODIFIED_1  ?= $$($(git) diff $(GIT_HEAD_COMMIT) $(GIT_TAG_COMMIT) --quiet && echo "" || echo ".dev")
GIT_MODIFIED_2  ?= $$($(git) diff --quiet && echo "" || echo ".dirty")
GIT_MODIFIED    ?= $$(echo "$(GIT_MODIFIED_1)$(GIT_MODIFIED_2)")
GIT_REPO        ?= $$($(git) config --get remote.origin.url)
BUILD_DATE      ?= $$($(git) log -1 --format="%at" | xargs -I{} date -d @{} +%Y-%m-%dT%H:%M:%S)

define declare_binpaths
$(1) = $(shell command -v 2>/dev/null $(1))
endef

$(foreach bin,$(bins),\
	$(eval $(call declare_binpaths,$(bin)))\
)

.PHONY: doc
doc:
	@go run docs/gen.go

.PHONY: build
build:
	@$(go) build .

.PHONY: run
run:
	@$(go) run .

.PHONY: test
test:
	@$(go) test -v -cover -gcflags=-l ./...

.PHONY: lint
lint: golangci-lint
	@$(golangci-lint) run ./...

.PHONY: golangci-lint
golangci-lint:
	@$(go) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0

.PHONY: gofumpt
gofumpt:
	@$(go) install mvdan.cc/gofumpt@v0.3.1

.PHONY: oci/build
oci/build: test
	@$(docker) build . -t $(oci_image):$(VERSION) -f Dockerfile \
		--build-arg GIT_HEAD_COMMIT=$(GIT_HEAD_COMMIT) \
 		--build-arg GIT_TAG_COMMIT=$(GIT_TAG_COMMIT) \
 		--build-arg GIT_MODIFIED=$(GIT_MODIFIED) \
 		--build-arg GIT_REPO=$(GIT_REPO) \
 		--build-arg GIT_LAST_TAG=$(VERSION) \
 		--build-arg BUILD_DATE=$(BUILD_DATE)

.PHONY: oci/push
oci/push: oci/build
	@$(docker) push $(oci_image):$(VERSION)

.PHONY: clean
clean:
	@rm -f $(APP)

.PHONY: help
help: list

.PHONY: list
list:
	@LC_ALL=C $(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

