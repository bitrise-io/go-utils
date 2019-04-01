package log

type logMessage struct{
	LogLevel string `json:"log_level"`
	Message string `json:"message"`
	Data map[string]interface{} `json:"data"`
}

func (lm logMessage) SendToInternal(stepID, tag string, data map[string]interface{}) {
	// todo: post to analytics server
}

