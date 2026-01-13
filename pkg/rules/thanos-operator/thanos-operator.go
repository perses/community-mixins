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

package thanosoperator

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
	runbookThanosOperatorDown                       = "#thanosoperatordown"
	runbookThanosOperatorHighReconcileErrorRate     = "#thanosoperatorhighreconcileerrorrate"
	runbookThanosOperatorReconcileStuck             = "#thanosoperatorreconcilestuck"
	runbookThanosOperatorWorkQueueGrowth            = "#thanosoperatorworkqueuegrowth"
	runbookThanosOperatorSlowReconciliation         = "#thanosoperatorslowreconciliation"
	runbookThanosQueryNoEndpointsConfigured         = "#thanosquerynoendpointsconfigured"
	runbookThanosQueryServiceWatchReconcileStorm    = "#thanosqueryservicewatchreconcilestorm"
	runbookThanosReceiveNoHashringsConfigured       = "#thanosreceivenohashringsconfigured"
	runbookThanosReceiveHashringNoEndpoints         = "#thanosreceivehashringnoendpoints"
	runbookThanosReceiveHashringConfigurationChange = "#thanosreceivehashringconfigurationchange"
	runbookThanosReceiveEndpointReconcileStorm      = "#thanosreceiveendpointreconcilestorm"
	runbookThanosRulerNoQueryEndpointsConfigured    = "#thanosrulernoqueryendpointsconfigured"
	runbookThanosRulerNoRulesConfigured             = "#thanosrulernorulesconfigured"
	runbookThanosRulerConfigMapCreationFailures     = "#thanosrulerconfigmapcreationfailures"
	runbookThanosRulerHighConfigMapCreationRate     = "#thanosrulerhighconfigmapcreationrate"
	runbookThanosRulerWatchReconcileStorm           = "#thanosrulerwatchreconcilestorm"
	runbookThanosStoreNoShardsConfigured            = "#thanosstenoreshardsconfigured"
	runbookThanosStoreShardCreationFailures         = "#thanosstoreshardcreationfailures"
	runbookThanosCompactNoShardsConfigured          = "#thanoscompactnoshardsconfigured"
	runbookThanosCompactShardCreationFailures       = "#thanoscompactshardcreationfailures"
	runbookThanosResourcePausedForLong              = "#thanosresourcepausedforlong"
	runbookThanosOperatorHighWorkqueueRetries       = "#thanosoperatorhighworkqueueretries"
	runbookThanosOperatorLongWorkqueueLatency       = "#thanosoperatorlongworkqueuelatency"
)

type ThanosOperatorRulesConfig struct {
	RunbookURL   string
	DashboardURL string

	ServiceLabelValue      string
	MetricsServiceSelector string

	AdditionalAlertLabels      map[string]string
	AdditionalAlertAnnotations map[string]string
}

type ThanosOperatorRulesConfigOption func(*ThanosOperatorRulesConfig)

func WithRunbookURL(runbookURL string) ThanosOperatorRulesConfigOption {
	return func(config *ThanosOperatorRulesConfig) {
		config.RunbookURL = runbookURL
	}
}

func WithMetricsServiceSelector(metricsServiceSelector string) ThanosOperatorRulesConfigOption {
	return func(config *ThanosOperatorRulesConfig) {
		if metricsServiceSelector == "" {
			metricsServiceSelector = "thanos-operator-controller-manager-metrics-service"
		}
		config.MetricsServiceSelector = metricsServiceSelector
	}
}

func WithDashboardURL(dashboardURL string) ThanosOperatorRulesConfigOption {
	return func(config *ThanosOperatorRulesConfig) {
		config.DashboardURL = dashboardURL
	}
}

func WithServiceLabelValue(serviceLabelValue string) ThanosOperatorRulesConfigOption {
	return func(config *ThanosOperatorRulesConfig) {
		config.ServiceLabelValue = serviceLabelValue
	}
}

func WithAdditionalAlertLabels(additionalAlertLabels map[string]string) ThanosOperatorRulesConfigOption {
	return func(config *ThanosOperatorRulesConfig) {
		config.AdditionalAlertLabels = additionalAlertLabels
	}
}

