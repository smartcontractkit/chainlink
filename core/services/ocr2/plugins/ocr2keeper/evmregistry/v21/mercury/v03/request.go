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

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury"
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

func (c *client) DoRequest(ctx context.Context, streamsLookup *mercury.StreamsLookup, pluginRetryKey string) (encoding.PipelineExecutionState, encoding.UpkeepFailureReason, [][]byte, bool, time.Duration, error) {
	if len(streamsLookup.Feeds) == 0 {
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonInvalidRevertDataInput, [][]byte{}, false, 0 * time.Second, fmt.Errorf("invalid revert data input: feed param key %s, time param key %s, feeds %s", streamsLookup.FeedParamKey, streamsLookup.TimeParamKey, streamsLookup.Feeds)
	}
	resultLen := 1 // Only 1 multi-feed request is made for all feeds
	ch := make(chan mercury.MercuryData, resultLen)
	c.threadCtrl.Go(func(ctx context.Context) {
		c.multiFeedsRequest(ctx, ch, streamsLookup)
	})

	var reqErr error
	var retryInterval time.Duration
	results := make([][]byte, len(streamsLookup.Feeds))
	retryable := false
	state := encoding.NoPipelineError

	m := <-ch
	if m.Error != nil {
		reqErr = m.Error
		retryable = m.Retryable
		state = m.State
		if retryable {
			retryInterval = mercury.CalculateRetryConfigFn(pluginRetryKey, c.mercuryConfig)
		}
	} else {
		results = m.Bytes
	}

	return state, encoding.UpkeepFailureReasonNone, results, retryable, retryInterval, reqErr
}

func (c *client) multiFeedsRequest(ctx context.Context, ch chan<- mercury.MercuryData, sl *mercury.StreamsLookup) {
	// this won't work bc q.Encode() will encode commas as '%2C' but the server is strictly expecting a comma separated list
	//q := url.Values{
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
		ch <- mercury.MercuryData{Index: 0, Error: err, Retryable: false, State: encoding.InvalidMercuryRequest}
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
	retryable := false
	sent := false
	retryErr := retry.Do(
		func() error {
			retryable = false
			resp, err := c.httpClient.Do(req)
			if err != nil {
				c.lggr.Warnf("at timestamp %s upkeep %s GET request fails from mercury v0.3: %v", sl.Time.String(), sl.UpkeepId.String(), err)
				retryable = true
				state = encoding.MercuryFlakyFailure
				return err
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				retryable = false
				state = encoding.InvalidMercuryResponse
				return err
			}

			c.lggr.Infof("at timestamp %s upkeep %s received status code %d from mercury v0.3", sl.Time.String(), sl.UpkeepId.String(), resp.StatusCode)
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				retryable = false
				state = encoding.UpkeepNotAuthorized
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3, most likely this is caused by unauthorized upkeep", sl.Time.String(), sl.UpkeepId.String(), resp.StatusCode)
			case http.StatusBadRequest:
				retryable = false
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3, most likely this is caused by invalid format of timestamp", sl.Time.String(), sl.UpkeepId.String(), resp.StatusCode)
			case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
				retryable = true
				state = encoding.MercuryFlakyFailure
				return fmt.Errorf("%d", resp.StatusCode)
			case http.StatusPartialContent:
				// TODO (AUTO-5044): handle response code 206 entirely with errors field parsing
				c.lggr.Warnf("at timestamp %s upkeep %s requested [%s] feeds but mercury v0.3 server returned 206 status, treating it as 404 and retrying", sl.Time.String(), sl.UpkeepId.String(), sl.Feeds)
				retryable = true
				state = encoding.MercuryFlakyFailure
				return fmt.Errorf("%d", http.StatusPartialContent)
			case http.StatusOK:
				// continue
			default:
				retryable = false
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3", sl.Time.String(), sl.UpkeepId.String(), resp.StatusCode)
			}
			c.lggr.Debugf("at block %s upkeep %s received status code %d from mercury v0.3 with BODY=%s", sl.Time.String(), sl.UpkeepId.String(), resp.StatusCode, hexutil.Encode(body))

			var response MercuryV03Response
			if err := json.Unmarshal(body, &response); err != nil {
				c.lggr.Warnf("at timestamp %s upkeep %s failed to unmarshal body to MercuryV03Response from mercury v0.3: %v", sl.Time.String(), sl.UpkeepId.String(), err)
				retryable = false
				state = encoding.MercuryUnmarshalError
				return err
			}

			// in v0.3, if some feeds are not available, the server will only return available feeds, but we need to make sure ALL feeds are retrieved before calling user contract
			// hence, retry in this case. retry will help when we send a very new timestamp and reports are not yet generated
			if len(response.Reports) != len(sl.Feeds) {
				var receivedFeeds []string
				for _, f := range response.Reports {
					receivedFeeds = append(receivedFeeds, f.FeedID)
				}
				c.lggr.Warnf("at timestamp %s upkeep %s mercury v0.3 server returned 206 status with [%s] reports while we requested [%s] feeds, retrying", sl.Time.String(), sl.UpkeepId.String(), receivedFeeds, sl.Feeds)
				retryable = true
				state = encoding.MercuryFlakyFailure
				return fmt.Errorf("%d", http.StatusNotFound)
			}
			var reportBytes [][]byte
			for _, rsp := range response.Reports {
				b, err := hexutil.Decode(rsp.FullReport)
				if err != nil {
					c.lggr.Warnf("at timestamp %s upkeep %s failed to decode reportBlob %s: %v", sl.Time.String(), sl.UpkeepId.String(), rsp.FullReport, err)
					retryable = false
					state = encoding.InvalidMercuryResponse
					return err
				}
				reportBytes = append(reportBytes, b)
			}
			ch <- mercury.MercuryData{
				Index:     0,
				Bytes:     reportBytes,
				Retryable: false,
				State:     encoding.NoPipelineError,
			}
			sent = true
			return nil
		},
		// only retry when the error is 206 Partial Content, 404 Not Found, 500 Internal Server Error, 502 Bad Gateway, 503 Service Unavailable, 504 Gateway Timeout
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusPartialContent) || err.Error() == fmt.Sprintf("%d", http.StatusNotFound) || err.Error() == fmt.Sprintf("%d", http.StatusInternalServerError) || err.Error() == fmt.Sprintf("%d", http.StatusBadGateway) || err.Error() == fmt.Sprintf("%d", http.StatusServiceUnavailable) || err.Error() == fmt.Sprintf("%d", http.StatusGatewayTimeout)
		}),
		retry.Context(ctx),
		retry.Delay(retryDelay),
		retry.Attempts(totalAttempt),
	)

	if !sent {
		ch <- mercury.MercuryData{
			Index:     0,
			Bytes:     [][]byte{},
			Retryable: retryable,
			Error:     retryErr,
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
