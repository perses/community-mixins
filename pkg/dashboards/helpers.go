package dashboards

import (
	"github.com/perses/community-dashboards/pkg/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/dashboard"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/prometheus/prometheus/model/labels"
)

var (
	SecondsUnit            = string(commonSdk.SecondsUnit)
	DecimalUnit            = string(commonSdk.DecimalUnit)
	BytesUnit              = string(commonSdk.BytesUnit)
	BytesPerSecondsUnit    = string(commonSdk.BytesPerSecondsUnit)
	MilliSecondsUnit       = string(commonSdk.MilliSecondsUnit)
	PercentDecimalUnit     = string(commonSdk.PercentDecimalUnit)
	RequestsPerSecondsUnit = string(commonSdk.RequestsPerSecondsUnit)
	OpsPerSecondsUnit      = string(commonSdk.OpsPerSecondsUnit)
	PercentUnit            = string(commonSdk.PercentUnit)
	CountsPerSecondsUnit   = string(commonSdk.CountsPerSecondsUnit)
	PacketsPerSecondsUnit  = string(commonSdk.PacketsPerSecondsUnit)
	ReadsPerSecondsUnit    = string(commonSdk.ReadsPerSecondsUnit)
	PercentMode            = string(commonSdk.PercentMode)
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
