package alerting

import (
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/prometheus/prometheus/promql/parser"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func AlertName(alertName string) Option {
	return func(builder *Builder) error {
		builder.Alert = alertName
		return nil
	}
}

func Expr(expr parser.Expr) Option {
	return func(builder *Builder) error {
		builder.Expr = intstr.FromString(expr.Pretty(0))
		return nil
	}
}

func Labels(labels map[string]string) Option {
	return func(builder *Builder) error {
		builder.Labels = labels
		return nil
	}
}

func Annotations(annotations map[string]string) Option {
	return func(builder *Builder) error {
		builder.Annotations = annotations
		return nil
	}
}

func For(forTime string) Option {
	return func(builder *Builder) error {
		if forTime != "" {
			duration := monitoringv1.Duration(forTime)
			builder.For = &duration
		}
		return nil
	}
}
