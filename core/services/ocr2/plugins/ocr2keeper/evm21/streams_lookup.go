package evm

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
	mercuryPathV02      = "/client?"              // only used to access mercury v0.2 server
	mercuryBatchPathV03 = "/api/v1/reports/bulk?" // only used to access mercury v0.3 server
	retryDelay          = 500 * time.Millisecond
	timestamp           = "timestamp" // valid for v0.3
	totalAttempt        = 3
)

type OnChainFeedDatum struct {
	price                  *big.Int
	observationsTimestamp  uint32
	stalenessSeconds       uint32
	deviationPercentagePPM uint32
}

type StreamsLookup struct {
	feedParamKey    string
	feeds           []string
	timeParamKey    string
	time            *big.Int
	extraData       []byte
	upkeepId        *big.Int
	block           uint64
	useFunctions    bool
	onChainFeedData []*OnChainFeedDatum
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
	FeedID                string `json:"feedID"` // feed id in hex
	ValidFromTimestamp    string `json:"validFromTimestamp"`
	ObservationsTimestamp string `json:"observationsTimestamp"`
	FullReport            string `json:"fullReport"` // the actual mercury report of this feed, can be sent to verifier
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
		lggr.Infof("at block %d upkeep %s trying to decodeStreamsLookup performData=%s", block, upkeepId, hexutil.Encode(checkResults[i].PerformData))
		l, err := r.decodeStreamsLookup(res.PerformData)
		if err != nil {
			lggr.Warnf("at block %d upkeep %s decodeStreamsLookup failed: %v", block, upkeepId, err)
			// user contract did not revert with StreamsLookup error
			continue
		}

		if len(l.feeds) == 0 {
			checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonInvalidRevertDataInput)
			lggr.Warnf("at block %s upkeep %s has empty feeds array", block, upkeepId)
			continue
		}
		// mercury permission checking for v0.3 is done by mercury server
		if l.feedParamKey == feedIdHex && l.timeParamKey == blockNumber {
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
				lggr.Warnf("at block %d upkeep %s NOT allowed to query Mercury server", block, upkeepId)
				checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonMercuryAccessNotAllowed)
				continue
			}
		} else if l.feedParamKey != feedIDs || l.timeParamKey != timestamp {
			// if mercury version cannot be determined, set failure reason
			lggr.Warnf("at block %d upkeep %s NOT allowed to query Mercury server", block, upkeepId)
			checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonInvalidRevertDataInput)
			continue
		}

		l.upkeepId = upkeepId.BigInt()
		// the block here is exclusively used to call checkCallback at this block, not to be confused with the block number
		// in the revert for mercury v0.2, which is denoted by time in the struct bc starting from v0.3, only timestamp will be supported
		l.block = uint64(block.Int64())
		lggr.Infof("at block %d upkeep %s decodeStreamsLookup feedKey=%s timeKey=%s feeds=%v time=%s extraData=%s", block, upkeepId, l.feedParamKey, l.timeParamKey, l.feeds, l.time, hexutil.Encode(l.extraData))
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

	values := [][]byte{common.Hex2Bytes("0006a2f7f9b6c10385739c687064aa1e457812927f59446cccddf7740cc025ad00000000000000000000000000000000000000000000000000000000014cb94e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001206962e629c3a0f5b7e3e9294b0c283c9b20f94f1c89c8ba8c1ee4650738f20fb20000000000000000000000000000000000000000000000000000000064e50c980000000000000000000000000000000000000000000000000000025a0864a8c00000000000000000000000000000000000000000000000000000025a063481720000000000000000000000000000000000000000000000000000025a0a94d00f000000000000000000000000000000000000000000000000000000000226181f4733a6d98892d1821771c041d5d69298210fdca9d643ad74477423b6a3045647000000000000000000000000000000000000000000000000000000000226181f0000000000000000000000000000000000000000000000000000000064e50c9700000000000000000000000000000000000000000000000000000000000000027f3056b1b71dd516037afd2e636f8afb39853f5cb3ccaa4b02d6f9a2a64622534e94aa1f794f6a72478deb7e0eb2942864b7fac76d6e120bd809530b1b74a32b00000000000000000000000000000000000000000000000000000000000000027bd3b385c0812dfcad2652d225410a014a0b836cd9635a6e7fb404f65f7a912f0b193db57e5c4f38ce71f29170f7eadfa94d972338858bacd59ab224245206db")}

	for j, v := range values {
		lggr.Infof("upkeep %s doMercuryRequest values[%d]: %s", lookup.upkeepId, j, hexutil.Encode(v))
	}

	var performData []byte
	var needed bool
	if !lookup.useFunctions {
		var state encoding.PipelineExecutionState
		var err error
		var failureReason uint8
		state, retryable, mercuryBytes, err := r.checkCallback(ctx, values, lookup)
		if err != nil {
			lggr.Errorf("at block %d upkeep %s checkCallback err: %v", lookup.block, lookup.upkeepId, err)
			checkResults[i].Retryable = retryable
			checkResults[i].PipelineExecutionState = uint8(state)
			return
		}
		lggr.Infof("checkCallback mercuryBytes=%s", hexutil.Encode(mercuryBytes))

		state, needed, performData, failureReason, _, err = r.packer.UnpackCheckCallbackResult(mercuryBytes)
		if err != nil {
			lggr.Errorf("at block %d upkeep %s UnpackCheckCallbackResult err: %v", lookup.block, lookup.upkeepId, err)
			checkResults[i].PipelineExecutionState = uint8(state)
			return
		}

		if failureReason == uint8(encoding.UpkeepFailureReasonMercuryCallbackReverted) {
			checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonMercuryCallbackReverted)
			lggr.Debugf("at block %d upkeep %s mercury callback reverts", lookup.block, lookup.upkeepId)
			return
		}
	} else {
		var state encoding.PipelineExecutionState
		var err error
		var retryable bool
		state, retryable, needed, performData, _, err = r.checkCallbackWithFunctions(ctx, values, lookup)
		if err != nil {
			lggr.Errorf("at block %d upkeep %s checkCallbackWithFunctions err: %v", lookup.block, lookup.upkeepId, err)
			checkResults[i].Retryable = retryable
			checkResults[i].PipelineExecutionState = uint8(state)
			return
		}
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

	cfg, err := r.registry.GetUpkeepPrivilegeConfig(opts, upkeepId)
	if err != nil {
		return encoding.RpcFlakyFailure, encoding.UpkeepFailureReasonNone, true, false, fmt.Errorf("failed to get upkeep privilege config: %v", err)
	}
	if len(cfg) == 0 {
		r.mercury.allowListCache.Set(upkeepId.String(), false, cache.DefaultExpiration)
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonMercuryAccessNotAllowed, false, false, fmt.Errorf("upkeep privilege config is empty")
	}

	var a UpkeepPrivilegeConfig
	err = json.Unmarshal(cfg, &a)
	if err != nil {
		return encoding.MercuryUnmarshalError, encoding.UpkeepFailureReasonNone, false, false, fmt.Errorf("failed to unmarshal privilege config: %v", err)
	}
	r.mercury.allowListCache.Set(upkeepId.String(), a.MercuryEnabled, cache.DefaultExpiration)
	return encoding.NoPipelineError, encoding.UpkeepFailureReasonNone, false, a.MercuryEnabled, nil
}

