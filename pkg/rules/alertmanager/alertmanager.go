package alertmanager

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
	runbookAlertmanagerFailedReload                     = "#alertmanagerfailedreload"
	runbookAlertmanagerMembersInconsistent              = "#alertmanagermembersinconsistent"
	runbookAlertmanagerFailedToSendAlerts               = "#alertmanagerfailedtosendalerts"
	runbookAlertmanagerClusterFailedToSendAlerts        = "#alertmanagerclusterfailedtosendalerts"
	runbookAlertmanagerClusterFailedToSendAlertsNonCrit = "#alertmanagerclusterfailedtosendalertsnonCrit"
	runbookAlertmanagerConfigInconsistent               = "#alertmanagerconfiginconsistent"
	runbookAlertmanagerClusterDown                      = "#alertmanagerclusterdown"
	runbookAlertmanagerClusterCrashlooping              = "#alertmanagerclustercrashlooping"
)

type AlertmanagerRulesConfig struct {
	RunbookURL   string
	DashboardURL string

	ServiceLabelValue              string
	AlertmanagerServiceSelector    string
	CriticalIntegrationSelector    string
	NonCriticalIntegrationSelector string

	AdditionalAlertLabels      map[string]string
	AdditionalAlertAnnotations map[string]string
}

type AlertmanagerRulesConfigOption func(*AlertmanagerRulesConfig)

func WithRunbookURL(runbookURL string) AlertmanagerRulesConfigOption {
	return func(alertmanagerRulesConfig *AlertmanagerRulesConfig) {
		alertmanagerRulesConfig.RunbookURL = runbookURL
	}
}

func WithAlertmanagerServiceSelector(alertmanagerServiceSelector string) AlertmanagerRulesConfigOption {
	return func(alertmanagerRulesConfig *AlertmanagerRulesConfig) {
		if alertmanagerServiceSelector == "" {
			alertmanagerServiceSelector = "alertmanager"
		}
		alertmanagerRulesConfig.AlertmanagerServiceSelector = alertmanagerServiceSelector
	}
}

func WithCriticalIntegrationSelectorRegexp(criticalIntegrationSelector string) AlertmanagerRulesConfigOption {
	return func(alertmanagerRulesConfig *AlertmanagerRulesConfig) {
		if criticalIntegrationSelector == "" {
			criticalIntegrationSelector = ".*"
		}
		alertmanagerRulesConfig.CriticalIntegrationSelector = criticalIntegrationSelector
	}
}

func WithNonCriticalIntegrationSelectorRegexp(nonCriticalIntegrationSelector string) AlertmanagerRulesConfigOption {
	return func(alertmanagerRulesConfig *AlertmanagerRulesConfig) {
		if nonCriticalIntegrationSelector == "" {
			nonCriticalIntegrationSelector = ".*"
		}
		alertmanagerRulesConfig.NonCriticalIntegrationSelector = nonCriticalIntegrationSelector
	}
}

func WithServiceLabelValue(serviceLabelValue string) AlertmanagerRulesConfigOption {
	return func(alertmanagerRulesConfig *AlertmanagerRulesConfig) {
		alertmanagerRulesConfig.ServiceLabelValue = serviceLabelValue
	}
}

func WithAdditionalAlertLabels(additionalAlertLabels map[string]string) AlertmanagerRulesConfigOption {
	return func(alertmanagerRulesConfig *AlertmanagerRulesConfig) {
		alertmanagerRulesConfig.AdditionalAlertLabels = additionalAlertLabels
	}
}

func WithAdditionalAlertAnnotations(additionalAlertAnnotations map[string]string) AlertmanagerRulesConfigOption {
	return func(alertmanagerRulesConfig *AlertmanagerRulesConfig) {
		alertmanagerRulesConfig.AdditionalAlertAnnotations = additionalAlertAnnotations
	}
}

func WithDashboardURL(dashboardURL string) AlertmanagerRulesConfigOption {
	return func(alertmanagerRulesConfig *AlertmanagerRulesConfig) {
		alertmanagerRulesConfig.DashboardURL = dashboardURL
	}
}

