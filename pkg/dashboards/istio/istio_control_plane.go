package istio

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/istio"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
)

func withControlPlaneMetrics(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Control Plane Metrics",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.PushSize(datasource, labelMatcher),
		panels.PushTime(datasource, labelMatcher),
		panels.Connections(datasource, labelMatcher),
	)
}

func withControlPlaneResources(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Resource Usage",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.CPUUsage(datasource, labelMatcher),
		panels.MemoryUsage(datasource, labelMatcher),
		panels.MemoryAllocations(datasource, labelMatcher),
	)
}

func withControlPlaneEvents(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Events",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.Events(datasource, labelMatcher),
		panels.Goroutines(datasource, labelMatcher),
	)
}

func withControlPlaneOperations(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Operations",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.Injection(datasource, labelMatcher),
		panels.Validation(datasource, labelMatcher),
	)
}

func withControlPlaneStatus(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Status",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.PilotVersions(datasource, labelMatcher),
		panels.PushErrors(datasource, labelMatcher),
		panels.XDSPushes(datasource, labelMatcher),
	)
}

func BuildIstioControlPlane(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("istio-control-plane",
			dashboard.ProjectName(project),
			dashboard.Name("Istio / Control Plane"),
			dashboard.AddVariable("cluster",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("cluster",
						labelValuesVar.Matchers("pilot_build_info"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("cluster"),
					listVar.AllowMultiple(true),
				),
			),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "pilot_build_info"),
			dashboard.AddVariable("namespace",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("namespace",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"pilot_build_info",
								[]promql.LabelMatcher{clusterLabelMatcher, {Name: "cluster", Type: "=", Value: "$cluster"}},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("namespace"),
				),
			),
			withControlPlaneMetrics(datasource, clusterLabelMatcher),
			withControlPlaneResources(datasource, clusterLabelMatcher),
			withControlPlaneEvents(datasource, clusterLabelMatcher),
			withControlPlaneOperations(datasource, clusterLabelMatcher),
			withControlPlaneStatus(datasource, clusterLabelMatcher),
		),
	).Component("istio")
}
