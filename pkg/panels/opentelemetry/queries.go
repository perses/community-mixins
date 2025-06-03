package opentelemetry

import (
	"github.com/perses/community-dashboards/pkg/promql"
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
}

// OverrideOpentelemetryPanelQueries overrides the OpentelemetryCommonPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
// The convention of naming followed, is to use Panel function name (with _suffix, in case panel has multiple queries)
func OverrideOpentelemetryPanelQueries(queries map[string]parser.Expr) {
	maps.Copy(OpentelemetryCommonPanelQueries, queries)
}
