package analytics

import (
	"encoding/json"
	"io"
	"time"

	"github.com/gofrs/uuid"
)

// EventDTO ...
type EventDTO struct {
	ID         string                 `json:"id"`
	EventName  string                 `json:"event_name"`
	Timestamp  int64                  `json:"timestamp"`
	Properties map[string]interface{} `json:"properties"`
}

// Event ...
type Event interface {
	toJSON(writer io.Writer, properties ...Property)
}

type event struct {
	ID         string
	EventName  string
	Timestamp  int64
	Properties []Property
}

// NewEvent ...
func NewEvent(name string, properties ...Property) Event {
	return event{
		ID:         uuid.Must(uuid.NewV4()).String(),
		EventName:  name,
		Timestamp:  time.Now().UnixNano() / int64(time.Millisecond),
		Properties: properties,
	}
}

func (e event) toJSON(writer io.Writer, shared ...Property) {
	dto := EventDTO{
		ID:        e.ID,
		EventName: e.EventName,
		Timestamp: e.Timestamp,
	}
	if len(shared) > 0 || len(e.Properties) > 0 {
		dto.Properties = map[string]interface{}{}
	}
	for _, property := range shared {
		dto.appendProperty(property)
	}
	for _, property := range e.Properties {
		dto.appendProperty(property)
	}
	if err := json.NewEncoder(writer).Encode(dto); err != nil {
		panic("Analytics event should be serializable to JSON")
	}
}

func (e EventDTO) appendProperty(property Property) {
	e.Properties[property.GetKey()] = property.GetValue()
}