// decodeStreamsLookup decodes the revert error StreamsLookup(string feedParamKey, string[] feeds, string feedParamKey, uint256 time, byte[] extraData)
func (r *EvmRegistry) decodeStreamsLookup(data []byte) (*StreamsLookup, error) {
	e := r.mercury.abi.Errors["StreamsLookup"]
	unpack, err := e.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("unpack error: %w", err)
	}
	errorParameters := unpack.([]interface{})

	feeds := *abi.ConvertType(errorParameters[1], new([]string)).(*[]string)
	extraData := *abi.ConvertType(errorParameters[4], new([]byte)).(*[]byte)

	var useFunctions bool = false
	var onChainFeedData []*OnChainFeedDatum
	if bytes.Equal(extraData[:20], []byte("CHECK_WITH_FUNCTIONS")) {
		useFunctions = true
		for i := range feeds {
			start := 20 + 32*i
			price := big.NewInt(0).SetBytes(extraData[start : start+20])
			observationsTimestamp := binary.BigEndian.Uint32(extraData[start+20 : start+24])
			stalenessSeconds := binary.BigEndian.Uint32(extraData[start+24 : start+28])
			deviationPercentagePPM := binary.BigEndian.Uint32(extraData[start+28 : start+32])
			onChainFeedData = append(onChainFeedData, &OnChainFeedDatum{
				price:                  price,
				observationsTimestamp:  observationsTimestamp,
				stalenessSeconds:       stalenessSeconds,
				deviationPercentagePPM: deviationPercentagePPM,
			})
		}
	}

	streamsLookup := StreamsLookup{
		feedParamKey:    *abi.ConvertType(errorParameters[0], new(string)).(*string),
		feeds:           feeds,
		timeParamKey:    *abi.ConvertType(errorParameters[2], new(string)).(*string),
		time:            *abi.ConvertType(errorParameters[3], new(*big.Int)).(**big.Int),
		extraData:       extraData,
		useFunctions:    useFunctions,
		onChainFeedData: onChainFeedData,
	}

	return &streamsLookup, nil
}

