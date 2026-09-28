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
	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
	"github.com/prometheus/prometheus/model/labels"
)

func timeSeriesChart(unit *string, stacked bool) panel.Option {
	visual := timeSeriesPanel.Visual{
		Display:      timeSeriesPanel.LineDisplay,
		ConnectNulls: false,
		LineWidth:    0.25,
		AreaOpacity:  0.5,
		Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
	}
	if stacked {
		visual.Stack = timeSeriesPanel.AllStack
		visual.AreaOpacity = 1
	}

	return timeSeriesPanel.Chart(
		timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
			Format: &commonSdk.Format{Unit: unit},
		}),
		timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
			Position: timeSeriesPanel.BottomPosition,
			Mode:     timeSeriesPanel.TableMode,
			Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
		}),
		timeSeriesPanel.WithVisual(visual),
	)
}

func promQLQuery(datasourceName string, queryName string, seriesName string, labelMatchers ...*labels.Matcher) panel.Option {
	return panel.AddQuery(
		query.PromQL(
			promql.SetLabelMatchersV2(
				AISIXCommonPanelQueries[queryName],
				labelMatchers,
			).Pretty(0),
			dashboards.AddQueryDataSource(datasourceName),
			query.SeriesNameFormat(seriesName),
		),
	)
}

func RequestRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Request rate",
		panel.Description("Rate of completed AISIX requests, grouped by provider, model, and outcome."),
		timeSeriesChart(&dashboards.RequestsPerSecondsUnit, true),
		promQLQuery(datasourceName, "RequestRate", "{{provider}} / {{model}} / {{outcome}}", labelMatchers...),
	)
}

func SuccessRatio(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Success ratio",
		panel.Description("Percentage of completed AISIX requests whose outcome is success."),
		timeSeriesChart(&dashboards.PercentDecimalUnit, false),
		promQLQuery(datasourceName, "SuccessRatio", "Success", labelMatchers...),
	)
}

func EndToEndLatencyP95(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("P95 end-to-end latency",
		panel.Description("P95 client-perceived request latency. Streaming requests are observed at stream completion."),
		timeSeriesChart(&dashboards.SecondsUnit, false),
		promQLQuery(datasourceName, "EndToEndLatencyP95", "P95 end-to-end", labelMatchers...),
	)
}

func TimeToFirstTokenP95(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("P95 time to first token",
		panel.Description("P95 time to first token for streaming requests. Non-streaming requests do not contribute."),
		timeSeriesChart(&dashboards.SecondsUnit, false),
		promQLQuery(datasourceName, "TimeToFirstTokenP95", "P95 TTFT", labelMatchers...),
	)
}

func TokenThroughput(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Token throughput by model",
		panel.Description("Input and output tokens per second reported by the upstream provider."),
		timeSeriesChart(&dashboards.CountsPerSecondsUnit, false),
		promQLQuery(datasourceName, "TokenThroughput_input", "{{model}} input", labelMatchers...),
		promQLQuery(datasourceName, "TokenThroughput_output", "{{model}} output", labelMatchers...),
	)
}

func CacheOutcomes(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Cache outcomes",
		panel.Description("Rate of cache-policy outcomes. Requests without an active cache policy do not contribute."),
		timeSeriesChart(&dashboards.RequestsPerSecondsUnit, true),
		promQLQuery(datasourceName, "CacheOutcomes", "{{outcome}}", labelMatchers...),
	)
}

func PolicyEnforcement(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Policy enforcement",
		panel.Description("Rates of rate-limit rejections, guardrail blocks, and fail-open guardrail bypass events."),
		timeSeriesChart(&dashboards.CountsPerSecondsUnit, false),
		promQLQuery(datasourceName, "PolicyEnforcement_rateLimit", "Rate limit / {{scope}}", labelMatchers...),
		promQLQuery(datasourceName, "PolicyEnforcement_guardrailBlocks", "Guardrail blocks", labelMatchers...),
		promQLQuery(datasourceName, "PolicyEnforcement_guardrailBypasses", "Guardrail bypass / {{reason}}", labelMatchers...),
	)
}

func OTLPDelivery(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("OTLP delivery failures and drops",
		panel.Description("Failure and drop event rates from AISIX OTLP HTTP exporters. A zero value is not proof that tracing is enabled."),
		timeSeriesChart(&dashboards.CountsPerSecondsUnit, false),
		promQLQuery(datasourceName, "OTLPDelivery_failures", "{{exporter}} failures", labelMatchers...),
		promQLQuery(datasourceName, "OTLPDelivery_drops", "{{exporter}} drops / {{reason}}", labelMatchers...),
	)
}
