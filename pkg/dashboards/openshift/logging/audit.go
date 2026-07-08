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

package logging

import (
	"fmt"

	"github.com/perses/community-mixins/pkg/dashboards"
	panels "github.com/perses/community-mixins/pkg/panels/openshift/logging"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listvariable "github.com/perses/perses/go-sdk/variable/list-variable"
	textvariable "github.com/perses/perses/go-sdk/variable/text-variable"
	markdown "github.com/perses/plugins/markdown/sdk/go"
	staticlist "github.com/perses/plugins/staticlistvariable/sdk/go"
)

const (
	// auditVarQueryFmt builds a metric expression returning unique values of a
	// field from audit logs. Uses count_over_time to get a matrix result (one
	// series per unique value) instead of a streams result, avoiding the default
	// 100-entry limit on log queries.
	//
	// Only suitable for low-cardinality fields (verb, resource type). High-cardinality
	// fields (username, namespace) use TextVariable because exceeding the 500
	// series limit causes Loki to return HTTP 400.
	auditVarQueryFmt = `count by (%[1]s) (count_over_time({log_type="audit", openshift_log_source="kubeAPI"} | json %[1]s [1h]))`

	// auditLogQuery is the LogQL query for OTLP-format audit logs.
	// openshift_log_source is a stream label (indexed), so it avoids a full scan.
	// Audit fields require | json until structured metadata extraction lands upstream.
	auditLogQuery = `{log_type="audit", openshift_log_source="kubeAPI"} | json | user_username!~"${exclude_sa}" | user_username=~"(?i).*(?:${username}).*" | verb=~"${verb}" | objectRef_resource=~"${resource}" | objectRef_resource!~"${exclude_resource}" | objectRef_namespace=~".*(?:${namespace}).*" | objectRef_name=~"(?i).*(?:${resource_name}).*" | responseStatus_code=~"${response_code}" | userAgent=~"(?i).*(?:${client}).*" | line_format "User={{.user_username}} | Verb={{.verb}} | Namespace={{.objectRef_namespace}} | Resource={{.objectRef_resource}} | Resource Name={{.objectRef_name}} | Status={{.responseStatus_code}} | Client={{.userAgent}}" ${filter}`

	excludeSACustomAllValue       = `system:serviceaccount:.*|system:node:.*|system:kube.*|system:openshift.*|system:apiserver.*|system:aggregator.*|system:open-cluster-management:.*|system:ovn-node:.*|system:authenticated.*|system:unauthenticated.*|system:monitoring.*|system:master.*|system:multus.*`
	excludeResourceCustomAllValue = `events|endpoints|endpointslices|leases|tokenreviews|subjectaccessreviews|selfsubjectaccessreviews|selfsubjectrulesreviews`

	auditHelpText = "**Requires:** OpenShift Logging with OTLP data model enabled.\n\n**Filters:** All text filters support regex. Leave empty = match all.\n- **Username:** e.g. `admin`, `.*@example.com`, `user1|user2`\n- **Resource Type:** e.g. `pods`, `deploy.*`, `configmaps|secrets`\n- **Resource Name:** e.g. `my-pod.*`, `nginx`, `etcd.*`\n- **Namespace:** e.g. `openshift-.*`, `my-app`, `kube-system|default`\n- **Client:** e.g. `oc`, `kubectl`, `console`, `argocd`\n\n**LogQL Filter:** Raw stage, e.g. `|~ \"error\"` (include), `!~ \"health\"` (exclude), `| user_username!~\"bot.*\"`\n\n**Tip:** Use shorter time ranges for faster queries."

	auditQueryDisplay = "Active Query: `" + auditLogQuery + "`"
)

// lokiDatasourceRef represents a reference to a Loki datasource in variable specs.
// The Go SDK does not yet have builders for LokiLogQLVariable
// (only the JS/TS UI supports them, added in perses/plugins#651).
type lokiDatasourceRef struct {
	Kind string `json:"kind"`
	Name string `json:"name,omitempty"`
}

type lokiLogQLVarSpec struct {
	Datasource *lokiDatasourceRef `json:"datasource,omitempty"`
	Expr       string             `json:"expr"`
	LabelName  string             `json:"labelName"`
}

