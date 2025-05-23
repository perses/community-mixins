package promql

import (
	"fmt"
	"sort"
	"strings"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

// Use SetLabelMatchersV2 instead, which uses perses/promql-builder with inbuilt support for variables.
type PersesVarProcessor struct {
	Replacements map[string]string
	SortedKeys   []string
}

// Use SetLabelMatchersV2 instead, which uses perses/promql-builder with inbuilt support for variables.
// NewPersesVarProcessor creates a new PersesVarProcessor which will help replace/restore Perses variables in a query
// to make it compatible with PromQL parser.
// Sort of a hack, by modifying the query in place with random values
func NewPersesVarProcessor() *PersesVarProcessor {
	// Keep the variables in this map super unique, so that we can replace them safely.
	replacements := map[string]string{
		"$__rate_interval": "2d20h8m7s",
		"$__interval":      "2d20h8m8s",
		"$__interval_ms":   "7d19h59m27s",
		"$__dashboard":     "CHEESECAKE",
		"$__project":       "CHEESECAKE-DEV",
		"$__from":          "1715222400000.000",
		"$__to":            "1715222400000.000",
		"$__range":         "2d20h8m9s",
		"$__range_s":       "1h2m17s",
		"$__range_ms":      "3737373",
	}

	var sortedKeys []string
	for key := range replacements {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		return len(sortedKeys[i]) > len(sortedKeys[j])
	})

	return &PersesVarProcessor{
		Replacements: replacements,
		SortedKeys:   sortedKeys,
	}
}

func (p *PersesVarProcessor) Replace(query string) (string, map[string]string) {
	original := make(map[string]string)
	for _, name := range p.SortedKeys {
		if strings.Contains(query, name) {
			placeholder := p.Replacements[name]
			original[name] = placeholder
			query = strings.ReplaceAll(query, name, placeholder)
		}
	}
	return query, original
}

func (p *PersesVarProcessor) Restore(query string, original map[string]string) string {
	// Build a list of varName/placeholder pairs
	type pair struct {
		varName     string
		placeholder string
	}
	var pairs []pair
	for varName, placeholder := range original {
		pairs = append(pairs, pair{varName, placeholder})
	}

	// Sort by placeholder length descending
	sort.Slice(pairs, func(i, j int) bool {
		return len(pairs[i].placeholder) > len(pairs[j].placeholder)
	})

	// Replace each placeholder with the original var
	for _, p := range pairs {
		query = strings.ReplaceAll(query, p.placeholder, p.varName)
	}

	return query
}

// Use SetLabelMatchersV2 with prometheus []*labels.Matcher instead.
type LabelMatcher struct {
	Name  string
	Value string
	Type  string
}

// Use SetLabelMatchersV2 instead.
func SetLabelMatchers(query string, labelMatchers []LabelMatcher) string {
	processor := NewPersesVarProcessor()

	for _, l := range labelMatchers {
		query = LabelsSetPromQL(query, l.Type, l.Name, l.Value, processor)
	}
	return query
}

// Use LabelsSetPromQLV2 instead.
func LabelsSetPromQL(query, labelMatchType, name, value string, processor *PersesVarProcessor) string {
	modifiedQuery, originalVars := processor.Replace(query)
	expr, err := parser.ParseExpr(modifiedQuery)
	if err != nil {
		fmt.Println("Error parsing query:", err, modifiedQuery)
		return ""
	}

	if name == "" || value == "" {
		// Get the modified query and restore Perses variables
		result := expr.Pretty(0)
		return processor.Restore(result, originalVars)
	}

	var matchType labels.MatchType
	switch labelMatchType {
	case parser.ItemType(parser.EQL).String():
		matchType = labels.MatchEqual
	case parser.ItemType(parser.NEQ).String():
		matchType = labels.MatchNotEqual
	case parser.ItemType(parser.EQL_REGEX).String():
		matchType = labels.MatchRegexp
	case parser.ItemType(parser.NEQ_REGEX).String():
		matchType = labels.MatchNotRegexp
	default:
		fmt.Println("Invalid match type:", labelMatchType)
		return ""
	}

	parser.Inspect(expr, func(node parser.Node, path []parser.Node) error {
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

	// Get the modified query and restore Perses variables
	result := expr.Pretty(0)
	return processor.Restore(result, originalVars)
}
