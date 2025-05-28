package thanos

import (
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/prometheus/prometheus/model/labels"

	panelsGostats "github.com/perses/community-dashboards/pkg/panels/gostats"
	panels "github.com/perses/community-dashboards/pkg/panels/thanos"
	"github.com/perses/community-dashboards/pkg/promql"
)

func withThanosResourcesGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	labelMatchersToUse := []*labels.Matcher{
		promql.NamespaceVarV2,
		promql.JobVarV2,
	}
	labelMatchersToUse = append(labelMatchersToUse, labelMatcher)

	return dashboard.AddPanelGroup("Resources",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(8),
		panelsGostats.CPUUsage(datasource, "pod", labelMatchersToUse...),
		panelsGostats.MemoryUsage(datasource, "pod", labelMatchersToUse...),
		panelsGostats.Goroutines(datasource, "pod", labelMatchersToUse...),
		panelsGostats.GarbageCollectionPauseTimeQuantiles(datasource, "pod", labelMatchersToUse...),
	)
}

func withThanosBucketOperationsGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Bucket Operations",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.BucketOperationRate(datasource, labelMatcher),
		panels.BucketOperationErrors(datasource, labelMatcher),
		panels.BucketOperationDurations(datasource, labelMatcher),
	)
}

func withThanosReadGRPCUnaryGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Read gRPC Unary (StoreAPI Info/Labels)",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.ReadGRPCUnaryRate(datasource, labelMatcher),
		panels.ReadGRPCUnaryErrors(datasource, labelMatcher),
		panels.ReadGPRCUnaryDurations(datasource, labelMatcher),
	)
}

func withThanosReadGRPCStreamGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Read gRPC Stream (StoreAPI Series/Exemplars)",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.ReadGRPCStreamRate(datasource, labelMatcher),
		panels.ReadGRPCStreamErrors(datasource, labelMatcher),
		panels.ReadGPRCStreamDurations(datasource, labelMatcher),
	)
}

func withThanosBucketUploadGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Last Bucket Upload",
		panelgroup.PanelsPerLine(1),
		panels.BucketUploadTable(datasource, labelMatcher),
	)
}
