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
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
)

type MercuryLookup struct {
	feedLabel  string
	feeds      []string
	queryLabel string
	query      *big.Int
	extraData  []byte
}

type MercuryResponse struct {
	ChainlinkBlob string `json:"chainlinkBlob"`
}

type MercuryBytes struct {
	Index int
	Error error
	Bytes []byte
}

const (
	RetryDelay   = 750 * time.Millisecond
	TotalAttempt = 3
)

// mercuryLookup looks through check upkeep results looking for any that need off chain lookup
func (r *EvmRegistry) mercuryLookup(ctx context.Context, upkeepResults []types.UpkeepResult) ([]types.UpkeepResult, error) {
	// return error only if there are errors which stops the process
	// don't surface Mercury API errors to plugin bc MercuryLookup process should be self-contained
	// TODO (AUTO-2862): parallelize the mercury lookup work for all upkeeps
	for i := range upkeepResults {
		block, upkeepId, err := blockAndIdFromKey(upkeepResults[i].Key)
		if err != nil {
			r.lggr.Error("MercuryLookup error getting block and upkeep id:", err)
			return nil, err
		}

		// if its another reason continue/skip
		r.lggr.Debugf("MercuryLookup upkeep ID %s block %s has failure reason %d status %d", upkeepId.String(), block.String(), upkeepResults[i].FailureReason, upkeepResults[i].State)
		if upkeepResults[i].FailureReason != UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED {
			r.lggr.Debugf("MercuryLookup %s failure reason is NOT UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED. Won't do mercury lookup", upkeepId.String())
			continue
		}

		// checking if this upkeep is in cooldown from api errors
		_, onIce := r.mercury.cooldownCache.Get(upkeepId.String())
		if onIce {
			r.lggr.Debugf("MercuryLookup upkeep %s block %s skipped bc of cool down", upkeepId.String(), block.String())
			continue
		}

		// if it doesn't decode to the mercury custom error continue/skip
		mercuryLookup, err := r.decodeMercuryLookup(upkeepResults[i].PerformData)
		if err != nil {
			r.lggr.Errorf("MercuryLookup upkeep %s block %s decodeMercuryLookup: %v", upkeepId.String(), block.String(), err)
			continue
		}
		r.lggr.Debugf("MercuryLookup upkeep %s block %s decodeMercuryLookup success", upkeepId.String(), block.String())

		opts, err := r.buildCallOpts(ctx, block)
		if err != nil {
			r.lggr.Errorf("MercuryLookup upkeep %s block %s buildCallOpts: %v", upkeepId.String(), block.String(), err)
			return nil, err
		}
		// need upkeep info for offchainConfig and to hit callback
		upkeepInfo, err := r.getUpkeepInfo(upkeepId, opts)
		if err != nil {
			r.lggr.Errorf("MercuryLookup upkeep %s block %s GetUpkeep: %v", upkeepId.String(), block.String(), err)
			return nil, err
		}

		// do the mercury lookup request
		values, err := r.doMercuryRequest(ctx, mercuryLookup, upkeepId)
		if err != nil {
			r.lggr.Errorf("MercuryLookup upkeep %s block %s doMercuryRequest: %v", upkeepId.String(), block.String(), err)
			continue
		}
		r.lggr.Debugf("MercuryLookup upkeep %s block %s doMercuryRequest success", upkeepId.String(), block.String())

		needed, performData, err := r.mercuryLookupCallback(ctx, mercuryLookup, values, upkeepInfo, opts)
		if err != nil {
			r.lggr.Errorf("MercuryLookup upkeep %s block %s mercuryLookupCallback err: %v", upkeepId.String(), block.String(), err)
			continue
		}
		r.lggr.Debugf("MercuryLookup upkeep %s block %s mercuryLookupCallback success", upkeepId.String(), block.String())
		if !needed {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_UPKEEP_NOT_NEEDED
			r.lggr.Debugf("MercuryLookup upkeep %s block %s callback reports upkeep not needed", upkeepId.String(), block.String())
			continue
		}

		upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_NONE
		upkeepResults[i].State = types.Eligible
		upkeepResults[i].PerformData = performData
		r.lggr.Debugf("MercuryLookup upkeep %s block %s successful with perform data: %+v", upkeepId.String(), block.String(), performData)
	}
	// don't surface error to plugin bc MercuryLookup process should be self-contained.
	return upkeepResults, nil
}

