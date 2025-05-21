package thanos

import (
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"

	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/thanos"
	"github.com/perses/community-dashboards/pkg/promql"
)

func withThanosQueryInstantQueryGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Instant Query",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.InstantQueryRequestRate(datasource, labelMatcher),
		panels.InstantQueryRequestErrors(datasource, labelMatcher),
		panels.InstantQueryRequestDurations(datasource, labelMatcher),
	)
}

func withThanosQueryRangeQueryGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Range Query",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.RangeQueryRequestRate(datasource, labelMatcher),
		panels.RangeQueryRequestErrors(datasource, labelMatcher),
		panels.RangeQueryRequestDurations(datasource, labelMatcher),
	)
}

func withThanosQueryConcurrencyGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Query Available Concurrency",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.QueryConcurrency(datasource, labelMatcher),
	)
}

func withThanosQueryDNSLookupGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("DNS Lookups",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.DNSLookups(datasource, labelMatcher),
		panels.DNSLookupsErrors(datasource, labelMatcher),
	)
}

func BuildThanosQueryOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	clusterLabelMatcherV2 := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("thanos-query-overview",
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
			withThanosQueryInstantQueryGroup(datasource, clusterLabelMatcher),
			withThanosQueryRangeQueryGroup(datasource, clusterLabelMatcher),
			withThanosReadGRPCUnaryGroup(datasource, clusterLabelMatcherV2),
			withThanosReadGRPCStreamGroup(datasource, clusterLabelMatcherV2),
			withThanosQueryConcurrencyGroup(datasource, clusterLabelMatcher),
			withThanosQueryDNSLookupGroup(datasource, clusterLabelMatcher),
			withThanosResourcesGroup(datasource, clusterLabelMatcher),
		),
	).Component("thanos")
}
