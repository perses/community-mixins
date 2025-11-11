package thanos

import (
	"time"

	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/matrix"
	"github.com/perses/promql-builder/vector"

	rulehelpers "github.com/perses/community-mixins/pkg/rules"
	"github.com/perses/community-mixins/pkg/rules/rule-sdk/alerting"
	"github.com/perses/community-mixins/pkg/rules/rule-sdk/common"
	"github.com/perses/community-mixins/pkg/rules/rule-sdk/promtheusrule"
	"github.com/perses/community-mixins/pkg/rules/rule-sdk/rulegroup"
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

type ThanosRulesConfig struct {
	RunbookURL          string
	CompactDashboardURL string
	QueryDashboardURL   string
	ReceiveDashboardURL string
	StoreDashboardURL   string
	RuleDashboardURL    string

	ServiceLabelValue              string
	ReceiveRouterServiceSelector   string
	ReceiveIngesterServiceSelector string
	RulerServiceSelector           string
	StoreServiceSelector           string
	CompactServiceSelector         string
	QueryServiceSelector           string

	AdditionalAlertLabels      map[string]string
	AdditionalAlertAnnotations map[string]string
}

type ThanosRulesConfigOption func(*ThanosRulesConfig)

func WithRunbookURL(runbookURL string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		thanosRulesConfig.RunbookURL = runbookURL
	}
}

func WithServiceLabelValue(serviceLabelValue string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		thanosRulesConfig.ServiceLabelValue = serviceLabelValue
	}
}

func WithReceiveRouterServiceSelector(receiveRouterServiceSelector string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		if receiveRouterServiceSelector == "" {
			receiveRouterServiceSelector = "thanos-receive-router.*"
		}
		thanosRulesConfig.ReceiveRouterServiceSelector = receiveRouterServiceSelector
	}
}

func WithReceiveIngesterServiceSelector(receiveIngesterServiceSelector string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		if receiveIngesterServiceSelector == "" {
			receiveIngesterServiceSelector = "thanos-receive-ingester.*"
		}
		thanosRulesConfig.ReceiveIngesterServiceSelector = receiveIngesterServiceSelector
	}
}

func WithRulerServiceSelector(rulerServiceSelector string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		if rulerServiceSelector == "" {
			rulerServiceSelector = "thanos-ruler.*"
		}
		thanosRulesConfig.RulerServiceSelector = rulerServiceSelector
	}
}

func WithStoreServiceSelector(storeServiceSelector string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		if storeServiceSelector == "" {
			storeServiceSelector = "thanos-store.*"
		}
		thanosRulesConfig.StoreServiceSelector = storeServiceSelector
	}
}

func WithCompactServiceSelector(compactServiceSelector string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		if compactServiceSelector == "" {
			compactServiceSelector = "thanos-compact.*"
		}
		thanosRulesConfig.CompactServiceSelector = compactServiceSelector
	}
}

func WithQueryServiceSelector(queryServiceSelector string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		if queryServiceSelector == "" {
			queryServiceSelector = "thanos-query.*"
		}
		thanosRulesConfig.QueryServiceSelector = queryServiceSelector
	}
}

func WithAdditionalAlertLabels(additionalAlertLabels map[string]string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		thanosRulesConfig.AdditionalAlertLabels = additionalAlertLabels
	}
}

func WithAdditionalAlertAnnotations(additionalAlertAnnotations map[string]string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		thanosRulesConfig.AdditionalAlertAnnotations = additionalAlertAnnotations
	}
}

func WithCompactDashboardURL(compactDashboardURL string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		thanosRulesConfig.CompactDashboardURL = compactDashboardURL
	}
}

func WithQueryDashboardURL(queryDashboardURL string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		thanosRulesConfig.QueryDashboardURL = queryDashboardURL
	}
}

func WithReceiveDashboardURL(receiveDashboardURL string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		thanosRulesConfig.ReceiveDashboardURL = receiveDashboardURL
	}
}

func WithStoreDashboardURL(storeDashboardURL string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		thanosRulesConfig.StoreDashboardURL = storeDashboardURL
	}
}

