package thanos

import (
	"time"

	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/matrix"
	"github.com/perses/promql-builder/vector"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	rulehelpers "github.com/perses/community-dashboards/pkg/rules"
)

// Runbook fragments
const (
	runbookThanosCompactIsDown                                     = "#thanoscompactisdown"
	runbookThanosQueryIsDown                                       = "#thanosqueryisdown"
	runbookThanosReceiveIsDown                                     = "#thanosreceiveisdown"
	runbookThanosRuleIsDown                                        = "#thanosruleisdown"
	runbookThanosStoreIsDown                                       = "#thanosstoreisdown"
	runbookThanosCompactMultipleRunning                            = "#thanoscompactmultiplerunning"
	runbookThanosCompactHalted                                     = "#thanoscompacthalted"
	runbookThanosCompactHighCompactionFailures                     = "#thanoscompacthighcompactionfailures"
	runbookThanosCompactBucketHighOperationFailures                = "#thanoscompactbuckethighoperationfailures"
	runbookThanosCompactHasNotRun                                  = "#thanoscompacthasnotrun"
	runbookThanosQueryHttpRequestQueryErrorRateHigh                = "#thanosqueryhttprequestqueryerrorratehigh"
	runbookThanosQueryGrpcServerErrorRate                          = "#thanosquerygrpcservererrorrate"
	runbookThanosQueryGrpcClientErrorRate                          = "#thanosquerygrpcclienterrorrate"
	runbookThanosQueryHighDNSFailures                              = "#thanosqueryhighdnsfailures"
	runbookThanosQueryInstantLatencyHigh                           = "#thanosqueryinstantlatencyhigh"
	runbookThanosReceiveHttpRequestErrorRateHigh                   = "#thanosreceivehttprequesterrorratehigh"
	runbookThanosReceiveHttpRequestLatencyHigh                     = "#thanosreceivehttprequestlatencyhigh"
	runbookThanosReceiveHighReplicationFailures                    = "#thanosreceivehighreplicationfailures"
	runbookThanosReceiveHighForwardRequestFailures                 = "#thanosreceivehighforwardrequestfailures"
	runbookThanosReceiveHighHashringFileRefreshFailures            = "#thanosreceivehighhashringfilerefreshfailures"
	runbookThanosReceiveConfigReloadFailure                        = "#thanosreceiveconfigreloadfailure"
	runbookThanosReceiveNoUpload                                   = "#thanosreceivenoupload"
	runbookThanosReceiveLimitsConfigReloadFailure                  = "#thanosreceivelimitsconfigreloadfailure"
	runbookThanosReceiveLimitsHighMetaMonitoringQueriesFailureRate = "#thanosreceivelimitshighmetamonitoringqueriesfailurerate"
	runbookThanosReceiveTenantLimitedByHeadSeries                  = "#thanosreceivetenantlimitedbyheadseries"
	runbookThanosStoreGrpcErrorRate                                = "#thanosstoregrpcerrorrate"
	runbookThanosStoreBucketHighOperationFailures                  = "#thanosstorebuckethighoperationfailures"
	runbookThanosStoreObjstoreOperationLatencyHigh                 = "#thanosstoreobjstoreoperationlatencyhigh"
	runbookThanosRuleQueueIsDroppingAlerts                         = "#thanosrulequeueisdroppingalerts"
	runbookThanosRuleSenderIsFailingAlerts                         = "#thanosrulesenderisfailingalerts"
	runbookThanosRuleHighRuleEvaluationFailures                    = "#thanosrulehighruleevaluationfailures"
	runbookThanosRuleHighRuleEvaluationWarnings                    = "#thanosrulehighruleevaluationwarnings"
	runbookThanosRuleRuleEvaluationLatencyHigh                     = "#thanosruleruleevaluationlatencyhigh"
	runbookThanosRuleGrpcErrorRate                                 = "#thanosrulegrpcerrorrate"
	runbookThanosRuleConfigReloadFailure                           = "#thanosruleconfigreloadfailure"
	runbookThanosRuleQueryHighDNSFailures                          = "#thanosrulequeryhighdnsfailures"
	runbookThanosRuleAlertmanagerHighDNSFailures                   = "#thanosrulealertmanagerhighdnsfailures"
	runbookThanosRuleNoEvaluationFor10Intervals                    = "#thanosrulenoevaluationfor10intervals"
	runbookThanosNoRuleEvaluations                                 = "#thanosnoruleevaluations"
)