// lokiLogQLVariable returns a listvariable.Option that populates the variable
// dropdown by running a LogQL expression and extracting unique values of labelName.
func lokiLogQLVariable(datasourceName, expr, labelName string) listvariable.Option {
	return func(builder *listvariable.Builder) error {
		builder.ListVariableSpec.Plugin.Kind = "LokiLogQLVariable"
		builder.ListVariableSpec.Plugin.Spec = lokiLogQLVarSpec{
			Datasource: &lokiDatasourceRef{Kind: "LokiDatasource", Name: datasourceName},
			Expr:       expr,
			LabelName:  labelName,
		}
		return nil
	}
}

// labeledValue pairs a query-time value with a human-readable display label.
// The Go SDK's staticlist.Values only accepts plain strings; this bypasses it
// to use the {value, label} object form supported by the StaticListVariable schema.
type labeledValue struct {
	Value string `json:"value"`
	Label string `json:"label,omitempty"`
}

// staticListWithLabels returns a listvariable.Option that uses StaticListVariable
// with {value, label} pairs so the dropdown shows friendly names while the query
// uses the raw value.
func staticListWithLabels(values ...labeledValue) listvariable.Option {
	return func(builder *listvariable.Builder) error {
		builder.ListVariableSpec.Plugin.Kind = "StaticListVariable"
		builder.ListVariableSpec.Plugin.Spec = struct {
			Values []labeledValue `json:"values"`
		}{Values: values}
		return nil
	}
}

