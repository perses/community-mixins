// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aisix

import (
	"maps"

	"github.com/perses/community-mixins/pkg/promql"
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/prometheus/prometheus/promql/parser"
)

var AISIXCommonPanelQueries = map[string]parser.Expr{
	"RequestRate": promql.SumByRate(
		"aisix_requests_total",
		[]string{"provider", "model", "outcome"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"SuccessRatio": promqlbuilder.Div(
		promql.SumRate(
			"aisix_requests_total",
			label.New("job").EqualRegexp("$job"),
			label.New("instance").EqualRegexp("$instance"),
			label.New("outcome").Equal("success"),
		),
		promql.SumRate(
			"aisix_requests_total",
			label.New("job").EqualRegexp("$job"),
			label.New("instance").EqualRegexp("$instance"),
		),
	),
	"EndToEndLatencyP95": promqlbuilder.HistogramQuantile(
		0.95,
		promql.SumByRate(
			"aisix_request_e2e_latency_seconds_bucket",
			[]string{"le"},
			label.New("job").EqualRegexp("$job"),
			label.New("instance").EqualRegexp("$instance"),
		),
	),
	"TimeToFirstTokenP95": promqlbuilder.HistogramQuantile(
		0.95,
		promql.SumByRate(
			"aisix_request_ttft_seconds_bucket",
			[]string{"le"},
			label.New("job").EqualRegexp("$job"),
			label.New("instance").EqualRegexp("$instance"),
		),
	),
	"TokenThroughput_input": promql.SumByRate(
		"aisix_llm_input_tokens_total",
		[]string{"model"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"TokenThroughput_output": promql.SumByRate(
		"aisix_llm_output_tokens_total",
		[]string{"model"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"CacheOutcomes": promql.SumByRate(
		"aisix_cache_requests_total",
		[]string{"outcome"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"PolicyEnforcement_rateLimit": promql.SumByRate(
		"aisix_ratelimit_rejections_total",
		[]string{"scope"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"PolicyEnforcement_guardrailBlocks": promql.SumRate(
		"aisix_guardrail_blocks_total",
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"PolicyEnforcement_guardrailBypasses": promql.SumByRate(
		"aisix_guardrail_bypasses_total",
		[]string{"reason"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"OTLPDelivery_failures": promql.SumByRate(
		"aisix_otlp_fanout_failures_total",
		[]string{"exporter"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"OTLPDelivery_drops": promql.SumByRate(
		"aisix_otlp_fanout_drops_total",
		[]string{"exporter", "reason"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
}

// OverrideAISIXPanelQueries overrides the AISIXCommonPanelQueries global.
// Refer to panel queries in the map that you would like to override.
func OverrideAISIXPanelQueries(queries map[string]parser.Expr) {
	maps.Copy(AISIXCommonPanelQueries, queries)
}
