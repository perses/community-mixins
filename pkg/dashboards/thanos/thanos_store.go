package thanos

import (
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/prometheus/prometheus/model/labels"

	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/thanos"
)

func withThanosBlockOperationsGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Block Operations",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(8),
		panels.BlockLoadRate(datasource, labelMatcher),
		panels.BlockLoadErrors(datasource, labelMatcher),
		panels.BlockDropRate(datasource, labelMatcher),
		panels.BlockDropErrors(datasource, labelMatcher),
	)
}

func withThanosCacheOperationsGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Cache Operations",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(8),
		panels.CacheRequestRate(datasource, labelMatcher),
		panels.CacheHitRate(datasource, labelMatcher),
		panels.CacheItemsAddRate(datasource, labelMatcher),
		panels.CacheItemsEvictRate(datasource, labelMatcher),
	)
}

func withThanosQueryOperationsGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Query Operations",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(8),
		panels.BlocksQueried(datasource, labelMatcher),
		panels.DataFetched(datasource, labelMatcher),
		panels.DataTouched(datasource, labelMatcher),
		panels.ResultSeries(datasource, labelMatcher),
	)
}

func withThanosQueryOperationDurationGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Query Operation Durations",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.GetAllSeriesDurations(datasource, labelMatcher),
		panels.MergeDurations(datasource, labelMatcher),
		panels.GateWaitingDurations(datasource, labelMatcher),
	)
}

func withThanosStoreSentGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Store Sent Chunk Size",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.StoreSentChunkSizes(datasource, labelMatcher),
	)
}

func BuildThanosStoreOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	clusterLabelMatcherV2 := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("thanos-store-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Thanos / Store Gateway / Overview"),
			dashboard.AddVariable("job",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("job",
						labelValuesVar.Matchers("thanos_build_info{container=\"thanos-store\"}"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("job"),
					listVar.AllowMultiple(true),
				),
			),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "thanos_build_info"),
			dashboard.AddVariable("namespace",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("namespace",
						labelValuesVar.Matchers("thanos_status"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("namespace"),
				),
			),
			withThanosReadGRPCUnaryGroup(datasource, clusterLabelMatcherV2),
			withThanosReadGRPCStreamGroup(datasource, clusterLabelMatcherV2),
			withThanosBucketOperationsGroup(datasource, clusterLabelMatcherV2),
			withThanosBlockOperationsGroup(datasource, clusterLabelMatcherV2),
			withThanosCacheOperationsGroup(datasource, clusterLabelMatcherV2),
			withThanosQueryOperationsGroup(datasource, clusterLabelMatcherV2),
			withThanosQueryOperationDurationGroup(datasource, clusterLabelMatcherV2),
			withThanosStoreSentGroup(datasource, clusterLabelMatcherV2),
			withThanosResourcesGroup(datasource, clusterLabelMatcher),
		),
	).Component("thanos")
}
