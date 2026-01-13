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

package opentelemetry

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/promql"
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
					Unit: &dashboards.DecimalUnit,
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
					Unit: &dashboards.DecimalUnit,
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
					Unit: &dashboards.DecimalUnit,
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
					Unit: &dashboards.DecimalUnit,
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
					Unit: &dashboards.DecimalUnit,
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
					Unit: &dashboards.DecimalUnit,
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
					Unit: &dashboards.DecimalUnit,
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
					Unit: &dashboards.DecimalUnit,
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
					Unit: &dashboards.DecimalUnit,
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

// SpanExporterRate creates a panel that displays the rate of spans being sent by the OpenTelemetry collector's span exporter.
// It shows both sent spans and failed spans (enqueued or sent).
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
func SpanExporterRate(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Span Rate",
		panel.Description("Rate of spans being sent by the OpenTelemetry collector's span exporter"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.DecimalUnit,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["ExporterRate_sent_spans"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Sent: {{exporter}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["ExporterRate_enqueue_failed_spans"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Enqueue failed: {{exporter}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["ExporterRate_send_failed_spans"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Send failed: {{exporter}}"),
			),
		),
	)
}

// MetricExporterRate creates a panel that displays the rate of metrics being sent by the OpenTelemetry collector's metric exporter.
// It shows both sent metrics and failed metrics (enqueued or sent).
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
func MetricExporterRate(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Metric Rate",
		panel.Description("Rate of metrics being sent by the OpenTelemetry collector's metric exporter"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.DecimalUnit,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["ExporterRate_sent_metrics"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Sent: {{exporter}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["ExporterRate_enqueue_failed_metrics"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Enqueue failed: {{exporter}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["ExporterRate_send_failed_metrics"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Send failed: {{exporter}}"),
			),
		),
	)
}

// LogExporterRate creates a panel that displays the rate of logs being sent by the OpenTelemetry collector's log exporter.
// It shows both sent logs and failed logs (enqueued or sent).
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
func LogExporterRate(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Log Rate",
		panel.Description("Rate of logs being sent by the OpenTelemetry collector's log exporter"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.DecimalUnit,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["ExporterRate_sent_logs"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Sent: {{exporter}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["ExporterRate_enqueue_failed_logs"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Enqueue failed: {{exporter}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["ExporterRate_send_failed_logs"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Send failed: {{exporter}}"),
			),
		),
	)
}

// QueueSizeExporterRate creates a panel that displays the size of the queue for each exporter.
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
func QueueSizeExporterRate(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Exporter Queue Size",
		panel.Description("Current size of the retry queue (in batches)"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.DecimalUnit,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["ExporterRate_queue_size"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Queue size: {{exporter}}"),
			),
		),
	)
}

// QueueCapacityExporterRate creates a panel that displays the capacity of the queue for each exporter.
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
func QueueCapacityExporterRate(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Exporter Queue Capacity",
		panel.Description("Fixed capacity of the retry queue (in batches)"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.DecimalUnit,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["ExporterRate_queue_capacity"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Queue capacity: {{exporter}}"),
			),
		),
	)
}

// QueueUtilizationExporterRate creates a panel that displays the utilization of the queue for each exporter.
//
// Parameters:
//   - datasource: The name of the Prometheus datasource to query
//   - labelMatchers: Optional label matchers to filter the metrics
func QueueUtilizationExporterRate(datasource string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Exporter Queue Utilization",
		panel.Description("Utilization of the retry queue (in batches)"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.PercentUnit,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					OpentelemetryCommonPanelQueries["ExporterRate_queue_utilization"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasource),
				query.SeriesNameFormat("Queue utilization: {{exporter}}"),
			),
		),
	)
}
