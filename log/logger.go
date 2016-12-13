package log

import "fmt"
import "encoding/json"

// Logger ...
type Logger interface {
	PrintO(f Formatable)
	Printf(format string, a ...interface{})
}

// Formatable ...
type Formatable interface {
	String() string
	JSON() string
}

// Message ...
type Message struct {
	Content  interface{} `json:"content,omitempty"`
	Error    string      `json:"error,omitempty"`
	Warnings []string    `json:"warnings,omitempty"`
}

// String ...
func (m Message) String() string {
	msg := ""
	if m.Error != "" {
		msg = fmt.Sprintf("Error: %s", m.Error)
	} else {
		msg = fmt.Sprintf("%v", m.Content)
	}

	if len(m.Warnings) > 0 {
		msg += fmt.Sprintf("\nWarnings:\n")
		for i, warning := range m.Warnings {
			msg = fmt.Sprintf("- %s", warning)
			if i != len(m.Warnings)-1 {
				msg += "\n"
			}
		}
	}

	return msg
}

// JSON ...
func (m Message) JSON() string {
	bytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf(`"Error: %s"`, err)
	}
	return string(bytes)
}
