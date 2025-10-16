package rules

import (
	"github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/prometheus/prometheus/promql/parser"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	promRuleTypeMeta = metav1.TypeMeta{
		APIVersion: monitoring.GroupName + "/" + monitoringv1.Version,
		Kind:       monitoringv1.PrometheusRuleKind,
	}
)

// NewPrometheusRule creates a new PrometheusRule object
func NewPrometheusRule(
	name, namespace string,
	labels map[string]string,
	annotations map[string]string,
	groups []monitoringv1.RuleGroup) *monitoringv1.PrometheusRule {

	return &monitoringv1.PrometheusRule{
		TypeMeta: promRuleTypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: monitoringv1.PrometheusRuleSpec{
			Groups: groups,
		},
	}
}

// NewRuleGroup creates a new RuleGroup object
func NewRuleGroup(name, interval string, labels map[string]string, rules []monitoringv1.Rule) monitoringv1.RuleGroup {
	intervalDuration := monitoringv1.Duration(interval)
	return monitoringv1.RuleGroup{
		Name:     name,
		Interval: &intervalDuration,
		Rules:    rules,
	}
}

// NewAlertingRule creates a new Rule object
func NewAlertingRule(alertName string, expr parser.Expr, forTime string, labels map[string]string, annotations map[string]string) monitoringv1.Rule {
	duration := monitoringv1.Duration(forTime)
	return monitoringv1.Rule{
		Alert:       alertName,
		Expr:        intstr.FromString(expr.Pretty(0)),
		For:         &duration,
		Labels:      labels,
		Annotations: annotations,
	}
}

// NewRecordingRule creates a new Rule object
func NewRecordingRule(recordName string, expr parser.Expr, labels map[string]string, annotations map[string]string) monitoringv1.Rule {
	return monitoringv1.Rule{
		Record:      recordName,
		Expr:        intstr.FromString(expr.Pretty(0)),
		Labels:      labels,
		Annotations: annotations,
	}
}

// BuildAnnotations creates annotations map with conditional dashboard and runbook labels
func BuildAnnotations(dashboardURL, runbookURL, runbookFragment, description, message, summary string) map[string]string {
	annotations := map[string]string{
		"description": description,
		"message":     message,
		"summary":     summary,
	}

	if dashboardURL != "" {
		annotations["dashboard"] = dashboardURL
	}
	if runbookURL != "" {
		annotations["runbook"] = runbookURL + runbookFragment
	}

	return annotations
}

// MergeMaps merges two maps
func MergeMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
