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
	"os/exec"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func Kubectl(args ...string) {
	cmdOutputWriter := gexec.NewPrefixedWriter("[kubectl] ", ginkgo.GinkgoWriter)

	cmdArgs := append([]string{"kubectl"}, args...)

	//nolint:gosec // G204: Subprocess launched with a potential tainted input or cmd arguments
	session, err := gexec.Start(exec.Command(cmdArgs[0], cmdArgs[1:]...), cmdOutputWriter, cmdOutputWriter)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Eventually(session, 5*time.Minute).Should(gexec.Exit(0))
}
