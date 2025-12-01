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
	"io"
	"os/exec"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type Helm struct {
	releaseName     string
	namespace       string
	cmdOutputWriter io.Writer
}

func HelmInstall(releaseName, chartPath, namespace string, additionalArgs ...string) *Helm {
	helm := &Helm{
		releaseName:     releaseName,
		cmdOutputWriter: gexec.NewPrefixedWriter("[helm] ", ginkgo.GinkgoWriter),
		namespace:       namespace,
	}

	cmdArgs := append([]string{
		"install", releaseName, chartPath,
		"--namespace", namespace,
		"--create-namespace",
		"--wait",
	}, additionalArgs...)

	session, err := gexec.Start(exec.Command("helm", cmdArgs...), helm.cmdOutputWriter, helm.cmdOutputWriter)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Eventually(session, 5*time.Minute).Should(gexec.Exit(0))

	return helm
}

func (h *Helm) Uninstall() {
	cmdArgs := []string{"uninstall", h.releaseName, "--namespace", h.namespace}
	session, err := gexec.Start(exec.Command("helm", cmdArgs...), h.cmdOutputWriter, h.cmdOutputWriter)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Eventually(session, 1*time.Minute).Should(gexec.Exit(0))
}
