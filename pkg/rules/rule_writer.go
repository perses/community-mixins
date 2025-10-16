package rules

import (
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

type RuleWriter struct {
	ruleResults []RuleResult
	executor    Exec
}

type RuleResult struct {
	rule      *monitoringv1.PrometheusRule
	component string
	err       error
}

func NewRuleResult(rule *monitoringv1.PrometheusRule, err error) RuleResult {
	return RuleResult{
		rule: rule,
		err:  err,
	}
}

// Components sets the component field of the RuleResult.
// This component field is used by RuleWriter, as the subdirectory name for the rule.
func (d RuleResult) Component(component string) RuleResult {
	d.component = component
	return d
}

func NewRuleWriter() *RuleWriter {
	return &RuleWriter{
		executor: NewExec(),
	}
}

// Add adds a rule to the writer.
func (w *RuleWriter) Add(dr RuleResult) {
	w.ruleResults = append(w.ruleResults, dr)
}

// Write writes the rules to the output directory.
func (w *RuleWriter) Write() {
	for _, result := range w.ruleResults {
		w.executor.BuildRule(result)
	}
}