func (r *EvmRegistry) checkCallback(ctx context.Context, values [][]byte, lookup *StreamsLookup) (encoding.PipelineExecutionState, bool, hexutil.Bytes, error) {
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

func (r *EvmRegistry) checkCallbackWithFunctions(
	ctx context.Context,
	values [][]byte,
	lookup *StreamsLookup,
) (encoding.PipelineExecutionState, bool, bool, []byte, uint8, error) {

	r.lggr.Infow("LOOKUP_REVERT_DATA",
		"deviation", lookup.onChainFeedData[0].deviationPercentagePPM,
		"timestamp", lookup.onChainFeedData[0].observationsTimestamp,
		"price", lookup.onChainFeedData[0].price.String(),
		"staleness", lookup.onChainFeedData[0].stalenessSeconds,
	)
	// Gateway is still blocked. We simply stub the functionality and return the values seen off-chain.
	// Below this logic is what could be run for gateway, commented out.
	abiEncode, err := utils.ABIEncode(`[{"type":"bytes[]"},{"type":"bytes"}]`, values, lookup.extraData)
	if err != nil {
		return encoding.PackUnpackDecodeFailed, false, false, nil, 0, nil
	}
	return 0, false, true, abiEncode, 0, nil

	/* url := "GATEWAY_URL"
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "request",
		"id":      "1001",
		"params": map[string]interface{}{
			"signature": "",
			"args": map[string]interface{}{
				"0": lookup.block, // provide block number agreed upon in-case of non-deterministic logic
				"1": json.marshal(lookup.OnChainFeedData), // provide on-chain data
				"2": values, // provide off-chain data
			},
			"body": map[string]interface{}{
				"don_id": "functions_experimental",
				"sender": "GATEWAY_DEMO",
				"payload": map[string]string{
					"source":       [ENTER_SCRIPT_HERE],
					"codeLocation": "0",
					"language":     "0",
				},
			},
		},
	} */
}

// doMercuryRequest sends requests to Mercury API to retrieve mercury data.
func (r *EvmRegistry) doMercuryRequest(ctx context.Context, sl *StreamsLookup, lggr logger.Logger) (encoding.PipelineExecutionState, encoding.UpkeepFailureReason, [][]byte, bool, error) {
	var isMercuryV03 bool
	resultLen := len(sl.feeds)
	ch := make(chan MercuryData, resultLen)
	if len(sl.feeds) == 0 {
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonInvalidRevertDataInput, [][]byte{}, false, fmt.Errorf("invalid revert data input: feed param key %s, time param key %s, feeds %s", sl.feedParamKey, sl.timeParamKey, sl.feeds)
	}
	if sl.feedParamKey == feedIdHex && sl.timeParamKey == blockNumber {
		// only mercury v0.2
		for i := range sl.feeds {
			go r.singleFeedRequest(ctx, ch, i, sl, lggr)
		}
	} else if sl.feedParamKey == feedIDs && sl.timeParamKey == timestamp {
		// only mercury v0.3
		resultLen = 1
		isMercuryV03 = true
		ch = make(chan MercuryData, resultLen)
		go r.multiFeedsRequest(ctx, ch, sl, lggr)
	} else {
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonInvalidRevertDataInput, [][]byte{}, false, fmt.Errorf("invalid revert data input: feed param key %s, time param key %s, feeds %s", sl.feedParamKey, sl.timeParamKey, sl.feeds)
	}

	var reqErr error
	results := make([][]byte, len(sl.feeds))
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
	lggr.Debugf("upkeep %s retryable %t reqErr %w", sl.upkeepId.String(), retryable && !allSuccess, reqErr)
	// only retry when not all successful AND none are not retryable
	return state, encoding.UpkeepFailureReasonNone, results, retryable && !allSuccess, reqErr
}

// singleFeedRequest sends a v0.2 Mercury request for a single feed report.
func (r *EvmRegistry) singleFeedRequest(ctx context.Context, ch chan<- MercuryData, index int, sl *StreamsLookup, lggr logger.Logger) {
	q := url.Values{
		sl.feedParamKey: {sl.feeds[index]},
		sl.timeParamKey: {sl.time.String()},
	}
	mercuryURL := r.mercury.cred.URL
	reqUrl := fmt.Sprintf("%s%s%s", mercuryURL, mercuryPathV02, q.Encode())
	lggr.Infow("request URL for upkeep %s feed %s: %s", sl.upkeepId.String(), sl.feeds[index], reqUrl)

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
				lggr.Warnf("at block %s upkeep %s GET request fails for feed %s: %v", sl.time.String(), sl.upkeepId.String(), sl.feeds[index], err1)
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
				lggr.Warnf("at block %s upkeep %s received status code %d for feed %s", sl.time.String(), sl.upkeepId.String(), resp.StatusCode, sl.feeds[index])
				retryable = true
				state = encoding.MercuryFlakyFailure
				return errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
			} else if resp.StatusCode != http.StatusOK {
				retryable = false
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("at block %s upkeep %s received status code %d for feed %s", sl.time.String(), sl.upkeepId.String(), resp.StatusCode, sl.feeds[index])
			}

			var m MercuryV02Response
			err1 = json.Unmarshal(body, &m)
			if err1 != nil {
				lggr.Warnf("at block %s upkeep %s failed to unmarshal body to MercuryV02Response for feed %s: %v", sl.time.String(), sl.upkeepId.String(), sl.feeds[index], err1)
				retryable = false
				state = encoding.MercuryUnmarshalError
				return err1
			}
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				lggr.Warnf("at block %s upkeep %s failed to decode chainlinkBlob %s for feed %s: %v", sl.time.String(), sl.upkeepId.String(), m.ChainlinkBlob, sl.feeds[index], err1)
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
			Error:     fmt.Errorf("failed to request feed for %s: %w", sl.feeds[index], retryErr),
			State:     state,
		}
		ch <- md
	}
}