// BuildThanosRules builds the Thanos rules for the given namespace, dashboard URLs, runbook URL, labels, and annotations.
func BuildThanosRules(
	namespace,
	componentDashboardURL,
	compactDashboardURL,
	queryDashboardURL,
	receiveDashboardURL,
	storeDashboardURL,
	ruleDashboardURL,
	runbookURL string,
	labels map[string]string,
	annotations map[string]string,
) *monitoringv1.PrometheusRule {

	groups := []monitoringv1.RuleGroup{
		ThanosComponentAbsentGroup(componentDashboardURL, runbookURL),
		ThanosCompactGroup(compactDashboardURL, runbookURL),
		ThanosQueryGroup(queryDashboardURL, runbookURL),
		ThanosReceiveGroup(receiveDashboardURL, runbookURL),
		ThanosStoreGroup(storeDashboardURL, runbookURL),
		ThanosRuleGroup(ruleDashboardURL, runbookURL),
	}

	return rulehelpers.NewPrometheusRule(
		"thanos-rules",
		namespace,
		labels,
		annotations,
		groups,
	)
}

func ThanosComponentAbsentGroup(dashboardURL, runbookURL string) monitoringv1.RuleGroup {
	rules := []monitoringv1.Rule{
		rulehelpers.NewAlertingRule(
			"ThanosCompactIsDown",
			promqlbuilder.Absent(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("up"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp("thanos-compact.*"),
						),
					),
					promqlbuilder.NewNumber(1),
				),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosCompactIsDown,
				"ThanosCompact has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
				"ThanosCompact has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
				"Thanos component has disappeared from {{$labels.namespace}}.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosQueryIsDown",
			promqlbuilder.Absent(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("up"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp("thanos-query.*"),
						),
					),
					promqlbuilder.NewNumber(1),
				),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosQueryIsDown,
				"ThanosQuery has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
				"ThanosQuery has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
				"Thanos component has disappeared from {{$labels.namespace}}.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosReceiveRouterIsDown",
			promqlbuilder.Absent(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("up"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp("thanos-receive-router.*"),
						),
					),
					promqlbuilder.NewNumber(1),
				),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosReceiveIsDown,
				"ThanosReceiveRouter has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
				"ThanosReceiveRouter has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
				"Thanos component has disappeared from {{$labels.namespace}}.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosReceiveIngesterIsDown",
			promqlbuilder.Absent(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("up"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp("thanos-receive-ingester.*"),
						),
					),
					promqlbuilder.NewNumber(1),
				),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosReceiveIsDown,
				"ThanosReceiveIngester has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
				"ThanosReceiveIngester has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
				"Thanos component has disappeared from {{$labels.namespace}}.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosRuleIsDown",
			promqlbuilder.Absent(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("up"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp("thanos-ruler.*"),
						),
					),
					promqlbuilder.NewNumber(1),
				),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosRuleIsDown,
				"ThanosRule has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
				"ThanosRule has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
				"Thanos component has disappeared from {{$labels.namespace}}.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosStoreIsDown",
			promqlbuilder.Absent(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("up"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp("thanos-store.*"),
						),
					),
					promqlbuilder.NewNumber(1),
				),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosStoreIsDown,
				"ThanosStore has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
				"ThanosStore has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
				"Thanos component has disappeared from {{$labels.namespace}}.",
			),
		),
	}

	return rulehelpers.NewRuleGroup("thanos-component-absent", "", nil, rules)
}

