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

var NODE_EXPORTER_JOB_VALUE = "node"

// GetNodeExporterJobValue returns the current job label value for node-exporter dashboards.
func GetNodeExporterJobValue() string {
	return NODE_EXPORTER_JOB_VALUE
}

// SetNodeExporterJobValue sets the job label value for node-exporter dashboards globally.
// WARNING: Ensure you only set this to a value that represents a Node Exporter specifically.
func SetNodeExporterJobValue(labelValue string) {
	NODE_EXPORTER_JOB_VALUE = labelValue
}
