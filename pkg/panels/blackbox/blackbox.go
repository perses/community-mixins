package blackbox

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	stat "github.com/perses/perses/go-sdk/panel/stat"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
	"github.com/perses/perses/go-sdk/prometheus/query"
)

func ProbeStatusMapfunc(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Status Map",
		panel.Description("Shows Probe success, either 1 if up, or 0 if down"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.DecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
				Steps: []commonSdk.StepOption{
					{
						Color: "red",
						Value: 0,
					},
					{
						Color: "green",
						Value: 1,
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max by (instance) (probe_success{job=~'$job'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

func ProbeSuccessCount(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Probes",
		panel.Description("Counts Probes Success"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.DecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"count(probe_success{job=~'$job'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ProbeSuccessPercent(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Probes Success",
		panel.Description("Percentage of Probes Success"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.PercentDecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "red",
				Steps: []commonSdk.StepOption{
					{
						Color: "yellow",
						Value: 0.99,
					},
					{
						Color: "green",
						Value: 0.999,
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(count(probe_success{job=~'$job'} == 1) OR vector(0)) / count(probe_success{job=~'$job'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ProbeHTTPSSL(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Probes SSL",
		panel.Description("Proportion of HTTP probes that successfully used SSL"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.PercentDecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "red",
				Steps: []commonSdk.StepOption{
					{
						Color: "yellow",
						Value: 0.99,
					},
					{
						Color: "green",
						Value: 0.999,
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"count(probe_http_ssl{job=~'$job'} == 1) / count(probe_http_version{job=~'$job'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ProbeAverageDuration(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Probe Average Duration",
		panel.Description("Duration in Seconds"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.SecondsUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"avg(probe_duration_seconds{job=~'$job'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ProbeUptimeSuccess(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Uptime",
		panel.Description("Max uptime by instance"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.PercentDecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "red",
				Steps: []commonSdk.StepOption{
					{
						Color: "yellow",
						Value: 0.99,
					},
					{
						Color: "green",
						Value: 0.999,
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max by (instance) (probe_success{job=~'$job',instance=~'$instance'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

func ProbeUptimeMonthly(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Uptime 30d",
		panel.Description("30 days uptime"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.PercentDecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "red",
				Steps: []commonSdk.StepOption{
					{
						Color: "yellow",
						Value: 0.99,
					},
					{
						Color: "green",
						Value: 0.999,
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"avg_over_time(probe_success{job=~'$job',instance=~'$instance'}[30d])",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

func ProbeDurationSeconds(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Probe Duration",
		panel.Description("Shows Probe duration in seconds"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.SecondsUnit),
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
				promql.SetLabelMatchers(
					"sum by (instance) (avg by (phase,instance) (probe_http_duration_seconds{job=~'$job',instance=~'$instance'}))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("HTTP duration"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"avg by (instance) (probe_duration_seconds{job=~'$job',instance=~'$instance'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Total probe duration"),
			),
		),
	)
}

func ProbePhases(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Probe Phases",
		panel.Description("Shows Probe duration in seconds"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.SecondsUnit),
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
				promql.SetLabelMatchers(
					"avg by (phase) (probe_http_duration_seconds{job=~'$job',instance=~'$instance'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{phase}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"avg by (phase) (probe_icmp_duration_seconds{job=~'$job',instance=~'$instance'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{phase}}"),
			),
		),
	)
}

func ProbeStatusCode(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Latest Response Code",
		panel.Description("Shows Probe Last Status Code"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.DecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
				Steps: []commonSdk.StepOption{
					{
						Color: "red",
						Value: 500,
					},
					{
						Color: "yellow",
						Value: 400,
					},
					{
						Color: "blue",
						Value: 300,
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max by (instance) (probe_http_status_code{job=~'$job',instance=~'$instance'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

func ProbeTLSVersion(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("SSL Version",
		panel.Description("Shows Probe TLS Version"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.DecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max by (instance, version) (probe_tls_version_info{job=~'$job',instance=~'$instance'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{version}}"),
			),
		),
	)
}

func ProbeSSLExpiry(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("SSL Certificate Expiry",
		panel.Description("Shows When SSL Cert Will Expire"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.DecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"min by (instance) (probe_ssl_earliest_cert_expiry{job=~'$job',instance=~'$instance'}) - time()",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ProbeRedirects(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Redirects",
		panel.Description("Shows Probes HTTP Redirects"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.DecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "blue",
				Steps: []commonSdk.StepOption{
					{
						Color: "green",
						Value: 0,
						Name:  "No",
					},
					{
						Color: "blue",
						Value: 1,
						Name:  "Yes",
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max by (instance) (probe_http_redirects{job=~'$job',instance=~'$instance'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ProbeHTTPVersion(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("HTTP Version",
		panel.Description("Shows Probes HTTP Version"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.DecimalUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "blue",
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max by (instance) (probe_http_version{job=~'$job',instance=~'$instance'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{version}}"),
			),
		),
	)
}

func ProbeAverageDurationInstance(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Average Latency",
		panel.Description("Average Duration in Seconds by Instance"),
		stat.Chart(
			stat.Calculation(commonSdk.MeanCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.SecondsUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"avg by (instance)(probe_duration_seconds{job=~'$job',instance=~'$instance'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ProbeAverageDNSLookupPerInstance(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Average DNS Lookup Time",
		panel.Description("Average DNS lookup Time per instance"),
		stat.Chart(
			stat.Calculation(commonSdk.MeanCalculation),
			stat.Format(commonSdk.Format{
				Unit: string(commonSdk.SecondsUnit),
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"avg by (instance)(probe_dns_lookup_time_seconds{job=~'$job',instance=~'$instance'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}
