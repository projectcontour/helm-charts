# Changelog
## 0.3.0
* Centralize Contour and Gateway API CRDs in a shared library dependency to avoid duplicating manifests between charts.

## 0.2.0
* Contour upgraded to 1.33.0
* Envoy upgraded to 1.35.2

## 0.1.0
* Forked from [bitnami/charts/contour](https://github.com/bitnami/charts/tree/main/bitnami/contour) version 21.1.0
* Remove `defaultBackend` functionality
* Use images directly from contour and envoy instead of bitnami