func WithRuleDashboardURL(ruleDashboardURL string) ThanosRulesConfigOption {
	return func(thanosRulesConfig *ThanosRulesConfig) {
		thanosRulesConfig.RuleDashboardURL = ruleDashboardURL
	}
}

// NewThanosRulesBuilder creates a new Thanos rules builder.
func NewThanosRulesBuilder(
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	options ...ThanosRulesConfigOption,
) (promtheusrule.Builder, error) {
	thanosRulesConfig := ThanosRulesConfig{
		ReceiveRouterServiceSelector:   "thanos-receive-router.*",
		ReceiveIngesterServiceSelector: "thanos-receive-ingester.*",
		RulerServiceSelector:           "thanos-ruler.*",
		StoreServiceSelector:           "thanos-store.*",
		CompactServiceSelector:         "thanos-compact.*",
		QueryServiceSelector:           "thanos-query.*",
	}
	for _, option := range options {
		option(&thanosRulesConfig)
	}

	promRule, err := promtheusrule.New(
		"thanos-rules",
		namespace,
		promtheusrule.Labels(labels),
		promtheusrule.Annotations(annotations),
		promtheusrule.AddRuleGroup(
			"thanos-component-absent",
			thanosRulesConfig.ThanosComponentAbsentGroup()...,
		),
		promtheusrule.AddRuleGroup(
			"thanos-compact",
			thanosRulesConfig.ThanosCompactGroup()...,
		),
		promtheusrule.AddRuleGroup(
			"thanos-query",
			thanosRulesConfig.ThanosQueryGroup()...,
		),
		promtheusrule.AddRuleGroup(
			"thanos-receive",
			thanosRulesConfig.ThanosReceiveGroup()...,
		),
		promtheusrule.AddRuleGroup(
			"thanos-store",
			thanosRulesConfig.ThanosStoreGroup()...,
		),
		promtheusrule.AddRuleGroup(
			"thanos-rule",
			thanosRulesConfig.ThanosRuleGroup()...,
		),
	)

	return promRule, err
}

// BuildThanosRules builds the Thanos rules for the given namespace, dashboard URLs, runbook URL, labels, and annotations.
func BuildThanosRules(
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	options ...ThanosRulesConfigOption,
) rulehelpers.RuleResult {
	promRule, err := NewThanosRulesBuilder(namespace, labels, annotations, options...)
	if err != nil {
		return rulehelpers.NewRuleResult(nil, err).Component("thanos")
	}

	return rulehelpers.NewRuleResult(
		&promRule.PrometheusRule,
		nil,
	).Component("thanos")
}

