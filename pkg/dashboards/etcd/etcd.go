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

package etcd

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	panels "github.com/perses/community-mixins/pkg/panels/etcd"
	panelsGostats "github.com/perses/community-mixins/pkg/panels/gostats"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/model/labels"
)

func withETCDStatsGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("etcd Status",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.EtcdUpStatus(datasource, labelMatcher),
	)
}

func withRPCGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("RPC and Streams",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.RPCRate(datasource, labelMatcher),
		panels.ActiveStreams(datasource, labelMatcher),
	)
}

func withDBGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("etcd DB",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.DBSize(datasource, labelMatcher),
		panels.DiskSyncDuration(datasource, labelMatcher),
	)
}

func withTrafficGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("etcd Traffic",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.ClientTrafficIn(datasource, labelMatcher),
		panels.ClientTrafficOut(datasource, labelMatcher),
		panels.PeerTrafficIn(datasource, labelMatcher),
		panels.PeerTrafficOut(datasource, labelMatcher),
	)
}

func withRaftGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("etcd Raft",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.RaftProposals(datasource, labelMatcher),
	)
}

func withRoundTripTimeGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("etcd Peer Round Trip Time",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.PeerRoundtripTime(datasource, labelMatcher),
	)
}

func withETCDResources(datasource string, clusterLabelMatcher *labels.Matcher) dashboard.Option {
	// TODO(saswatamcode): Add a way to configure these.
	labelMatchersToUse := []*labels.Matcher{
		promql.ClusterVarV2,
		{
			Name:  "job",
			Value: ".*etcd.*",
			Type:  labels.MatchRegexp,
		},
	}

	labelMatchersToUse = append(labelMatchersToUse, clusterLabelMatcher)

	return dashboard.AddPanelGroup("Resource Usage",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panelsGostats.MemoryUsage(datasource, "instance", labelMatchersToUse...),
		panelsGostats.CPUUsage(datasource, "instance", labelMatchersToUse...),
		panelsGostats.Goroutines(datasource, "instance", labelMatchersToUse...),
		panelsGostats.GarbageCollectionPauseTimeQuantiles(datasource, "instance", labelMatchersToUse...),
	)
}

func BuildETCDOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("etcd-overview",
			dashboard.ProjectName(project),
			dashboard.Name("etcd / Overview"),
			dashboard.AddVariable("cluster",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("cluster",
						labelValuesVar.Matchers(
							promql.SetLabelMatchersV2(
								vector.New(vector.WithMetricName("etcd_server_has_leader")),
								[]*labels.Matcher{clusterLabelMatcher, {Name: "job", Type: labels.MatchRegexp, Value: ".*etcd.*"}},
							).Pretty(0),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("cluster"),
				),
			),
			withETCDStatsGroup(datasource, clusterLabelMatcher),
			withRPCGroup(datasource, clusterLabelMatcher),
			withDBGroup(datasource, clusterLabelMatcher),
			withRaftGroup(datasource, clusterLabelMatcher),
			withTrafficGroup(datasource, clusterLabelMatcher),
			withRoundTripTimeGroup(datasource, clusterLabelMatcher),
			withETCDResources(datasource, clusterLabelMatcher),
		),
	).Component("etcd")
}
