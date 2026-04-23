package httputil

import (
	"net/http"
	"net/http/httputil"

	"github.com/bitrise-io/go-utils/v2/log"
)

// PrintRequest dumps request (including body) to logger at debug level.
// A nil request is a no-op.
func PrintRequest(logger log.Logger, request *http.Request) error {
	if request == nil {
		return nil
	}

	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		return err
	}

	logger.Debugf("%s", dump)
	return nil
}

// PrintResponse dumps response (including body) to logger at debug level.
// A nil response is a no-op.
func PrintResponse(logger log.Logger, response *http.Response) error {
	if response == nil {
		return nil
	}

	dump, err := httputil.DumpResponse(response, true)
	if err != nil {
		return err
	}

	logger.Debugf("%s", dump)
	return nil
}

// IsUserFixable reports whether statusCode indicates an error the caller
// can correct (bad request or unauthorized).
func IsUserFixable(statusCode int) bool {
	return statusCode == http.StatusBadRequest || statusCode == http.StatusUnauthorized
}
