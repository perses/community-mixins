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

package dashboards

import (
	"github.com/perses/community-mixins/pkg/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/dashboard"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/prometheus/prometheus/model/labels"
)

var (
	BytesPerSecondsUnit    = string(commonSdk.BytesPerSecondsUnit)
	BytesUnit              = string(commonSdk.BinaryBytesUnit)
	CountsPerSecondsUnit   = string(commonSdk.CountsPerSecondsUnit)
	DecimalUnit            = string(commonSdk.DecimalUnit)
	MilliSecondsUnit       = string(commonSdk.MilliSecondsUnit)
	OpsPerSecondsUnit      = string(commonSdk.OpsPerSecondsUnit)
	PacketsPerSecondsUnit  = string(commonSdk.PacketsPerSecondsUnit)
	PercentDecimalUnit     = string(commonSdk.PercentDecimalUnit)
	PercentMode            = string(commonSdk.PercentMode)
	PercentUnit            = string(commonSdk.PercentUnit)
	ReadsPerSecondsUnit    = string(commonSdk.ReadsPerSecondsUnit)
	RequestsPerSecondsUnit = string(commonSdk.RequestsPerSecondsUnit)
	SecondsUnit            = string(commonSdk.SecondsUnit)
)

func AddVariableDatasource(datasourceName string) labelValuesVar.Option {
	if datasourceName == "" {
		return func(plugin *labelValuesVar.Builder) error {
			return nil
		}
	}
	return labelValuesVar.Datasource(datasourceName)
}

func AddQueryDataSource(datasourceName string) query.Option {
	if datasourceName == "" {
		return func(plugin *query.Builder) error {
			return nil
		}
	}
	return query.Datasource(datasourceName)
}

func AddClusterVariable(datasource, clusterLabelName, matcher string) dashboard.Option {
	if clusterLabelName == "" {
		return func(builder *dashboard.Builder) error {
			return nil
		}
	}
	return dashboard.AddVariable("cluster",
		listVar.List(
			labelValuesVar.PrometheusLabelValues(clusterLabelName,
				labelValuesVar.Matchers(matcher),
				AddVariableDatasource(datasource),
			),
			listVar.DisplayName(clusterLabelName),
		),
	)
}

func GetClusterLabelMatcher(clusterLabelName string) promql.LabelMatcher {
	return promql.LabelMatcher{
		Name:  clusterLabelName,
		Value: "$cluster",
		Type:  "=",
	}
}

func GetClusterLabelMatcherV2(clusterLabelName string) *labels.Matcher {
	return &labels.Matcher{
		Name:  clusterLabelName,
		Value: "$cluster",
		Type:  labels.MatchEqual,
	}
}
