package gostats

import (
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/matrix"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/promql/parser"
)

var GoCommonPanelQueries = map[string]parser.Expr{
	"MemoryUsage_allocAll": vector.New(
		vector.WithMetricName("go_memstats_alloc_bytes"),
	),
	"MemoryUsage_allocHeap": vector.New(
		vector.WithMetricName("go_memstats_heap_alloc_bytes"),
	),
	"MemoryUsage_allocRateAll": promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName("go_memstats_alloc_bytes_total"),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	),
	"MemoryUsage_allocRateHeap": promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName("go_memstats_heap_alloc_bytes"),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	),
	"MemoryUsage_inuseStack": vector.New(
		vector.WithMetricName("go_memstats_stack_inuse_bytes"),
	),
	"MemoryUsage_inuseHeap": vector.New(
		vector.WithMetricName("go_memstats_heap_inuse_bytes"),
	),
	"MemoryUsage_processResident": vector.New(
		vector.WithMetricName("process_resident_memory_bytes"),
	),
	"Goroutines": vector.New(
		vector.WithMetricName("go_goroutines"),
	),
	"GarbageCollectionPauseTimeQuantiles": vector.New(
		vector.WithMetricName("go_gc_duration_seconds"),
	),
	"CPUUsage": promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName("process_cpu_seconds_total"),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	),
}

// OverrideGoPanelQueries overrides the GoCommonPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
// The convention of naming followed, is to use Panel function name (with _suffix, in case panel has multiple queries)
func OverrideGoPanelQueries(queries map[string]parser.Expr) {
	for k, v := range queries {
		GoCommonPanelQueries[k] = v
	}
}
