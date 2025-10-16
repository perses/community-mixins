package compute_resources

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	panels "github.com/perses/community-mixins/pkg/panels/kubernetes"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
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

func BuildKubernetesMultiClusterOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("kubernetes-multi-cluster-resources-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Kubernetes / Compute Resources / Multi-Cluster"),
			withMultiClusterStatsGroup(datasource, clusterLabelMatcher),
			withMultiClusterCPUUsageGroup(datasource, clusterLabelMatcher),
			withMultiClusterCPUUsageQuotaGroup(datasource, clusterLabelMatcher),
			withMultiClusterMemoryUsageGroup(datasource, clusterLabelMatcher),
			withMultiClusterMemoryUsageQuotaGroup(datasource, clusterLabelMatcher),
		),
	).Component("kubernetes")
}
