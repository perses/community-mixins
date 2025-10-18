package alerting

import monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

type Option func(recordingRule *Builder) error

type Builder struct {
	monitoringv1.Rule `json:",inline" yaml:",inline"`
}

func New(alertName string, options ...Option) (Builder, error) {
	builder := &Builder{
		Rule: monitoringv1.Rule{},
	}

	defaults := []Option{
		AlertName(alertName),
	}

	for _, opt := range append(defaults, options...) {
		if err := opt(builder); err != nil {
			return *builder, err
		}
	}

	return *builder, nil
}
