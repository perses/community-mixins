package compute_resources

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/kubernetes"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"

	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
)

func withNamespaceStatsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Namespace Stats",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(4),
		panels.KubernetesCPUUtilizationStat("namespace-requests", datasource, labelMatcher),
		panels.KubernetesCPUUtilizationStat("namespace-limits", datasource, labelMatcher),
		panels.KubernetesMemoryUtilizationStat("namespace-requests", datasource, labelMatcher),
		panels.KubernetesMemoryUtilizationStat("namespace-limits", datasource, labelMatcher),
	)
}

func withNamespaceCPUUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.KubernetesCPUUsage("namespace-pod", datasource, labelMatcher),
	)
}

func withNamespaceCPUUsageQuotaGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU Usage Quota",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.NamespaceCPUUsageQuota(datasource, labelMatcher),
	)
}

func withNamespaceMemoryUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.KubernetesMemoryUsage("namespace-pod", datasource, labelMatcher),
	)
}

func withNamespaceMemoryUsageQuotaGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory Usage Quota",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.NamespaceMemoryUsageQuota(datasource, labelMatcher),
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
		panels.KubernetesReceiveBandwidth("namespace-pod", datasource, labelMatcher),
		panels.KubernetesTransmitBandwidth("namespace-pod", datasource, labelMatcher),
	)
}

func withNamespaceRateOfPacketsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPackets("namespace-pod", datasource, labelMatcher),
		panels.KubernetesTransmittedPackets("namespace-pod", datasource, labelMatcher),
	)
}

func withNamespaceRateOfPacketsDroppedGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets Dropped",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPacketsDropped("namespace-pod", datasource, labelMatcher),
		panels.KubernetesTransmittedPacketsDropped("namespace-pod", datasource, labelMatcher),
	)
}

func withNamespaceStorageIOGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage IO",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesIOPS("namespace-pod", datasource, labelMatcher),
		panels.KubernetesThroughput("namespace-pod", datasource, labelMatcher),
	)
}

func withNamespaceCurrentStorageIOGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage IO - Distribution",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.NamespaceCurrentStorageIO(datasource, labelMatcher),
	)
}

func BuildKubernetesNamespaceOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("kubernetes-namespace-resources-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Kubernetes / Compute Resources / Namespace (Pods)"),
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
			withNamespaceStatsGroup(datasource, clusterLabelMatcher),
			withNamespaceCPUUsageGroup(datasource, clusterLabelMatcher),
			withNamespaceCPUUsageQuotaGroup(datasource, clusterLabelMatcher),
			withNamespaceMemoryUsageGroup(datasource, clusterLabelMatcher),
			withNamespaceMemoryUsageQuotaGroup(datasource, clusterLabelMatcher),
			withNamespaceNetworkUsageGroup(datasource, clusterLabelMatcher),
			withNamespaceBandwidthGroup(datasource, clusterLabelMatcher),
			withNamespaceRateOfPacketsGroup(datasource, clusterLabelMatcher),
			withNamespaceRateOfPacketsDroppedGroup(datasource, clusterLabelMatcher),
			withNamespaceStorageIOGroup(datasource, clusterLabelMatcher),
			withNamespaceCurrentStorageIOGroup(datasource, clusterLabelMatcher),
		),
	).Component("kubernetes")
}
