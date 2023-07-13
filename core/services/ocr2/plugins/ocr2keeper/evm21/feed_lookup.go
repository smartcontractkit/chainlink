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

type MercuryBytes struct {
	Index     int
	Error     error
	Retryable bool
	Bytes     []byte
}

// AdminOffchainConfig represents the administrative offchain config for each upkeep. It can be set by s_upkeepManager
// role on the registry. Upkeeps allowed to use Mercury server will have this set to true.
type AdminOffchainConfig struct {
	MercuryEnabled bool `json:"mercuryEnabled"`
}

// feedLookup looks through check upkeep results looking for any that need off chain lookup
func (r *EvmRegistry) feedLookup(ctx context.Context, upkeepResults []EVMAutomationUpkeepResult21) ([]EVMAutomationUpkeepResult21, error) {
	lookups := map[int]*FeedLookup{}
	for i := range upkeepResults {
		if upkeepResults[i].FailureReason != UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED {
			continue
		}

		block := upkeepResults[i].Block
		upkeepId := upkeepResults[i].ID
		opts, err := r.buildCallOpts(ctx, big.NewInt(int64(block)))
		if err != nil {
			r.lggr.Errorf("[FeedLookup] upkeep %s block %d buildCallOpts: %v", upkeepId, block, err)
			return nil, err
		}

		allowed, err := r.allowedToUseMercury(opts, upkeepId)
		if err != nil {
			r.lggr.Errorf("[FeedLookup] upkeep %s block %d failed to time mercury allow list: %v", upkeepId, block, err)
			continue
		}

		if !allowed {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_MERCURY_ACCESS_NOT_ALLOWED
			r.lggr.Errorf("[FeedLookup] upkeep %s block %d NOT allowed to time Mercury server", upkeepId, block)
			continue
		}

		r.lggr.Infof("[FeedLookup] upkeep %s block %d decodeFeedLookup performData=%s", upkeepId, block, hexutil.Encode(upkeepResults[i].PerformData))
		lookup, err := r.decodeFeedLookup(upkeepResults[i].PerformData)
		if err != nil {
			r.lggr.Errorf("[FeedLookup] upkeep %s block %d decodeFeedLookup: %v", upkeepId, block, err)
			continue
		}
		lookup.upkeepId = upkeepId
		// the block here is exclusively used to call checkCallback at this block, not to be confused with the block number
		// in the revert for mercury v0.2, which is denoted by time in the struct bc starting from v0.3, only timestamp will be supported
		lookup.block = uint64(block)
		r.lggr.Infof("[FeedLookup] upkeep %s block %d decodeFeedLookup feedKey=%s timeKey=%s feeds=%v time=%s extraData=%s", upkeepId, block, lookup.feedParamKey, lookup.timeParamKey, lookup.feeds, lookup.time, hexutil.Encode(lookup.extraData))
		lookups[i] = lookup
	}

	var wg sync.WaitGroup
	for i, lookup := range lookups {
		wg.Add(1)
		go r.doLookup(ctx, &wg, lookup, i, upkeepResults)
	}
	wg.Wait()

	// don't surface error to plugin bc FeedLookup process should be self-contained.
	return upkeepResults, nil
}

func (r *EvmRegistry) doLookup(ctx context.Context, wg *sync.WaitGroup, lookup *FeedLookup, i int, upkeepResults []EVMAutomationUpkeepResult21) {
	defer wg.Done()

	values, retryable, err := r.doMercuryRequest(ctx, lookup)
	if err != nil {
		r.lggr.Errorf("[FeedLookup] upkeep %s retryable %v doMercuryRequest: %v", lookup.upkeepId, retryable, err)
		upkeepResults[i].Retryable = retryable
		return
	}
	for j, v := range values {
		r.lggr.Infof("[FeedLookup] checkCallback values[%d]=%s", j, hexutil.Encode(v))
	}

	mercuryBytes, err := r.checkCallback(ctx, values, lookup)
	if err != nil {
		r.lggr.Errorf("[FeedLookup] upkeep %s block %d checkCallback err: %v", lookup.upkeepId, lookup.block, err)
		return
	}
	r.lggr.Infof("[FeedLookup] checkCallback mercuryBytes=%s", hexutil.Encode(mercuryBytes))

	needed, performData, failureReason, _, err := r.packer.UnpackCheckCallbackResult(mercuryBytes)
	if err != nil {
		r.lggr.Errorf("[FeedLookup] upkeep %s block %d UnpackCheckCallbackResult err: %v", lookup.upkeepId, lookup.block, err)
		return
	}

	if int(failureReason) == UPKEEP_FAILURE_REASON_MERCURY_CALLBACK_REVERTED {
		upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_MERCURY_CALLBACK_REVERTED
		r.lggr.Debugf("[FeedLookup] upkeep %s block %d mercury callback reverts", lookup.upkeepId, lookup.block)
		return
	}

	if !needed {
		upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_UPKEEP_NOT_NEEDED
		r.lggr.Debugf("[FeedLookup] upkeep %s block %d callback reports upkeep not needed", lookup.upkeepId, lookup.block)
		return
	}

	upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_NONE
	upkeepResults[i].Eligible = true
	upkeepResults[i].PerformData = performData
	r.lggr.Infof("[FeedLookup] upkeep %s block %d successful with perform data: %s", lookup.upkeepId, lookup.block, hexutil.Encode(performData))
}

