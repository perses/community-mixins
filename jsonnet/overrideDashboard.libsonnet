local overrideDashboard(dashboard, namespace, commonLabels, newDatasource) =
  dashboard {
    metadata+: {
      namespace: namespace,
      labels: dashboard.metadata.labels + commonLabels,
    },

    spec+: {
      panels: {
        [panelKey]: dashboard.spec.panels[panelKey] {
          spec+: if std.objectHas(dashboard.spec.panels[panelKey].spec, 'queries') then {
            queries: [
              query {
                spec+: {
                  plugin+: {
                    spec+: {
                      datasource+: {
                        name: newDatasource,
                      },
                    },
                  },
                },
              }
              for query in dashboard.spec.panels[panelKey].spec.queries
            ],
          } else {},
        }
        for panelKey in std.objectFields(dashboard.spec.panels)
      },

      variables: if std.objectHas(dashboard.spec, 'variables') then [
        variable {
          spec+: {
            plugin+: {
              spec+: {
                datasource+: {
                  name: newDatasource,
                },
              },
            },
          },
        }
        for variable in dashboard.spec.variables
      ] else [],
    },
  };
{
  overrideDashboard: overrideDashboard,
}