// multiFeedsRequest sends a Mercury v0.3 request for a multi-feed report
func (r *EvmRegistry) multiFeedsRequest(ctx context.Context, ch chan<- MercuryData, sl *StreamsLookup, lggr logger.Logger) {
	q := url.Values{
		feedIDs:   {strings.Join(sl.feeds, ",")},
		timestamp: {sl.time.String()},
	}

	reqUrl := fmt.Sprintf("%s%s%s", r.mercury.cred.URL, mercuryBatchPathV03, q.Encode())
	lggr.Debugf("request URL for upkeep %s feed %s: %s", sl.upkeepId.String(), strings.Join(sl.feeds, ","), reqUrl)

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
				lggr.Warnf("at block %s upkeep %s GET request fails from mercury v0.3: %v", sl.time.String(), sl.upkeepId.String(), err1)
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
				lggr.Warnf("at block %s upkeep %s received status code %d from mercury v0.3", sl.time.String(), sl.upkeepId.String(), resp.StatusCode)
				retryable = true
				state = encoding.MercuryFlakyFailure
				return errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
			} else if resp.StatusCode != http.StatusOK {
				retryable = false
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("at block %s upkeep %s received status code %d from mercury v0.3", sl.time.String(), sl.upkeepId.String(), resp.StatusCode)
			}

			var response MercuryV03Response
			err1 = json.Unmarshal(body, &response)
			if err1 != nil {
				lggr.Warnf("at block %s upkeep %s failed to unmarshal body to MercuryV03Response from mercury v0.3: %v", sl.time.String(), sl.upkeepId.String(), err1)
				retryable = false
				state = encoding.MercuryUnmarshalError
				return err1
			}
			if len(response.Reports) != len(sl.feeds) {
				// this should never happen. if this upkeep does not have permission for any feeds it requests, or if certain feeds are
				// missing in mercury server, the mercury server v0.3 should respond with 400s, rather than returning partial results
				retryable = false
				state = encoding.InvalidMercuryResponse
				return fmt.Errorf("at block %s upkeep %s requested %d feeds but received %d reports from mercury v0.3", sl.time.String(), sl.upkeepId.String(), len(sl.feeds), len(response.Reports))
			}
			var reportBytes [][]byte
			var b []byte
			for _, rsp := range response.Reports {
				b, err1 = hexutil.Decode(rsp.FullReport)
				if err1 != nil {
					lggr.Warnf("upkeep %s block %s failed to decode fullReport %s from mercury v0.3: %v", sl.upkeepId.String(), sl.time.String(), rsp.FullReport, err1)
					retryable = false
					state = encoding.InvalidMercuryResponse
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
