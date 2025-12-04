# Contour CRDs Helm Chart

A lightweight Helm chart that installs only the CustomResourceDefinitions required by [Contour](https://projectcontour.io) and (optionally) the Gateway API. Use it when you want to manage CRDs centrally instead of bundling them with the main Contour release.

## Prerequisites

- Kubernetes 1.23+
- Helm 3.8.0+

## Installing the Chart

To install the chart with the release name `contour-crds`:

```console
helm repo add contour https://projectcontour.github.io/helm-charts/
helm repo update
helm install contour-crds contour/contour-crds
```

## Using with the `contour` chart

If you install CRDs through this chart, disable CRD management in the `contour` chart to avoid conflicts:

```yaml
contour:
  manageCRDs: false
```

## Parameters

| Name | Description | Value |
| --- | --- | --- |
| `contour.manageCRDs` | Manage the creation, upgrade and deletion of Contour CRDs. | `true` |
| `gatewayAPI.manageCRDs` | Manage the creation, upgrade and deletion of Gateway API CRDs. | `false` |

## Local testing

Run basic linting and rendering checks before publishing:

```console
helm lint charts/contour-crds
helm template charts/contour-crds --namespace default
```

Apply to a local cluster like minikube

```console
helm install contour-crds charts/contour-crds --namespace default
# to install with the gateway api crds
helm install contour-crds ./charts/contour-crds --namespace default --set gatewayAPI.manageCRDs=true
# to remove
helm uninstall contour-crds --namespace default
```
