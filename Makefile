TEST?=$$(go list ./...)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=netbox
DIR_NAME=netbox
PLATFORM=linux
XC_ARCH=amd64 arm64
#XC_OS=linux darwin windows
XC_OS=darwin
GIT_COMMIT=$$(git rev-parse HEAD)
RELEASE_VERSION=$$(git for-each-ref refs/tags --sort=-taggerdate --format='%(refname:short)' --count=1)
PKG := github.com/fenglyu/terraform-provider-netbox
version ?= 0.1.5
LD_FLAGS=-X ${PKG}/version.ProviderVersion=${version} -X ${PKG}/version.GitCommit=${GIT_COMMIT}
#LD_FLAGS=" -s -w "
TESTARGS=-v
TAG=v2.8.9-testing
#TEST_NETBOX_IMAGE ?= docker-hub.battle.net/cloud/netbox:${TAG}
TEST_NETBOX_IMAGE ?= netboxcommunity/netbox:${TAG}


default: build

prep:
	@test ! -d pkg && mkdir pkg || true
	@command zip -h >/dev/null 2>&1 || ( echo "==> Please install zip command" )

gox:
	@echo "==> Installing gox..."
	@go install github.com/mitchellh/gox@v1.0.1

## comment out for terraform 0.12
#build-dev: fmtcheck generate
#	@[ -z "${version}" ] || ( echo "==> please use 'make build-dev version=vX.Y.Z'" )
#	go build  -ldflags="-X main.GitCommit=${GIT_COMMIT}" -o ~/.terraform.d/plugins/terraform-provider-$(PKG_NAME)_${version} .

build-dev: fmtcheck generate
	@[ -z "${version}" ] || ( echo "==> please use 'make build-dev version=vX.Y.Z'" )
	#go build  -ldflags="-X main.GitCommit=${GIT_COMMIT}" -o ~/.terraform.d/plugins/registry.terraform.io/-/netbox/${version}/${PLATFORM}_${XC_ARCH}/terraform-provider-$(PKG_NAME)_v${version} .
	go build  -ldflags="-X main.GitCommit=${GIT_COMMIT}" -o /usr/local/terraform/plugins/terraform.cloud.blizzard.net/cf/netbox/${version}/${PLATFORM}_${XC_ARCH}/terraform-provider-$(PKG_NAME)_v${version} .

build: fmtcheck generate prep gox
	@echo "==> Building..."
	@CGO_ENABLED=0 gox -os="$(XC_OS)" -arch="$(XC_ARCH)" -ldflags "$(LD_FLAGS)" -output "pkg/{{.OS}}_{{.Arch}}/terraform-provider-$(PKG_NAME)_v${RELEASE_VERSION}" .

release: build $(eval SHELL:=/bin/bash)
	@for PLATFORM in $$(find ./pkg -mindepth 1 -maxdepth 1 -type d); do \
		OSARCH=$$(basename $$PLATFORM); \
		echo "--> $$OSARCH"; \
		pushd $$PLATFORM >/dev/null 2>&1; \
		zip ../$$OSARCH.zip ./*; \
		popd >/dev/null 2>&1; \
	done
	@echo -ne "\n==> Results:\n"
	@find pkg/ -type f -exec ls -sh '{}' \;
	@## upload to platform, TBD

clean:
	@rm -rf pkg/

test: fmtcheck generate
	go test $(TESTARGS) -timeout=30s $(TEST)

test-netbox-env-up: $(eval SHELL:=/bin/bash)
	pushd tests/netbox-docker >/dev/null 2>&1; \
	NETBOX_IMAGE="$(TEST_NETBOX_IMAGE)" docker-compose up -d ; \
	#NETBOX_IMAGE="$(TEST_NETBOX_IMAGE)" docker-compose up ; \
	popd >/dev/null 2>&1; \

test-netbox-env-down: $(eval SHELL:=/bin/bash)
	pushd tests/netbox-docker >/dev/null 2>&1; \
	NETBOX_IMAGE=$(TEST_NETBOX_IMAGE) docker-compose down; \
	popd >/dev/null 2>&1; \
	docker volume rm netbox-docker_netbox-static-files; \
	docker volume rm netbox-docker_netbox-nginx-config; \
	docker volume rm netbox-docker_netbox-media-files; \
	docker volume rm netbox-docker_netbox-postgres-data; \
	docker volume rm netbox-docker_netbox-redis-data


testacc: fmtcheck
	TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test $(TEST) $(TESTARGS) -timeout 240m -ldflags="-X=github.com/fenglyu/terraform-provider-netbox/version.ProviderVersion=acc"

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -w -s ./$(DIR_NAME)

# Currently required by tf-deploy compile
fmtcheck:
	@echo "==> Checking source code against gofmt..."
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run ./$(DIR_NAME)

tools:
	@echo "==> installing required tooling..."
	go install github.com/client9/misspell/cmd/misspell
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

generate:
	go generate  ./...

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(DIR_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

docscheck:
	@sh -c "'$(CURDIR)/scripts/docscheck.sh'"

.PHONY: build-dev build test prep release vet fmt fmtcheck lint tools errcheck test-compile generate website website-test docscheck generate