// allowedToUseMercury retrieves upkeep's administrative offchain config and decode a mercuryEnabled bool to indicate if
// this upkeep is allowed to use Mercury service.
func (r *EvmRegistry) allowedToUseMercury(opts *bind.CallOpts, upkeepId *big.Int) (bool, error) {
	allowed, ok := r.mercury.allowListCache.Get(upkeepId.String())
	if ok {
		return allowed.(bool), nil
	}

	cfg, err := r.registry.GetUpkeepPrivilegeConfig(opts, upkeepId)
	if err != nil {
		return false, fmt.Errorf("failed to get upkeep privilege config for upkeep ID %s: %v", upkeepId, err)
	}

	var a AdminOffchainConfig
	err = json.Unmarshal(cfg, &a)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal privilege config for upkeep ID %s: %v", upkeepId, err)
	}
	r.mercury.allowListCache.Set(upkeepId.String(), a.MercuryEnabled, cache.DefaultExpiration)
	return a.MercuryEnabled, nil
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

func (r *EvmRegistry) checkCallback(ctx context.Context, values [][]byte, lookup *FeedLookup) (hexutil.Bytes, error) {
	payload, err := r.abi.Pack("checkCallback", lookup.upkeepId, values, lookup.extraData)
	if err != nil {
		return nil, err
	}

	var b hexutil.Bytes
	args := map[string]interface{}{
		"to":   r.addr.Hex(),
		"data": hexutil.Bytes(payload),
	}

	// call checkCallback function at the block which OCR3 has agreed upon
	err = r.client.CallContext(ctx, &b, "eth_call", args, hexutil.EncodeUint64(lookup.block))
	if err != nil {
		return nil, err
	}
	return b, nil
}

// doMercuryRequest sends requests to Mercury API to retrieve ChainlinkBlob.
func (r *EvmRegistry) doMercuryRequest(ctx context.Context, ml *FeedLookup) ([][]byte, bool, error) {
	// TODO (AUTO-3253): if no feed labels are provided in v0.3, request for all feeds
	resultLen := len(ml.feeds)
	ch := make(chan MercuryBytes, resultLen)
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
			ch = make(chan MercuryBytes, resultLen)
			go r.multiFeedsRequest(ctx, ch, ml)
		}
	} else {
		return nil, false, fmt.Errorf("invalid label combination: feed param key %s and time param key %s", ml.feedParamKey, ml.timeParamKey)
	}

	var reqErr error
	results := make([][]byte, len(ml.feeds))
	retryable := true
	allSuccess := true
	for i := 0; i < len(results); i++ {
		m := <-ch
		if m.Error != nil {
			reqErr = errors.Join(reqErr, m.Error)
			retryable = retryable && m.Retryable
			allSuccess = false
		}
		results[m.Index] = m.Bytes
	}
	r.lggr.Debugf("FeedLookup upkeep %s retryable %s reqErr %w", ml.upkeepId.String(), retryable && !allSuccess, reqErr)
	// only retry when not all successful AND none are not retryable
	return results, retryable && !allSuccess, reqErr
}

