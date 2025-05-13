package perses

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panelsGostats "github.com/perses/community-dashboards/pkg/panels/gostats"
	"github.com/perses/community-dashboards/pkg/panels/perses"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listvariable "github.com/perses/perses/go-sdk/variable/list-variable"
	labelvalues "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
)

func BuildPersesOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("perses-overview",
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
			withPersesPlugins(datasource, clusterLabelMatcher),
		),
	).Component("perses")
}

func withPersesOverviewStatsGroup(datasource string, clusterLabelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Perses Stats", panelgroup.PanelsPerLine(1),
		perses.StatsTable(datasource, clusterLabelMatcher))
}

func withPersesAPiRequestGroup(datasource string, clusterLabelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("API Requests", panelgroup.PanelsPerLine(2),
		perses.HTTPRequestsRatePanel(datasource, clusterLabelMatcher),
		perses.HTTPErrorPercentagePanel(datasource, clusterLabelMatcher),
		perses.HTTPRequestsLatencyPanel(datasource, clusterLabelMatcher),
	)
}

func withPersesResources(datasource string, clusterLabelMatcher promql.LabelMatcher) dashboard.Option {
	labelMatchersToUse := []promql.LabelMatcher{
		promql.InstanceVar,
		promql.JobVar,
	}
	labelMatchersToUse = append(labelMatchersToUse, clusterLabelMatcher)

	return dashboard.AddPanelGroup("Resource Usage",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(10),
		panelsGostats.MemoryUsage(datasource, "pod", labelMatchersToUse...),
		panelsGostats.CPUUsage(datasource, "pod", labelMatchersToUse...),
		panelsGostats.Goroutines(datasource, "pod", labelMatchersToUse...),
		panelsGostats.GarbageCollectionPauseTimeQuantiles(datasource, "pod", labelMatchersToUse...),
		perses.FileDescriptors(datasource, clusterLabelMatcher),
	)
}

func withPersesPlugins(datasource string, clusterLabelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Plugins Usage", panelgroup.PanelsPerLine(1), panelgroup.PanelHeight(8),
		perses.PluginSchemaLoadAttempts(datasource, clusterLabelMatcher))
}
