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
	"maps"

	"github.com/perses/community-mixins/pkg/promql"
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/matrix"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

var OpentelemetryCommonPanelQueries = map[string]parser.Expr{
	"SpanRate_accepted": counterSumByRate(
		"otelcol_receiver_accepted_spans",
		[]string{"job", "receiver"},
		label.New("job").EqualRegexp("$job"),
		label.New("receiver").EqualRegexp("$receiver"),
	),
	"SpanRate_refused": counterSumByRate(
		"otelcol_receiver_refused_spans",
		[]string{"job", "receiver"},
		label.New("job").EqualRegexp("$job"),
		label.New("receiver").EqualRegexp("$receiver"),
	),
	"MeticPointsRate_accepted": counterSumByRate(
		"otelcol_receiver_accepted_metric_points",
		[]string{"job", "receiver"},
		label.New("job").EqualRegexp("$job"),
		label.New("receiver").EqualRegexp("$receiver"),
	),
	"MeticPointsRate_refused": counterSumByRate(
		"otelcol_receiver_refused_metric_points",
		[]string{"job", "receiver"},
		label.New("job").EqualRegexp("$job"),
		label.New("receiver").EqualRegexp("$receiver"),
	),
	"LogRecordsRate_accepted": counterSumByRate(
		"otelcol_receiver_accepted_log_records",
		[]string{"job", "receiver"},
		label.New("job").EqualRegexp("$job"),
		label.New("receiver").EqualRegexp("$receiver"),
	),
	"LogRecordsRate_refused": counterSumByRate(
		"otelcol_receiver_refused_log_records",
		[]string{"job", "receiver"},
		label.New("job").EqualRegexp("$job"),
		label.New("receiver").EqualRegexp("$receiver"),
	),
	"SpanProcessorRate_incoming_items":   processorItemsRate("otelcol_processor_incoming_items", "traces"),
	"SpanProcessorRate_outgoing_items":   processorItemsRate("otelcol_processor_outgoing_items", "traces"),
	"MetricProcessorRate_incoming_items": processorItemsRate("otelcol_processor_incoming_items", "metrics"),
	"MetricProcessorRate_outgoing_items": processorItemsRate("otelcol_processor_outgoing_items", "metrics"),
	"LogProcessorRate_incoming_items":    processorItemsRate("otelcol_processor_incoming_items", "logs"),
	"LogProcessorRate_outgoing_items":    processorItemsRate("otelcol_processor_outgoing_items", "logs"),
	"BatchProcessorRate_batch_send_size": promql.SumByIncrease(
		"otelcol_processor_batch_batch_send_size_bucket",
		[]string{"job", "processor", "le"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
	),
	"BatchProcessorRate_batch_send_size_count": promql.SumByRate(
		"otelcol_processor_batch_batch_send_size_count",
		[]string{"job", "processor"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
	),
	"BatchProcessorRate_batch_send_size_sum": promql.SumByRate(
		"otelcol_processor_batch_batch_send_size_sum",
		[]string{"job", "processor"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
	),
	"BatchProcessorRate_batch_size_trigger_send": counterSumByRate(
		"otelcol_processor_batch_batch_size_trigger_send",
		[]string{"job", "processor"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
	),
	"BatchProcessorRate_batch_timeout_trigger_send": counterSumByRate(
		"otelcol_processor_batch_timeout_trigger_send",
		[]string{"job", "processor"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
	),
	"ExporterRate_sent_spans": counterSumByRate(
		"otelcol_exporter_sent_spans",
		[]string{"job", "exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_enqueue_failed_spans": counterSumByRate(
		"otelcol_exporter_enqueue_failed_spans",
		[]string{"job", "exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_send_failed_spans": counterSumByRate(
		"otelcol_exporter_send_failed_spans",
		[]string{"job", "exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_sent_metrics": counterSumByRate(
		"otelcol_exporter_sent_metric_points",
		[]string{"job", "exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_enqueue_failed_metrics": counterSumByRate(
		"otelcol_exporter_enqueue_failed_metric_points",
		[]string{"job", "exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_send_failed_metrics": counterSumByRate(
		"otelcol_exporter_send_failed_metric_points",
		[]string{"job", "exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_sent_logs": counterSumByRate(
		"otelcol_exporter_sent_log_records",
		[]string{"job", "exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_enqueue_failed_logs": counterSumByRate(
		"otelcol_exporter_enqueue_failed_log_records",
		[]string{"job", "exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_send_failed_logs": counterSumByRate(
		"otelcol_exporter_send_failed_log_records",
		[]string{"job", "exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_queue_size": promql.MaxBy(
		"otelcol_exporter_queue_size",
		[]string{"job", "exporter", "data_type"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_queue_capacity": promql.MinBy(
		"otelcol_exporter_queue_capacity",
		[]string{"job", "exporter", "data_type"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_queue_utilization": promqlbuilder.Mul(
		&parser.NumberLiteral{Val: 100},
		promqlbuilder.Div(
			promql.MaxBy(
				"otelcol_exporter_queue_size",
				[]string{"job", "exporter", "data_type"},
				label.New("job").EqualRegexp("$job"),
				label.New("exporter").EqualRegexp("$exporter"),
			),
			promql.MaxBy(
				"otelcol_exporter_queue_capacity",
				[]string{"job", "exporter", "data_type"},
				label.New("job").EqualRegexp("$job"),
				label.New("exporter").EqualRegexp("$exporter"),
			),
		),
	),
}

// counterName matches an otelcol counter under both naming conventions found in Prometheus.
// The collector's own Prometheus endpoint exposes counters without the _total suffix
// (without_type_suffix defaults to true), while the OTLP to Prometheus translation adds it.
func counterName(name string) *labels.Matcher {
	return label.New("__name__").EqualRegexp(name + "(_total)?")
}

func counterSumByRate(name string, byLabels []string, labelMatchers ...*labels.Matcher) parser.Expr {
	return promql.SumByRate("", byLabels, append([]*labels.Matcher{counterName(name)}, labelMatchers...)...)
}

// processorItemsRate returns the rate of a processor items counter for one signal.
// Prometheus 3 keeps the original otel.signal label name when it scrapes the collector with
// UTF-8 names, whereas Prometheus 2, underscore escaping and the OTLP translation store it as
// otel_signal. label_replace copies otel.signal into otel_signal so both layouts can be filtered.
func processorItemsRate(name string, signal string) parser.Expr {
	return promqlbuilder.And(
		promqlbuilder.Sum(
			promqlbuilder.LabelReplace(
				promqlbuilder.Rate(
					matrix.New(
						vector.New(
							vector.WithLabelMatchers(
								counterName(name),
								label.New("job").EqualRegexp("$job"),
								label.New("processor").EqualRegexp("$processor"),
							),
						),
						matrix.WithRangeAsVariable("$__rate_interval"),
					),
				),
				"otel_signal", "$1", "otel.signal", "(.+)",
			),
		).By("job", "processor", "otel_signal"),
		promqlbuilder.LabelReplace(promqlbuilder.Vector(1), "otel_signal", signal, "", ""),
	).On("otel_signal")
}

// OverrideOpentelemetryPanelQueries overrides the OpentelemetryCommonPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
// The convention of naming followed, is to use Panel function name (with _suffix, in case panel has multiple queries)
func OverrideOpentelemetryPanelQueries(queries map[string]parser.Expr) {
	maps.Copy(OpentelemetryCommonPanelQueries, queries)
}
