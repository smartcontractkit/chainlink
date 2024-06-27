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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/encoding"
)

const (
	applicationJson                = "application/json"
	blockNumber                    = "blockNumber" // valid for v0.2
	feedIDs                        = "feedIDs"     // valid for v0.3
	feedIdHex                      = "feedIdHex"   // valid for v0.2
	headerAuthorization            = "Authorization"
	headerContentType              = "Content-Type"
	headerTimestamp                = "X-Authorization-Timestamp"
	headerSignature                = "X-Authorization-Signature-SHA256"
	headerUpkeepId                 = "X-Authorization-Upkeep-Id"
	mercuryPathV02                 = "/client?"                 // only used to access mercury v0.2 server
	mercuryBatchPathV03            = "/api/v1/reports/bulk?"    // only used to access mercury v0.3 server
	mercuryBatchPathV03BlockNumber = "/api/v1gmx/reports/bulk?" // only used to access mercury v0.3 server with blockNumber
	retryDelay                     = 500 * time.Millisecond
	timestamp                      = "timestamp" // valid for v0.3
	totalAttempt                   = 3
)

type StreamsLookup struct {
	*encoding.StreamsLookupError
	upkeepId *big.Int
	block    uint64
}

// MercuryV02Response represents a JSON structure used by Mercury v0.2
type MercuryV02Response struct {
	ChainlinkBlob string `json:"chainlinkBlob"`
}

// MercuryV03Response represents a JSON structure used by Mercury v0.3
type MercuryV03Response struct {
	Reports []MercuryV03Report `json:"reports"`
}

type MercuryV03Report struct {
	FeedID                string `json:"feedID"` // feed id in hex encoded
	ValidFromTimestamp    uint32 `json:"validFromTimestamp"`
	ObservationsTimestamp uint32 `json:"observationsTimestamp"`
	FullReport            string `json:"fullReport"` // the actual hex encoded mercury report of this feed, can be sent to verifier
}

type MercuryData struct {
	Index     int
	Error     error
	Retryable bool
	Bytes     [][]byte
	State     encoding.PipelineExecutionState
}

// UpkeepPrivilegeConfig represents the administrative offchain config for each upkeep. It can be set by s_upkeepPrivilegeManager
// role on the registry. Upkeeps allowed to use Mercury server will have this set to true.
type UpkeepPrivilegeConfig struct {
	MercuryEnabled bool `json:"mercuryEnabled"`
}

// streamsLookup looks through check upkeep results looking for any that need off chain lookup
func (r *EvmRegistry) streamsLookup(ctx context.Context, checkResults []ocr2keepers.CheckResult) []ocr2keepers.CheckResult {
	lggr := r.lggr.With("where", "StreamsLookup")
	lookups := map[int]*StreamsLookup{}
	for i, res := range checkResults {
		if res.IneligibilityReason != uint8(encoding.UpkeepFailureReasonTargetCheckReverted) {
			// Streams Lookup only works when upkeep target check reverts
			continue
		}

		block := big.NewInt(int64(res.Trigger.BlockNumber))
		upkeepId := res.UpkeepID

		// Try to decode the revert error into streams lookup format. User upkeeps can revert with any reason, see if they
		// tried to call mercury
		lggr.Infof("at block %d upkeep %s trying to DecodeStreamsLookupRequest performData=%s", block, upkeepId, hexutil.Encode(checkResults[i].PerformData))
		streamsLookupErr, err := r.packer.DecodeStreamsLookupRequest(res.PerformData)
		if err != nil {
			lggr.Debugf("at block %d upkeep %s DecodeStreamsLookupRequest failed: %v", block, upkeepId, err)
			// user contract did not revert with StreamsLookup error
			continue
		}
		l := &StreamsLookup{StreamsLookupError: streamsLookupErr}
		if r.mercury.cred == nil {
			lggr.Errorf("at block %d upkeep %s tries to access mercury server but mercury credential is not configured", block, upkeepId)
			continue
		}

		if len(l.Feeds) == 0 {
			checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonInvalidRevertDataInput)
			lggr.Debugf("at block %s upkeep %s has empty feeds array", block, upkeepId)
			continue
		}
		// mercury permission checking for v0.3 is done by mercury server
		if l.FeedParamKey == feedIdHex && l.TimeParamKey == blockNumber {
			// check permission on the registry for mercury v0.2
			opts := r.buildCallOpts(ctx, block)
			state, reason, retryable, allowed, err := r.allowedToUseMercury(opts, upkeepId.BigInt())
			if err != nil {
				lggr.Warnf("at block %s upkeep %s failed to query mercury allow list: %s", block, upkeepId, err)
				checkResults[i].PipelineExecutionState = uint8(state)
				checkResults[i].IneligibilityReason = uint8(reason)
				checkResults[i].Retryable = retryable
				continue
			}

			if !allowed {
				lggr.Debugf("at block %d upkeep %s NOT allowed to query Mercury server", block, upkeepId)
				checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonMercuryAccessNotAllowed)
				continue
			}
		} else if l.FeedParamKey != feedIDs {
			// if mercury version cannot be determined, set failure reason
			lggr.Debugf("at block %d upkeep %s NOT allowed to query Mercury server", block, upkeepId)
			checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonInvalidRevertDataInput)
			continue
		}

		l.upkeepId = upkeepId.BigInt()
		// the block here is exclusively used to call checkCallback at this block, not to be confused with the block number
		// in the revert for mercury v0.2, which is denoted by time in the struct bc starting from v0.3, only timestamp will be supported
		l.block = uint64(block.Int64())
		lggr.Infof("at block %d upkeep %s DecodeStreamsLookupRequest feedKey=%s timeKey=%s feeds=%v time=%s extraData=%s", block, upkeepId, l.FeedParamKey, l.TimeParamKey, l.Feeds, l.Time, hexutil.Encode(l.ExtraData))
		lookups[i] = l
	}

	var wg sync.WaitGroup
	for i, lookup := range lookups {
		wg.Add(1)
		go r.doLookup(ctx, &wg, lookup, i, checkResults, lggr)
	}
	wg.Wait()

	// don't surface error to plugin bc StreamsLookup process should be self-contained.
	return checkResults
}

