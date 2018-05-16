package log

// Logger ...
type Logger interface {
	Print(f Formatable)
}

// Formatable ...
type Formatable interface {
	String() string
	JSON() string
}

// DefaultLogger ...
type DefaultLogger struct {
	ts bool
}

type logfunc func(string, ...interface{})

// Timestamp ...
func (dl DefaultLogger) Timestamp() DefaultLogger {
	return DefaultLogger{ts: true}
}

// Donef ...
func (dl DefaultLogger) Donef(format string, v ...interface{}) {
	fSelect(dl.ts, TDonef, Donef)(format, v...)
}

// Successf ...
func (dl DefaultLogger) Successf(format string, v ...interface{}) {
	fSelect(dl.ts, TSuccessf, Successf)(format, v...)
}

// Infof ...
func (dl DefaultLogger) Infof(format string, v ...interface{}) {
	fSelect(dl.ts, TInfof, Infof)(format, v...)
}

// Printf ...
func (dl DefaultLogger) Printf(format string, v ...interface{}) {
	fSelect(dl.ts, TPrintf, Printf)(format, v...)
}

// Debugf ...
func (dl DefaultLogger) Debugf(format string, v ...interface{}) {
	if enableDebugLog {
		fSelect(dl.ts, TDebugf, Debugf)(format, v...)
	}
}

// Warnf ...
func (dl DefaultLogger) Warnf(format string, v ...interface{}) {
	fSelect(dl.ts, TWarnf, Warnf)(format, v...)
}

// Errorf ...
func (dl DefaultLogger) Errorf(format string, v ...interface{}) {
	fSelect(dl.ts, TErrorf, Errorf)(format, v...)
}

func fSelect(t bool, tf logfunc, f logfunc) logfunc {
	if t {
		return tf
	}
	return f
}