func ThanosCompactGroup(dashboardURL, runbookURL string) monitoringv1.RuleGroup {
	rules := []monitoringv1.Rule{
		rulehelpers.NewAlertingRule(
			"ThanosCompactMultipleRunning",
			promqlbuilder.Sum(
				promqlbuilder.Gtr(
					vector.New(
						vector.WithMetricName("up"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp("thanos-compact.*"),
						),
					),
					promqlbuilder.NewNumber(1),
				),
			).By("namespace", "job"),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosCompactMultipleRunning,
				"No more than one Thanos Compact instance should be running at once. There are {{$value}} in {{$labels.namespace}} instances running.",
				"No more than one Thanos Compact instance should be running at once. There are {{$value}} in {{$labels.namespace}} instances running.",
				"Thanos Compact has multiple instances running.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosCompactHalted",
			promqlbuilder.Eqlc(
				vector.New(
					vector.WithMetricName("thanos_compact_halted"),
					vector.WithLabelMatchers(
						label.New("job").EqualRegexp("thanos-compact.*"),
					),
				),
				promqlbuilder.NewNumber(1),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosCompactHalted,
				"Thanos Compact {{$labels.job}} in {{$labels.namespace}} has failed to run and now is halted.",
				"Thanos Compact {{$labels.job}} in {{$labels.namespace}} has failed to run and now is halted.",
				"Thanos Compact has failed to run and is now halted.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosCompactHighCompactionFailures",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_compact_group_compactions_failures_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-compact.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_compact_group_compactions_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-compact.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(5),
			),
			"15m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosCompactHighCompactionFailures,
				"Thanos Compact {{$labels.job}} in {{$labels.namespace}} is failing to execute {{$value | humanize}}% of compactions.",
				"Thanos Compact {{$labels.job}} in {{$labels.namespace}} is failing to execute {{$value | humanize}}% of compactions.",
				"Thanos Compact is failing to execute compactions.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosCompactBucketHighOperationFailures",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_objstore_bucket_operation_failures_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-compact.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_objstore_bucket_operations_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-compact.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(5),
			),
			"15m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosCompactBucketHighOperationFailures,
				"Thanos Compact {{$labels.job}} in {{$labels.namespace}} Bucket is failing to execute {{$value | humanize}}% of operations.",
				"Thanos Compact {{$labels.job}} in {{$labels.namespace}} Bucket is failing to execute {{$value | humanize}}% of operations.",
				"Thanos Compact Bucket is having a high number of operation failures.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosCompactHasNotRun",
			promqlbuilder.Gtr(
				promqlbuilder.Div(
					promqlbuilder.Div(
						promqlbuilder.Sub(
							promqlbuilder.Time(),
							promqlbuilder.Max(
								promqlbuilder.MaxOverTime(
									matrix.New(
										vector.New(
											vector.WithMetricName("thanos_objstore_bucket_last_successful_upload_time"),
											vector.WithLabelMatchers(
												label.New("job").EqualRegexp("thanos-compact.*"),
											),
										),
										matrix.WithRange(24*time.Hour),
									),
								),
							).By("namespace", "job"),
						),
						promqlbuilder.NewNumber(60),
					),
					promqlbuilder.NewNumber(60),
				),
				promqlbuilder.NewNumber(24),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosCompactHasNotRun,
				"Thanos Compact {{$labels.job}} in {{$labels.namespace}} has not uploaded anything for 24 hours.",
				"Thanos Compact {{$labels.job}} in {{$labels.namespace}} has not uploaded anything for 24 hours.",
				"Thanos Compact has not uploaded anything for last 24 hours.",
			),
		),
	}

	return rulehelpers.NewRuleGroup("thanos-compact", "", nil, rules)
}

