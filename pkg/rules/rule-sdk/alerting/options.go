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
