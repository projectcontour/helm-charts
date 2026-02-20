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

// This script determines the latest stable Contour version from the Contour repository
// it then bumps versions accordingly:
//
// - charts/contour/Chart.yaml appVersion to the latest stable Contour version.
// - charts/contour/values.yaml Contour and Envoy image tags to match the latest stable versions.
// - charts/contour/Chart.yaml minor version is incremented by one.
//
// If the current chart appVersion is already the latest stable Contour version, no changes are made.
//
// Usage:
//
//	go run hack/actions/bump-chart-versions/main.go
package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	log = logrus.StandardLogger()

	contourVersionsURL = "https://raw.githubusercontent.com/projectcontour/contour/refs/heads/main/versions.yaml"

	chartPath  = "./charts/contour/Chart.yaml"
	valuesPath = "./charts/contour/values.yaml"
)

func main() {
	log.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	// Define global HTTP client timeout.
	http.DefaultClient.Timeout = 2 * time.Minute

	// Read current chart and app versions.
	currentChartVersion, currentChartAppVersion, err := getCurrentChartVersions(chartPath)
	if err != nil {
		log.Fatalf("Failed to get current chart versions: %v", err)
	}
	log.Infof("Current chart version: %s, appVersion: %s", currentChartVersion, currentChartAppVersion)

	// Get latest stable Contour and Envoy versions from Github.
	contourVersion, envoyVersion, err := getLatestStableVersions()
	if err != nil {
		log.Fatalf("Failed to get latest stable versions: %v", err)
	}
	log.Infof("Latest stable Contour: %s, Envoy: %s", contourVersion, envoyVersion)

	// Compare if versions are the same.
	if contourVersion == currentChartAppVersion {
		log.Infof("Contour version %s is already up to date", contourVersion)
		return
	}

	// Update Chart.yaml with new minor chart version and appVersion based on latest Contour version info.
	nextChartVersion, err := nextMinorVersion(currentChartVersion)
	if err != nil {
		log.Fatalf("Failed to get next minor version: %v", err)
	}
	err = setYAMLField(chartPath, "version", nextChartVersion)
	if err != nil {
		log.Fatalf("Failed to update Contour chart version: %v", err)
	}
	log.Infof("Updated Contour chart version to %s in %s", nextChartVersion, chartPath)

	err = setYAMLField(chartPath, "appVersion", contourVersion)
	if err != nil {
		log.Fatalf("Failed to update Contour chart appVersion: %v", err)
	}
	log.Infof("Updated Contour chart appVersion to %s in %s", contourVersion, chartPath)

	// Update values.yaml with new Contour and Envoy versions.
	contourImageTag := fmt.Sprintf("v%s", contourVersion)
	err = setYAMLField(valuesPath, "contour.image.tag", contourImageTag)
	if err != nil {
		log.Fatalf("Failed to update Contour version: %v", err)
	}
	log.Infof("Updated Contour image tag to %s in %s", contourImageTag, valuesPath)

	envoyImageTag := fmt.Sprintf("v%s", envoyVersion)
	err = setYAMLField(valuesPath, "envoy.image.tag", envoyImageTag)
	if err != nil {
		log.Fatalf("Failed to update Envoy version: %v", err)
	}
	log.Infof("Updated Envoy image tag to %s in %s", envoyImageTag, valuesPath)

	log.Infof("Successfully bumped versions.")
}

func getLatestStableVersions() (string, string, error) {
	resp, err := http.Get(contourVersionsURL) //nolint:gosec // G107: URL is constructed from a hardcoded constant
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to fetch versions.yaml: status code %d", resp.StatusCode)
	}

	type VersionEntry struct {
		Version      string `yaml:"version"`
		Supported    string `yaml:"supported"`
		Dependencies struct {
			Envoy string `yaml:"envoy"`
		} `yaml:"dependencies"`
	}
	type VersionsYaml struct {
		Versions []VersionEntry `yaml:"versions"`
	}

	var versions VersionsYaml
	if err := yaml.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return "", "", err
	}

	for _, entry := range versions.Versions {
		// Skip "main" version.
		if entry.Version == "main" {
			continue
		}

		// Pick the first supported version as the latest stable.
		// Note: assumes versions.yaml is ordered from latest to oldest!
		if entry.Supported == "true" {
			contourVersion := strings.TrimPrefix(entry.Version, "v")
			return contourVersion, entry.Dependencies.Envoy, nil
		}
	}

	return "", "", fmt.Errorf("no supported versions found")
}

// setYAMLField updates a specific field in a YAML file.
func setYAMLField(filePath, fieldPath, newValue string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("failed to unmarshal yaml from %s: %w", filePath, err)
	}

	node := &root
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return fmt.Errorf("empty document")
		}
		node = node.Content[0]
	}

	parts := strings.Split(fieldPath, ".")
	if err := updateNode(node, parts, newValue); err != nil {
		return fmt.Errorf("failed to update field %s in %s: %w", fieldPath, filePath, err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	if err := encoder.Encode(&root); err != nil {
		return fmt.Errorf("failed to encode yaml to %s: %w", filePath, err)
	}

	return nil
}

// updateNode recursively updates the YAML node at the specified path.
func updateNode(node *yaml.Node, path []string, newValue string) error {
	if len(path) == 0 {
		node.Value = newValue
		return nil
	}

	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping node")
	}

	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == path[0] {
			return updateNode(node.Content[i+1], path[1:], newValue)
		}
	}

	return fmt.Errorf("field %s not found", path[0])
}

// getCurrentChartVersions reads the current chart and app version.
func getCurrentChartVersions(filePath string) (string, string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	type ChartYaml struct {
		Version    string `yaml:"version"`
		AppVersion string `yaml:"appVersion"`
	}

	var chart ChartYaml
	if err := yaml.Unmarshal(data, &chart); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal yaml from %s: %w", filePath, err)
	}

	return chart.Version, chart.AppVersion, nil
}

// nextMinorVersion calculates the next minor version given a version string.
func nextMinorVersion(version string) (string, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid version format: %s", version)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid minor version: %s", parts[1])
	}
	return fmt.Sprintf("%s.%d.0", parts[0], minor+1), nil
}
