package thanos

import (
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/perses/go-sdk/prometheus/query"

	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"

	commonSdk "github.com/perses/perses/go-sdk/common"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
)

func RuleEvaluationRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rule Group Evaluation Rate",
		panel.Description("Shows the rate of rule group evaluations."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.CountsPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, rule_group, strategy) (rate(prometheus_rule_evaluations_total{namespace='$namespace', job=~'$job'}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{rule_group}} {{strategy}} - {{job}} {{namespace}}"),
			),
		),
	)
}

func RuleEvaluationFailureRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rule Group Evaluation Failures",
		panel.Description("Shows the rate of failed rule group evaluations."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.CountsPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, rule_group, strategy) (rate(prometheus_rule_evaluation_failures_total{namespace='$namespace', job=~'$job'}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{rule_group}} {{strategy}} - {{job}} {{namespace}}"),
			),
		),
	)
}

func RuleGroupEvaluationsMissRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rule Group Evaluations Missed",
		panel.Description("Shows rate of rule group evaluations missed due to slow rule evaluations."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.CountsPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, rule_group, strategy) (rate(prometheus_rule_group_iterations_missed_total{namespace='$namespace', job=~'$job'}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{rule_group}} {{strategy}} - {{job}} {{namespace}}"),
			),
		),
	)
}

func RuleGroupEvaluationsTooSlow(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rule Group Evaluations Too Slow",
		panel.Description("Shows rule groups with evaluations taking longer than their set interval."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.SecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					`(
  sum by(namespace, job, rule_group) (prometheus_rule_group_last_duration_seconds{namespace='$namespace', job=~'$job'})
  >
  sum by(namespace, job, rule_group) (prometheus_rule_group_interval_seconds{namespace='$namespace', job=~'$job'})
)`,
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{rule_group}} - {{job}} {{namespace}}"),
			),
		),
	)
}

func AlertsDroppedRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Alerts Dropped",
		panel.Description("Shows rate of alerts dropped by Ruler when sending."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.CountsPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, alertmanager) (rate(thanos_alert_sender_alerts_dropped_total{namespace='$namespace', job=~'$job'}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{alertmanager}} - {{job}} {{namespace}}"),
			),
		),
	)
}

func AlertsSentRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Alerts Sent",
		panel.Description("Shows rate of alerts sent by ruler."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.CountsPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, alertmanager) (rate(thanos_alert_sender_alerts_sent_total{namespace='$namespace', job=~'$job'}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{alertmanager}} - {{job}} {{namespace}}"),
			),
		),
	)
}

func AlertSendingErrors(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Alert Sending Errors",
		panel.Description("Shows percentage of alert sending operations that have resulted in errors."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					`sum by (namespace, job) (rate(thanos_alert_sender_errors_total{namespace='$namespace', job=~'$job'}[5m])) / sum by (namespace, job) (rate(thanos_alert_sender_alerts_sent_total{namespace='$namespace', job=~'$job'}[5m]))`,
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} {{namespace}}"),
			),
		),
	)
}

func AlertSendingDurations(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Alert Sending Durations",
		panel.Description("Shows p50, p90 and p99 durations for alerts being sent from the ruler."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.MilliSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum by (namespace, job, le) (rate(thanos_alert_sender_latency_seconds_bucket{namespace='$namespace', job=~'$job'}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p50 {{job}} {{namespace}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum by (namespace, job, le) (rate(thanos_alert_sender_latency_seconds_bucket{namespace='$namespace', job=~'$job'}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p90 {{job}} {{namespace}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum by (namespace, job, le) (rate(thanos_alert_sender_latency_seconds_bucket{namespace='$namespace', job=~'$job'}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p99 {{job}} {{namespace}}"),
			),
		),
	)
}

func AlertQueuePushedRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Alert Queue Pushed",
		panel.Description("Shows the rate of alerts being pushed to sender queue."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.CountsPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job) (rate(thanos_alert_queue_alerts_pushed_total{namespace='$namespace', job=~'$job'}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} {{namespace}}"),
			),
		),
	)
}

func AlertQueuePoppedRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Alert Queue Popped",
		panel.Description("Shows rate of alerts popped from queue, to be sent to Alertmanagers."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.CountsPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job) (rate(thanos_alert_queue_alerts_popped_total{namespace='$namespace', job=~'$job'}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} {{namespace}}"),
			),
		),
	)
}

func DroppedRatio(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Dropped Percentage",
		panel.Description("Shows percentage of alerts dropped as compared to alerts being queued for sending."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					`sum by (namespace, job) (rate(thanos_alert_queue_alerts_dropped_total{namespace='$namespace', job=~'$job'}[5m])) / sum by (namespace, job) (rate(thanos_alert_queue_alerts_pushed_total{namespace='$namespace', job=~'$job'}[5m]))`,
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} {{namespace}}"),
			),
		),
	)
}
