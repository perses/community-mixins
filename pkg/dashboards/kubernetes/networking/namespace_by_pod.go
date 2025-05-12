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

func withNamespaceCurrentRateBytesGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Current Rate of Bytes",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesCurrentRateOfBytesReceived("namespace-pod-networking", datasource, labelMatcher),
		panels.KubernetesCurrentRateOfBytesTransmitted("namespace-pod-networking", datasource, labelMatcher),
	)
}

func withNamespaceNetworkUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Network Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.NamespaceCurrentNetworkUsage(datasource, labelMatcher),
	)
}

func withNamespaceBandwidthGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Bandwidth",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceiveBandwidth("namespace-pod-networking", datasource, labelMatcher),
		panels.KubernetesTransmitBandwidth("namespace-pod-networking", datasource, labelMatcher),
	)
}

func withNamespaceRateOfPacketsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPackets("namespace-pod-networking", datasource, labelMatcher),
		panels.KubernetesTransmittedPackets("namespace-pod-networking", datasource, labelMatcher),
	)
}

func withNamespaceRateOfPacketsDroppedGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets Dropped",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPacketsDropped("namespace-pod-networking", datasource, labelMatcher),
		panels.KubernetesTransmittedPacketsDropped("namespace-pod-networking", datasource, labelMatcher),
	)
}

func BuildKubernetesNamespaceByPodOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("kubernetes-namespace-networking-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Kubernetes / Networking / Namespace (Pods)"),
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
			withNamespaceCurrentRateBytesGroup(datasource, clusterLabelMatcher),
			withNamespaceNetworkUsageGroup(datasource, clusterLabelMatcher),
			withNamespaceBandwidthGroup(datasource, clusterLabelMatcher),
			withNamespaceRateOfPacketsGroup(datasource, clusterLabelMatcher),
			withNamespaceRateOfPacketsDroppedGroup(datasource, clusterLabelMatcher),
		),
	).Component("kubernetes")
}
