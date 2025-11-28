#! /usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DRY_RUN=true # Default to dry run mode: do not make commits or create a pull request unless --real-run is provided.
for arg in "$@"; do
    if [[ "$arg" == "--real-run" ]]; then
        DRY_RUN=false
    fi
done

git::exec() {
    if [ "$DRY_RUN" = true ]; then
        echo "[DRY RUN] git $*"
    else
        git "$@"
    fi
}

gh::exec() {
    if [ "$DRY_RUN" = true ]; then
        echo "[DRY RUN] gh $*"
    else
        gh "$@"
    fi
}

echo "Updating Helm chart versions"
go run ./hack/actions/bump-chart-versions/main.go

echo "Synchronizing CRDs"
go run ./hack/actions/synchronize-crds/main.go

if git diff --quiet; then
    echo "No update needed"
    exit 0
fi

# Read new chart and app versions.
readonly NEW_CHART_VERSION=$(grep '^version:' ./charts/contour/Chart.yaml | awk '{print $2}')
readonly NEW_APP_VERSION=$(grep '^appVersion:' ./charts/contour/Chart.yaml | awk '{print $2}')
readonly PR_BRANCH_NAME="github-actions/contour-${NEW_APP_VERSION}"

if git ls-remote --quiet origin "refs/heads/${PR_BRANCH_NAME}" | grep -q .; then
    echo "Pull request branch for ${PR_BRANCH_NAME} already exists on remote, skipping update"
    exit 0
fi

echo "Creating branch ${PR_BRANCH_NAME}"
git::exec checkout -b "${PR_BRANCH_NAME}"
git status --short      # Show the change files for logging purposes.
git::exec add --update  # Stage all modified files.

echo "Committing and pushing changes"
git::exec commit --signoff --message "Update Contour Helm chart to Contour ${NEW_APP_VERSION}"
git::exec push origin "${PR_BRANCH_NAME}"

echo "Creating pull request"
gh::exec pr create \
    --title "Update Contour Helm chart to Contour ${NEW_APP_VERSION}" \
    --body "This PR updates the Contour Helm chart to Contour version ${NEW_APP_VERSION} and chart version ${NEW_CHART_VERSION}." \
    --head "${PR_BRANCH_NAME}"
