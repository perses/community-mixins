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

var NODE_EXPORTER_LABEL_VALUE = "node"

// GetNodeExporterLabelValue returns the current job label value for node-exporter dashboards.
func GetNodeExporterLabelValue() string {
	return NODE_EXPORTER_LABEL_VALUE
}

// SetNodeExporterLabelValue sets the job label value for node-exporter dashboards globally.
// WARNING: Ensure you only set this to a value that represents a Node Exporter specifically.
func SetNodeExporterLabelValue(labelValue string) {
	NODE_EXPORTER_LABEL_VALUE = labelValue
}
