package promql

import (
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/matrix"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

func SumByRate(metricName string, byLabels []string, labelMatchers ...*labels.Matcher) parser.Expr {
	return promqlbuilder.Sum(
		promqlbuilder.Rate(
			matrix.New(
				vector.New(
					vector.WithMetricName(metricName),
					vector.WithLabelMatchers(labelMatchers...),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
	).By(byLabels...)
}

func SumByIncrease(metricName string, byLabels []string, labelMatchers ...*labels.Matcher) parser.Expr {
	return promqlbuilder.Sum(
		promqlbuilder.Increase(
			matrix.New(
				vector.New(
					vector.WithMetricName(metricName),
					vector.WithLabelMatchers(labelMatchers...),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
	).By(byLabels...)
}

func SumBy(metricName string, byLabels []string, labelMatchers ...*labels.Matcher) parser.Expr {
	return promqlbuilder.Sum(
		matrix.New(
			vector.New(vector.WithMetricName(metricName), vector.WithLabelMatchers(labelMatchers...)),
		),
	).By(byLabels...)
}

func ErrorCaseRatio(
	numeratorMetricName string,
	numeratorByLabels []string,
	numeratorLabelMatchers []*labels.Matcher,
	denominatorMetricName string,
	denominatorByLabels []string,
	denominatorLabelMatchers []*labels.Matcher,
) *promqlbuilder.BinaryBuilder {
	return promqlbuilder.Div(
		promqlbuilder.Sum(
			promqlbuilder.Rate(
				matrix.New(
					vector.New(
						vector.WithMetricName(numeratorMetricName),
						vector.WithLabelMatchers(numeratorLabelMatchers...),
					),
					matrix.WithRangeAsVariable("$__rate_interval"),
				),
			),
		).By(numeratorByLabels...),
		promqlbuilder.Sum(
			promqlbuilder.Rate(
				matrix.New(
					vector.New(
						vector.WithMetricName(denominatorMetricName),
						vector.WithLabelMatchers(denominatorLabelMatchers...),
					),
					matrix.WithRangeAsVariable("$__rate_interval"),
				),
			),
		).By(denominatorByLabels...),
	)
}

func ErrorCasePercentage(
	numeratorMetricName string,
	numeratorByLabels []string,
	numeratorLabelMatchers []*labels.Matcher,
	denominatorMetricName string,
	denominatorByLabels []string,
	denominatorLabelMatchers []*labels.Matcher,
) parser.Expr {
	return promqlbuilder.Mul(
		ErrorCaseRatio(numeratorMetricName, numeratorByLabels, numeratorLabelMatchers, denominatorMetricName, denominatorByLabels, denominatorLabelMatchers),
		&parser.NumberLiteral{Val: 100},
	)
}

func IgnoringGroupLeft(binaryOp *promqlbuilder.BinaryBuilder, ignoringLabels []string, groupLeftLabels ...string) parser.Expr {
	return binaryOp.Ignoring(ignoringLabels...).GroupLeft(groupLeftLabels...)
}

func OnGroupLeft(binaryOp *promqlbuilder.BinaryBuilder, onLabels []string, groupLeftLabels ...string) parser.Expr {
	return binaryOp.On(onLabels...).GroupLeft(groupLeftLabels...)
}
