package tempo

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	markdown "github.com/perses/plugins/markdown/sdk/go"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
)

func WritesGatewayQPS(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("QPS",
		panel.Description("Rate of HTTP request durations for Tempo Gateway"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (status) (\n  label_replace(label_replace(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\"}[$__rate_interval]),\n  \"status\", \"${1}xx\", \"status_code\", \"([0-9])..\"),\n  \"status\", \"${1}\", \"status_code\", \"([a-zA-Z]+)\"))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{status}}"),
			),
		),
	)
}

func WritesGatewayLatency(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Latency",
		panel.Description("Shows the 99th and 50th quantile latency of Gateway."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.MilliSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(rate(tempo_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 99th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum(rate(tempo_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 50th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_request_duration_seconds_sum{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\"}[$__rate_interval])) by () * 1e3 / sum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\"}[$__rate_interval])) by ()",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} Average"),
			),
		),
	)
}

func WritesEnvoyProxyQPS(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("QPS",
		panel.Description("Rate of gRPC response statuses from Envoy"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (grpc_status) (\n    rate(\n        label_replace(\n            {cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", __name__=~\"envoy_cluster_grpc_proto_collector_trace_v1_TraceService_[0-9]+\"},\n            \"grpc_status\", \"$1\", \"__name__\", \"envoy_cluster_grpc_proto_collector_trace_v1_TraceService_(.+)\"\n        )\n        [$__rate_interval:30s]\n    )\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{grpc_status}}"),
			),
		),
	)
}

func WritesEnvoygRPCStatusCodes(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("gRPC status codes",
		markdown.Markdown("gRPC status codes",
			markdown.Text(`Visit [Status codes and their use in gRPC](https://github.com/grpc/grpc/blob/master/doc/statuscodes.md)

Code | Number | Description
---|---|---
OK | 0 | Not an error; returned on success.
CANCELLED | 1 | The operation was cancelled, typically by the caller.
UNKNOWN | 2 | Unknown error. For example, this error may be returned when a Status value received from another address space belongs to an error space that is not known in this address space. Also errors raised by APIs that do not return enough error information may be converted to this error.
INVALID_ARGUMENT | 3 | The client specified an invalid argument. Note that this differs from FAILED_PRECONDITION.
DEADLINE_EXCEEDED | 4 | The deadline expired before the operation could complete.
NOT_FOUND | 5 | Some requested entity (e.g., file or directory) was not found.
ALREADY_EXISTS | 6 | The entity that a client attempted to create (e.g., file or directory) already exists.
PERMISSION_DENIED | 7 | The caller does not have permission to execute the specified operation.
RESOURCE_EXHAUSTED | 8 | Some resource has been exhausted, perhaps a per-user quota, or perhaps the entire file system is out of space.
FAILED_PRECONDITION | 9 | The operation was rejected because the system is not in a state required for the operation's execution.
ABORTED | 10 | The operation was aborted, typically due to a concurrency issue such as a sequencer check failure or transaction abort.
OUT_OF_RANGE | 11 | The operation was attempted past the valid range.
UNIMPLEMENTED | 12 | The operation is not implemented or is not supported/enabled in this service.
INTERNAL | 13 | Internal errors. This means that some invariants expected by the underlying system have been broken.
UNAVAILABLE | 14 | The service is currently unavailable. This is most likely a transient condition.
DATA_LOSS | 15 | Unrecoverable data loss or corruption.
UNAUTHENTICATED | 16 | The request does not have valid authentication credentials for the operation.`)),
	)
}

func WritesDistributorSpansSecond(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Spans / sec",
		panel.Description("Rate of Spans per Second for Tempo Distributor"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_receiver_accepted_spans{cluster=~\"$cluster\", job=~\"($namespace)/distributor\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("accepted"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_receiver_refused_spans{cluster=~\"$cluster\", job=~\"($namespace)/distributor\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("refused"),
			),
		),
	)
}

func WritesDistributorBytesPerSecond(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Bytes / sec",
		panel.Description("Rate of bytes received by Tempo distributors"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_distributor_bytes_received_total{cluster=~\"$cluster\", job=~\"($namespace)/distributor\"}[$__rate_interval])) by (status)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("received"),
			),
		),
	)
}

func WritesDistributorLatency(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Latency",
		panel.Description("Shows the 99th and 50th quantile latency of Distributor."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.MilliSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(rate(tempo_distributor_push_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/distributor\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 99th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.5, sum(rate(tempo_distributor_push_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/distributor\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 50th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_distributor_push_duration_seconds_sum{cluster=~\"$cluster\", job=~\"($namespace)/distributor\"}[$__rate_interval])) by () * 1e3 / sum(rate(tempo_distributor_push_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/distributor\"}[$__rate_interval])) by ()",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} Average"),
			),
		),
	)
}

func WritesDistributorKafkaAppendRecords(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Kafka append records / sec",
		panel.Description("Rate of bytes received by Tempo Distributors"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_distributor_kafka_appends_total{cluster=~\"$cluster\", job=~\"($namespace)/distributor\", status=\"success\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("appends"),
			),
		),
	)
}

func WritesDistributorKafkaAppendFail(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Kafka failed append records / sec",
		panel.Description("Rate of failed bytes received by Tempo distributors"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_distributor_kafka_appends_total{cluster=~\"$cluster\", job=~\"($namespace)/distributor\", status=\"fail\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("failed"),
			),
		),
	)
}

func WritesDistributorKafkaWrite(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Kafka write bytes / sec",
		panel.Description("Rate of append (write) operations the Tempo Distributor to Kafkas"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_distributor_kafka_write_bytes_total{cluster=~\"$cluster\", job=~\"($namespace)/distributor\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("writes"),
			),
		),
	)
}

func WritesDistributorKafkaWriteLatency(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Kafka write latency (sec)",
		panel.Description("Shows the 99th and 50th quantile latency of Distributor Kafka Write."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.MilliSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum by (le) (rate(tempo_distributor_kafka_write_latency_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/distributor\"}[$__rate_interval])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("50th percentile"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum by (le) (rate(tempo_distributor_kafka_write_latency_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/distributor\"}[$__rate_interval])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("99th percentile"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_distributor_kafka_write_latency_seconds_sum{cluster=~\"$cluster\", job=~\"($namespace)/distributor\"}[$__rate_interval])) / sum(rate(tempo_distributor_kafka_write_latency_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/distributor\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Average"),
			),
		),
	)
}

func WritesIngesterQPS(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("QPS",
		panel.Description("Rate of HTTP request durations for Tempo Ingester"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (status) (\n  label_replace(label_replace(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/ingester\", route=~\"/tempopb.Pusher/Push.*\"}[$__rate_interval]),\n  \"status\", \"${1}xx\", \"status_code\", \"([0-9])..\"),\n  \"status\", \"${1}\", \"status_code\", \"([a-zA-Z]+)\"))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{status}}"),
			),
		),
	)
}

func WritesIngesterLatency(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Latency",
		panel.Description("Shows the 99th and 50th quantile latency of Ingester."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.MilliSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(rate(tempo_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",route=~\"/tempopb.Pusher/Push.*\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 99th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum(rate(tempo_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",route=~\"/tempopb.Pusher/Push.*\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 50th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_request_duration_seconds_sum{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",route=~\"/tempopb.Pusher/Push.*\"}[$__rate_interval])) by () * 1e3 / sum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",route=~\"/tempopb.Pusher/Push.*\"}[$__rate_interval])) by ()",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} Average"),
			),
		),
	)
}

func WritesMemcachedIngesterQPS(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("QPS",
		panel.Description("Rate of HTTP request durations for Tempo Memcached Ingester"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (status) (\n  label_replace(label_replace(rate(tempo_memcache_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",method=\"Memcache.Put\"}[$__rate_interval]),\n  \"status\", \"${1}xx\", \"status_code\", \"([0-9])..\"),\n  \"status\", \"${1}\", \"status_code\", \"([a-zA-Z]+)\"))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{status}}"),
			),
		),
	)
}

func WritesMemcachedIngesterLatency(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Latency",
		panel.Description("Shows the 99th and 50th quantile latency of Memcached Ingester."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.MilliSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(rate(tempo_memcache_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",method=\"Memcache.Put\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 99th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum(rate(tempo_memcache_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",method=\"Memcache.Put\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 50th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_memcache_request_duration_seconds_sum{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",method=\"Memcache.Put\"}[$__rate_interval])) by () * 1e3 / sum(rate(tempo_memcache_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",method=\"Memcache.Put\"}[$__rate_interval])) by ()",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} Average"),
			),
		),
	)
}

func WritesBackendIngesterQPS(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("QPS",
		panel.Description("Rate of HTTP request durations for Tempo Backend Ingester"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (status) (\n  label_replace(label_replace(rate(tempodb_backend_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",operation=~\"(PUT|POST)\"}[$__rate_interval]),\n  \"status\", \"${1}xx\", \"status_code\", \"([0-9])..\"),\n  \"status\", \"${1}\", \"status_code\", \"([a-zA-Z]+)\"))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{status}}"),
			),
		),
	)
}

func WritesBackendIngesterLatency(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Latency",
		panel.Description("Shows the 99th and 50th quantile latency of Backend Ingester."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.MilliSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(rate(tempodb_backend_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",operation=~\"(PUT|POST)\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 99th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(rate(tempodb_backend_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",operation=~\"(PUT|POST)\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 50th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempodb_backend_request_duration_seconds_sum{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",operation=~\"(PUT|POST)\"}[$__rate_interval])) by () * 1e3 / sum(rate(tempodb_backend_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",operation=~\"(PUT|POST)\"}[$__rate_interval])) by ()",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} Average"),
			),
		),
	)
}

func WritesMemcachedCompactorQPS(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("QPS",
		panel.Description("Rate of HTTP request durations for Memcached Compactor"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (status) (\n  label_replace(label_replace(rate(tempo_memcache_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",method=\"Memcache.Put\"}[$__rate_interval]),\n  \"status\", \"${1}xx\", \"status_code\", \"([0-9])..\"),\n  \"status\", \"${1}\", \"status_code\", \"([a-zA-Z]+)\"))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{status}}"),
			),
		),
	)
}

func WritesMemcachedCompactorLatency(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Latency",
		panel.Description("Shows the 99th and 50th quantile latency of Backend Ingester."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.MilliSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(rate(tempo_memcache_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",method=\"Memcache.Put\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 99th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum(rate(tempo_memcache_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",method=\"Memcache.Put\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 50th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_memcache_request_duration_seconds_sum{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",method=\"Memcache.Put\"}[$__rate_interval])) by () * 1e3 / sum(rate(tempo_memcache_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",method=\"Memcache.Put\"}[$__rate_interval])) by ()",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} Average"),
			),
		),
	)
}

func WritesBackendCompactorQPS(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("QPS",
		panel.Description("Rate of HTTP request durations for Backend Compactor"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (status) (\n  label_replace(label_replace(rate(tempodb_backend_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",operation=~\"(PUT|POST)\"}[$__rate_interval]),\n  \"status\", \"${1}xx\", \"status_code\", \"([0-9])..\"),\n  \"status\", \"${1}\", \"status_code\", \"([a-zA-Z]+)\"))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{status}}"),
			),
		),
	)
}

func WritesBackendCompactorLatency(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Latency",
		panel.Description("Shows the 99th and 50th quantile latency of Backend Ingester."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.MilliSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(rate(tempodb_backend_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",operation=~\"(PUT|POST)\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 99th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(rate(tempodb_backend_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",operation=~\"(PUT|POST)\"}[$__rate_interval])) by (le,)) * 1e3",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} 50th"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempodb_backend_request_duration_seconds_sum{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",operation=~\"(PUT|POST)\"}[$__rate_interval])) by () * 1e3 / sum(rate(tempodb_backend_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",operation=~\"(PUT|POST)\"}[$__rate_interval])) by ()",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{route}} Average"),
			),
		),
	)
}
