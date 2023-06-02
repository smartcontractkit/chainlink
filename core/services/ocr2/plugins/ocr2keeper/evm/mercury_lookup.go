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
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

const (
	BlockNumber        = "blockNumber" // valid for v0.2
	FeedID             = "feedID"      // valid for v0.3
	FeedIDHex          = "feedIDHex"   // valid for v0.2
	MercuryHostV2      = ""
	MercuryHostV3      = ""
	MercuryPathV2      = "/client?"
	MercuryPathV3      = "/v1/reports?"
	MercuryBatchPathV3 = "/v1/reports/bulk?"
	RetryDelay         = 600 * time.Millisecond
	Timestamp          = "timestamp" // valid for v0.3
	TotalAttempt       = 4
	UserId             = "userId"
)

type MercuryLookup struct {
	feedLabel  string
	feeds      []string
	queryLabel string
	query      *big.Int
	extraData  []byte
}

// MercuryResponse is used in both single feed endpoint and bulk endpoint because bulk endpoint will return ONE
// chainlinkBlob which contains multiple reports instead of multiple blobs.
type MercuryResponse struct {
	ChainlinkBlob string `json:"chainlinkBlob"`
}

/**
 *           retryable   error    bytes
 *   200        N          N        Y
 *   404        Y          Y        N
 *   other      N          Y        N
 */

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

// mercuryLookup looks through check upkeep results looking for any that need off chain lookup
func (r *EvmRegistry) mercuryLookup(ctx context.Context, upkeepResults []EVMAutomationUpkeepResult20) ([]EVMAutomationUpkeepResult20, error) {
	// return error only if there are errors which stops the process
	// don't surface Mercury API errors to plugin bc MercuryLookup process should be self-contained
	// TODO (AUTO-2862): parallelize the mercury lookup work for all upkeeps
	for i := range upkeepResults {
		// if its another reason continue/skip
		if upkeepResults[i].FailureReason != UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED {
			continue
		}

		block := upkeepResults[i].Block
		upkeepId := upkeepResults[i].ID

		opts, err := r.buildCallOpts(ctx, big.NewInt(int64(block)))
		if err != nil {
			r.lggr.Errorf("[MercuryLookup] upkeep %s block %d buildCallOpts: %v", upkeepId, block, err)
			return nil, err
		}

		allowed, err := r.allowedToUseMercury(opts, upkeepId)
		r.lggr.Info(allowed)
		if err != nil {
			r.lggr.Errorf("[MercuryLookup] upkeep %s block %d failed to query mercury allow list: %v", upkeepId, block, err)
			continue
		}

		// do anything? failure reason not allowed?
		if !allowed {
			r.lggr.Errorf("[MercuryLookup] upkeep %s block %d NOT allowed to query Mercury server", upkeepId, block)
			continue
		}

		r.lggr.Debugf("[MercuryLookup] upkeep %s block %d perform data: %v", upkeepId, block, upkeepResults[i].PerformData)
		// if it doesn't decode to the mercury custom error continue/skip
		mercuryLookup, err := r.decodeMercuryLookup(upkeepResults[i].PerformData)
		if err != nil {
			r.lggr.Errorf("[MercuryLookup] upkeep %s block %d decodeMercuryLookup: %v", upkeepId, block, err)
			continue
		}

		// do the mercury lookup request
		values, retryable, err := r.doMercuryRequest(ctx, mercuryLookup, upkeepId)
		if err != nil {
			r.lggr.Errorf("[MercuryLookup] upkeep %s block %d doMercuryRequest: %v", upkeepId, block, err)
			r.lggr.Infof("[MercuryLookup] upkeep %s block %d retryable: %s", upkeepId, block, retryable)
			upkeepResults[i].Retryable = retryable
			continue
		}

		r.lggr.Debugf("[MercuryLookup] upkeep %s block %d values: %v", values)
		r.lggr.Debugf("[MercuryLookup] upkeep %s block %d extraData: %v", mercuryLookup.extraData)
		needed, performData, failureReason, _, err := r.mercuryCallback21(ctx, upkeepId, values, mercuryLookup.extraData, block)
		if err != nil {
			r.lggr.Errorf("[MercuryLookup] upkeep %s block %d mercuryLookupCallback err: %v", upkeepId, block, err)
			continue
		}
		r.lggr.Debugf("[MercuryLookup] upkeep %s block %d performData: %v", performData)
		r.lggr.Debugf("[MercuryLookup] upkeep %s block %d failureReason: %v", failureReason)

		if int(failureReason) == UPKEEP_FAILURE_REASON_MERCURY_CALLBACK_REVERTED {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_MERCURY_CALLBACK_REVERTED
			r.lggr.Debugf("[MercuryLookup] upkeep %s block %d mercury callback reverts", upkeepId, block)
			continue
		}

		if !needed {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_UPKEEP_NOT_NEEDED
			r.lggr.Debugf("[MercuryLookup] upkeep %s block %d callback reports upkeep not needed", upkeepId, block)
			continue
		}

		upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_NONE
		upkeepResults[i].Eligible = true
		upkeepResults[i].PerformData = performData
		r.lggr.Infof("[MercuryLookup] upkeep %s block %d successful with perform data: %+v", upkeepId, block, performData)
	}
	// don't surface error to plugin bc MercuryLookup process should be self-contained.
	return upkeepResults, nil
}

