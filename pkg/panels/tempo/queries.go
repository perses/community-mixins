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

package tempo

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

var clusterMatcher = label.New("cluster").EqualRegexp("$cluster")

func jobMatcher(job string) *labels.Matcher {
	return label.New("job").EqualRegexp(job)
}

// toLabelMatchersV1 adapts *labels.Matcher (used by every other panel query, backed by
// SetLabelMatchersV2) to the deprecated promql.LabelMatcher shape expected by the string-based
// promql.SetLabelMatchers, needed for WritesEnvoyProxyQPSQuery
func toLabelMatchersV1(matchers []*labels.Matcher) []promql.LabelMatcher {
	converted := make([]promql.LabelMatcher, len(matchers))
	for i, m := range matchers {
		converted[i] = promql.LabelMatcher{Name: m.Name, Value: m.Value, Type: m.Type.String()}
	}
	return converted
}

// statusFromCode builds the classic `sum by (status) (label_replace(label_replace(rate(<metric>{...}[$__rate_interval]), ...)))`
// pattern used to turn a numeric status_code label into a class like "2xx"/"5xx".
func statusFromCode(metricName string, matchers ...*labels.Matcher) parser.Expr {
	rate := promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName(metricName),
				vector.WithLabelMatchers(matchers...),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	)
	numeric := promqlbuilder.LabelReplace(rate, "status", "${1}xx", "status_code", "([0-9])..")
	alpha := promqlbuilder.LabelReplace(numeric, "status", "${1}", "status_code", "([a-zA-Z]+)")
	return promqlbuilder.Sum(alpha).By("status")
}

// latencyQuantile builds `histogram_quantile(<quantile>, sum(rate(<metric>_bucket{...}[$__rate_interval])) by (le,)) * 1e3`.
func latencyQuantile(quantile float64, bucketMetricName string, matchers ...*labels.Matcher) parser.Expr {
	return promqlbuilder.Mul(
		promqlbuilder.HistogramQuantile(
			quantile,
			promql.SumByRate(bucketMetricName, []string{"le"}, matchers...),
		),
		promqlbuilder.NewNumber(1e3),
	)
}

// latencyAverage builds `sum(rate(<metric>_sum{...})) by () * 1e3 / sum(rate(<metric>_count{...})) by ()`.
func latencyAverage(sumMetricName, countMetricName string, matchers ...*labels.Matcher) parser.Expr {
	return promqlbuilder.Div(
		promqlbuilder.Mul(
			promql.SumByRate(sumMetricName, []string{}, matchers...),
			promqlbuilder.NewNumber(1e3),
		),
		promql.SumByRate(countMetricName, []string{}, matchers...),
	)
}

const WritesEnvoyProxyQPSQuery = `sum by (grpc_status) (
    rate(
        label_replace(
            {cluster=~"$cluster", job=~"($namespace)/cortex-gw(-internal)?", __name__=~"envoy_cluster_grpc_proto_collector_trace_v1_TraceService_[0-9]+"},
            "grpc_status", "$1", "__name__", "envoy_cluster_grpc_proto_collector_trace_v1_TraceService_(.+)"
        )
        [$__rate_interval:30s]
    )
)
`

