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

	"github.com/bombsimon/logrusr/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Framework provides a collection of helpful functions for
// writing end-to-end (E2E) tests for Contour.
type Framework struct {
	// HTTP provides helpers for making HTTP/HTTPS requests.
	HTTP *HTTP
}

func NewFramework() *Framework {
	t := ginkgo.GinkgoT()

	// Deferring GinkgoRecover() provides better error messages in case of panic
	// e.g. when CONTOUR_E2E_LOCAL_HOST environment variable is not set.
	defer ginkgo.GinkgoRecover()

	log.SetLogger(logrusr.New(logrus.StandardLogger()))

	ipV6Cluster := os.Getenv("IPV6_CLUSTER") == "true"

	httpURLBase := os.Getenv("CONTOUR_E2E_HTTP_URL_BASE")
	if httpURLBase == "" {
		if ipV6Cluster {
			httpURLBase = "http://[::1]:9080"
		} else {
			httpURLBase = "http://127.0.0.1:9080"
		}
	}

	httpsURLBase := os.Getenv("CONTOUR_E2E_HTTPS_URL_BASE")
	if httpsURLBase == "" {
		if ipV6Cluster {
			httpsURLBase = "https://[::1]:9443"
		} else {
			httpsURLBase = "https://127.0.0.1:9443"
		}
	}

	return &Framework{
		HTTP: &HTTP{
			HTTPURLBase:   httpURLBase,
			HTTPSURLBase:  httpsURLBase,
			RetryInterval: time.Second,
			RetryTimeout:  60 * time.Second,
			t:             t,
		},
	}
}

type NamespacedTestBody func(string)

func (f *Framework) NamespacedTest(namespace string, body NamespacedTestBody, additionalNamespaces ...string) {
	ginkgo.Context("with namespace: "+namespace, func() {
		ginkgo.BeforeEach(func() {
			RecreateKindCluster()
			for _, ns := range append(additionalNamespaces, namespace) {
				f.CreateNamespace(ns)
			}
		})
		ginkgo.AfterEach(func() {
			DeleteKindCluster()
		})

		body(namespace)
	})
}

// CreateNamespace creates a namespace with the given name in the
// Kubernetes API or fails the test if it encounters an error.
func (f *Framework) CreateNamespace(name string) {
	Kubectl("create", "namespace", name)
}

// DeleteNamespace deletes the namespace with the given name in the
// Kubernetes API or fails the test if it encounters an error.
func (f *Framework) DeleteNamespace(name string, waitForDeletion bool) {
	Kubectl("delete", "namespace", name)

	if waitForDeletion {
		Kubectl("wait", "--for=delete", "--timeout=2m", "namespace/"+name)
	}
}
