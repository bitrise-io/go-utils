package retryhttp

import (
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/zachgrayio/go-retryablehttp"
)

// NewClient returns a retryable HTTP client with common defaults
func NewClientExperimental(logger log.Logger) *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.Logger = &httpLogAdaptor{logger: logger}
	client.ErrorHandler = retryablehttp.PassthroughErrorHandler

	return client
}
