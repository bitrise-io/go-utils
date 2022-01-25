package analytics

// Event ...
type Event struct {
	ID         string                 `json:"id"`
	EventName  string                 `json:"event_name"`
	Timestamp  int64                  `json:"timestamp"`
	Properties map[string]interface{} `json:"properties"`
}
