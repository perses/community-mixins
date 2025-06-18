package opentelemetry

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
	"github.com/prometheus/prometheus/model/labels"
)

// SpanRate creates a panel that displays the rate of spans being processed by the OpenTelemetry collector.
// It shows both accepted spans (successfully pushed into the pipeline) and refused spans (failed to push).
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
//
// Returns a panelgroup.Option that can be used to add this panel to a dashboard.
func SpanRate(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Span Rate",
		panel.Description(`Accepted: rate of spans successfully pushed into the pipeline.
		Refused: rate of spans that could not be pushed into the pipeline.`),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["SpanRate_accepted"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Accepted: {{receiver}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["SpanRate_refused"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Refused: {{receiver}}"),
			),
		),
	)
}

// MeticPointsRate creates a panel that displays the rate of metric points being processed by the OpenTelemetry collector.
// It shows both accepted metric points (successfully pushed into the pipeline) and refused metric points (failed to push).
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
//
// Returns a panelgroup.Option that can be used to add this panel to a dashboard.
func MeticPointsRate(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Metic Points Rate",
		panel.Description(`Accepted: rate of metric points successfully pushed into the pipeline.
		Refused: rate of metric points that could not be pushed into the pipeline.`),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["MeticPointsRate_accepted"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Accepted: {{receiver}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["MeticPointsRate_refused"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Refused: {{receiver}}"),
			),
		),
	)
}

// LogRecordsRate creates a panel that displays the rate of log records being processed by the OpenTelemetry collector.
// It shows both accepted log records (successfully pushed into the pipeline) and refused log records (failed to push).
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
//
// Returns a panelgroup.Option that can be used to add this panel to a dashboard.
func LogRecordsRate(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Log Records Rate",
		panel.Description(`Accepted: rate of log records successfully pushed into the pipeline.
		Refused: rate of log records that could not be pushed into the pipeline.`),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["LogRecordsRate_accepted"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Accepted: {{receiver}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["LogRecordsRate_refused"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Refused: {{receiver}}"),
			),
		),
	)
}

// SpanProcessorRate creates a panel that displays the rate of spans being processed by individual processors
// in the OpenTelemetry collector pipeline. It shows both incoming and outgoing spans for each processor.
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
//
// Returns a panelgroup.Option that can be used to add this panel to a dashboard.
func SpanProcessorRate(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Span Rate",
		panel.Description("Rate of span processors"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["SpanProcessorRate_incoming_items"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Incoming: {{processor}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["SpanProcessorRate_outgoing_items"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Outgoing: {{processor}}"),
			),
		),
	)
}

// MetricProcessorRate creates a panel that displays the rate of metrics being processed by individual processors
// in the OpenTelemetry collector pipeline. It shows both incoming and outgoing metrics for each processor.
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
//
// Returns a panelgroup.Option that can be used to add this panel to a dashboard.
func MetricProcessorRate(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Metric Rate",
		panel.Description("Rate of metric processors"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["MetricProcessorRate_incoming_items"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Incoming: {{processor}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["MetricProcessorRate_outgoing_items"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Outgoing: {{processor}}"),
			),
		),
	)
}

// LogProcessorRate creates a panel that displays the rate of logs being processed by individual processors
// in the OpenTelemetry collector pipeline. It shows both incoming and outgoing logs for each processor.
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
//
// Returns a panelgroup.Option that can be used to add this panel to a dashboard.
func LogProcessorRate(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Log Rate",
		panel.Description("Rate of log processors"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["LogProcessorRate_incoming_items"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Incoming: {{processor}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["LogProcessorRate_outgoing_items"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Outgoing: {{processor}}"),
			),
		),
	)
}

// BatchProcessorBatchSendSize creates a panel that displays the size of batches being sent by the OpenTelemetry collector's batch processor.
// This metric helps monitor the efficiency of batch processing by showing the number of items (spans, metrics, or logs) in each batch.
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
//
// Returns a panelgroup.Option that can be used to add this panel to a dashboard.
func BatchProcessorBatchSendSize(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Batch Send Size",
		panel.Description("Number of items in each batch being processed"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["BatchProcessorRate_batch_send_size"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("{{le}}"),
			),
		),
	)
}

// BatchProcessorBatchSendSizeCount creates a panel that displays the number of batches being processed by the OpenTelemetry collector's batch processor.
// This metric helps monitor the efficiency of batch processing by showing the number of batches being processed.
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
func BatchProcessorBatchSendSizeCount(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Batch Send Size Count",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["BatchProcessorRate_batch_send_size_count"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Batch send size count: {{processor}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["BatchProcessorRate_batch_send_size_sum"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Batch send size sum: {{processor}}"),
			),
		),
	)
}

// BatchProcessorBatchSizeTriggerSend creates a panel that displays the number of batches being processed by the OpenTelemetry collector's batch processor due to a size trigger.
// This metric helps monitor the efficiency of batch processing by showing the number of batches being processed due to a size trigger.
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
func BatchProcessorBatchSizeTriggerSend(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Batch Size Trigger Send",
		panel.Description("Number of batches being processed by the OpenTelemetry collector's batch processor"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.DecimalUnit),
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["BatchProcessorRate_batch_size_trigger_send"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Batch sent due to a timeout trigger: {{processor}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["BatchProcessorRate_batch_timeout_trigger_send"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Batch sent due to a timeout trigger: {{processor}}"),
			),
		),
	)
}
