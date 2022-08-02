package http

import (
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/hashicorp/go-retryablehttp"
)

// HTTPLogAdaptor adapts the retryablehttp.Logger interface to the go-utils logger.
type HTTPLogAdaptor struct {
	logger log.Logger
}

// Printf implements the retryablehttp.Logger interface
func (a *HTTPLogAdaptor) Printf(fmtStr string, vars ...interface{}) {
	a.logger.Debugf(fmtStr, vars...)
}

// NewHTTPClient returns a retryable HTTP client with common defaults
func NewHTTPClient(logger log.Logger) *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.Logger = &HTTPLogAdaptor{}
	client.ErrorHandler = retryablehttp.PassthroughErrorHandler

	return client
}
