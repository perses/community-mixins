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

func withRolloutProgres(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rollout progress",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.RolloutProgress(datasource, labelMatcher),
	)
}

func withWritesStats(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Tempo Writes stats",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(8),
		panels.TempoWrites2xx(datasource, labelMatcher),
		panels.TempoWrites4xx(datasource, labelMatcher),
		panels.TempoWrites5xx(datasource, labelMatcher),
		panels.TempoWritesLatency(datasource, labelMatcher),
	)
}

func withReadsStats(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Tempo Reads stats",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(8),
		panels.TempoReads2xx(datasource, labelMatcher),
		panels.TempoReads4xx(datasource, labelMatcher),
		panels.TempoReads5xx(datasource, labelMatcher),
		panels.TempoReadsLatency(datasource, labelMatcher),
	)
}

func withTempoPods(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Tempo Pods",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.TempoUnhealthyPods(datasource, labelMatcher),
		panels.TempoPodsCount(datasource, labelMatcher),
		panels.TempoLatencyHistory(datasource, labelMatcher),
	)
}

func BuildTempoRolloutOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("tempo-rollout-progress-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Tempo / Rollout progress"),
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
			withRolloutProgres(datasource, clusterLabelMatcher),
			withWritesStats(datasource, clusterLabelMatcher),
			withReadsStats(datasource, clusterLabelMatcher),
			withTempoPods(datasource, clusterLabelMatcher),
		),
	).Component("tempo")
}
