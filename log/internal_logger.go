package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type logMessage struct{
	LogLevel string `json:"log_level"`
	Message string `json:"message"`
	Data map[string]interface{} `json:"data"`
}

func (lm logMessage) SendToInternal(stepID, tag string, data map[string]interface{}) {
	for k, v := range data {
		lm.Data[k] = v
	}

	lm.Data["step_id"] = stepID
	lm.Data["tag"] = tag

	b, err := json.Marshal(lm)
	if err != nil {
		fmt.Printf("marshal log message: %s\n", err)
	}

	resp, err := http.Post("https://bitrise-step-analytics.herokuapp.com/logs", "application/json", bytes.NewReader(b))
	if err != nil {
		fmt.Printf("post log message: %s\n", err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("post log message response: %s %s\n", resp.Status, resp.Body)
	}
}

