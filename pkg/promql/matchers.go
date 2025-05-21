package promql

import (
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

func SetLabelMatchersV2(query parser.Expr, matchers []*labels.Matcher) parser.Expr {
	for _, l := range matchers {
		query = LabelsSetPromQLV2(query, l.Type, l.Name, l.Value)
	}
	return query
}

func LabelsSetPromQLV2(query parser.Expr, matchType labels.MatchType, name, value string) parser.Expr {
	if name == "" || value == "" {
		return query
	}

	parser.Inspect(query, func(node parser.Node, path []parser.Node) error {
		if n, ok := node.(*parser.VectorSelector); ok {
			var found bool
			for i, l := range n.LabelMatchers {
				if l.Name == name {
					n.LabelMatchers[i].Type = matchType
					n.LabelMatchers[i].Value = value
					found = true
				}
			}
			if !found {
				n.LabelMatchers = append(n.LabelMatchers, &labels.Matcher{
					Type:  matchType,
					Name:  name,
					Value: value,
				})
			}
		}
		return nil
	})

	return query
}
