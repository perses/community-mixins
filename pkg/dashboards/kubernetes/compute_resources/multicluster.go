// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
