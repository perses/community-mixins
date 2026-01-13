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

package persistent_volume

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	panels "github.com/perses/community-mixins/pkg/panels/kubernetes"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"

	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
)

func withPVVolumeUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Volume Space Usage",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.VolumeSpaceUsage(datasource, labelMatcher),
		panels.VolumeSpaceUsageGauge(datasource, labelMatcher),
	)
}

func withPVInodesUsageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Volume Inodes Usage",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.VolumeInodesUsage(datasource, labelMatcher),
		panels.VolumeInodesUsageGauge(datasource, labelMatcher),
	)
}

func BuildKubernetesPersistentVolumeOverview(project string, datasource string, clusterLabelName string, variableOverrides ...dashboard.Option) dashboards.DashboardResult {
	defaultVars := []dashboard.Option{
		dashboard.AddVariable("cluster",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("cluster",
					labelValuesVar.Matchers("kubelet_volume_stats_capacity_bytes{"+panels.GetKubeletMatcher()+"}"),
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
							"kubelet_volume_stats_capacity_bytes{"+panels.GetKubeletMatcher()+"}",
							[]promql.LabelMatcher{{Name: "cluster", Type: "=", Value: "$cluster"}},
						),
					),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("namespace"),
			),
		),
		dashboard.AddVariable("volume",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("persistentvolumeclaim",
					labelValuesVar.Matchers(
						promql.SetLabelMatchers(
							"kubelet_volume_stats_capacity_bytes{"+panels.GetKubeletMatcher()+"}",
							[]promql.LabelMatcher{
								{Name: "cluster", Type: "=", Value: "$cluster"},
								{Name: "namespace", Type: "=", Value: "$namespace"},
							},
						),
					),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("PersistentVolume"),
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
		dashboard.Name("Kubernetes / Persistent Volume"),
	}, vars...)
	options = append(options,
		withPVVolumeUsageGroup(datasource, clusterLabelMatcher),
		withPVInodesUsageGroup(datasource, clusterLabelMatcher),
	)
	return dashboards.NewDashboardResult(
		dashboard.New("kubernetes-persistent-volume-overview", options...),
	).Component("kubernetes")
}