func (r *EvmRegistry) doLookup(ctx context.Context, wg *sync.WaitGroup, lookup *StreamsLookup, i int, checkResults []ocr2keepers.CheckResult, lggr logger.Logger) {
	defer wg.Done()

	state, reason, values, retryable, err := r.doMercuryRequest(ctx, lookup, lggr)
	if err != nil {
		lggr.Errorf("upkeep %s retryable %v doMercuryRequest: %s", lookup.upkeepId, retryable, err.Error())
		checkResults[i].Retryable = retryable
		checkResults[i].PipelineExecutionState = uint8(state)
		checkResults[i].IneligibilityReason = uint8(reason)
		return
	}

	for j, v := range values {
		lggr.Infof("upkeep %s doMercuryRequest values[%d]: %s", lookup.upkeepId, j, hexutil.Encode(v))
	}

	state, retryable, mercuryBytes, err := r.checkCallback(ctx, values, lookup)
	if err != nil {
		lggr.Errorf("at block %d upkeep %s checkCallback err: %s", lookup.block, lookup.upkeepId, err.Error())
		checkResults[i].Retryable = retryable
		checkResults[i].PipelineExecutionState = uint8(state)
		return
	}
	lggr.Infof("checkCallback mercuryBytes=%s", hexutil.Encode(mercuryBytes))

	state, needed, performData, failureReason, _, err := r.packer.UnpackCheckCallbackResult(mercuryBytes)
	if err != nil {
		lggr.Errorf("at block %d upkeep %s UnpackCheckCallbackResult err: %s", lookup.block, lookup.upkeepId, err.Error())
		checkResults[i].PipelineExecutionState = uint8(state)
		return
	}

	if failureReason == uint8(encoding.UpkeepFailureReasonMercuryCallbackReverted) {
		checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonMercuryCallbackReverted)
		lggr.Debugf("at block %d upkeep %s mercury callback reverts", lookup.block, lookup.upkeepId)
		return
	}

	if !needed {
		checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonUpkeepNotNeeded)
		lggr.Debugf("at block %d upkeep %s callback reports upkeep not needed", lookup.block, lookup.upkeepId)
		return
	}

	checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonNone)
	checkResults[i].Eligible = true
	checkResults[i].PerformData = performData
	lggr.Infof("at block %d upkeep %s successful with perform data: %s", lookup.block, lookup.upkeepId, hexutil.Encode(performData))
}