// allowedToUseMercury retrieves upkeep's administrative offchain config and decode a mercuryEnabled bool to indicate if
// this upkeep is allowed to use Mercury service.
func (r *EvmRegistry) allowedToUseMercury(opts *bind.CallOpts, upkeepId *big.Int) (bool, error) {
	allowed, ok := r.mercury.mercuryAllowListCache.Get(upkeepId.String())
	r.lggr.Info(allowed)
	if ok {
		return allowed.(bool), nil
	}

	cfg, err := r.registry21.GetUpkeepAdminOffchainConfig(opts, upkeepId)
	if err != nil {
		return false, fmt.Errorf("failed to get upkeep admin offchain config for upkeep ID %s: %v", upkeepId, err)
	}

	var a AdminOffchainConfig
	err = json.Unmarshal(cfg, &a)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal admin offchain config for upkeep ID %s: %v", upkeepId, err)
	}
	r.mercury.mercuryAllowListCache.Set(upkeepId.String(), a.MercuryEnabled, cache.DefaultExpiration)
	return a.MercuryEnabled, nil
}

// decodeMercuryLookup decodes the revert error MercuryLookup(string feedLabel, string[] feeds, string feedLabel, uint256 query, byte[] extraData)
func (r *EvmRegistry) decodeMercuryLookup(data []byte) (*MercuryLookup, error) {
	e := r.mercury.abi.Errors["MercuryLookup"]
	unpack, err := e.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("unpack error: %w", err)
	}
	errorParameters := unpack.([]interface{})

	return &MercuryLookup{
		feedLabel:  *abi.ConvertType(errorParameters[0], new(string)).(*string),
		feeds:      *abi.ConvertType(errorParameters[1], new([]string)).(*[]string),
		queryLabel: *abi.ConvertType(errorParameters[2], new(string)).(*string),
		query:      *abi.ConvertType(errorParameters[3], new(*big.Int)).(**big.Int),
		extraData:  *abi.ConvertType(errorParameters[4], new([]byte)).(*[]byte),
	}, nil
}

// mercuryCallback21 calls mercuryCallback function on registry 2.1. The registry will forward the call to the mercuryCallback function
// on user's contract with proper check gas limit.
func (r *EvmRegistry) mercuryCallback21(ctx context.Context, upkeepID *big.Int, values [][]byte, ed []byte, block uint32) (bool, []byte, uint8, *big.Int, error) {
	payload, err := r.abi21.Pack("mercuryCallback", upkeepID, values, ed)
	if err != nil {
		return false, nil, 0, nil, err
	}

	var b hexutil.Bytes
	args := map[string]interface{}{
		"to":   r.addr.Hex(),
		"data": hexutil.Bytes(payload),
	}

	err = r.client.CallContext(ctx, &b, "eth_call", args, hexutil.EncodeUint64(uint64(block)))
	if err != nil {
		return false, nil, 0, nil, err
	}

	return r.packer.UnpackMercuryLookupResult(b)
}

