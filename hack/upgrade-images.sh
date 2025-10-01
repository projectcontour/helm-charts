#!/usr/bin/env bash
set -euo pipefail

curl -Lo versions.yaml https://raw.githubusercontent.com/projectcontour/contour/refs/heads/main/versions.yaml
# TODO: yq doesn't have any native support for semver so let's assume the file is sorted for now
yq '.versions | filter(.version != "main") | filter(.supported == true) | .[0]' versions.yaml > latest.yaml

contourVersion="$(yq '.version | sub("v", "")' latest.yaml)"
envoyVersion="$(yq .dependencies.envoy latest.yaml)"
root="$(git rev-parse --show-toplevel)"

sed -I '' "s/appVersion: .*/appVersion: $contourVersion/" "${root}/charts/contour/Chart.yaml"

contourLineNum="$(grep -n "    repository: projectcontour/contour" "${root}/charts/contour/values.yaml" | cut -d ':' -f1)"
sed -I '' "$((contourLineNum+1))s/tag: .*/tag: v${contourVersion}/" "${root}/charts/contour/values.yaml"

envoyLineNum="$(grep -n "    repository: envoyproxy/envoy" "${root}/charts/contour/values.yaml" | cut -d ':' -f1)"
sed -I '' "$((envoyLineNum+1))s/tag: .*/tag: v${envoyVersion}/" "${root}/charts/contour/values.yaml"
