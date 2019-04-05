package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	analyticsServerURL = "https://bitrise-step-analytics.herokuapp.com"
	netClient = http.Client{
		Timeout: time.Second * 5,
	}
)

// Entry represents a line in a log
type Entry struct{
	LogLevel string `json:"log_level"`
	Message string `json:"message"`
	Data map[string]interface{} `json:"data"`
}

// SetAnalyticsServerURL updates the the analytics server collecting the
// logs. It is intended for use during tests. Warning: current implementation
// is not thread safe, do not call the function during runtime.
func SetAnalyticsServerURL(url string) {
	analyticsServerURL = url
}

// Internal sends the log message to the configured analytics server
func (lm Entry) Internal(stepID, tag string, data map[string]interface{}) {
	lm.Data = make(map[string]interface{})
	for k, v := range data {
		lm.Data[k] = v
	}

	if v, ok := lm.Data["step_id"]; ok {
		fmt.Printf("internal logger: data.step_id (%s) will be overriden with (%s) ", v, stepID)
	}
	if v, ok := lm.Data["tag"]; ok {
		fmt.Printf("internal logger: data.tag (%s) will be overriden with (%s) ", v, tag)
	}

	lm.Data["step_id"] = stepID
	lm.Data["tag"] = tag

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(lm); err != nil {
		return
	}

	go func() {
		if _, err := netClient.Post(analyticsServerURL + "/logs", "application/json", &b); err != nil {}
	}()

}

