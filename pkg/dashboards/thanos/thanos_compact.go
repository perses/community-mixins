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

func withThanosCompactTODOGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("TODO Operations",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(10),
		panels.TodoCompactionBlocks(datasource, labelMatcher),
		panels.TodoCompactions(datasource, labelMatcher),
		panels.TodoDeletions(datasource, labelMatcher),
		panels.TodoDownsamples(datasource, labelMatcher),
	)
}

func withThanosCompactGroupCompactionGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Group Compactions",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(10),
		panels.GroupCompactionRate(datasource, labelMatcher),
		panels.GroupCompactionErrors(datasource, labelMatcher),
	)
}

func withThanosCompactDownsampleGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Downsample Operations",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(10),
		panels.DownsampleRate(datasource, labelMatcher),
		panels.DownsampleErrors(datasource, labelMatcher),
		panels.DownsampleDurations(datasource, labelMatcher),
	)
}

func withThanosCompactSyncMetaGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Sync Meta",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(10),
		panels.SyncMetaRate(datasource, labelMatcher),
		panels.SyncMetaErrors(datasource, labelMatcher),
		panels.SyncMetaDurations(datasource, labelMatcher),
	)
}

func withThanosCompactBlockDeletionGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Block Deletion",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(10),
		panels.DeletionRate(datasource, labelMatcher),
		panels.DeletionErrors(datasource, labelMatcher),
		panels.MarkingRate(datasource, labelMatcher),
	)
}

func withThanosCompactHaltedGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Halted Compactors",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.HaltedCompactors(datasource, labelMatcher),
	)
}

func withThanosCompactGarbageCollectionGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Garbage Collection",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(10),
		panels.GarbageCollectionRate(datasource, labelMatcher),
		panels.GarbageCollectionErrors(datasource, labelMatcher),
		panels.GarbageCollectionDurations(datasource, labelMatcher),
	)
}

func BuildThanosCompactOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("thanos-compact-overview",
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
			withThanosCompactTODOGroup(datasource, clusterLabelMatcher),
			withThanosCompactGroupCompactionGroup(datasource, clusterLabelMatcher),
			withThanosCompactDownsampleGroup(datasource, clusterLabelMatcher),
			withThanosCompactSyncMetaGroup(datasource, clusterLabelMatcher),
			withThanosCompactBlockDeletionGroup(datasource, clusterLabelMatcher),
			withThanosBucketOperationsGroup(datasource, clusterLabelMatcher),
			withThanosCompactHaltedGroup(datasource, clusterLabelMatcher),
			withThanosCompactGarbageCollectionGroup(datasource, clusterLabelMatcher),
			withThanosResourcesGroup(datasource, clusterLabelMatcher),
		),
	).Component("thanos")
}
