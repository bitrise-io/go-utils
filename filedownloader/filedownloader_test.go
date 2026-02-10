package filedownloader

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockHTTPClient struct {
	mock.Mock
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

// Test helper to create downloader with mock client
// Wraps the mock in a way that's compatible with *http.Client
func newTestDownloader(mockClient *mockHTTPClient, logger log.Logger) Downloader {
	// Create a custom http.Client that uses the mock's Do method
	client := &http.Client{
		Transport: &mockTransport{mock: mockClient},
	}
	d := &downloader{
		client: client,
		logger: logger,
	}
	return d
}

// mockTransport adapts mockHTTPClient to http.RoundTripper
type mockTransport struct {
	mock *mockHTTPClient
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.mock.Do(req)
}

type mockLogger struct {
	warnings []string
}

func (m *mockLogger) Infof(format string, v ...interface{})   {}
func (m *mockLogger) Warnf(format string, v ...interface{})   { m.warnings = append(m.warnings, format) }
func (m *mockLogger) Printf(format string, v ...interface{})  {}
func (m *mockLogger) Donef(format string, v ...interface{})   {}
func (m *mockLogger) Debugf(format string, v ...interface{})  {}
func (m *mockLogger) Errorf(format string, v ...interface{})  {}
func (m *mockLogger) TInfof(format string, v ...interface{})  {}
func (m *mockLogger) TWarnf(format string, v ...interface{})  {}
func (m *mockLogger) TPrintf(format string, v ...interface{}) {}
func (m *mockLogger) TDonef(format string, v ...interface{})  {}
func (m *mockLogger) TDebugf(format string, v ...interface{}) {}
func (m *mockLogger) TErrorf(format string, v ...interface{}) {}
func (m *mockLogger) Println()                                {}
func (m *mockLogger) PrintWithoutNewline(msg string)          {}
func (m *mockLogger) EnableDebugLog(enable bool)              {}

func TestDownloader_Get_Success(t *testing.T) {
	client := new(mockHTTPClient)
	logger := &mockLogger{}
	downloader := newTestDownloader(client, logger)

	expectedContent := "test file content"
	client.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(expectedContent)),
	}, nil)

	ctx := context.Background()
	reader, err := downloader.Get(ctx, "https://example.com/file.txt")

	require.NoError(t, err)
	require.NotNil(t, reader)
	defer func() { require.NoError(t, reader.Close()) }()

	content, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, expectedContent, string(content))
	client.AssertExpectations(t)
}

func TestDownloader_Get_NonOKStatus(t *testing.T) {
	client := new(mockHTTPClient)
	logger := &mockLogger{}
	downloader := newTestDownloader(client, logger)

	client.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil)

	ctx := context.Background()
	reader, err := downloader.Get(ctx, "https://example.com/notfound.txt")

	require.Error(t, err)
	assert.Nil(t, reader)
	assert.Contains(t, err.Error(), "status code 404")
	assert.Contains(t, err.Error(), "https://example.com/notfound.txt")
	client.AssertExpectations(t)
}

func TestDownloader_Get_NetworkError(t *testing.T) {
	client := new(mockHTTPClient)
	logger := &mockLogger{}
	downloader := newTestDownloader(client, logger)

	expectedError := errors.New("network connection failed")
	client.On("Do", mock.Anything).Return(nil, expectedError)

	ctx := context.Background()
	reader, err := downloader.Get(ctx, "https://example.com/file.txt")

	require.Error(t, err)
	assert.Nil(t, reader)
	assert.Contains(t, err.Error(), "network connection failed")
	client.AssertExpectations(t)
}

func TestDownloader_Get_ContextTimeout(t *testing.T) {
	client := new(mockHTTPClient)
	logger := &mockLogger{}
	downloader := newTestDownloader(client, logger)

	client.On("Do", mock.Anything).Return(nil, context.DeadlineExceeded)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	reader, err := downloader.Get(ctx, "https://example.com/file.txt")

	require.Error(t, err)
	assert.Nil(t, reader)
	client.AssertExpectations(t)
}

