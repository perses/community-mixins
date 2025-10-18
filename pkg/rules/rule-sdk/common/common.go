package common

import "maps"

// BuildAnnotations creates annotations map with conditional dashboard and runbook labels
func BuildAnnotations(dashboardURL, runbookURL, runbookFragment, description, message, summary string) map[string]string {
	annotations := map[string]string{
		"description": description,
		"message":     message,
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
