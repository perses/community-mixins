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

package alertmanager

import (
	"maps"

	"github.com/perses/community-mixins/pkg/promql"
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
	maps.Copy(AlertmanagerCommonPanelQueries, queries)
}
