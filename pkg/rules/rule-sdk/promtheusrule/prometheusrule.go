package promtheusrule

import (
	"github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Option func(prometheusRule *Builder) error

type Builder struct {
	monitoringv1.PrometheusRule `json:",inline" yaml:",inline"`
}

var promRuleTypeMeta = metav1.TypeMeta{
	APIVersion: monitoring.GroupName + "/" + monitoringv1.Version,
	Kind:       monitoringv1.PrometheusRuleKind,
}

func New(name, namespace string, options ...Option) (Builder, error) {
	builder := &Builder{
		PrometheusRule: monitoringv1.PrometheusRule{
			TypeMeta: promRuleTypeMeta,
		},
	}

	defaults := []Option{
		Name(name),
		Namespace(namespace),
	}

	for _, opt := range append(defaults, options...) {
		if err := opt(builder); err != nil {
			return *builder, err
		}
	}

	return *builder, nil
}
