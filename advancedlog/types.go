package logger

// Producer ...
type Producer int

func (p Producer) String() string {
	switch p {
	case Step:
		return "step"
	default:
		return "cli"
	}
}

const (
	// CLI ...
	CLI Producer = iota
	// Step ...
	Step
)

// Level ...
type Level int

func (l Level) String() string {
	switch l {
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warn"
	case DoneLevel:
		return "done"
	case NormalLevel:
		return "normal"
	case DebugLevel:
		return "debug"
	default:
		return "info"
	}
}

const (
	// ErrorLevel ...
	ErrorLevel Level = iota
	// WarnLevel ...
	WarnLevel
	// InfoLevel ...
	InfoLevel
	// DoneLevel ...
	DoneLevel
	// NormalLevel ...
	NormalLevel
	// DebugLevel ...
	DebugLevel
)