// allowedToUseMercury retrieves upkeep's administrative offchain config and decode a mercuryEnabled bool to indicate if
// this upkeep is allowed to use Mercury service.
func (r *EvmRegistry) allowedToUseMercury(opts *bind.CallOpts, upkeepId *big.Int) (state encoding.PipelineExecutionState, reason encoding.UpkeepFailureReason, retryable bool, allow bool, err error) {
	allowed, ok := r.mercury.allowListCache.Get(upkeepId.String())
	if ok {
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonNone, false, allowed.(bool), nil
	}

	payload, err := r.packer.PackGetUpkeepPrivilegeConfig(upkeepId)
	if err != nil {
		// pack error, no retryable
		r.lggr.Warnf("failed to pack getUpkeepPrivilegeConfig data for upkeepId %s: %s", upkeepId, err)

		return encoding.PackUnpackDecodeFailed, encoding.UpkeepFailureReasonNone, false, false, fmt.Errorf("failed to pack upkeepId: %w", err)
	}

	var resultBytes hexutil.Bytes
	args := map[string]interface{}{
		"to":   r.addr.Hex(),
		"data": hexutil.Bytes(payload),
	}

	// call checkCallback function at the block which OCR3 has agreed upon
	err = r.client.CallContext(opts.Context, &resultBytes, "eth_call", args, hexutil.EncodeBig(opts.BlockNumber))
	if err != nil {
		return encoding.RpcFlakyFailure, encoding.UpkeepFailureReasonNone, true, false, fmt.Errorf("failed to get upkeep privilege config: %v", err)
	}

	cfg, err := r.packer.UnpackGetUpkeepPrivilegeConfig(resultBytes)
	if err != nil {
		return encoding.PackUnpackDecodeFailed, encoding.UpkeepFailureReasonNone, false, false, fmt.Errorf("failed to get upkeep privilege config: %v", err)
	}

	if len(cfg) == 0 {
		r.mercury.allowListCache.Set(upkeepId.String(), false, cache.DefaultExpiration)
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonMercuryAccessNotAllowed, false, false, fmt.Errorf("upkeep privilege config is empty")
	}

	var privilegeConfig UpkeepPrivilegeConfig
	if err := json.Unmarshal(cfg, &privilegeConfig); err != nil {
		return encoding.MercuryUnmarshalError, encoding.UpkeepFailureReasonNone, false, false, fmt.Errorf("failed to unmarshal privilege config: %v", err)
	}

	r.mercury.allowListCache.Set(upkeepId.String(), privilegeConfig.MercuryEnabled, cache.DefaultExpiration)

	return encoding.NoPipelineError, encoding.UpkeepFailureReasonNone, false, privilegeConfig.MercuryEnabled, nil
}

