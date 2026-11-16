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

package opentelemetry

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
)

func buildCollectorDashboardJSON(t *testing.T, clusterLabelName string) string {
	t.Helper()
	result := BuildOpenTelemetryCollector("default", "prometheus", clusterLabelName)
	if result.Err() != nil {
		t.Fatalf("BuildOpenTelemetryCollector() returned error: %v", result.Err())
	}
	dashboardJSON, err := json.Marshal(result.Builder().Dashboard)
	if err != nil {
		t.Fatalf("failed to marshal dashboard: %v", err)
	}
	var decoded any
	if err := json.Unmarshal(dashboardJSON, &decoded); err != nil {
		t.Fatalf("failed to unmarshal dashboard: %v", err)
	}
	// Collect every string value with whitespace collapsed, so that pretty-printed
	// PromQL can be matched on a single line.
	var values []string
	var walk func(v any)
	walk = func(v any) {
		switch v := v.(type) {
		case map[string]any:
			for _, child := range v {
				walk(child)
			}
		case []any:
			for _, child := range v {
				walk(child)
			}
		case string:
			values = append(values, strings.Join(strings.Fields(v), " "))
		}
	}
	walk(decoded)
	return strings.Join(values, "\n") + "\n" + string(dashboardJSON)
}

func TestBuildOpenTelemetryCollector(t *testing.T) {
	output := buildCollectorDashboardJSON(t, "cluster")

	if got := strings.Count(output, `"kind":"Panel"`); got != 15 {
		t.Errorf("dashboard contains %d panels, want 15", got)
	}

	// The collector exposes counters without the _total suffix on its Prometheus endpoint,
	// while the OTLP to Prometheus translation adds it: counters must match both names.
	if hardcoded := regexp.MustCompile(`otelcol_[a-z_]+_total[{\[]`).FindAllString(output, -1); len(hardcoded) > 0 {
		t.Errorf("dashboard selects counters by their _total name only: %v", hardcoded)
	}

	for _, expected := range []string{
		`otelcol_receiver_accepted_spans(_total)?`,
		`otelcol_receiver_refused_log_records(_total)?`,
		`otelcol_processor_incoming_items(_total)?`,
		`otelcol_processor_batch_timeout_trigger_send(_total)?`,
		`otelcol_exporter_send_failed_metric_points(_total)?`,
		`otelcol_exporter_enqueue_failed_spans(_total)?`,
		// Prometheus 3 stores the otel.signal attribute with its dotted name.
		`label_replace(`,
		`"otel_signal", "$1", "otel.signal", "(.+)"`,
		`and on (otel_signal) label_replace(vector(1), "otel_signal", "traces", "", "")`,
		// Exporter queues are per signal, so data_type must be kept.
		`max by (job, exporter, data_type)`,
		`min by (job, exporter, data_type)`,
		`Metric Points Rate`,
		`Batch sent due to a size trigger: {{processor}}`,
		`{__name__=~'otelcol_process_uptime.*', job='$job'}`,
		`cluster="$cluster"`,
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("dashboard does not contain %q", expected)
		}
	}

	for _, unexpected := range []string{
		`Metic Points Rate`,
		`otel_signal=~"traces"`,
		`up{job='$job'}`,
	} {
		if strings.Contains(output, unexpected) {
			t.Errorf("dashboard still contains %q", unexpected)
		}
	}
}

func TestBuildOpenTelemetryCollectorWithoutClusterLabel(t *testing.T) {
	output := buildCollectorDashboardJSON(t, "")
	if strings.Contains(output, `$cluster`) {
		t.Error("dashboard should not contain a cluster variable when clusterLabelName is empty")
	}
}
