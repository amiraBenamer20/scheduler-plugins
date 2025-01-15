# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO_VERSION := 1.22.0
INTEGTESTENVVAR=SCHED_PLUGINS_TEST_VERBOSE=1

# Manage platform and builders
PLATFORMS ?= linux/amd64
BUILDER ?= docker
ifeq ($(BUILDER),podman)
	ALL_FLAG=--all
else
	ALL_FLAG=
endif

# REGISTRY is the container registry to push
REGISTRY ?= audhub

# Set a valid RELEASE_VERSION for testing
RELEASE_VERSION := v0.30.6
RELEASE_IMAGE := fspaas:kube-scheduler-$(RELEASE_VERSION)
RELEASE_CONTROLLER_IMAGE := controller-$(RELEASE_VERSION)

GO_BASE_IMAGE ?= golang:$(GO_VERSION)
DISTROLESS_BASE_IMAGE ?= gcr.io/distroless/static:nonroot
EXTRA_ARGS := ""

# VERSION is the scheduler's version
VERSION := $(shell echo $(RELEASE_VERSION) | awk -F - '{print $$2}')
VERSION := $(or $(VERSION),v0.0.$(shell date +%Y%m%d))

.PHONY: all
all: build

.PHONY: build
build: build-controller build-scheduler

.PHONY: build-controller
build-controller:
	@echo "Building controller binary..."
	@set -x
	$(GO_BUILD_ENV) go build -v -ldflags '-X k8s.io/component-base/version.gitVersion=$(VERSION) -w' -o bin/controller cmd/controller/controller.go || (echo "Controller build failed"; exit 1)

.PHONY: build-scheduler
build-scheduler:
	@echo "Building scheduler binary..."
	@echo "Using GO_BUILD_ENV=$(GO_BUILD_ENV)"
	@set -x
	$(GO_BUILD_ENV) go build -v -ldflags '-X k8s.io/component-base/version.gitVersion=$(VERSION) -w' -o bin/kube-scheduler cmd/scheduler/main.go || (echo "Scheduler build failed"; exit 1)

.PHONY: build-images
build-images:
	@echo "Building container images..."
	BUILDER=$(BUILDER) \
	PLATFORMS=$(PLATFORMS) \
	RELEASE_VERSION=$(RELEASE_VERSION) \
	REGISTRY=$(REGISTRY) \
	IMAGE=$(RELEASE_IMAGE) \
	CONTROLLER_IMAGE=$(RELEASE_CONTROLLER_IMAGE) \
	GO_BASE_IMAGE=$(GO_BASE_IMAGE) \
	DISTROLESS_BASE_IMAGE=$(DISTROLESS_BASE_IMAGE) \
	DOCKER_BUILDX_CMD=$(DOCKER_BUILDX_CMD) \
	EXTRA_ARGS=$(EXTRA_ARGS) hack/build-images.sh || (echo "Image build failed"; exit 1)

.PHONY: local-image
local-image: PLATFORMS="linux/$$(uname -m)"
local-image: RELEASE_VERSION="v0.0.0"
local-image: REGISTRY="localhost:5000/scheduler-plugins"
local-image: EXTRA_ARGS="--load"
local-image: clean build-images

.PHONY: release-images
push-images: EXTRA_ARGS="--push"
push-images: build-images

.PHONY: update-gomod
update-gomod:
	@echo "Updating Go modules..."
	hack/update-gomod.sh || (echo "Go module update failed"; exit 1)

.PHONY: unit-test
unit-test: install-envtest
	@echo "Running unit tests..."
	hack/unit-test.sh $(ARGS) || (echo "Unit tests failed"; exit 1)

.PHONY: install-envtest
install-envtest:
	@echo "Installing envtest..."
	hack/install-envtest.sh || (echo "Envtest installation failed"; exit 1)

.PHONY: integration-test
integration-test: install-envtest
	@echo "Running integration tests..."
	$(INTEGTESTENVVAR) hack/integration-test.sh $(ARGS) || (echo "Integration tests failed"; exit 1)

.PHONY: verify
verify:
	@echo "Verifying code format and structure..."
	hack/verify-gomod.sh || (echo "Go module verification failed"; exit 1)
	hack/verify-gofmt.sh || (echo "Code formatting verification failed"; exit 1)
	hack/verify-crdgen.sh || (echo "CRD generation verification failed"; exit 1)
	hack/verify-structured-logging.sh || (echo "Structured logging verification failed"; exit 1)
	hack/verify-toc.sh || (echo "Table of contents verification failed"; exit 1)

.PHONY: clean
clean:
	@echo "Cleaning up build artifacts..."
	rm -rf ./bin || (echo "Failed to clean up"; exit 1)