func (r *EvmRegistry) checkCallback(ctx context.Context, values [][]byte, lookup *StreamsLookup) (encoding.PipelineExecutionState, bool, hexutil.Bytes, error) {
	payload, err := r.abi.Pack("checkCallback", lookup.upkeepId, values, lookup.ExtraData)
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

// doMercuryRequest sends requests to Mercury API to retrieve mercury data.
func (r *EvmRegistry) doMercuryRequest(ctx context.Context, sl *StreamsLookup, lggr logger.Logger) (encoding.PipelineExecutionState, encoding.UpkeepFailureReason, [][]byte, bool, error) {
	var isMercuryV03 bool
	resultLen := len(sl.Feeds)
	ch := make(chan MercuryData, resultLen)
	if len(sl.Feeds) == 0 {
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonInvalidRevertDataInput, [][]byte{}, false, fmt.Errorf("invalid revert data input: feed param key %s, time param key %s, feeds %s", sl.FeedParamKey, sl.TimeParamKey, sl.Feeds)
	}
	if sl.FeedParamKey == feedIdHex && sl.TimeParamKey == blockNumber {
		// only mercury v0.2
		for i := range sl.Feeds {
			go r.singleFeedRequest(ctx, ch, i, sl, lggr)
		}
	} else if sl.FeedParamKey == feedIDs {
		// only mercury v0.3
		resultLen = 1
		isMercuryV03 = true
		ch = make(chan MercuryData, resultLen)
		go r.multiFeedsRequest(ctx, ch, sl, lggr)
	} else {
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonInvalidRevertDataInput, [][]byte{}, false, fmt.Errorf("invalid revert data input: feed param key %s, time param key %s, feeds %s", sl.FeedParamKey, sl.TimeParamKey, sl.Feeds)
	}

	var reqErr error
	results := make([][]byte, len(sl.Feeds))
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
	// only retry when not all successful AND none are not retryable
	return state, encoding.UpkeepFailureReasonNone, results, retryable && !allSuccess, reqErr
}

// singleFeedRequest sends a v0.2 Mercury request for a single feed report.
func (r *EvmRegistry) singleFeedRequest(ctx context.Context, ch chan<- MercuryData, index int, sl *StreamsLookup, lggr logger.Logger) {
	q := url.Values{
		sl.FeedParamKey: {sl.Feeds[index]},
		sl.TimeParamKey: {sl.Time.String()},
	}
	mercuryURL := r.mercury.cred.LegacyURL
	reqUrl := fmt.Sprintf("%s%s%s", mercuryURL, mercuryPathV02, q.Encode())
	lggr.Debugf("request URL for upkeep %s feed %s: %s", sl.upkeepId.String(), sl.Feeds[index], reqUrl)

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
				lggr.Warnf("at block %s upkeep %s GET request fails for feed %s: %v", sl.Time.String(), sl.upkeepId.String(), sl.Feeds[index], err1)
				retryable = true
				state = encoding.MercuryFlakyFailure
				return err1
			}
			defer func(Body io.ReadCloser) {
				err = Body.Close()
				if err != nil {
					lggr.Warnf("failed to close mercury response Body: %s", err)
				}
			}(resp.Body)

			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				retryable = false
				state = encoding.InvalidMercuryResponse
				return err1
			}

			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
				lggr.Warnf("at block %s upkeep %s received status code %d for feed %s", sl.Time.String(), sl.upkeepId.String(), resp.StatusCode, sl.Feeds[index])
				retryable = true
				state = encoding.MercuryFlakyFailure
				return errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
			} else if resp.StatusCode != http.StatusOK {
				retryable = false
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("at block %s upkeep %s received status code %d for feed %s", sl.Time.String(), sl.upkeepId.String(), resp.StatusCode, sl.Feeds[index])
			}

			lggr.Debugf("at block %s upkeep %s received status code %d from mercury v0.2 with BODY=%s", sl.Time.String(), sl.upkeepId.String(), resp.StatusCode, hexutil.Encode(body))

			var m MercuryV02Response
			err1 = json.Unmarshal(body, &m)
			if err1 != nil {
				lggr.Warnf("at block %s upkeep %s failed to unmarshal body to MercuryV02Response for feed %s: %v", sl.Time.String(), sl.upkeepId.String(), sl.Feeds[index], err1)
				retryable = false
				state = encoding.MercuryUnmarshalError
				return err1
			}
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				lggr.Warnf("at block %s upkeep %s failed to decode chainlinkBlob %s for feed %s: %v", sl.Time.String(), sl.upkeepId.String(), m.ChainlinkBlob, sl.Feeds[index], err1)
				retryable = false
				state = encoding.InvalidMercuryResponse
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
			Error:     fmt.Errorf("failed to request feed for %s: %w", sl.Feeds[index], retryErr),
			State:     state,
		}
		ch <- md
	}
}

