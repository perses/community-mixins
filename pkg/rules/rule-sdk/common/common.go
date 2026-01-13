// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import "maps"

// BuildAnnotations creates annotations map with conditional dashboard and runbook labels
func BuildAnnotations(dashboardURL, runbookURL, runbookFragment, description, summary string) map[string]string {
	annotations := map[string]string{
		"description": description,
		"summary":     summary,
	}

	if dashboardURL != "" {
		annotations["dashboard"] = dashboardURL
	}
	if runbookURL != "" {
		annotations["runbook"] = runbookURL + runbookFragment
	}

	return annotations
}

// MergeMaps merges two maps
func MergeMaps(mapsToMerge ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range mapsToMerge {
		maps.Copy(result, m)
	}
	return result
}
