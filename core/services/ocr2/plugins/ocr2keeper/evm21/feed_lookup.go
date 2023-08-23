package evm

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/encoding"
)

const (
	applicationJson     = "application/json"
	blockNumber         = "blockNumber" // valid for v0.2
	feedIDs             = "feedIDs"     // valid for v0.3
	feedIdHex           = "feedIdHex"   // valid for v0.2
	headerAuthorization = "Authorization"
	headerContentType   = "Content-Type"
	headerTimestamp     = "X-Authorization-Timestamp"
	headerSignature     = "X-Authorization-Signature-SHA256"
	headerUpkeepId      = "X-Authorization-Upkeep-Id"
	mercuryPathV02      = "/client?"          // only used to access mercury v0.2 server
	mercuryBatchPathV03 = "/v1/reports/bulk?" // only used to access mercury v0.3 server
	retryDelay          = 500 * time.Millisecond
	timestamp           = "timestamp" // valid for v0.3
	totalAttempt        = 3
)

type MercuryVersion string

type FeedLookup struct {
	feedParamKey string
	feeds        []string
	timeParamKey string
	time         *big.Int
	extraData    []byte
	upkeepId     *big.Int
	block        uint64
}

// MercuryV02Response represents a JSON structure used by Mercury v0.2
type MercuryV02Response struct {
	ChainlinkBlob string `json:"chainlinkBlob"`
}

// MercuryV03Response represents a JSON structure used by Mercury v0.3
type MercuryV03Response struct {
	FeedID                string `json:"feedID"`
	ValidFromTimestamp    string `json:"validFromTimestamp"`
	ObservationsTimestamp string `json:"observationsTimestamp"`
	FullReport            string `json:"fullReport"`
}

type MercuryData struct {
	Index     int
	Error     error
	Retryable bool
	Bytes     [][]byte
	State     encoding.PipelineExecutionState
}

// UpkeepPrivilegeConfig represents the administrative offchain config for each upkeep. It can be set by s_upkeepManager
// role on the registry. Upkeeps allowed to use Mercury server will have this set to true.
type UpkeepPrivilegeConfig struct {
	MercuryEnabled bool `json:"mercuryEnabled"`
}

// feedLookup looks through check upkeep results looking for any that need off chain lookup
func (r *EvmRegistry) feedLookup(ctx context.Context, checkResults []ocr2keepers.CheckResult) []ocr2keepers.CheckResult {
	lggr := r.lggr.With("where", "FeedLookup")
	lookups := map[int]*FeedLookup{}
	for i, res := range checkResults {
		if res.IneligibilityReason != uint8(encoding.UpkeepFailureReasonTargetCheckReverted) {
			// Feedlookup only works when upkeep target check reverts
			continue
		}

		block := big.NewInt(int64(res.Trigger.BlockNumber))
		upkeepId := res.UpkeepID

		// Try to decode the revert error into feed lookup format. User upkeeps can revert with any reason, see if they
		// tried to call mercury
		lggr.Infof("upkeep %s block %d trying to decodeFeedLookup performData=%s", upkeepId, block, hexutil.Encode(checkResults[i].PerformData))
		lookup, err := r.decodeFeedLookup(res.PerformData)
		if err != nil {
			lggr.Warnf("upkeep %s block %d decodeFeedLookup failed: %v", upkeepId, block, err)
			// Not feed lookup error, nothing to do here
			continue
		}

		opts := r.buildCallOpts(ctx, block)
		// TODO: Remove allowlist for v0.3
		state, retryable, allowed, err := r.allowedToUseMercury(opts, upkeepId.BigInt())
		if err != nil {
			lggr.Warnf("upkeep %s block %d failed to query mercury allow list: %s", upkeepId, block, err)
			checkResults[i].PipelineExecutionState = uint8(state)
			checkResults[i].Retryable = retryable
			continue
		}

		if !allowed {
			lggr.Warnf("upkeep %s block %d NOT allowed to query Mercury server", upkeepId, block)
			checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonMercuryAccessNotAllowed)
			continue
		}

		lookup.upkeepId = upkeepId.BigInt()
		// the block here is exclusively used to call checkCallback at this block, not to be confused with the block number
		// in the revert for mercury v0.2, which is denoted by time in the struct bc starting from v0.3, only timestamp will be supported
		lookup.block = uint64(block.Int64())
		lggr.Infof("upkeep %s block %d decodeFeedLookup feedKey=%s timeKey=%s feeds=%v time=%s extraData=%s", upkeepId, block, lookup.feedParamKey, lookup.timeParamKey, lookup.feeds, lookup.time, hexutil.Encode(lookup.extraData))
		lookups[i] = lookup
	}

	var wg sync.WaitGroup
	for i, lookup := range lookups {
		wg.Add(1)
		go r.doLookup(ctx, &wg, lookup, i, checkResults, lggr)
	}
	wg.Wait()

	// don't surface error to plugin bc FeedLookup process should be self-contained.
	return checkResults
}

