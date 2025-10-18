package rulegroup

import (
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

type Option func(recordingRule *Builder) error

type Builder struct {
	monitoringv1.RuleGroup `json:",inline" yaml:",inline"`
}

func New(name string, options ...Option) (Builder, error) {
	builder := &Builder{
		RuleGroup: monitoringv1.RuleGroup{},
	}

	defaults := []Option{
		Name(name),
	}

	for _, opt := range append(defaults, options...) {
		if err := opt(builder); err != nil {
			return *builder, err
		}
	}

	return *builder, nil
}