// doMercuryRequest
func (r *EvmRegistry) doMercuryRequest(ctx context.Context, ml *MercuryLookup, upkeepId *big.Int) ([][]byte, bool, error) {
	// TODO (AUTO-3253): if no feed labels are provided, request for all feeds
	resultLen := len(ml.feeds)
	ch := make(chan MercuryBytes, resultLen)
	if ml.feedLabel == FeedIDHex && ml.queryLabel == BlockNumber {
		// only mercury v0.2
		for i := range ml.feeds {
			go r.singleFeedRequest(ctx, ch, upkeepId, i, ml, job.MercuryV02)
		}
	} else if ml.feedLabel == FeedID && ml.queryLabel == Timestamp {
		// only mercury v0.3
		if resultLen == 1 {
			go r.singleFeedRequest(ctx, ch, upkeepId, 0, ml, job.MercuryV03)
		} else {
			// create a new channel with buffer size 1 since the batch endpoint will only return 1 blob
			resultLen = 1
			ch = make(chan MercuryBytes, resultLen)
			go r.multiFeedsRequest(ctx, ch, upkeepId, ml)
		}
	} else {
		return nil, false, fmt.Errorf("invalid label combination: feed label %s and query label %s", ml.feedLabel, ml.queryLabel)
	}

	var reqErr error
	retryable := true
	results := make([][]byte, resultLen)
	for i := 0; i < len(results); i++ {
		m := <-ch
		if m.Error != nil {
			retryable = false
			reqErr = errors.Join(reqErr, m.Error)
		}
		results[m.Index] = m.Bytes
	}
	r.lggr.Debugf("MercuryLookup upkeep %s retryable %s reqErr %w", upkeepId.String(), retryable, reqErr)

	// set retryable properly
	return results, retryable, reqErr
}

// singleFeedRequest sends a Mercury request for a single feed report.
func (r *EvmRegistry) singleFeedRequest(ctx context.Context, ch chan<- MercuryBytes, upkeepId *big.Int, index int, ml *MercuryLookup, mv job.MercuryVersion) {
	q := url.Values{
		ml.feedLabel:  {ml.feeds[index]},
		ml.queryLabel: {ml.query.String()},
		UserId:        {upkeepId.String()},
	}
	mercuryURL := MercuryHostV2
	path := MercuryPathV2
	if mv == job.MercuryV03 {
		mercuryURL = MercuryHostV3
		path = MercuryPathV3
	}
	reqUrl := fmt.Sprintf("%s%s%s", mercuryURL, path, q.Encode())
	r.lggr.Debugf("MercuryLookup request URL: %s", reqUrl)

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
			resp, err1 := r.hc.Do(req)
			if err1 != nil {
				r.lggr.Errorf("MercuryLookup upkeep %s block %s GET request fails for feed %s: %v", upkeepId.String(), ml.query.String(), ml.feeds[index], err1)
				return err1
			}
			defer resp.Body.Close()
			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				r.lggr.Errorf("MercuryLookup upkeep %s block %s fails to read response body for feed %s: %v", upkeepId.String(), ml.query.String(), ml.feeds[index], err1)
				return err1
			}

			if resp.StatusCode == http.StatusNotFound {
				// there are 2 possible causes for 404: incorrect URL and querying a block where report has not been generated
				r.lggr.Errorf("MercuryLookup upkeep %s block %s received status code %d for feed %s", upkeepId.String(), ml.query.String(), resp.StatusCode, ml.feeds[index])
				// return 404 for retry
				retryable = true
				return fmt.Errorf("%d", http.StatusNotFound)
			} else if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("MercuryLookup upkeep %s block %s received status code %d for feed %s", upkeepId.String(), ml.query.String(), resp.StatusCode, ml.feeds[index])
			}

			var m MercuryResponse
			err1 = json.Unmarshal(body, &m)
			if err1 != nil {
				r.lggr.Errorf("MercuryLookup upkeep %s block %s failed to unmarshal body to MercuryResponse for feed %s: %v", upkeepId.String(), ml.query.String(), ml.feeds[index], err1)
				return err1
			}
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				r.lggr.Errorf("MercuryLookup upkeep %s block %s failed to decode chainlinkBlob %s for feed %s: %v", upkeepId.String(), ml.query.String(), m.ChainlinkBlob, ml.feeds[index], err1)
				return err1
			}
			ch <- MercuryBytes{
				Index: index,
				Bytes: blobBytes,
			}
			return nil
		},
		// only retry when the error is 404 Not Found
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound)
		}),
		retry.Context(ctx),
		retry.Delay(RetryDelay),
		retry.Attempts(TotalAttempt))

	// if all retries fail, return the error and ask the caller to handle cool down and heavyweight retry
	if retryErr != nil || retryable {
		mb := MercuryBytes{
			Index:     index,
			Retryable: retryable,
			Error:     retryErr,
		}
		ch <- mb
	}
}

