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

func withPerformanceDataPlane(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Data Plane Performance",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.DataPlaneMemory(datasource, labelMatcher),
		panels.DataPlaneCPU(datasource, labelMatcher),
	)
}

func withPerformanceControlPlane(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Control Plane Performance",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.ProxyPushTime(datasource, labelMatcher),
		panels.ProxyQueueSize(datasource, labelMatcher),
	)
}

func withPerformanceOperations(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Operations Performance",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.ConfigurationValidation(datasource, labelMatcher),
		panels.SidecarInjection(datasource, labelMatcher),
	)
}

func BuildIstioPerformance(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("istio-performance",
			dashboard.ProjectName(project),
			dashboard.Name("Istio / Performance"),
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
			withPerformanceDataPlane(datasource, clusterLabelMatcher),
			withPerformanceControlPlane(datasource, clusterLabelMatcher),
			withPerformanceOperations(datasource, clusterLabelMatcher),
		),
	).Component("istio")
}
