package kubernetes

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/kubernetes"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"

	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
)

func withClusterStatsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Cluster Stats",
		panelgroup.PanelsPerLine(6),
		panelgroup.PanelHeight(4),
		panels.CPUUtilizationStat(datasource, labelMatcher),
		panels.CPURequestsCommitmentStat(datasource, labelMatcher),
		panels.CPULimitsCommitmentStat(datasource, labelMatcher),
		panels.MemoryUtilizationStat(datasource, labelMatcher),
		panels.MemoryRequestsCommitmentStat(datasource, labelMatcher),
		panels.MemoryLimitsCommitmentStat(datasource, labelMatcher),
	)
}

func withClusterCPUUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.KubernetesCPUUsage("cluster", datasource, labelMatcher),
	)
}

func withClusterCPUUsageQuotaGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU Usage Quota",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.ClusterCPUUsageQuota(datasource, labelMatcher),
	)
}

func withClusterMemoryUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.KubernetesMemoryUsage("cluster", datasource, labelMatcher),
	)
}

func withClusterMemoryUsageQuotaGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory Usage Quota",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.ClusterMemoryUsageQuota(datasource, labelMatcher),
	)
}

func withClusterNetworkUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Network Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.ClusterCurrentNetworkUsage(datasource, labelMatcher),
	)
}

func withClusterBandwidthGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Bandwidth",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceiveBandwidth("cluster", datasource, labelMatcher),
		panels.KubernetesTransmitBandwidth("cluster", datasource, labelMatcher),
	)
}

func withClusterAvgBandwidthGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Average Container Bandwidth",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesAvgContainerBandwidthReceived("cluster", datasource, labelMatcher),
		panels.KubernetesAvgContainerBandwidthTransmitted("cluster", datasource, labelMatcher),
	)
}

func withClusterRateOfPacketsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPackets("cluster", datasource, labelMatcher),
		panels.KubernetesTransmittedPackets("cluster", datasource, labelMatcher),
	)
}

func withClusterRateOfPacketsDroppedGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets Dropped",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPacketsDropped("cluster", datasource, labelMatcher),
		panels.KubernetesTransmittedPacketsDropped("cluster", datasource, labelMatcher),
	)
}

func withClusterStorageIOGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage IO",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesIOPS("cluster", datasource, labelMatcher),
		panels.KubernetesThroughput("cluster", datasource, labelMatcher),
	)
}

func withClusterCurrentStorageIOGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage IO - Distribution",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.ClusterCurrentStorageIO(datasource, labelMatcher),
	)
}

func BuildKubernetesClusterOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("kubernetes-cluster-overview",
		dashboard.ProjectName(project),
		dashboard.Name("Kubernetes / Compute Resources / Cluster"),
		dashboard.AddVariable("cluster",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("cluster",
					labelValuesVar.Matchers("up{job=\"kubelet\", metrics_path=\"/metrics/cadvisor\"}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("cluster"),
			),
		),
		withClusterStatsGroup(datasource, clusterLabelMatcher),
		withClusterCPUUsageGroup(datasource, clusterLabelMatcher),
		withClusterCPUUsageQuotaGroup(datasource, clusterLabelMatcher),
		withClusterMemoryUsageGroup(datasource, clusterLabelMatcher),
		withClusterMemoryUsageQuotaGroup(datasource, clusterLabelMatcher),
		withClusterNetworkUsageGroup(datasource, clusterLabelMatcher),
		withClusterBandwidthGroup(datasource, clusterLabelMatcher),
		withClusterAvgBandwidthGroup(datasource, clusterLabelMatcher),
		withClusterRateOfPacketsGroup(datasource, clusterLabelMatcher),
		withClusterRateOfPacketsDroppedGroup(datasource, clusterLabelMatcher),
		withClusterStorageIOGroup(datasource, clusterLabelMatcher),
		withClusterCurrentStorageIOGroup(datasource, clusterLabelMatcher),
	)
}