func (r *EvmRegistry) doLookup(ctx context.Context, wg *sync.WaitGroup, lookup *FeedLookup, i int, checkResults []ocr2keepers.CheckResult, lggr logger.Logger) {
	defer wg.Done()

	state, values, retryable, err := r.doMercuryRequest(ctx, lookup, lggr)
	if err != nil {
		lggr.Errorf("upkeep %s retryable %v doMercuryRequest: %v", lookup.upkeepId, retryable, err)
		checkResults[i].Retryable = retryable
		checkResults[i].PipelineExecutionState = uint8(state)
		return
	}
	for j, v := range values {
		lggr.Infof("checkCallback values[%d]=%s", j, hexutil.Encode(v))
	}

	state, retryable, mercuryBytes, err := r.checkCallback(ctx, values, lookup)
	if err != nil {
		lggr.Errorf("upkeep %s block %d checkCallback err: %v", lookup.upkeepId, lookup.block, err)
		checkResults[i].Retryable = retryable
		checkResults[i].PipelineExecutionState = uint8(state)
		return
	}
	lggr.Infof("checkCallback mercuryBytes=%s", hexutil.Encode(mercuryBytes))

	state, needed, performData, failureReason, _, err := r.packer.UnpackCheckCallbackResult(mercuryBytes)
	if err != nil {
		lggr.Errorf("upkeep %s block %d UnpackCheckCallbackResult err: %v", lookup.upkeepId, lookup.block, err)
		checkResults[i].PipelineExecutionState = uint8(state)
		return
	}

	if failureReason == uint8(encoding.UpkeepFailureReasonMercuryCallbackReverted) {
		checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonMercuryCallbackReverted)
		lggr.Debugf("upkeep %s block %d mercury callback reverts", lookup.upkeepId, lookup.block)
		return
	}

	if !needed {
		checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonUpkeepNotNeeded)
		lggr.Debugf("upkeep %s block %d callback reports upkeep not needed", lookup.upkeepId, lookup.block)
		return
	}

	checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonNone)
	checkResults[i].Eligible = true
	checkResults[i].PerformData = performData
	lggr.Infof("upkeep %s block %d successful with perform data: %s", lookup.upkeepId, lookup.block, hexutil.Encode(performData))
}

// allowedToUseMercury retrieves upkeep's administrative offchain config and decode a mercuryEnabled bool to indicate if
// this upkeep is allowed to use Mercury service.
func (r *EvmRegistry) allowedToUseMercury(opts *bind.CallOpts, upkeepId *big.Int) (state encoding.PipelineExecutionState, retryable bool, allow bool, err error) {
	allowed, ok := r.mercury.allowListCache.Get(upkeepId.String())
	if ok {
		return encoding.NoPipelineError, false, allowed.(bool), nil
	}

	cfg, err := r.registry.GetUpkeepPrivilegeConfig(opts, upkeepId)
	if err != nil {
		return encoding.RpcFlakyFailure, true, false, fmt.Errorf("failed to get upkeep privilege config for upkeep ID %s: %v", upkeepId, err)
	}

	var a UpkeepPrivilegeConfig
	err = json.Unmarshal(cfg, &a)
	if err != nil {
		return encoding.MercuryUnmarshalError, false, false, fmt.Errorf("failed to unmarshal privilege config for upkeep ID %s: %v", upkeepId, err)
	}
	r.mercury.allowListCache.Set(upkeepId.String(), a.MercuryEnabled, cache.DefaultExpiration)
	return encoding.NoPipelineError, false, a.MercuryEnabled, nil
}

