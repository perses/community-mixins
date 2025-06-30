package blackbox

import (
	"github.com/perses/community-dashboards/pkg/promql"
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/matrix"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/promql/parser"
	"golang.org/x/exp/maps"
)

var BlackboxCommonPanelQueries = map[string]parser.Expr{
	"BlackboxProbeSucess": promql.MaxBy(
		"probe_success",
		[]string{"instance"},
		label.New("job").EqualRegexp("$job"),
	),
	"BlackboxProbeSucessCount": promqlbuilder.Count(
		vector.New(
			vector.WithMetricName("probe_success"),
			vector.WithLabelMatchers(
				label.New("job").EqualRegexp("$job"),
			),
		),
	),
	"BlackboxProbeSucessPercent": promqlbuilder.Div(
		promqlbuilder.Or(
			promqlbuilder.Count(
				promqlbuilder.Eql(
					vector.New(
						vector.WithMetricName("probe_success"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp("$job"),
						),
					),
					&parser.NumberLiteral{Val: 1},
				),
			),
			promqlbuilder.Vector(0),
		),
		promqlbuilder.Count(
			vector.New(
				vector.WithMetricName("probe_success"),
				vector.WithLabelMatchers(
					label.New("job").EqualRegexp("$job"),
				),
			),
		),
	),
	"BlackboxProbeHTTPSSL": promqlbuilder.Div(
		promqlbuilder.Count(
			promqlbuilder.Eql(
				vector.New(
					vector.WithMetricName("probe_http_ssl"),
					vector.WithLabelMatchers(
						label.New("job").EqualRegexp("$job"),
					),
				),
				&parser.NumberLiteral{Val: 1},
			),
		),
		promqlbuilder.Count(
			vector.New(
				vector.WithMetricName("probe_http_version"),
				vector.WithLabelMatchers(
					label.New("job").EqualRegexp("$job"),
				),
			),
		),
	),
	"BlackboxAvgProbeDuration": promqlbuilder.Avg(
		vector.New(
			vector.WithMetricName("probe_duration_seconds"),
			vector.WithLabelMatchers(
				label.New("job").EqualRegexp("$job"),
			),
		),
	),
	"BlackboxProbeUptime": promql.MaxBy(
		"probe_success",
		[]string{"instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"BlackboxProbeUptimeMonthly": promqlbuilder.AvgOverTime(
		matrix.New(
			vector.New(
				vector.WithMetricName("probe_success"),
				vector.WithLabelMatchers(
					label.New("job").EqualRegexp("$job"),
					label.New("instance").EqualRegexp("$instance"),
				),
			),
			matrix.WithRangeAsString("30d"),
		),
	),
	"BlackboxProbeHttpDuration": promqlbuilder.Sum(
		promql.AvgBy(
			"probe_http_duration_seconds",
			[]string{"phase", "instance"},
			label.New("job").EqualRegexp("$job"),
			label.New("instance").EqualRegexp("$instance"),
		),
	).By("instance"),
	"BlackboxAvgProbeDurationSeconds": promql.AvgBy(
		"probe_duration_seconds",
		[]string{"instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"BlackboxProbeHttpPhases": promql.AvgBy(
		"probe_http_duration_seconds",
		[]string{"phase"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"BlackboxProbeIcmpPhases": promql.AvgBy(
		"probe_icmp_duration_seconds",
		[]string{"phase"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"BlackboxProbeStatusCode": promql.MaxBy(
		"probe_http_status_code",
		[]string{"instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"BlackboxProbeTLSVersion": promql.MaxBy(
		"probe_tls_version_info",
		[]string{"instance", "version"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"BlackboxProbeSSLExpiry": promql.MinBy(
		"probe_ssl_earliest_cert_expiry",
		[]string{"instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"BlackboxProbeRedirects": promql.MaxBy(
		"probe_http_redirects",
		[]string{"instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"BlackboxProbeHTTPVersion": promql.MaxBy(
		"probe_http_version",
		[]string{"instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"BlackboxProbeAverageDuration": promql.AvgBy(
		"probe_duration_seconds",
		[]string{"instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"BlackboxProbeAverageDNSLookupPerInstance": promql.AvgBy(
		"probe_dns_lookup_time_seconds",
		[]string{"instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
}

// OverrideBlackboxPanelQueries overrides the BlackboxCommonPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
// The convention of naming followed, is to use Panel function name (with _suffix, in case panel has multiple queries)
func OverrideBlackboxPanelQueries(queries map[string]parser.Expr) {
	maps.Copy(BlackboxCommonPanelQueries, queries)
}
