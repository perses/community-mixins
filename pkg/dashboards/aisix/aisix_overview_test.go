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

package aisix

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildAISIXOverview(t *testing.T) {
	result := BuildAISIXOverview("default", "prometheus", "cluster")
	if result.Err() != nil {
		t.Fatalf("BuildAISIXOverview() returned error: %v", result.Err())
	}

	dashboardJSON, err := json.Marshal(result.Builder().Dashboard)
	if err != nil {
		t.Fatalf("failed to marshal dashboard: %v", err)
	}

	output := string(dashboardJSON)
	if got := strings.Count(output, `"kind":"Panel"`); got != 8 {
		t.Errorf("dashboard contains %d panels, want 8", got)
	}
	for _, expected := range []string{
		"aisix-ai-gateway-overview",
		"aisix_requests_total",
		"aisix_request_e2e_latency_seconds_bucket",
		"aisix_request_ttft_seconds_bucket",
		"aisix_llm_input_tokens_total",
		"aisix_cache_requests_total",
		"aisix_ratelimit_rejections_total",
		"aisix_otlp_fanout_failures_total",
		`job=~\"$job\"`,
		`instance=~\"$instance\"`,
		`cluster=\"$cluster\"`,
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("dashboard does not contain %q", expected)
		}
	}

	if strings.Contains(output, `* 100`) {
		t.Error("success ratio must stay in the 0-1 range used by the percent-decimal unit")
	}
}

func TestBuildAISIXOverviewWithoutClusterLabel(t *testing.T) {
	result := BuildAISIXOverview("default", "", "")
	if result.Err() != nil {
		t.Fatalf("BuildAISIXOverview() returned error: %v", result.Err())
	}

	dashboardJSON, err := json.Marshal(result.Builder().Dashboard)
	if err != nil {
		t.Fatalf("failed to marshal dashboard: %v", err)
	}

	if strings.Contains(string(dashboardJSON), `$cluster`) {
		t.Error("dashboard should not contain a cluster variable when clusterLabelName is empty")
	}
}