// decodeFeedLookup decodes the revert error FeedLookup(string feedParamKey, string[] feeds, string feedParamKey, uint256 time, byte[] extraData)
func (r *EvmRegistry) decodeFeedLookup(data []byte) (*FeedLookup, error) {
	e := r.mercury.abi.Errors["FeedLookup"]
	unpack, err := e.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("unpack error: %w", err)
	}
	errorParameters := unpack.([]interface{})

	return &FeedLookup{
		feedParamKey: *abi.ConvertType(errorParameters[0], new(string)).(*string),
		feeds:        *abi.ConvertType(errorParameters[1], new([]string)).(*[]string),
		timeParamKey: *abi.ConvertType(errorParameters[2], new(string)).(*string),
		time:         *abi.ConvertType(errorParameters[3], new(*big.Int)).(**big.Int),
		extraData:    *abi.ConvertType(errorParameters[4], new([]byte)).(*[]byte),
	}, nil
}

func (r *EvmRegistry) checkCallback(ctx context.Context, values [][]byte, lookup *FeedLookup) (encoding.PipelineExecutionState, bool, hexutil.Bytes, error) {
	payload, err := r.abi.Pack("checkCallback", lookup.upkeepId, values, lookup.extraData)
	if err != nil {
		return encoding.PackUnpackDecodeFailed, false, nil, err
	}

	var b hexutil.Bytes
	args := map[string]interface{}{
		"to":   r.addr.Hex(),
		"data": hexutil.Bytes(payload),
	}

	// call checkCallback function at the block which OCR3 has agreed upon
	err = r.client.CallContext(ctx, &b, "eth_call", args, hexutil.EncodeUint64(lookup.block))
	if err != nil {
		return encoding.RpcFlakyFailure, true, nil, err
	}
	return encoding.NoPipelineError, false, b, nil
}

// doMercuryRequest sends requests to Mercury API to retrieve ChainlinkBlob.
func (r *EvmRegistry) doMercuryRequest(ctx context.Context, ml *FeedLookup, lggr logger.Logger) (encoding.PipelineExecutionState, [][]byte, bool, error) {
	var isMercuryV03 bool
	resultLen := len(ml.feeds)
	ch := make(chan MercuryData, resultLen)
	if ml.feedParamKey == feedIdHex && ml.timeParamKey == blockNumber {
		// only mercury v0.2
		for i := range ml.feeds {
			go r.singleFeedRequest(ctx, ch, i, ml, lggr)
		}
	} else if ml.feedParamKey == feedIDs && ml.timeParamKey == timestamp {
		// only mercury v0.3
		resultLen = 1
		isMercuryV03 = true
		ch = make(chan MercuryData, resultLen)
		go r.multiFeedsRequest(ctx, ch, ml, lggr)
	} else {
		// TODO: This should result in upkeep ineligible since upkeep gave a wrong mercury input
		return encoding.InvalidRevertDataInput, nil, false, fmt.Errorf("invalid label combination: feed param key %s and time param key %s", ml.feedParamKey, ml.timeParamKey)
	}

	var reqErr error
	results := make([][]byte, len(ml.feeds))
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
		if isMercuryV03 {
			results = m.Bytes
		} else {
			results[m.Index] = m.Bytes[0]
		}
	}
	lggr.Debugf("upkeep %s retryable %s reqErr %w", ml.upkeepId.String(), retryable && !allSuccess, reqErr)
	// only retry when not all successful AND none are not retryable
	return state, results, retryable && !allSuccess, reqErr
}

