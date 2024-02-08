package v02

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	mercuryPathV02      = "/client?" // only used to access mercury v0.2 server
	retryDelay          = 500 * time.Millisecond
	totalAttempt        = 3
	contentTypeHeader   = "Content-Type"
	authorizationHeader = "Authorization"
	timestampHeader     = "X-Authorization-Timestamp"
	signatureHeader     = "X-Authorization-Signature-SHA256"
)

type MercuryV02Response struct {
	ChainlinkBlob string `json:"chainlinkBlob"`
}

type client struct {
	services.StateMachine
	mercuryConfig mercury.MercuryConfigProvider
	httpClient    mercury.HttpClient
	threadCtrl    utils.ThreadControl
	lggr          logger.Logger
}

func NewClient(mercuryConfig mercury.MercuryConfigProvider, httpClient mercury.HttpClient, threadCtrl utils.ThreadControl, lggr logger.Logger) *client {
	return &client{
		mercuryConfig: mercuryConfig,
		httpClient:    httpClient,
		threadCtrl:    threadCtrl,
		lggr:          lggr,
	}
}

func (c *client) DoRequest(ctx context.Context, streamsLookup *mercury.StreamsLookup, pluginRetryKey string) (encoding.PipelineExecutionState, encoding.UpkeepFailureReason, [][]byte, bool, time.Duration, error) {
	resultLen := len(streamsLookup.Feeds)
	ch := make(chan mercury.MercuryData, resultLen)
	if len(streamsLookup.Feeds) == 0 {
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonInvalidRevertDataInput, [][]byte{}, false, 0 * time.Second, fmt.Errorf("invalid revert data input: feed param key %s, time param key %s, feeds %s", streamsLookup.FeedParamKey, streamsLookup.TimeParamKey, streamsLookup.Feeds)
	}
	for i := range streamsLookup.Feeds {
		// TODO (AUTO-7209): limit the number of concurrent requests
		i := i
		c.threadCtrl.Go(func(ctx context.Context) {
			c.singleFeedRequest(ctx, ch, i, streamsLookup)
		})
	}

	var reqErr error
	var retryInterval time.Duration
	results := make([][]byte, len(streamsLookup.Feeds))
	retryable := true
	allSuccess := true
	// in v0.2, use the last execution error as the state, if no execution errors, state will be no error
	state := encoding.NoPipelineError
	for i := 0; i < resultLen; i++ {
		m := <-ch
		if m.Error != nil {
			reqErr = errors.Join(reqErr, m.Error)
			retryable = retryable && m.Retryable
			allSuccess = false
			if m.State != encoding.NoPipelineError {
				state = m.State
			}
			continue
		}
		results[m.Index] = m.Bytes[0]
	}
	if retryable && !allSuccess {
		retryInterval = mercury.CalculateRetryConfigFn(pluginRetryKey, c.mercuryConfig)
	}
	// only retry when not all successful AND none are not retryable
	return state, encoding.UpkeepFailureReasonNone, results, retryable && !allSuccess, retryInterval, reqErr
}

