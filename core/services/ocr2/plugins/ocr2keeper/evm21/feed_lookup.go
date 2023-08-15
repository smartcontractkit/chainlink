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
)

const (
	BlockNumber        = "blockNumber" // valid for v0.2
	FeedId             = "feedId"      // valid for v0.3
	FeedIdHex          = "feedIdHex"   // valid for v0.2
	MercuryPathV2      = "/client?"
	MercuryPathV3      = "/v1/reports?"
	MercuryBatchPathV3 = "/v1/reports/bulk?"
	RetryDelay         = 500 * time.Millisecond
	Timestamp          = "timestamp" // valid for v0.3
	TotalAttempt       = 3
	UserId             = "userId"
	MercuryV02         = MercuryVersion("v0.2")
	MercuryV03         = MercuryVersion("v0.3")
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

// MercuryResponse is used in both single feed endpoint and bulk endpoint because bulk endpoint will return ONE
// chainlinkBlob which contains multiple reports instead of multiple blobs.
type MercuryResponse struct {
	ChainlinkBlob string `json:"chainlinkBlob"`
}

type MercuryData struct {
	Index     int
	Error     error
	Retryable bool
	Bytes     []byte
	State     PipelineExecutionState
}

// AdminOffchainConfig represents the administrative offchain config for each upkeep. It can be set by s_upkeepManager
// role on the registry. Upkeeps allowed to use Mercury server will have this set to true.
type AdminOffchainConfig struct {
	MercuryEnabled bool `json:"mercuryEnabled"`
}

// feedLookup looks through check upkeep results looking for any that need off chain lookup
func (r *EvmRegistry) feedLookup(ctx context.Context, checkResults []ocr2keepers.CheckResult) []ocr2keepers.CheckResult {
	lookups := map[int]*FeedLookup{}
	for i, res := range checkResults {
		if res.IneligibilityReason != uint8(UpkeepFailureReasonTargetCheckReverted) {
			continue
		}

		block := big.NewInt(int64(res.Trigger.BlockNumber))
		upkeepId := res.UpkeepID

		opts := r.buildCallOpts(ctx, block)

		state, retryable, allowed, err := r.allowedToUseMercury(opts, upkeepId.BigInt())
		if err != nil {
			r.lggr.Warnf("[FeedLookup] upkeep %s block %d failed to query mercury allow list: %s", upkeepId, block, err)
			checkResults[i].PipelineExecutionState = uint8(state)
			checkResults[i].Retryable = retryable
			continue
		}

		if !allowed {
			r.lggr.Warnf("[FeedLookup] upkeep %s block %d NOT allowed to query Mercury server", upkeepId, block)
			checkResults[i].IneligibilityReason = uint8(UpkeepFailureReasonMercuryAccessNotAllowed)
			checkResults[i].Retryable = retryable
			continue
		}

		r.lggr.Infof("[FeedLookup] upkeep %s block %d decodeFeedLookup performData=%s", upkeepId, block, hexutil.Encode(checkResults[i].PerformData))
		state, lookup, err := r.decodeFeedLookup(res.PerformData)
		if err != nil {
			r.lggr.Warnf("[FeedLookup] upkeep %s block %d decodeFeedLookup: %v", upkeepId, block, err)
			checkResults[i].PipelineExecutionState = uint8(state)
			continue
		}
		lookup.upkeepId = upkeepId.BigInt()
		// the block here is exclusively used to call checkCallback at this block, not to be confused with the block number
		// in the revert for mercury v0.2, which is denoted by time in the struct bc starting from v0.3, only timestamp will be supported
		lookup.block = uint64(block.Int64())
		r.lggr.Infof("[FeedLookup] upkeep %s block %d decodeFeedLookup feedKey=%s timeKey=%s feeds=%v time=%s extraData=%s", upkeepId, block, lookup.feedParamKey, lookup.timeParamKey, lookup.feeds, lookup.time, hexutil.Encode(lookup.extraData))
		lookups[i] = lookup
	}

	var wg sync.WaitGroup
	for i, lookup := range lookups {
		wg.Add(1)
		go r.doLookup(ctx, &wg, lookup, i, checkResults)
	}
	wg.Wait()

	// don't surface error to plugin bc FeedLookup process should be self-contained.
	return checkResults
}

func (r *EvmRegistry) doLookup(ctx context.Context, wg *sync.WaitGroup, lookup *FeedLookup, i int, checkResults []ocr2keepers.CheckResult) {
	defer wg.Done()

	state, values, retryable, err := r.doMercuryRequest(ctx, lookup)
	if err != nil {
		r.lggr.Errorf("[FeedLookup] upkeep %s retryable %v doMercuryRequest: %v", lookup.upkeepId, retryable, err)
		checkResults[i].Retryable = retryable
		checkResults[i].PipelineExecutionState = uint8(state)
		return
	}
	for j, v := range values {
		r.lggr.Infof("[FeedLookup] checkCallback values[%d]=%s", j, hexutil.Encode(v))
	}

	state, retryable, mercuryBytes, err := r.checkCallback(ctx, values, lookup)
	if err != nil {
		r.lggr.Errorf("[FeedLookup] upkeep %s block %d checkCallback err: %v", lookup.upkeepId, lookup.block, err)
		checkResults[i].Retryable = retryable
		checkResults[i].PipelineExecutionState = uint8(state)
		return
	}
	r.lggr.Infof("[FeedLookup] checkCallback mercuryBytes=%s", hexutil.Encode(mercuryBytes))

	state, needed, performData, failureReason, _, err := r.packer.UnpackCheckCallbackResult(mercuryBytes)
	if err != nil {
		r.lggr.Errorf("[FeedLookup] upkeep %s block %d UnpackCheckCallbackResult err: %v", lookup.upkeepId, lookup.block, err)
		checkResults[i].PipelineExecutionState = uint8(state)
		return
	}

	if failureReason == uint8(UpkeepFailureReasonMercuryCallbackReverted) {
		checkResults[i].IneligibilityReason = uint8(UpkeepFailureReasonMercuryCallbackReverted)
		r.lggr.Debugf("[FeedLookup] upkeep %s block %d mercury callback reverts", lookup.upkeepId, lookup.block)
		return
	}

	if !needed {
		checkResults[i].IneligibilityReason = uint8(UpkeepFailureReasonUpkeepNotNeeded)
		r.lggr.Debugf("[FeedLookup] upkeep %s block %d callback reports upkeep not needed", lookup.upkeepId, lookup.block)
		return
	}

	checkResults[i].IneligibilityReason = uint8(UpkeepFailureReasonNone)
	checkResults[i].Eligible = true
	checkResults[i].PerformData = performData
	r.lggr.Infof("[FeedLookup] upkeep %s block %d successful with perform data: %s", lookup.upkeepId, lookup.block, hexutil.Encode(performData))
}

// allowedToUseMercury retrieves upkeep's administrative offchain config and decode a mercuryEnabled bool to indicate if
// this upkeep is allowed to use Mercury service.
func (r *EvmRegistry) allowedToUseMercury(opts *bind.CallOpts, upkeepId *big.Int) (state PipelineExecutionState, retryable bool, allow bool, err error) {
	allowed, ok := r.mercury.allowListCache.Get(upkeepId.String())
	if ok {
		return NoPipelineError, false, allowed.(bool), nil
	}

	cfg, err := r.registry.GetUpkeepPrivilegeConfig(opts, upkeepId)
	if err != nil {
		return RpcFlakyFailure, true, false, fmt.Errorf("failed to get upkeep privilege config for upkeep ID %s: %v", upkeepId, err)
	}

	var a AdminOffchainConfig
	err = json.Unmarshal(cfg, &a)
	if err != nil {
		return MercuryUnmarshalError, false, false, fmt.Errorf("failed to unmarshal privilege config for upkeep ID %s: %v", upkeepId, err)
	}
	r.mercury.allowListCache.Set(upkeepId.String(), a.MercuryEnabled, cache.DefaultExpiration)
	return NoPipelineError, false, a.MercuryEnabled, nil
}

// decodeFeedLookup decodes the revert error FeedLookup(string feedParamKey, string[] feeds, string feedParamKey, uint256 time, byte[] extraData)
func (r *EvmRegistry) decodeFeedLookup(data []byte) (PipelineExecutionState, *FeedLookup, error) {
	e := r.mercury.abi.Errors["FeedLookup"]
	unpack, err := e.Unpack(data)
	if err != nil {
		return PackUnpackDecodeFailed, nil, fmt.Errorf("unpack error: %w", err)
	}
	errorParameters := unpack.([]interface{})

	return NoPipelineError, &FeedLookup{
		feedParamKey: *abi.ConvertType(errorParameters[0], new(string)).(*string),
		feeds:        *abi.ConvertType(errorParameters[1], new([]string)).(*[]string),
		timeParamKey: *abi.ConvertType(errorParameters[2], new(string)).(*string),
		time:         *abi.ConvertType(errorParameters[3], new(*big.Int)).(**big.Int),
		extraData:    *abi.ConvertType(errorParameters[4], new([]byte)).(*[]byte),
	}, nil
}

func (r *EvmRegistry) checkCallback(ctx context.Context, values [][]byte, lookup *FeedLookup) (PipelineExecutionState, bool, hexutil.Bytes, error) {
	payload, err := r.abi.Pack("checkCallback", lookup.upkeepId, values, lookup.extraData)
	if err != nil {
		return PackUnpackDecodeFailed, false, nil, err
	}

	var b hexutil.Bytes
	args := map[string]interface{}{
		"to":   r.addr.Hex(),
		"data": hexutil.Bytes(payload),
	}

	// call checkCallback function at the block which OCR3 has agreed upon
	err = r.client.CallContext(ctx, &b, "eth_call", args, hexutil.EncodeUint64(lookup.block))
	if err != nil {
		return RpcFlakyFailure, true, nil, err
	}
	return NoPipelineError, false, b, nil
}

// doMercuryRequest sends requests to Mercury API to retrieve ChainlinkBlob.
func (r *EvmRegistry) doMercuryRequest(ctx context.Context, ml *FeedLookup) (PipelineExecutionState, [][]byte, bool, error) {
	resultLen := len(ml.feeds)
	ch := make(chan MercuryData, resultLen)
	if ml.feedParamKey == FeedIdHex && ml.timeParamKey == BlockNumber {
		// only mercury v0.2
		for i := range ml.feeds {
			go r.singleFeedRequest(ctx, ch, i, ml, MercuryV02)
		}
	} else if ml.feedParamKey == FeedId && ml.timeParamKey == Timestamp {
		// only mercury v0.3
		if resultLen == 1 {
			go r.singleFeedRequest(ctx, ch, 0, ml, MercuryV03)
		} else {
			// create a new channel with buffer size 1 since the batch endpoint will only return 1 blob
			resultLen = 1
			ch = make(chan MercuryData, resultLen)
			go r.multiFeedsRequest(ctx, ch, ml)
		}
	} else {
		return InvalidRevertDataInput, nil, false, fmt.Errorf("invalid label combination: feed param key %s and time param key %s", ml.feedParamKey, ml.timeParamKey)
	}

	var reqErr error
	results := make([][]byte, len(ml.feeds))
	retryable := true
	allSuccess := true
	// use the last execution error as the state, if no execution errors, state will be no error
	state := NoPipelineError
	for i := 0; i < len(results); i++ {
		m := <-ch
		if m.Error != nil {
			reqErr = errors.Join(reqErr, m.Error)
			retryable = retryable && m.Retryable
			allSuccess = false
			if m.State != NoPipelineError {
				state = m.State
			}
		}
		results[m.Index] = m.Bytes
	}
	r.lggr.Debugf("FeedLookup upkeep %s retryable %s reqErr %w", ml.upkeepId.String(), retryable && !allSuccess, reqErr)
	// only retry when not all successful AND none are not retryable
	return state, results, retryable && !allSuccess, reqErr
}

// singleFeedRequest sends a Mercury request for a single feed report.
func (r *EvmRegistry) singleFeedRequest(ctx context.Context, ch chan<- MercuryData, index int, ml *FeedLookup, mv MercuryVersion) {
	q := url.Values{
		ml.feedParamKey: {ml.feeds[index]},
		ml.timeParamKey: {ml.time.String()},
	}
	mercuryURL := r.mercury.cred.URL
	path := MercuryPathV2
	if mv == MercuryV03 {
		path = MercuryPathV3
	}
	reqUrl := fmt.Sprintf("%s%s%s", mercuryURL, path, q.Encode())
	r.lggr.Debugf("FeedLookup request URL: %s", reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		ch <- MercuryData{Index: index, Error: err, Retryable: false, State: InvalidMercuryRequest}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := r.generateHMAC(http.MethodGet, path+q.Encode(), []byte{}, r.mercury.cred.Username, r.mercury.cred.Password, ts)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.mercury.cred.Username)
	req.Header.Set("X-Authorization-Timestamp", strconv.FormatInt(ts, 10))
	req.Header.Set("X-Authorization-Signature-SHA256", signature)

	// in the case of multiple retries here, use the last attempt's data
	state := NoPipelineError
	retryable := false
	sent := false
	retryErr := retry.Do(
		func() error {
			retryable = false
			resp, err1 := r.hc.Do(req)
			if err1 != nil {
				r.lggr.Warnf("FeedLookup upkeep %s block %s GET request fails for feed %s: %v", ml.upkeepId.String(), ml.time.String(), ml.feeds[index], err1)
				retryable = true
				state = MercuryFlakyFailure
				return err1
			}
			defer resp.Body.Close()
			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				retryable = false
				state = FailedToReadMercuryResponse
				return err1
			}

			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
				r.lggr.Warnf("FeedLookup upkeep %s block %s received status code %d for feed %s", ml.upkeepId.String(), ml.time.String(), resp.StatusCode, ml.feeds[index])
				retryable = true
				state = MercuryFlakyFailure
				return errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
			} else if resp.StatusCode != http.StatusOK {
				retryable = false
				state = InvalidMercuryRequest
				return fmt.Errorf("FeedLookup upkeep %s block %s received status code %d for feed %s", ml.upkeepId.String(), ml.time.String(), resp.StatusCode, ml.feeds[index])
			}

			var m MercuryResponse
			err1 = json.Unmarshal(body, &m)
			if err1 != nil {
				r.lggr.Warnf("FeedLookup upkeep %s block %s failed to unmarshal body to MercuryResponse for feed %s: %v", ml.upkeepId.String(), ml.time.String(), ml.feeds[index], err1)
				retryable = false
				state = MercuryUnmarshalError
				return err1
			}
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				r.lggr.Warnf("FeedLookup upkeep %s block %s failed to decode chainlinkBlob %s for feed %s: %v", ml.upkeepId.String(), ml.time.String(), m.ChainlinkBlob, ml.feeds[index], err1)
				retryable = false
				state = FailedToReadMercuryResponse
				return err1
			}
			ch <- MercuryData{
				Index:     index,
				Bytes:     blobBytes,
				Retryable: false,
				State:     NoPipelineError,
			}
			sent = true
			return nil
		},
		// only retry when the error is 404 Not Found or 500 Internal Server Error
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound) || err.Error() == fmt.Sprintf("%d", http.StatusInternalServerError)
		}),
		retry.Context(ctx),
		retry.Delay(RetryDelay),
		retry.Attempts(TotalAttempt))

	if !sent {
		md := MercuryData{
			Index:     index,
			Retryable: retryable,
			Error:     retryErr,
			State:     state,
		}
		ch <- md
	}
}

