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

import "testing"

func TestGetNodeExporterLabelValueDefault(t *testing.T) {
	NODE_EXPORTER_LABEL_VALUE = "node"
	got := GetNodeExporterLabelValue()
	if got != "node" {
		t.Errorf("GetNodeExporterLabelValue() = %q, want %q", got, "node")
	}
}

func TestSetNodeExporterLabelValue(t *testing.T) {
	defer func() { NODE_EXPORTER_LABEL_VALUE = "node" }()

	tests := []struct {
		name  string
		value string
	}{
		{"kube-prometheus-stack default", "node-exporter"},
		{"custom value", "my-node-exporter"},
		{"reset to default", "node"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetNodeExporterLabelValue(tt.value)
			got := GetNodeExporterLabelValue()
			if got != tt.value {
				t.Errorf("after SetNodeExporterLabelValue(%q), GetNodeExporterLabelValue() = %q", tt.value, got)
			}
		})
	}
}
