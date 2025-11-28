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
	f.NamespacedTest("test-helm-installation", func(namespace string) {
		releaseName := "contour"
		chartPath := "../../charts/contour"

		It("should deploy contour using helm", func() {
			helm := HelmInstall(releaseName, chartPath, namespace,
				"--set", "envoy.useHostPort.http=true",
				"--set", "envoy.useHostPort.https=true",
			)

			NewEchoDeploy(namespace)

			res, ok := f.HTTP.RequestUntil(&HTTPRequestOpts{
				Host:      "echoserver.projectcontour.io",
				Path:      "/",
				Condition: HasStatusCode(200),
			})

			Expect(ok).To(BeTrue(), "expected to receive 200 OK from echoserver")
			Expect(res.StatusCode).To(Equal(200))

			helm.Uninstall()
		})
	})
})
