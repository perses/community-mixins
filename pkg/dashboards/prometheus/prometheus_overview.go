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

package prometheus

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	panelsGostats "github.com/perses/community-mixins/pkg/panels/gostats"
	panels "github.com/perses/community-mixins/pkg/panels/prometheus"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/model/labels"

	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
)

func withPrometheusOverviewStatsGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Prometheus Stats",
		panelgroup.PanelsPerLine(1),
		panels.PrometheusStatsTable(datasource, labelMatcher),
	)
}

func withPrometheusOverviewDiscoveryGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Discovery",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusTargetSync(datasource, labelMatcher),
		panels.PrometheusTargets(datasource, labelMatcher),
	)
}

func withPrometheusRetrievalGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Retrieval",
		panelgroup.PanelsPerLine(3),
		panels.PrometheusAverageScrapeIntervalDuration(datasource, labelMatcher),
		panels.PrometheusScrapeFailures(datasource, labelMatcher),
		panels.PrometheusAppendedSamples(datasource, labelMatcher),
	)
}

func withPrometheusStorageGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusHeadSeries(datasource, labelMatcher),
		panels.PrometheusHeadChunks(datasource, labelMatcher),
	)
}

func withPrometheusQueryGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Query",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusQueryRate(datasource, labelMatcher),
		panels.PrometheusQueryStateDuration(datasource, labelMatcher),
	)
}

func withPrometheusResourcesGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	labelMatchersToUse := []*labels.Matcher{
		promql.InstanceVarV2,
		promql.JobVarV2,
	}
	labelMatchersToUse = append(labelMatchersToUse, labelMatcher)

	return dashboard.AddPanelGroup("Resources",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(10),
		panelsGostats.CPUUsage(datasource, "pod", labelMatchersToUse...),
		panelsGostats.MemoryUsage(datasource, "pod", labelMatchersToUse...),
		panelsGostats.Goroutines(datasource, "pod", labelMatchersToUse...),
		panelsGostats.GarbageCollectionPauseTimeQuantiles(datasource, "pod", labelMatchersToUse...),
	)
}

func BuildPrometheusOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("prometheus-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Prometheus / Overview"),
			dashboard.AddVariable("job",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("job",
						labelValuesVar.Matchers("prometheus_build_info{}"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("job"),
				),
			),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "prometheus_build_info"),
			dashboard.AddVariable("instance",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("instance",
						labelValuesVar.Matchers(
							promql.SetLabelMatchersV2(
								vector.New(vector.WithMetricName("prometheus_build_info")),
								[]*labels.Matcher{clusterLabelMatcher, {Name: "job", Type: labels.MatchEqual, Value: "$job"}},
							).Pretty(0),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("instance"),
				),
			),
			withPrometheusOverviewStatsGroup(datasource, clusterLabelMatcher),
			withPrometheusOverviewDiscoveryGroup(datasource, clusterLabelMatcher),
			withPrometheusRetrievalGroup(datasource, clusterLabelMatcher),
			withPrometheusStorageGroup(datasource, clusterLabelMatcher),
			withPrometheusQueryGroup(datasource, clusterLabelMatcher),
			withPrometheusResourcesGroup(datasource, clusterLabelMatcher),
		),
	).Component("prometheus")
}
