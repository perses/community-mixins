package nodeexporter

import (
	"maps"

	"github.com/perses/community-mixins/pkg/promql"
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/matrix"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/promql/parser"
)

var NodeExporterCommonPanelQueries = map[string]parser.Expr{
	"NodeExporterCPUUsagePercentage": promql.IgnoringGroupLeft(
		promqlbuilder.Div(
			promqlbuilder.Sub(
				&parser.NumberLiteral{Val: 1},
				promqlbuilder.Sum(
					promqlbuilder.Rate(
						matrix.New(
							vector.New(
								vector.WithMetricName("node_cpu_seconds_total"),
								vector.WithLabelMatchers(
									label.New("job").Equal("node"),
									label.New("mode").EqualRegexp("idle|iowait|steal"),
									label.New("instance").Equal("$instance"),
								),
							),
							matrix.WithRangeAsVariable("$__rate_interval"),
						),
					),
				).Without("mode"),
			),
			promqlbuilder.Count(
				vector.New(
					vector.WithMetricName("node_cpu_seconds_total"),
					vector.WithLabelMatchers(
						label.New("job").Equal("node"),
						label.New("mode").Equal("idle"),
						label.New("instance").Equal("$instance"),
					),
				),
			).Without("cpu", "mode"),
		),
		[]string{"cpu"},
	),
	"NodeExporterClusterNodeCPUUsagePercentage": promqlbuilder.Div(
		promqlbuilder.Neq(
			promqlbuilder.Mul(
				vector.New(
					vector.WithMetricName("instance:node_cpu_utilisation:rate5m"),
					vector.WithLabelMatchers(
						label.New("job").Equal("node"),
					),
				),
				vector.New(
					vector.WithMetricName("instance:node_num_cpu:sum"),
					vector.WithLabelMatchers(
						label.New("job").Equal("node"),
					),
				),
			),
			&parser.NumberLiteral{Val: 0},
		),
		promqlbuilder.Scalar(
			promqlbuilder.Sum(
				vector.New(
					vector.WithMetricName("instance:node_num_cpu:sum"),
					vector.WithLabelMatchers(
						label.New("job").Equal("node"),
						label.New("instance").EqualRegexp("$instance"),
					),
				),
			),
		),
	),
	"NodeExporterClusterNodeCPUSaturationPercentage": promqlbuilder.Neq(
		promqlbuilder.Div(
			vector.New(
				vector.WithMetricName("instance:node_load1_per_cpu:ratio"),
				vector.WithLabelMatchers(
					label.New("job").Equal("node"),
				),
			),
			promqlbuilder.Scalar(
				promqlbuilder.Count(
					vector.New(
						vector.WithMetricName("instance:node_load1_per_cpu:ratio"),
						vector.WithLabelMatchers(
							label.New("job").Equal("node"),
							label.New("instance").EqualRegexp("$instance"),
						),
					),
				),
			),
		),
		&parser.NumberLiteral{Val: 0},
	),
	"NodeExporterClusterNodeMemoryUsagePercentage": promqlbuilder.Neq(
		promqlbuilder.Div(
			vector.New(
				vector.WithMetricName("instance:node_memory_utilisation:ratio"),
				vector.WithLabelMatchers(
					label.New("job").Equal("node"),
				),
			),
			promqlbuilder.Scalar(
				promqlbuilder.Count(
					vector.New(
						vector.WithMetricName("instance:node_memory_utilisation:ratio"),
						vector.WithLabelMatchers(
							label.New("job").Equal("node"),
							label.New("instance").EqualRegexp("$instance"),
						),
					),
				),
			),
		),
		&parser.NumberLiteral{Val: 0},
	),
	"NodeExporterClusterNodeMemorySaturationPercentage": vector.New(
		vector.WithMetricName("instance:node_vmstat_pgmajfault:rate5m"),
		vector.WithLabelMatchers(
			label.New("job").Equal("node"),
		),
	),
	"NodeExporterClusterNodeDiskUsagePercentage": promqlbuilder.Neq(
		promqlbuilder.Div(
			vector.New(
				vector.WithMetricName("instance_device:node_disk_io_time_seconds:rate5m"),
				vector.WithLabelMatchers(
					label.New("job").Equal("node"),
				),
			),
			promqlbuilder.Scalar(
				promqlbuilder.Count(
					vector.New(
						vector.WithMetricName("instance_device:node_disk_io_time_seconds:rate5m"),
						vector.WithLabelMatchers(
							label.New("job").Equal("node"),
							label.New("instance").EqualRegexp("$instance"),
						),
					),
				),
			),
		),
		&parser.NumberLiteral{Val: 0},
	),
	"NodeExporterClusterNodeDiskSaturationPercentage": promqlbuilder.Neq(
		promqlbuilder.Div(
			vector.New(
				vector.WithMetricName("instance_device:node_disk_io_time_weighted_seconds:rate5m"),
				vector.WithLabelMatchers(
					label.New("job").Equal("node"),
				),
			),
			promqlbuilder.Scalar(
				promqlbuilder.Count(
					vector.New(
						vector.WithMetricName("instance_device:node_disk_io_time_weighted_seconds:rate5m"),
						vector.WithLabelMatchers(
							label.New("job").Equal("node"),
							label.New("instance").EqualRegexp("$instance"),
						),
					),
				),
			),
		),
		&parser.NumberLiteral{Val: 0},
	),
	"NodeExporterClusterNodeDiskSpacePercentage": promqlbuilder.Div(
		promqlbuilder.Sum(
			promqlbuilder.Max(
				promqlbuilder.Neq(
					promqlbuilder.Sub(
						vector.New(
							vector.WithMetricName("node_filesystem_size_bytes"),
							vector.WithLabelMatchers(
								label.New("fstype").NotEqual(""),
								label.New("instance").EqualRegexp("$instance"),
								label.New("job").Equal("node"),
								label.New("mountpoint").NotEqual(""),
							),
						),
						vector.New(
							vector.WithMetricName("node_filesystem_avail_bytes"),
							vector.WithLabelMatchers(
								label.New("fstype").NotEqual(""),
								label.New("instance").EqualRegexp("$instance"),
								label.New("job").Equal("node"),
								label.New("mountpoint").NotEqual(""),
							),
						),
					),
					&parser.NumberLiteral{Val: 0},
				),
			).Without("fstype", "mountpoint"),
		).Without("device"),
		promqlbuilder.Scalar(
			promqlbuilder.Sum(
				promqlbuilder.Max(
					vector.New(
						vector.WithMetricName("node_filesystem_size_bytes"),
						vector.WithLabelMatchers(
							label.New("fstype").NotEqual(""),
							label.New("instance").EqualRegexp("$instance"),
							label.New("job").Equal("node"),
							label.New("mountpoint").NotEqual(""),
						),
					),
				).Without("fstype", "mountpoint"),
			),
		),
	),
	"NodeExporterClusterNodeNetworkSaturationBytes": promqlbuilder.Neq(
		vector.New(
			vector.WithMetricName("instance:node_network_receive_drop_excluding_lo:rate5m"),
			vector.WithLabelMatchers(
				label.New("job").Equal("node"),
			),
		),
		&parser.NumberLiteral{Val: 0},
	),
	"NodeExporterClusterNodeNetworkUsageBytesRecv": promqlbuilder.Neq(
		vector.New(
			vector.WithMetricName("instance:node_network_receive_bytes_excluding_lo:rate5m"),
			vector.WithLabelMatchers(
				label.New("job").Equal("node"),
			),
		),
		&parser.NumberLiteral{Val: 0},
	),
	"NodeExporterClusterNodeNetworkUsageBytesTrams": promqlbuilder.Neq(
		vector.New(
			vector.WithMetricName("instance:node_network_transmit_bytes_excluding_lo:rate5m"),
			vector.WithLabelMatchers(
				label.New("job").Equal("node"),
			),
		),
		&parser.NumberLiteral{Val: 0},
	),
	"NodeExporterNodeAverageLoad1": vector.New(
		vector.WithMetricName("node_load1"),
		vector.WithLabelMatchers(
			label.New("job").Equal("node"),
			label.New("instance").Equal("$instance"),
		),
	),
	"NodeExporterNodeAverageLoad5": vector.New(
		vector.WithMetricName("node_load5"),
		vector.WithLabelMatchers(
			label.New("job").Equal("node"),
			label.New("instance").Equal("$instance"),
		),
	),
	"NodeExporterNodeAverageLoad15": vector.New(
		vector.WithMetricName("node_load15"),
		vector.WithLabelMatchers(
			label.New("job").Equal("node"),
			label.New("instance").Equal("$instance"),
		),
	),
	"NodeExporterNodeAverageCountCPU": promqlbuilder.Count(
		vector.New(
			vector.WithMetricName("node_cpu_seconds_total"),
			vector.WithLabelMatchers(
				label.New("job").Equal("node"),
				label.New("instance").Equal("$instance"),
				label.New("mode").Equal("idle"),
			),
		),
	),
	"NodeExporterNodeMemoryUsageBytesBuffers": vector.New(
		vector.WithMetricName("node_memory_Buffers_bytes"),
		vector.WithLabelMatchers(
			label.New("job").Equal("node"),
			label.New("instance").Equal("$instance"),
		),
	),
	"NodeExporterNodeMemoryUsageBytesCached": vector.New(
		vector.WithMetricName("node_memory_Cached_bytes"),
		vector.WithLabelMatchers(
			label.New("job").Equal("node"),
			label.New("instance").Equal("$instance"),
		),
	),
	"NodeExporterNodeMemoryUsageBytesMemFree": vector.New(
		vector.WithMetricName("node_memory_MemFree_bytes"),
		vector.WithLabelMatchers(
			label.New("job").Equal("node"),
			label.New("instance").Equal("$instance"),
		),
	),
	"NodeExporterNodeMemoryUsagePercentage": promqlbuilder.Sub(
		&parser.NumberLiteral{Val: 100},
		promqlbuilder.Div(
			promqlbuilder.Avg(
				vector.New(
					vector.WithMetricName("node_memory_MemAvailable_bytes"),
					vector.WithLabelMatchers(
						label.New("job").Equal("node"),
						label.New("instance").Equal("$instance"),
					),
				),
			),
			promqlbuilder.Mul(
				promqlbuilder.Avg(
					vector.New(
						vector.WithMetricName("node_memory_MemTotal_bytes"),
						vector.WithLabelMatchers(
							label.New("job").Equal("node"),
							label.New("instance").Equal("$instance"),
						),
					),
				),
				&parser.NumberLiteral{Val: 100},
			),
		),
	),
	"NodeExporterNodeDiskIOBytesTotal": promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName("node_disk_read_bytes_total"),
				vector.WithLabelMatchers(
					label.New("job").Equal("node"),
					label.New("instance").Equal("$instance"),
					label.New("device").NotEqual(""),
				),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	),
	"NodeExporterNodeDiskIOBytesTime": promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName("node_disk_io_time_seconds_total"),
				vector.WithLabelMatchers(
					label.New("job").Equal("node"),
					label.New("instance").Equal("$instance"),
					label.New("device").NotEqual(""),
				),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	),
	"NodeExporterNodeDiskIOSeconds": promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName("node_disk_io_time_seconds_total"),
				vector.WithLabelMatchers(
					label.New("job").Equal("node"),
					label.New("instance").Equal("$instance"),
					label.New("device").NotEqual(""),
				),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	),
	"NodeExporterNodeNetworkReceivedBytes": promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName("node_network_receive_bytes_total"),
				vector.WithLabelMatchers(
					label.New("job").Equal("node"),
					label.New("instance").Equal("$instance"),
					label.New("device").NotEqual("lo"),
				),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	),
	"NodeExporterNodeNetworkTransmitedBytes": promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName("node_network_transmit_bytes_total"),
				vector.WithLabelMatchers(
					label.New("job").Equal("node"),
					label.New("instance").Equal("$instance"),
					label.New("device").NotEqual("lo"),
				),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	),
}

// OverrideNodeExporterPanelQueries overrides the NodeExporterCommonPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
// The convention of naming followed, is to use Panel function name (with _suffix, in case panel has multiple queries)
func OverrideNodeExporterPanelQueries(queries map[string]parser.Expr) {
	maps.Copy(NodeExporterCommonPanelQueries, queries)
}
