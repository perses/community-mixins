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

package promql

import (
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

func SetLabelMatchersV2(query parser.Expr, matchers []*labels.Matcher) parser.Expr {
	copy := promqlbuilder.DeepCopyExpr(query)
	for _, l := range matchers {
		copy = LabelsSetPromQLV2(copy, l.Type, l.Name, l.Value)
	}
	return copy
}

func LabelsSetPromQLV2(query parser.Expr, matchType labels.MatchType, name, value string) parser.Expr {
	if name == "" || value == "" {
		return query
	}

	promqlbuilder.Inspect(query, func(node parser.Node, path []parser.Node) error {
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
