CONTOUR_E2E_PACKAGE_FOCUS ?= ./test/e2e
# Optional variables
# Run specific test specs (matched by regex)
CONTOUR_E2E_TEST_FOCUS ?=

.PHONY: all
all: help

.PHONY: checkall
checkall: lint ## Run all checks

.PHONY: lint
lint: ## Run lint checks
lint: lint-golint

.PHONY: lint-golint
lint-golint: ## Run Go linter
	@echo Running Go linter ...
	@./hack/golangci-lint run --build-tags=e2e,none

.PHONY: setup-kind-cluster
setup-kind-cluster: ## Make a kind cluster for testing
	./test/scripts/make-kind-cluster.sh

.PHONY: e2e
e2e: | setup-kind-cluster run-e2e cleanup-kind ## Run E2E tests against a real k8s cluster

.PHONY: run-e2e
run-e2e:
	go run github.com/onsi/ginkgo/v2/ginkgo -tags=e2e -mod=readonly -keep-going -randomize-suites -randomize-all -poll-progress-after=120s --focus '$(CONTOUR_E2E_TEST_FOCUS)' -r $(CONTOUR_E2E_PACKAGE_FOCUS)

.PHONY: cleanup-kind
cleanup-kind: ## Delete the kind cluster
	./test/scripts/cleanup.sh

help: ## Display this help
	@echo Targets:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9._-]+:.*?## / {printf "  %-25s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
