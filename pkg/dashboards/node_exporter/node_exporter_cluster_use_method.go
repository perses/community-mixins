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

package nodeexporter

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	panels "github.com/perses/community-mixins/pkg/panels/node_exporter"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/model/labels"
)

func withClusterCPU(datasource string, clusterLabelMatcher *labels.Matcher, instanceLabelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU",
		panelgroup.PanelsPerLine(2),
		panels.ClusterNodeCPUUsagePercentage(datasource, clusterLabelMatcher, instanceLabelMatcher),
		panels.ClusterNodeCPUSaturationPercentage(datasource, clusterLabelMatcher, instanceLabelMatcher),
	)
}

func withClusterMemory(datasource string, clusterLabelMatcher *labels.Matcher, instanceLabelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory",
		panelgroup.PanelsPerLine(2),
		panels.ClusterNodeMemoryUsagePercentage(datasource, clusterLabelMatcher, instanceLabelMatcher),
		panels.ClusterNodeMemorySaturationPercentage(datasource, clusterLabelMatcher, instanceLabelMatcher),
	)
}

func withClusterNetwork(datasource string, clusterLabelMatcher *labels.Matcher, instanceLabelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Network",
		panelgroup.PanelsPerLine(2),
		panels.ClusterNodeNetworkUsageBytes(datasource, clusterLabelMatcher, instanceLabelMatcher),
		panels.ClusterNodeNetworkSaturationBytes(datasource, clusterLabelMatcher, instanceLabelMatcher),
	)
}

func withClusterDiskIO(datasource string, clusterLabelMatcher *labels.Matcher, instanceLabelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Disk IO",
		panelgroup.PanelsPerLine(2),
		panels.ClusterNodeDiskUsagePercentage(datasource, clusterLabelMatcher, instanceLabelMatcher),
		panels.ClusterNodeDiskSaturationPercentage(datasource, clusterLabelMatcher, instanceLabelMatcher),
	)
}

func withClusterDiskSpace(datasource string, clusterLabelMatcher *labels.Matcher, instanceLabelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Disk Space",
		panelgroup.PanelsPerLine(1),
		panels.ClusterNodeDiskSpacePercentage(datasource, clusterLabelMatcher, instanceLabelMatcher),
	)
}

func BuildNodeExporterClusterUseMethod(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	instanceLabelMatcher := &labels.Matcher{
		Name:  "instance",
		Value: "$instance",
		Type:  labels.MatchRegexp,
	}
	return dashboards.NewDashboardResult(
		dashboard.New("node-exporter-cluster-use-method",
			dashboard.ProjectName(project),
			dashboard.Name("Node Exporter / USE Method / Cluster"),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "node_uname_info{job='node', sysname!='Darwin'}"),
			dashboard.AddVariable("instance",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("instance",
						dashboards.AddVariableDatasource(datasource),
						labelValuesVar.Matchers(
							promql.SetLabelMatchersV2(
								vector.New(vector.WithMetricName("node_uname_info")),
								[]*labels.Matcher{clusterLabelMatcher,
									{Name: "job", Type: labels.MatchEqual, Value: "node"},
									{Name: "sysname", Type: labels.MatchNotEqual, Value: "Darwin"}},
							).Pretty(0),
						),
					),
					listVar.DisplayName("instance"),
					listVar.AllowAllValue(true),
					listVar.AllowMultiple(true),
				),
			),
			withClusterCPU(datasource, clusterLabelMatcher, instanceLabelMatcher),
			withClusterMemory(datasource, clusterLabelMatcher, instanceLabelMatcher),
			withClusterNetwork(datasource, clusterLabelMatcher, instanceLabelMatcher),
			withClusterDiskIO(datasource, clusterLabelMatcher, instanceLabelMatcher),
			withClusterDiskSpace(datasource, clusterLabelMatcher, instanceLabelMatcher),
		),
	).Component("node-exporter")
}
