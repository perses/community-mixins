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