// BuildAlertmanagerRules builds the Alertmanager rules for the given namespace, dashboard URLs, runbook URL, labels, and annotations.
func BuildAlertmanagerRules(
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	options ...AlertmanagerRulesConfigOption,
) rulehelpers.RuleResult {
	alertmanagerRulesConfig := AlertmanagerRulesConfig{
		AlertmanagerServiceSelector:    "alertmanager",
		CriticalIntegrationSelector:    ".*",
		NonCriticalIntegrationSelector: ".*",
	}
	for _, option := range options {
		option(&alertmanagerRulesConfig)
	}

	promRule, err := promtheusrule.New(
		"alertmanager-rules",
		namespace,
		promtheusrule.Labels(labels),
		promtheusrule.Annotations(annotations),
		promtheusrule.AddRuleGroup(
			"alertmanager.rules",
			alertmanagerRulesConfig.AlertmanagerRulesGroup()...,
		),
	)

	if err != nil {
		return rulehelpers.NewRuleResult(nil, err).Component("alertmanager")
	}

	return rulehelpers.NewRuleResult(
		&promRule.PrometheusRule,
		nil,
	).Component("alertmanager")
}

func (a AlertmanagerRulesConfig) AlertmanagerRulesGroup() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.AddRule[alerting.Option](
			"AlertmanagerFailedReload",
			alerting.Expr(
				promqlbuilder.Eqlc(
					promqlbuilder.MaxOverTime(
						matrix.New(
							vector.New(
								vector.WithMetricName("alertmanager_config_last_reload_successful"),
								vector.WithLabelMatchers(
									label.New("job").Equal(a.AlertmanagerServiceSelector),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("10m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  a.ServiceLabelValue,
						"severity": "critical",
					},
					a.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						a.DashboardURL,
						a.RunbookURL,
						runbookAlertmanagerFailedReload,
						"Configuration has failed to load for {{$labels.instance}}.",
						"Configuration has failed to load for {{$labels.instance}}.",
						"Reloading an Alertmanager configuration has failed.",
					),
					a.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule[alerting.Option](
			"AlertmanagerMembersInconsistent",
			alerting.Expr(
				promqlbuilder.Lss(
					promqlbuilder.MaxOverTime(
						matrix.New(
							vector.New(
								vector.WithMetricName("alertmanager_cluster_members"),
								vector.WithLabelMatchers(
									label.New("job").Equal(a.AlertmanagerServiceSelector),
								),
							),
							matrix.WithRange(5*time.Minute),
						),
					),
					promqlbuilder.Count(
						promqlbuilder.MaxOverTime(
							matrix.New(
								vector.New(
									vector.WithMetricName("alertmanager_cluster_members"),
									vector.WithLabelMatchers(
										label.New("job").Equal(a.AlertmanagerServiceSelector),
									),
								),
								matrix.WithRange(5*time.Minute),
							),
						),
					).By("job"),
				).On("job").GroupLeft(),
			),
			alerting.For("15m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  a.ServiceLabelValue,
						"severity": "critical",
					},
					a.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						a.DashboardURL,
						a.RunbookURL,
						runbookAlertmanagerMembersInconsistent,
						"Alertmanager {{$labels.instance}} has only found {{ $value }} members of the {{$labels.job}} cluster.",
						"Alertmanager {{$labels.instance}} has only found {{ $value }} members of the {{$labels.job}} cluster.",
						"A member of an Alertmanager cluster has not found all other cluster members.",
					),
					a.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule[alerting.Option](
			"AlertmanagerFailedToSendAlerts",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Div(
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("alertmanager_notifications_failed_total"),
									vector.WithLabelMatchers(
										label.New("job").Equal(a.AlertmanagerServiceSelector),
									),
								),
								matrix.WithRange(15*time.Minute),
							),
						),
						promqlbuilder.Rate(
							matrix.New(
								vector.New(
									vector.WithMetricName("alertmanager_notifications_total"),
									vector.WithLabelMatchers(
										label.New("job").Equal(a.AlertmanagerServiceSelector),
									),
								),
								matrix.WithRange(15*time.Minute),
							),
						),
					).Ignoring("reason").GroupLeft(),
					promqlbuilder.NewNumber(0.01),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  a.ServiceLabelValue,
						"severity": "warning",
					},
					a.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						a.DashboardURL,
						a.RunbookURL,
						runbookAlertmanagerFailedToSendAlerts,
						"Alertmanager {{$labels.instance}} failed to send {{ $value | humanizePercentage }} of notifications to {{ $labels.integration }}.",
						"Alertmanager {{$labels.instance}} failed to send {{ $value | humanizePercentage }} of notifications to {{ $labels.integration }}.",
						"An Alertmanager instance failed to send notifications.",
					),
					a.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule[alerting.Option](
			"AlertmanagerClusterFailedToSendAlerts",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Min(
						promqlbuilder.Div(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("alertmanager_notifications_failed_total"),
										vector.WithLabelMatchers(
											label.New("job").Equal(a.AlertmanagerServiceSelector),
											label.New("integration").EqualRegexp(a.CriticalIntegrationSelector),
										),
									),
									matrix.WithRange(15*time.Minute),
								),
							),
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("alertmanager_notifications_total"),
										vector.WithLabelMatchers(
											label.New("job").Equal(a.AlertmanagerServiceSelector),
											label.New("integration").EqualRegexp(a.CriticalIntegrationSelector),
										),
									),
									matrix.WithRange(15*time.Minute),
								),
							),
						).Ignoring("reason").GroupLeft(),
					).By("job", "integration"),
					promqlbuilder.NewNumber(0.01),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  a.ServiceLabelValue,
						"severity": "critical",
					},
					a.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						a.DashboardURL,
						a.RunbookURL,
						runbookAlertmanagerClusterFailedToSendAlerts,
						"The minimum notification failure rate to {{ $labels.integration }} sent from any instance in the {{$labels.job}} cluster is {{ $value | humanizePercentage }}.",
						"The minimum notification failure rate to {{ $labels.integration }} sent from any instance in the {{$labels.job}} cluster is {{ $value | humanizePercentage }}.",
						"All Alertmanager instances in a cluster failed to send notifications to a critical integration.",
					),
					a.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule[alerting.Option](
			"AlertmanagerClusterFailedToSendAlerts",
			alerting.Expr(
				promqlbuilder.Gtr(
					promqlbuilder.Min(
						promqlbuilder.Div(
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("alertmanager_notifications_failed_total"),
										vector.WithLabelMatchers(
											label.New("job").Equal(a.AlertmanagerServiceSelector),
											label.New("integration").NotEqualRegexp(a.NonCriticalIntegrationSelector),
										),
									),
									matrix.WithRange(15*time.Minute),
								),
							),
							promqlbuilder.Rate(
								matrix.New(
									vector.New(
										vector.WithMetricName("alertmanager_notifications_total"),
										vector.WithLabelMatchers(
											label.New("job").Equal(a.AlertmanagerServiceSelector),
											label.New("integration").NotEqualRegexp(a.NonCriticalIntegrationSelector),
										),
									),
									matrix.WithRange(15*time.Minute),
								),
							),
						).Ignoring("reason").GroupLeft(),
					).By("job", "integration"),
					promqlbuilder.NewNumber(0.01),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  a.ServiceLabelValue,
						"severity": "warning",
					},
					a.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						a.DashboardURL,
						a.RunbookURL,
						runbookAlertmanagerClusterFailedToSendAlertsNonCrit,
						"The minimum notification failure rate to {{ $labels.integration }} sent from any instance in the {{$labels.job}} cluster is {{ $value | humanizePercentage }}.",
						"The minimum notification failure rate to {{ $labels.integration }} sent from any instance in the {{$labels.job}} cluster is {{ $value | humanizePercentage }}.",
						"All Alertmanager instances in a cluster failed to send notifications to a non-critical integration.",
					),
					a.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule[alerting.Option](
			"AlertmanagerConfigInconsistent",
			alerting.Expr(
				promqlbuilder.Neq(
					promqlbuilder.Count(
						promqlbuilder.CountValues(
							"config_hash",
							vector.New(
								vector.WithMetricName("alertmanager_config_hash"),
								vector.WithLabelMatchers(
									label.New("job").Equal(a.AlertmanagerServiceSelector),
								),
							),
						).By("job"),
					).By("job"),
					promqlbuilder.NewNumber(1),
				),
			),
			alerting.For("20m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  a.ServiceLabelValue,
						"severity": "critical",
					},
					a.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						a.DashboardURL,
						a.RunbookURL,
						runbookAlertmanagerConfigInconsistent,
						"Alertmanager instances within the {{$labels.job}} cluster have different configurations.",
						"Alertmanager instances within the {{$labels.job}} cluster have different configurations.",
						"Alertmanager instances within the same cluster have different configurations.",
					),
					a.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule[alerting.Option](
			"AlertmanagerClusterDown",
			alerting.Expr(
				promqlbuilder.Gte(
					promqlbuilder.Div(
						promqlbuilder.Count(
							promqlbuilder.Lss(
								promqlbuilder.AvgOverTime(
									matrix.New(
										vector.New(
											vector.WithMetricName("up"),
											vector.WithLabelMatchers(
												label.New("job").Equal(a.AlertmanagerServiceSelector),
											),
										),
										matrix.WithRange(5*time.Minute),
									),
								),
								promqlbuilder.NewNumber(0.5),
							),
						).By("job"),
						promqlbuilder.Count(
							vector.New(
								vector.WithMetricName("up"),
								vector.WithLabelMatchers(
									label.New("job").Equal(a.AlertmanagerServiceSelector),
								),
							),
						).By("job"),
					),
					promqlbuilder.NewNumber(0.5),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  a.ServiceLabelValue,
						"severity": "critical",
					},
					a.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						a.DashboardURL,
						a.RunbookURL,
						runbookAlertmanagerClusterDown,
						"{{ $value | humanizePercentage }} of Alertmanager instances within the {{$labels.job}} cluster have been up for less than half of the last 5m.",
						"{{ $value | humanizePercentage }} of Alertmanager instances within the {{$labels.job}} cluster have been up for less than half of the last 5m.",
						"Half or more of the Alertmanager instances within the same cluster are down.",
					),
					a.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule[alerting.Option](
			"AlertmanagerClusterCrashlooping",
			alerting.Expr(
				promqlbuilder.Gte(
					promqlbuilder.Div(
						promqlbuilder.Count(
							promqlbuilder.Gtr(
								promqlbuilder.Changes(
									matrix.New(
										vector.New(
											vector.WithMetricName("process_start_time_seconds"),
											vector.WithLabelMatchers(
												label.New("job").Equal(a.AlertmanagerServiceSelector),
											),
										),
										matrix.WithRange(10*time.Minute),
									),
								),
								promqlbuilder.NewNumber(4),
							),
						).By("job"),
						promqlbuilder.Count(
							vector.New(
								vector.WithMetricName("up"),
								vector.WithLabelMatchers(
									label.New("job").Equal(a.AlertmanagerServiceSelector),
								),
							),
						).By("job"),
					),
					promqlbuilder.NewNumber(0.5),
				),
			),
			alerting.For("5m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"service":  a.ServiceLabelValue,
						"severity": "critical",
					},
					a.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					common.BuildAnnotations(
						a.DashboardURL,
						a.RunbookURL,
						runbookAlertmanagerClusterCrashlooping,
						"{{ $value | humanizePercentage }} of Alertmanager instances within the {{$labels.job}} cluster have restarted at least 5 times in the last 10m.",
						"{{ $value | humanizePercentage }} of Alertmanager instances within the {{$labels.job}} cluster have restarted at least 5 times in the last 10m.",
						"Half or more of the Alertmanager instances within the same cluster are crashlooping.",
					),
					a.AdditionalAlertAnnotations,
				),
			),
		),
	}
}
