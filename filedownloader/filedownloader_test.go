package filedownloader

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
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
	mockedHTTPClient := givenMultiResponseHTTPClient(mockResponses)
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
			url: "http://url2.com",
			response: http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("filecontent2"))},
		},
	}
	mockedHTTPClient := givenMultiResponseHTTPClient(mockResponses)
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
	mockedHTTPClient := givenMultiResponseHTTPClient(mockResponses)
	downloader := givenFileDownloader(mockedHTTPClient)

	// // When
	actualErr := downloader.GetWithFallback(path, mockResponses[0].url, mockResponses[1].url)

	// // Then
	require.Equal(t, expectedErr, actualErr)
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Get(source string) (*http.Response, error) {
	args := m.Called(source)
	return args.Get(0).(*http.Response), args.Error(1)
}

func givenHTTPClient(response http.Response) *MockHTTPClient {
	mockedHTTPClient := new(MockHTTPClient)
	mockedHTTPClient.On("Get", mock.Anything).Return(&response, nil)
	return mockedHTTPClient
}

func givenFailingHTTPClient(err error) *MockHTTPClient {
	response := http.Response{StatusCode: 500}

	mockedHTTPClient := new(MockHTTPClient)
	mockedHTTPClient.On("Get", mock.Anything).Return(&response, err)
	return mockedHTTPClient
}

func givenMultiResponseHTTPClient(responses []mockResponse) *MockHTTPClient {
	mockedHTTPClient := new(MockHTTPClient)

	for _, mockResp := range responses {
		mockResp := mockResp
		mockedHTTPClient.On("Get", mockResp.url).Return(&mockResp.response, nil)
	}
	return mockedHTTPClient
}

type mockResponse struct {
	url      string
	response http.Response
}

func givenFileDownloader(client HTTPClient) HTTPFileDownloader {
	return HTTPFileDownloader{client}
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
