package rulegroup

import (
	"fmt"

	"github.com/perses/community-mixins/pkg/rules/rule-sdk/alerting"
	"github.com/perses/community-mixins/pkg/rules/rule-sdk/recording"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

func Name(name string) Option {
	return func(builder *Builder) error {
		builder.Name = name
		return nil
	}
}

func Labels(labels map[string]string) Option {
	return func(builder *Builder) error {
		builder.Labels = labels
		return nil
	}
}

func Interval(interval string) Option {
	return func(builder *Builder) error {
		if interval != "" {
			duration := monitoringv1.Duration(interval)
			builder.Interval = &duration
		}
		return nil
	}
}

func AddRule[O recording.Option | alerting.Option](name string, options ...O) Option {
	return func(builder *Builder) error {
		var rule monitoringv1.Rule

		var o O
		switch (any(o)).(type) {
		case recording.Option:
			recOptions := make([]recording.Option, len(options))
			for i, opt := range options {
				recOptions[i] = any(opt).(recording.Option)
			}

			r, recErr := recording.New(name, recOptions...)
			if recErr != nil {
				return recErr
			}
			rule = r.Rule
		case alerting.Option:
			altOptions := make([]alerting.Option, len(options))
			for i, opt := range options {
				altOptions[i] = any(opt).(alerting.Option)
			}

			a, altErr := alerting.New(name, altOptions...)
			if altErr != nil {
				return altErr
			}
			rule = a.Rule

		default:
			return fmt.Errorf("unsupported option type: %T", o)
		}

		builder.Rules = append(builder.Rules, rule)
		return nil
	}
}
