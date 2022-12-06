package command

type errorCollector struct {
	errorLines  []string
	errorFinder ErrorFinder
}

func (e *errorCollector) Write(p []byte) (n int, err error) {
	e.CollectErrors(string(p))
	return len(p), nil
}

func (e *errorCollector) CollectErrors(output string) {
	lines := e.errorFinder(output)
	if len(lines) > 0 {
		e.errorLines = append(e.errorLines, lines...)
	}
}