func (r *EvmRegistry) getUpkeepInfo(upkeepId *big.Int, opts *bind.CallOpts) (keeper_registry_wrapper2_0.UpkeepInfo, error) {
	zero := common.Address{}
	var err error
	var upkeepInfo keeper_registry_wrapper2_0.UpkeepInfo
	u, found := r.mercury.upkeepCache.Get(upkeepId.String())
	if found {
		upkeepInfo = u.(keeper_registry_wrapper2_0.UpkeepInfo)
		r.lggr.Debugf("MercuryLookup upkeep %s block %s cache hit UpkeepInfo: %+v", upkeepId.String(), opts.BlockNumber.String(), upkeepInfo)
	} else {
		upkeepInfo, err = r.registry.GetUpkeep(opts, upkeepId)
		if err != nil {
			return upkeepInfo, err
		}
		if upkeepInfo.Target == zero {
			return upkeepInfo, errors.New("upkeepInfo should not be nil")
		}
		r.lggr.Debugf("MercuryLookup upkeep %s block %s cache miss UpkeepInfo: %+v", upkeepId.String(), opts.BlockNumber.String(), upkeepInfo)
		r.mercury.upkeepCache.Set(upkeepId.String(), upkeepInfo, cache.DefaultExpiration)
	}
	return upkeepInfo, nil
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

// mercuryLookupCallback calls the callback(string[] memory chainlinkBlobHex, bytes memory extraData) specified by the
// 4-byte selector from the revert. the return will match check telling us if the upkeep is needed and what the perform data is
func (r *EvmRegistry) mercuryLookupCallback(ctx context.Context, mercuryLookup *MercuryLookup, values [][]byte, upkeepInfo keeper_registry_wrapper2_0.UpkeepInfo, opts *bind.CallOpts) (bool, []byte, error) {
	payload, err := r.mercury.abi.Pack("mercuryCallback", values, mercuryLookup.extraData)
	if err != nil {
		return false, nil, fmt.Errorf("callback args pack error: %w", err)
	}

	// use gas 0 to provide infinite gas and from empty address
	callbackMsg := ethereum.CallMsg{
		To:   &upkeepInfo.Target,
		Data: payload,
	}

	callbackResp, err := r.client.CallContract(ctx, callbackMsg, opts.BlockNumber)
	if err != nil {
		return false, nil, fmt.Errorf("call contract callback error: %w", err)
	}

	return r.packer.UnpackMercuryLookupResult(callbackResp)
}

func (r *EvmRegistry) doMercuryRequest(ctx context.Context, ml *MercuryLookup, upkeepId *big.Int) ([][]byte, error) {
	ch := make(chan MercuryBytes, len(ml.feeds))
	for i := range ml.feeds {
		go r.singleFeedRequest(ctx, ch, upkeepId, i, ml)
	}
	var reqErr error
	results := make([][]byte, len(ml.feeds))
	for i := 0; i < len(results); i++ {
		m := <-ch
		if m.Error != nil {
			reqErr = errors.Join(reqErr, fmt.Errorf("feed[%s]: %w", ml.feeds[i], m.Error))
		}
		results[m.Index] = m.Bytes
	}
	r.lggr.Debugf("MercuryLookup upkeep %s reqErr %w", upkeepId.String(), reqErr)
	return results, reqErr
}

func (r *EvmRegistry) singleFeedRequest(ctx context.Context, ch chan<- MercuryBytes, upkeepId *big.Int, index int, ml *MercuryLookup) {
	q := url.Values{
		ml.feedLabel:  {ml.feeds[index]},
		ml.queryLabel: {ml.query.String()},
	}
	reqUrl := fmt.Sprintf("%s/client?%s", r.mercury.cred.URL, q.Encode())
	r.lggr.Debugf("MercuryLookup request URL: %s", reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		ch <- MercuryBytes{Index: index, Error: err}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := r.generateHMAC(http.MethodGet, "/client?"+q.Encode(), []byte{}, r.mercury.cred.Username, r.mercury.cred.Password, ts)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.mercury.cred.Username)
	req.Header.Set("X-Authorization-Timestamp", strconv.FormatInt(ts, 10))
	req.Header.Set("X-Authorization-Signature-SHA256", signature)

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
			empty := MercuryResponse{}
			if m == empty || m.ChainlinkBlob == "" {
				r.lggr.Errorf("MercuryLookup upkeep %s block %s response is empty", upkeepId.String(), ml.query.String())
			}
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				r.lggr.Errorf("MercuryLookup upkeep %s block %s failed to decode chainlinkBlob %s for feed %s: %v", upkeepId.String(), ml.query.String(), m.ChainlinkBlob, ml.feeds[index], err1)
				return err1
			}
			ch <- MercuryBytes{Index: index, Bytes: blobBytes}
			return nil
		},
		// only retry when the error is 404 Not Found
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound)
		}),
		retry.Context(ctx),
		retry.Delay(RetryDelay),
		retry.Attempts(TotalAttempt))

	// if all retries fail, it's very likely the feed IDs are incorrect or block number is too old or too new, put into cooldownf
	if retryErr != nil {
		ch <- MercuryBytes{Index: index, Error: retryErr}
		r.setCachesOnAPIErr(upkeepId)
	}
}

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

// setCachesOnAPIErr when an off chain look up request fails or gets a 4xx/5xx response code we increment error count and put the upkeep in cooldown state
func (r *EvmRegistry) setCachesOnAPIErr(upkeepId *big.Int) {
	r.lggr.Debugf("MercuryLookup adding %s to API error cache", upkeepId.String())
	errCount := 1
	cacheKey := upkeepId.String()
	e, ok := r.mercury.apiErrCache.Get(cacheKey)
	if ok {
		errCount = e.(int) + 1
	}

	// With a 10m Error Cache Window every error sets the error count and resets the TTL to 10m
	// On every error that hits during this rolling 10m the error count is increased and the cooldown period by associate is increased
	// This means the user will suffer a max cooldown of 17m on the 10th error at which point the error cache will have expired since its window is 10m
	// After that the user will reset to 0 and start over after a combined total of 34m in cooldown state.

	// increment error count and reset expiration to shift window with last seen error
	r.mercury.apiErrCache.Set(cacheKey, errCount, cache.DefaultExpiration)
	// put upkeep in cooldown state for 2^errors seconds.
	r.mercury.cooldownCache.Set(cacheKey, nil, time.Second*time.Duration(2^errCount))
}
