package thanos

import (
	"github.com/perses/community-dashboards/pkg/promql"
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

var ThanosCommonPanelQueries = map[string]parser.Expr{
	// Common Thanos Panels
	"BucketOperationRate": promql.SumByRate(
		"thanos_objstore_bucket_operations_total",
		[]string{"namespace", "job", "operation"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"BucketOperationErrors": promql.ErrorCasePercentage(
		"thanos_objstore_bucket_operation_failures_total",
		[]string{"namespace", "job", "operation"},
		[]*labels.Matcher{
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		},
		"thanos_objstore_bucket_operations_total",
		[]string{"namespace", "job", "operation"},
		[]*labels.Matcher{
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		},
	),
	"BucketOperationDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"thanos_objstore_bucket_operation_duration_seconds_bucket",
			[]string{"namespace", "job", "operation", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"BucketOperationDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"thanos_objstore_bucket_operation_duration_seconds_bucket",
			[]string{"namespace", "job", "operation", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"BucketOperationDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"thanos_objstore_bucket_operation_duration_seconds_bucket",
			[]string{"namespace", "job", "operation", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"ReadGRPCUnaryRate": promql.SumByRate(
		"grpc_server_handled_total",
		[]string{"namespace", "job", "grpc_method", "grpc_code"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("grpc_type").Equal("unary"),
		label.New("grpc_method").NotEqual("RemoteWrite"),
	),
	"ReadGRPCUnaryErrors": promql.IgnoringGroupLeft(
		promqlbuilder.Div(
			promql.SumByRate(
				"grpc_server_handled_total",
				[]string{"namespace", "job", "grpc_code"},
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
				label.New("grpc_type").Equal("unary"),
				label.New("grpc_method").NotEqual("RemoteWrite"),
				label.New("grpc_code").EqualRegexp("Unknown|ResourceExhausted|Internal|Unavailable|DataLoss"),
			),
			promql.SumByRate(
				"grpc_server_handled_total",
				[]string{"namespace", "job"},
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
				label.New("grpc_type").Equal("unary"),
				label.New("grpc_method").NotEqual("RemoteWrite"),
			),
		),
		[]string{"grpc_code"},
	),
	"ReadGPRCUnaryDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"grpc_server_handling_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("grpc_type").Equal("unary"),
			label.New("grpc_method").NotEqual("RemoteWrite"),
		),
	),
	"ReadGPRCUnaryDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"grpc_server_handling_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("grpc_type").Equal("unary"),
			label.New("grpc_method").NotEqual("RemoteWrite"),
		),
	),
	"ReadGPRCUnaryDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"grpc_server_handling_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("grpc_type").Equal("unary"),
			label.New("grpc_method").NotEqual("RemoteWrite"),
		),
	),
	"ReadGRPCStreamRate": promql.SumByRate(
		"grpc_server_handled_total",
		[]string{"namespace", "job", "grpc_method", "grpc_code"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("grpc_type").Equal("server_stream"),
	),
	"ReadGRPCStreamErrors": promql.IgnoringGroupLeft(
		promqlbuilder.Div(
			promql.SumByRate(
				"grpc_server_handled_total",
				[]string{"namespace", "job", "grpc_code"},
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
				label.New("grpc_type").Equal("server_stream"),
				label.New("grpc_code").EqualRegexp("Unknown|ResourceExhausted|Internal|Unavailable|DataLoss"),
			),
			promql.SumByRate(
				"grpc_server_handled_total",
				[]string{"namespace", "job"},
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
				label.New("grpc_type").Equal("server_stream"),
			),
		),
		[]string{"grpc_code"},
	),
	"ReadGPRCStreamDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"grpc_server_handling_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("grpc_type").Equal("server_stream"),
		),
	),
	"ReadGPRCStreamDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"grpc_server_handling_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("grpc_type").Equal("server_stream"),
		),
	),
	"ReadGPRCStreamDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"grpc_server_handling_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("grpc_type").Equal("server_stream"),
		),
	),

	// Thanos Compact Panels
	"GroupCompactionRate": promql.SumByRate(
		"thanos_compact_group_compactions_total",
		[]string{"namespace", "job", "resolution"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"GroupCompactionErrors": promqlbuilder.Mul(promqlbuilder.Div(
		promql.SumByRate(
			"thanos_compact_group_compactions_failures_total",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
		promql.SumByRate(
			"thanos_compact_group_compactions_total",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
		&parser.NumberLiteral{Val: 100},
	),
	"DownsampleRate": promql.SumByRate(
		"thanos_compact_downsample_total",
		[]string{"namespace", "job", "resolution"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"DownsampleErrors": promqlbuilder.Mul(promqlbuilder.Div(
		promql.SumByRate(
			"thanos_compact_downsample_failed_total",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
		promql.SumByRate(
			"thanos_compact_downsample_total",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
		&parser.NumberLiteral{Val: 100},
	),
	"DownsampleDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"thanos_compact_downsample_duration_seconds_bucket",
			[]string{"namespace", "job", "resolution", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"DownsampleDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"thanos_compact_downsample_duration_seconds_bucket",
			[]string{"namespace", "job", "resolution", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"DownsampleDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"thanos_compact_downsample_duration_seconds_bucket",
			[]string{"namespace", "job", "resolution", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"SyncMetaRate": promql.SumByRate(
		"thanos_blocks_meta_syncs_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"SyncMetaErrors": promqlbuilder.Mul(promqlbuilder.Div(
		promql.SumByRate(
			"thanos_blocks_meta_sync_failures_total",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
		promql.SumByRate(
			"thanos_blocks_meta_syncs_total",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
		&parser.NumberLiteral{Val: 100},
	),
	"SyncMetaDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"thanos_blocks_meta_sync_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"SyncMetaDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"thanos_blocks_meta_sync_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"SyncMetaDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"thanos_blocks_meta_sync_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"DeletionRate": promql.SumByRate(
		"thanos_compact_blocks_cleaned_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"DeletionErrors": promql.SumByRate(
		"thanos_compact_block_cleanup_failures_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"MarkingRate": promql.SumByRate(
		"thanos_compact_blocks_marked_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("marker").Equal("deletion-mark.json"),
	),
	"GarbageCollectionRate": promql.SumByRate(
		"thanos_compact_garbage_collection_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"GarbageCollectionErrors": promqlbuilder.Mul(promqlbuilder.Div(
		promql.SumByRate(
			"thanos_compact_garbage_collection_failures_total",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
		promql.SumByRate(
			"thanos_compact_garbage_collection_total",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
		&parser.NumberLiteral{Val: 100},
	),
	"GarbageCollectionDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"thanos_compact_garbage_collection_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"GarbageCollectionDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"thanos_compact_garbage_collection_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"GarbageCollectionDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"thanos_compact_garbage_collection_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"TodoCompactionBlocks": promqlbuilder.Sum(
		vector.New(
			vector.WithMetricName("thanos_compact_todo_compaction_blocks"),
			vector.WithLabelMatchers(
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
			),
		),
	).By("namespace", "job"),
	"TodoCompactions": promqlbuilder.Sum(
		vector.New(
			vector.WithMetricName("thanos_compact_todo_compactions"),
			vector.WithLabelMatchers(
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
			),
		),
	).By("namespace", "job"),
	"TodoDeletions": promqlbuilder.Sum(
		vector.New(
			vector.WithMetricName("thanos_compact_todo_deletion_blocks"),
			vector.WithLabelMatchers(
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
			),
		),
	).By("namespace", "job"),
	"TodoDownsamples": promqlbuilder.Sum(
		vector.New(
			vector.WithMetricName("thanos_compact_todo_downsample_blocks"),
			vector.WithLabelMatchers(
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
			),
		),
	).By("namespace", "job"),
	"HaltedCompactors": promqlbuilder.Sum(
		vector.New(
			vector.WithMetricName("thanos_compact_halted"),
			vector.WithLabelMatchers(
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
			),
		),
	).By("namespace", "job"),
}

// OverrideThanosPanelQueries overrides the ThanosPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
func OverrideThanosPanelQueries(queries map[string]parser.Expr) {
	for k, v := range queries {
		ThanosCommonPanelQueries[k] = v
	}
}
