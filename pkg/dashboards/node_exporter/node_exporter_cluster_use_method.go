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
	"fmt"

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

func withClusterCPU(datasource string, labelMatchers ...*labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU",
		panelgroup.PanelsPerLine(2),
		panels.ClusterNodeCPUUsagePercentage(datasource, labelMatchers...),
		panels.ClusterNodeCPUSaturationPercentage(datasource, labelMatchers...),
	)
}

func withClusterMemory(datasource string, labelMatchers ...*labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory",
		panelgroup.PanelsPerLine(2),
		panels.ClusterNodeMemoryUsagePercentage(datasource, labelMatchers...),
		panels.ClusterNodeMemorySaturationPercentage(datasource, labelMatchers...),
	)
}

func withClusterNetwork(datasource string, labelMatchers ...*labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Network",
		panelgroup.PanelsPerLine(2),
		panels.ClusterNodeNetworkUsageBytes(datasource, labelMatchers...),
		panels.ClusterNodeNetworkSaturationBytes(datasource, labelMatchers...),
	)
}

func withClusterDiskIO(datasource string, labelMatchers ...*labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Disk IO",
		panelgroup.PanelsPerLine(2),
		panels.ClusterNodeDiskUsagePercentage(datasource, labelMatchers...),
		panels.ClusterNodeDiskSaturationPercentage(datasource, labelMatchers...),
	)
}

func withClusterDiskSpace(datasource string, labelMatchers ...*labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Disk Space",
		panelgroup.PanelsPerLine(1),
		panels.ClusterNodeDiskSpacePercentage(datasource, labelMatchers...),
	)
}

func BuildNodeExporterClusterUseMethod(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	instanceLabelMatcher := &labels.Matcher{
		Name:  "instance",
		Value: "$instance",
		Type:  labels.MatchRegexp,
	}
	jobValue := panels.GetNodeExporterLabelValue()
	jobMatcher := &labels.Matcher{Name: "job", Type: labels.MatchEqual, Value: jobValue}
	return dashboards.NewDashboardResult(
		dashboard.New("node-exporter-cluster-use-method",
			dashboard.ProjectName(project),
			dashboard.Name("Node Exporter / USE Method / Cluster"),
			dashboards.AddClusterVariable(datasource, clusterLabelName, fmt.Sprintf("node_uname_info{job='%s', sysname!='Darwin'}", jobValue)),
			dashboard.AddVariable("instance",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("instance",
						dashboards.AddVariableDatasource(datasource),
						labelValuesVar.Matchers(
							promql.SetLabelMatchersV2(
								vector.New(vector.WithMetricName("node_uname_info")),
								[]*labels.Matcher{clusterLabelMatcher,
									{Name: "job", Type: labels.MatchEqual, Value: jobValue},
									{Name: "sysname", Type: labels.MatchNotEqual, Value: "Darwin"}},
							).Pretty(0),
						),
					),
					listVar.DisplayName("instance"),
					listVar.AllowAllValue(true),
					listVar.AllowMultiple(true),
				),
			),
			withClusterCPU(datasource, clusterLabelMatcher, instanceLabelMatcher, jobMatcher),
			withClusterMemory(datasource, clusterLabelMatcher, instanceLabelMatcher, jobMatcher),
			withClusterNetwork(datasource, clusterLabelMatcher, instanceLabelMatcher, jobMatcher),
			withClusterDiskIO(datasource, clusterLabelMatcher, instanceLabelMatcher, jobMatcher),
			withClusterDiskSpace(datasource, clusterLabelMatcher, instanceLabelMatcher, jobMatcher),
		),
	).Component("node-exporter")
}