// singleFeedRequest sends a v0.2 Mercury request for a single feed report.
func (r *EvmRegistry) singleFeedRequest(ctx context.Context, ch chan<- MercuryData, index int, ml *FeedLookup, lggr logger.Logger) {
	q := url.Values{
		ml.feedParamKey: {ml.feeds[index]},
		ml.timeParamKey: {ml.time.String()},
	}
	mercuryURL := r.mercury.cred.URL
	reqUrl := fmt.Sprintf("%s%s%s", mercuryURL, mercuryPathV02, q.Encode())
	lggr.Debugf("request URL: %s", reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		ch <- MercuryData{Index: index, Error: err, Retryable: false, State: encoding.InvalidMercuryRequest}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := r.generateHMAC(http.MethodGet, mercuryPathV02+q.Encode(), []byte{}, r.mercury.cred.Username, r.mercury.cred.Password, ts)
	req.Header.Set(headerContentType, applicationJson)
	req.Header.Set(headerAuthorization, r.mercury.cred.Username)
	req.Header.Set(headerTimestamp, strconv.FormatInt(ts, 10))
	req.Header.Set(headerSignature, signature)

	// in the case of multiple retries here, use the last attempt's data
	state := encoding.NoPipelineError
	retryable := false
	sent := false
	retryErr := retry.Do(
		func() error {
			retryable = false
			resp, err1 := r.hc.Do(req)
			if err1 != nil {
				lggr.Warnf("upkeep %s block %s GET request fails for feed %s: %v", ml.upkeepId.String(), ml.time.String(), ml.feeds[index], err1)
				retryable = true
				state = encoding.MercuryFlakyFailure
				return err1
			}
			defer resp.Body.Close()
			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				retryable = false
				state = encoding.FailedToDecodeMercuryResponse
				return err1
			}

			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
				lggr.Warnf("upkeep %s block %s received status code %d for feed %s", ml.upkeepId.String(), ml.time.String(), resp.StatusCode, ml.feeds[index])
				retryable = true
				state = encoding.MercuryFlakyFailure
				return errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
			} else if resp.StatusCode != http.StatusOK {
				retryable = false
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("FeedLookup upkeep %s block %s received status code %d for feed %s", ml.upkeepId.String(), ml.time.String(), resp.StatusCode, ml.feeds[index])
			}

			var m MercuryV02Response
			err1 = json.Unmarshal(body, &m)
			if err1 != nil {
				lggr.Warnf("upkeep %s block %s failed to unmarshal body to MercuryV02Response for feed %s: %v", ml.upkeepId.String(), ml.time.String(), ml.feeds[index], err1)
				retryable = false
				state = encoding.MercuryUnmarshalError
				return err1
			}
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				lggr.Warnf("upkeep %s block %s failed to decode chainlinkBlob %s for feed %s: %v", ml.upkeepId.String(), ml.time.String(), m.ChainlinkBlob, ml.feeds[index], err1)
				retryable = false
				state = encoding.FailedToDecodeMercuryResponse
				return err1
			}
			ch <- MercuryData{
				Index:     index,
				Bytes:     [][]byte{blobBytes},
				Retryable: false,
				State:     encoding.NoPipelineError,
			}
			sent = true
			return nil
		},
		// only retry when the error is 404 Not Found or 500 Internal Server Error
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound) || err.Error() == fmt.Sprintf("%d", http.StatusInternalServerError)
		}),
		retry.Context(ctx),
		retry.Delay(retryDelay),
		retry.Attempts(totalAttempt))

	if !sent {
		md := MercuryData{
			Index:     index,
			Bytes:     [][]byte{},
			Retryable: retryable,
			Error:     retryErr,
			State:     state,
		}
		ch <- md
	}
}

