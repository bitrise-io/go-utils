package analytics

import (
	"bytes"
	"net/http"
	"time"

	"github.com/bitrise-io/go-utils/retry"
)

const trackEndpoint = "https://bitrise-step-analytics.herokuapp.com/track"
const timeOutInSecs = 30

// Client ...
type Client interface {
	Send(buffer *bytes.Buffer)
}

type client struct {
	httpClient *http.Client
	endpoint   string
}

// NewDefaultClient ...
func NewDefaultClient() Client {
	httpClient := retry.NewHTTPClient().StandardClient()
	httpClient.Timeout = time.Second * timeOutInSecs
	return NewClient(httpClient, trackEndpoint)
}

// NewClient ...
func NewClient(httpClient *http.Client, endpoint string) Client {
	return client{httpClient: httpClient, endpoint: endpoint}
}

// Send ...
func (t client) Send(buffer *bytes.Buffer) {
	_, _ = t.httpClient.Post(t.endpoint, "application/json", buffer)
}
