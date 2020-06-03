package gojspipe

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sebvautour/gojspipe/utils"
)

var promScriptDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "script_duration",
	Help: "Duration of script execution",
},
	[]string{"script"},
)

var promScriptReturnedFalseCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "script_returned_false_count",
	Help: "count of times scripts retruned a bool false, which stops the pipeline",
},
	[]string{"script"},
)

var promScriptExecutionCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "script_execution_count",
	Help: "count of script execution",
},
	[]string{"script"},
)

// PromCollectors returns promtheus metric Collectors used by this package
// prometheus.MustRegister(u.PromCollectors()...)
func PromCollectors() []prometheus.Collector {
	return []prometheus.Collector{promScriptDuration, promScriptReturnedFalseCount, promScriptExecutionCount, utils.PromScriptErrorCount}
}
