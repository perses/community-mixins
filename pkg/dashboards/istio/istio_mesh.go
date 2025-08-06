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

func withMeshOverview(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Mesh Overview",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(6),
		panels.GlobalRequestVolume(datasource, labelMatcher),
		panels.GlobalSuccessRate(datasource, labelMatcher),
		panels.Global4xxRate(datasource, labelMatcher),
		panels.Global5xxRate(datasource, labelMatcher),
	)
}

func withMeshWorkloads(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("HTTP/gRPC Workloads",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(12),
		panels.HTTPGRPCWorkloads(datasource, labelMatcher),
	)
}

func withMeshTCPServices(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("TCP Services",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.TCPServices(datasource, labelMatcher),
	)
}

func BuildIstioMesh(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("istio-mesh",
			dashboard.ProjectName(project),
			dashboard.Name("Istio / Mesh"),
			dashboard.AddVariable("cluster",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("cluster",
						labelValuesVar.Matchers("istio_requests_total"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("cluster"),
					listVar.AllowMultiple(true),
				),
			),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "istio_requests_total"),
			withMeshOverview(datasource, clusterLabelMatcher),
			withMeshWorkloads(datasource, clusterLabelMatcher),
			withMeshTCPServices(datasource, clusterLabelMatcher),
		),
	).Component("istio")
}