var TempoCommonPanelQueries = map[string]parser.Expr{
	// Writes / Gateway
	"WritesGatewayQPS": statusFromCode(
		"tempo_request_duration_seconds_count",
		clusterMatcher,
		jobMatcher("($namespace)/cortex-gw(-internal)?"),
		label.New("route").EqualRegexp("(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)"),
	),
	"WritesGatewayLatency_p99": latencyQuantile(0.99, "tempo_request_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/cortex-gw(-internal)?"),
		label.New("route").EqualRegexp("(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)"),
	),
	"WritesGatewayLatency_p50": latencyQuantile(0.50, "tempo_request_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/cortex-gw(-internal)?"),
		label.New("route").EqualRegexp("(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)"),
	),
	"WritesGatewayLatency_avg": latencyAverage("tempo_request_duration_seconds_sum", "tempo_request_duration_seconds_count",
		clusterMatcher,
		jobMatcher("($namespace)/cortex-gw(-internal)?"),
		label.New("route").EqualRegexp("(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)"),
	),

	// Writes / Distributor
	"WritesDistributorSpansSecond_accepted": promql.SumRate(
		"tempo_receiver_accepted_spans",
		clusterMatcher,
		jobMatcher("($namespace)/distributor"),
	),
	"WritesDistributorSpansSecond_refused": promql.SumRate(
		"tempo_receiver_refused_spans",
		clusterMatcher,
		jobMatcher("($namespace)/distributor"),
	),
	"WritesDistributorBytesPerSecond": promql.SumByRate(
		"tempo_distributor_bytes_received_total",
		[]string{"status"},
		clusterMatcher,
		jobMatcher("($namespace)/distributor"),
	),
	"WritesDistributorLatency_p99": latencyQuantile(0.99, "tempo_distributor_push_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/distributor"),
	),
	"WritesDistributorLatency_p50": latencyQuantile(0.5, "tempo_distributor_push_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/distributor"),
	),
	"WritesDistributorLatency_avg": latencyAverage("tempo_distributor_push_duration_seconds_sum", "tempo_distributor_push_duration_seconds_count",
		clusterMatcher,
		jobMatcher("($namespace)/distributor"),
	),
	"WritesDistributorKafkaAppendRecords": promql.SumRate(
		"tempo_distributor_kafka_appends_total",
		clusterMatcher,
		jobMatcher("($namespace)/distributor"),
		label.New("status").Equal("success"),
	),
	"WritesDistributorKafkaAppendFail": promql.SumRate(
		"tempo_distributor_kafka_appends_total",
		clusterMatcher,
		jobMatcher("($namespace)/distributor"),
		label.New("status").Equal("fail"),
	),
	"WritesDistributorKafkaWrite": promql.SumRate(
		"tempo_distributor_kafka_write_bytes_total",
		clusterMatcher,
		jobMatcher("($namespace)/distributor"),
	),
	"WritesDistributorKafkaWriteLatency_p50": promqlbuilder.HistogramQuantile(
		0.50,
		promql.SumByRate("tempo_distributor_kafka_write_latency_seconds_bucket", []string{"le"},
			clusterMatcher,
			jobMatcher("($namespace)/distributor"),
		),
	),
	"WritesDistributorKafkaWriteLatency_p99": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate("tempo_distributor_kafka_write_latency_seconds_bucket", []string{"le"},
			clusterMatcher,
			jobMatcher("($namespace)/distributor"),
		),
	),
	"WritesDistributorKafkaWriteLatency_avg": promqlbuilder.Div(
		promql.SumRate("tempo_distributor_kafka_write_latency_seconds_sum",
			clusterMatcher,
			jobMatcher("($namespace)/distributor"),
		),
		promql.SumRate("tempo_distributor_kafka_write_latency_seconds_count",
			clusterMatcher,
			jobMatcher("($namespace)/distributor"),
		),
	),

	// Writes / Ingester
	"WritesIngesterQPS": statusFromCode(
		"tempo_request_duration_seconds_count",
		clusterMatcher,
		jobMatcher("($namespace)/ingester"),
		label.New("route").EqualRegexp("/tempopb.Pusher/Push.*"),
	),
	"WritesIngesterLatency_p99": latencyQuantile(0.99, "tempo_request_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/ingester"),
		label.New("route").EqualRegexp("/tempopb.Pusher/Push.*"),
	),
	"WritesIngesterLatency_p50": latencyQuantile(0.50, "tempo_request_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/ingester"),
		label.New("route").EqualRegexp("/tempopb.Pusher/Push.*"),
	),
	"WritesIngesterLatency_avg": latencyAverage("tempo_request_duration_seconds_sum", "tempo_request_duration_seconds_count",
		clusterMatcher,
		jobMatcher("($namespace)/ingester"),
		label.New("route").EqualRegexp("/tempopb.Pusher/Push.*"),
	),

	// Writes / Memcached Ingester
	"WritesMemcachedIngesterQPS": statusFromCode(
		"tempo_memcache_request_duration_seconds_count",
		clusterMatcher,
		jobMatcher("($namespace)/ingester"),
		label.New("method").Equal("Memcache.Put"),
	),
	"WritesMemcachedIngesterLatency_p99": latencyQuantile(0.99, "tempo_memcache_request_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/ingester"),
		label.New("method").Equal("Memcache.Put"),
	),
	"WritesMemcachedIngesterLatency_p50": latencyQuantile(0.50, "tempo_memcache_request_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/ingester"),
		label.New("method").Equal("Memcache.Put"),
	),
	"WritesMemcachedIngesterLatency_avg": latencyAverage("tempo_memcache_request_duration_seconds_sum", "tempo_memcache_request_duration_seconds_count",
		clusterMatcher,
		jobMatcher("($namespace)/ingester"),
		label.New("method").Equal("Memcache.Put"),
	),

	// Writes / Backend Ingester
	"WritesBackendIngesterQPS": statusFromCode(
		"tempodb_backend_request_duration_seconds_count",
		clusterMatcher,
		jobMatcher("($namespace)/ingester"),
		label.New("operation").EqualRegexp("(PUT|POST)"),
	),
	// NB: kept identical to the pre-existing (buggy) dashboard, both series use the 99th percentile.
	"WritesBackendIngesterLatency_p99": latencyQuantile(0.99, "tempodb_backend_request_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/ingester"),
		label.New("operation").EqualRegexp("(PUT|POST)"),
	),
	"WritesBackendIngesterLatency_p50": latencyQuantile(0.99, "tempodb_backend_request_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/ingester"),
		label.New("operation").EqualRegexp("(PUT|POST)"),
	),
	"WritesBackendIngesterLatency_avg": latencyAverage("tempodb_backend_request_duration_seconds_sum", "tempodb_backend_request_duration_seconds_count",
		clusterMatcher,
		jobMatcher("($namespace)/ingester"),
		label.New("operation").EqualRegexp("(PUT|POST)"),
	),

	// Writes / Memcached Compactor
	"WritesMemcachedCompactorQPS": statusFromCode(
		"tempo_memcache_request_duration_seconds_count",
		clusterMatcher,
		jobMatcher("($namespace)/compactor"),
		label.New("method").Equal("Memcache.Put"),
	),
	"WritesMemcachedCompactorLatency_p99": latencyQuantile(0.99, "tempo_memcache_request_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/compactor"),
		label.New("method").Equal("Memcache.Put"),
	),
	"WritesMemcachedCompactorLatency_p50": latencyQuantile(0.50, "tempo_memcache_request_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/compactor"),
		label.New("method").Equal("Memcache.Put"),
	),
	"WritesMemcachedCompactorLatency_avg": latencyAverage("tempo_memcache_request_duration_seconds_sum", "tempo_memcache_request_duration_seconds_count",
		clusterMatcher,
		jobMatcher("($namespace)/compactor"),
		label.New("method").Equal("Memcache.Put"),
	),

	// Writes / Backend Compactor
	"WritesBackendCompactorQPS": statusFromCode(
		"tempodb_backend_request_duration_seconds_count",
		clusterMatcher,
		jobMatcher("($namespace)/compactor"),
		label.New("operation").EqualRegexp("(PUT|POST)"),
	),
	// NB: kept identical to the pre-existing (buggy) dashboard, both series use the 99th percentile.
	"WritesBackendCompactorLatency_p99": latencyQuantile(0.99, "tempodb_backend_request_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/compactor"),
		label.New("operation").EqualRegexp("(PUT|POST)"),
	),
	"WritesBackendCompactorLatency_p50": latencyQuantile(0.99, "tempodb_backend_request_duration_seconds_bucket",
		clusterMatcher,
		jobMatcher("($namespace)/compactor"),
		label.New("operation").EqualRegexp("(PUT|POST)"),
	),
	"WritesBackendCompactorLatency_avg": latencyAverage("tempodb_backend_request_duration_seconds_sum", "tempodb_backend_request_duration_seconds_count",
		clusterMatcher,
		jobMatcher("($namespace)/compactor"),
		label.New("operation").EqualRegexp("(PUT|POST)"),
	),

	// Tenants
	"TenantInfo": promqlbuilder.Max(
		promqlbuilder.Or(
			promql.MaxBy("tempo_limits_overrides", []string{"cluster", "namespace", "limit_name"},
				clusterMatcher,
				jobMatcher("($namespace)/compactor"),
				label.New("user").Equal("$tenant"),
			),
			promql.MaxBy("tempo_limits_defaults", []string{"cluster", "namespace", "limit_name"},
				clusterMatcher,
				jobMatcher("($namespace)/compactor"),
			),
		),
	).By("limit_name"),

	"TenantDistributorBytes_received": promql.SumRate(
		"tempo_distributor_bytes_received_total",
		clusterMatcher,
		jobMatcher("($namespace)/distributor"),
		label.New("tenant").Equal("$tenant"),
	),
	"TenantDistributorBytes_limit": promqlbuilder.Max(
		promqlbuilder.Or(
			promql.MaxBy("tempo_limits_overrides", []string{"cluster", "namespace", "limit_name"},
				clusterMatcher,
				jobMatcher("($namespace)/compactor"),
				label.New("user").Equal("$tenant"),
				label.New("limit_name").Equal("ingestion_rate_limit_bytes"),
			),
			promql.MaxBy("tempo_limits_defaults", []string{"cluster", "namespace", "limit_name"},
				clusterMatcher,
				jobMatcher("($namespace)/compactor"),
				label.New("limit_name").Equal("ingestion_rate_limit_bytes"),
			),
		),
	).By("ingestion_rate_limit_bytes"),
	"TenantDistributorBytes_burstLimit": promqlbuilder.Max(
		promqlbuilder.Or(
			promql.MaxBy("tempo_limits_overrides", []string{"cluster", "namespace", "limit_name"},
				clusterMatcher,
				jobMatcher("($namespace)/compactor"),
				label.New("user").Equal("$tenant"),
				label.New("limit_name").Equal("ingestion_burst_size_bytes"),
			),
			promql.MaxBy("tempo_limits_defaults", []string{"cluster", "namespace", "limit_name"},
				clusterMatcher,
				jobMatcher("($namespace)/compactor"),
				label.New("limit_name").Equal("ingestion_burst_size_bytes"),
			),
		),
	).By("ingestion_burst_size_bytes"),

	"TenantDistributorSpan_accepted": promql.SumRate(
		"tempo_distributor_spans_received_total",
		clusterMatcher,
		jobMatcher("($namespace)/distributor"),
		label.New("tenant").Equal("$tenant"),
	),
	"TenantDistributorSpan_refused": promql.SumByRate(
		"tempo_discarded_spans_total",
		[]string{"reason"},
		clusterMatcher,
		jobMatcher("($namespace)/distributor"),
		label.New("tenant").Equal("$tenant"),
	),

	"TenantLiveTraces_liveTraces": promqlbuilder.Max(
		vector.New(
			vector.WithMetricName("tempo_ingester_live_traces"),
			vector.WithLabelMatchers(
				clusterMatcher,
				jobMatcher("($namespace)/ingester"),
				label.New("tenant").Equal("$tenant"),
			),
		),
	),
	"TenantLiveTraces_globalLimit": promqlbuilder.Max(
		promqlbuilder.Or(
			promql.MaxBy("tempo_limits_overrides", []string{"cluster", "namespace", "limit_name"},
				clusterMatcher,
				jobMatcher("($namespace)/compactor"),
				label.New("user").Equal("$tenant"),
				label.New("limit_name").Equal("max_global_traces_per_user"),
			),
			promql.MaxBy("tempo_limits_defaults", []string{"cluster", "namespace", "limit_name"},
				clusterMatcher,
				jobMatcher("($namespace)/compactor"),
				label.New("limit_name").Equal("max_global_traces_per_user"),
			),
		),
	).By("max_global_traces_per_user"),
	"TenantLiveTraces_localLimit": promqlbuilder.Max(
		promqlbuilder.Or(
			promql.MaxBy("tempo_limits_overrides", []string{"cluster", "namespace", "limit_name"},
				clusterMatcher,
				jobMatcher("($namespace)/compactor"),
				label.New("user").Equal("$tenant"),
				label.New("limit_name").Equal("max_local_traces_per_user"),
			),
			promql.MaxBy("tempo_limits_defaults", []string{"cluster", "namespace", "limit_name"},
				clusterMatcher,
				jobMatcher("($namespace)/compactor"),
				label.New("limit_name").Equal("max_local_traces_per_user"),
			),
		),
	).By("max_local_traces_per_user"),

	"TenantQueriesID": promql.SumByRate(
		"tempo_query_frontend_queries_total",
		[]string{"status"},
		clusterMatcher,
		jobMatcher("($namespace)/query-frontend"),
		label.New("tenant").Equal("$tenant"),
		label.New("op").Equal("traces"),
	),
	"TenantQueriesSearch": promql.SumByRate(
		"tempo_query_frontend_queries_total",
		[]string{"status"},
		clusterMatcher,
		jobMatcher("($namespace)/query-frontend"),
		label.New("tenant").Equal("$tenant"),
		label.New("op").Equal("search"),
	),

	"TenantBlockslistLength": promqlbuilder.Avg(
		vector.New(
			vector.WithMetricName("tempodb_blocklist_length"),
			vector.WithLabelMatchers(
				clusterMatcher,
				jobMatcher("($namespace)/compactor"),
				label.New("tenant").Equal("$tenant"),
			),
		),
	),

	"TenantOutstandingCompactions": promqlbuilder.Div(
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("tempodb_compaction_outstanding_blocks"),
				vector.WithLabelMatchers(
					clusterMatcher,
					jobMatcher("($namespace)/compactor"),
					label.New("tenant").Equal("$tenant"),
				),
			),
		),
		promqlbuilder.Count(
			vector.New(
				vector.WithMetricName("tempo_build_info"),
				vector.WithLabelMatchers(
					clusterMatcher,
					jobMatcher("($namespace)/compactor"),
				),
			),
		),
	),

	"TenantMetricGeneratorBytes": promql.SumRate(
		"tempo_metrics_generator_bytes_received_total",
		clusterMatcher,
		jobMatcher("($namespace)/metrics-generator"),
		label.New("tenant").Equal("$tenant"),
	),

	"TenantMetricGeneratorActiveSeries_series": promqlbuilder.Sum(
		vector.New(
			vector.WithMetricName("tempo_metrics_generator_registry_active_series"),
			vector.WithLabelMatchers(
				clusterMatcher,
				jobMatcher("($namespace)/metrics-generator"),
				label.New("tenant").Equal("$tenant"),
			),
		),
	),
	"TenantMetricGeneratorActiveSeries_limit": promqlbuilder.Max(
		promqlbuilder.Or(
			promql.MaxBy("tempo_limits_overrides", []string{"cluster", "namespace", "limit_name"},
				clusterMatcher,
				jobMatcher("($namespace)/compactor"),
				label.New("user").Equal("$tenant"),
				label.New("limit_name").Equal("metrics_generator_max_active_series"),
			),
			promql.MaxBy("tempo_limits_defaults", []string{"cluster", "namespace", "limit_name"},
				clusterMatcher,
				jobMatcher("($namespace)/compactor"),
				label.New("limit_name").Equal("metrics_generator_max_active_series"),
			),
		),
	).By("metrics_generator_max_active_series"),
}

// OverrideTempoPanelQueries overrides the TempoCommonPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
// The convention of naming followed, is to use Panel function name (with _suffix, in case panel has multiple queries)
func OverrideTempoPanelQueries(queries map[string]parser.Expr) {
	maps.Copy(TempoCommonPanelQueries, queries)
}
