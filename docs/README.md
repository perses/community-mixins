# Dashboard Design Recommendations

We recommend following some standards when building dashboards, to keep the visualization looking good and making them better for usability. 

These are optional recommendations, there is no linting for this. 
If porting upstream existing dashboards, feel free to defer to original designs.

## Time Series Panels

* For generic non-histogram metric panels (counter/gauge): **Fully opaque area charts.** This ensures full visibility of the data, even against dark backgrounds, and enhances contrast when selecting individual result series.

```go
timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
```

* For request rate, error rates, error ratio/percentage panels: **Stacked fully opaque charts.** The opacity adds visual contrast, while stacking allows users to compare result data side by side and in relation to each other. This is especially helpful when visually distinguishing successful responses (2xx) from error responses (4xx, 5xx).

```go
timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
```

* For duration/quantile metric panels or metrics panels where there is no benefit in seeing result series relations: **Semi-opaque non-stacked.** This approach enables users to view all resulting data simultaneously, offering the flexibility to selectively explore specific ones while also making it easier to identify outliers in the data.

```go
timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
```
