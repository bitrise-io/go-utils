package analytics

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gofrs/uuid"
)

// Event ...
type Event interface {
	toJSON(writer io.Writer, properties ...Property)
}

type event struct {
	ID         string                 `json:"id"`
	EventName  string                 `json:"event_name"`
	Timestamp  int64                  `json:"timestamp"`
	Properties map[string]interface{} `json:"properties"`
}

// NewEvent ...
func NewEvent(name string, properties ...Property) Event {
	return event{
		ID:         uuid.Must(uuid.NewV4()).String(),
		EventName:  name,
		Timestamp:  time.Now().UnixNano() / int64(time.Millisecond),
		Properties: unwrap(properties),
	}
}

func (e event) toJSON(writer io.Writer, shared ...Property) {
	if len(shared) > 0 {
		if e.Properties == nil {
			e.Properties = map[string]interface{}{}
		}
		for _, property := range shared {
			if _, ok := e.Properties[property.GetKey()]; ok {
				continue
			}
			e.Properties[property.GetKey()] = property.GetValue()
		}
	}
	if err := json.NewEncoder(writer).Encode(e); err != nil {
		panic(fmt.Sprintf("Analytics event should be serializable to JSON: %s", err.Error()))
	}
}

func unwrap(properties []Property) map[string]interface{} {
	if len(properties) == 0 {
		return nil
	}
	m := map[string]interface{}{}
	for _, property := range properties {
		m[property.GetKey()] = property.GetValue()
	}
	return m
}
