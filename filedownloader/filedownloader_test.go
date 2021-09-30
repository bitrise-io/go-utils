package filedownloader

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_get_Success(t *testing.T) {
	// Given
	path := givenTempPath(t)
	mockedHTTPClient := givenHTTPClient(
		http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader("filecontent1")),
		})
	downloader := givenFileDownloader(mockedHTTPClient)

	// When
	err := downloader.Get(path, "http://url.com")

	// Then
	require.NoError(t, err)
	assertFileContent(t, path, "filecontent1")
}

func Test_get_InvalidStatusCode(t *testing.T) {
	// Given
	path := givenTempPath(t)
	url := "http://url.com"
	statusCode := 404
	expectedErr := fmt.Errorf("unable to download file from: %s. Status code: %d", url, statusCode)
	mockedHTTPClient := givenHTTPClient(
		http.Response{
			StatusCode: statusCode,
		})
	downloader := givenFileDownloader(mockedHTTPClient)

	// When
	err := downloader.Get(path, url)

	// Then
	require.Equal(t, expectedErr, err)
	assertFileNotExists(t, path)
}

func Test_get_HTTPError(t *testing.T) {
	// Given
	path := givenTempPath(t)
	expectedErr := errors.New("failed")
	downloader := givenFileDownloader(givenFailingHTTPClient(expectedErr))

	// When
	actualErr := downloader.Get(path, "http://url.com")

	// Then
	require.Equal(t, expectedErr, actualErr)
	assertFileNotExists(t, path)
}

func Test_GetWithFallback_FirstSuccess(t *testing.T) {
	// Given
	path := givenTempPath(t)
	mockResponses := []mockResponse{
		{
			url: "http://url1.com",
			response: http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("filecontent1"))},
		},
		{
			url:      "http://url2.com",
			response: http.Response{StatusCode: 400},
		},
	}
	mockedHTTPClient := givenMultiResponseClient(mockResponses)
	downloader := givenFileDownloader(mockedHTTPClient)

	// // When
	err := downloader.GetWithFallback(path, mockResponses[0].url, mockResponses[1].url)

	// // Then
	require.NoError(t, err)
	assertFileContent(t, path, "filecontent1")
}

func Test_GetWithFallback_SecondSuccess(t *testing.T) {
	// Given
	path := givenTempPath(t)
	mockResponses := []mockResponse{
		{
			url:      "http://url1.com",
			response: http.Response{StatusCode: 400},
		},
		{
			url: "http://url1.com",
			response: http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("filecontent2")),
			},
		},
	}
	mockedHTTPClient := givenMultiResponseClient(mockResponses)
	downloader := givenFileDownloader(mockedHTTPClient)

	// // When
	err := downloader.GetWithFallback(path, mockResponses[0].url, mockResponses[1].url)

	// // Then
	require.NoError(t, err)
	assertFileContent(t, path, "filecontent2")
}

func Test_GetWithFallback_NoneSuccess(t *testing.T) {
	// Given
	path := givenTempPath(t)
	expectedErr := errors.New("None of the sources returned 200 OK status")
	mockResponses := []mockResponse{
		{
			url:      "http://url1.com",
			response: http.Response{StatusCode: 400},
		},
		{
			url:      "http://url2.com",
			response: http.Response{StatusCode: 400},
		},
	}
	mockedHTTPClient := givenMultiResponseClient(mockResponses)
	downloader := givenFileDownloader(mockedHTTPClient)

	// // When
	actualErr := downloader.GetWithFallback(path, mockResponses[0].url, mockResponses[1].url)

	// // Then
	require.Equal(t, expectedErr, actualErr)
}

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockClient) GetRemoteContents(URL string) ([]byte, error) {
	args := m.Called(URL)
	arg0, ok := args.Get(0).([]byte)
	if !ok {
		panic("unexpected type")
	}

	return arg0, args.Error(1)
}

func (m *MockClient) ReadLocalFile(path string) ([]byte, error) {
	args := m.Called(path)
	arg0, ok := args.Get(0).([]byte)
	if !ok {
		panic("unexpected type")
	}

	return arg0, args.Error(1)
}

func givenHTTPClient(response http.Response) *MockClient {
	mockedClient := new(MockClient)
	mockedClient.On("Do", mock.Anything).Return(&response, nil)
	return mockedClient
}

func givenFailingHTTPClient(err error) *MockClient {
	response := http.Response{StatusCode: 500}

	mockedHTTPClient := new(MockClient)
	mockedHTTPClient.On("Do", mock.Anything).Return(&response, err)
	return mockedHTTPClient
}

func givenMultiResponseClient(responses []mockResponse) *MockClient {
	mockedHTTPClient := new(MockClient)

	for _, mockResp := range responses {
		mockResp := mockResp
		mockedHTTPClient.On("Do", mock.Anything).Return(&mockResp.response, nil).Once()
	}
	return mockedHTTPClient
}

type mockResponse struct {
	url      string
	response http.Response
}

func givenFileDownloader(client HTTPClient) FileDownloader {
	return FileDownloader{
		client: client,
	}
}

func givenTempPath(t *testing.T) string {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
	if err != nil {
		t.Errorf("Could not create tempDir: %s, error: %s", tmpDir, err)
	}
	return filepath.Join(tmpDir, "file.extension")
}

func assertFileNotExists(t *testing.T, path string) {
	if _, err := os.Stat(path); os.IsExist(err) {
		t.Fatalf("File should not exist at: %s", path)
	}
}

func assertFileContent(t *testing.T, path, expectedContent string) {
	actualContent, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Could not read file content at: %s. Error: %s", path, err)
	}

	require.Equal(t, expectedContent, string(actualContent))
}

func TestFileDownloader_GetRemoteContents(t *testing.T) {
	want := []byte{1}
	storage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write(want)
		if err != nil {
			t.Errorf("failed to write response: %s", err)
		}
	}))

	downloader := FileDownloader{
		client:  http.DefaultClient,
		context: nil,
	}
	got, err := downloader.GetRemoteContents(storage.URL)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
