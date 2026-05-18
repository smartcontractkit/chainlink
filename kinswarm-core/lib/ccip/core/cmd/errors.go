package cmd

import (
	"fmt"
	"io"
	"net/http"
)

func httpError(resp *http.Response) error {
	errResult, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return fmt.Errorf("status %d %q: error reading body %w", resp.StatusCode, http.StatusText(resp.StatusCode), err2)
	}
	return fmt.Errorf("status %d %q: %s", resp.StatusCode, http.StatusText(resp.StatusCode), string(errResult))
}
