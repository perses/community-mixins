package perses

import (
	"fmt"

	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	labelvalues "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listvariable "github.com/perses/perses/go-sdk/variable/list-variable"
)

func BuildPersesOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	fmt.Println("clusterLabelMatcher", clusterLabelMatcher.Name)
	return dashboard.New("perses-overview",
		dashboard.ProjectName(project),
		dashboard.Name("Perses / Overview"),
		dashboard.AddVariable("job",
			listvariable.List(
				labelvalues.PrometheusLabelValues("job",
					labelvalues.Matchers("prometheus_build_info{}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listvariable.DisplayName("job"),
			),
		),
		dashboard.AddVariable("instance",
			listvariable.List(
				labelvalues.PrometheusLabelValues("instance",
					labelvalues.Matchers(
						promql.SetLabelMatchers(
							"prometheus_build_info",
							[]promql.LabelMatcher{clusterLabelMatcher, {Name: "job", Type: "=", Value: "$job"}},
						),
					),
					dashboards.AddVariableDatasource(datasource),
				),
				listvariable.DisplayName("instance"),
			),
		),
	)
}
