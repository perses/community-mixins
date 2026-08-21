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
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	logstable "github.com/perses/plugins/logstable/sdk/go"
	lokiquery "github.com/perses/plugins/loki/sdk/go/query/log"
)

// AuditLogsPanel creates a LogsTable panel that displays audit log entries from Loki.
//
// Parameters:
//   - datasourceName: The name of the Loki data source.
//   - query: The LogQL query expression for fetching audit logs.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func AuditLogsPanel(datasourceName, query string) panelgroup.Option {
	return panelgroup.AddPanel("Audit Logs",
		logstable.LogsTable(
			logstable.EnableDetails(true),
			logstable.ShowTime(true),
		),
		panel.AddQuery(
			lokiquery.LokiLogQuery(query,
				lokiquery.Datasource(datasourceName),
			),
		),
	)
}
