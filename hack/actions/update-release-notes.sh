#! /usr/bin/env bash

# update-release-notes.sh: Write the release notes for a new chart release.
# chart-releaser uses the chart description as release notes. This script
# changes them to a short intro with the Contour version, and the list of
# changes that GitHub generates.

set -o errexit
set -o nounset
set -o pipefail

readonly HERE=$(cd "$(dirname "$0")" && pwd)
readonly REPO=$(cd "${HERE}/../.." && pwd)

DRY_RUN=true # Do not change the release unless --real-run is given.
for arg in "$@"; do
    if [[ "$arg" == "--real-run" ]]; then
        DRY_RUN=false
    fi
done

cd "${REPO}"

readonly CHART_VERSION=$(grep '^version:' ./charts/contour/Chart.yaml | awk '{print $2}')
readonly APP_VERSION=$(grep '^appVersion:' ./charts/contour/Chart.yaml | awk '{print $2}')
readonly TAG="contour-${CHART_VERSION}"

readonly INTRODUCTION="Helm chart for [Contour ${APP_VERSION}](https://github.com/projectcontour/contour/releases/tag/v${APP_VERSION}) - an open source Kubernetes ingress controller that works by deploying the Envoy proxy as a reverse proxy and load balancer."

# chart-releaser creates the tag on the current commit. If the tag is on some
# other commit, this is an old release and we do not touch it.
# Do not use readonly with gh calls. It would hide errors from gh.
TAG_COMMIT=$(gh api "repos/{owner}/{repo}/commits/${TAG}" --jq .sha)
if [[ "${TAG_COMMIT}" != "$(git rev-parse HEAD)" ]]; then
    echo "Release ${TAG} was not created from this commit, skipping"
    exit 0
fi

echo "Generating release notes for ${TAG}"
GENERATED_NOTES=$(gh api "repos/{owner}/{repo}/releases/generate-notes" -f tag_name="${TAG}" --jq .body)

readonly NEW_NOTES="${INTRODUCTION}

${GENERATED_NOTES}"

echo "----"
echo "${NEW_NOTES}"
echo "----"

if [[ "$DRY_RUN" == true ]]; then
    echo "[DRY RUN] gh release edit ${TAG} --notes <notes above>"
else
    gh release edit "${TAG}" --notes "${NEW_NOTES}"
fi
