package httputil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/stretchr/testify/require"
)

// debugCapture collects Debugf output; other Logger methods are no-ops.
type debugCapture struct {
	log.Logger
	buf bytes.Buffer
}

func (d *debugCapture) Debugf(format string, v ...any) {
	fmt.Fprintf(&d.buf, format, v...)
}

func TestIsUserFixable(t *testing.T) {
	require.True(t, IsUserFixable(http.StatusBadRequest))
	require.True(t, IsUserFixable(http.StatusUnauthorized))
	require.False(t, IsUserFixable(http.StatusOK))
	require.False(t, IsUserFixable(http.StatusForbidden))
	require.False(t, IsUserFixable(http.StatusInternalServerError))
}

func TestPrintRequest_nil(t *testing.T) {
	logger := &debugCapture{}
	require.NoError(t, PrintRequest(logger, nil))
	require.Empty(t, logger.buf.String())
}

func TestPrintRequest_dumpsMethodPathAndBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/widgets", strings.NewReader(`{"name":"foo"}`))
	req.Header.Set("X-Test", "yes")

	logger := &debugCapture{}
	require.NoError(t, PrintRequest(logger, req))

	out := logger.buf.String()
	require.Contains(t, out, "POST /v1/widgets")
	require.Contains(t, out, "X-Test: yes")
	require.Contains(t, out, `{"name":"foo"}`)
}

func TestPrintResponse_nil(t *testing.T) {
	logger := &debugCapture{}
	require.NoError(t, PrintResponse(logger, nil))
	require.Empty(t, logger.buf.String())
}

func TestPrintResponse_dumpsStatusHeadersAndBody(t *testing.T) {
	resp := &http.Response{
		Status:     "201 Created",
		StatusCode: http.StatusCreated,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{
			"Content-Type": {"application/json"},
		},
		Body:          io.NopCloser(strings.NewReader(`{"id":1}`)),
		ContentLength: int64(len(`{"id":1}`)),
	}

	logger := &debugCapture{}
	require.NoError(t, PrintResponse(logger, resp))

	out := logger.buf.String()
	require.Contains(t, out, "201 Created")
	require.Contains(t, out, "Content-Type: application/json")
	require.Contains(t, out, `{"id":1}`)
}
