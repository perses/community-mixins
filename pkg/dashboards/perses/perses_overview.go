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

package perses

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	panelsGostats "github.com/perses/community-mixins/pkg/panels/gostats"
	"github.com/perses/community-mixins/pkg/panels/perses"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listvariable "github.com/perses/perses/go-sdk/variable/list-variable"
	labelvalues "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/model/labels"
)

func BuildPersesOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("perses-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Perses / Overview"),
			dashboard.AddVariable("job",
				listvariable.List(
					labelvalues.PrometheusLabelValues("job",
						labelvalues.Matchers("perses_build_info{}"),
						dashboards.AddVariableDatasource(datasource),
					),
					listvariable.DisplayName("job"),
				),
			),
			dashboard.AddVariable("instance",
				listvariable.List(
					labelvalues.PrometheusLabelValues("instance",
						labelvalues.Matchers(
							promql.SetLabelMatchersV2(
								vector.New(vector.WithMetricName("perses_build_info")),
								[]*labels.Matcher{clusterLabelMatcher, {Name: "job", Type: labels.MatchEqual, Value: "$job"}},
							).Pretty(0),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listvariable.DisplayName("instance"),
				),
			),
			withPersesOverviewStatsGroup(datasource, clusterLabelMatcher),
			withPersesAPiRequestGroup(datasource, clusterLabelMatcher),
			withPersesResources(datasource, clusterLabelMatcher),
			withPersesPlugins(datasource, clusterLabelMatcher),
		),
	).Component("perses")
}

func withPersesOverviewStatsGroup(datasource string, clusterLabelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Perses Stats", panelgroup.PanelsPerLine(1),
		perses.StatsTable(datasource, clusterLabelMatcher))
}

func withPersesAPiRequestGroup(datasource string, clusterLabelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("API Requests", panelgroup.PanelsPerLine(2),
		perses.HTTPRequestsRatePanel(datasource, clusterLabelMatcher),
		perses.HTTPErrorPercentagePanel(datasource, clusterLabelMatcher),
		perses.HTTPRequestsLatencyPanel(datasource, clusterLabelMatcher),
	)
}

func withPersesResources(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	labelMatchersToUse := []*labels.Matcher{
		promql.InstanceVarV2,
		promql.JobVarV2,
	}
	labelMatchersToUse = append(labelMatchersToUse, labelMatcher)

	return dashboard.AddPanelGroup("Resource Usage",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(10),
		panelsGostats.MemoryUsage(datasource, "pod", labelMatchersToUse...),
		panelsGostats.CPUUsage(datasource, "pod", labelMatchersToUse...),
		panelsGostats.Goroutines(datasource, "pod", labelMatchersToUse...),
		panelsGostats.GarbageCollectionPauseTimeQuantiles(datasource, "pod", labelMatchersToUse...),
		perses.FileDescriptors(datasource, labelMatchersToUse...),
	)
}

func withPersesPlugins(datasource string, clusterLabelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Plugins Usage", panelgroup.PanelsPerLine(1), panelgroup.PanelHeight(8),
		perses.PluginSchemaLoadAttempts(datasource, clusterLabelMatcher))
}
