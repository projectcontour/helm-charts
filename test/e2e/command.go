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
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func runCommand(tool string, timeout time.Duration, allowFailure bool, stdout io.Writer, args ...string) *gexec.Session {
	writer := gexec.NewPrefixedWriter("["+tool+"] ", ginkgo.GinkgoWriter)
	cmdArgs := append([]string{tool}, args...)

	fmt.Fprintf(writer, "Running: %s\n", strings.Join(cmdArgs, " "))

	commandStdout := io.Writer(writer)
	if stdout != nil {
		commandStdout = io.MultiWriter(stdout, writer)
	}

	//nolint:gosec // G204: Subprocess launched with dynamic command arguments in controlled test helpers.
	session, err := gexec.Start(exec.Command(cmdArgs[0], cmdArgs[1:]...), commandStdout, writer)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Eventually(session, timeout).Should(gexec.Exit())

	if !allowFailure {
		gomega.Expect(session.ExitCode()).To(gomega.Equal(0),
			"%s command failed: %s\nOutput:\n%s%s",
			tool,
			strings.Join(cmdArgs, " "),
			string(session.Out.Contents()),
			string(session.Err.Contents()),
		)
	}

	return session
}
