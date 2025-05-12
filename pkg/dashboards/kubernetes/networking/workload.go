package networking

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/kubernetes"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"

	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
)

func withWorkloadCurrentRateOfBytesGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Current Rate of Bytes",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesCurrentRateOfBytesReceived("workload-networking", datasource, labelMatcher),
		panels.KubernetesCurrentRateOfBytesTransmitted("workload-networking", datasource, labelMatcher),
	)
}

func withWorkloadAverageRateOfBytesGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Average Rate of Bytes",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesAverageRateOfBytesReceived("workload-networking", datasource, labelMatcher),
		panels.KubernetesAverageRateOfBytesTransmitted("workload-networking", datasource, labelMatcher),
	)
}

func withWorkloadBandwidthGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Bandwidth",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceiveBandwidth("workload-networking", datasource, labelMatcher),
		panels.KubernetesTransmitBandwidth("workload-networking", datasource, labelMatcher),
	)
}

func withWorkloadRateOfPacketsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPackets("workload-networking", datasource, labelMatcher),
		panels.KubernetesTransmittedPackets("workload-networking", datasource, labelMatcher),
	)
}

func withWorkloadRateOfPacketsDroppedGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets Dropped",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPacketsDropped("workload-networking", datasource, labelMatcher),
		panels.KubernetesTransmittedPacketsDropped("workload-networking", datasource, labelMatcher),
	)
}

func BuildKubernetesWorkloadOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("kubernetes-workload-networking-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Kubernetes / Networking / Workload"),
			dashboard.AddVariable("cluster",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("cluster",
						labelValuesVar.Matchers("up{"+panels.GetKubeletMatcher()+", metrics_path=\"/metrics/cadvisor\"}"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("cluster"),
				),
			),
			dashboard.AddVariable("namespace",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("namespace",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"kube_namespace_status_phase{"+panels.GetKubeStateMetricsMatcher()+"}",
								[]promql.LabelMatcher{{Name: "cluster", Type: "=", Value: "$cluster"}},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("namespace"),
				),
			),
			dashboard.AddVariable("type",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("workload_type",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"namespace_workload_pod:kube_pod_owner:relabel",
								[]promql.LabelMatcher{
									{Name: "cluster", Type: "=", Value: "$cluster"},
									{Name: "namespace", Type: "=", Value: "$namespace"},
								},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("workload_type"),
				),
			),
			dashboard.AddVariable("workload",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("workload",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"namespace_workload_pod:kube_pod_owner:relabel",
								[]promql.LabelMatcher{
									{Name: "cluster", Type: "=", Value: "$cluster"},
									{Name: "namespace", Type: "=", Value: "$namespace"},
									{Name: "workload_type", Type: "=", Value: "$type"},
								},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("workload"),
				),
			),
			withWorkloadCurrentRateOfBytesGroup(datasource, clusterLabelMatcher),
			withWorkloadAverageRateOfBytesGroup(datasource, clusterLabelMatcher),
			withWorkloadBandwidthGroup(datasource, clusterLabelMatcher),
			withWorkloadRateOfPacketsGroup(datasource, clusterLabelMatcher),
			withWorkloadRateOfPacketsDroppedGroup(datasource, clusterLabelMatcher),
		),
	).Component("kubernetes")
}
