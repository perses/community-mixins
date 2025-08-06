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

func withZtunnelTraffic(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Traffic",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.ZtunnelBytesTransmitted(datasource, labelMatcher),
		panels.ZtunnelConnections(datasource, labelMatcher),
	)
}

func withZtunnelResources(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Resource Usage",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.ZtunnelCPUUsage(datasource, labelMatcher),
		panels.ZtunnelMemoryUsage(datasource, labelMatcher),
	)
}

func withZtunnelOperations(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Operations",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.ZtunnelDNSRequest(datasource, labelMatcher),
		panels.ZtunnelWorkloadManager(datasource, labelMatcher),
	)
}

func withZtunnelXDS(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("XDS",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.ZtunnelXDSConnections(datasource, labelMatcher),
		panels.ZtunnelXDSPushes(datasource, labelMatcher),
	)
}

func withZtunnelStatus(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Status",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(6),
		panels.ZtunnelVersions(datasource, labelMatcher),
	)
}

func BuildIstioZtunnel(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("istio-ztunnel",
			dashboard.ProjectName(project),
			dashboard.Name("Istio / Ztunnel"),
			dashboard.AddVariable("cluster",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("cluster",
						labelValuesVar.Matchers("ztunnel_build_info"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("cluster"),
					listVar.AllowMultiple(true),
				),
			),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "ztunnel_build_info"),
			withZtunnelTraffic(datasource, clusterLabelMatcher),
			withZtunnelResources(datasource, clusterLabelMatcher),
			withZtunnelOperations(datasource, clusterLabelMatcher),
			withZtunnelXDS(datasource, clusterLabelMatcher),
			withZtunnelStatus(datasource, clusterLabelMatcher),
		),
	).Component("istio")
}
