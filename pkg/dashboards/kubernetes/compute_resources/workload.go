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

func withWorkloadCPUUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.KubernetesCPUUsage("workload", datasource, labelMatcher),
	)
}

func withWorkloadCPUUsageQuotaGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU Usage Quota",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.WorkloadCPUUsageQuota(datasource, labelMatcher),
	)
}

func withWorkloadMemoryUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.KubernetesMemoryUsage("workload", datasource, labelMatcher),
	)
}

func withWorkloadMemoryUsageQuotaGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory Usage Quota",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.WorkloadMemoryUsageQuota(datasource, labelMatcher),
	)
}

func withWorkloadNetworkUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Network Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.WorkloadCurrentNetworkUsage(datasource, labelMatcher),
	)
}

func withWorkloadBandwidthGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Bandwidth",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceiveBandwidth("workload", datasource, labelMatcher),
		panels.KubernetesTransmitBandwidth("workload", datasource, labelMatcher),
	)
}

func withWorkloadAvgContainerBandwidthGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Average Container Bandwidth",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesAvgContainerBandwidthReceived("workload", datasource, labelMatcher),
		panels.KubernetesAvgContainerBandwidthTransmitted("workload", datasource, labelMatcher),
	)
}

func withWorkloadRateOfPacketsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPackets("workload", datasource, labelMatcher),
		panels.KubernetesTransmittedPackets("workload", datasource, labelMatcher),
	)
}

func withWorkloadRateOfPacketsDroppedGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets Dropped",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPacketsDropped("workload", datasource, labelMatcher),
		panels.KubernetesTransmittedPacketsDropped("workload", datasource, labelMatcher),
	)
}

func BuildKubernetesWorkloadOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("kubernetes-workload-resources-overview",
		dashboard.ProjectName(project),
		dashboard.Name("Kubernetes / Compute Resources / Workload"),
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
		withWorkloadCPUUsageGroup(datasource, clusterLabelMatcher),
		withWorkloadCPUUsageQuotaGroup(datasource, clusterLabelMatcher),
		withWorkloadMemoryUsageGroup(datasource, clusterLabelMatcher),
		withWorkloadMemoryUsageQuotaGroup(datasource, clusterLabelMatcher),
		withWorkloadNetworkUsageGroup(datasource, clusterLabelMatcher),
		withWorkloadBandwidthGroup(datasource, clusterLabelMatcher),
		withWorkloadAvgContainerBandwidthGroup(datasource, clusterLabelMatcher),
		withWorkloadRateOfPacketsGroup(datasource, clusterLabelMatcher),
		withWorkloadRateOfPacketsDroppedGroup(datasource, clusterLabelMatcher),
	)
}
