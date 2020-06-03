package utils

import "github.com/prometheus/client_golang/prometheus"

// PromScriptErrorCount is used by the ReportError method
var PromScriptErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "script_error_count",
	Help: "count of errors reported by scripts",
},
	[]string{"script", "cmp"},
)
