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

func BuildKubernetesNamespaceOverview(project string,
	datasource string,
	clusterLabelName string,
	variableOverrides ...dashboard.Option,
) dashboards.DashboardResult {
	defaultVars := []dashboard.Option{
		dashboard.AddVariable("cluster",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("cluster",
					labelValuesVar.Matchers("up{"+panels.GetKubeletMatcher()+"}"),
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
	}

	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)

	vars := defaultVars
	if len(variableOverrides) > 0 {
		vars = variableOverrides
	}
	options := append([]dashboard.Option{
		dashboard.ProjectName(project),
		dashboard.Name("Kubernetes / Compute Resources / Namespace (Pods)"),
	}, vars...)
	options = append(options,
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
	)
	return dashboards.NewDashboardResult(
		dashboard.New("kubernetes-namespace-resources-overview", options...),
	).Component("kubernetes")
}