func WithAdditionalAlertAnnotations(additionalAlertAnnotations map[string]string) ThanosOperatorRulesConfigOption {
	return func(config *ThanosOperatorRulesConfig) {
		config.AdditionalAlertAnnotations = additionalAlertAnnotations
	}
}

// NewThanosOperatorRulesBuilder creates a new Thanos Operator rules builder.
func NewThanosOperatorRulesBuilder(
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	options ...ThanosOperatorRulesConfigOption,
) (promtheusrule.Builder, error) {
	config := ThanosOperatorRulesConfig{
		MetricsServiceSelector: "thanos-operator-controller-manager-metrics-service",
	}

	for _, option := range options {
		option(&config)
	}

	promRule, err := promtheusrule.New(
		"thanos-operator-alerts",
		namespace,
		promtheusrule.Labels(labels),
		promtheusrule.Annotations(annotations),
		promtheusrule.AddRuleGroup(
			"thanos-operator.general",
			config.ThanosOperatorGeneralGroup()...,
		),
		promtheusrule.AddRuleGroup(
			"thanos-operator.query",
			config.ThanosOperatorQueryGroup()...,
		),
		promtheusrule.AddRuleGroup(
			"thanos-operator.receive",
			config.ThanosOperatorReceiveGroup()...,
		),
		promtheusrule.AddRuleGroup(
			"thanos-operator.ruler",
			config.ThanosOperatorRulerGroup()...,
		),
		promtheusrule.AddRuleGroup(
			"thanos-operator.store",
			config.ThanosOperatorStoreGroup()...,
		),
		promtheusrule.AddRuleGroup(
			"thanos-operator.compact",
			config.ThanosOperatorCompactGroup()...,
		),
		promtheusrule.AddRuleGroup(
			"thanos-operator.paused",
			config.ThanosOperatorPausedGroup()...,
		),
		promtheusrule.AddRuleGroup(
			"thanos-operator.workqueue",
			config.ThanosOperatorWorkqueueGroup()...,
		),
	)

	return promRule, err
}

// BuildThanosOperatorRules builds the Thanos Operator rules
func BuildThanosOperatorRules(
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	options ...ThanosOperatorRulesConfigOption,
) rulehelpers.RuleResult {
	promRule, err := NewThanosOperatorRulesBuilder(namespace, labels, annotations, options...)
	if err != nil {
		return rulehelpers.NewRuleResult(nil, err).Component("thanos-operator")
	}

	return rulehelpers.NewRuleResult(
		&promRule.PrometheusRule,
		nil,
	).Component("thanos-operator")
}

