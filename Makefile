CONTOUR_E2E_PACKAGE_FOCUS ?= ./test/e2e
# Optional variables
# Run specific test specs (matched by regex)
# Example: CONTOUR_E2E_TEST_FOCUS="should deploy contour using helm"
CONTOUR_E2E_TEST_FOCUS ?=
# Override Envoy ingress address when not using kind with host ports.
# Example: CONTOUR_E2E_HTTP_URL_BASE=http://192.168.1.100:80 CONTOUR_E2E_HTTPS_URL_BASE=https://192.168.1.100:443
CONTOUR_E2E_HTTP_URL_BASE ?=
CONTOUR_E2E_HTTPS_URL_BASE ?=
# Additional ginkgo args, for example for verbose logs.
CONTOUR_E2E_GINKGO_ARGS ?= --vv --output-interceptor-mode=none --fail-fast
KIND ?= kind
CLUSTERNAME ?= contour-e2e

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
	$(KIND) create cluster --name $(CLUSTERNAME) --config test/scripts/kind-expose-port.yaml

.PHONY: e2e
e2e: | setup-kind-cluster run-e2e cleanup-kind ## Run E2E tests against a real k8s cluster

.PHONY: run-e2e
run-e2e:
	CONTOUR_E2E_HTTP_URL_BASE=$(CONTOUR_E2E_HTTP_URL_BASE) \
	CONTOUR_E2E_HTTPS_URL_BASE=$(CONTOUR_E2E_HTTPS_URL_BASE) \
	go run github.com/onsi/ginkgo/v2/ginkgo -tags=e2e -mod=readonly -keep-going -randomize-suites -randomize-all -poll-progress-after=120s --focus '$(CONTOUR_E2E_TEST_FOCUS)' $(CONTOUR_E2E_GINKGO_ARGS) -r $(CONTOUR_E2E_PACKAGE_FOCUS)

.PHONY: cleanup-kind
cleanup-kind: ## Delete the kind cluster
	$(KIND) delete cluster --name $(CLUSTERNAME)

help: ## Display this help
	@echo Targets:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9._-]+:.*?## / {printf "  %-25s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
