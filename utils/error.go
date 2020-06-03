package utils

// ReportError logs a given error and increase the Prometheus script_error_count metric
// cmp should be a one-worded word describing the caller (function name for example)
func (u *Utils) ReportError(cmp string, err string, warning ...bool) {
	if cmp != "" {
		cmp += ": "
	}

	PromScriptErrorCount.WithLabelValues(u.ScriptName, cmp).Inc()

	if len(warning) != 0 && warning[0] == true {
		u.LogWarning(cmp + err)
		return
	}
	u.LogError(cmp + err)
}
