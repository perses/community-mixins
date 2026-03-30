// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nodeexporter

import (
	"encoding/json"
	"strings"
	"testing"

	panels "github.com/perses/community-mixins/pkg/panels/node_exporter"
)

func TestBuildNodeExporterNodes_DefaultJobLabel(t *testing.T) {
	panels.SetNodeExporterLabelValue("node")
	defer panels.SetNodeExporterLabelValue("node")

	result := BuildNodeExporterNodes("default", "", "")
	if result.Err() != nil {
		t.Fatalf("BuildNodeExporterNodes() returned error: %v", result.Err())
	}

	dashboardJSON, err := json.Marshal(result.Builder().Dashboard)
	if err != nil {
		t.Fatalf("failed to marshal dashboard: %v", err)
	}

	output := string(dashboardJSON)
	if !strings.Contains(output, `job=\"node\"`) && !strings.Contains(output, `job="node"`) && !strings.Contains(output, `job='node'`) {
		t.Error("expected dashboard to contain job=\"node\" with default config")
	}
}

func TestBuildNodeExporterNodes_CustomJobLabel(t *testing.T) {
	panels.SetNodeExporterLabelValue("node-exporter")
	defer panels.SetNodeExporterLabelValue("node")

	result := BuildNodeExporterNodes("default", "", "")
	if result.Err() != nil {
		t.Fatalf("BuildNodeExporterNodes() returned error: %v", result.Err())
	}

	dashboardJSON, err := json.Marshal(result.Builder().Dashboard)
	if err != nil {
		t.Fatalf("failed to marshal dashboard: %v", err)
	}

	output := string(dashboardJSON)
	if !strings.Contains(output, "node-exporter") {
		t.Error("expected dashboard to contain 'node-exporter' after setting custom job label")
	}

	// The variable matcher string should use the custom value, not the default
	if strings.Contains(output, `job='node'`) {
		t.Error("dashboard should not contain job='node' after setting job label to 'node-exporter'")
	}
}

func TestBuildNodeExporterClusterUseMethod_CustomJobLabel(t *testing.T) {
	panels.SetNodeExporterLabelValue("node-exporter")
	defer panels.SetNodeExporterLabelValue("node")

	result := BuildNodeExporterClusterUseMethod("default", "", "")
	if result.Err() != nil {
		t.Fatalf("BuildNodeExporterClusterUseMethod() returned error: %v", result.Err())
	}

	dashboardJSON, err := json.Marshal(result.Builder().Dashboard)
	if err != nil {
		t.Fatalf("failed to marshal dashboard: %v", err)
	}

	output := string(dashboardJSON)
	if !strings.Contains(output, "node-exporter") {
		t.Error("expected dashboard to contain 'node-exporter' after setting custom job label")
	}

	if strings.Contains(output, `job='node'`) {
		t.Error("dashboard should not contain job='node' after setting job label to 'node-exporter'")
	}
}
