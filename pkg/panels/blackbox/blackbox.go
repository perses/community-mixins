package blackbox

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	stat "github.com/perses/plugins/statchart/sdk/go"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
)

// ProbeStatusMap creates a panel option for displaying Blackbox Probe Success for all instances
//
// The panel uses the following Prometheus metrics:
// - probe_success: indicates if the probe succeeded
//
// The panel shows:
// - Probe success, either 1 if up, or 0 if down
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.

func ProbeStatusMap(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Status Map",
		panel.Description("Shows Probe success, either 1 if up, or 0 if down"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit:          string(commonSdk.DecimalUnit),
				DecimalPlaces: 0,
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

// ProbeSuccessCount creates a panel option for displaying a counter of all Probe Success
//
// The panel uses the following Prometheus metrics:
// - probe_success: indicates if the probe succeeded
//
// The panel shows:
// - Counter of all existing probe_success, so the number of endpoints being probed
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ProbeSuccessCount(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Probes",
		panel.Description("Counts Probes Success"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit:          string(commonSdk.DecimalUnit),
				DecimalPlaces: 0,
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.ValueFontSize(50),
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

// ProbeSuccessPercent creates a panel option for displaying a Percentage of success probes from all executed probes
//
// The panel uses the following Prometheus metrics:
// - probe_success: indicates if the probe succeeded
//
// The panel shows:
// - Calculates the percentage of count(probe_success == 1) / count(probe_success)
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
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
			stat.ValueFontSize(50),
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

// ProbeHTTPSSL creates a panel option for displaying a Proportion HTTP probes that successfully used SSL
//
// The panel uses the following Prometheus metrics:
// - probe_http_ssl: whether an SSL/TLS connection was successfully established
// - probe_http_version: reports the HTTP version returned by the probed endpoint
//
// The panel shows:
// - Calculates the percentage of count(probe_http_ssl) == 1 / count(probe_http_version)
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
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
			stat.ValueFontSize(50),
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

// ProbeAverageDuration creates a panel option for displaying the Average of probe duration in seconds
//
// The panel uses the following Prometheus metrics:
// - probe_duration_seconds: measures the total time it takes to execute the probe
//
// The panel shows:
// - Calculates the Average of time the probe took to execute
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
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
			stat.ValueFontSize(50),
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

// ProbeUptimeSuccess creates The Probe Uptime
//
// The panel uses the following Prometheus metrics:
// - probe_success: indicates if the probe succeeded
//
// The panel shows:
// - Calculates the Max value for the probe_success by instance
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ProbeUptimeSuccess(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Uptime",
		panel.Description("Max uptime by instance"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit:          string(commonSdk.PercentDecimalUnit),
				DecimalPlaces: 0,
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.ValueFontSize(50),
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

// ProbeUptimeMonthly creates The Probe Uptime for 30 days
//
// The panel uses the following Prometheus metrics:
// - probe_success: indicates if the probe succeeded
//
// The panel shows:
// - Calculates the Max value for the probe_success by instance for 30 days
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
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
			stat.ValueFontSize(50),
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

// ProbeDurationSeconds creates Probes duration states for http duration and total probes duration
//
// The panel uses the following Prometheus metrics:
// - probe_duration_seconds: measures the total time it takes to execute the probe
// - probe_http_duration_seconds: measures the duration of the HTTP request/response phase of a probe.
//
// The panel shows:
// - The sum of average HTTP durations per phase, grouped by each instance.
// - The average total probe duration per instance for all matching targets
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
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

// ProbePhases creates Probes duration states for http duration and total probes duration
//
// The panel uses the following Prometheus metrics:
// - probe_icmp_duration_seconds: measures the duration of the ICMP probe.
// - probe_http_duration_seconds: measures the duration of the HTTP request/response phase of a probe.
//
// The panel shows:
// - The average total probe duration per instance for all matching targets
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
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

// ProbeStatusCode creates a panel that will show the last status code returned by the probe target
//
// The panel uses the following Prometheus metrics:
// - probe_http_status_code: HTTP status code returned by the probed target
//
// The panel shows:
// - The latest value for status code
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ProbeStatusCode(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Latest Response Code",
		panel.Description("Shows Probe Last Status Code"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit:          string(commonSdk.DecimalUnit),
				DecimalPlaces: 0,
			}),
			stat.ValueFontSize(50),
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

// ProbeTLSVersion creates a panel that will show the TLS version
//
// The panel uses the following Prometheus metrics:
// - probe_tls_version_info: indicates which TLS version was negotiated with the target
//
// The panel shows:
// - For each instance and TLS version, it returns the most recent value.
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
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
			stat.ValueFontSize(50),
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

// ProbeSSLExpiry creates a panel that will show when SSL Cert Will Expire
//
// The panel uses the following Prometheus metrics:
// - probe_tls_version_info: shows the Unix timestamp (in seconds) of the earliest expiring certificate in the SSL/TLS
//
// The panel shows:
// - The time until the certificate expirates
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
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
			stat.ValueFontSize(50),
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

// ProbeRedirects creates a panel that will show the number of HTTP redirects followed by the Blackbox Exporter during probes
//
// The panel uses the following Prometheus metrics:
// - probe_http_redirects: count of HTTP redirect responses
//
// The panel shows:
// -  HTTP redirects followed by the Blackbox Exporter during probes, grouped by instance
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ProbeRedirects(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Redirects",
		panel.Description("Shows Probes HTTP Redirects"),
		stat.Chart(
			stat.Calculation(commonSdk.LastCalculation),
			stat.Format(commonSdk.Format{
				Unit:          string(commonSdk.DecimalUnit),
				DecimalPlaces: 0,
			}),
			stat.WithSparkline(stat.Sparkline{
				Width: 1,
			}),
			stat.ValueFontSize(50),
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

// ProbeHTTPVersion creates a panel that will show the HTTP version returned by the probed endpoint
//
// The panel uses the following Prometheus metrics:
// - probe_http_version: reports the HTTP version returned by the probed endpoint
//
// The panel shows:
// - The most recent HTTP version used in probes
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
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
			stat.ValueFontSize(50),
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

// ProbeAverageDurationInstance creates a panel that will show Average Probe Duration in seconds
//
// The panel uses the following Prometheus metrics:
// - probe_duration_seconds: measures the total time it takes to execute the probe
//
// The panel shows:
// - The Average Probe Duration in Seconds by Instance
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
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
			stat.ValueFontSize(50),
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

// ProbeAverageDNSLookupPerInstance creates a panel that will show the time (in seconds) spent on DNS resolution during a probe.
//
// The panel uses the following Prometheus metrics:
// - probe_dns_lookup_time_seconds: time (in seconds) spent on DNS resolution during a probe.
//
// The panel shows:
// - The average time in seconds spent on DNS resolution during a probe.
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
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
			stat.ValueFontSize(50),
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
