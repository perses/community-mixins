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

func withServiceClientTraffic(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Client Traffic",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.ClientRequestVolume(datasource, labelMatcher),
		panels.ClientSuccessRate(datasource, labelMatcher),
		panels.ClientRequestDuration(datasource, labelMatcher),
	)
}

func withServiceServerTraffic(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Server Traffic",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.ServerRequestVolume(datasource, labelMatcher),
	)
}

func withServiceTCP(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("TCP Traffic",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.ServiceTCPBytesReceived(datasource, labelMatcher),
		panels.ServiceTCPBytesSent(datasource, labelMatcher),
	)
}

func BuildIstioService(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("istio-service-dashboard",
			dashboard.ProjectName(project),
			dashboard.Name("Istio Service Dashboard"),
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
			dashboard.AddVariable("qrep",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("reporter",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"istio_requests_total",
								[]promql.LabelMatcher{clusterLabelMatcher, {Name: "cluster", Type: "=", Value: "$cluster"}},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Reporter"),
					listVar.DefaultValue("destination"),
				),
			),
			dashboard.AddVariable("service",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("destination_service",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"istio_requests_total",
								[]promql.LabelMatcher{clusterLabelMatcher, {Name: "cluster", Type: "=", Value: "$cluster"}},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Service"),
				),
			),
			dashboard.AddVariable("srcns",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("source_workload_namespace",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"istio_requests_total",
								[]promql.LabelMatcher{
									clusterLabelMatcher,
									{Name: "cluster", Type: "=", Value: "$cluster"},
									{Name: "destination_service", Type: "=", Value: "$service"},
								},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Client Workload Namespace"),
					listVar.AllowAllValue(true),
				),
			),
			dashboard.AddVariable("srcwl",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("source_workload",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"istio_requests_total",
								[]promql.LabelMatcher{
									clusterLabelMatcher,
									{Name: "cluster", Type: "=", Value: "$cluster"},
									{Name: "destination_service", Type: "=", Value: "$service"},
									{Name: "source_workload_namespace", Type: "=~", Value: "$srcns"},
								},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Client Workload"),
					listVar.AllowAllValue(true),
				),
			),
			dashboard.AddVariable("dstns",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("destination_workload_namespace",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"istio_requests_total",
								[]promql.LabelMatcher{
									clusterLabelMatcher,
									{Name: "cluster", Type: "=", Value: "$cluster"},
									{Name: "destination_service", Type: "=", Value: "$service"},
								},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Service Workload Namespace"),
					listVar.AllowAllValue(true),
				),
			),
			dashboard.AddVariable("dstwl",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("destination_workload",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"istio_requests_total",
								[]promql.LabelMatcher{
									clusterLabelMatcher,
									{Name: "cluster", Type: "=", Value: "$cluster"},
									{Name: "destination_service", Type: "=", Value: "$service"},
									{Name: "destination_workload_namespace", Type: "=~", Value: "$dstns"},
								},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Service Workload"),
					listVar.AllowAllValue(true),
				),
			),
			withServiceClientTraffic(datasource, clusterLabelMatcher),
			withServiceServerTraffic(datasource, clusterLabelMatcher),
			withServiceTCP(datasource, clusterLabelMatcher),
		),
	).Component("istio")
}
