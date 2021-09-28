package filedownloader

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bitrise-io/go-utils/log"
)

// HTTPClient ...
type HTTPClient interface {
	Get(source string) (*http.Response, error)
}

// DefaultFileDownloader ...
type DefaultFileDownloader struct {
	client HTTPClient
}

// New ...
func New(client HTTPClient) DefaultFileDownloader {
	return DefaultFileDownloader{client}
}

// GetWithFallback downloads a file from a given source. Provided destination should be a file that does not exist.
// You can specify fallback sources which will be used in order if downloading fails from either source.
func (downloader DefaultFileDownloader) GetWithFallback(destination, source string, fallbackSources ...string) error {
	sources := append([]string{source}, fallbackSources...)
	for _, source := range sources {
		err := downloader.Get(destination, source)
		if err != nil {
			log.Errorf("Could not download file from: %s", err)
		} else {
			log.Infof("URL used to download file: %s", source)
			return nil
		}
	}
	return fmt.Errorf("None of the sources returned 200 OK status")
}

// Get downloads a file from a given source. Provided destination should be a file that does not exist.
func (downloader DefaultFileDownloader) Get(destination, source string) error {
	f, err := os.Create(destination)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Errorf("Failed to close file, error: %s", err)
		}
	}()

	return download(downloader.client, source, f)
}

// GetRemoteContents fetches a remote URL contents
func (downloader DefaultFileDownloader) GetRemoteContents(URL string) ([]byte, error) {
	var contents []byte
	if err := download(downloader.client, URL, bytes.NewBuffer(contents)); err != nil {
		return nil, err
	}

	return contents, nil
}

// ReadLocalFile returns a local file contents
func (downloader DefaultFileDownloader) ReadLocalFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func download(client HTTPClient, source string, destination io.Writer) error {
	resp, err := client.Get(source)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to download file from: %s. Status code: %d", source, resp.StatusCode)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("Failed to close body, error: %s", err)
		}
	}()

	if _, err = io.Copy(destination, resp.Body); err != nil {
		return err
	}

	return nil
}
