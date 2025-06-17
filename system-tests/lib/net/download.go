package net

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"
)

func Download(ctx context.Context, url string) ([]byte, error) {
	requestCtx, cancelFn := context.WithTimeout(ctx, 120*time.Second)
	defer cancelFn()

	req, err := http.NewRequestWithContext(requestCtx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

func DownloadAndDecodeBase64(ctx context.Context, url string) ([]byte, error) {
	data, err := Download(ctx, url)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 content: %w", err)
	}

	return decoded, nil
}
