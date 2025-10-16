package blackbox

import (
	"time"

	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/matrix"
	"github.com/perses/promql-builder/vector"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	rulehelpers "github.com/perses/community-mixins/pkg/rules"
)

type BlackboxRulesConfig struct {
	AdditionalAlertLabels      map[string]string
	AdditionalAlertAnnotations map[string]string
	DashboardURL               string
}

func NewBlackboxRulesConfig(
	additionalAlertLabels map[string]string,
	additionalAlertAnnotations map[string]string,
	dashboardURL string,
) BlackboxRulesConfig {
	return BlackboxRulesConfig{
		AdditionalAlertLabels:      additionalAlertLabels,
		AdditionalAlertAnnotations: additionalAlertAnnotations,
		DashboardURL:               dashboardURL,
	}
}

// BuildBlackboxRules builds blackbox exporter rules
func BuildBlackboxRules(
	namespace string,
	blackboxRulesConfig BlackboxRulesConfig,
	labels map[string]string,
	annotations map[string]string,
) rulehelpers.RuleResult {

	groups := []monitoringv1.RuleGroup{
		blackboxRulesConfig.BlackboxExporterGroup(),
	}

	return rulehelpers.NewRuleResult(
		rulehelpers.NewPrometheusRule(
			"blackbox-exporter-rules",
			namespace,
			labels,
			annotations,
			groups,
		),
		nil,
	).Component("blackbox-exporter")
}

func (b BlackboxRulesConfig) BlackboxExporterGroup() monitoringv1.RuleGroup {
	rules := []monitoringv1.Rule{}
	rules = append(rules, rulehelpers.NewAlertingRule(
		"BlackboxProbeFailed",
		promqlbuilder.Eqlc(
			vector.New(
				vector.WithMetricName("probe_success"),
				vector.WithLabelMatchers(
					label.New("job").Equal("blackbox-exporter"),
				),
			),
			promqlbuilder.NewNumber(0),
		),
		"1m",
		rulehelpers.MergeMaps(
			map[string]string{
				"severity": "critical",
			},
			b.AdditionalAlertLabels,
		),
		rulehelpers.MergeMaps(
			map[string]string{
				"summary":       "Probe has failed for the past 1m interval.",
				"description":   "The probe failed for the instance {{ $labels.instance }}.",
				"dashboard_url": b.DashboardURL + "?var-instance={{ $labels.instance }}",
			},
			b.AdditionalAlertAnnotations,
		),
	))
	rules = append(rules, rulehelpers.NewAlertingRule(
		"BlackboxLowUptime30d",
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
		"",
		rulehelpers.MergeMaps(
			map[string]string{
				"severity": "info",
			},
			b.AdditionalAlertLabels,
		),
		rulehelpers.MergeMaps(
			map[string]string{
				"summary":       "Probe uptime is lower than 99.9% for the last 30 days.",
				"description":   "The probe has a lower uptime than 99.9% the last 30 days for the instance {{ $labels.instance }}.",
				"dashboard_url": b.DashboardURL + "?var-instance={{ $labels.instance }}",
			},
			b.AdditionalAlertAnnotations,
		),
	))
	rules = append(rules, rulehelpers.NewAlertingRule(
		"BlackboxSslCertificateWillExpireSoon",
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
		"",
		rulehelpers.MergeMaps(
			map[string]string{
				"severity": "warning",
			},
			b.AdditionalAlertLabels,
		),
		rulehelpers.MergeMaps(
			map[string]string{
				"summary":       "SSL certificate will expire soon.",
				"description":   "The SSL certificate of the instance {{ $labels.instance }} is expiring within 21 days.\nActual time left: {{ $value | humanizeDuration }}.",
				"dashboard_url": b.DashboardURL + "?var-instance={{ $labels.instance }}",
			},
			b.AdditionalAlertAnnotations,
		),
	))

	return rulehelpers.NewRuleGroup("blackbox-exporter.rules", "", nil, rules)
}
