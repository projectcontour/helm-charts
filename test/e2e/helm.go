// Copyright Project Contour Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/onsi/gomega"
)

type Helm struct {
	releaseName string
	namespace   string
}

const (
	helmInstallTimeout   = 5 * time.Minute
	helmUpgradeTimeout   = 5 * time.Minute
	helmUninstallTimeout = 1 * time.Minute
	helmRepoTimeout      = 3 * time.Minute
)

func HelmInstall(releaseName, chartPath, namespace string, additionalArgs ...string) *Helm {
	helm := &Helm{
		releaseName: releaseName,
		namespace:   namespace,
	}

	cmdArgs := append([]string{
		"helm", "install", releaseName, chartPath,
		"--namespace", namespace,
		"--create-namespace",
		"--wait",
	}, additionalArgs...)

	helm.runWithTimeout(cmdArgs, helmInstallTimeout)

	return helm
}

func (h *Helm) Upgrade(chartPath string, additionalArgs ...string) {
	cmdArgs := append([]string{
		"helm", "upgrade", h.releaseName, chartPath,
		"--namespace", h.namespace,
		"--wait",
	}, additionalArgs...)

	h.runWithTimeout(cmdArgs, helmUpgradeTimeout)
}

func (h *Helm) Uninstall() {
	h.runWithTimeout([]string{"helm", "uninstall", h.releaseName, "--namespace", h.namespace}, helmUninstallTimeout)
}

// run executes a helm command and fails the test if it exits non-zero.
func (h *Helm) runWithTimeout(cmdArgs []string, timeout time.Duration) {
	runCommand(cmdArgs[0], timeout, false, nil, cmdArgs[1:]...)
}

// HelmRepoAdd adds a Helm repository and updates its index.
func HelmRepoAdd(repoName, repoURL string) {
	runCommand("helm", helmRepoTimeout, false, nil, "repo", "add", repoName, repoURL)
	runCommand("helm", helmRepoTimeout, false, nil, "repo", "update", repoName)
}

// HelmSearchLatestVersion returns the latest version of a chart in a Helm repository.
func HelmSearchLatestVersion(repoName, chartName string) string {
	var stdout bytes.Buffer
	runCommand("helm", helmRepoTimeout, false, &stdout, "search", "repo", repoName+"/"+chartName, "--output", "json")

	var results []struct {
		Version string `json:"version"`
	}
	gomega.Expect(json.Unmarshal(stdout.Bytes(), &results)).To(gomega.Succeed())
	gomega.Expect(results).NotTo(gomega.BeEmpty(), "no versions found for %s/%s", repoName, chartName)

	return results[0].Version
}
