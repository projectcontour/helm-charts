# Contributing

Thanks for taking the time to join our community and start contributing.
These guidelines will help you get started with the Contour Helm Charts project.
Please note that we require [DCO sign off](https://github.com/projectcontour/contour/blob/main/CONTRIBUTING.md#dco-sign-off).

## Building from source

### Prerequisites

To run e2e tests, you will need to have the following tools installed:

1. Install [Go](https://go.dev/doc/install)
2. Install [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation), [kubectl](https://kubernetes.io/docs/tasks/tools/) and [Helm](https://helm.sh/docs/intro/install/).

### Fetch the source

1. [Fork](https://docs.github.com/en/github/getting-started-with-github/fork-a-repo#fork-an-example-repository) the `projectcontour/helm-charts` repository.
2. Clone your fork:

   ```bash
   git clone git@github.com:YOUR-USERNAME/helm-charts.git
   ```
### Make targets

Run `make help` to see all available targets.

Common targets:
- `make lint` - Run all lint checks
- `make lint-helm` - Run Helm lint only
- `make lint-golint` - Run Go lint only

### Running E2E tests

E2E tests run automatically in CI when you submit a pull request, but you can run them locally.

The test suite verifies chart installation and upgrade scenarios against a kind cluster:

```bash
make e2e
```

This command creates a cluster, runs the full end-to-end test suite, and cleans up afterwards.
