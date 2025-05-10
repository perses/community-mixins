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

func withNodeCPUUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.KubernetesCPUUsage("node", datasource, labelMatcher),
	)
}

func withNodeCPUQuotaGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU Quota",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(12),
		panels.CPUUsageQuota(datasource, labelMatcher),
	)
}

func withNodeMemoryUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory Usage with Cache",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.KubernetesMemoryUsage("node-with-cache", datasource, labelMatcher),
	)
}

func withNodeMemoryUsageWithoutCacheGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory Usage without Cache",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.KubernetesMemoryUsage("node-without-cache", datasource, labelMatcher),
	)
}

func withNodeMemoryQuotaGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory Quota",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(12),
		panels.MemoryQuota(datasource, labelMatcher),
	)
}

func BuildKubernetesNodeResourcesOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("kubernetes-node-resources-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Kubernetes / Compute Resources / Node (Pods)"),
			dashboard.AddVariable("cluster",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("cluster",
						labelValuesVar.Matchers("up{"+panels.GetKubeletMatcher()+", metrics_path=\"/metrics/cadvisor\"}"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("cluster"),
				),
			),
			dashboard.AddVariable("node",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("node",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"kube_pod_info",
								[]promql.LabelMatcher{{Name: "cluster", Type: "=", Value: "$cluster"}},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("node"),
				),
			),
			withNodeCPUUsageGroup(datasource, clusterLabelMatcher),
			withNodeCPUQuotaGroup(datasource, clusterLabelMatcher),
			withNodeMemoryUsageGroup(datasource, clusterLabelMatcher),
			withNodeMemoryUsageWithoutCacheGroup(datasource, clusterLabelMatcher),
			withNodeMemoryQuotaGroup(datasource, clusterLabelMatcher),
		),
	).Component("kubernetes")
}