// multiFeedsRequest sends a Mercury v0.3 request for a multi-feed report
func (r *EvmRegistry) multiFeedsRequest(ctx context.Context, ch chan<- MercuryData, sl *StreamsLookup, lggr logger.Logger) {
	// this won't work bc q.Encode() will encode commas as '%2C' but the server is strictly expecting a comma separated list
	//q := url.Values{
	//	feedIDs:   {strings.Join(sl.Feeds, ",")},
	//	timestamp: {sl.Time.String()},
	//}
	params := fmt.Sprintf("%s=%s&%s=%s", feedIDs, strings.Join(sl.Feeds, ","), sl.TimeParamKey, sl.Time.String())
	batchPathV03 := mercuryBatchPathV03
	if sl.TimeParamKey == blockNumber {
		batchPathV03 = mercuryBatchPathV03BlockNumber
	}
	reqUrl := fmt.Sprintf("%s%s%s", r.mercury.cred.URL, batchPathV03, params)
	lggr.Debugf("request URL for upkeep %s userId %s: %s", sl.upkeepId.String(), r.mercury.cred.Username, reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		ch <- MercuryData{Index: 0, Error: err, Retryable: false, State: encoding.InvalidMercuryRequest}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := r.generateHMAC(http.MethodGet, batchPathV03+params, []byte{}, r.mercury.cred.Username, r.mercury.cred.Password, ts)
	req.Header.Set(headerContentType, applicationJson)
	// username here is often referred to as user id
	req.Header.Set(headerAuthorization, r.mercury.cred.Username)
	req.Header.Set(headerTimestamp, strconv.FormatInt(ts, 10))
	req.Header.Set(headerSignature, signature)
	// mercury will inspect authorization headers above to make sure this user (in automation's context, this node) is eligible to access mercury
	// and if it has an automation role. it will then look at this upkeep id to check if it has access to all the requested feeds.
	req.Header.Set(headerUpkeepId, sl.upkeepId.String())

	// in the case of multiple retries here, use the last attempt's data
	state := encoding.NoPipelineError
	retryable := false
	sent := false
	retryErr := retry.Do(
		func() error {
			retryable = false
			resp, err1 := r.hc.Do(req)
			if err1 != nil {
				lggr.Warnf("at timestamp %s upkeep %s GET request fails from mercury v0.3: %v", sl.Time.String(), sl.upkeepId.String(), err1)
				retryable = true
				state = encoding.MercuryFlakyFailure
				return err1
			}
			defer func(Body io.ReadCloser) {
				err = Body.Close()
				if err != nil {
					lggr.Warnf("failed to close mercury response Body: %s", err)
				}
			}(resp.Body)

			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				retryable = false
				state = encoding.InvalidMercuryResponse
				return err1
			}

			lggr.Infof("at timestamp %s upkeep %s received status code %d from mercury v0.3", sl.Time.String(), sl.upkeepId.String(), resp.StatusCode)
			if resp.StatusCode == http.StatusUnauthorized {
				retryable = false
				state = encoding.UpkeepNotAuthorized
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3, most likely this is caused by unauthorized upkeep", sl.Time.String(), sl.upkeepId.String(), resp.StatusCode)
			} else if resp.StatusCode == http.StatusBadRequest {
				retryable = false
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3 with message: %s", sl.Time.String(), sl.upkeepId.String(), resp.StatusCode, string(body))
			} else if resp.StatusCode == http.StatusInternalServerError {
				retryable = true
				state = encoding.MercuryFlakyFailure
				return fmt.Errorf("%d", http.StatusInternalServerError)
			} else if resp.StatusCode == 420 {
				// in 0.3, this will happen when missing/malformed query args, missing or bad required headers, non-existent feeds, or no permissions for feeds
				retryable = false
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3, most likely this is caused by missing/malformed query args, missing or bad required headers, non-existent feeds, or no permissions for feeds", sl.Time.String(), sl.upkeepId.String(), resp.StatusCode)
			} else if resp.StatusCode != http.StatusOK {
				retryable = false
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3", sl.Time.String(), sl.upkeepId.String(), resp.StatusCode)
			}

			lggr.Debugf("at block %s upkeep %s received status code %d from mercury v0.3 with BODY=%s", sl.Time.String(), sl.upkeepId.String(), resp.StatusCode, hexutil.Encode(body))

			var response MercuryV03Response
			err1 = json.Unmarshal(body, &response)
			if err1 != nil {
				lggr.Warnf("at timestamp %s upkeep %s failed to unmarshal body to MercuryV03Response from mercury v0.3: %v", sl.Time.String(), sl.upkeepId.String(), err1)
				retryable = false
				state = encoding.MercuryUnmarshalError
				return err1
			}
			// in v0.3, if some feeds are not available, the server will only return available feeds, but we need to make sure ALL feeds are retrieved before calling user contract
			// hence, retry in this case. retry will help when we send a very new timestamp and reports are not yet generated
			if len(response.Reports) != len(sl.Feeds) {
				// TODO: AUTO-5044: calculate what reports are missing and log a warning
				lggr.Warnf("at timestamp %s upkeep %s mercury v0.3 server retruned 200 status with %d reports while we requested %d feeds, treating as 404 (not found) and retrying", sl.Time.String(), sl.upkeepId.String(), len(response.Reports), len(sl.Feeds))
				retryable = true
				state = encoding.MercuryFlakyFailure
				return fmt.Errorf("%d", http.StatusNotFound)
			}
			var reportBytes [][]byte
			for _, rsp := range response.Reports {
				b, err := hexutil.Decode(rsp.FullReport)
				if err != nil {
					lggr.Warnf("at timestamp %s upkeep %s failed to decode reportBlob %s: %v", sl.Time.String(), sl.upkeepId.String(), rsp.FullReport, err)
					retryable = false
					state = encoding.InvalidMercuryResponse
					return err
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
