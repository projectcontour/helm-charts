#!/usr/bin/env sh

here="$(cd -- "$(dirname "$0")" > /dev/null 2>&1 || exit; pwd -P)"
# replace this with `yq .appVersion "${here}/Chart.yaml"` if yq is reliably available
version="$(grep 'appVersion:' "${here}/Chart.yaml" | cut -d':' -f2 | tr -d ' ' | cut -d'"' -f2 | cut -d"'" -f2)"

curl -L -o "${here}/src/raw-gateway-crds.yaml" "https://raw.githubusercontent.com/projectcontour/contour/v${version}/examples/gateway/00-crds.yaml"
curl -L -o "${here}/templates/contour-crds.yaml" "https://raw.githubusercontent.com/projectcontour/contour/v${version}/examples/contour/01-crds.yaml"