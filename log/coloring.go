package log

// FormatWithSeverityColor ...
func FormatWithSeverityColor(s Severity, format string, v ...interface{}) string {
	return severityColorFuncMap[s](format, v...)
}