// multiFeedsRequest sends a Mercury request for a multi-feed report
func (r *EvmRegistry) multiFeedsRequest(ctx context.Context, ch chan<- MercuryData, ml *FeedLookup) {
	q := url.Values{
		FeedId:    {strings.Join(ml.feeds, ",")},
		Timestamp: {ml.time.String()},
	}

	reqUrl := fmt.Sprintf("%s%s%s", r.mercury.cred.URL, MercuryBatchPathV3, q.Encode())
	r.lggr.Debugf("FeedLookup request URL: %s", reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		ch <- MercuryData{Index: 0, Error: err, Retryable: false, State: InvalidMercuryRequest}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := r.generateHMAC(http.MethodGet, MercuryBatchPathV3+q.Encode(), []byte{}, r.mercury.cred.Username, r.mercury.cred.Password, ts)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.mercury.cred.Username)
	req.Header.Set("X-Authorization-Timestamp", strconv.FormatInt(ts, 10))
	req.Header.Set("X-Authorization-Signature-SHA256", signature)

	// in the case of multiple retries here, use the last attempt's data
	state := NoPipelineError
	retryable := false
	sent := false
	retryErr := retry.Do(
		func() error {
			retryable = false
			resp, err1 := r.hc.Do(req)
			if err1 != nil {
				r.lggr.Warnf("FeedLookup upkeep %s block %s GET request fails for multi feed: %v", ml.upkeepId.String(), ml.time.String(), err1)
				retryable = true
				state = MercuryFlakyFailure
				return err1
			}
			defer resp.Body.Close()
			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				retryable = false
				state = FailedToReadMercuryResponse
				return err1
			}

			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
				r.lggr.Warnf("FeedLookup upkeep %s block %s received status code %d for multi feed", ml.upkeepId.String(), ml.time.String(), resp.StatusCode)
				retryable = true
				state = MercuryFlakyFailure
				return errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
			} else if resp.StatusCode != http.StatusOK {
				retryable = false
				state = InvalidMercuryRequest
				return fmt.Errorf("FeedLookup upkeep %s block %s received status code %d for multi feed", ml.upkeepId.String(), ml.time.String(), resp.StatusCode)
			}

			var m MercuryResponse
			err1 = json.Unmarshal(body, &m)
			if err1 != nil {
				r.lggr.Warnf("FeedLookup upkeep %s block %s failed to unmarshal body to MercuryResponse for multi feed: %v", ml.upkeepId.String(), ml.time.String(), err1)
				retryable = false
				state = MercuryUnmarshalError
				return err1
			}
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				r.lggr.Warnf("FeedLookup upkeep %s block %s failed to decode chainlinkBlob %s for multi feed: %v", ml.upkeepId.String(), ml.time.String(), m.ChainlinkBlob, err1)
				retryable = false
				state = FailedToReadMercuryResponse
				return err1
			}
			ch <- MercuryData{
				Index:     0,
				Bytes:     blobBytes,
				Retryable: false,
				State:     NoPipelineError,
			}
			sent = true
			return nil
		},
		// only retry when the error is 404 Not Found or 500 Internal Server Error
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound) || err.Error() == fmt.Sprintf("%d", http.StatusInternalServerError)
		}),
		retry.Context(ctx),
		retry.Delay(RetryDelay),
		retry.Attempts(TotalAttempt))

	if !sent {
		md := MercuryData{
			Index:     0,
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
