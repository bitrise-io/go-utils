package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"net/http"
)

var analyticsServerURL = "https://bitrise-step-analytics.herokuapp.com"

// Message represents a line in a log
type Message struct{
	LogLevel string `json:"log_level"`
	Message string `json:"message"`
	Data map[string]interface{} `json:"data"`
}

func init() {
	if url := os.Getenv("ANALYTICS_SERVER_URL"); url != "" {
		analyticsServerURL = url
	}
}

// SendToInternal sends the log message to the configured analytics server
func (lm Message) SendToInternal(stepID, tag string, data map[string]interface{}) {
	for k, v := range data {
		lm.Data[k] = v
	}

	lm.Data["step_id"] = stepID
	lm.Data["tag"] = tag

	b, err := json.Marshal(lm)
	if err != nil {
		fmt.Printf("marshal log message: %s\n", err)
	}

	resp, err := http.Post(analyticsServerURL + "/logs", "application/json", bytes.NewReader(b))
	if err != nil {
		fmt.Printf("post log message: %s\n", err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("post log message response: %s %s\n", resp.Status, resp.Body)
	}
}

