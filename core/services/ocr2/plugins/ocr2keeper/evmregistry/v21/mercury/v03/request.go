package v03

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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
	mercuryBatchPathV03            = "/api/v1/reports/bulk?"    // only used to access mercury v0.3 server
	mercuryBatchPathV03BlockNumber = "/api/v1gmx/reports/bulk?" // only used to access mercury v0.3 server with blockNumber
	retryDelay                     = 500 * time.Millisecond
	totalAttempt                   = 3
	contentTypeHeader              = "Content-Type"
	authorizationHeader            = "Authorization"
	timestampHeader                = "X-Authorization-Timestamp"
	signatureHeader                = "X-Authorization-Signature-SHA256"
	upkeepIDHeader                 = "X-Authorization-Upkeep-Id"
)

type MercuryV03Response struct {
	Reports []MercuryV03Report `json:"reports"`
}

type MercuryV03Report struct {
	FeedID                string `json:"feedID"` // feed id in hex encoded
	ValidFromTimestamp    uint32 `json:"validFromTimestamp"`
	ObservationsTimestamp uint32 `json:"observationsTimestamp"`
	FullReport            string `json:"fullReport"` // the actual hex encoded mercury report of this feed, can be sent to verifier
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
	resultLen := 1 // Only 1 multi-feed request is made for all feeds
	ch := make(chan mercury.MercuryData, resultLen)
	c.threadCtrl.GoCtx(ctx, func(ctx context.Context) {
		c.multiFeedsRequest(ctx, ch, streamsLookup)
	})

	ctx, cancel := context.WithTimeout(ctx, mercury.RequestTimeout)
	defer cancel()
	select {
	case <-ctx.Done():
		// Request Timed out, return timeout error
		c.lggr.Errorf("at timestamp %s upkeep %s, streams lookup v0.3 timed out", streamsLookup.Time.String(), streamsLookup.UpkeepId.String())
		return encoding.NoPipelineError, nil, encoding.ErrCodeStreamsTimeout, false, 0 * time.Second, nil
	case m := <-ch:
		if m.Error != nil {
			// There was a pipeline error during execution
			// If error was non retryable then just return the state and error
			if !m.Retryable {
				return m.State, nil, m.ErrCode, m.Retryable, 0 * time.Second, m.Error
			}
			// If errors were retryable then calculate retry interval
			retryInterval := mercury.CalculateStreamsRetryConfigFn(upkeepType, pluginRetryKey, c.mercuryConfig)
			if retryInterval != mercury.RetryIntervalTimeout {
				// Return the retyrable state with appropriate retry interval
				return m.State, nil, m.ErrCode, m.Retryable, retryInterval, m.Error
			}

			// Now we have exhausted all our retries. We treat it as not a pipeline error
			// and expose error code to the user
			return encoding.NoPipelineError, nil, m.ErrCode, false, 0 * time.Second, nil
		}

		// No pipeline error, return bytes and error code out of which one should be null
		return encoding.NoPipelineError, m.Bytes, m.ErrCode, false, 0 * time.Second, nil
	}
}

