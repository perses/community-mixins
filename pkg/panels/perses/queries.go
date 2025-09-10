package perses

import (
	"maps"

	"github.com/perses/community-dashboards/pkg/promql"
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

var PersesCommonPanelQueries = map[string]parser.Expr{
	"PersesStatsTable": promqlbuilder.Count(
		vector.New(
			vector.WithMetricName("perses_build_info"),
			vector.WithLabelMatchers(
				label.New("job").EqualRegexp("$job"),
				label.New("instance").EqualRegexp("$instance"),
			),
		),
	).By("job", "instance", "version", "namespace", "pod"),
	"PersesHTTPLatency": promqlbuilder.Div(
		promql.SumByRate(
			"perses_http_request_duration_second_sum",
			[]string{"handler", "method"},
			label.New("job").EqualRegexp("$job"),
			label.New("instance").EqualRegexp("$instance"),
		),
		promql.SumByRate(
			"perses_http_request_duration_second_count",
			[]string{"handler", "method"},
			label.New("job").EqualRegexp("$job"),
			label.New("instance").EqualRegexp("$instance"),
		),
	),
	"PersesHTTPRequestRate": promql.SumByRate(
		"perses_http_request_total",
		[]string{"handler", "code"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"PersesHTTPErrorsPercentage": promqlbuilder.Mul(
		promql.IgnoringGroupLeft(
			promql.ErrorCaseRatio(
				"perses_http_request_total",
				[]string{"handler", "code"},
				[]*labels.Matcher{
					label.New("job").EqualRegexp("$job"),
					label.New("instance").EqualRegexp("$instance"),
					label.New("code").EqualRegexp("4..|5.."),
				},
				"perses_http_request_total",
				[]string{"handler"},
				[]*labels.Matcher{
					label.New("job").EqualRegexp("$job"),
					label.New("instance").EqualRegexp("$instance"),
				},
			),
			[]string{"code"},
		),
		&parser.NumberLiteral{Val: 100},
	),
	"PersesFileDescriptorsOpenFDS": vector.New(
		vector.WithMetricName("process_open_fds"),
		vector.WithLabelMatchers(
			label.New("job").EqualRegexp("$job"),
			label.New("instance").EqualRegexp("$instance"),
		),
	),
	"PersesFileDescriptorsMaxFDS": vector.New(
		vector.WithMetricName("process_max_fds"),
		vector.WithLabelMatchers(
			label.New("job").EqualRegexp("$job"),
			label.New("instance").EqualRegexp("$instance"),
		),
	),
	"PersesPluginSchemaLoadAttempts": vector.New(
		vector.WithMetricName("perses_plugin_schemas_load_attempts"),
		vector.WithLabelMatchers(
			label.New("job").EqualRegexp("$job"),
			label.New("instance").EqualRegexp("$instance"),
		),
	),
}

// OverridePersesPanelQueries overrides the PersesCommonPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
// The convention of naming followed, is to use Panel function name (with _suffix, in case panel has multiple queries)
func OverridePersesPanelQueries(queries map[string]parser.Expr) {
	maps.Copy(PersesCommonPanelQueries, queries)
}
