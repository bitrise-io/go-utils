package filedownloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/retryhttp"
)

// Downloader provides methods for downloading files from remote URLs with automatic retries.
type Downloader interface {
	// Download fetches a remote file and writes it to the specified destination path.
	Download(ctx context.Context, destination, source string) error

	// DownloadWithFallback tries to download from the primary source first,
	// then falls back to alternative sources in order if the primary fails.
	DownloadWithFallback(ctx context.Context, destination, source string, fallbackSources ...string) error

	// Get returns a streaming reader for the remote content without buffering the entire response in memory.
	// Caller is responsible for closing the returned io.ReadCloser.
	Get(ctx context.Context, source string) (io.ReadCloser, error)
}

type downloader struct {
	client *http.Client
	logger log.Logger
}

// NewDownloader creates a new Downloader with automatic retry support.
func NewDownloader(logger log.Logger) Downloader {
	retryClient := retryhttp.NewClient(logger)
	return &downloader{
		client: retryClient.StandardClient(),
		logger: logger,
	}
}

// NewDownloaderWithClient creates a new Downloader with a custom HTTP client.
// Use this for advanced scenarios where you need custom client configuration.
// For most use cases, prefer NewDownloader which includes retry support by default.
func NewDownloaderWithClient(client *http.Client, logger log.Logger) Downloader {
	return &downloader{
		client: client,
		logger: logger,
	}
}

// Download fetches a remote file and writes it to the specified destination path.
func (d *downloader) Download(ctx context.Context, destination, source string) error {
	reader, err := d.Get(ctx, source)
	if err != nil {
		return err
	}
	defer reader.Close()

	f, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("create destination file %s: %w", destination, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return fmt.Errorf("write to destination file %s: %w", destination, err)
	}

	return nil
}

// DownloadWithFallback tries to download from the primary source first,
// then falls back to alternative sources in order if the primary fails.
func (d *downloader) DownloadWithFallback(ctx context.Context, destination, source string, fallbackSources ...string) error {
	sources := append([]string{source}, fallbackSources...)

	for _, src := range sources {
		err := d.Download(ctx, destination, src)
		if err != nil {
			d.logger.Warnf("Could not download file from %s: %s", src, err)
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to download from any source (tried %d sources)", len(sources))
}

// Get returns a streaming reader for the remote content without buffering the entire response in memory.
// Caller is responsible for closing the returned io.ReadCloser.
func (d *downloader) Get(ctx context.Context, source string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if err != nil {
		return nil, fmt.Errorf("create request for %s: %w", source, err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download from %s: %w", source, err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("download from %s: status code %d", source, resp.StatusCode)
	}

	return resp.Body, nil
}
