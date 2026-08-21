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

package kubernetes

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"

	commonSdk "github.com/perses/perses/go-sdk/common"
	tablePanel "github.com/perses/plugins/table/sdk/go"
	"github.com/prometheus/prometheus/model/labels"
)

func ClusterCPUUsageQuota(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Quota",
		panel.Description("Shows the CPU requests, limits, and usage of workloads by namespace in tabular format."),
		tablePanel.Table(
			tablePanel.Transform([]commonSdk.Transform{
				{
					Kind: commonSdk.MergeSeriesKind,
					Spec: commonSdk.MergeSeriesSpec{
						Disabled: false,
					},
				},
			}),
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "namespace",
					Header: "Namespace",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value #1",
					Header: "Pods",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.DecimalUnit,
					},
				},
				{
					Name:   "value #2",
					Header: "Workloads",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.DecimalUnit,
					},
				},
				{
					Name:   "value #3",
					Header: "CPU Usage",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.DecimalUnit,
						DecimalPlaces: 4,
					},
				},
				{
					Name:   "value #4",
					Header: "CPU Requests",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.DecimalUnit,
						DecimalPlaces: 4,
					},
				},
				{
					Name:   "value #5",
					Header: "CPU Requests %",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.PercentDecimalUnit,
						DecimalPlaces: 4,
					},
				},
				{
					Name:   "value #6",
					Header: "CPU Limits",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.DecimalUnit,
						DecimalPlaces: 4,
					},
				},
				{
					Name:   "value #7",
					Header: "CPU Limits %",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.PercentDecimalUnit,
						DecimalPlaces: 4,
					},
				},
				{
					Name: "timestamp",
					Hide: true,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCPUUsageQuotaPodOwn"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCPUUsageQuotaNSwWorkload"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCPUUsageQuotaNodeNSCPU"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCPUUsageQuotaNSKubePodContainer"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCPUUsageQuotaNSKubePodContainerDiv"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCPUUsageQuotaNSCPUKubePodResources"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCPUUsageQuotaNSCPUKubePodResourcesDiv"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ClusterMemoryUsageQuota(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Requests by Namespace",
		panel.Description("Shows the memory requests, limits, and usage of workloads by namespace in tabular format."),
		tablePanel.Table(
			tablePanel.Transform([]commonSdk.Transform{
				{
					Kind: commonSdk.MergeSeriesKind,
					Spec: commonSdk.MergeSeriesSpec{
						Disabled: false,
					},
				},
			}),
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "namespace",
					Header: "Namespace",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value #1",
					Header: "Pods",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.DecimalUnit,
					},
				},
				{
					Name:   "value #2",
					Header: "Workloads",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.DecimalUnit,
					},
				},
				{
					Name:   "value #3",
					Header: "Memory Usage",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.BytesUnit,
					},
				},
				{
					Name:   "value #4",
					Header: "Memory Requests",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.BytesUnit,
					},
				},
				{
					Name:   "value #5",
					Header: "Memory Requests %",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.PercentDecimalUnit,
						DecimalPlaces: 4,
					},
				},
				{
					Name:   "value #6",
					Header: "Memory Limits",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.BytesUnit,
					},
				},
				{
					Name:   "value #7",
					Header: "Memory Limits %",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.PercentDecimalUnit,
						DecimalPlaces: 4,
					},
				},
				{
					Name: "timestamp",
					Hide: true,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterMemoryUsageQuotaPodOwner"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterMemoryUsageQuotaNSWorkloadPodOwner"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterMemoryUsageQuotaContainerMem"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterMemoryUsageQuotaContainerResourceReqSum"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterMemoryUsageQuotaContainerResourceReqSumDiv"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterMemoryUsageQuotaContainerReqLimits"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterMemoryUsageQuotaContainerReqLimitsDiv"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ClusterCurrentNetworkUsage(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Current Network Usage",
		panel.Description("Shows the current network usage of the cluster by namespace."),
		tablePanel.Table(
			tablePanel.Transform([]commonSdk.Transform{
				{
					Kind: commonSdk.MergeSeriesKind,
					Spec: commonSdk.MergeSeriesSpec{
						Disabled: false,
					},
				},
			}),
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "namespace",
					Header: "Namespace",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value #1",
					Header: "Current Receive Bandwidth",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.BytesPerSecondsUnit,
					},
				},
				{
					Name:   "value #2",
					Header: "Current Transmit Bandwidth",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.BytesPerSecondsUnit,
					},
				},
				{
					Name:   "value #3",
					Header: "Rate of Received Packets",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.PacketsPerSecondsUnit,
					},
				},
				{
					Name:   "value #4",
					Header: "Rate of Transmitted Packets",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.PacketsPerSecondsUnit,
					},
				},
				{
					Name:   "value #5",
					Header: "Rate of Received Packets Dropped",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.PacketsPerSecondsUnit,
					},
				},
				{
					Name:   "value #6",
					Header: "Rate of Transmitted Packets Dropped",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.PacketsPerSecondsUnit,
					},
				},
				{
					Name: "timestamp",
					Hide: true,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCurrentNetworkUsageBytesTotal"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCurrentNetworkTransmitBytesTotal"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCurrentNetworkReceivedTotal"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCurrentNetworkTransmitPacketsTotal"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCurrentNetworkReceivedPacketsDroppedTotal"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCurrentNetworkTransmitPacketsDroppedTotal"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ClusterCurrentStorageIO(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Current Storage IO",
		panel.Description("Shows the current storage IO of the cluster in tabular form, by namespace."),
		tablePanel.Table(
			tablePanel.Transform([]commonSdk.Transform{
				{
					Kind: commonSdk.MergeSeriesKind,
					Spec: commonSdk.MergeSeriesSpec{
						Disabled: false,
					},
				},
			}),
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "namespace",
					Header: "Namespace",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value #1",
					Header: "IOPS(Reads)",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.OpsPerSecondsUnit,
					},
				},
				{
					Name:   "value #2",
					Header: "IOPS(Writes)",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.OpsPerSecondsUnit,
					},
				},
				{
					Name:   "value #3",
					Header: "IOPS(Reads + Writes)",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.OpsPerSecondsUnit,
					},
				},
				{
					Name:   "value #4",
					Header: "Throughput(Reads)",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.BytesPerSecondsUnit,
					},
				},
				{
					Name:   "value #5",
					Header: "Throughput(Writes)",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.BytesPerSecondsUnit,
					},
				},
				{
					Name:   "value #6",
					Header: "Throughput(Reads + Writes)",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.BytesPerSecondsUnit,
					},
				},
				{
					Name: "timestamp",
					Hide: true,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCurrentStorageIOFsReadsTotal"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCurrentStorageIOFsWritesTotal"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCurrentStorageIOFsReadsWritesTotal"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCurrentStorageIOFsReadsBytesTotal"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCurrentStorageIOFsWritesBytesTotal"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					KubernetesCommonPanelQueries["ClusterCurrentStorageIOFsReadsWritesBytesTotal"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}