func (c *client) multiFeedsRequest(ctx context.Context, ch chan<- mercury.MercuryData, sl *mercury.StreamsLookup) {
	// this won't work bc q.Encode() will encode commas as '%2C' but the server is strictly expecting a comma separated list
	// q := url.Values{
	//	feedIDs:   {strings.Join(sl.Feeds, ",")},
	//	timestamp: {sl.Time.String()},
	//}

	params := fmt.Sprintf("%s=%s&%s=%s", mercury.FeedIDs, strings.Join(sl.Feeds, ","), mercury.Timestamp, sl.Time.String())
	batchPathV03 := mercuryBatchPathV03
	if sl.IsMercuryV03UsingBlockNumber() {
		batchPathV03 = mercuryBatchPathV03BlockNumber
	}
	reqUrl := fmt.Sprintf("%s%s%s", c.mercuryConfig.Credentials().URL, batchPathV03, params)

	c.lggr.Debugf("request URL for upkeep %s userId %s: %s", sl.UpkeepId.String(), c.mercuryConfig.Credentials().Username, reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		// Not a pipeline error, a bad streams request
		ch <- mercury.MercuryData{Index: 0, ErrCode: encoding.ErrCodeStreamsBadRequest, State: encoding.NoPipelineError}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := mercury.GenerateHMACFn(http.MethodGet, mercuryBatchPathV03+params, []byte{}, c.mercuryConfig.Credentials().Username, c.mercuryConfig.Credentials().Password, ts)

	req.Header.Set(contentTypeHeader, "application/json")
	// username here is often referred to as user id
	req.Header.Set(authorizationHeader, c.mercuryConfig.Credentials().Username)
	req.Header.Set(timestampHeader, strconv.FormatInt(ts, 10))
	req.Header.Set(signatureHeader, signature)
	// mercury will inspect authorization headers above to make sure this user (in automation's context, this node) is eligible to access mercury
	// and if it has an automation role. it will then look at this upkeep id to check if it has access to all the requested feeds.
	req.Header.Set(upkeepIDHeader, sl.UpkeepId.String())

	// in the case of multiple retries here, use the last attempt's data
	state := encoding.NoPipelineError
	errCode := encoding.ErrCodeNil
	retryable := false
	sent := false
	retryCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	retryErr := retry.Do(
		func() error {
			prommetrics.AutomationStreamsRetries.WithLabelValues(prommetrics.StreamsVersion03).Inc()
			retryable = false
			resp, err := c.httpClient.Do(req)
			if err != nil {
				c.lggr.Errorf("at timestamp %s upkeep %s GET request fails from mercury v0.3: %v", sl.Time.String(), sl.UpkeepId.String(), err)
				errCode = encoding.ErrCodeStreamsUnknownError
				if ctx.Err() != nil {
					errCode = encoding.ErrCodeStreamsTimeout
				}
				ch <- mercury.MercuryData{
					Index:   0,
					ErrCode: errCode,
					State:   encoding.NoPipelineError,
				}
				sent = true
				return nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				// Not a pipeline error, a bad streams response, send back error code
				ch <- mercury.MercuryData{
					Index:   0,
					ErrCode: encoding.ErrCodeStreamsBadResponse,
					State:   encoding.NoPipelineError,
				}
				sent = true
				return nil
			}

			c.lggr.Infof("at timestamp %s upkeep %s received status code %d from mercury v0.3", sl.Time.String(), sl.UpkeepId.String(), resp.StatusCode)
			prommetrics.AutomationStreamsResponses.WithLabelValues(prommetrics.StreamsVersion03, strconv.Itoa(resp.StatusCode)).Inc()
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				c.lggr.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3, most likely this is caused by unauthorized upkeep", sl.Time.String(), sl.UpkeepId.String(), resp.StatusCode)
				ch <- mercury.MercuryData{
					Index:   0,
					ErrCode: encoding.HttpToStreamsErrCode(resp.StatusCode),
					State:   encoding.NoPipelineError,
				}
				sent = true
				return nil
			case http.StatusBadRequest:
				c.lggr.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3, most likely this is caused by invalid format of timestamp", sl.Time.String(), sl.UpkeepId.String(), resp.StatusCode)
				ch <- mercury.MercuryData{
					Index:   0,
					ErrCode: encoding.HttpToStreamsErrCode(resp.StatusCode),
					State:   encoding.NoPipelineError,
				}
				sent = true
				return nil
			case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
				retryable = true
				state = encoding.MercuryFlakyFailure
				errCode = encoding.HttpToStreamsErrCode(resp.StatusCode)
				return fmt.Errorf("%d", resp.StatusCode)
			case http.StatusPartialContent:
				// TODO (AUTO-5044): handle response code 206 entirely with errors field parsing
				c.lggr.Warnf("at timestamp %s upkeep %s requested [%s] feeds but mercury v0.3 server returned 206 status, treating it as 404 and retrying", sl.Time.String(), sl.UpkeepId.String(), sl.Feeds)
				retryable = true
				state = encoding.MercuryFlakyFailure
				errCode = encoding.HttpToStreamsErrCode(resp.StatusCode)
				return fmt.Errorf("%d", http.StatusPartialContent)
			case http.StatusOK:
				// continue
			default:
				// Not considered as a pipeline error, a bad streams response with unknown status code. Send back to user as error code
				c.lggr.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3", sl.Time.String(), sl.UpkeepId.String(), resp.StatusCode)
				ch <- mercury.MercuryData{
					Index:   0,
					ErrCode: encoding.HttpToStreamsErrCode(resp.StatusCode),
					State:   encoding.NoPipelineError,
				}
				sent = true
				return nil
			}
			c.lggr.Debugf("at block %s upkeep %s received status code %d from mercury v0.3 with BODY=%s", sl.Time.String(), sl.UpkeepId.String(), resp.StatusCode, hexutil.Encode(body))

			var response MercuryV03Response
			if err := json.Unmarshal(body, &response); err != nil {
				c.lggr.Warnf("at timestamp %s upkeep %s failed to unmarshal body to MercuryV03Response from mercury v0.3: %v", sl.Time.String(), sl.UpkeepId.String(), err)
				ch <- mercury.MercuryData{
					Index:   0,
					ErrCode: encoding.ErrCodeStreamsBadResponse,
					State:   encoding.NoPipelineError,
				}
				sent = true
				return nil
			}

			// in v0.3, if some feeds are not available, the server will only return available feeds, but we need to make sure ALL feeds are retrieved before calling user contract
			// hence, retry in this case. retry will help when we send a very new timestamp and reports are not yet generated
			if len(response.Reports) != len(sl.Feeds) {
				var receivedFeeds []string
				for _, f := range response.Reports {
					receivedFeeds = append(receivedFeeds, f.FeedID)
				}
				c.lggr.Warnf("at timestamp %s upkeep %s mercury v0.3 server returned less reports [%s] while we requested [%s] feeds, retrying", sl.Time.String(), sl.UpkeepId.String(), receivedFeeds, sl.Feeds)
				retryable = true
				state = encoding.MercuryFlakyFailure
				errCode = encoding.HttpToStreamsErrCode(http.StatusPartialContent)
				return fmt.Errorf("%d", http.StatusPartialContent)
			}
			var reportBytes [][]byte
			for _, rsp := range response.Reports {
				b, err := hexutil.Decode(rsp.FullReport)
				if err != nil {
					c.lggr.Warnf("at timestamp %s upkeep %s failed to decode reportBlob %s: %v", sl.Time.String(), sl.UpkeepId.String(), rsp.FullReport, err)
					ch <- mercury.MercuryData{
						Index:   0,
						ErrCode: encoding.ErrCodeStreamsBadResponse,
						State:   encoding.NoPipelineError,
					}
					sent = true
					return nil
				}
				reportBytes = append(reportBytes, b)
			}
			ch <- mercury.MercuryData{
				Index: 0,
				Bytes: reportBytes,
				State: encoding.NoPipelineError,
			}
			sent = true
			return nil
		},
		// only retry when the error is 206 Partial Content, 404 Not Found, 500 Internal Server Error, 502 Bad Gateway, 503 Service Unavailable, 504 Gateway Timeout
		retry.RetryIf(func(err error) bool {
			return err.Error() == strconv.Itoa(http.StatusPartialContent) || err.Error() == strconv.Itoa(http.StatusNotFound) || err.Error() == strconv.Itoa(http.StatusInternalServerError) || err.Error() == strconv.Itoa(http.StatusBadGateway) || err.Error() == strconv.Itoa(http.StatusServiceUnavailable) || err.Error() == strconv.Itoa(http.StatusGatewayTimeout)
		}),
		retry.Context(retryCtx),
		retry.Delay(retryDelay),
		retry.Attempts(totalAttempt),
	)

	if !sent {
		ch <- mercury.MercuryData{
			Index:     0,
			Bytes:     nil,
			Retryable: retryable,
			Error:     retryErr,
			ErrCode:   errCode,
			State:     state,
		}
	}
}

func (c *client) Close() error {
	return c.StopOnce("v03_request", func() error {
		c.threadCtrl.Close()
		return nil
	})
}