// multiFeedsRequest sends a Mercury v0.3 request for a multi-feed report
func (r *EvmRegistry) multiFeedsRequest(ctx context.Context, ch chan<- MercuryData, ml *FeedLookup, lggr logger.Logger) {
	q := url.Values{
		feedIDs:   {strings.Join(ml.feeds, ",")},
		timestamp: {ml.time.String()},
	}

	reqUrl := fmt.Sprintf("%s%s%s", r.mercury.cred.URL, mercuryBatchPathV03, q.Encode())
	lggr.Debugf("request URL: %s", reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		ch <- MercuryData{Index: 0, Error: err, Retryable: false, State: encoding.InvalidMercuryRequest}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := r.generateHMAC(http.MethodGet, mercuryBatchPathV03+q.Encode(), []byte{}, r.mercury.cred.Username, r.mercury.cred.Password, ts)
	req.Header.Set(headerContentType, applicationJson)
	// username here is often referred to as user id
	req.Header.Set(headerAuthorization, r.mercury.cred.Username)
	req.Header.Set(headerTimestamp, strconv.FormatInt(ts, 10))
	req.Header.Set(headerSignature, signature)
	// mercury will inspect authorization headers above to make sure this user (in automation's context, this node) is eligible to access mercury
	// and if it has an automation role. it will then look at this upkeep id to check if it has access to all the requested feeds.
	req.Header.Set(headerUpkeepId, ml.upkeepId.String())

	// in the case of multiple retries here, use the last attempt's data
	state := encoding.NoPipelineError
	retryable := false
	sent := false
	retryErr := retry.Do(
		func() error {
			retryable = false
			resp, err1 := r.hc.Do(req)
			if err1 != nil {
				lggr.Warnf("upkeep %s block %s GET request fails for multi feed: %v", ml.upkeepId.String(), ml.time.String(), err1)
				retryable = true
				state = encoding.MercuryFlakyFailure
				return err1
			}
			defer resp.Body.Close()
			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				retryable = false
				state = encoding.FailedToDecodeMercuryResponse
				return err1
			}

			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
				lggr.Warnf("upkeep %s block %s received status code %d for multi feed", ml.upkeepId.String(), ml.time.String(), resp.StatusCode)
				retryable = true
				state = encoding.MercuryFlakyFailure
				return errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
			} else if resp.StatusCode != http.StatusOK {
				retryable = false
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("FeedLookup upkeep %s block %s received status code %d for multi feed", ml.upkeepId.String(), ml.time.String(), resp.StatusCode)
			}

			var responses []MercuryV03Response
			err1 = json.Unmarshal(body, &responses)
			if err1 != nil {
				lggr.Warnf("upkeep %s block %s failed to unmarshal body to MercuryV03Response for multi feed: %v", ml.upkeepId.String(), ml.time.String(), err1)
				retryable = false
				state = encoding.MercuryUnmarshalError
				return err1
			}
			var reportBytes [][]byte
			var b []byte
			for _, rsp := range responses {
				b, err1 = hexutil.Decode(rsp.FullReport)
				if err1 != nil {
					lggr.Warnf("upkeep %s block %s failed to decode reportBlob %s for multi feed: %v", ml.upkeepId.String(), ml.time.String(), rsp.FullReport, err1)
					retryable = false
					state = encoding.FailedToDecodeMercuryResponse
					return err1
				}
				reportBytes = append(reportBytes, b)
			}
			ch <- MercuryData{
				Index:     0,
				Bytes:     reportBytes,
				Retryable: false,
				State:     encoding.NoPipelineError,
			}
			sent = true
			return nil
		},
		// only retry when the error is 404 Not Found or 500 Internal Server Error
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound) || err.Error() == fmt.Sprintf("%d", http.StatusInternalServerError)
		}),
		retry.Context(ctx),
		retry.Delay(retryDelay),
		retry.Attempts(totalAttempt))

	if !sent {
		md := MercuryData{
			Index:     0,
			Bytes:     [][]byte{},
			Retryable: retryable,
			Error:     retryErr,
			State:     state,
		}
		ch <- md
	}
}

// generateHMAC calculates a user HMAC for Mercury server authentication.
func (r *EvmRegistry) generateHMAC(method string, path string, body []byte, clientId string, secret string, ts int64) string {
	bodyHash := sha256.New()
	bodyHash.Write(body)
	hashString := fmt.Sprintf("%s %s %s %s %d",
		method,
		path,
		hex.EncodeToString(bodyHash.Sum(nil)),
		clientId,
		ts)
	signedMessage := hmac.New(sha256.New, []byte(secret))
	signedMessage.Write([]byte(hashString))
	userHmac := hex.EncodeToString(signedMessage.Sum(nil))
	return userHmac
}
