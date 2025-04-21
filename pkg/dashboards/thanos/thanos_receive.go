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

func withThanosRemoteWriteGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Remote Write v1 - Incoming Requests",
		panelgroup.PanelsPerLine(3),
		panels.RemoteWriteRequestRate(datasource, labelMatcher),
		panels.RemoteWriteRequestErrors(datasource, labelMatcher),
		panels.RemoteWriteRequestDuration(datasource, labelMatcher),
	)
}

func withThanosRemoteWriteTenantedGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Remote Write v1 - Incoming Requests (tenanted)",
		panelgroup.PanelsPerLine(3),
		panels.TenantedRemoteWriteRequestRate(datasource, labelMatcher),
		panels.TenantedRemoteWriteRequestErrors(datasource, labelMatcher),
		panels.TenantedRemoteWriteRequestDuration(datasource, labelMatcher),
	)
}

func withThanosRemoteWriteHTTPGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Remote Write v1 - HTTP Requests",
		panelgroup.PanelsPerLine(3),
		panels.AvgRemoteWriteRequestSize(datasource, labelMatcher),
		panels.AvgFailedRemoteWriteRequestSize(datasource, labelMatcher),
		panels.InflightRemoteWriteRequests(datasource, labelMatcher),
	)
}

func withThanosRemoteWriteSeriesSampleGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Series and Samples (tenanted)",
		panelgroup.PanelsPerLine(4),
		panels.RemoteWriteSeriesRate(datasource, labelMatcher),
		panels.RemoteWriteSeriesNotWrittenRate(datasource, labelMatcher),
		panels.RemoteWriteSamplesRate(datasource, labelMatcher),
		panels.RemoteWriteSamplesNotWrittenRate(datasource, labelMatcher),
	)
}

func withThanosRemoteWriteReplicationGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Remote Write Replication",
		panelgroup.PanelsPerLine(2),
		panels.RemoteWriteReplicationRate(datasource, labelMatcher),
		panels.RemoteWriteReplicationErrorRate(datasource, labelMatcher),
	)
}

func withThanosRemoteWriteForwardGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Remote Write Forward",
		panelgroup.PanelsPerLine(2),
		panels.RemoteWriteForwardRate(datasource, labelMatcher),
		panels.RemoteWriteForwardErrRate(datasource, labelMatcher),
	)
}

func withThanosWriteGRPCUnaryGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Write gRPC Unary (WritableStore)",
		panelgroup.PanelsPerLine(3),
		panels.WriteGRPCUnaryRate(datasource, labelMatcher),
		panels.WriteGRPCUnaryErrRate(datasource, labelMatcher),
		panels.WriteGPRCUnaryDuration(datasource, labelMatcher),
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

func withThanosBucketUploadGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Last Bucket Upload",
		panelgroup.PanelsPerLine(1),
		panels.BucketUploadTable(datasource, labelMatcher),
	)
}

func withPrometheusStorageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage",
		panelgroup.PanelsPerLine(3),
		panels.ReceiveAppendedSamples(datasource, labelMatcher),
		panels.ReceiveHeadSeries(datasource, labelMatcher),
		panels.ReceiveHeadChunks(datasource, labelMatcher),
	)
}

func BuildThanosReceiveOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("thanos-receive-overview",
		dashboard.ProjectName(project),
		dashboard.Name("Thanos / Receive / Overview"),
		dashboard.AddVariable("job",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("job",
					labelValuesVar.Matchers("thanos_build_info{container=\"thanos-receive\"}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("job"),
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
		dashboard.AddVariable("tenant",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("tenant",
					labelValuesVar.Matchers(
						promql.SetLabelMatchers(
							"prometheus_tsdb_head_max_time{container=\"thanos-receive\"}",
							[]promql.LabelMatcher{clusterLabelMatcher, {Name: "job", Type: "=", Value: "$job"}, {Name: "namespace", Type: "=", Value: "$namespace"}},
						),
					),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("tenant"),
				listVar.AllowMultiple(true),
			),
		),
		withThanosRemoteWriteGroup(datasource, clusterLabelMatcher),
		withThanosRemoteWriteTenantedGroup(datasource, clusterLabelMatcher),
		withThanosRemoteWriteHTTPGroup(datasource, clusterLabelMatcher),
		withThanosRemoteWriteSeriesSampleGroup(datasource, clusterLabelMatcher),
		withThanosRemoteWriteReplicationGroup(datasource, clusterLabelMatcher),
		withThanosRemoteWriteForwardGroup(datasource, clusterLabelMatcher),
		withThanosWriteGRPCUnaryGroup(datasource, clusterLabelMatcher),
		withThanosReadGRPCUnaryGroup(datasource, clusterLabelMatcher),
		withThanosReadGRPCStreamGroup(datasource, clusterLabelMatcher),
		withThanosBucketUploadGroup(datasource, clusterLabelMatcher),
		withPrometheusStorageGroup(datasource, clusterLabelMatcher),
	)
}
