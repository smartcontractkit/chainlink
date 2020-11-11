package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jpillora/backoff"
	"github.com/smartcontractkit/chainlink/core/logger"
)

type HTTPRequest struct {
	Request *http.Request
	Config  HTTPRequestConfig
}

// HTTPRequestConfig holds the configurable settings for an http request
type HTTPRequestConfig struct {
	Timeout                        time.Duration
	MaxAttempts                    uint
	SizeLimit                      int64
	AllowUnrestrictedNetworkAccess bool
}

func (h *HTTPRequest) SendRequest(ctx context.Context) (responseBody []byte, statusCode int, err error) {
	tr := &http.Transport{
		DisableCompression: true,
	}
	if !h.Config.AllowUnrestrictedNetworkAccess {
		tr.DialContext = restrictedDialContext
	}
	client := &http.Client{Transport: tr}

	return withRetry(ctx, client, h.Request, h.Config)
}

// withRetry executes the http request in a retry. Timeout is controlled with a context
// Retry occurs if the request timeout, or there is any kind of connection or transport-layer error
// Retry also occurs on remote server 5xx errors
func withRetry(
	ctx context.Context,
	client *http.Client,
	originalRequest *http.Request,
	config HTTPRequestConfig,
) (responseBody []byte, statusCode int, err error) {
	bb := &backoff.Backoff{
		Min:    100,
		Max:    20 * time.Minute, // We stop retrying on the number of attempts!
		Jitter: true,
	}
	for {
		responseBody, statusCode, err = makeHTTPCall(ctx, client, originalRequest, config)
		if err == nil {
			return responseBody, statusCode, nil
		}
		if uint(bb.Attempt())+1 >= config.MaxAttempts { // Stop retrying.
			return responseBody, statusCode, err
		}
		switch err.(type) {
		// There is no point in retrying a request if the response was
		// too large since it's likely that all retries will suffer the
		// same problem
		case *HTTPResponseTooLargeError:
			return responseBody, statusCode, err
		}
		// Sleep and retry.
		select {
		case <-ctx.Done():
			return responseBody, statusCode, ctx.Err()
		case <-time.After(bb.Duration()):
		}
		logger.Debugw("http adapter error, will retry", "error", err.Error(), "attempt", bb.Attempt(), "timeout", config.Timeout)
	}
}

func makeHTTPCall(
	ctx context.Context,
	client *http.Client,
	originalRequest *http.Request,
	config HTTPRequestConfig,
) (responseBody []byte, statusCode int, err error) {
	ctx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()
	requestWithTimeout := originalRequest.Clone(ctx)

	// XXX: Workaround for https://github.com/golang/go/issues/36095
	// http.Request#Clone actually only does a shallow copy
	originalRequestBody, err := originalRequest.GetBody()
	if err != nil {
		return nil, 0, err
	}
	var b bytes.Buffer
	_, err = b.ReadFrom(originalRequestBody)
	if err != nil {
		return nil, 0, err
	}
	requestWithTimeout.Body = ioutil.NopCloser(&b)

	start := time.Now()

	r, err := client.Do(requestWithTimeout)
	if err != nil {
		logger.Warnw("http adapter got error", "error", err)
		return nil, 0, err
	}
	defer logger.ErrorIfCalling(r.Body.Close)

	statusCode = r.StatusCode
	elapsed := time.Since(start)
	logger.Debugw(fmt.Sprintf("http adapter got %v in %s", statusCode, elapsed), "statusCode", statusCode, "timeElapsedSeconds", elapsed)

	source := NewMaxBytesReader(r.Body, config.SizeLimit)
	bytes, err := ioutil.ReadAll(source)
	if err != nil {
		logger.Errorw("http adapter error reading body", "error", err)
		return nil, statusCode, err
	}
	elapsed = time.Since(start)
	logger.Debugw(fmt.Sprintf("http adapter finished after %s", elapsed), "statusCode", statusCode, "timeElapsedSeconds", elapsed)

	responseBody = bytes

	// Retry on 5xx since this might give a different result
	if 500 <= r.StatusCode && r.StatusCode < 600 {
		return responseBody, statusCode, &RemoteServerError{responseBody, statusCode}
	}

	return responseBody, statusCode, nil
}

type RemoteServerError struct {
	responseBody []byte
	statusCode   int
}

func (e *RemoteServerError) Error() string {
	return fmt.Sprintf("remote server error: %v\nResponse body: %v", e.statusCode, string(e.responseBody))
}

// MaxBytesReader is inspired by
// https://github.com/gin-contrib/size/blob/master/size.go
type MaxBytesReader struct {
	rc               io.ReadCloser
	limit, remaining int64
	sawEOF           bool
}

func NewMaxBytesReader(rc io.ReadCloser, limit int64) *MaxBytesReader {
	return &MaxBytesReader{
		rc:        rc,
		limit:     limit,
		remaining: limit,
	}
}

func (mbr *MaxBytesReader) Read(p []byte) (n int, err error) {
	toRead := mbr.remaining
	if mbr.remaining == 0 {
		if mbr.sawEOF {
			return mbr.tooLarge()
		}
		// The underlying io.Reader may not return (0, io.EOF)
		// at EOF if the requested size is 0, so read 1 byte
		// instead. The io.Reader docs are a bit ambiguous
		// about the return value of Read when 0 bytes are
		// requested, and {bytes,strings}.Reader gets it wrong
		// too (it returns (0, nil) even at EOF).
		toRead = 1
	}
	if int64(len(p)) > toRead {
		p = p[:toRead]
	}
	n, err = mbr.rc.Read(p)
	if err == io.EOF {
		mbr.sawEOF = true
	}
	if mbr.remaining == 0 {
		// If we had zero bytes to read remaining (but hadn't seen EOF)
		// and we get a byte here, that means we went over our limit.
		if n > 0 {
			return mbr.tooLarge()
		}
		return 0, err
	}
	mbr.remaining -= int64(n)
	if mbr.remaining < 0 {
		mbr.remaining = 0
	}
	return
}

type HTTPResponseTooLargeError struct {
	limit int64
}

func (e *HTTPResponseTooLargeError) Error() string {
	return fmt.Sprintf("HTTP response too large, must be less than %d bytes", e.limit)
}

func (mbr *MaxBytesReader) tooLarge() (int, error) {
	return 0, &HTTPResponseTooLargeError{mbr.limit}
}

func (mbr *MaxBytesReader) Close() error {
	return mbr.rc.Close()
}
