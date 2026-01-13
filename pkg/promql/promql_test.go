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
	"testing"
)

func TestSetLabelMatchers(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		labelMatchers []LabelMatcher
		want          string
	}{
		{
			name:  "simple metric query",
			query: "metric{job=\"test\"}",
			labelMatchers: []LabelMatcher{
				{Name: "instance", Value: "localhost:9090", Type: "="},
			},
			want: "metric{instance=\"localhost:9090\",job=\"test\"}",
		},
		{
			name:  "query with rate and interval variable",
			query: "rate(metric{job=\"test\"}[$__rate_interval])",
			labelMatchers: []LabelMatcher{
				{Name: "instance", Value: "localhost:9090", Type: "="},
			},
			want: "rate(metric{instance=\"localhost:9090\",job=\"test\"}[$__rate_interval])",
		},
		{
			name:  "query with multiple variables",
			query: "metric{job=\"test\"}[$__rate_interval] offset $__range",
			labelMatchers: []LabelMatcher{
				{Name: "instance", Value: "localhost:9090", Type: "="},
			},
			want: "metric{instance=\"localhost:9090\",job=\"test\"}[$__rate_interval] offset $__range",
		},
		{
			name:  "query with regex match",
			query: "metric{job=~\"test.*\"}",
			labelMatchers: []LabelMatcher{
				{Name: "instance", Value: "localhost.*", Type: "=~"},
			},
			want: "metric{instance=~\"localhost.*\",job=~\"test.*\"}",
		},
		{
			name:  "query with negative match",
			query: "metric{job!=\"test\"}",
			labelMatchers: []LabelMatcher{
				{Name: "instance", Value: "localhost:9090", Type: "!="},
			},
			want: "metric{instance!=\"localhost:9090\",job!=\"test\"}",
		},
		{
			name:  "query with multiple label matchers",
			query: "metric{job=\"test\"}",
			labelMatchers: []LabelMatcher{
				{Name: "instance", Value: "localhost:9090", Type: "="},
				{Name: "env", Value: "prod", Type: "="},
			},
			want: "metric{env=\"prod\",instance=\"localhost:9090\",job=\"test\"}",
		},
		{
			name:  "query with complex expression and variables",
			query: "sum(rate(metric{job=\"test\"}[$__rate_interval])) by (instance) / 3",
			labelMatchers: []LabelMatcher{
				{Name: "env", Value: "prod", Type: "="},
			},
			want: "sum by (instance) (rate(metric{env=\"prod\",job=\"test\"}[$__rate_interval])) / 3",
		},
		{
			name:  "query with all Perses variables",
			query: "metric{job=\"test\"}[$__interval] offset $__range",
			labelMatchers: []LabelMatcher{
				{Name: "instance", Value: "localhost:9090", Type: "="},
			},
			want: "metric{instance=\"localhost:9090\",job=\"test\"}[$__interval] offset $__range",
		},
		{
			name:  "query with dashboard and project variables",
			query: "metric{job=\"test\",dashboard=\"$__dashboard\",project=\"$__project\"}",
			labelMatchers: []LabelMatcher{
				{Name: "instance", Value: "localhost:9090", Type: "="},
			},
			want: "metric{dashboard=\"$__dashboard\",instance=\"localhost:9090\",job=\"test\",project=\"$__project\"}",
		},
		{
			name:  "query with time range variables",
			query: "metric{job=\"test\"}[$__interval] @ $__from",
			labelMatchers: []LabelMatcher{
				{Name: "instance", Value: "localhost:9090", Type: "="},
			},
			want: "metric{instance=\"localhost:9090\",job=\"test\"}[$__interval] @ $__from",
		},
		{
			name:  "query with ms variables",
			query: "metric{job=\"test\"}[$__interval_ms]",
			labelMatchers: []LabelMatcher{
				{Name: "instance", Value: "localhost:9090", Type: "="},
			},
			want: "metric{instance=\"localhost:9090\",job=\"test\"}[$__interval_ms]",
		},
		{
			name:  "query with range variables",
			query: "metric{job=\"test\"}[$__range_s]",
			labelMatchers: []LabelMatcher{
				{Name: "instance", Value: "localhost:9090", Type: "="},
			},
			want: "metric{instance=\"localhost:9090\",job=\"test\"}[$__range_s]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetLabelMatchers(tt.query, tt.labelMatchers)
			if got != tt.want {
				t.Errorf("SetLabelMatchers() = %v, want %v", got, tt.want)
			}
		})
	}
}
