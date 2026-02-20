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
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContourHelmChart(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Contour Helm Chart E2E Suite")
}

var f = NewFramework()

var _ = Describe("Contour", func() {
	const (
		releaseName = "contour"
		chartPath   = "../../charts/contour"
		repoName    = "contour"
		repoURL     = "https://projectcontour.github.io/helm-charts/"
	)

	// Required values for Contour to become ready in Kind cluster and be reachable for HTTP requests.
	mandatoryInstallArgs := []string{
		"--set", "envoy.service.type=NodePort",
		"--set", "envoy.useHostPort.http=true",
		"--set", "envoy.useHostPort.https=true",
	}

	assertEchoServesTraffic := func(message string) {
		res, ok := f.HTTP.RequestUntil(&HTTPRequestOpts{
			Host:      "echoserver.projectcontour.io",
			Path:      "/",
			Condition: HasStatusCode(200),
		})

		Expect(ok).To(BeTrue(), message)
		Expect(res.StatusCode).To(Equal(200))
	}

	f.NamespacedTest("test-helm-installation", func(namespace string) {
		It("should deploy contour using helm", func() {
			helmRelease := HelmInstall(releaseName, chartPath, namespace, mandatoryInstallArgs...)
			defer helmRelease.Uninstall()

			DeployEcho(namespace)
			assertEchoServesTraffic("expected to receive 200 OK from echoserver")
		})
	})

	f.NamespacedTest("test-helm-upgrade", func(namespace string) {
		It("should upgrade contour from previous chart version", func() {
			HelmRepoAdd(repoName, repoURL)
			previousVersion := HelmSearchLatestVersion(repoName, "contour")

			By("installing previous version " + previousVersion + " from Helm repo")
			upgradeBaseInstallArgs := append([]string{"--version", previousVersion}, mandatoryInstallArgs...)
			helmRelease := HelmInstall(releaseName, repoName+"/contour", namespace, upgradeBaseInstallArgs...)
			defer helmRelease.Uninstall()

			DeployEcho(namespace)
			assertEchoServesTraffic("expected 200 OK from echoserver before upgrade")

			By("upgrading to current local chart")
			helmRelease.Upgrade(chartPath, mandatoryInstallArgs...)
			assertEchoServesTraffic("expected 200 OK from echoserver after upgrade")
		})
	})
})
