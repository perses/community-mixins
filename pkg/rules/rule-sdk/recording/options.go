package recording

import (
	"github.com/prometheus/prometheus/promql/parser"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func RecordName(recordName string) Option {
	return func(builder *Builder) error {
		builder.Record = recordName
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
