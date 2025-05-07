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

func withMultiClusterStatsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Multi-Cluster Stats",
		panelgroup.PanelsPerLine(6),
		panelgroup.PanelHeight(4),
		panels.KubernetesCPUUtilizationStat("multicluster", datasource, labelMatcher),
		panels.KubernetesCPURequestsCommitmentStat("multicluster", datasource, labelMatcher),
		panels.KubernetesCPULimitsCommitmentStat("multicluster", datasource, labelMatcher),
		panels.KubernetesMemoryUtilizationStat("multicluster", datasource, labelMatcher),
		panels.KubernetesMemoryRequestsCommitmentStat("multicluster", datasource, labelMatcher),
		panels.KubernetesMemoryLimitsCommitmentStat("multicluster", datasource, labelMatcher),
	)
}

func withMultiClusterCPUUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Multi-Cluster CPU Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.KubernetesCPUUsage("multicluster", datasource, labelMatcher),
	)
}

func withMultiClusterCPUUsageQuotaGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Multi-Cluster CPU Usage Quota",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.MultiClusterCPUUsageQuota(datasource, labelMatcher),
	)
}

func withMultiClusterMemoryUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Multi-Cluster Memory Usage",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.KubernetesMemoryUsage("multicluster", datasource, labelMatcher),
	)
}

func withMultiClusterMemoryUsageQuotaGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Multi-Cluster Memory Usage Quota",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.MultiClusterMemoryUsageQuota(datasource, labelMatcher),
	)
}

func BuildKubernetesMultiClusterOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("kubernetes-multi-cluster-resources-overview",
		dashboard.ProjectName(project),
		dashboard.Name("Kubernetes / Compute Resources / Multi-Cluster"),
		dashboard.AddVariable("cluster",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("cluster",
					labelValuesVar.Matchers("up{"+panels.GetKubeletMatcher()+", metrics_path=\"/metrics/cadvisor\"}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("cluster"),
			),
		),
		withMultiClusterStatsGroup(datasource, clusterLabelMatcher),
		withMultiClusterCPUUsageGroup(datasource, clusterLabelMatcher),
		withMultiClusterCPUUsageQuotaGroup(datasource, clusterLabelMatcher),
		withMultiClusterMemoryUsageGroup(datasource, clusterLabelMatcher),
		withMultiClusterMemoryUsageQuotaGroup(datasource, clusterLabelMatcher),
	)
}