// multiFeedsRequest sends a Mercury request for a multi-feed report
func (r *EvmRegistry) multiFeedsRequest(ctx context.Context, ch chan<- MercuryBytes, upkeepId *big.Int, ml *MercuryLookup) {
	q := url.Values{
		ml.queryLabel: {ml.query.String()},
		UserId:        {upkeepId.String()},
	}
	// verify array params
	for _, f := range ml.feeds {
		q.Add(ml.feedLabel, f)
	}

	reqUrl := fmt.Sprintf("%s%s%s", MercuryHostV3, MercuryBatchPathV3, q.Encode())
	r.lggr.Debugf("MercuryLookup request URL: %s", reqUrl)

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
			resp, err1 := r.hc.Do(req)
			if err1 != nil {
				r.lggr.Errorf("MercuryLookup upkeep %s block %s GET request fails for multi feed: %v", upkeepId.String(), ml.query.String(), err1)
				return err1
			}
			defer resp.Body.Close()
			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				r.lggr.Errorf("MercuryLookup upkeep %s block %s fails to read response body for multi feed: %v", upkeepId.String(), ml.query.String(), err1)
				return err1
			}

			if resp.StatusCode == http.StatusNotFound {
				// there are 2 possible causes for 404: incorrect URL and querying a block where report has not been generated
				r.lggr.Errorf("MercuryLookup upkeep %s block %s received status code %d for multi feed", upkeepId.String(), ml.query.String(), resp.StatusCode)
				// return 404 for retry
				retryable = true
				return fmt.Errorf("%d", http.StatusNotFound)
			} else if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("MercuryLookup upkeep %s block %s received status code %d for multi feed", upkeepId.String(), ml.query.String(), resp.StatusCode)
			}

			var m MercuryResponse
			err1 = json.Unmarshal(body, &m)
			if err1 != nil {
				r.lggr.Errorf("MercuryLookup upkeep %s block %s failed to unmarshal body to MercuryResponse for multi feed: %v", upkeepId.String(), ml.query.String(), err1)
				return err1
			}
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				r.lggr.Errorf("MercuryLookup upkeep %s block %s failed to decode chainlinkBlob %s for multi feed: %v", upkeepId.String(), ml.query.String(), m.ChainlinkBlob, err1)
				return err1
			}
			ch <- MercuryBytes{
				Index: 0,
				Bytes: blobBytes,
			}
			return nil
		},
		// only retry when the error is 404 Not Found
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound)
		}),
		retry.Context(ctx),
		retry.Delay(RetryDelay),
		retry.Attempts(TotalAttempt))

	// if all retries fail, return the error and ask the caller to handle cool down and heavyweight retry
	if retryErr != nil || retryable {
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
