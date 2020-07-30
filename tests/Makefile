NAMESPACE ?= cloud
git_describe = $(shell git describe)
vcs_ref := $(shell git rev-parse HEAD)
build_date := $(shell date -u +%FT%T)
hadolint_available := $(shell hadolint --help > /dev/null 2>&1; echo $$?)
hadolint_command := hadolint --ignore DL3003 --ignore DL3008 --ignore DL3013 --ignore DL3018 --ignore DL3028 --ignore DL4000 --ignore DL4001
hadolint_container := hadolint/hadolint:latest

VERSION ?= v2.8.8
IS_LATEST := false
dockerfile := Dockerfile

prep:
	@git fetch --unshallow ||:
	@git fetch origin 'refs/tags/*:refs/tags/*'

lint:
	@docker pull $(hadolint_container)
	@docker run --rm -v $(PWD)/netbox/nginx/$(dockerfile):/Dockerfile \
		-i $(hadolint_container) $(hadolint_command) Dockerfile
	@docker run --rm -v $(PWD)/netbox/netbox/Dockerfile:/Dockerfile \
		-i $(hadolint_container) $(hadolint_command) Dockerfile

build: prep
	docker build \
		--pull \
		--build-arg vcs_ref=$(vcs_ref) \
		--build-arg build_date=$(build_date) \
		--build-arg version=$(VERSION) \
		--file netbox/nginx/Dockerfile \
		--tag $(NAMESPACE)/nginx:$(VERSION) $(PWD)/netbox/nginx/
	docker build \
		--build-arg namespace=$(NAMESPACE) \
		--build-arg vcs_ref=$(vcs_ref) \
		--build-arg build_date=$(build_date) \
		--build-arg version=$(VERSION) \
		--file netbox/netbox/Dockerfile \
		--tag $(NAMESPACE)/netbox:$(VERSION) $(PWD)/netbox/netbox/

ifeq ($(IS_LATEST),true)
	@docker tag $(NAMESPACE)/nginx:$(VERSION) \
		$(NAMESPACE)/nginx:$(LATEST_VERSION)
	@docker tag $(NAMESPACE)/netbox:$(VERSION) \
    	$(NAMESPACE)/netbox:$(LATEST_VERSION)
endif

compose-up:
	@(cd netbox; docker-compose -p netbox_v2.8.8  up -d)

compose-down:
	@(cd netbox; docker-compose -p netbox_v2.8.8 down)

push-image: prep
	@docker push $(NAMESPACE)/netbox:$(VERSION)
ifeq ($(IS_LATEST),true)
	@docker push $(NAMESPACE)/netbox:$(LATEST_VERSION)
endif

publish: push-image

.PHONY: prep lint build publish push-image compose-up compose-down
