// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package thanos

import (
	"maps"

	"github.com/perses/community-mixins/pkg/promql"
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/matrix"
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

	// Thanos Ruler Panels
	"RuleEvaluationRate": promql.SumByRate(
		"prometheus_rule_evaluations_total",
		[]string{"namespace", "job", "rule_group", "strategy"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"RuleEvaluationFailureRate": promql.SumByRate(
		"prometheus_rule_evaluation_failures_total",
		[]string{"namespace", "job", "rule_group", "strategy"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"RuleGroupEvaluationsMissRate": promql.SumByRate(
		"prometheus_rule_group_iterations_missed_total",
		[]string{"namespace", "job", "rule_group", "strategy"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"RuleGroupEvaluationsTooSlow": promqlbuilder.Gtr(
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("prometheus_rule_group_last_duration_seconds"),
				vector.WithLabelMatchers(
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
				),
			),
		).By("namespace", "job", "rule_group"),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("prometheus_rule_group_interval_seconds"),
				vector.WithLabelMatchers(
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
				),
			),
		).By("namespace", "job", "rule_group"),
	),
	"AlertsDroppedRate": promql.SumByRate(
		"thanos_alert_sender_alerts_dropped_total",
		[]string{"namespace", "job", "alertmanager"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"AlertsSentRate": promql.SumByRate(
		"thanos_alert_sender_alerts_sent_total",
		[]string{"namespace", "job", "alertmanager"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"AlertSendingErrors": promql.ErrorCaseRatio(
		"thanos_alert_sender_errors_total",
		[]string{"namespace", "job"},
		[]*labels.Matcher{
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		},
		"thanos_alert_sender_alerts_sent_total",
		[]string{"namespace", "job"},
		[]*labels.Matcher{
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		},
	),
	"AlertSendingDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"thanos_alert_sender_latency_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"AlertSendingDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"thanos_alert_sender_latency_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"AlertSendingDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"thanos_alert_sender_latency_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"AlertQueuePushedRate": promql.SumByRate(
		"thanos_alert_queue_alerts_pushed_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"AlertQueuePoppedRate": promql.SumByRate(
		"thanos_alert_queue_alerts_popped_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"DroppedRatio": promqlbuilder.Div(
		promql.SumByRate(
			"thanos_alert_queue_alerts_dropped_total",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
		promql.SumByRate(
			"thanos_alert_queue_alerts_pushed_total",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),

	// Thanos Query Panels
	"InstantQueryRequestRate": promql.SumByRate(
		"http_requests_total",
		[]string{"namespace", "job", "handler", "code"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("handler").Equal("query"),
	),
	"InstantQueryRequestErrors": promqlbuilder.Mul(
		promql.IgnoringGroupLeft(
			promql.ErrorCaseRatio(
				"http_requests_total",
				[]string{"namespace", "job", "code"},
				[]*labels.Matcher{
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
					label.New("handler").Equal("query"),
					label.New("code").EqualRegexp("5.."),
				},
				"http_requests_total",
				[]string{"namespace", "job"},
				[]*labels.Matcher{
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
					label.New("handler").Equal("query"),
				},
			),
			[]string{"code"},
		),
		&parser.NumberLiteral{Val: 100},
	),
	"InstantQueryRequestDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"http_request_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("handler").Equal("query"),
		),
	),
	"InstantQueryRequestDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"http_request_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("handler").Equal("query"),
		),
	),
	"InstantQueryRequestDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"http_request_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("handler").Equal("query"),
		),
	),
	"RangeQueryRequestRate": promql.SumByRate(
		"http_requests_total",
		[]string{"namespace", "job", "handler", "code"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("handler").Equal("query_range"),
	),
	"RangeQueryRequestErrors": promqlbuilder.Mul(
		promql.IgnoringGroupLeft(
			promql.ErrorCaseRatio(
				"http_requests_total",
				[]string{"namespace", "job", "code"},
				[]*labels.Matcher{
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
					label.New("handler").Equal("query_range"),
					label.New("code").EqualRegexp("5.."),
				},
				"http_requests_total",
				[]string{"namespace", "job"},
				[]*labels.Matcher{
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
					label.New("handler").Equal("query_range"),
				},
			),
			[]string{"code"},
		),
		&parser.NumberLiteral{Val: 100},
	),
	"RangeQueryRequestDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"http_request_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("handler").Equal("query_range"),
		),
	),
	"RangeQueryRequestDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"http_request_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("handler").Equal("query_range"),
		),
	),
	"RangeQueryRequestDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"http_request_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("handler").Equal("query_range"),
		),
	),
	"QueryConcurrency": promqlbuilder.Sub(
		promqlbuilder.MaxOverTime(
			matrix.New(
				vector.New(
					vector.WithMetricName("thanos_query_concurrent_gate_queries_max"),
					vector.WithLabelMatchers(
						label.New("namespace").Equal("$namespace"),
						label.New("job").EqualRegexp("$job"),
					),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
		promqlbuilder.AvgOverTime(
			matrix.New(
				vector.New(
					vector.WithMetricName("thanos_query_concurrent_gate_queries_in_flight"),
					vector.WithLabelMatchers(
						label.New("namespace").Equal("$namespace"),
						label.New("job").EqualRegexp("$job"),
					),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
	),
	"DNSLookups": promql.SumByRate(
		"thanos_query_store_apis_dns_lookups_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"DNSLookupsErrors": promql.ErrorCaseRatio(
		"thanos_query_store_apis_dns_failures_total",
		[]string{"namespace", "job"},
		[]*labels.Matcher{
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		},
		"thanos_query_store_apis_dns_lookups_total",
		[]string{"namespace", "job"},
		[]*labels.Matcher{
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		},
	),

	// Thanos Query Frontend Panels
	"QueryFrontendRequestRate": promql.SumByRate(
		"http_requests_total",
		[]string{"namespace", "job", "handler", "code"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("handler").Equal("query-frontend"),
	),
	"QueryFrontendQueryRate": promql.SumByRate(
		"thanos_query_frontend_queries_total",
		[]string{"namespace", "job", "handler", "code"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("op").Equal("query_range"),
	),
	"QueryFrontendErrors": promqlbuilder.Mul(
		promql.IgnoringGroupLeft(
			promql.ErrorCaseRatio(
				"http_requests_total",
				[]string{"namespace", "job", "code"},
				[]*labels.Matcher{
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
					label.New("handler").Equal("query-frontend"),
					label.New("code").EqualRegexp("5.."),
				},
				"http_requests_total",
				[]string{"namespace", "job"},
				[]*labels.Matcher{
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
					label.New("handler").Equal("query-frontend"),
				},
			),
			[]string{"code"},
		),
		&parser.NumberLiteral{Val: 100},
	),
	"QueryFrontendDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"http_request_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("handler").Equal("query-frontend"),
		),
	),
	"QueryFrontendDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"http_request_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("handler").Equal("query-frontend"),
		),
	),
	"QueryFrontendDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"http_request_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("handler").Equal("query-frontend"),
		),
	),
	"QueryFrontendCacheRequestRate": promql.SumByRate(
		"cortex_cache_request_duration_seconds_count",
		[]string{"namespace", "job", "tripperware"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"QueryFrontendCacheHitRate": promql.SumByRate(
		"cortex_cache_hits_total",
		[]string{"namespace", "job", "tripperware"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"QueryFrontendCacheMissRate": promql.SumByRate(
		"querier_cache_misses_total",
		[]string{"namespace", "job", "tripperware"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"QueryFrontendFetchedKeyRate": promql.SumByRate(
		"cortex_cache_fetched_keys_total",
		[]string{"namespace", "job", "tripperware"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),

	// Thanos Receive Panels
	"RemoteWriteRequestRate": promql.SumByRate(
		"http_requests_total",
		[]string{"namespace", "job", "handler", "code"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("handler").Equal("receive"),
	),
	"RemoteWriteRequestErrors": promqlbuilder.Mul(
		promql.IgnoringGroupLeft(
			promql.ErrorCaseRatio(
				"http_requests_total",
				[]string{"namespace", "job", "code"},
				[]*labels.Matcher{
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
					label.New("handler").Equal("receive"),
					label.New("code").EqualRegexp("5.."),
				},
				"http_requests_total",
				[]string{"namespace", "job"},
				[]*labels.Matcher{
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
					label.New("handler").Equal("receive"),
				},
			),
			[]string{"code"},
		),
		&parser.NumberLiteral{Val: 100},
	),
	"RemoteWriteRequestDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"http_request_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("handler").Equal("receive"),
		),
	),
	"RemoteWriteRequestDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"http_request_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("handler").Equal("receive"),
		),
	),
	"RemoteWriteRequestDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"http_request_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("handler").Equal("receive"),
		),
	),
	"TenantedRemoteWriteRequestRate": promql.SumByRate(
		"http_requests_total",
		[]string{"tenant", "job", "handler", "code"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("tenant").EqualRegexp("$tenant"),
		label.New("handler").Equal("receive"),
	),
	"TenantedRemoteWriteRequestErrors": promqlbuilder.Mul(
		promql.IgnoringGroupLeft(
			promql.ErrorCaseRatio(
				"http_requests_total",
				[]string{"tenant", "namespace", "job", "code"},
				[]*labels.Matcher{
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
					label.New("tenant").EqualRegexp("$tenant"),
					label.New("handler").Equal("receive"),
					label.New("code").NotEqualRegexp("2.."),
				},
				"http_requests_total",
				[]string{"tenant", "namespace", "job"},
				[]*labels.Matcher{
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
					label.New("tenant").EqualRegexp("$tenant"),
					label.New("handler").Equal("receive"),
				},
			),
			[]string{"code"},
		),
		&parser.NumberLiteral{Val: 100},
	),
	"TenantedRemoteWriteRequestDurations": promqlbuilder.Div(
		promql.SumByRate(
			"http_request_duration_seconds_sum",
			[]string{"tenant", "namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("tenant").EqualRegexp("$tenant"),
			label.New("handler").Equal("receive"),
		),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("http_request_duration_seconds_count"),
				vector.WithLabelMatchers(
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
					label.New("tenant").EqualRegexp("$tenant"),
					label.New("handler").Equal("receive"),
				),
			),
		).By("tenant", "namespace", "job"),
	),
	"AvgRemoteWriteRequestSize": promqlbuilder.Div(
		promql.SumByRate(
			"http_request_size_bytes_sum",
			[]string{"tenant", "namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("tenant").EqualRegexp("$tenant"),
			label.New("handler").Equal("receive"),
			label.New("code").EqualRegexp("2.."),
		),
		promql.SumByRate(
			"http_request_size_bytes_count",
			[]string{"tenant", "namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("tenant").EqualRegexp("$tenant"),
			label.New("handler").Equal("receive"),
			label.New("code").EqualRegexp("2.."),
		),
	),
	"AvgFailedRemoteWriteRequestSize": promqlbuilder.Div(
		promql.SumByRate(
			"http_request_size_bytes_sum",
			[]string{"tenant", "namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("tenant").EqualRegexp("$tenant"),
			label.New("handler").Equal("receive"),
			label.New("code").NotEqualRegexp("2.."),
		),
		promql.SumByRate(
			"http_request_size_bytes_count",
			[]string{"tenant", "namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("tenant").EqualRegexp("$tenant"),
			label.New("handler").Equal("receive"),
			label.New("code").NotEqualRegexp("2.."),
		),
	),
	"InflightRemoteWriteRequests": promqlbuilder.Sum(
		vector.New(
			vector.WithMetricName("http_inflight_requests"),
			vector.WithLabelMatchers(
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
				label.New("tenant").EqualRegexp("$tenant"),
				label.New("handler").Equal("receive"),
			),
		),
	).By("tenant", "namespace", "job", "method"),
	"RemoteWriteSeriesRate": promql.SumByRate(
		"thanos_receive_write_timeseries_sum",
		[]string{"tenant", "namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("tenant").EqualRegexp("$tenant"),
		label.New("code").EqualRegexp("2.."),
	),
	"RemoteWriteSeriesNotWrittenRate": promql.SumByRate(
		"thanos_receive_write_timeseries_sum",
		[]string{"tenant", "namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("tenant").EqualRegexp("$tenant"),
		label.New("code").NotEqualRegexp("2.."),
	),
	"RemoteWriteSamplesRate": promql.SumByRate(
		"thanos_receive_write_samples_sum",
		[]string{"tenant", "namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("tenant").EqualRegexp("$tenant"),
		label.New("code").EqualRegexp("2.."),
	),
	"RemoteWriteSamplesNotWrittenRate": promql.SumByRate(
		"thanos_receive_write_samples_sum",
		[]string{"tenant", "namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("tenant").EqualRegexp("$tenant"),
		label.New("code").NotEqualRegexp("2.."),
	),
	"RemoteWriteReplicationRate": promql.SumByRate(
		"thanos_receive_replications_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"RemoteWriteReplicationErrorRate": promql.SumByRate(
		"thanos_receive_replications_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("result").Equal("error"),
	),
	"RemoteWriteForwardRate": promql.SumByRate(
		"thanos_receive_forward_requests_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"RemoteWriteForwardErrorRate": promql.SumByRate(
		"thanos_receive_forward_requests_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("result").Equal("error"),
	),
	"WriteGRPCUnaryRate": promql.SumByRate(
		"grpc_server_handled_total",
		[]string{"namespace", "job", "grpc_method", "grpc_code"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
		label.New("grpc_method").Equal("RemoteWrite"),
		label.New("grpc_type").Equal("unary"),
	),
	"WriteGRPCUnaryErrors": promql.IgnoringGroupLeft(
		promql.ErrorCaseRatio(
			"grpc_server_handled_total",
			[]string{"namespace", "job", "grpc_code"},
			[]*labels.Matcher{
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
				label.New("grpc_method").Equal("RemoteWrite"),
				label.New("grpc_type").Equal("unary"),
				label.New("grpc_code").EqualRegexp("Unknown|ResourceExhausted|Internal|Unavailable|DataLoss"),
			},
			"grpc_server_handled_total",
			[]string{"namespace", "job"},
			[]*labels.Matcher{
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
				label.New("grpc_method").Equal("RemoteWrite"),
				label.New("grpc_type").Equal("unary"),
			},
		),
		[]string{"grpc_code"},
	),
	"WriteGPRCUnaryDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"grpc_server_handling_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("grpc_method").Equal("RemoteWrite"),
			label.New("grpc_type").Equal("unary"),
		),
	),
	"WriteGPRCUnaryDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"grpc_server_handling_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("grpc_method").Equal("RemoteWrite"),
			label.New("grpc_type").Equal("unary"),
		),
	),
	"WriteGPRCUnaryDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"grpc_server_handling_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
			label.New("grpc_method").Equal("RemoteWrite"),
			label.New("grpc_type").Equal("unary"),
		),
	),
	"ReceiveAppendedSampleRate": promql.SumByRate(
		"prometheus_tsdb_head_samples_appended_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"ReceiveHeadSeries": promqlbuilder.Sum(
		vector.New(
			vector.WithMetricName("prometheus_tsdb_head_series"),
			vector.WithLabelMatchers(
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
			),
		),
	).By("namespace", "job"),
	"ReceiveHeadChunks": promqlbuilder.Sum(
		vector.New(
			vector.WithMetricName("prometheus_tsdb_head_chunks"),
			vector.WithLabelMatchers(
				label.New("namespace").Equal("$namespace"),
				label.New("job").EqualRegexp("$job"),
			),
		),
	).By("namespace", "job"),
	"BucketUploadTable_uploadedAgo": promqlbuilder.Sub(
		promqlbuilder.Time(),
		promqlbuilder.Max(
			vector.New(
				vector.WithMetricName("thanos_objstore_bucket_last_successful_upload_time"),
				vector.WithLabelMatchers(
					label.New("namespace").Equal("$namespace"),
					label.New("job").EqualRegexp("$job"),
				),
			),
		).By("namespace", "job", "bucket"),
	),

	// Thanos Store Panels
	"BlockLoadRate": promql.SumByRate(
		"thanos_bucket_store_block_loads_total",
		[]string{"namespace", "job"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"BlockLoadErrors": promql.ErrorCasePercentage(
		"thanos_bucket_store_block_load_failures_total",
		[]string{"namespace", "job"},
		[]*labels.Matcher{
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		},
		"thanos_bucket_store_block_loads_total",
		[]string{"namespace", "job"},
		[]*labels.Matcher{
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		},
	),
	"BlockDropRate": promql.SumByRate(
		"thanos_bucket_store_block_drops_total",
		[]string{"namespace", "job", "operation"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"BlockDropErrors": promql.ErrorCasePercentage(
		"thanos_bucket_store_block_drop_failures_total",
		[]string{"namespace", "job"},
		[]*labels.Matcher{
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		},
		"thanos_bucket_store_block_drops_total",
		[]string{"namespace", "job"},
		[]*labels.Matcher{
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		},
	),
	"CacheRequestRate": promql.SumByRate(
		"thanos_store_index_cache_requests_total",
		[]string{"namespace", "job", "item_type"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"CacheHitRate": promql.SumByRate(
		"thanos_store_index_cache_hits_total",
		[]string{"namespace", "job", "item_type"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"CacheItemsAddRate": promql.SumByRate(
		"thanos_store_index_cache_items_added_total",
		[]string{"namespace", "job", "item_type"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"CacheItemsEvictRate": promql.SumByRate(
		"thanos_store_index_cache_items_evicted_total",
		[]string{"namespace", "job", "item_type"},
		label.New("namespace").Equal("$namespace"),
		label.New("job").EqualRegexp("$job"),
	),
	"BlocksQueried_mean": promqlbuilder.Div(
		promql.SumByRate(
			"thanos_bucket_store_series_blocks_queried_sum",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
		promql.SumByRate(
			"thanos_bucket_store_series_blocks_queried_count",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"BlocksQueried_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"thanos_bucket_store_series_blocks_queried_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"BlocksQueried_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"thanos_bucket_store_series_blocks_queried_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"BlocksQueried_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"thanos_bucket_store_series_blocks_queried_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"DataFetched_mean": promqlbuilder.Div(
		promql.SumByRate(
			"thanos_bucket_store_series_data_size_fetched_bytes_sum",
			[]string{"namespace", "job", "data_type"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
		promql.SumByRate(
			"thanos_bucket_store_series_data_size_fetched_bytes_count",
			[]string{"namespace", "job", "data_type"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"DataFetched_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"thanos_bucket_store_series_data_size_fetched_bytes_bucket",
			[]string{"namespace", "job", "data_type", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"DataFetched_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"thanos_bucket_store_series_data_size_fetched_bytes_bucket",
			[]string{"namespace", "job", "data_type", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"DataFetched_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"thanos_bucket_store_series_data_size_fetched_bytes_bucket",
			[]string{"namespace", "job", "data_type", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"DataTouched_mean": promqlbuilder.Div(
		promql.SumByRate(
			"thanos_bucket_store_series_data_size_touched_bytes_sum",
			[]string{"namespace", "job", "data_type"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
		promql.SumByRate(
			"thanos_bucket_store_series_data_size_touched_bytes_count",
			[]string{"namespace", "job", "data_type"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"DataTouched_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"thanos_bucket_store_series_data_size_touched_bytes_bucket",
			[]string{"namespace", "job", "data_type", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"DataTouched_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"thanos_bucket_store_series_data_size_touched_bytes_bucket",
			[]string{"namespace", "job", "data_type", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"DataTouched_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"thanos_bucket_store_series_data_size_touched_bytes_bucket",
			[]string{"namespace", "job", "data_type", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"ResultSeries_mean": promqlbuilder.Div(
		promql.SumByRate(
			"thanos_bucket_store_series_result_series_sum",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
		promql.SumByRate(
			"thanos_bucket_store_series_result_series_count",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"ResultSeries_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"thanos_bucket_store_series_result_series_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"ResultSeries_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"thanos_bucket_store_series_result_series_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"ResultSeries_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"thanos_bucket_store_series_result_series_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"GetAllSeriesDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"thanos_bucket_store_series_get_all_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"GetAllSeriesDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"thanos_bucket_store_series_get_all_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"GetAllSeriesDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"thanos_bucket_store_series_get_all_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"MergeDurations_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"thanos_bucket_store_series_merge_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"MergeDurations_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"thanos_bucket_store_series_merge_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"MergeDurations_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"thanos_bucket_store_series_merge_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"GateWaitingDuration_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"thanos_bucket_store_series_gate_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"GateWaitingDuration_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"thanos_bucket_store_series_gate_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"GateWaitingDuration_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"thanos_bucket_store_series_gate_duration_seconds_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"StoreSentChunkSizes_mean": promqlbuilder.Div(
		promql.SumByRate(
			"thanos_bucket_store_sent_chunk_size_bytes_sum",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
		promql.SumByRate(
			"thanos_bucket_store_sent_chunk_size_bytes_count",
			[]string{"namespace", "job"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"StoreSentChunkSizes_50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate(
			"thanos_bucket_store_sent_chunk_size_bytes_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"StoreSentChunkSizes_90": promqlbuilder.HistogramQuantile(
		0.90,
		promql.SumByRate(
			"thanos_bucket_store_sent_chunk_size_bytes_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
	"StoreSentChunkSizes_99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"thanos_bucket_store_sent_chunk_size_bytes_bucket",
			[]string{"namespace", "job", "le"},
			label.New("namespace").Equal("$namespace"),
			label.New("job").EqualRegexp("$job"),
		),
	),
}

// OverrideThanosPanelQueries overrides the ThanosCommonPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
// The convention of naming followed, is to use Panel function name (with _suffix, in case panel has multiple queries)
func OverrideThanosPanelQueries(queries map[string]parser.Expr) {
	maps.Copy(ThanosCommonPanelQueries, queries)
}
