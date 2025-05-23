package tempo

import (
	panels "github.com/perses/community-dashboards/pkg/panels/tempo"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"

	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
)

func withTenantInfo(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Tenant Info",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.TenantInfo(datasource, labelMatcher),
	)
}

func withTenantIngestion(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Ingestion",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.TenantDistributorBytes(datasource, labelMatcher),
		panels.TenantDistributorSpan(datasource, labelMatcher),
		panels.TenantLiveTraces(datasource, labelMatcher),
	)
}

func withTenantReads(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Reads",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.TenantQueriesID(datasource, labelMatcher),
		panels.TenantQueriesSearch(datasource, labelMatcher),
	)
}

func withTenantStorage(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.TenantBlockslistLength(datasource, labelMatcher),
		panels.TenantOutstandingCompactions(datasource, labelMatcher),
	)
}

func withTenantMetricGenerator(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Metrics Generator",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.TenantMetricGeneratorBytes(datasource, labelMatcher),
		panels.TenantMetricGeneratorActiveSeries(datasource, labelMatcher),
	)
}

func BuildTempoTenantOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("thanos-tenant-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Tempo / Tenant"),
			dashboard.AddVariable("cluster",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("cluster",
						labelValuesVar.Matchers("tempo_build_info"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("cluster"),
					listVar.AllowMultiple(true),
				),
			),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "tempo_build_info"),
			dashboard.AddVariable("namespace",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("namespace",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"tempo_build_info",
								[]promql.LabelMatcher{clusterLabelMatcher, {Name: "cluster", Type: "=", Value: "$cluster"}},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("namespace"),
				),
			),
			dashboard.AddVariable("tenant",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("tenant",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"tempodb_blocklist_length",
								[]promql.LabelMatcher{
									{Name: "cluster", Type: "=", Value: "$cluster"},
									{Name: "job", Type: "=~", Value: "($namespace)/compactor"},
								}),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("tenant"),
				),
			),
			withTenantInfo(datasource, clusterLabelMatcher),
			withTenantIngestion(datasource, clusterLabelMatcher),
			withTenantReads(datasource, clusterLabelMatcher),
			withTenantStorage(datasource, clusterLabelMatcher),
			withTenantMetricGenerator(datasource, clusterLabelMatcher),
		),
	).Component("tempo")
}
