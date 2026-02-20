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
	"os"
	"time"
)

const kindClusterName = "contour-e2e"

func CreateKindCluster() {
	args := []string{"create", "cluster", "--name", kindClusterName}
	if config := kindConfigPath(); config != "" {
		args = append(args, "--config", config)
	}
	runKind(args...)
}

func DeleteKindCluster() {
	runKindAllowFailure("delete", "cluster", "--name", kindClusterName)
}

func RecreateKindCluster() {
	DeleteKindCluster()
	CreateKindCluster()
}

func kindConfigPath() string {
	ipv6Cluster := os.Getenv("IPV6_CLUSTER") == "true"
	configFile := "kind-expose-port.yaml"
	if ipv6Cluster {
		configFile = "kind-ipv6.yaml"
	}

	for _, p := range []string{
		"test/scripts/" + configFile,
		"../../test/scripts/" + configFile,
	} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

func runKind(args ...string) {
	runCommand("kind", 10*time.Minute, false, nil, args...)
}

func runKindAllowFailure(args ...string) {
	runCommand("kind", 2*time.Minute, true, nil, args...)
}
