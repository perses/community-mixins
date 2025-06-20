package opentelemetry

import (
	"github.com/perses/community-dashboards/pkg/promql"
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/prometheus/prometheus/promql/parser"
	"golang.org/x/exp/maps"
)

var OpentelemetryCommonPanelQueries = map[string]parser.Expr{
	"SpanRate_accepted": promql.SumByRate(
		"otelcol_receiver_accepted_spans_total",
		[]string{"job", "receiver"},
		label.New("job").EqualRegexp("$job"),
		label.New("receiver").EqualRegexp("$receiver"),
	),
	"SpanRate_refused": promql.SumByRate(
		"otelcol_receiver_refused_spans_total",
		[]string{"job", "receiver"},
		label.New("job").EqualRegexp("$job"),
		label.New("receiver").EqualRegexp("$receiver"),
	),
	"MeticPointsRate_accepted": promql.SumByRate(
		"otelcol_receiver_accepted_metric_points_total",
		[]string{"job", "receiver"},
		label.New("job").EqualRegexp("$job"),
		label.New("receiver").EqualRegexp("$receiver"),
	),
	"MeticPointsRate_refused": promql.SumByRate(
		"otelcol_receiver_refused_metric_points_total",
		[]string{"job", "receiver"},
		label.New("job").EqualRegexp("$job"),
		label.New("receiver").EqualRegexp("$receiver"),
	),
	"LogRecordsRate_accepted": promql.SumByRate(
		"otelcol_receiver_accepted_log_records_total",
		[]string{"job", "receiver"},
		label.New("job").EqualRegexp("$job"),
		label.New("receiver").EqualRegexp("$receiver"),
	),
	"LogRecordsRate_refused": promql.SumByRate(
		"otelcol_receiver_refused_log_records_total",
		[]string{"job", "receiver"},
		label.New("job").EqualRegexp("$job"),
		label.New("receiver").EqualRegexp("$receiver"),
	),
	"SpanProcessorRate_incoming_items": promql.SumByRate(
		"otelcol_processor_incoming_items_total",
		[]string{"job", "processor", "otel_signal"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
		label.New("otel_signal").EqualRegexp("traces"),
	),
	"SpanProcessorRate_outgoing_items": promql.SumByRate(
		"otelcol_processor_outgoing_items_total",
		[]string{"job", "processor", "otel_signal"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
		label.New("otel_signal").EqualRegexp("traces"),
	),
	"MetricProcessorRate_incoming_items": promql.SumByRate(
		"otelcol_processor_incoming_items_total",
		[]string{"job", "processor", "otel_signal"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
		label.New("otel_signal").EqualRegexp("metrics"),
	),
	"MetricProcessorRate_outgoing_items": promql.SumByRate(
		"otelcol_processor_outgoing_items_total",
		[]string{"job", "processor", "otel_signal"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
		label.New("otel_signal").EqualRegexp("metrics"),
	),
	"LogProcessorRate_incoming_items": promql.SumByRate(
		"otelcol_processor_incoming_items_total",
		[]string{"job", "processor", "otel_signal"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
		label.New("otel_signal").EqualRegexp("logs"),
	),
	"LogProcessorRate_outgoing_items": promql.SumByRate(
		"otelcol_processor_outgoing_items_total",
		[]string{"job", "processor", "otel_signal"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
		label.New("otel_signal").EqualRegexp("logs"),
	),
	"BatchProcessorRate_batch_send_size": promql.SumByIncrease(
		"otelcol_processor_batch_batch_send_size_bucket",
		[]string{"job", "processor", "le"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
	),
	"BatchProcessorRate_batch_send_size_count": promql.SumByRate(
		"otelcol_processor_batch_batch_send_size_count",
		[]string{"processor"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
	),
	"BatchProcessorRate_batch_send_size_sum": promql.SumByRate(
		"otelcol_processor_batch_batch_send_size_sum",
		[]string{"processor"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
	),
	"BatchProcessorRate_batch_size_trigger_send": promql.SumByRate(
		"otelcol_processor_batch_batch_size_trigger_send_total",
		[]string{"processor"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
	),
	"BatchProcessorRate_batch_timeout_trigger_send": promql.SumByRate(
		"otelcol_processor_batch_timeout_trigger_send_total",
		[]string{"processor"},
		label.New("job").EqualRegexp("$job"),
		label.New("processor").EqualRegexp("$processor"),
	),
	"ExporterRate_sent_spans": promql.SumByRate(
		"otelcol_exporter_sent_spans_total",
		[]string{"exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_enqueue_failed_spans": promql.SumByRate(
		"otelcol_exporter_enqueue_failed_spans_total",
		[]string{"exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_send_failed_spans": promql.SumByRate(
		"otelcol_exporter_send_failed_spans_total",
		[]string{"exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_sent_metrics": promql.SumByRate(
		"otelcol_exporter_sent_metric_points_total",
		[]string{"exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_enqueue_failed_metrics": promql.SumByRate(
		"otelcol_exporter_enqueue_failed_metric_points_total",
		[]string{"exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_send_failed_metrics": promql.SumByRate(
		"otelcol_exporter_send_failed_metric_points_total",
		[]string{"exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_sent_logs": promql.SumByRate(
		"otelcol_exporter_sent_log_records_total",
		[]string{"exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_enqueue_failed_logs": promql.SumByRate(
		"otelcol_exporter_enqueue_failed_log_records_total",
		[]string{"exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_send_failed_logs": promql.SumByRate(
		"otelcol_exporter_send_failed_log_records_total",
		[]string{"exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_queue_size": promql.MaxBy(
		"otelcol_exporter_queue_size",
		[]string{"exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_queue_capacity": promql.MinBy(
		"otelcol_exporter_queue_capacity",
		[]string{"exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("exporter").EqualRegexp("$exporter"),
	),
	"ExporterRate_queue_utilization": promqlbuilder.Mul(
		&parser.NumberLiteral{Val: 100},
		promqlbuilder.Div(
			promql.MaxBy(
				"otelcol_exporter_queue_size",
				[]string{"exporter"},
				label.New("job").EqualRegexp("$job"),
				label.New("exporter").EqualRegexp("$exporter"),
			),
			promql.MaxBy(
				"otelcol_exporter_queue_capacity",
				[]string{"exporter"},
				label.New("job").EqualRegexp("$job"),
				label.New("exporter").EqualRegexp("$exporter"),
			),
		),
	),
}

// OverrideOpentelemetryPanelQueries overrides the OpentelemetryCommonPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
// The convention of naming followed, is to use Panel function name (with _suffix, in case panel has multiple queries)
func OverrideOpentelemetryPanelQueries(queries map[string]parser.Expr) {
	maps.Copy(OpentelemetryCommonPanelQueries, queries)
}
