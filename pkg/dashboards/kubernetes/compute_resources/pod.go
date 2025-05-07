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

func withPodCPUUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.KubernetesCPUUsage("pod", datasource, labelMatcher),
	)
}

func withPodCPUThrottlingGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU Throttling",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.PodCPUThrottling(datasource, labelMatcher),
	)
}

func withPodCPUUsageQuotaGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU Usage Quota",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.PodCPUUsageQuota(datasource, labelMatcher),
	)
}

func withPodMemoryUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.KubernetesMemoryUsage("pod", datasource, labelMatcher),
	)
}

func withPodMemoryUsageQuotaGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory Usage Quota",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.PodMemoryUsageQuota(datasource, labelMatcher),
	)
}

func withPodBandwidthGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Bandwidth",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceiveBandwidth("pod", datasource, labelMatcher),
		panels.KubernetesTransmitBandwidth("pod", datasource, labelMatcher),
	)
}

func withPodRateOfPacketsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPackets("pod", datasource, labelMatcher),
		panels.KubernetesTransmittedPackets("pod", datasource, labelMatcher),
	)
}

func withPodRateOfPacketsDroppedGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets Dropped",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPacketsDropped("pod", datasource, labelMatcher),
		panels.KubernetesTransmittedPacketsDropped("pod", datasource, labelMatcher),
	)
}

func withPodStorageIOGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage IO",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesIOPS("pod", datasource, labelMatcher),
		panels.KubernetesThroughput("pod", datasource, labelMatcher),
	)
}

func withPodStorageIOContainerGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage IO - Container",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesIOPS("pod-container", datasource, labelMatcher),
		panels.KubernetesThroughput("pod-container", datasource, labelMatcher),
	)
}

func withPodCurrentStorageIOGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage IO - Distribution",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.PodCurrentStorageIO(datasource, labelMatcher),
	)
}

func BuildKubernetesPodOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("kubernetes-pod-overview",
		dashboard.ProjectName(project),
		dashboard.Name("Kubernetes / Compute Resources / Pod"),
		dashboard.AddVariable("cluster",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("cluster",
					labelValuesVar.Matchers("up{job=\"kubelet\", metrics_path=\"/metrics/cadvisor\"}"),
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
							"kube_namespace_status_phase{job=\"kube-state-metrics\"}",
							[]promql.LabelMatcher{{Name: "cluster", Type: "=", Value: "$cluster"}},
						),
					),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("namespace"),
			),
		),
		dashboard.AddVariable("pod",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("pod",
					labelValuesVar.Matchers(
						promql.SetLabelMatchers(
							"kube_pod_info{job=\"kube-state-metrics\"}",
							[]promql.LabelMatcher{{Name: "cluster", Type: "=", Value: "$cluster"}, {Name: "namespace", Type: "=", Value: "$namespace"}},
						),
					),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("pod"),
			),
		),
		withPodCPUUsageGroup(datasource, clusterLabelMatcher),
		withPodCPUThrottlingGroup(datasource, clusterLabelMatcher),
		withPodCPUUsageQuotaGroup(datasource, clusterLabelMatcher),
		withPodMemoryUsageGroup(datasource, clusterLabelMatcher),
		withPodMemoryUsageQuotaGroup(datasource, clusterLabelMatcher),
		withPodBandwidthGroup(datasource, clusterLabelMatcher),
		withPodRateOfPacketsGroup(datasource, clusterLabelMatcher),
		withPodRateOfPacketsDroppedGroup(datasource, clusterLabelMatcher),
		withPodStorageIOGroup(datasource, clusterLabelMatcher),
		withPodStorageIOContainerGroup(datasource, clusterLabelMatcher),
		withPodCurrentStorageIOGroup(datasource, clusterLabelMatcher),
	)
}
