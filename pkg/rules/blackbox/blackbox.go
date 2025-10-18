package blackbox

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

type BlackboxRulesConfig struct {
	AdditionalAlertLabels      map[string]string
	AdditionalAlertAnnotations map[string]string
	DashboardURL               string
}

type BlackboxRulesConfigOption func(*BlackboxRulesConfig)

func WithAdditionalAlertLabels(additionalAlertLabels map[string]string) BlackboxRulesConfigOption {
	return func(blackboxRulesConfig *BlackboxRulesConfig) {
		blackboxRulesConfig.AdditionalAlertLabels = additionalAlertLabels
	}
}

func WithAdditionalAlertAnnotations(additionalAlertAnnotations map[string]string) BlackboxRulesConfigOption {
	return func(blackboxRulesConfig *BlackboxRulesConfig) {
		blackboxRulesConfig.AdditionalAlertAnnotations = additionalAlertAnnotations
	}
}

func WithDashboardURL(dashboardURL string) BlackboxRulesConfigOption {
	return func(blackboxRulesConfig *BlackboxRulesConfig) {
		blackboxRulesConfig.DashboardURL = dashboardURL
	}
}

// BuildBlackboxRules builds blackbox exporter rules
func BuildBlackboxRules(
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	options ...BlackboxRulesConfigOption,
) rulehelpers.RuleResult {
	blackboxRulesConfig := BlackboxRulesConfig{}
	for _, option := range options {
		option(&blackboxRulesConfig)
	}

	promRule, err := promtheusrule.New(
		"blackbox-exporter-rules",
		namespace,
		promtheusrule.Labels(labels),
		promtheusrule.Annotations(annotations),
		promtheusrule.AddRuleGroup(
			"blackbox-exporter.rules",
			blackboxRulesConfig.BlackboxExporterRuleGroupOptions()...,
		),
	)

	if err != nil {
		return rulehelpers.NewRuleResult(nil, err).Component("blackbox-exporter")
	}

	return rulehelpers.NewRuleResult(
		&promRule.PrometheusRule,
		nil,
	).Component("blackbox-exporter")
}

func (b BlackboxRulesConfig) BlackboxExporterRuleGroupOptions() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.AddRule[alerting.Option](
			"BlackboxProbeFailed",
			alerting.Expr(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("probe_success"),
						vector.WithLabelMatchers(
							label.New("job").Equal("blackbox-exporter"),
						),
					),
					promqlbuilder.NewNumber(0),
				),
			),
			alerting.For("1m"),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"severity": "critical",
					},
					b.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					map[string]string{
						"summary":       "Probe has failed for the past 1m interval.",
						"description":   "The probe failed for the instance {{ $labels.instance }}.",
						"dashboard_url": b.DashboardURL + "?var-instance={{ $labels.instance }}",
					},
					b.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule[alerting.Option](
			"BlackboxLowUptime30d",
			alerting.Expr(
				promqlbuilder.Lss(
					promqlbuilder.Mul(
						promqlbuilder.AvgOverTime(
							matrix.New(
								vector.New(
									vector.WithMetricName("probe_success"),
									vector.WithLabelMatchers(
										label.New("job").Equal("blackbox-exporter"),
									),
								),
								matrix.WithRange(30*24*time.Hour),
							),
						),
						promqlbuilder.NewNumber(100),
					),
					promqlbuilder.NewNumber(99.9),
				),
			),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"severity": "info",
					},
					b.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					map[string]string{
						"summary":       "Probe uptime is lower than 99.9% for the last 30 days.",
						"description":   "The probe has a lower uptime than 99.9% the last 30 days for the instance {{ $labels.instance }}.",
						"dashboard_url": b.DashboardURL + "?var-instance={{ $labels.instance }}",
					},
					b.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule[alerting.Option](
			"BlackboxSslCertificateWillExpireSoon",
			alerting.Expr(
				promqlbuilder.Lss(
					promqlbuilder.Sub(
						vector.New(
							vector.WithMetricName("probe_ssl_earliest_cert_expiry"),
							vector.WithLabelMatchers(
								label.New("job").Equal("blackbox-exporter"),
							),
						),
						promqlbuilder.Time(),
					),
					promqlbuilder.Mul(
						promqlbuilder.Mul(
							promqlbuilder.NewNumber(21),
							promqlbuilder.NewNumber(24),
						),
						promqlbuilder.NewNumber(3600),
					),
				),
			),
			alerting.Labels(
				common.MergeMaps(
					map[string]string{
						"severity": "warning",
					},
					b.AdditionalAlertLabels,
				),
			),
			alerting.Annotations(
				common.MergeMaps(
					map[string]string{
						"summary":       "SSL certificate will expire soon.",
						"description":   "The SSL certificate of the instance {{ $labels.instance }} is expiring within 21 days.\nActual time left: {{ $value | humanizeDuration }}.",
						"dashboard_url": b.DashboardURL + "?var-instance={{ $labels.instance }}",
					},
					b.AdditionalAlertAnnotations,
				),
			),
		),
	}
}
