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

	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/matrix"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetLabelMatchersV2(t *testing.T) {
	t.Run("adds label matcher and validates", func(t *testing.T) {
		expr := vector.New(vector.WithMetricName("up"))
		got := SetLabelMatchersV2(expr, []*labels.Matcher{
			label.New("job").Equal("prometheus"),
		})
		assert.Equal(t, `up{job="prometheus"}`, got.String())
	})

	t.Run("preserves rate interval variable", func(t *testing.T) {
		expr := promqlbuilder.Rate(
			matrix.New(
				vector.New(vector.WithMetricName("http_requests_total")),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		)
		got := SetLabelMatchersV2(expr, []*labels.Matcher{
			label.New("namespace").Equal("$namespace"),
		})
		assert.Equal(t, `rate(http_requests_total{namespace="$namespace"}[$__rate_interval])`, got.String())
	})

	t.Run("overwrites existing matcher", func(t *testing.T) {
		expr := vector.New(
			vector.WithMetricName("up"),
			vector.WithLabelMatchers(label.New("job").Equal("old")),
		)
		got := SetLabelMatchersV2(expr, []*labels.Matcher{
			label.New("job").Equal("new"),
		})
		assert.Equal(t, `up{job="new"}`, got.String())
	})

	t.Run("invalid expression panics on validate", func(t *testing.T) {
		expr := &parser.BinaryExpr{
			Op:  parser.ADD,
			LHS: vector.New(vector.WithMetricName("foo")),
		}
		assert.Panics(t, func() {
			SetLabelMatchersV2(expr, nil)
		})
	})
}

func TestLabelsSetPromQLV2(t *testing.T) {
	expr := vector.New(vector.WithMetricName("up"))
	got := LabelsSetPromQLV2(expr, labels.MatchEqual, "job", "prometheus")
	require.Equal(t, `up{job="prometheus"}`, got.String())

	// empty name/value is a no-op
	unchanged := LabelsSetPromQLV2(expr, labels.MatchEqual, "", "x")
	assert.Equal(t, expr, unchanged)
}
