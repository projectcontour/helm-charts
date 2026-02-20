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

.PHONY: lint
lint: ## Run all lint checks
lint: lint-golint lint-helm

.PHONY: lint-golint
lint-golint: ## Run Go linter
	@echo Running Go linter ...
	@./hack/golangci-lint run --build-tags=e2e,none

.PHONY: lint-helm
lint-helm: ## Run Helm linter
	@echo Running Helm linter ...
	@helm lint --strict charts/contour/

.PHONY: e2e
e2e: ## Run e2e tests against Kind cluster
	CONTOUR_E2E_HTTP_URL_BASE=$(CONTOUR_E2E_HTTP_URL_BASE) \
	CONTOUR_E2E_HTTPS_URL_BASE=$(CONTOUR_E2E_HTTPS_URL_BASE) \
	go run github.com/onsi/ginkgo/v2/ginkgo -tags=e2e -mod=readonly -keep-going -randomize-suites -randomize-all -poll-progress-after=120s --focus '$(CONTOUR_E2E_TEST_FOCUS)' $(CONTOUR_E2E_GINKGO_ARGS) -r $(CONTOUR_E2E_PACKAGE_FOCUS)

help: ## Display this help
	@echo Targets:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9._-]+:.*?## / {printf "  %-25s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
