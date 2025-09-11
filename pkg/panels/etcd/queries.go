package etcd

import (
	"maps"

	"github.com/perses/community-dashboards/pkg/promql"
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/matrix"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/promql/parser"
)

var EtcdCommonPanelQueries = map[string]parser.Expr{
	"EtcdUpStatus": promqlbuilder.Sum(
		vector.New(
			vector.WithMetricName("etcd_server_has_leader"),
			vector.WithLabelMatchers(
				label.New("job").EqualRegexp(".*etcd.*"),
				label.New("cluster").Equal("$cluster"),
			),
		),
	),
	"EtcdgRPCRateStarted": promqlbuilder.Sum(
		promqlbuilder.Rate(
			matrix.New(
				vector.New(
					vector.WithMetricName("grpc_server_started_total"),
					vector.WithLabelMatchers(
						label.New("job").EqualRegexp(".*etcd.*"),
						label.New("cluster").Equal("$cluster"),
						label.New("grpc_type").Equal("unary"),
					),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
	),
	"EtcdgRPCRateTotal": promql.SumRate(
		"grpc_server_handled_total",
		label.New("job").EqualRegexp(".*etcd.*"),
		label.New("cluster").Equal("$cluster"),
		label.New("grpc_type").Equal("unary"),
		label.New("grpc_code").EqualRegexp("Unknown|FailedPrecondition|ResourceExhausted|Internal|Unavailable|DataLoss|DeadlineExceeded"),
	),
	"EtcdActiveStreamsWatch": promqlbuilder.Sub(
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("grpc_server_started_total"),
				vector.WithLabelMatchers(
					label.New("job").EqualRegexp(".*etcd.*"),
					label.New("cluster").Equal("$cluster"),
					label.New("grpc_service").Equal("etcdserverpb.Watch"),
					label.New("grpc_type").Equal("bidi_stream"),
				),
			),
		),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("grpc_server_handled_total"),
				vector.WithLabelMatchers(
					label.New("cluster").Equal("$cluster"),
					label.New("grpc_service").Equal("etcdserverpb.Watch"),
					label.New("grpc_type").Equal("bidi_stream"),
				),
			),
		),
	),
	"EtcdActiveStreamsLease": promqlbuilder.Sub(
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("grpc_server_started_total"),
				vector.WithLabelMatchers(
					label.New("job").EqualRegexp(".*etcd.*"),
					label.New("cluster").Equal("$cluster"),
					label.New("grpc_service").Equal("etcdserverpb.Lease"),
					label.New("grpc_type").Equal("bidi_stream"),
				),
			),
		),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("grpc_server_handled_total"),
				vector.WithLabelMatchers(
					label.New("cluster").Equal("$cluster"),
					label.New("grpc_service").Equal("etcdserverpb.Lease"),
					label.New("grpc_type").Equal("bidi_stream"),
				),
			),
		),
	),
	"EtcdDBSize": vector.New(
		vector.WithMetricName("etcd_mvcc_db_total_size_in_bytes"),
		vector.WithLabelMatchers(
			label.New("job").EqualRegexp(".*etcd.*"),
			label.New("cluster").Equal("$cluster"),
		),
	),
	"EtcdDiskSyncWalFsyncDuration": promqlbuilder.HistogramQuantile(0.99,
		promql.SumByRate(
			"etcd_disk_wal_fsync_duration_seconds_bucket",
			[]string{"instance", "le"},
			label.New("job").EqualRegexp(".*etcd.*"),
			label.New("cluster").Equal("$cluster"),
		),
	),
	"EtcdDiskSyncBackendDuration": promqlbuilder.HistogramQuantile(0.99,
		promql.SumByRate(
			"etcd_disk_backend_commit_duration_seconds_bucket",
			[]string{"instance", "le"},
			label.New("job").EqualRegexp(".*etcd.*"),
			label.New("cluster").Equal("$cluster"),
		),
	),
	"EtcdClientTrafficIn": promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName("etcd_network_client_grpc_received_bytes_total"),
				vector.WithLabelMatchers(
					label.New("job").EqualRegexp(".*etcd.*"),
					label.New("cluster").Equal("$cluster"),
				),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	),
	"EtcdClientTrafficOut": promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName("etcd_network_client_grpc_sent_bytes_total"),
				vector.WithLabelMatchers(
					label.New("job").EqualRegexp(".*etcd.*"),
					label.New("cluster").Equal("$cluster"),
				),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	),
	"EtcdPeerTrafficIn": promql.SumByRate(
		"etcd_network_peer_received_bytes_total",
		[]string{"instance"},
		label.New("job").EqualRegexp(".*etcd.*"),
		label.New("cluster").Equal("$cluster"),
	),
	"EtcdPeerTrafficOut": promql.SumByRate(
		"etcd_network_peer_sent_bytes_total",
		[]string{"instance"},
		label.New("job").EqualRegexp(".*etcd.*"),
		label.New("cluster").Equal("$cluster"),
	),
	"EtcdRaftProposals": promqlbuilder.Changes(
		matrix.New(
			vector.New(
				vector.WithMetricName("etcd_server_leader_changes_seen_total"),
				vector.WithLabelMatchers(
					label.New("job").EqualRegexp(".*etcd.*"),
					label.New("cluster").Equal("$cluster"),
				),
			),
			matrix.WithRangeAsString("1d"),
		),
	),
	"EtcdPeerRoundtripTime": promqlbuilder.HistogramQuantile(0.99,
		promql.SumByRate(
			"etcd_network_peer_round_trip_time_seconds_bucket",
			[]string{"instance", "le"},
			label.New("job").EqualRegexp(".*etcd.*"),
			label.New("cluster").Equal("$cluster"),
		),
	),
}

// OverrideEtcdPanelQueries overrides the EtcdCommonPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
// The convention of naming followed, is to use Panel function name (with _suffix, in case panel has multiple queries)
func OverrideEtcdPanelQueries(queries map[string]parser.Expr) {
	maps.Copy(EtcdCommonPanelQueries, queries)
}