func (t ThanosOperatorRulesConfig) ThanosOperatorGeneralGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.Interval("30s"),
		rulegroup.AddRule(
			"ThanosOperatorDown",
			alerting.Expr(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("up"),
						vector.WithLabelMatchers(
							label.New("job").Equal(t.MetricsServiceSelector),
						),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "critical",
						"component": "thanos-operator",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosOperatorDown,
						"The Thanos Operator has been down for more than 5 minutes. No reconciliation is happening.",
						"Thanos resources are not being reconciled. Configuration changes will not be applied.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosOperatorHighReconcileErrorRate",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Div(
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("controller_runtime_reconcile_errors_total"),
										vector.WithLabelMatchers(
											label.New("job").Equal(t.MetricsServiceSelector),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("controller"),
						promqlbuilder.Sum(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("controller_runtime_reconcile_total"),
										vector.WithLabelMatchers(
											label.New("job").Equal(t.MetricsServiceSelector),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
						).By("controller"),
					),
					promqlbuilder.NewNumber(0.1),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":    t.ServiceLabelValue,
						"severity":   "warning",
						"component":  "thanos-operator",
						"controller": "{{ $labels.controller }}",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosOperatorHighReconcileErrorRate,
						"Controller {{ $labels.controller }} for resource {{ $labels.namespace }}/{{ $labels.resource }} has a high reconciliation error rate of {{ $value | humanizePercentage }} over the last 10 minutes.",
						"Resources managed by this controller may not be correctly configured or updated.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosOperatorReconcileStuck",
			alerting.Expr(
				promqlbuilder.And(
					promqlbuilder.Eqlc(
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("controller_runtime_reconcile_total"),
									vector.WithLabelMatchers(
										label.New("job").Equal(t.MetricsServiceSelector),
									),
								),
								matrix.WithRange(10*time.Minute),
							),
						),
						promqlbuilder.NewNumber(0),
					),
					promqlbuilder.Gtr(
						vector.New(
							vector.WithMetricName("workqueue_depth"),
							vector.WithLabelMatchers(
								label.New("job").Equal(t.MetricsServiceSelector),
							),
						),
						promqlbuilder.NewNumber(0),
					),
				).On("controller", "job"),
			),
			alerting.For("15m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":    t.ServiceLabelValue,
						"severity":   "warning",
						"component":  "thanos-operator",
						"controller": "{{ $labels.controller }}",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosOperatorReconcileStuck,
						"Workqueue for {{ $labels.name }} has items but no reconciliations are happening. Controller appears stuck.",
						"Changes to Thanos resources are not being processed.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosOperatorWorkQueueGrowth",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.LastOverTime(
						matrix.New(
							vector.New(
								vector.WithMetricName("workqueue_depth"),
								vector.WithLabelMatchers(
									label.New("job").Equal(t.MetricsServiceSelector),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
					promqlbuilder.NewNumber(100),
				),
			),
			alerting.For("15m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":    t.ServiceLabelValue,
						"severity":   "warning",
						"component":  "thanos-operator",
						"controller": "{{ $labels.controller }}",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosOperatorWorkQueueGrowth,
						"Workqueue depth for {{ $labels.name }} is {{ $value }}, indicating the controller cannot keep up with events.",
						"Reconciliation is falling behind. Configuration updates may be delayed.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosOperatorSlowReconciliation",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.HistogramQuantile(
						0.99,
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("controller_runtime_reconcile_time_seconds_bucket"),
									vector.WithLabelMatchers(
										label.New("job").Equal(t.MetricsServiceSelector),
									),
								),
								matrix.WithRange(10*time.Minute),
							),
						),
					),
					promqlbuilder.NewNumber(60),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":    t.ServiceLabelValue,
						"severity":   "warning",
						"component":  "thanos-operator",
						"controller": "{{ $labels.controller }}",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosOperatorSlowReconciliation,
						"P99 reconciliation time for {{ $labels.controller }} is {{ $value | humanizeDuration }}, which is slow",
						"Configuration changes are taking longer than expected to apply.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}

func (t ThanosOperatorRulesConfig) ThanosOperatorQueryGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.Interval("30s"),
		rulegroup.AddRule(
			"ThanosQueryNoEndpointsConfigured",
			alerting.Expr(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("thanos_operator_query_endpoints_configured"),
						vector.WithLabelMatchers(
							label.New("job").Equal(t.MetricsServiceSelector),
						),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "warning",
						"component": "thanos-query",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosQueryNoEndpointsConfigured,
						"ThanosQuery resource {{ $labels.namespace }}/{{ $labels.resource }} has no store endpoints configured.",
						"Query component cannot retrieve data from any stores. Queries will return no results.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),

		rulegroup.AddRule(
			"ThanosQueryServiceWatchReconcileStorm",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Rate(
						matrix.New(
							vector.New(
								vector.WithMetricName("thanos_operator_query_service_event_reconciliations_total"),
								vector.WithLabelMatchers(
									label.New("job").Equal(t.MetricsServiceSelector),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
					promqlbuilder.NewNumber(2),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "warning",
						"component": "thanos-query",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosQueryServiceWatchReconcileStorm,
						"ThanosQuery {{ $labels.namespace }}/{{ $labels.resource }} is reconciling {{ $value | humanize }} times/sec due to service events.",
						"Excessive reconciliations may indicate service churn or configuration issues.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}

func (t ThanosOperatorRulesConfig) ThanosOperatorReceiveGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.Interval("30s"),
		rulegroup.AddRule(
			"ThanosReceiveNoHashringsConfigured",
			alerting.Expr(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("thanos_operator_receive_hashrings_configured"),
						vector.WithLabelMatchers(
							label.New("job").Equal(t.MetricsServiceSelector),
						),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "warning",
						"component": "thanos-receive",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosReceiveNoHashringsConfigured,
						"ThanosReceive resource {{ $labels.namespace }}/{{ $labels.resource }} has no hashrings configured.",
						"Receive component cannot accept remote write data without hashring configuration.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveHashringNoEndpoints",
			alerting.Expr(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("thanos_operator_receive_hashring_endpoints_configured"),
						vector.WithLabelMatchers(
							label.New("job").Equal(t.MetricsServiceSelector),
						),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "critical",
						"component": "thanos-receive",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosReceiveHashringNoEndpoints,
						"Hashring {{ $labels.hashring }} for ThanosReceive {{ $labels.namespace }}/{{ $labels.resource }} has no endpoints configured.",
						"Data cannot be distributed to this hashring. Remote write data may be lost or rejected.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveHashringConfigurationChange",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Abs(
						promqlbuilder.Delta(
							matrix.New(
								vector.New(
									vector.WithMetricName("thanos_operator_receive_hashring_hash"),
									vector.WithLabelMatchers(
										label.New("job").Equal(t.MetricsServiceSelector),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("0m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "info",
						"component": "thanos-receive",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosReceiveHashringConfigurationChange,
						"The hashring configuration for ThanosReceive {{ $labels.namespace }}/{{ $labels.resource }} has changed.",
						"Data distribution pattern has changed. This may cause temporary inconsistencies.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosReceiveEndpointReconcileStorm",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Rate(
						matrix.New(
							vector.New(
								vector.WithMetricName("thanos_operator_receive_endpoint_event_reconciliations_total"),
								vector.WithLabelMatchers(
									label.New("job").Equal(t.MetricsServiceSelector),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
					promqlbuilder.NewNumber(2),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "warning",
						"component": "thanos-receive",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosReceiveEndpointReconcileStorm,
						"ThanosReceive {{ $labels.namespace }}/{{ $labels.resource }} is reconciling {{ $value | humanize }} times/sec due to endpoint events.",
						"Excessive reconciliations may indicate endpoint churn or configuration issues.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}

func (t ThanosOperatorRulesConfig) ThanosOperatorRulerGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.Interval("30s"),
		rulegroup.AddRule(
			"ThanosRulerNoQueryEndpointsConfigured",
			alerting.Expr(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("thanos_operator_ruler_query_endpoints_configured"),
						vector.WithLabelMatchers(
							label.New("job").Equal(t.MetricsServiceSelector),
						),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "warning",
						"component": "thanos-ruler",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosRulerNoQueryEndpointsConfigured,
						"ThanosRuler resource {{ $labels.namespace }}/{{ $labels.resource }} has no query endpoints configured.",
						"Ruler cannot query data for rule evaluation. Recording and alerting rules will not work.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRulerNoPrometheusRulesConfigured",
			alerting.Expr(
				promqlbuilder.Or(
					promqlbuilder.Eqlc(
						vector.New(
							vector.WithMetricName("thanos_operator_ruler_promrules_found"),
							vector.WithLabelMatchers(
								label.New("job").Equal(t.MetricsServiceSelector),
							),
						),
						promqlbuilder.NewNumber(0),
					),
					promqlbuilder.Absent(
						vector.New(
							vector.WithMetricName("thanos_operator_ruler_promrules_found"),
							vector.WithLabelMatchers(
								label.New("job").Equal(t.MetricsServiceSelector),
							),
						),
					),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "warning",
						"component": "thanos-ruler",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosRulerNoRulesConfigured,
						"No PrometheusRules found for ThanosRuler {{ $labels.namespace }}/{{ $labels.resource }}.",
						"No rules from PrometheusRules are being evaluated. This may be expected if no rules have been defined yet.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRulerNoRulesConfigured",
			alerting.Expr(
				promqlbuilder.Or(
					promqlbuilder.Eqlc(
						vector.New(
							vector.WithMetricName("thanos_operator_ruler_rulefiles_configured"),
							vector.WithLabelMatchers(
								label.New("job").Equal(t.MetricsServiceSelector),
							),
						),
						promqlbuilder.NewNumber(0),
					),
					promqlbuilder.Absent(
						vector.New(
							vector.WithMetricName("thanos_operator_ruler_rulefiles_configured"),
							vector.WithLabelMatchers(
								label.New("job").Equal(t.MetricsServiceSelector),
							),
						),
					),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "warning",
						"component": "thanos-ruler",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosRulerNoRulesConfigured,
						"No rule configmaps found for ThanosRuler {{ $labels.namespace }}/{{ $labels.resource }}.",
						"No rules are being evaluated as no rule configmaps were found. This may be expected if no rules have been defined yet.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRulerConfigMapCreationFailures",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Rate(
						matrix.New(
							vector.New(
								vector.WithMetricName("thanos_operator_ruler_cfgmaps_creation_failures_total"),
								vector.WithLabelMatchers(
									label.New("job").Equal(t.MetricsServiceSelector),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "critical",
						"component": "thanos-ruler",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosRulerConfigMapCreationFailures,
						"ThanosRuler controller is failing to create ConfigMaps for {{ $labels.namespace }}/{{ $labels.resource }}.",
						"PrometheusRules cannot be loaded into the Ruler. Rules will not be evaluated.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRulerHighConfigMapCreationRate",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Rate(
						matrix.New(
							vector.New(
								vector.WithMetricName("thanos_operator_ruler_cfgmaps_created_total"),
								vector.WithLabelMatchers(
									label.New("job").Equal(t.MetricsServiceSelector),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
					promqlbuilder.NewNumber(1),
				),
			),
			alerting.For("15m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "warning",
						"component": "thanos-ruler",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosRulerHighConfigMapCreationRate,
						"ThanosRuler {{ $labels.namespace }}/{{ $labels.resource }} is creating ConfigMaps at {{ $value | humanize }}/sec.",
						"Excessive ConfigMap updates may indicate rule churn and cause unnecessary Ruler reloads.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosRulerWatchReconcileStorm",
			alerting.Expr(
				promqlbuilder.Or(
					promqlbuilder.Or(
						promqlbuilder.Gtr(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_operator_ruler_service_event_reconciliations_total"),
										vector.WithLabelMatchers(
											label.New("job").Equal(t.MetricsServiceSelector),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
							promqlbuilder.NewNumber(2),
						),
						promqlbuilder.Gtr(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("thanos_operator_ruler_cfgmap_event_reconciliations_total"),
										vector.WithLabelMatchers(
											label.New("job").Equal(t.MetricsServiceSelector),
										),
									),
									matrix.WithRange(5*time.Minute),
								),
							),
							promqlbuilder.NewNumber(2),
						),
					),
					promqlbuilder.Gtr(
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("thanos_operator_ruler_promrule_event_reconciliations_total"),
									vector.WithLabelMatchers(
										label.New("job").Equal(t.MetricsServiceSelector),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
						promqlbuilder.NewNumber(2),
					),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "warning",
						"component": "thanos-ruler",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosRulerWatchReconcileStorm,
						"ThanosRuler {{ $labels.namespace }}/{{ $labels.resource }} is experiencing high reconciliation rate.",
						"Excessive reconciliations may indicate resource churn or configuration issues.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}

func (t ThanosOperatorRulesConfig) ThanosOperatorStoreGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.Interval("30s"),
		rulegroup.AddRule(
			"ThanosStoreNoShardsConfigured",
			alerting.Expr(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("thanos_operator_store_shards_configured"),
						vector.WithLabelMatchers(
							label.New("job").Equal(t.MetricsServiceSelector),
						),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "info",
						"component": "thanos-store",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosStoreNoShardsConfigured,
						"ThanosStore resource {{ $labels.namespace }}/{{ $labels.resource }} has 0 shards configured.",
						"This may be expected for single-instance stores. For sharded deployments, data queries may fail.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosStoreShardCreationFailures",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Rate(
						matrix.New(
							vector.New(
								vector.WithMetricName("thanos_operator_store_shards_creation_update_failures_total"),
								vector.WithLabelMatchers(
									label.New("job").Equal(t.MetricsServiceSelector),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "critical",
						"component": "thanos-store",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosStoreShardCreationFailures,
						"ThanosStore controller is failing to create/update shards for {{ $labels.namespace }}/{{ $labels.resource }}.",
						"Store shards cannot be created. Historical data queries may be incomplete or fail.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}

func (t ThanosOperatorRulesConfig) ThanosOperatorCompactGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.Interval("30s"),
		rulegroup.AddRule(
			"ThanosCompactNoShardsConfigured",
			alerting.Expr(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("thanos_operator_compact_shards_configured"),
						vector.WithLabelMatchers(
							label.New("job").Equal(t.MetricsServiceSelector),
						),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "info",
						"component": "thanos-compact",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosCompactNoShardsConfigured,
						"ThanosCompact resource {{ $labels.namespace }}/{{ $labels.resource }} has 0 shards configured.",
						"This may be expected for single-instance compactors. For sharded deployments, compaction may not work.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosCompactShardCreationFailures",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Rate(
						matrix.New(
							vector.New(
								vector.WithMetricName("thanos_operator_compact_shards_creation_update_failures_total"),
								vector.WithLabelMatchers(
									label.New("job").Equal(t.MetricsServiceSelector),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":   t.ServiceLabelValue,
						"severity":  "critical",
						"component": "thanos-compact",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosCompactShardCreationFailures,
						"ThanosCompact controller is failing to create/update shards for {{ $labels.namespace }}/{{ $labels.resource }}.",
						"Compactor shards cannot be created. Data compaction will not occur, leading to increased storage costs and slower queries.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}

func (t ThanosOperatorRulesConfig) ThanosOperatorPausedGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.Interval("30s"),
		rulegroup.AddRule(
			"ThanosResourcePausedForLong",
			alerting.Expr(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("thanos_operator_paused"),
						vector.WithLabelMatchers(
							label.New("job").Equal(t.MetricsServiceSelector),
						),
					),
					promqlbuilder.NewNumber(1),
				),
			),
			alerting.For("24h"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":    t.ServiceLabelValue,
						"severity":   "info",
						"component":  "{{ $labels.component }}",
						"controller": "{{ $labels.controller }}",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosResourcePausedForLong,
						"{{ $labels.component }} resource {{ $labels.namespace }}/{{ $labels.resource }} has been in paused state for over 24 hours.",
						"No reconciliation is happening for this resource. Configuration changes are not being applied.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}

func (t ThanosOperatorRulesConfig) ThanosOperatorWorkqueueGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.Interval("30s"),
		rulegroup.AddRule(
			"ThanosOperatorHighWorkqueueRetries",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Rate(
						matrix.New(
							vector.New(
								vector.WithMetricName("workqueue_retries_total"),
								vector.WithLabelMatchers(
									label.New("job").Equal(t.MetricsServiceSelector),
								),
							),
							matrix.WithRange(10*time.Minute),
						),
					),
					promqlbuilder.NewNumber(0.5),
				),
			),
			alerting.For("15m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":    t.ServiceLabelValue,
						"severity":   "warning",
						"component":  "thanos-operator",
						"controller": "{{ $labels.controller }}",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosOperatorHighWorkqueueRetries,
						"Workqueue {{ $labels.name }} has {{ $value | humanize }} retries/sec.",
						"Items are being retried frequently, indicating persistent errors or resource issues.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"ThanosOperatorLongWorkqueueLatency",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.HistogramQuantile(
						0.99,
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("workqueue_queue_duration_seconds_bucket"),
									vector.WithLabelMatchers(
										label.New("job").Equal(t.MetricsServiceSelector),
									),
								),
								matrix.WithRange(10*time.Minute),
							),
						),
					),
					promqlbuilder.NewNumber(60),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":    t.ServiceLabelValue,
						"severity":   "warning",
						"component":  "thanos-operator",
						"controller": "{{ $labels.controller }}",
					},
					t.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						t.DashboardURL,
						t.RunbookURL,
						runbookThanosOperatorLongWorkqueueLatency,
						"P99 queue wait time for {{ $labels.name }} is {{ $value | humanizeDuration }}.",
						"Items are waiting too long in queue before processing. Reconciliation is delayed.",
					),
					t.AdditionalAlertAnnotations,
				),
			),
		),
	}
}
