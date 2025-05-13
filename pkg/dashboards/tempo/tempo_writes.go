package tempo

import (
	"github.com/perses/perses/go-sdk/dashboard"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"

	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
)

func BuildTempoWritesOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("thanos-writes-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Tempo / Writes"),
			dashboard.AddVariable("cluster",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("cluster",
						labelValuesVar.Matchers("tempo_build_info"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("cluster"),
					listVar.AllowMultiple(true),
				),
			),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "tempo_build_info"),
			dashboard.AddVariable("namespace",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("namespace",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"tempo_build_info",
								[]promql.LabelMatcher{clusterLabelMatcher, {Name: "cluster", Type: "=", Value: "$cluster"}},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("namespace"),
				),
			),
		),
	).Component("tempo")
}
