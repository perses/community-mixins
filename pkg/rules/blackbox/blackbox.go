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

// Runbook fragments
const (
	runbookBlackboxProbeFailed                  = "#blackboxprobefailed"
	runbookBlackboxLowUptime30d                 = "#blackboxlowuptime30d"
	runbookBlackboxSslCertificateWillExpireSoon = "#blackboxsslcertificatewillexpiresoon"
)

type BlackboxRulesConfig struct {
	DashboardURL string
	RunbookURL   string

	BlackboxExporterServiceSelector string

	AdditionalAlertLabels      map[string]string
	AdditionalAlertAnnotations map[string]string
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

func WithRunbookURL(runbookURL string) BlackboxRulesConfigOption {
	return func(blackboxRulesConfig *BlackboxRulesConfig) {
		blackboxRulesConfig.RunbookURL = runbookURL
	}
}

func WithBlackboxExporterServiceSelector(blackboxExporterServiceSelector string) BlackboxRulesConfigOption {
	return func(blackboxRulesConfig *BlackboxRulesConfig) {
		if blackboxExporterServiceSelector == "" {
			blackboxExporterServiceSelector = "blackbox-exporter"
		}
		blackboxRulesConfig.BlackboxExporterServiceSelector = blackboxExporterServiceSelector
	}
}

// NewBlackboxRulesBuilder creates a new Blackbox rules builder.
func NewBlackboxRulesBuilder(
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	options ...BlackboxRulesConfigOption,
) (promtheusrule.Builder, error) {
	blackboxRulesConfig := BlackboxRulesConfig{
		BlackboxExporterServiceSelector: "blackbox-exporter",
	}
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

	return promRule, err
}

// BuildBlackboxRules builds blackbox exporter rules
func BuildBlackboxRules(
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	options ...BlackboxRulesConfigOption,
) rulehelpers.RuleResult {
	promRule, err := NewBlackboxRulesBuilder(namespace, labels, annotations, options...)
	if err != nil {
		return rulehelpers.NewRuleResult(nil, err).Component("blackbox-exporter")
	}

	return rulehelpers.NewRuleResult(
		&promRule.PrometheusRule,
		nil,
	).Component("blackbox-exporter")
}

// BuildBlackboxRulesDefault builds the Blackbox Exporter rules with default configuration.
func BuildBlackboxRulesDefault(project string) rulehelpers.RuleResult {
	labels := map[string]string{
		"app.kubernetes.io/component": "blackbox-exporter",
		"app.kubernetes.io/name":      "blackbox-exporter-rules",
		"app.kubernetes.io/part-of":   "blackbox-exporter",
		"app.kubernetes.io/version":   "main",
	}
	annotations := map[string]string{}
	options := []BlackboxRulesConfigOption{
		WithDashboardURL("https://demo.perses.dev/projects/perses/dashboards/blackboxexporter"),
	}
	return BuildBlackboxRules(project, labels, annotations, options...)
}

func (b BlackboxRulesConfig) BlackboxExporterRuleGroupOptions() []rulegroup.Option {
	return []rulegroup.Option{
		rulegroup.AddRule(
			"BlackboxProbeFailed",
			alerting.Expr(
				promqlbuilder.Eqlc(
					vector.New(
						vector.WithMetricName("probe_success"),
						vector.WithLabelMatchers(
							label.New("job").Equal(b.BlackboxExporterServiceSelector),
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
					common.BuildAnnotations(
						b.DashboardURL,
						b.RunbookURL,
						runbookBlackboxProbeFailed,
						"The probe failed for the instance {{ $labels.instance }}.",
						"Probe has failed for the past 1m interval.",
					),
					b.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"BlackboxLowUptime30d",
			alerting.Expr(
				promqlbuilder.Lss(
					promqlbuilder.Mul(
						promqlbuilder.AvgOverTime(
							matrix.New(
								vector.New(
									vector.WithMetricName("probe_success"),
									vector.WithLabelMatchers(
										label.New("job").Equal(b.BlackboxExporterServiceSelector),
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
					common.BuildAnnotations(
						b.DashboardURL,
						b.RunbookURL,
						runbookBlackboxLowUptime30d,
						"The probe has a lower uptime than 99.9% the last 30 days for the instance {{ $labels.instance }}.",
						"Probe uptime is lower than 99.9% for the last 30 days.",
					),
					b.AdditionalAlertAnnotations,
				),
			),
		),
		rulegroup.AddRule(
			"BlackboxSslCertificateWillExpireSoon",
			alerting.Expr(
				promqlbuilder.Lss(
					promqlbuilder.Sub(
						vector.New(
							vector.WithMetricName("probe_ssl_earliest_cert_expiry"),
							vector.WithLabelMatchers(
								label.New("job").Equal(b.BlackboxExporterServiceSelector),
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
					common.BuildAnnotations(
						b.DashboardURL,
						b.RunbookURL,
						runbookBlackboxSslCertificateWillExpireSoon,
						"The SSL certificate of the instance {{ $labels.instance }} is expiring within 21 days.\nActual time left: {{ $value | humanizeDuration }}.",
						"SSL certificate will expire soon.",
					),
					b.AdditionalAlertAnnotations,
				),
			),
		),
	}
}
