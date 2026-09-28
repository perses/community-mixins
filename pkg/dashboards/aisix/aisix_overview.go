// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aisix

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	panels "github.com/perses/community-mixins/pkg/panels/aisix"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/model/labels"
)

func withOverviewGroup(datasource string, labelMatchers ...*labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Traffic and latency",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.RequestRate(datasource, labelMatchers...),
		panels.SuccessRatio(datasource, labelMatchers...),
		panels.EndToEndLatencyP95(datasource, labelMatchers...),
		panels.TimeToFirstTokenP95(datasource, labelMatchers...),
	)
}

func withUsageAndPolicyGroup(datasource string, labelMatchers ...*labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Usage and policy signals",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.TokenThroughput(datasource, labelMatchers...),
		panels.CacheOutcomes(datasource, labelMatchers...),
		panels.PolicyEnforcement(datasource, labelMatchers...),
		panels.OTLPDelivery(datasource, labelMatchers...),
	)
}

func BuildAISIXOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	labelMatchers := []*labels.Matcher{promql.JobVarV2, promql.InstanceVarV2, clusterLabelMatcher}

	return dashboards.NewDashboardResult(
		dashboard.New("aisix-ai-gateway-overview",
			dashboard.ProjectName(project),
			dashboard.Name("AISIX AI Gateway / Overview"),
			dashboard.AddVariable("job",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("job",
						labelValuesVar.Matchers("aisix_requests_total{}"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.AllowAllValue(true),
					listVar.AllowMultiple(true),
					listVar.DisplayName("job"),
				),
			),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "aisix_requests_total{}"),
			dashboard.AddVariable("instance",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("instance",
						labelValuesVar.Matchers(
							promql.SetLabelMatchersV2(
								vector.New(vector.WithMetricName("aisix_requests_total")),
								[]*labels.Matcher{clusterLabelMatcher, {Name: "job", Type: labels.MatchRegexp, Value: "$job"}},
							).Pretty(0),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.AllowAllValue(true),
					listVar.AllowMultiple(true),
					listVar.DisplayName("instance"),
				),
			),
			withOverviewGroup(datasource, labelMatchers...),
			withUsageAndPolicyGroup(datasource, labelMatchers...),
		),
	).Component("aisix")
}
