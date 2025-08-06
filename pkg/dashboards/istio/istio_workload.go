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

func withWorkloadIncoming(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Incoming Traffic",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.IncomingRequestVolume(datasource, labelMatcher),
		panels.IncomingSuccessRate(datasource, labelMatcher),
		panels.IncomingRequestDuration(datasource, labelMatcher),
	)
}

func withWorkloadOutgoing(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Outgoing Traffic",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.OutgoingRequestVolume(datasource, labelMatcher),
		panels.OutgoingSuccessRate(datasource, labelMatcher),
	)
}

func withWorkloadTCP(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("TCP Traffic",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.TCPBytesReceived(datasource, labelMatcher),
		panels.TCPBytesSent(datasource, labelMatcher),
	)
}

func BuildIstioWorkload(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("istio-workload",
			dashboard.ProjectName(project),
			dashboard.Name("Istio / Workload"),
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
			dashboard.AddVariable("namespace",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("destination_workload_namespace",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"istio_requests_total",
								[]promql.LabelMatcher{clusterLabelMatcher, {Name: "cluster", Type: "=", Value: "$cluster"}},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Namespace"),
				),
			),
			dashboard.AddVariable("workload",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("destination_workload",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"istio_requests_total",
								[]promql.LabelMatcher{
									clusterLabelMatcher,
									{Name: "cluster", Type: "=", Value: "$cluster"},
									{Name: "destination_workload_namespace", Type: "=", Value: "$namespace"},
								},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Workload"),
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
									{Name: "destination_workload", Type: "=", Value: "$workload"},
									{Name: "destination_workload_namespace", Type: "=", Value: "$namespace"},
								},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Source Workload Namespace"),
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
									{Name: "destination_workload", Type: "=", Value: "$workload"},
									{Name: "destination_workload_namespace", Type: "=", Value: "$namespace"},
									{Name: "source_workload_namespace", Type: "=~", Value: "$srcns"},
								},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Source Workload"),
					listVar.AllowAllValue(true),
				),
			),
			dashboard.AddVariable("dstsvc",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("destination_service",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"istio_requests_total",
								[]promql.LabelMatcher{
									clusterLabelMatcher,
									{Name: "cluster", Type: "=", Value: "$cluster"},
									{Name: "source_workload", Type: "=", Value: "$workload"},
									{Name: "source_workload_namespace", Type: "=", Value: "$namespace"},
								},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Destination Service"),
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
									{Name: "source_workload", Type: "=", Value: "$workload"},
									{Name: "source_workload_namespace", Type: "=", Value: "$namespace"},
									{Name: "destination_service", Type: "=~", Value: "$dstsvc"},
								},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Destination Workload Namespace"),
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
									{Name: "source_workload", Type: "=", Value: "$workload"},
									{Name: "source_workload_namespace", Type: "=", Value: "$namespace"},
									{Name: "destination_service", Type: "=~", Value: "$dstsvc"},
									{Name: "destination_workload_namespace", Type: "=~", Value: "$dstns"},
								},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Destination Workload"),
					listVar.AllowAllValue(true),
				),
			),
			withWorkloadIncoming(datasource, clusterLabelMatcher),
			withWorkloadOutgoing(datasource, clusterLabelMatcher),
			withWorkloadTCP(datasource, clusterLabelMatcher),
		),
	).Component("istio")
}
