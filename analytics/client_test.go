package analytics

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/v2/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const cientTimeout = 10 * time.Second

func Test_trackerClient_send_success(t *testing.T) {
	mockLogger := new(mocks.Logger)
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(b), "{}")
		assert.Equal(t, req.Method, http.MethodPost)
		assert.Equal(t, req.Header.Get("Content-Type"), "application/json")
		res.WriteHeader(200)
		_, err = res.Write([]byte("ok"))
		assert.NoError(t, err)
	}))
	defer func() { testServer.Close() }()
	client := NewClient(http.DefaultClient, testServer.URL, mockLogger, cientTimeout)
	client.Send(bytes.NewBufferString("{}"))
	mockLogger.AssertNotCalled(t, "Debugf", mock.Anything, mock.Anything)
}

func Test_trackerClient_send_failure(t *testing.T) {
	mockLogger := new(mocks.Logger)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(b), "{}")
		assert.Equal(t, req.Method, http.MethodPost)
		assert.Equal(t, req.Header.Get("Content-Type"), "application/json")
		res.WriteHeader(500)
		_, err = res.Write([]byte("failure"))
		assert.NoError(t, err)
	}))
	defer func() { testServer.Close() }()
	client := NewClient(http.DefaultClient, testServer.URL, mockLogger, cientTimeout)
	client.Send(bytes.NewBufferString("{}"))
	mockLogger.AssertCalled(t, "Debugf", "Couldn't send analytics event, status code: %d", 500)
}