// singleFeedRequest sends a Mercury request for a single feed report.
func (r *EvmRegistry) singleFeedRequest(ctx context.Context, ch chan<- MercuryBytes, index int, ml *FeedLookup, mv MercuryVersion) {
	q := url.Values{
		ml.feedParamKey: {ml.feeds[index]},
		ml.timeParamKey: {ml.time.String()},
		UserId:          {ml.upkeepId.String()},
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
		ch <- MercuryBytes{Index: index, Error: err}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := r.generateHMAC(http.MethodGet, path+q.Encode(), []byte{}, r.mercury.cred.Username, r.mercury.cred.Password, ts)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.mercury.cred.Username)
	req.Header.Set("X-Authorization-Timestamp", strconv.FormatInt(ts, 10))
	req.Header.Set("X-Authorization-Signature-SHA256", signature)

	retryable := false
	retryErr := retry.Do(
		func() error {
			retryable = false
			resp, err1 := r.hc.Do(req)
			if err1 != nil {
				r.lggr.Errorf("FeedLookup upkeep %s block %s GET request fails for feed %s: %v", ml.upkeepId.String(), ml.time.String(), ml.feeds[index], err1)
				return err1
			}
			defer resp.Body.Close()
			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				return err1
			}

			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
				r.lggr.Errorf("FeedLookup upkeep %s block %s received status code %d for feed %s", ml.upkeepId.String(), ml.time.String(), resp.StatusCode, ml.feeds[index])
				retryable = true
				return errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
			} else if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("FeedLookup upkeep %s block %s received status code %d for feed %s", ml.upkeepId.String(), ml.time.String(), resp.StatusCode, ml.feeds[index])
			}

			var m MercuryResponse
			err1 = json.Unmarshal(body, &m)
			if err1 != nil {
				r.lggr.Errorf("FeedLookup upkeep %s block %s failed to unmarshal body to MercuryResponse for feed %s: %v", ml.upkeepId.String(), ml.time.String(), ml.feeds[index], err1)
				return err1
			}
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				r.lggr.Errorf("FeedLookup upkeep %s block %s failed to decode chainlinkBlob %s for feed %s: %v", ml.upkeepId.String(), ml.time.String(), m.ChainlinkBlob, ml.feeds[index], err1)
				return err1
			}
			ch <- MercuryBytes{Index: index, Bytes: blobBytes}
			return nil
		},
		// only retry when the error is 404 Not Found or 500 Internal Server Error
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound) || err.Error() == fmt.Sprintf("%d", http.StatusInternalServerError)
		}),
		retry.Context(ctx),
		retry.Delay(RetryDelay),
		retry.Attempts(TotalAttempt))

	// if all retries fail, return the error and ask the caller to handle cool down and heavyweight retry
	if retryErr != nil {
		mb := MercuryBytes{
			Index:     index,
			Retryable: retryable,
			Error:     retryErr,
		}
		ch <- mb
	}
}

// multiFeedsRequest sends a Mercury request for a multi-feed report
func (r *EvmRegistry) multiFeedsRequest(ctx context.Context, ch chan<- MercuryBytes, ml *FeedLookup) {
	q := url.Values{
		FeedId:    {strings.Join(ml.feeds, ",")},
		Timestamp: {ml.time.String()},
		UserId:    {ml.upkeepId.String()},
	}

	reqUrl := fmt.Sprintf("%s%s%s", r.mercury.cred.URL, MercuryBatchPathV3, q.Encode())
	r.lggr.Debugf("FeedLookup request URL: %s", reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		ch <- MercuryBytes{Index: 0, Error: err}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := r.generateHMAC(http.MethodGet, MercuryBatchPathV3+q.Encode(), []byte{}, r.mercury.cred.Username, r.mercury.cred.Password, ts)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.mercury.cred.Username)
	req.Header.Set("X-Authorization-Timestamp", strconv.FormatInt(ts, 10))
	req.Header.Set("X-Authorization-Signature-SHA256", signature)

	retryable := false
	retryErr := retry.Do(
		func() error {
			retryable = false
			resp, err1 := r.hc.Do(req)
			if err1 != nil {
				r.lggr.Errorf("FeedLookup upkeep %s block %s GET request fails for multi feed: %v", ml.upkeepId.String(), ml.time.String(), err1)
				return err1
			}
			defer resp.Body.Close()
			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				return err1
			}

			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
				r.lggr.Errorf("FeedLookup upkeep %s block %s received status code %d for multi feed", ml.upkeepId.String(), ml.time.String(), resp.StatusCode)
				retryable = true
				return errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
			} else if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("FeedLookup upkeep %s block %s received status code %d for multi feed", ml.upkeepId.String(), ml.time.String(), resp.StatusCode)
			}

			var m MercuryResponse
			err1 = json.Unmarshal(body, &m)
			if err1 != nil {
				r.lggr.Errorf("FeedLookup upkeep %s block %s failed to unmarshal body to MercuryResponse for multi feed: %v", ml.upkeepId.String(), ml.time.String(), err1)
				return err1
			}
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				r.lggr.Errorf("FeedLookup upkeep %s block %s failed to decode chainlinkBlob %s for multi feed: %v", ml.upkeepId.String(), ml.time.String(), m.ChainlinkBlob, err1)
				return err1
			}
			ch <- MercuryBytes{
				Index: 0,
				Bytes: blobBytes,
			}
			return nil
		},
		// only retry when the error is 404 Not Found or 500 Internal Server Error
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound) || err.Error() == fmt.Sprintf("%d", http.StatusInternalServerError)
		}),
		retry.Context(ctx),
		retry.Delay(RetryDelay),
		retry.Attempts(TotalAttempt))

	// if all retries fail, return the error and ask the caller to handle cool down and heavyweight retry
	if retryErr != nil {
		mb := MercuryBytes{
			Index:     0,
			Retryable: retryable,
			Error:     retryErr,
		}
		ch <- mb
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
