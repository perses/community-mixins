package perses

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/panels/perses"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	labelvalues "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listvariable "github.com/perses/perses/go-sdk/variable/list-variable"
)

func BuildPersesOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("perses-overview",
		dashboard.ProjectName(project),
		dashboard.Name("Perses / Overview"),
		dashboard.AddVariable("job",
			listvariable.List(
				labelvalues.PrometheusLabelValues("job",
					labelvalues.Matchers("perses_build_info{}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listvariable.DisplayName("job"),
			),
		),
		dashboard.AddVariable("instance",
			listvariable.List(
				labelvalues.PrometheusLabelValues("instance",
					labelvalues.Matchers(
						promql.SetLabelMatchers(
							"perses_build_info",
							[]promql.LabelMatcher{clusterLabelMatcher, {Name: "job", Type: "=", Value: "$job"}},
						),
					),
					dashboards.AddVariableDatasource(datasource),
				),
				listvariable.DisplayName("instance"),
			),
		),
		withPersesOverviewStatsGroup(datasource, clusterLabelMatcher),
		withPersesAPiRequestGroup(datasource, clusterLabelMatcher),
		withPersesResources(datasource, clusterLabelMatcher),
	)
}

func withPersesOverviewStatsGroup(datasource string, clusterLabelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Perses Stats", panelgroup.PanelsPerLine(1), perses.PersesStatsTable(datasource, clusterLabelMatcher))
}

func withPersesAPiRequestGroup(datasource string, clusterLabelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("API Requests", panelgroup.PanelsPerLine(2), perses.PersesHTTPRequestsLatency(datasource, clusterLabelMatcher), perses.PersesTotalHTTPRequests(datasource, clusterLabelMatcher))
}

func withPersesResources(datasource string, clusterLabelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Resource Usage", panelgroup.PanelsPerLine(2),
		perses.PersesMemoryUsage(datasource, clusterLabelMatcher),
		perses.PersesCPUUsage(datasource, clusterLabelMatcher))
}
