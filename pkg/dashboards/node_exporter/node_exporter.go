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
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/promql-builder/vector"

	panels "github.com/perses/community-mixins/pkg/panels/node_exporter"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/prometheus/prometheus/model/labels"
)

func withNodeExporterNodesCPU(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU",
		panelgroup.PanelsPerLine(2),
		panels.NodeCPUUsagePercentage(datasource, labelMatcher),
		panels.NodeAverage(datasource, labelMatcher),
	)
}

func withNodeExporterNodesMemory(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory",
		panelgroup.PanelsPerLine(2),
		panels.NodeMemoryUsageBytes(datasource, labelMatcher),
		panels.NodeMemoryUsagePercentage(datasource, labelMatcher),
	)
}

func withNodeExporterNodesDisk(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Disk",
		panelgroup.PanelsPerLine(2),
		panels.NodeDiskIOBytes(datasource, labelMatcher),
		panels.NodeDiskIOSeconds(datasource, labelMatcher),
	)
}

func withNodeExporterNodesNetwork(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Network",
		panelgroup.PanelsPerLine(2),
		panels.NodeNetworkReceivedBytes(datasource, labelMatcher),
		panels.NodeNetworkTransmitedBytes(datasource, labelMatcher),
	)
}

func BuildNodeExporterNodes(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("node-exporter-nodes",
			dashboard.ProjectName(project),
			dashboard.Name("Node Exporter / Nodes"),
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
						// promql.SetLabelMatchers(
						// 	"node_uname_info{job='node', sysname!='Darwin'}",
						// 	[]*labels.Matcher{clusterLabelMatcher},
						// )),
					),
					listVar.DisplayName("instance"),
					listVar.AllowAllValue(true),
				),
			),
			withNodeExporterNodesCPU(datasource, clusterLabelMatcher),
			withNodeExporterNodesMemory(datasource, clusterLabelMatcher),
			withNodeExporterNodesDisk(datasource, clusterLabelMatcher),
			withNodeExporterNodesNetwork(datasource, clusterLabelMatcher),
		),
	).Component("node-exporter")
}
