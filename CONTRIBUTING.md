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

### Running checks

```bash
make checkall
```

See `make help` for a full list of available targets.

### Running E2E tests

E2E tests are automatically run by CI when you submit a pull request, so running them locally is optional.

The test suite verifies chart installation and upgrade scenarios against a kind cluster.
To create a cluster, run the full end-to-end test suite, and clean up afterwards, execute:

```bash
make e2e
```
