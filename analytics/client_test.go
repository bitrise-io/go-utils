package analytics

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_trackerClient_send(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		b, _ := io.ReadAll(req.Body)
		assert.Equal(t, string(b), "{}")
		assert.Equal(t, req.Method, http.MethodPost)
		res.WriteHeader(200)
		_, _ = res.Write([]byte("ok"))
	}))
	defer func() { testServer.Close() }()
	client := NewClient(http.DefaultClient, testServer.URL)
	client.Send(bytes.NewBufferString("{}"))
}
