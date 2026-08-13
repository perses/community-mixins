# Variables de Datasource en Dashboards de Istio

## Resumen

Sí, es posible usar variables para el datasource en los dashboards de Istio. Esta funcionalidad permite a los usuarios seleccionar dinámicamente la fuente de datos desde la interfaz del dashboard.

## Implementación

### 1. Función Helper

Se ha creado una función helper en `pkg/dashboards/helpers.go`:

```go
// AddDatasourceVariable creates a datasource variable that allows users to select
// the datasource dynamically in the dashboard
func AddDatasourceVariable(variableName string, defaultValue string) dashboard.Option {
	return dashboard.AddVariable(variableName,
		datasourceVar.Datasource(
			datasourceVar.DisplayName("Datasource"),
			datasourceVar.DefaultValue(defaultValue),
		),
	)
}
```

### 2. Dashboards con Variables de Datasource

Se han creado versiones de los dashboards de Istio que incluyen variables de datasource:

#### Istio Mesh Dashboard con Variable de Datasource

```go
func BuildIstioMeshWithDatasourceVariable(project string, defaultDatasource string, clusterLabelName string) dashboards.DashboardResult {
	emptyLabelMatcher := &labels.Matcher{}
	return dashboards.NewDashboardResult(
		dashboard.New("istio-mesh-with-datasource-var",
			dashboard.ProjectName(project),
			dashboard.Name("Istio Mesh Dashboard (with Datasource Variable)"),
			// Add datasource variable
			dashboards.AddDatasourceVariable("datasource", defaultDatasource),
			// Use the datasource variable in panels
			withMeshOverview("$datasource", emptyLabelMatcher),
			withMeshWorkloads("$datasource", emptyLabelMatcher),
			withIstioComponentVersions("$datasource", emptyLabelMatcher),
		),
	).Component("istio")
}
```

#### Istio Service Dashboard con Variable de Datasource

```go
func BuildIstioServiceWithDatasourceVariable(project string, defaultDatasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("istio-service-dashboard-with-datasource-var",
			dashboard.ProjectName(project),
			dashboard.Name("Istio Service Dashboard (with Datasource Variable)"),
			// Add datasource variable
			dashboards.AddDatasourceVariable("datasource", defaultDatasource),
			// All variables and panels now use $datasource
			// ... resto de la configuración
		),
	).Component("istio")
}
```

## Uso

### 1. En el código Go

```go
// Crear dashboard con variable de datasource
dashboardWriter.Add(istio.BuildIstioMeshWithDatasourceVariable(project, "prometheus-datasource", clusterLabelName))
dashboardWriter.Add(istio.BuildIstioServiceWithDatasourceVariable(project, "prometheus-datasource", clusterLabelName))
```

### 2. En la interfaz del dashboard

Una vez generado el dashboard, los usuarios verán:

1. **Variable de Datasource**: Un dropdown en la parte superior del dashboard que permite seleccionar entre diferentes fuentes de datos disponibles
2. **Valor por defecto**: Se establece el datasource por defecto especificado en la función
3. **Actualización automática**: Todos los paneles y variables se actualizan automáticamente cuando se cambia el datasource

## Ventajas

1. **Flexibilidad**: Los usuarios pueden cambiar entre diferentes fuentes de datos sin modificar el dashboard
2. **Reutilización**: Un mismo dashboard puede usarse con múltiples fuentes de datos
3. **Facilidad de uso**: No requiere conocimiento técnico para cambiar el datasource
4. **Compatibilidad**: Funciona con cualquier fuente de datos compatible con Prometheus

## Ejemplo de Configuración

```yaml
# En un archivo de configuración
datasources:
  - name: "prometheus-prod"
    url: "http://prometheus-prod:9090"
  - name: "prometheus-staging"
    url: "http://prometheus-staging:9090"
  - name: "thanos-query"
    url: "http://thanos-query:9090"
```

El dashboard permitirá seleccionar entre estas opciones dinámicamente.

## Notas Técnicas

- La variable de datasource se referencia como `$datasource` en las consultas
- Todos los paneles y variables del dashboard deben usar esta variable
- El valor por defecto se establece al crear el dashboard
- La funcionalidad es compatible con el ecosistema completo de Perses


