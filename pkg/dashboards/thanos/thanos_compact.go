package thanos

import (
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"

	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/thanos"
	"github.com/perses/community-dashboards/pkg/promql"
)

func withThanosTODOGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("TODO Operations",
		panelgroup.PanelsPerLine(4),
		panels.TodoCompactionBlocks(datasource, labelMatcher),
		panels.TodoCompactions(datasource, labelMatcher),
		panels.TodoDeletions(datasource, labelMatcher),
		panels.TodoDownsamples(datasource, labelMatcher),
	)
}

func withThanosGroupCompactionGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Group Compactions",
		panelgroup.PanelsPerLine(2),
		panels.GroupCompactions(datasource, labelMatcher),
		panels.GroupCompactionErrors(datasource, labelMatcher),
	)
}

func withThanosDownsampleGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Downsample Operations",
		panelgroup.PanelsPerLine(3),
		panels.DownsampleRate(datasource, labelMatcher),
		panels.DownsampleErrors(datasource, labelMatcher),
		panels.DownsampleDurations(datasource, labelMatcher),
	)
}

func withThanosSyncMetaGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Sync Meta",
		panelgroup.PanelsPerLine(3),
		panels.SyncMetaRate(datasource, labelMatcher),
		panels.SyncMetaErrors(datasource, labelMatcher),
		panels.SyncMetaDurations(datasource, labelMatcher),
	)
}

func withThanosBlockDeletionGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Block Deletion",
		panelgroup.PanelsPerLine(3),
		panels.DeletionRate(datasource, labelMatcher),
		panels.DeletionErrors(datasource, labelMatcher),
		panels.MarkingRate(datasource, labelMatcher),
	)
}

func withThanosHaltedGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Halted Compactors",
		panelgroup.PanelsPerLine(1),
		panels.HaltedCompactors(datasource, labelMatcher),
	)
}

func withThanosGarbageCollectionGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Garbage Collection",
		panelgroup.PanelsPerLine(3),
		panels.GarbageCollection(datasource, labelMatcher),
		panels.GarbageCollectionErrors(datasource, labelMatcher),
		panels.GarbageCollectionDurations(datasource, labelMatcher),
	)
}

func BuildThanosCompactOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("thanos-compact-overview",
		dashboard.ProjectName(project),
		dashboard.Name("Thanos / Compact / Overview"),
		dashboard.AddVariable("job",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("job",
					labelValuesVar.Matchers("thanos_build_info{container=\"thanos-compact\"}"),
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
		withThanosTODOGroup(datasource, clusterLabelMatcher),
		withThanosGroupCompactionGroup(datasource, clusterLabelMatcher),
		withThanosDownsampleGroup(datasource, clusterLabelMatcher),
		withThanosSyncMetaGroup(datasource, clusterLabelMatcher),
		withThanosBlockDeletionGroup(datasource, clusterLabelMatcher),
		withThanosBucketOperationsGroup(datasource, clusterLabelMatcher),
		withThanosHaltedGroup(datasource, clusterLabelMatcher),
		withThanosGarbageCollectionGroup(datasource, clusterLabelMatcher),
		withThanosResourcesGroup(datasource, clusterLabelMatcher),
	)
}
