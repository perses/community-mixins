package alertmanager

import (
	"github.com/perses/community-dashboards/pkg/promql"
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/promql/parser"
)

var AlertmanagerCommonPanelQueries = map[string]parser.Expr{
	"Alerts": promqlbuilder.Sum(
		vector.New(
			vector.WithMetricName("alertmanager_alerts"),
			vector.WithLabelMatchers(
				label.New("job").EqualRegexp("$job"),
			),
		),
	).By("instance"),
	"AlertsReceiveRate_received": promql.SumByRate(
		"alertmanager_alerts_received_total",
		[]string{"job", "instance"},
		label.New("job").EqualRegexp("$job"),
	),
	"AlertsReceiveRate_invalid": promql.SumByRate(
		"alertmanager_alerts_invalid_total",
		[]string{"job", "instance"},
		label.New("job").EqualRegexp("$job"),
	),
	"NotificationsSendRate_total": promql.SumByRate(
		"alertmanager_notifications_total",
		[]string{"integration", "instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("integration").EqualRegexp("$integration"),
	),
	"NotificationsSendRate_failed": promql.SumByRate(
		"alertmanager_notifications_failed_total",
		[]string{"integration", "instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("integration").EqualRegexp("$integration"),
	),
	"NotificationDuration_p99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"alertmanager_notification_latency_seconds_bucket",
			[]string{"le", "integration", "instance"},
			label.New("job").EqualRegexp("$job"),
			label.New("integration").EqualRegexp("$integration"),
		),
	),
	"NotificationDuration_p50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"alertmanager_notification_latency_seconds_bucket",
			[]string{"le", "integration", "instance"},
			label.New("job").EqualRegexp("$job"),
			label.New("integration").EqualRegexp("$integration"),
		),
	),
	"NotificationDuration_avg": promqlbuilder.Div(
		promql.SumByRate(
			"alertmanager_notification_latency_seconds_sum",
			[]string{"integration", "instance"},
			label.New("job").EqualRegexp("$job"),
			label.New("integration").EqualRegexp("$integration"),
		),
		promql.SumByRate(
			"alertmanager_notification_latency_seconds_count",
			[]string{"integration", "instance"},
			label.New("job").EqualRegexp("$job"),
			label.New("integration").EqualRegexp("$integration"),
		),
	),
}

// OverrideAlertmanagerPanelQueries overrides the AlertmanagerCommonPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
// The convention of naming followed, is to use Panel function name (with _suffix, in case panel has multiple queries)
func OverrideAlertmanagerPanelQueries(queries map[string]parser.Expr) {
	for k, v := range queries {
		AlertmanagerCommonPanelQueries[k] = v
	}
}