func (c *client) singleFeedRequest(ctx context.Context, ch chan<- mercury.MercuryData, index int, sl *mercury.StreamsLookup) {
	var httpRequest *http.Request
	var err error

	q := url.Values{
		sl.FeedParamKey: {sl.Feeds[index]},
		sl.TimeParamKey: {sl.Time.String()},
	}
	mercuryURL := c.mercuryConfig.Credentials().LegacyURL
	reqUrl := fmt.Sprintf("%s%s%s", mercuryURL, mercuryPathV02, q.Encode())
	c.lggr.Debugf("request URL for upkeep %s feed %s: %s", sl.UpkeepId.String(), sl.Feeds[index], reqUrl)

	httpRequest, err = http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		ch <- mercury.MercuryData{Index: index, Error: err, Retryable: false, State: encoding.InvalidMercuryRequest}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := mercury.GenerateHMACFn(http.MethodGet, mercuryPathV02+q.Encode(), []byte{}, c.mercuryConfig.Credentials().Username, c.mercuryConfig.Credentials().Password, ts)
	httpRequest.Header.Set(contentTypeHeader, "application/json")
	httpRequest.Header.Set(authorizationHeader, c.mercuryConfig.Credentials().Username)
	httpRequest.Header.Set(timestampHeader, strconv.FormatInt(ts, 10))
	httpRequest.Header.Set(signatureHeader, signature)

	// in the case of multiple retries here, use the last attempt's data
	state := encoding.NoPipelineError
	retryable := false
	sent := false
	retryErr := retry.Do(
		func() error {
			var httpResponse *http.Response
			var responseBody []byte
			var blobBytes []byte

			retryable = false
			if httpResponse, err = c.httpClient.Do(httpRequest); err != nil {
				c.lggr.Warnf("at block %s upkeep %s GET request fails for feed %s: %v", sl.Time.String(), sl.UpkeepId.String(), sl.Feeds[index], err)
				retryable = true
				state = encoding.MercuryFlakyFailure
				return err
			}
			defer httpResponse.Body.Close()

			if responseBody, err = io.ReadAll(httpResponse.Body); err != nil {
				state = encoding.InvalidMercuryResponse
				return err
			}

			switch httpResponse.StatusCode {
			case http.StatusNotFound, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
				c.lggr.Warnf("at block %s upkeep %s received status code %d for feed %s", sl.Time.String(), sl.UpkeepId.String(), httpResponse.StatusCode, sl.Feeds[index])
				retryable = true
				state = encoding.MercuryFlakyFailure
				return errors.New(strconv.FormatInt(int64(httpResponse.StatusCode), 10))
			case http.StatusOK:
				// continue
			default:
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("at block %s upkeep %s received status code %d for feed %s", sl.Time.String(), sl.UpkeepId.String(), httpResponse.StatusCode, sl.Feeds[index])
			}

			c.lggr.Debugf("at block %s upkeep %s received status code %d from mercury v0.2 with BODY=%s", sl.Time.String(), sl.UpkeepId.String(), httpResponse.StatusCode, hexutil.Encode(responseBody))

			var m MercuryV02Response
			if err = json.Unmarshal(responseBody, &m); err != nil {
				c.lggr.Warnf("at block %s upkeep %s failed to unmarshal body to MercuryV02Response for feed %s: %v", sl.Time.String(), sl.UpkeepId.String(), sl.Feeds[index], err)
				state = encoding.MercuryUnmarshalError
				return err
			}
			if blobBytes, err = hexutil.Decode(m.ChainlinkBlob); err != nil {
				c.lggr.Warnf("at block %s upkeep %s failed to decode chainlinkBlob %s for feed %s: %v", sl.Time.String(), sl.UpkeepId.String(), m.ChainlinkBlob, sl.Feeds[index], err)
				state = encoding.InvalidMercuryResponse
				return err
			}
			ch <- mercury.MercuryData{
				Index:     index,
				Bytes:     [][]byte{blobBytes},
				Retryable: false,
				State:     encoding.NoPipelineError,
			}
			sent = true
			return nil
		},
		// only retry when the error is 404 Not Found, 500 Internal Server Error, 502 Bad Gateway, 503 Service Unavailable, 504 Gateway Timeout
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound) || err.Error() == fmt.Sprintf("%d", http.StatusInternalServerError) || err.Error() == fmt.Sprintf("%d", http.StatusBadGateway) || err.Error() == fmt.Sprintf("%d", http.StatusServiceUnavailable) || err.Error() == fmt.Sprintf("%d", http.StatusGatewayTimeout)
		}),
		retry.Context(ctx),
		retry.Delay(retryDelay),
		retry.Attempts(totalAttempt),
	)

	if !sent {
		ch <- mercury.MercuryData{
			Index:     index,
			Bytes:     [][]byte{},
			Retryable: retryable,
			Error:     fmt.Errorf("failed to request feed for %s: %w", sl.Feeds[index], retryErr),
			State:     state,
		}
	}
}

func (c *client) Close() error {
	return c.StopOnce("v02_request", func() error {
		c.threadCtrl.Close()
		return nil
	})
}
