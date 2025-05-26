package tempo

import (
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"

	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"

	commonSdk "github.com/perses/perses/go-sdk/common"
	barChartPanel "github.com/perses/plugins/barchart/sdk/go"
	stat "github.com/perses/plugins/statchart/sdk/go"
	tablePanel "github.com/perses/plugins/table/sdk/go"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
)

func RolloutProgress(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rollout progress",
		panel.Description("Shows progress of Tempo deployment"),
		barChartPanel.Chart(
			barChartPanel.WithMode(barChartPanel.PercentageMode),
			barChartPanel.Calculation(commonSdk.LastCalculation),
			barChartPanel.Format(commonSdk.Format{
				Unit: string(commonSdk.PercentUnit)}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(\n  sum by(tempo_service) (\n    label_replace(\n      kube_statefulset_status_replicas_updated{cluster=~\"$cluster\", namespace=~\"$namespace\",statefulset=~\".*(cortex-gw|distributor|ingester|query-frontend|querier|compactor|metrics-generator).*\"},\n      \"tempo_service\", \"$1\", \"statefulset\", \"(.*?)(?:-zone-[a-z])?\"\n    )\n  )\n  /\n  sum by(tempo_service) (\n    label_replace(\n      kube_statefulset_replicas{cluster=~\"$cluster\", namespace=~\"$namespace\"},\n      \"tempo_service\", \"$1\", \"statefulset\", \"(.*?)(?:-zone-[a-z])?\"\n    )\n  )\n) and (\n  sum by(tempo_service) (\n    label_replace(\n      kube_statefulset_replicas{cluster=~\"$cluster\", namespace=~\"$namespace\"},\n      \"tempo_service\", \"$1\", \"statefulset\", \"(.*?)(?:-zone-[a-z])?\"\n    )\n  )\n  > 0\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tempo_service}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(\n  sum by(tempo_service) (\n    label_replace(\n      kube_deployment_status_replicas_updated{cluster=~\"$cluster\", namespace=~\"$namespace\",deployment=~\".*(cortex-gw|distributor|ingester|query-frontend|querier|compactor|metrics-generator).*\"},\n      \"tempo_service\", \"$1\", \"deployment\", \"(.*?)(?:-zone-[a-z])?\"\n    )\n  )\n  /\n  sum by(tempo_service) (\n    label_replace(\n      kube_deployment_spec_replicas{cluster=~\"$cluster\", namespace=~\"$namespace\"},\n      \"tempo_service\", \"$1\", \"deployment\", \"(.*?)(?:-zone-[a-z])?\"\n    )\n  )\n) and (\n  sum by(tempo_service) (\n    label_replace(\n      kube_deployment_spec_replicas{cluster=~\"$cluster\", namespace=~\"$namespace\"},\n      \"tempo_service\", \"$1\", \"deployment\", \"(.*?)(?:-zone-[a-z])?\"\n    )\n  )\n  > 0\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tempo_service}}"),
			),
		),
	)
}

func TempoWrites2xx(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Writes - 2xx",
		panel.Description("Success ratio of trace requests handled by Tempo"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.PercentDecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.ValueFontSize(50),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\",status_code=~\"2.+\"}[$__rate_interval])) /\nsum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\"}[$__rate_interval]))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func TempoWrites4xx(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Writes - 4xx",
		panel.Description("Ratio of trace requests errors of the request handled by Tempo"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.PercentDecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.ValueFontSize(50),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\",status_code=~\"4.+\"}[$__rate_interval])) /\nsum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\"}[$__rate_interval]))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func TempoWrites5xx(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Writes - 5xx",
		panel.Description("Server error rate of the request handled by Tempo"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.PercentDecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.ValueFontSize(50),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\",status_code=~\"5.+\"}[$__rate_interval])) /\nsum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\"}[$__rate_interval]))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func TempoWritesLatency(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Writes 99th latency",
		panel.Description("Shows the 99th quantile latency for Tempo"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.SecondsUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.ValueFontSize(50),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
				Steps: []commonSdk.StepOption{
					{
						Color: "orange",
						Value: 0.2,
					},
					{
						Color: "red",
						Value: 0.5,
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum by (le) (tempo_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\"}))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func TempoReads2xx(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Reads - 2xx",
		panel.Description("Success ratio of trace requests handled by Tempo"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.PercentDecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.ValueFontSize(50),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"tempo_api_.*\",status_code=~\"2.+\"}[$__rate_interval])) /\nsum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"tempo_api_.*\"}[$__rate_interval]))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func TempoReads4xx(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Reads - 4xx",
		panel.Description("Ratio of trace requests errors of the request handled by Tempo"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.PercentDecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.ValueFontSize(50),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"tempo_api_.*\",status_code=~\"4.+\"}[$__rate_interval])) /\nsum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"tempo_api_.*\"}[$__rate_interval]))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func TempoReads5xx(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Reads - 5xx",
		panel.Description("Server error rate of the request handled by Tempo"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.PercentDecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.ValueFontSize(50),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"tempo_api_.*\",status_code=~\"5.+\"}[$__rate_interval])) /\nsum(rate(tempo_request_duration_seconds_count{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"tempo_api_.*\"}[$__rate_interval]))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func TempoReadsLatency(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Writes 99th latency",
		panel.Description("Shows the 99th quantile latency for Tempo"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.SecondsUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.ValueFontSize(50),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
				Steps: []commonSdk.StepOption{
					{
						Color: "orange",
						Value: 1,
					},
					{
						Color: "red",
						Value: 2.5,
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum by (le) (tempo_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"tempo_api_.*\"}))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func TempoUnhealthyPods(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Unhealthy pods",
		panel.Description("Shows the status of Tempo deployment"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          string(commonSdk.DecimalUnit),
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
					"kube_deployment_status_replicas_unavailable{cluster=~\"$cluster\", namespace=~\"$namespace\", deployment=~\".*(cortex-gw|distributor|ingester|query-frontend|querier|compactor|metrics-generator).*\"}\n> 0\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{deployment}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"kube_statefulset_status_replicas_current{cluster=~\"$cluster\", namespace=~\"$namespace\", statefulset=~\".*(cortex-gw|distributor|ingester|query-frontend|querier|compactor|metrics-generator).*\"} -\nkube_statefulset_status_replicas_ready {cluster=~\"$cluster\", namespace=~\"$namespace\", statefulset=~\".*(cortex-gw|distributor|ingester|query-frontend|querier|compactor|metrics-generator).*\"}\n> 0\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{statefulset}}"),
			),
		),
	)
}

func TempoPodsCount(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Pods count per version",
		panel.Description("Table with the count of tempo pods by version"),
		tablePanel.Table(
			tablePanel.Transform([]commonSdk.Transform{
				{
					Kind: commonSdk.MergeSeriesKind,
					Spec: commonSdk.MergeSeriesSpec{
						Disabled: false,
					},
				},
			}),
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "container",
					Header: "Container",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "version",
					Header: "Version",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value",
					Header: "number of pods",
					Align:  tablePanel.RightAlign,
				},
				{
					Name: "timestamp",
					Hide: true,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"count by(container, version) (\n  label_replace(\n    kube_pod_container_info{cluster=~\"$cluster\", namespace=~\"$namespace\",container=~\".*(cortex-gw|distributor|ingester|query-frontend|querier|compactor|metrics-generator).*\"},\n    \"version\", \"$1\", \"image\", \".*:(.+)-.*\"\n  )\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func TempoLatencyHistory(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Latency vs 24h ago",
		panel.Description("Shows the 99th quantile latency of Tempo Read and Writes."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.MilliSecondsUnit),
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
					"1 - (\n  avg_over_time(histogram_quantile(0.99, sum by (le) (tempo_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\"} offset 24h))[1h:])\n  /\n  avg_over_time(histogram_quantile(0.99, sum by (le) (tempo_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"(opentelemetry_proto_collector_trace_v1_traceservice_export|otlp_v1_traces)\"}))[1h:])\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("writes"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"1 - (\n  avg_over_time(histogram_quantile(0.99, sum by (le) (tempo_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"tempo_api_.*\"} offset 24h))[1h:])\n  /\n  avg_over_time(histogram_quantile(0.99, sum by (le) (tempo_request_duration_seconds_bucket{cluster=~\"$cluster\", job=~\"($namespace)/cortex-gw(-internal)?\", route=~\"tempo_api_.*\"}))[1h:])\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("reads"),
			),
		),
	)
}
