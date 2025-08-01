package networking

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/kubernetes"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"

	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
)

func withWorkloadNamespaceCurrentRateBytesGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Current Rate of Bytes",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesCurrentRateOfBytesReceived("namespace-workload-networking", datasource, labelMatcher),
		panels.KubernetesCurrentRateOfBytesTransmitted("namespace-workload-networking", datasource, labelMatcher),
	)
}

func withWorkloadNamespaceNetworkStatusGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Current Status",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.WorkloadNamespaceCurrentNetworkStatus(datasource, labelMatcher),
	)
}

func withWorkloadNamespaceBandwidthGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Bandwidth",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceiveBandwidth("namespace-workload-networking", datasource, labelMatcher),
		panels.KubernetesTransmitBandwidth("namespace-workload-networking", datasource, labelMatcher),
	)
}

func withWorkloadNamespaceAvgContainerBandwidthGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Average Container Bandwidth",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesAvgContainerBandwidthReceived("namespace-workload-networking", datasource, labelMatcher),
		panels.KubernetesAvgContainerBandwidthTransmitted("namespace-workload-networking", datasource, labelMatcher),
	)
}

func withWorkloadNamespaceRateOfPacketsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPackets("namespace-workload-networking", datasource, labelMatcher),
		panels.KubernetesTransmittedPackets("namespace-workload-networking", datasource, labelMatcher),
	)
}

func withWorkloadNamespaceRateOfPacketsDroppedGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets Dropped",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPacketsDropped("namespace-workload-networking", datasource, labelMatcher),
		panels.KubernetesTransmittedPacketsDropped("namespace-workload-networking", datasource, labelMatcher),
	)
}

func BuildKubernetesNamespaceByWorkloadOverview(project string, datasource string, clusterLabelName string, variableOverrides []dashboard.Option) dashboards.DashboardResult {
	defaultVars := []dashboard.Option{
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
							"container_network_receive_packets_total",
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
							"namespace_workload_pod:kube_pod_owner:relabel{workload=~\".+\"}",
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
	}

	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)

	vars := defaultVars
	if len(variableOverrides) > 0 {
		vars = variableOverrides
	}
	options := append([]dashboard.Option{
		dashboard.ProjectName(project),
		dashboard.Name("Kubernetes / Networking / Namespace (Workloads)"),
	}, vars...)
	options = append(options,
		withWorkloadNamespaceCurrentRateBytesGroup(datasource, clusterLabelMatcher),
		withWorkloadNamespaceNetworkStatusGroup(datasource, clusterLabelMatcher),
		withWorkloadNamespaceBandwidthGroup(datasource, clusterLabelMatcher),
		withWorkloadNamespaceAvgContainerBandwidthGroup(datasource, clusterLabelMatcher),
		withWorkloadNamespaceRateOfPacketsGroup(datasource, clusterLabelMatcher),
		withWorkloadNamespaceRateOfPacketsDroppedGroup(datasource, clusterLabelMatcher),
	)
	return dashboards.NewDashboardResult(
		dashboard.New("kubernetes-workload-ns-networking-overview", options...),
	).Component("kubernetes")
}
