package thanos

import (
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"

	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/thanos"
	"github.com/perses/community-dashboards/pkg/promql"
)

func withThanosQueryGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Instant Query",
		panelgroup.PanelsPerLine(3),
		panels.InstantQueryRequestRate(datasource, labelMatcher),
		panels.InstantQueryRequestErrors(datasource, labelMatcher),
		panels.InstantQueryRequestDuration(datasource, labelMatcher),
	)
}

func withThanosRangeQueryGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Range Query",
		panelgroup.PanelsPerLine(3),
		panels.RangeQueryRequestRate(datasource, labelMatcher),
		panels.RangeQueryRequestErrors(datasource, labelMatcher),
		panels.RangeQueryRequestDuration(datasource, labelMatcher),
	)
}

func withThanosQueryConcurrencyGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Query Available Concurrency",
		panelgroup.PanelsPerLine(1),
		panels.QueryConcurrency(datasource, labelMatcher),
	)
}

func withThanosDNSLookupGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("DNS Lookups",
		panelgroup.PanelsPerLine(2),
		panels.DNSLookups(datasource, labelMatcher),
		panels.DNSLookupsErrors(datasource, labelMatcher),
	)
}

func BuildThanosQueryOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("thanos-query-overview",
		dashboard.ProjectName(project),
		dashboard.Name("Thanos / Query / Overview"),
		dashboard.AddVariable("namespace",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("namespace",
					labelValuesVar.Matchers("thanos_status"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("namespace"),
			),
		),
		dashboard.AddVariable("job",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("job",
					labelValuesVar.Matchers("thanos_build_info{container=\"thanos-query\"}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("job"),
				listVar.AllowMultiple(true),
			),
		),
		dashboards.AddClusterVariable(datasource, clusterLabelName, "thanos_build_info"),
		withThanosQueryGroup(datasource, clusterLabelMatcher),
		withThanosRangeQueryGroup(datasource, clusterLabelMatcher),
		withThanosReadGRPCUnaryGroup(datasource, clusterLabelMatcher),
		withThanosReadGRPCStreamGroup(datasource, clusterLabelMatcher),
		withThanosQueryConcurrencyGroup(datasource, clusterLabelMatcher),
		withThanosDNSLookupGroup(datasource, clusterLabelMatcher),
	)
}
