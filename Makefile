TEST?=$$(go list ./...)
PKG_NAME=netbox
DIR_NAME=netbox

default: build

build-dev:
#  @[ "${version}" ] || ( echo ">> please provide version=vX.Y.Z"; exit 1 )
	go build -o ~/.terraform.d/plugins/terraform-provider-$(PKG_NAME)_${version} .

#testacc: fmtcheck
#	TF_ACC=1 go test $(TEST) -v -count $(TEST_COUNT) -parallel 20 $(TESTARGS) -timeout 120m

build: fmtcheck generate
	go install

test: fmtcheck generate
	go test $(TESTARGS) -timeout=30s $(TEST)

testacc: fmtcheck
	TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test $(TEST) -v $(TESTARGS) -timeout 240m -ldflags="-X=github.com/fenglyu/terraform-provider-netbox/version.ProviderVersion=acc"

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



.PHONY: build-dev build test  vet fmt fmtcheck lint tools errcheck test-compile generate

