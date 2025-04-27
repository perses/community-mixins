package thanos

import (
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"

	panels "github.com/perses/community-dashboards/pkg/panels/thanos"
	"github.com/perses/community-dashboards/pkg/promql"
)

func withThanosResourcesGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Resources",
		panelgroup.PanelsPerLine(3),
		panels.MemoryUsage(datasource, labelMatcher),
		panels.Goroutines(datasource, labelMatcher),
		panels.GCDurationQuantiles(datasource, labelMatcher),
	)
}

func withThanosBucketOperationsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Bucket Operations",
		panelgroup.PanelsPerLine(3),
		panels.BucketOperations(datasource, labelMatcher),
		panels.BucketOperationErrors(datasource, labelMatcher),
		panels.BucketOperationLatency(datasource, labelMatcher),
	)
}

func withThanosReadGRPCUnaryGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Read gRPC Unary (StoreAPI Info/Labels)",
		panelgroup.PanelsPerLine(3),
		panels.ReadGRPCUnaryRate(datasource, labelMatcher),
		panels.ReadGRPCUnaryErrRate(datasource, labelMatcher),
		panels.ReadGPRCUnaryDuration(datasource, labelMatcher),
	)
}

func withThanosReadGRPCStreamGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Read gRPC Stream (StoreAPI Series/Exemplars)",
		panelgroup.PanelsPerLine(3),
		panels.ReadGRPCStreamRate(datasource, labelMatcher),
		panels.ReadGRPCStreamErrRate(datasource, labelMatcher),
		panels.ReadGPRCStreamDuration(datasource, labelMatcher),
	)
}