func (t ThanosRulesConfig) ThanosComponentAbsentGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.AddRule(
			"ThanosCompactIsDown",
			alerting.Expr(
				promqlbuilder.Absent(
					promqlbuilder.Eqlc(
						vector.New(
							vector.WithMetricName("up"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp(t.CompactServiceSelector),
							),
						),
						promqlbuilder.NewNumber(1),
					),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.CompactDashboardURL,
						t.RunbookURL,
						runbookThanosCompactIsDown,
						"ThanosCompact has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
						"ThanosCompact has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
						"Thanos component has disappeared from {{$labels.namespace}}.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosQueryIsDown",
			alerting.Expr(
				promqlbuilder.Absent(
					promqlbuilder.Eqlc(
						vector.New(
							vector.WithMetricName("up"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp(t.QueryServiceSelector),
							),
						),
						promqlbuilder.NewNumber(1),
					),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.QueryDashboardURL,
						t.RunbookURL,
						runbookThanosQueryIsDown,
						"ThanosQuery has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
						"ThanosQuery has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
						"Thanos component has disappeared from {{$labels.namespace}}.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveRouterIsDown",
			alerting.Expr(
				promqlbuilder.Absent(
					promqlbuilder.Eqlc(
						vector.New(
							vector.WithMetricName("up"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
							),
						),
						promqlbuilder.NewNumber(1),
					),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.ReceiveDashboardURL,
						t.RunbookURL,
						runbookThanosReceiveIsDown,
						"ThanosReceiveRouter has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
						"ThanosReceiveRouter has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
						"Thanos component has disappeared from {{$labels.namespace}}.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveIngesterIsDown",
			alerting.Expr(
				promqlbuilder.Absent(
					promqlbuilder.Eqlc(
						vector.New(
							vector.WithMetricName("up"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp(t.ReceiveIngesterServiceSelector),
							),
						),
						promqlbuilder.NewNumber(1),
					),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.ReceiveDashboardURL,
						t.RunbookURL,
						runbookThanosReceiveIsDown,
						"ThanosReceiveIngester has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
						"ThanosReceiveIngester has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
						"Thanos component has disappeared from {{$labels.namespace}}.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRuleIsDown",
			alerting.Expr(
				promqlbuilder.Absent(
					promqlbuilder.Eqlc(
						vector.New(
							vector.WithMetricName("up"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp(t.RulerServiceSelector),
							),
						),
						promqlbuilder.NewNumber(1),
					),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.RuleDashboardURL,
						t.RunbookURL,
						runbookThanosRuleIsDown,
						"ThanosRule has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
						"ThanosRule has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
						"Thanos component has disappeared from {{$labels.namespace}}.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosStoreIsDown",
			alerting.Expr(
				promqlbuilder.Absent(
					promqlbuilder.Eqlc(
						vector.New(
							vector.WithMetricName("up"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp(t.StoreServiceSelector),
							),
						),
						promqlbuilder.NewNumber(1),
					),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.StoreDashboardURL,
						t.RunbookURL,
						runbookThanosStoreIsDown,
						"ThanosStore has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
						"ThanosStore has disappeared from {{$labels.namespace}}. Prometheus target for the component cannot be discovered.",
						"Thanos component has disappeared from {{$labels.namespace}}.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}

func (t ThanosRulesConfig) ThanosCompactGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.AddRule(
			"ThanosCompactMultipleRunning",
			alerting.Expr(
				promqlbuilder.Sum(
					promqlbuilder.Gtr(
						vector.New(
							vector.WithMetricName("up"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp(t.CompactServiceSelector),
							),
						),
						promqlbuilder.NewNumber(1),
					),
				).By("namespace", "job"),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.CompactDashboardURL,
						t.RunbookURL,
						runbookThanosCompactMultipleRunning,
						"No more than one Thanos Compact instance should be running at once. There are {{$value}} in {{$labels.namespace}} instances running.",
						"No more than one Thanos Compact instance should be running at once. There are {{$value}} in {{$labels.namespace}} instances running.",
						"Thanos Compact has multiple instances running.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosCompactHalted",
			alerting.Expr(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("thanos_compact_halted"),
						vector.WithLabelMatchers(
							label.New("job").EqualRegexp(t.CompactServiceSelector),
						),
					),
					promqlbuilder.NewNumber(1),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.CompactDashboardURL,
						t.RunbookURL,
						runbookThanosCompactHalted,
						"Thanos Compact {{$labels.job}} in {{$labels.namespace}} has failed to run and now is halted.",
						"Thanos Compact {{$labels.job}} in {{$labels.namespace}} has failed to run and now is halted.",
						"Thanos Compact has failed to run and is now halted.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosCompactHighCompactionFailures",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Mul(
						promqlbuilder.Div(
							promqlbuilder.Sum(
								promqlbuilder.Rate(
									matrix.New(
										vector.New(
											vector.WithMetricName("thanos_compact_group_compactions_failures_total"),
											vector.WithLabelMatchers(
												label.New("job").EqualRegexp(t.CompactServiceSelector),
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
												label.New("job").EqualRegexp(t.CompactServiceSelector),
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
			),
			alerting.For("15m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.CompactDashboardURL,
						t.RunbookURL,
						runbookThanosCompactHighCompactionFailures,
						"Thanos Compact {{$labels.job}} in {{$labels.namespace}} is failing to execute {{$value | humanize}}% of compactions.",
						"Thanos Compact {{$labels.job}} in {{$labels.namespace}} is failing to execute {{$value | humanize}}% of compactions.",
						"Thanos Compact is failing to execute compactions.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosCompactBucketHighOperationFailures",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Mul(
						promqlbuilder.Div(
							promqlbuilder.Sum(
								promqlbuilder.Rate(
									matrix.New(
										vector.New(
											vector.WithMetricName("thanos_objstore_bucket_operation_failures_total"),
											vector.WithLabelMatchers(
												label.New("job").EqualRegexp(t.CompactServiceSelector),
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
												label.New("job").EqualRegexp(t.CompactServiceSelector),
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
			),
			alerting.For("15m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.CompactDashboardURL,
						t.RunbookURL,
						runbookThanosCompactBucketHighOperationFailures,
						"Thanos Compact {{$labels.job}} in {{$labels.namespace}} Bucket is failing to execute {{$value | humanize}}% of operations.",
						"Thanos Compact {{$labels.job}} in {{$labels.namespace}} Bucket is failing to execute {{$value | humanize}}% of operations.",
						"Thanos Compact Bucket is having a high number of operation failures.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosCompactHasNotRun",
			alerting.Expr(
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
													label.New("job").EqualRegexp(t.CompactServiceSelector),
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
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.CompactDashboardURL,
						t.RunbookURL,
						runbookThanosCompactHasNotRun,
						"Thanos Compact {{$labels.job}} in {{$labels.namespace}} has not uploaded anything for 24 hours.",
						"Thanos Compact {{$labels.job}} in {{$labels.namespace}} has not uploaded anything for 24 hours.",
						"Thanos Compact has not uploaded anything for last 24 hours.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}

func (t ThanosRulesConfig) ThanosQueryGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.AddRule(
			"ThanosQueryHttpRequestQueryErrorRateHigh",
			alerting.Expr(
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
												label.New("job").EqualRegexp(t.QueryServiceSelector),
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
												label.New("job").EqualRegexp(t.QueryServiceSelector),
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
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.QueryDashboardURL,
						t.RunbookURL,
						runbookThanosQueryHttpRequestQueryErrorRateHigh,
						"Thanos Query {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of \"query\" requests.",
						"Thanos Query {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of \"query\" requests.",
						"Thanos Query is failing to handle requests.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosQueryGrpcServerErrorRate",
			alerting.Expr(
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
												label.New("job").EqualRegexp(t.QueryServiceSelector),
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
												label.New("job").EqualRegexp(t.QueryServiceSelector),
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
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.QueryDashboardURL,
						t.RunbookURL,
						runbookThanosQueryGrpcServerErrorRate,
						"Thanos Query {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
						"Thanos Query {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
						"Thanos Query is failing to handle requests.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosQueryGrpcClientErrorRate",
			alerting.Expr(
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
												label.New("job").EqualRegexp(t.QueryServiceSelector),
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
												label.New("job").EqualRegexp(t.QueryServiceSelector),
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
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.QueryDashboardURL,
						t.RunbookURL,
						runbookThanosQueryGrpcClientErrorRate,
						"Thanos Query {{$labels.job}} in {{$labels.namespace}} is failing to send {{$value | humanize}}% of requests.",
						"Thanos Query {{$labels.job}} in {{$labels.namespace}} is failing to send {{$value | humanize}}% of requests.",
						"Thanos Query is failing to send requests.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosQueryHighDNSFailures",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Mul(
						promqlbuilder.Div(
							promqlbuilder.Sum(
								promqlbuilder.Rate(
									matrix.New(
										vector.New(
											vector.WithMetricName("thanos_query_store_apis_dns_failures_total"),
											vector.WithLabelMatchers(
												label.New("job").EqualRegexp(t.QueryServiceSelector),
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
												label.New("job").EqualRegexp(t.QueryServiceSelector),
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
			),
			alerting.For("15m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.QueryDashboardURL,
						t.RunbookURL,
						runbookThanosQueryHighDNSFailures,
						"Thanos Query {{$labels.job}} in {{$labels.namespace}} have {{$value | humanize}}% of failing DNS queries for store endpoints.",
						"Thanos Query {{$labels.job}} in {{$labels.namespace}} have {{$value | humanize}}% of failing DNS queries for store endpoints.",
						"Thanos Query is having high number of DNS failures.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosQueryInstantLatencyHigh",
			alerting.Expr(
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
												label.New("job").EqualRegexp(t.QueryServiceSelector),
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
											label.New("job").EqualRegexp(t.QueryServiceSelector),
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
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.QueryDashboardURL,
						t.RunbookURL,
						runbookThanosQueryInstantLatencyHigh,
						"Thanos Query {{$labels.job}} in {{$labels.namespace}} has a 99th percentile latency of {{$value}} seconds for instant queries.",
						"Thanos Query {{$labels.job}} in {{$labels.namespace}} has a 99th percentile latency of {{$value}} seconds for instant queries.",
						"Thanos Query has high latency for queries.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}

func (t ThanosRulesConfig) ThanosReceiveGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.AddRule(
			"ThanosReceiveHttpRequestErrorRateHigh",
			alerting.Expr(
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
												label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
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
												label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
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
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.ReceiveDashboardURL,
						t.RunbookURL,
						runbookThanosReceiveHttpRequestErrorRateHigh,
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
						"Thanos Receive is failing to handle requests.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveHttpRequestLatencyHigh",
			alerting.Expr(
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
												label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
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
											label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
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
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.ReceiveDashboardURL,
						t.RunbookURL,
						runbookThanosReceiveHttpRequestLatencyHigh,
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} has a 99th percentile latency of {{ $value }} seconds for requests.",
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} has a 99th percentile latency of {{ $value }} seconds for requests.",
						"Thanos Receive has high HTTP requests latency.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveHighReplicationFailures",
			alerting.Expr(
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
													label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
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
													label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
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
														label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
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
											label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
										),
									),
								).By("namespace", "job"),
							),
						),
						promqlbuilder.NewNumber(100),
					),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.ReceiveDashboardURL,
						t.RunbookURL,
						runbookThanosReceiveHighReplicationFailures,
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to replicate {{$value | humanize}}% of requests.",
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to replicate {{$value | humanize}}% of requests.",
						"Thanos Receive is having high number of replication failures.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveHighForwardRequestFailures",
			alerting.Expr(
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
												label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
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
												label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
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
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.ReceiveDashboardURL,
						t.RunbookURL,
						runbookThanosReceiveHighForwardRequestFailures,
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to forward {{$value | humanize}}% of requests.",
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to forward {{$value | humanize}}% of requests.",
						"Thanos Receive is failing to forward requests.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveHighHashringFileRefreshFailures",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_receive_hashrings_file_errors_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
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
											label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("15m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.ReceiveDashboardURL,
						t.RunbookURL,
						runbookThanosReceiveHighHashringFileRefreshFailures,
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to refresh hashring file, {{$value | humanize}} of attempts failed.",
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing to refresh hashring file, {{$value | humanize}} of attempts failed.",
						"Thanos Receive is failing to refresh hasring file.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveConfigReloadFailure",
			alerting.Expr(
				promqlbuilder.Neq(
					promqlbuilder.Avg(
						vector.New(
							vector.WithMetricName("thanos_receive_config_last_reload_successful"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
							),
						),
					).By("namespace", "job"),
					promqlbuilder.NewNumber(1),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.ReceiveDashboardURL,
						t.RunbookURL,
						runbookThanosReceiveConfigReloadFailure,
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} has not been able to reload hashring configurations.",
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} has not been able to reload hashring configurations.",
						"Thanos Receive is failing to reload hashring configurations.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveNoUpload",
			alerting.Expr(
				promqlbuilder.Add(
					promqlbuilder.Sub(
						vector.New(
							vector.WithMetricName("up"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp(t.ReceiveIngesterServiceSelector),
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
											label.New("job").EqualRegexp(t.ReceiveIngesterServiceSelector),
										),
									),
									matrix.WithRange(3*time.Hour),
								),
							),
						).By("namespace", "job", "instance"),
						promqlbuilder.NewNumber(0),
					),
				).On("namespace", "job", "instance"),
			),
			alerting.For("3h"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.ReceiveDashboardURL,
						t.RunbookURL,
						runbookThanosReceiveNoUpload,
						"Thanos Receive {{$labels.instance}} in {{$labels.namespace}} has not uploaded latest data to object storage.",
						"Thanos Receive {{$labels.instance}} in {{$labels.namespace}} has not uploaded latest data to object storage.",
						"Thanos Receive has not uploaded latest data to object storage.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveLimitsConfigReloadFailure",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Sum(
						promqlbuilder.Increase(
							matrix.New(
								vector.New(
									vector.WithMetricName("thanos_receive_limits_config_reload_err_total"),
									vector.WithLabelMatchers(
										label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
					).By("namespace", "job"),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.ReceiveDashboardURL,
						t.RunbookURL,
						runbookThanosReceiveLimitsConfigReloadFailure,
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} has not been able to reload the limits configuration.",
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} has not been able to reload the limits configuration.",
						"Thanos Receive has not been able to reload the limits configuration.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveLimitsHighMetaMonitoringQueriesFailureRate",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Mul(
						promqlbuilder.Div(
							promqlbuilder.Sum(
								promqlbuilder.Increase(
									matrix.New(
										vector.New(
											vector.WithMetricName("thanos_receive_metamonitoring_failed_queries_total"),
											vector.WithLabelMatchers(
												label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
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
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.ReceiveDashboardURL,
						t.RunbookURL,
						runbookThanosReceiveLimitsHighMetaMonitoringQueriesFailureRate,
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing for {{$value | humanize}}% of meta monitoring queries.",
						"Thanos Receive {{$labels.job}} in {{$labels.namespace}} is failing for {{$value | humanize}}% of meta monitoring queries.",
						"Thanos Receive has not been able to update the number of head series.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveTenantLimitedByHeadSeries",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Sum(
						promqlbuilder.Increase(
							matrix.New(
								vector.New(
									vector.WithMetricName("thanos_receive_head_series_limited_requests_total"),
									vector.WithLabelMatchers(
										label.New("job").EqualRegexp(t.ReceiveRouterServiceSelector),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
					).By("namespace", "job", "tenant"),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.ReceiveDashboardURL,
						t.RunbookURL,
						runbookThanosReceiveTenantLimitedByHeadSeries,
						"Thanos Receive tenant {{$labels.tenant}} in {{$labels.namespace}} is limited by head series.",
						"Thanos Receive tenant {{$labels.tenant}} in {{$labels.namespace}} is limited by head series.",
						"Thanos Receive tenant is limited by head series.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}

func (t ThanosRulesConfig) ThanosStoreGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.AddRule(
			"ThanosStoreGrpcErrorRate",
			alerting.Expr(
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
												label.New("job").EqualRegexp(t.StoreServiceSelector),
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
												label.New("job").EqualRegexp(t.StoreServiceSelector),
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
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.StoreDashboardURL,
						t.RunbookURL,
						runbookThanosStoreGrpcErrorRate,
						"Thanos Store {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
						"Thanos Store {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
						"Thanos Store is failing to handle gRPC requests.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosStoreBucketHighOperationFailures",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Mul(
						promqlbuilder.Div(
							promqlbuilder.Sum(
								promqlbuilder.Rate(
									matrix.New(
										vector.New(
											vector.WithMetricName("thanos_objstore_bucket_operation_failures_total"),
											vector.WithLabelMatchers(
												label.New("job").EqualRegexp(t.StoreServiceSelector),
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
												label.New("job").EqualRegexp(t.StoreServiceSelector),
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
			),
			alerting.For("15m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.StoreDashboardURL,
						t.RunbookURL,
						runbookThanosStoreBucketHighOperationFailures,
						"Thanos Store {{$labels.job}} in {{$labels.namespace}} Bucket is failing to execute {{$value | humanize}}% of operations.",
						"Thanos Store {{$labels.job}} in {{$labels.namespace}} Bucket is failing to execute {{$value | humanize}}% of operations.",
						"Thanos Store Bucket is failing to execute operations.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosStoreObjstoreOperationLatencyHigh",
			alerting.Expr(
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
												label.New("job").EqualRegexp(t.StoreServiceSelector),
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
											label.New("job").EqualRegexp(t.StoreServiceSelector),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("namespace", "job"),
						promqlbuilder.NewNumber(0),
					),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.StoreDashboardURL,
						t.RunbookURL,
						runbookThanosStoreObjstoreOperationLatencyHigh,
						"Thanos Store {{$labels.job}} in {{$labels.namespace}} Bucket has a 99th percentile latency of {{$value}} seconds for the bucket operations.",
						"Thanos Store {{$labels.job}} in {{$labels.namespace}} Bucket has a 99th percentile latency of {{$value}} seconds for the bucket operations.",
						"Thanos Store is having high latency for bucket operations.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}

func (t ThanosRulesConfig) ThanosRuleGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.AddRule(
			"ThanosRuleQueueIsDroppingAlerts",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Sum(
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("thanos_alert_queue_alerts_dropped_total"),
									vector.WithLabelMatchers(
										label.New("job").EqualRegexp(t.RulerServiceSelector),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
					).By("namespace", "job", "instance"),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.RuleDashboardURL,
						t.RunbookURL,
						runbookThanosRuleQueueIsDroppingAlerts,
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} is failing to queue rulehelpers.",
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} is failing to queue rulehelpers.",
						"Thanos Rule is failing to queue rulehelpers.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRuleSenderIsFailingAlerts",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Sum(
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("thanos_alert_sender_alerts_dropped_total"),
									vector.WithLabelMatchers(
										label.New("job").EqualRegexp(t.RulerServiceSelector),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
					).By("namespace", "job", "instance"),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.RuleDashboardURL,
						t.RunbookURL,
						runbookThanosRuleSenderIsFailingAlerts,
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} is failing to send alerts to alertmanager.",
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} is failing to send alerts to alertmanager.",
						"Thanos Rule is failing to send alerts to alertmanager.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRuleHighRuleEvaluationFailures",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Mul(
						promqlbuilder.Div(
							promqlbuilder.Sum(
								promqlbuilder.Rate(
									matrix.New(
										vector.New(
											vector.WithMetricName("prometheus_rule_evaluation_failures_total"),
											vector.WithLabelMatchers(
												label.New("job").EqualRegexp(t.RulerServiceSelector),
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
												label.New("job").EqualRegexp(t.RulerServiceSelector),
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
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.RuleDashboardURL,
						t.RunbookURL,
						runbookThanosRuleHighRuleEvaluationFailures,
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} is failing to evaluate rules.",
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} is failing to evaluate rules.",
						"Thanos Rule is failing to evaluate rules.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRuleHighRuleEvaluationWarnings",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Sum(
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("thanos_rule_evaluation_with_warnings_total"),
									vector.WithLabelMatchers(
										label.New("job").EqualRegexp(t.RulerServiceSelector),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
					).By("namespace", "job", "instance"),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("15m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.RuleDashboardURL,
						t.RunbookURL,
						runbookThanosRuleHighRuleEvaluationWarnings,
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} has high number of evaluation warnings.",
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} has high number of evaluation warnings.",
						"Thanos Rule has high number of evaluation warnings.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRuleRuleEvaluationLatencyHigh",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Sum(
						vector.New(
							vector.WithMetricName("prometheus_rule_group_last_duration_seconds"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp(t.RulerServiceSelector),
							),
						),
					).By("namespace", "job", "instance", "rule_group"),
					promqlbuilder.Sum(
						vector.New(
							vector.WithMetricName("prometheus_rule_group_interval_seconds"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp(t.RulerServiceSelector),
							),
						),
					).By("namespace", "job", "instance", "rule_group"),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.RuleDashboardURL,
						t.RunbookURL,
						runbookThanosRuleRuleEvaluationLatencyHigh,
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} has higher evaluation latency than interval for {{$labels.rule_group}}.",
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} has higher evaluation latency than interval for {{$labels.rule_group}}.",
						"Thanos Rule has high rule evaluation latency.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRuleGrpcErrorRate",
			alerting.Expr(
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
												label.New("job").EqualRegexp(t.RulerServiceSelector),
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
												label.New("job").EqualRegexp(t.RulerServiceSelector),
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
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.RuleDashboardURL,
						t.RunbookURL,
						runbookThanosRuleGrpcErrorRate,
						"Thanos Rule {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
						"Thanos Rule {{$labels.job}} in {{$labels.namespace}} is failing to handle {{$value | humanize}}% of requests.",
						"Thanos Rule is failing to handle grpc requests.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRuleConfigReloadFailure",
			alerting.Expr(
				promqlbuilder.Neq(
					promqlbuilder.Avg(
						vector.New(
							vector.WithMetricName("thanos_rule_config_last_reload_successful"),
							vector.WithLabelMatchers(
								label.New("job").EqualRegexp(t.RulerServiceSelector),
							),
						),
					).By("namespace", "job", "instance"),
					promqlbuilder.NewNumber(1),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.RuleDashboardURL,
						t.RunbookURL,
						runbookThanosRuleConfigReloadFailure,
						"Thanos Rule {{$labels.job}} in {{$labels.namespace}} has not been able to reload its configuration.",
						"Thanos Rule {{$labels.job}} in {{$labels.namespace}} has not been able to reload its configuration.",
						"Thanos Rule has not been able to reload configuration.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRuleQueryHighDNSFailures",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Mul(
						promqlbuilder.Div(
							promqlbuilder.Sum(
								promqlbuilder.Rate(
									matrix.New(
										vector.New(
											vector.WithMetricName("thanos_rule_query_apis_dns_failures_total"),
											vector.WithLabelMatchers(
												label.New("job").EqualRegexp(t.RulerServiceSelector),
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
												label.New("job").EqualRegexp(t.RulerServiceSelector),
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
			),
			alerting.For("15m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.RuleDashboardURL,
						t.RunbookURL,
						runbookThanosRuleQueryHighDNSFailures,
						"Thanos Rule {{$labels.job}} in {{$labels.namespace}} has {{$value | humanize}}% of failing DNS queries for query endpoints.",
						"Thanos Rule {{$labels.job}} in {{$labels.namespace}} has {{$value | humanize}}% of failing DNS queries for query endpoints.",
						"Thanos Rule is having high number of DNS failures.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRuleAlertmanagerHighDNSFailures",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Mul(
						promqlbuilder.Div(
							promqlbuilder.Sum(
								promqlbuilder.Rate(
									matrix.New(
										vector.New(
											vector.WithMetricName("thanos_rule_alertmanagers_dns_failures_total"),
											vector.WithLabelMatchers(
												label.New("job").EqualRegexp(t.RulerServiceSelector),
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
												label.New("job").EqualRegexp(t.RulerServiceSelector),
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
			),
			alerting.For("15m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "medium",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.RuleDashboardURL,
						t.RunbookURL,
						runbookThanosRuleAlertmanagerHighDNSFailures,
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} has {{$value | humanize}}% of failing DNS queries for Alertmanager endpoints.",
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} has {{$value | humanize}}% of failing DNS queries for Alertmanager endpoints.",
						"Thanos Rule is having high number of DNS failures.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRuleNoEvaluationFor10Intervals",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Sub(
						promqlbuilder.Time(),
						promqlbuilder.Max(
							vector.New(
								vector.WithMetricName("prometheus_rule_group_last_evaluation_timestamp_seconds"),
								vector.WithLabelMatchers(
									label.New("job").EqualRegexp(t.RulerServiceSelector),
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
									label.New("job").EqualRegexp(t.RulerServiceSelector),
								),
							),
						).By("namespace", "job", "instance", "group"),
					),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "high",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.RuleDashboardURL,
						t.RunbookURL,
						runbookThanosRuleNoEvaluationFor10Intervals,
						"Thanos Rule {{$labels.job}} in {{$labels.namespace}} has rule groups that did not evaluate for at least 10x of their expected interval.",
						"Thanos Rule {{$labels.job}} in {{$labels.namespace}} has rule groups that did not evaluate for at least 10x of their expected interval.",
						"Thanos Rule has rule groups that did not evaluate for 10 intervals.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosNoRuleEvaluations",
			alerting.Expr(
				promqlbuilder.And(
					promqlbuilder.Lte(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("prometheus_rule_evaluations_total"),
										vector.WithLabelMatchers(
											label.New("job").EqualRegexp(t.RulerServiceSelector),
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
									label.New("job").EqualRegexp(t.RulerServiceSelector),
								),
							),
						).By("namespace", "job", "instance"),
						promqlbuilder.NewNumber(0),
					),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  t.ServiceLabelValue,
						"severity": "critical",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.RuleDashboardURL,
						t.RunbookURL,
						runbookThanosNoRuleEvaluations,
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} did not perform any rule evaluations in the past 10 minutes.",
						"Thanos Rule {{$labels.instance}} in {{$labels.namespace}} did not perform any rule evaluations in the past 10 minutes.",
						"Thanos Rule did not perform any rule evaluations.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}