func ThanosQueryGroup(dashboardURL, runbookURL string) monitoringv1.RuleGroup {
	rules := []monitoringv1.Rule{
		rulehelpers.NewAlertingRule(
			"ThanosQueryHttpRequestQueryErrorRateHigh",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("http_requests_total"),
										vector.WithLabelMatchers(
											label.New("code").EqualRegexp("5.."),
											label.New("job").EqualRegexp("thanos-query.*"),
											label.New("handler").Equal("query"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("http_requests_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-query.*"),
											label.New("handler").Equal("query"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(5),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosQueryHttpRequestQueryErrorRateHigh,
				"Thanos Query {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of \"query\" requests.",
				"Thanos Query {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of \"query\" requests.",
				"Thanos Query is failing to handle requests.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosQueryGrpcServerErrorRate",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("grpc_server_handled_total"),
										vector.WithLabelMatchers(
											label.New("grpc_code").EqualRegexp("Unknown|ResourceExhausted|Internal|Unavailable|DataLoss|DeadlineExceeded"),
											label.New("job").EqualRegexp("thanos-query.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("grpc_server_started_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-query.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(5),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosQueryGrpcServerErrorRate,
				"Thanos Query {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
				"Thanos Query {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
				"Thanos Query is failing to handle requests.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosQueryGrpcClientErrorRate",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("grpc_client_handled_total"),
										vector.WithLabelMatchers(
											label.New("grpc_code").NotEqual("OK"),
											label.New("job").EqualRegexp("thanos-query.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("grpc_client_started_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-query.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(5),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosQueryGrpcClientErrorRate,
				"Thanos Query {{$labels.job}} in {{$labels.namespace}} is failing to send {{$value | humanize}}% of requests.",
				"Thanos Query {{$labels.job}} in {{$labels.namespace}} is failing to send {{$value | humanize}}% of requests.",
				"Thanos Query is failing to send requests.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosQueryHighDNSFailures",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_query_store_apis_dns_failures_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-query.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_query_store_apis_dns_lookups_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-query.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(1),
			),
			"15m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosQueryHighDNSFailures,
				"Thanos Query {{$labels.job}} in {{$labels.namespace}} have {{$value | humanize}}% of failing DNS queries for store endpoints.",
				"Thanos Query {{$labels.job}} in {{$labels.namespace}} have {{$value | humanize}}% of failing DNS queries for store endpoints.",
				"Thanos Query is having high number of DNS failures.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosQueryInstantLatencyHigh",
			promqlbuilder.And(
				promqlbuilder.Gtr(
					promqlbuilder.HistogramQuantile(
						0.99,
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("http_request_duration_seconds_bucket"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-query.*"),
											label.New("handler").Equal("query"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job", "le"),
					),
					promqlbuilder.NewNumber(90),
				),
				promqlbuilder.Gtr(
					promqlbuilder.Sum(
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("http_request_duration_seconds_count"),
									vector.WithLabelMatchers(
										label.New("job").EqualRegexp("thanos-query.*"),
										label.New("handler").Equal("query"),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
					).By("namespace", "job"),
					promqlbuilder.NewNumber(0),
				),
			),
			"10m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosQueryInstantLatencyHigh,
				"Thanos Query {{$labels.job}} in {{$labels.namespace}} has a 99th percentile latency of {{$value}} seconds for instant queries.",
				"Thanos Query {{$labels.job}} in {{$labels.namespace}} has a 99th percentile latency of {{$value}} seconds for instant queries.",
				"Thanos Query has high latency for queries.",
			),
		),
	}

	return rulehelpers.NewRuleGroup("thanos-query", "", nil, rules)
}

func ThanosReceiveGroup(dashboardURL, runbookURL string) monitoringv1.RuleGroup {
	rules := []monitoringv1.Rule{
		rulehelpers.NewAlertingRule(
			"ThanosReceiveHttpRequestErrorRateHigh",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("http_requests_total"),
										vector.WithLabelMatchers(
											label.New("code").EqualRegexp("5.."),
											label.New("job").EqualRegexp("thanos-receive-router.*"),
											label.New("handler").Equal("receive"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("http_requests_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-receive-router.*"),
											label.New("handler").Equal("receive"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(5),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosReceiveHttpRequestErrorRateHigh,
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
				"Thanos Receive is failing to handle requests.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosReceiveHttpRequestLatencyHigh",
			promqlbuilder.And(
				promqlbuilder.Gtr(
					promqlbuilder.HistogramQuantile(
						0.99,
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("http_request_duration_seconds_bucket"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-receive-router.*"),
											label.New("handler").Equal("receive"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job", "le"),
					),
					promqlbuilder.NewNumber(10),
				),
				promqlbuilder.Gtr(
					promqlbuilder.Sum(
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("http_request_duration_seconds_count"),
									vector.WithLabelMatchers(
										label.New("job").EqualRegexp("thanos-receive-router.*"),
										label.New("handler").Equal("receive"),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
					).By("namespace", "job"),
					promqlbuilder.NewNumber(0),
				),
			),
			"10m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosReceiveHttpRequestLatencyHigh,
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} has a 99th percentile latency of {{ $value }} seconds for requests.",
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} has a 99th percentile latency of {{ $value }} seconds for requests.",
				"Thanos Receive has high HTTP requests latency.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosReceiveHighReplicationFailures",
			promqlbuilder.And(
				promqlbuilder.Gtr(
					vector.New(
						vector.WithMetricName("thanos_receive_replication_factor"),
					),
					promqlbuilder.NewNumber(1),
				),
				promqlbuilder.Mul(
					promqlbuilder.Gtr(
						promqlbuilder.Div(
							promqlbuilder.Sum(
								promqlbuilder.Rate(
									matrix.New(
										vector.New(
											vector.WithMetricName("thanos_receive_replications_total"),
											vector.WithLabelMatchers(
												label.New("result").Equal("error"),
												label.New("job").EqualRegexp("thanos-receive-router.*"),
											),
										),
										matrix.WithRange(5*time.Minute),
									),
								),
							).By("namespace", "job"),
							promqlbuilder.Sum(
								promqlbuilder.Rate(
									matrix.New(
										vector.New(
											vector.WithMetricName("thanos_receive_replications_total"),
											vector.WithLabelMatchers(
												label.New("job").EqualRegexp("thanos-receive-router.*"),
											),
										),
										matrix.WithRange(5*time.Minute),
									),
								),
							).By("namespace", "job"),
						),
						promqlbuilder.Div(
							promqlbuilder.Max(
								promqlbuilder.Floor(
									promqlbuilder.Div(
										promqlbuilder.Add(
											vector.New(
												vector.WithMetricName("thanos_receive_replication_factor"),
												vector.WithLabelMatchers(
													label.New("job").EqualRegexp("thanos-receive-router.*"),
												),
											),
											promqlbuilder.NewNumber(1),
										),
										promqlbuilder.NewNumber(2),
									),
								),
							).By("namespace", "job"),
							promqlbuilder.Max(
								vector.New(
									vector.WithMetricName("thanos_receive_hashring_nodes"),
									vector.WithLabelMatchers(
										label.New("job").EqualRegexp("thanos-receive-router.*"),
									),
								),
							).By("namespace", "job"),
						),
					),
					promqlbuilder.NewNumber(100),
				),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosReceiveHighReplicationFailures,
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to replicate {{$value | humanize}}% of requests.",
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to replicate {{$value | humanize}}% of requests.",
				"Thanos Receive is having high number of replication failures.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosReceiveHighForwardRequestFailures",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_receive_forward_requests_total"),
										vector.WithLabelMatchers(
											label.New("result").Equal("error"),
											label.New("job").EqualRegexp("thanos-receive-router.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_receive_forward_requests_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-receive-router.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(20),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosReceiveHighForwardRequestFailures,
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to forward {{$value | humanize}}% of requests.",
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to forward {{$value | humanize}}% of requests.",
				"Thanos Receive is failing to forward requests.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosReceiveHighHashringFileRefreshFailures",
			promqlbuilder.Gtr(
				promqlbuilder.Div(
					promqlbuilder.Sum(
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("thanos_receive_hashrings_file_errors_total"),
									vector.WithLabelMatchers(
										label.New("job").EqualRegexp("thanos-receive-router.*"),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
					).By("namespace", "job"),
					promqlbuilder.Sum(
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("thanos_receive_hashrings_file_refreshes_total"),
									vector.WithLabelMatchers(
										label.New("job").EqualRegexp("thanos-receive-router.*"),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
					).By("namespace", "job"),
				),
				promqlbuilder.NewNumber(0),
			),
			"15m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosReceiveHighHashringFileRefreshFailures,
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to refresh hashring file, {{$value | humanize}} of attempts failed.",
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to refresh hashring file, {{$value | humanize}} of attempts failed.",
				"Thanos Receive is failing to refresh hasring file.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosReceiveConfigReloadFailure",
			promqlbuilder.Neq(
				promqlbuilder.Avg(
					vector.New(
						vector.WithMetricName("thanos_receive_config_last_reload_successful"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp("thanos-receive-router.*"),
						),
					),
				).By("namespace", "job"),
				promqlbuilder.NewNumber(1),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosReceiveConfigReloadFailure,
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} has not been able to reload hashring configurations.",
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} has not been able to reload hashring configurations.",
				"Thanos Receive is failing to reload hashring configurations.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosReceiveNoUpload",
			promqlbuilder.Add(
				promqlbuilder.Sub(
					vector.New(
						vector.WithMetricName("up"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp("thanos-receive-ingester.*"),
						),
					),
					promqlbuilder.NewNumber(1),
				),
				promqlbuilder.Eqlc(
					promqlbuilder.Sum(
						promqlbuilder.Increase(
							matrix.New(
								vector.New(
									vector.WithMetricName("thanos_shipper_uploads_total"),
									vector.WithLabelMatchers(
										label.New("job").EqualRegexp("thanos-receive-ingester.*"),
									),
								),
								matrix.WithRange(3*time.Hour),
							),
						),
					).By("namespace", "job", "instance"),
					promqlbuilder.NewNumber(0),
				),
			).On("namespace", "job", "instance"),
			"3h",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosReceiveNoUpload,
				"Thanos Receive {{$labels.instance}} in {{$labels.namespace}} has not uploaded latest data to object storage.",
				"Thanos Receive {{$labels.instance}} in {{$labels.namespace}} has not uploaded latest data to object storage.",
				"Thanos Receive has not uploaded latest data to object storage.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosReceiveLimitsConfigReloadFailure",
			promqlbuilder.Gtr(
				promqlbuilder.Sum(
					promqlbuilder.Increase(
						matrix.New(
							vector.New(
								vector.WithMetricName("thanos_receive_limits_config_reload_err_total"),
								vector.WithLabelMatchers(
									label.New("job").EqualRegexp("thanos-receive-router.*"),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
				).By("namespace", "job"),
				promqlbuilder.NewNumber(0),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosReceiveLimitsConfigReloadFailure,
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} has not been able to reload the limits configuration.",
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} has not been able to reload the limits configuration.",
				"Thanos Receive has not been able to reload the limits configuration.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosReceiveLimitsHighMetaMonitoringQueriesFailureRate",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Increase(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_receive_metamonitoring_failed_queries_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-receive-router.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
						promqlbuilder.NewNumber(20),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(20),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosReceiveLimitsHighMetaMonitoringQueriesFailureRate,
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing for {{$value | humanize}}% of meta monitoring queries.",
				"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing for {{$value | humanize}}% of meta monitoring queries.",
				"Thanos Receive has not been able to update the number of head series.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosReceiveTenantLimitedByHeadSeries",
			promqlbuilder.Gtr(
				promqlbuilder.Sum(
					promqlbuilder.Increase(
						matrix.New(
							vector.New(
								vector.WithMetricName("thanos_receive_head_series_limited_requests_total"),
								vector.WithLabelMatchers(
									label.New("job").EqualRegexp("thanos-receive-router.*"),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
				).By("namespace", "job", "tenant"),
				promqlbuilder.NewNumber(0),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosReceiveTenantLimitedByHeadSeries,
				"Thanos Receive tenant {{$labels.tenant}} in {{$labels.namespace}} is limited by head series.",
				"Thanos Receive tenant {{$labels.tenant}} in {{$labels.namespace}} is limited by head series.",
				"Thanos Receive tenant is limited by head series.",
			),
		),
	}

	return rulehelpers.NewRuleGroup("thanos-receive", "", nil, rules)
}

func ThanosStoreGroup(dashboardURL, runbookURL string) monitoringv1.RuleGroup {
	rules := []monitoringv1.Rule{
		rulehelpers.NewAlertingRule(
			"ThanosStoreGrpcErrorRate",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("grpc_server_handled_total"),
										vector.WithLabelMatchers(
											label.New("grpc_code").EqualRegexp("Unknown|ResourceExhausted|Internal|Unavailable|DataLoss|DeadlineExceeded"),
											label.New("job").EqualRegexp("thanos-store.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("grpc_server_started_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-store.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(5),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosStoreGrpcErrorRate,
				"Thanos Store {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
				"Thanos Store {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
				"Thanos Store is failing to handle gRPC requests.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosStoreBucketHighOperationFailures",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_objstore_bucket_operation_failures_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-store.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_objstore_bucket_operations_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-store.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(5),
			),
			"15m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosStoreBucketHighOperationFailures,
				"Thanos Store {{$labels.job}} in {{$labels.namespace}} Bucket is failing to execute {{$value | humanize}}% of operations.",
				"Thanos Store {{$labels.job}} in {{$labels.namespace}} Bucket is failing to execute {{$value | humanize}}% of operations.",
				"Thanos Store Bucket is failing to execute operations.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosStoreObjstoreOperationLatencyHigh",
			promqlbuilder.And(
				promqlbuilder.Gtr(
					promqlbuilder.HistogramQuantile(
						0.99,
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_objstore_bucket_operation_duration_seconds_bucket"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-store.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job", "le"),
					),
					promqlbuilder.NewNumber(7),
				),
				promqlbuilder.Gtr(
					promqlbuilder.Sum(
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("thanos_objstore_bucket_operation_duration_seconds_count"),
									vector.WithLabelMatchers(
										label.New("job").EqualRegexp("thanos-store.*"),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
					).By("namespace", "job"),
					promqlbuilder.NewNumber(0),
				),
			),
			"10m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosStoreObjstoreOperationLatencyHigh,
				"Thanos Store {{$labels.job}} in {{$labels.namespace}} Bucket has a 99th percentile latency of {{$value}} seconds for the bucket operations.",
				"Thanos Store {{$labels.job}} in {{$labels.namespace}} Bucket has a 99th percentile latency of {{$value}} seconds for the bucket operations.",
				"Thanos Store is having high latency for bucket operations.",
			),
		),
	}

	return rulehelpers.NewRuleGroup("thanos-store", "", nil, rules)
}

func ThanosRuleGroup(dashboardURL, runbookURL string) monitoringv1.RuleGroup {
	rules := []monitoringv1.Rule{
		rulehelpers.NewAlertingRule(
			"ThanosRuleQueueIsDroppingAlerts",
			promqlbuilder.Gtr(
				promqlbuilder.Sum(
					promqlbuilder.Rate(
						matrix.New(
							vector.New(
								vector.WithMetricName("thanos_alert_queue_alerts_dropped_total"),
								vector.WithLabelMatchers(
									label.New("job").EqualRegexp("thanos-ruler.*"),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
				).By("namespace", "job", "instance"),
				promqlbuilder.NewNumber(0),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosRuleQueueIsDroppingAlerts,
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} is failing to queue rulehelpers.",
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} is failing to queue rulehelpers.",
				"Thanos Rule is failing to queue rulehelpers.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosRuleSenderIsFailingAlerts",
			promqlbuilder.Gtr(
				promqlbuilder.Sum(
					promqlbuilder.Rate(
						matrix.New(
							vector.New(
								vector.WithMetricName("thanos_alert_sender_alerts_dropped_total"),
								vector.WithLabelMatchers(
									label.New("job").EqualRegexp("thanos-ruler.*"),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
				).By("namespace", "job", "instance"),
				promqlbuilder.NewNumber(0),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosRuleSenderIsFailingAlerts,
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} is failing to send alerts to alertmanager.",
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} is failing to send alerts to alertmanager.",
				"Thanos Rule is failing to send alerts to alertmanager.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosRuleHighRuleEvaluationFailures",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("prometheus_rule_evaluation_failures_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-ruler.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job", "instance"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("prometheus_rule_evaluations_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-ruler.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job", "instance"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(5),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosRuleHighRuleEvaluationFailures,
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} is failing to evaluate rules.",
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} is failing to evaluate rules.",
				"Thanos Rule is failing to evaluate rules.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosRuleHighRuleEvaluationWarnings",
			promqlbuilder.Gtr(
				promqlbuilder.Sum(
					promqlbuilder.Rate(
						matrix.New(
							vector.New(
								vector.WithMetricName("thanos_rule_evaluation_with_warnings_total"),
								vector.WithLabelMatchers(
									label.New("job").EqualRegexp("thanos-ruler.*"),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
				).By("namespace", "job", "instance"),
				promqlbuilder.NewNumber(0),
			),
			"15m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosRuleHighRuleEvaluationWarnings,
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} has high number of evaluation warnings.",
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} has high number of evaluation warnings.",
				"Thanos Rule has high number of evaluation warnings.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosRuleRuleEvaluationLatencyHigh",
			promqlbuilder.Gtr(
				promqlbuilder.Sum(
					vector.New(
						vector.WithMetricName("prometheus_rule_group_last_duration_seconds"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp("thanos-ruler.*"),
						),
					),
				).By("namespace", "job", "instance", "rule_group"),
				promqlbuilder.Sum(
					vector.New(
						vector.WithMetricName("prometheus_rule_group_interval_seconds"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp("thanos-ruler.*"),
						),
					),
				).By("namespace", "job", "instance", "rule_group"),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosRuleRuleEvaluationLatencyHigh,
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} has higher evaluation latency than interval for {{$labels.rule_group}}.",
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} has higher evaluation latency than interval for {{$labels.rule_group}}.",
				"Thanos Rule has high rule evaluation latency.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosRuleGrpcErrorRate",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("grpc_server_handled_total"),
										vector.WithLabelMatchers(
											label.New("grpc_code").EqualRegexp("Unknown|ResourceExhausted|Internal|Unavailable|DataLoss|DeadlineExceeded"),
											label.New("job").EqualRegexp("thanos-ruler.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job", "instance"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("grpc_server_started_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-ruler.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job", "instance"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(5),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosRuleGrpcErrorRate,
				"Thanos Rule {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
				"Thanos Rule {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
				"Thanos Rule is failing to handle grpc requests.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosRuleConfigReloadFailure",
			promqlbuilder.Neq(
				promqlbuilder.Avg(
					vector.New(
						vector.WithMetricName("thanos_rule_config_last_reload_successful"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp("thanos-ruler.*"),
						),
					),
				).By("namespace", "job", "instance"),
				promqlbuilder.NewNumber(1),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosRuleConfigReloadFailure,
				"Thanos Rule {{$labels.job}} in {{$labels.namespace}} has not been able to reload its configuration.",
				"Thanos Rule {{$labels.job}} in {{$labels.namespace}} has not been able to reload its configuration.",
				"Thanos Rule has not been able to reload configuration.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosRuleQueryHighDNSFailures",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_rule_query_apis_dns_failures_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-ruler.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job", "instance"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_rule_query_apis_dns_lookups_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-ruler.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job", "instance"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(1),
			),
			"15m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosRuleQueryHighDNSFailures,
				"Thanos Rule {{$labels.job}} in {{$labels.namespace}} has {{$value | humanize}}% of failing DNS queries for query endpoints.",
				"Thanos Rule {{$labels.job}} in {{$labels.namespace}} has {{$value | humanize}}% of failing DNS queries for query endpoints.",
				"Thanos Rule is having high number of DNS failures.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosRuleAlertmanagerHighDNSFailures",
			promqlbuilder.Gtr(
				promqlbuilder.Mul(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_rule_alertmanagers_dns_failures_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-ruler.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job", "instance"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_rule_alertmanagers_dns_lookups_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp("thanos-ruler.*"),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job", "instance"),
					),
					promqlbuilder.NewNumber(100),
				),
				promqlbuilder.NewNumber(1),
			),
			"15m",
			map[string]string{
				"service":  "thanos",
				"severity": "medium",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosRuleAlertmanagerHighDNSFailures,
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} has {{$value | humanize}}% of failing DNS queries for Alertmanager endpoints.",
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} has {{$value | humanize}}% of failing DNS queries for Alertmanager endpoints.",
				"Thanos Rule is having high number of DNS failures.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosRuleNoEvaluationFor10Intervals",
			promqlbuilder.Gtr(
				promqlbuilder.Sub(
					promqlbuilder.Time(),
					promqlbuilder.Max(
						vector.New(
							vector.WithMetricName("prometheus_rule_group_last_evaluation_timestamp_seconds"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp("thanos-ruler.*"),
							),
						),
					).By("namespace", "job", "instance", "group"),
				),
				promqlbuilder.Mul(
					promqlbuilder.NewNumber(10),
					promqlbuilder.Max(
						vector.New(
							vector.WithMetricName("prometheus_rule_group_interval_seconds"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp("thanos-ruler.*"),
							),
						),
					).By("namespace", "job", "instance", "group"),
				),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "high",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosRuleNoEvaluationFor10Intervals,
				"Thanos Rule {{$labels.job}} in {{$labels.namespace}} has rule groups that did not evaluate for at least 10x of their expected interval.",
				"Thanos Rule {{$labels.job}} in {{$labels.namespace}} has rule groups that did not evaluate for at least 10x of their expected interval.",
				"Thanos Rule has rule groups that did not evaluate for 10 intervals.",
			),
		),
		rulehelpers.NewAlertingRule(
			"ThanosNoRuleEvaluations",
			promqlbuilder.And(
				promqlbuilder.Lte(
					promqlbuilder.Sum(
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("prometheus_rule_evaluations_total"),
									vector.WithLabelMatchers(
										label.New("job").EqualRegexp("thanos-ruler.*"),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
					).By("namespace", "job", "instance"),
					promqlbuilder.NewNumber(0),
				),
				promqlbuilder.Gtr(
					promqlbuilder.Sum(
						vector.New(
							vector.WithMetricName("thanos_rule_loaded_rules"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp("thanos-ruler.*"),
							),
						),
					).By("namespace", "job", "instance"),
					promqlbuilder.NewNumber(0),
				),
			),
			"5m",
			map[string]string{
				"service":  "thanos",
				"severity": "critical",
			},
			rulehelpers.BuildAnnotations(
				dashboardURL,
				runbookURL,
				runbookThanosNoRuleEvaluations,
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} did not perform any rule evaluations in the past 10 minutes.",
				"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} did not perform any rule evaluations in the past 10 minutes.",
				"Thanos Rule did not perform any rule evaluations.",
			),
		),
	}

	return rulehelpers.NewRuleGroup("thanos-rule", "", nil, rules)
}
