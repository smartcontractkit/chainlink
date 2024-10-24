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

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	automationTypes "github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/prommetrics"
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

func (c *client) DoRequest(ctx context.Context, streamsLookup *mercury.StreamsLookup, upkeepType automationTypes.UpkeepType, pluginRetryKey string) (encoding.PipelineExecutionState, [][]byte, encoding.ErrCode, bool, time.Duration, error) {
	if len(streamsLookup.Feeds) == 0 {
		return encoding.NoPipelineError, nil, encoding.ErrCodeStreamsBadRequest, false, 0 * time.Second, nil
	}
	resultLen := len(streamsLookup.Feeds)
	ch := make(chan mercury.MercuryData, resultLen)
	for i := range streamsLookup.Feeds {
		// TODO (AUTO-7209): limit the number of concurrent requests
		i := i
		c.threadCtrl.GoCtx(ctx, func(ctx context.Context) {
			c.singleFeedRequest(ctx, ch, i, streamsLookup)
		})
	}

	ctx, cancel := context.WithTimeout(ctx, mercury.RequestTimeout)
	defer cancel()

	state := encoding.NoPipelineError
	var reqErr error
	retryable := true
	allFeedsPipelineSuccess := true
	allFeedsReturnedValues := true
	var errCode encoding.ErrCode
	results := make([][]byte, len(streamsLookup.Feeds))
	// in v0.2, when combining results for multiple feed requests
	// if any request resulted in pipeline execution error then use the last execution error as the state
	// if no execution errors, then check if any feed returned an error code, if so use the last error code
	for i := 0; i < resultLen; i++ {
		select {
		case <-ctx.Done():
			// Request Timed out, return timeout error
			c.lggr.Errorf("at block %s upkeep %s, streams lookup v0.2 timed out", streamsLookup.Time.String(), streamsLookup.UpkeepId.String())
			return encoding.NoPipelineError, nil, encoding.ErrCodeStreamsTimeout, false, 0 * time.Second, nil
		case m := <-ch:
			if m.Error != nil {
				state = m.State
				reqErr = errors.Join(reqErr, m.Error)
				retryable = retryable && m.Retryable
				if m.ErrCode != encoding.ErrCodeNil {
					// Some pipeline errors can get converted to error codes if retries are exhausted
					errCode = m.ErrCode
				}
				allFeedsPipelineSuccess = false
				continue
			}
			if m.ErrCode != encoding.ErrCodeNil {
				errCode = m.ErrCode
				allFeedsReturnedValues = false
				continue
			}
			// Feed request didn't face a pipeline error and didn't return an error code
			results[m.Index] = m.Bytes[0]
		}
	}

	if !allFeedsPipelineSuccess {
		// Some feeds faced a pipeline error during execution
		// If any error was non retryable then just return the state and error
		if !retryable {
			return state, nil, errCode, retryable, 0 * time.Second, reqErr
		}
		// If errors were retryable then calculate retry interval
		retryInterval := mercury.CalculateStreamsRetryConfigFn(upkeepType, pluginRetryKey, c.mercuryConfig)
		if retryInterval != mercury.RetryIntervalTimeout {
			// Return the retyrable state with appropriate retry interval
			return state, nil, errCode, retryable, retryInterval, reqErr
		}

		// Now we have exhausted all our retries. We treat it as not a pipeline error
		// and expose error code to the user
		return encoding.NoPipelineError, nil, errCode, false, 0 * time.Second, nil
	}

	// All feeds faced no pipeline error
	// If any feed request returned an error code, return the error code with empty values, else return the values
	if !allFeedsReturnedValues {
		return encoding.NoPipelineError, nil, errCode, false, 0 * time.Second, nil
	}

	// All success, return the results
	return encoding.NoPipelineError, results, encoding.ErrCodeNil, false, 0 * time.Second, nil
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
		// Not a pipeline error, a bad streams request
		ch <- mercury.MercuryData{Index: index, ErrCode: encoding.ErrCodeStreamsBadRequest, Bytes: nil, State: encoding.NoPipelineError}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := mercury.GenerateHMACFn(http.MethodGet, mercuryPathV02+q.Encode(), []byte{}, c.mercuryConfig.Credentials().Username, c.mercuryConfig.Credentials().Password, ts)
	httpRequest.Header.Set(contentTypeHeader, "application/json")
	httpRequest.Header.Set(authorizationHeader, c.mercuryConfig.Credentials().Username)
	httpRequest.Header.Set(timestampHeader, strconv.FormatInt(ts, 10))
	httpRequest.Header.Set(signatureHeader, signature)

	// in the case of multiple retries here, use the last attempt's data
	errCode := encoding.ErrCodeNil
	state := encoding.NoPipelineError
	retryable := false
	sent := false
	retryErr := retry.Do(
		func() error {
			prommetrics.AutomationStreamsRetries.WithLabelValues(prommetrics.StreamsVersion02).Inc()
			var httpResponse *http.Response
			var responseBody []byte
			var blobBytes []byte

			retryable = false
			if httpResponse, err = c.httpClient.Do(httpRequest); err != nil {
				c.lggr.Errorf("at block %s upkeep %s GET request fails for feed %s: %v", sl.Time.String(), sl.UpkeepId.String(), sl.Feeds[index], err)
				errCode = encoding.ErrCodeStreamsUnknownError
				if ctx.Err() != nil {
					errCode = encoding.ErrCodeStreamsTimeout
				}
				ch <- mercury.MercuryData{
					Index:   index,
					Bytes:   nil,
					ErrCode: errCode,
					State:   encoding.NoPipelineError,
				}
				sent = true
				return nil
			}
			defer httpResponse.Body.Close()

			if responseBody, err = io.ReadAll(httpResponse.Body); err != nil {
				// Not a pipeline error, a bad streams response, send back error code
				ch <- mercury.MercuryData{
					Index:   index,
					Bytes:   nil,
					ErrCode: encoding.ErrCodeStreamsBadResponse,
					State:   encoding.NoPipelineError,
				}
				sent = true
				return nil
			}

			prommetrics.AutomationStreamsResponses.WithLabelValues(prommetrics.StreamsVersion02, fmt.Sprintf("%d", httpResponse.StatusCode)).Inc()
			switch httpResponse.StatusCode {
			case http.StatusNotFound, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
				// Considered as pipeline error, but if retry attempts go over threshold, is changed upstream to ErrCode
				c.lggr.Warnf("at block %s upkeep %s received status code %d for feed %s", sl.Time.String(), sl.UpkeepId.String(), httpResponse.StatusCode, sl.Feeds[index])
				state = encoding.MercuryFlakyFailure
				retryable = true
				errCode = encoding.HttpToStreamsErrCode(httpResponse.StatusCode)
				return errors.New(strconv.FormatInt(int64(httpResponse.StatusCode), 10))
			case http.StatusOK:
				// continue
			default:
				// Not considered as a pipeline error, a bad streams response with unknown status code. Send back to user as error code
				c.lggr.Errorf("at block %s upkeep %s received unhandled status code %d for feed %s", sl.Time.String(), sl.UpkeepId.String(), httpResponse.StatusCode, sl.Feeds[index])
				ch <- mercury.MercuryData{
					Index:   index,
					Bytes:   nil,
					ErrCode: encoding.HttpToStreamsErrCode(httpResponse.StatusCode),
					State:   encoding.NoPipelineError,
				}
				sent = true
				return nil
			}

			c.lggr.Debugf("at block %s upkeep %s received status code %d from mercury v0.2 with BODY=%s", sl.Time.String(), sl.UpkeepId.String(), httpResponse.StatusCode, hexutil.Encode(responseBody))

			var m MercuryV02Response
			if err = json.Unmarshal(responseBody, &m); err != nil {
				c.lggr.Warnf("at block %s upkeep %s failed to unmarshal body to MercuryV02Response for feed %s: %v", sl.Time.String(), sl.UpkeepId.String(), sl.Feeds[index], err)
				ch <- mercury.MercuryData{
					Index:   index,
					Bytes:   nil,
					ErrCode: encoding.ErrCodeStreamsBadResponse,
					State:   encoding.NoPipelineError,
				}
				sent = true
				return nil
			}
			if blobBytes, err = hexutil.Decode(m.ChainlinkBlob); err != nil {
				c.lggr.Warnf("at block %s upkeep %s failed to decode chainlinkBlob %s for feed %s: %v", sl.Time.String(), sl.UpkeepId.String(), m.ChainlinkBlob, sl.Feeds[index], err)
				ch <- mercury.MercuryData{
					Index:   index,
					Bytes:   nil,
					ErrCode: encoding.ErrCodeStreamsBadResponse,
					State:   encoding.NoPipelineError,
				}
				sent = true
				return nil
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
			Bytes:     nil,
			ErrCode:   errCode,
			State:     state,
			Retryable: retryable,
			Error:     fmt.Errorf("failed to request feed for %s: %w", sl.Feeds[index], retryErr),
		}
	}
}

func (c *client) Close() error {
	return c.StopOnce("v02_request", func() error {
		c.threadCtrl.Close()
		return nil
	})
}
