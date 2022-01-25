package analytics

import (
	"encoding/json"
	"io"
	"time"

	"github.com/gofrs/uuid"
)

// EventBuilder ...
type EventBuilder interface {
	AddProperty(key string, value interface{}) EventBuilder
	Build() Event
}

type eventBuilder struct {
	eventName  string
	properties map[string]interface{}
}

// NewEventBuilder ...
func NewEventBuilder(name string) EventBuilder {
	return eventBuilder{eventName: name, properties: map[string]interface{}{}}
}

// AddProperty ...
func (e eventBuilder) AddProperty(key string, value interface{}) EventBuilder {
	e.properties[key] = value
	return e
}

// Build ...
func (e eventBuilder) Build() Event {
	return newEvent(e.eventName, e.properties)
}

// Event ...
type Event interface {
	toJson(writer io.Writer, shared map[string]interface{})
}

type event struct {
	ID         string                 `json:"id"`
	EventName  string                 `json:"event_name"`
	Timestamp  int64                  `json:"timestamp"`
	Properties map[string]interface{} `json:"properties"`
}

func newEvent(name string, properties map[string]interface{}) Event {
	return event{
		ID:         uuid.Must(uuid.NewV4()).String(),
		EventName:  name,
		Timestamp:  time.Now().UnixNano() / int64(time.Millisecond),
		Properties: properties,
	}
}

func (e event) toJson(writer io.Writer, shared map[string]interface{}) {
	if e.Properties == nil {
		e.Properties = map[string]interface{}{}
	}
	for k, v := range shared {
		_, ok := e.Properties[k] // More specific property takes precedence
		if !ok {
			e.Properties[k] = v
		}
	}
	if err := json.NewEncoder(writer).Encode(e); err != nil {
		panic("Analytics event should be serializable to JSON")
	}
}