func TestDownloader_Download_Success(t *testing.T) {
	client := new(mockHTTPClient)
	logger := &mockLogger{}
	downloader := newTestDownloader(client, logger)

	expectedContent := "downloaded file content"
	client.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(expectedContent)),
	}, nil)

	tmpFile := t.TempDir() + "/downloaded.txt"
	ctx := context.Background()
	err := downloader.Download(ctx, tmpFile, "https://example.com/file.txt")

	require.NoError(t, err)

	content, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	assert.Equal(t, expectedContent, string(content))
	client.AssertExpectations(t)
}

func TestDownloader_Download_CreateFileError(t *testing.T) {
	client := new(mockHTTPClient)
	logger := &mockLogger{}
	downloader := newTestDownloader(client, logger)

	client.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("content")),
	}, nil)

	ctx := context.Background()
	err := downloader.Download(ctx, "/invalid/path/that/does/not/exist/file.txt", "https://example.com/file.txt")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "create destination file")
}

func TestDownloader_DownloadWithFallback_FirstSourceSuccess(t *testing.T) {
	client := new(mockHTTPClient)
	logger := &mockLogger{}
	downloader := newTestDownloader(client, logger)

	expectedContent := "content from first source"
	client.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(expectedContent)),
	}, nil).Once()

	tmpFile := t.TempDir() + "/fallback.txt"
	ctx := context.Background()
	err := downloader.DownloadWithFallback(ctx, tmpFile, "https://example.com/primary.txt", "https://example.com/fallback1.txt", "https://example.com/fallback2.txt")

	require.NoError(t, err)

	content, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	assert.Equal(t, expectedContent, string(content))
	assert.Empty(t, logger.warnings, "no warnings should be logged on first source success")
	client.AssertExpectations(t)
}

func TestDownloader_DownloadWithFallback_SecondSourceSuccess(t *testing.T) {
	client := new(mockHTTPClient)
	logger := &mockLogger{}
	downloader := newTestDownloader(client, logger)

	client.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil).Once()

	expectedContent := "content from second source"
	client.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(expectedContent)),
	}, nil).Once()

	tmpFile := t.TempDir() + "/fallback.txt"
	ctx := context.Background()
	err := downloader.DownloadWithFallback(ctx, tmpFile, "https://example.com/primary.txt", "https://example.com/fallback.txt")

	require.NoError(t, err)

	content, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	assert.Equal(t, expectedContent, string(content))
	assert.Len(t, logger.warnings, 1, "one warning should be logged for first source failure")
	client.AssertExpectations(t)
}

func TestDownloader_DownloadWithFallback_AllSourcesFail(t *testing.T) {
	client := new(mockHTTPClient)
	logger := &mockLogger{}
	downloader := newTestDownloader(client, logger)

	client.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil).Times(3)

	tmpFile := t.TempDir() + "/fallback.txt"
	ctx := context.Background()
	err := downloader.DownloadWithFallback(ctx, tmpFile, "https://example.com/primary.txt", "https://example.com/fallback1.txt", "https://example.com/fallback2.txt")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to download from any source")
	assert.Contains(t, err.Error(), "3 sources")
	assert.Len(t, logger.warnings, 3, "three warnings should be logged for all source failures")
	client.AssertExpectations(t)
}

func TestDownloader_Get_DoesNotBufferEntireResponse(t *testing.T) {
	client := new(mockHTTPClient)
	logger := &mockLogger{}
	downloader := newTestDownloader(client, logger)

	largeContent := strings.Repeat("a", 10*1024*1024) // 10 MB
	client.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(largeContent)),
	}, nil)

	ctx := context.Background()
	reader, err := downloader.Get(ctx, "https://example.com/largefile.bin")

	require.NoError(t, err)
	require.NotNil(t, reader)
	defer func() { require.NoError(t, reader.Close()) }()

	buf := make([]byte, 1024)
	n, err := reader.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, 1024, n)
	assert.Equal(t, strings.Repeat("a", 1024), string(buf))
	client.AssertExpectations(t)
}
