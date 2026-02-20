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

//go:build none

// This script synchronizes the CRDs in the Helm chart with the ones from the Contour source code.
// It uses Chart.yaml appVersion to determine which Contour version to download.
//
// Usage:
//
//	go run hack/actions/synchronize-crds/main.go
package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/mholt/archives"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	log = logrus.StandardLogger()

	contourSourceURL = "https://github.com/projectcontour/contour/archive/refs/tags/v%s.tar.gz"

	chartPath = "./charts/contour/Chart.yaml"

	contourCRDSourcePath = "examples/contour/01-crds.yaml"
	contourCRDDestPath   = "./charts/contour/templates/crds/contour-crds.yaml"

	gatewayCRDSourcePath = "examples/gateway/00-crds.yaml"
	gatewayCRDDestPath   = "./charts/contour/templates/crds/gateway-api-crds.yaml"
)

func main() {
	log.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	// Define global HTTP client timeout.
	http.DefaultClient.Timeout = 2 * time.Minute

	// Read current app version.
	currentChartAppVersion, err := getCurrentChartAppVersion(chartPath)
	if err != nil {
		log.Fatalf("Failed to get current chart appVersion: %v", err)
	}
	log.Infof("Current chart appVersion: %s", currentChartAppVersion)

	// Download source code for the current Contour version.
	tmpDir, err := os.MkdirTemp("", "contour-source-")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}

	if err := syncCRDs(tmpDir, currentChartAppVersion); err != nil {
		os.RemoveAll(tmpDir)
		log.Fatalf("Failed to synchronize CRDs: %v", err)
	}
	os.RemoveAll(tmpDir)

	log.Infof("Successfully synchronized CRDs.")
}

func syncCRDs(tmpDir, currentChartAppVersion string) error {
	downloadPath := path.Join(tmpDir, "contour.tar.gz")
	log.Infof("Downloading Contour source tarball for version %s to %s", currentChartAppVersion, downloadPath)
	err := downloadFile(fmt.Sprintf(contourSourceURL, currentChartAppVersion), downloadPath)
	if err != nil {
		return fmt.Errorf("failed to download release: %w", err)
	}

	ctx := context.Background()
	tar, err := archives.FileSystem(ctx, downloadPath, nil)
	if err != nil {
		return fmt.Errorf("failed to open archive %s: %w", downloadPath, err)
	}

	fullContourCRDSourcePath := fmt.Sprintf("contour-%s/%s", currentChartAppVersion, contourCRDSourcePath)
	if err := copyCRD(tar, fullContourCRDSourcePath, contourCRDDestPath, ".Values.contour.manageCRDs"); err != nil {
		return fmt.Errorf("failed to copy Contour CRDs: %w", err)
	}
	log.Infof("Wrote Contour CRDs to %s", contourCRDDestPath)

	fullGatewayCRDSourcePath := fmt.Sprintf("contour-%s/%s", currentChartAppVersion, gatewayCRDSourcePath)
	if err := copyCRD(tar, fullGatewayCRDSourcePath, gatewayCRDDestPath, ".Values.gatewayAPI.manageCRDs"); err != nil {
		return fmt.Errorf("failed to copy Gateway API CRDs: %w", err)
	}
	log.Infof("Wrote Gateway API CRDs to %s", gatewayCRDDestPath)

	return nil
}

// getCurrentChartAppVersion reads the current chart and app version.
func getCurrentChartAppVersion(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	type ChartYaml struct {
		AppVersion string `yaml:"appVersion"`
	}

	var chart ChartYaml
	if err := yaml.Unmarshal(data, &chart); err != nil {
		return "", fmt.Errorf("failed to unmarshal yaml from %s: %w", filePath, err)
	}

	return chart.AppVersion, nil
}

// downloadFile downloads a file from the given URL to the specified destination path.
func downloadFile(sourceURL, destPath string) error {
	resp, err := http.Get(sourceURL) //nolint:gosec // G107: URL is constructed from a hardcoded constant
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", sourceURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s: status code %d", sourceURL, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if err := os.WriteFile(destPath, data, 0o644); err != nil { //nolint:gosec // G306: chart files are intentionally world-readable
		return fmt.Errorf("failed to write file %s: %w", destPath, err)
	}

	return nil
}

func copyCRD(fsys fs.FS, srcPath, destPath, conditional string) error {
	f, err := fsys.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", srcPath, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read source file %s: %w", srcPath, err)
	}

	if err := os.WriteFile(destPath, injectConditional(conditional, data), 0o644); err != nil { //nolint:gosec // G306: chart files are intentionally world-readable
		return fmt.Errorf("failed to write destination file %s: %w", destPath, err)
	}

	return nil
}

// injectConditional wraps the given data with Helm conditional statement.
func injectConditional(condition string, data []byte) []byte {
	return []byte(fmt.Sprintf("# Conditional: %s\n{{- if %s }}\n%s{{- end }}\n", condition, condition, string(data)))
}
