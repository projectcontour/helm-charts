# Helm Chart for Contour Gateway Provisioner

Deploys the Contour Gateway API provisioner controller using the upstream example manifest for dynamically provisioned Gateways.

## Installing the Chart

```console
helm repo add contour https://projectcontour.github.io/helm-charts/
helm repo update
helm install my-provisioner contour/contour-gateway-provisioner
```

> **Tip**: List all releases using `helm list` or `helm ls --all-namespaces`

## Local testing

Render manifests locally from this repo:

```console
helm template test ./charts/contour-gateway-provisioner
```

Override values inline for quick checks (example disabling RBAC):

```console
helm template test ./charts/contour-gateway-provisioner \
  --set rbac.create=false
```

## Configuration

| Name                     | Description                                                         | Value          |
| ------------------------ | ------------------------------------------------------------------- | -------------- |
| `image.registry`         | Contour image registry                                              | `ghcr.io`      |
| `image.repository`       | Contour image name                                                  | `projectcontour/contour` |
| `image.tag`              | Contour image tag                                                   | `v1.33.0`      |
| `image.pullPolicy`       | Image pull policy                                                   | `IfNotPresent` |
| `image.pullSecrets`      | Image pull secrets                                                  | `[]`           |
| `replicaCount`           | Provisioner controller replicas                                     | `1`            |
| `metricsAddress`         | Metrics bind address                                                | `127.0.0.1:8080` |
| `serviceAccount.create`  | Create a ServiceAccount for the provisioner                         | `true`         |
| `serviceAccount.name`    | Override ServiceAccount name                                        | `""`           |
| `serviceAccount.namespace` | Override ServiceAccount namespace used in RBAC subjects           | `""`           |
| `serviceAccount.annotations` | Annotations for the ServiceAccount                              | `{}`           |
| `serviceAccount.automountServiceAccountToken` | Automount ServiceAccount token                        | `true`         |
| `resources.requests`     | Resource requests for the controller                                | `cpu: 100m`, `memory: 70Mi` |
| `resources.limits`       | Resource limits for the controller                                  | `{}`           |
| `rbac.create`            | Create RBAC resources                                               | `true`         |
| `extraArgs`              | Extra CLI args appended to `contour gateway-provisioner`            | `[]`           |