// BuildAuditLogViewer creates the OCP Audit Log Viewer dashboard.
//
// This dashboard enables cluster admins to investigate Kubernetes API audit logs
// stored in Loki via OpenShift Logging with OTLP data model.
//
// For deployment instructions, recommended Loki audit log filtering,
// performance tuning, and known limitations, see:
// https://github.com/sradco/ocp-audit-log-perses-dashboard
//
// Prerequisites:
//   - OpenShift Logging with LokiStack configured to collect audit logs
//   - OTLP data model with LokiStack schema v13 for structured metadata support
//   - A LokiDatasource configured in Perses pointing at the audit tenant
//
// Parameters:
//   - project: The Perses project name.
//   - lokiDatasource: The name of the Loki data source for audit logs.
func BuildAuditLogViewer(project string, lokiDatasource string) dashboards.DashboardResult {
	return dashboards.NewDashboardResult(
		dashboard.New("ocp-audit-log-viewer",
			dashboard.ProjectName(project),
			dashboard.Name("OCP Audit Log Viewer"),
			dashboard.DurationAsString("1h"),
			dashboard.RefreshIntervalAsString("0s"),

			dashboard.AddVariable("username",
				textvariable.Text("",
					textvariable.DisplayName("Username"),
					textvariable.Description("Regex filter. Examples: username, .*@redhat.com, admin|user"),
				),
			),
			dashboard.AddVariable("exclude_sa",
				listvariable.List(
					staticlist.StaticList(
						staticlist.Values(
							"^$",
							"system:serviceaccount:.*",
							"system:node:.*",
							"system:kube.*",
							"system:openshift.*",
							"system:apiserver.*",
							"system:aggregator.*",
							"system:open-cluster-management:.*",
							"system:ovn-node:.*",
							"system:authenticated.*",
							"system:unauthenticated.*",
							"system:monitoring.*",
							"system:master.*",
							"system:multus.*",
						),
					),
					listvariable.DisplayName("Exclude System Users"),
					listvariable.Description("Select None to show all users including system accounts"),
					listvariable.AllowAllValue(true),
					listvariable.AllowMultiple(true),
					listvariable.CustomAllValue(excludeSACustomAllValue),
					listvariable.DefaultValue("$__all"),
				),
			),
			dashboard.AddVariable("verb",
				listvariable.List(
					staticlist.StaticList(
						staticlist.Values("create", "update", "patch", "delete", "deletecollection", "get", "list", "watch"),
					),
					listvariable.DisplayName("Verb"),
					listvariable.Description("Filter by API verb"),
					listvariable.AllowAllValue(true),
					listvariable.AllowMultiple(true),
					listvariable.CustomAllValue(".*"),
					listvariable.DefaultValue("$__all"),
				),
			),
			dashboard.AddVariable("resource",
				listvariable.List(
					lokiLogQLVariable(lokiDatasource, fmt.Sprintf(auditVarQueryFmt, "objectRef_resource"), "objectRef_resource"),
					listvariable.DisplayName("Resource"),
					listvariable.Description("Filter by resource type (populated from audit logs)"),
					listvariable.AllowAllValue(true),
					listvariable.AllowMultiple(true),
					listvariable.CustomAllValue(".*"),
					listvariable.DefaultValue("$__all"),
				),
			),
			dashboard.AddVariable("resource_name",
				textvariable.Text("",
					textvariable.DisplayName("Resource Name"),
					textvariable.Description("Regex filter. Examples: my-pod.*, nginx, etcd.*"),
				),
			),
			dashboard.AddVariable("namespace",
				textvariable.Text("",
					textvariable.DisplayName("Namespace"),
					textvariable.Description("Regex filter. Examples: openshift-.*, my-app, kube-system|default"),
				),
			),
			dashboard.AddVariable("response_code",
				listvariable.List(
					staticListWithLabels(
						labeledValue{Value: "200", Label: "200 OK"},
						labeledValue{Value: "201", Label: "201 Created"},
						labeledValue{Value: "204", Label: "204 No Content"},
						labeledValue{Value: "304", Label: "304 Not Modified"},
						labeledValue{Value: "400", Label: "400 Bad Request"},
						labeledValue{Value: "401", Label: "401 Unauthorized"},
						labeledValue{Value: "403", Label: "403 Forbidden"},
						labeledValue{Value: "404", Label: "404 Not Found"},
						labeledValue{Value: "409", Label: "409 Conflict"},
						labeledValue{Value: "422", Label: "422 Unprocessable"},
						labeledValue{Value: "500", Label: "500 Internal Error"},
						labeledValue{Value: "503", Label: "503 Unavailable"},
					),
					listvariable.DisplayName("Response Code"),
					listvariable.Description("Filter by HTTP response code"),
					listvariable.AllowAllValue(true),
					listvariable.AllowMultiple(true),
					listvariable.CustomAllValue(".*"),
					listvariable.DefaultValue("$__all"),
				),
			),
			dashboard.AddVariable("exclude_resource",
				listvariable.List(
					staticlist.StaticList(
						staticlist.Values(
							"^$",
							"events",
							"endpoints",
							"endpointslices",
							"leases",
							"tokenreviews",
							"subjectaccessreviews",
							"selfsubjectaccessreviews",
							"selfsubjectrulesreviews",
						),
					),
					listvariable.DisplayName("Exclude Resources"),
					listvariable.Description("Select None to show all resource types"),
					listvariable.AllowAllValue(true),
					listvariable.AllowMultiple(true),
					listvariable.CustomAllValue(excludeResourceCustomAllValue),
					listvariable.DefaultValue("$__all"),
				),
			),
			dashboard.AddVariable("client",
				textvariable.Text("",
					textvariable.DisplayName("Client"),
					textvariable.Description("User agent regex. Examples: oc, kubectl, console, argocd"),
				),
			),
			dashboard.AddVariable("filter",
				textvariable.Text("",
					textvariable.DisplayName("LogQL Filter"),
					textvariable.Description(`Raw LogQL stage. Examples: |~ "error" (include), !~ "health" (exclude), | user_username!~"bot.*"`),
				),
			),

			dashboard.AddPanelGroup("Help",
				panelgroup.Collapsed(true),
				panelgroup.PanelHeight(5),
				panelgroup.AddPanel("Usage Guide",
					markdown.Markdown(auditHelpText),
				),
			),

			dashboard.AddPanelGroup("Active Query",
				panelgroup.Collapsed(true),
				panelgroup.PanelHeight(6),
				panelgroup.AddPanel("Current LogQL Query",
					markdown.Markdown(auditQueryDisplay),
				),
			),

			dashboard.AddPanelGroup("",
				panelgroup.PanelHeight(20),
				panels.AuditLogsPanel(lokiDatasource, auditLogQuery),
			),
		),
	).Component("openshift/logging")
}
