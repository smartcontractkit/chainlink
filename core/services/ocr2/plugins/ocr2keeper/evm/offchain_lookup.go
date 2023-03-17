package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itchyny/gojq"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"
)

const ParseFieldError = "field error"

type OffchainLookup struct {
	url              string
	extraData        []byte
	fields           []string
	callbackFunction [4]byte
}

// offchainLookup looks through check upkeep results looking for any that need off chain lookup
func (r *EvmRegistry) offchainLookup(ctx context.Context, upkeepResults []types.UpkeepResult) ([]types.UpkeepResult, error) {
	for i := range upkeepResults {
		// if its another reason continue/skip
		if upkeepResults[i].FailureReason != UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED {
			continue
		}
		block, upkeepId, err := blockAndIdFromKey(upkeepResults[i].Key)
		if err != nil {
			r.lggr.Error("[OffchainLookup] error getting block and upkeep id:", err)
			continue
		}

		// checking if this upkeep is in cooldown from api errors
		_, onIce := r.cooldownCache.Get(upkeepId.String())
		if onIce {
			r.lggr.Infof("[OffchainLookup] cooldown Skipping UpkeepId: %s\n", upkeepId)
			continue
		}

		// if it doesn't decode to the offchain custom error continue/skip
		offchainLookup, err := r.decodeOffchainLookup(upkeepResults[i].PerformData)
		if err != nil {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_OFFCHAIN_LOOKUP_ERROR
			r.lggr.Debug("[OffchainLookup] not an offchain revert decodeOffchainLookup:", err)
			continue
		}
		r.lggr.Debugf("[OffchainLookup]: %+v\n", offchainLookup)

		opts, err := r.buildCallOpts(ctx, block)
		if err != nil {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_OFFCHAIN_LOOKUP_ERROR
			r.lggr.Error("[OffchainLookup] buildCallOpts:", err)
			continue
		}
		// need upkeep info for offchainConfig and to hit callback
		upkeepInfo, err := r.getUpkeepInfo(upkeepId, opts)
		if err != nil {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_OFFCHAIN_LOOKUP_ERROR
			r.lggr.Error("[OffchainLookup] GetUpkeep:", err)
			continue
		}

		// 	do the http request
		body, statusCode, err := r.doRequest(offchainLookup, upkeepId, upkeepInfo.OffchainConfig)
		if err != nil {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_OFFCHAIN_LOOKUP_ERROR
			r.lggr.Error("[OffchainLookup] doRequest:", err)
			continue
		}
		r.lggr.Debugf("[OffchainLookup] StatusCode: %d\n", statusCode)
		//r.lggr.Debugf("[OffchainLookup] Body: %s\n", string(body))

		values, err := offchainLookup.parseJson(body)
		if err != nil {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_OFFCHAIN_LOOKUP_ERROR
			r.lggr.Error("[OffchainLookup] parseJson:", err)
			continue
		}
		r.lggr.Debugf("[OffchainLookup] Parsed values: %+v\n", values)

		needed, performData, err := r.offchainLookupCallback(ctx, offchainLookup, values, statusCode, upkeepInfo, opts)
		if err != nil {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_OFFCHAIN_LOOKUP_ERROR
			r.lggr.Error("[OffchainLookup] offchainLookupCallback=", err)
			continue
		}
		if !needed {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_UPKEEP_NOT_NEEDED
			r.lggr.Debug("[OffchainLookup] callback reports upkeep not needed")
			continue
		}

		// success!
		upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_NONE
		upkeepResults[i].State = types.Eligible
		upkeepResults[i].PerformData = performData
		r.lggr.Debugf("[OffchainLookup] Success: %+v\n", upkeepResults[i])
	}
	return upkeepResults, nil
}

func (r *EvmRegistry) getUpkeepInfo(upkeepId *big.Int, opts *bind.CallOpts) (keeper_registry_wrapper2_0.UpkeepInfo, error) {
	zero := common.Address{}
	var err error
	var upkeepInfo keeper_registry_wrapper2_0.UpkeepInfo
	u, found := r.upkeepCache.Get(upkeepId.String())
	if found {
		upkeepInfo = u.(keeper_registry_wrapper2_0.UpkeepInfo)
		r.lggr.Debugf("[OffchainLookup] cache hit UpkeepInfo: %+v\n", upkeepInfo)
	} else {
		upkeepInfo, err = r.registry.GetUpkeep(opts, upkeepId)
		if err != nil {
			return upkeepInfo, err
		}
		if upkeepInfo.Target == zero {
			return upkeepInfo, errors.New("upkeepInfo should not be nil")
		}
		r.lggr.Debugf("[OffchainLookup] cache miss UpkeepInfo: %+v\n", upkeepInfo)
		r.upkeepCache.Set(upkeepId.String(), upkeepInfo, cache.DefaultExpiration)
	}
	return upkeepInfo, nil
}

// decodeOffchainLookup decodes the revert error ChainlinkAPIFetch(string url, bytes extraData, string[] jsonFields, bytes4 callbackSelector)
func (r *EvmRegistry) decodeOffchainLookup(data []byte) (OffchainLookup, error) {
	e := r.apiFetchABI.Errors["ChainlinkAPIFetch"]
	unpack, err := e.Unpack(data)
	if err != nil {
		return OffchainLookup{}, errors.Wrapf(err, "unpack error")
	}
	errorParameters := unpack.([]interface{})

	return OffchainLookup{
		url:              *abi.ConvertType(errorParameters[0], new(string)).(*string),
		extraData:        *abi.ConvertType(errorParameters[1], new([]byte)).(*[]byte),
		fields:           *abi.ConvertType(errorParameters[2], new([]string)).(*[]string),
		callbackFunction: *abi.ConvertType(errorParameters[3], new([4]byte)).(*[4]byte),
	}, nil
}

// offchainLookupCallback calls the callback(bytes calldata extraData, string[] calldata values, uint256 statusCode) specified by the
// 4-byte selector from the revert. the return will match check telling us if the upkeep is needed and what the perform data is
func (r *EvmRegistry) offchainLookupCallback(ctx context.Context, offchainLookup OffchainLookup, values []string, statusCode int, upkeepInfo keeper_registry_wrapper2_0.UpkeepInfo, opts *bind.CallOpts) (bool, []byte, error) {
	// call to the contract function specified by the 4-byte selector callbackFunction
	typBytes, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return false, nil, errors.Wrapf(err, "abi new bytes type error")
	}
	typStrings, err := abi.NewType("string[]", "", nil)
	if err != nil {
		return false, nil, errors.Wrapf(err, "abi new strings type error")
	}
	typUint, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return false, nil, errors.Wrapf(err, "abi new uint256 type error")
	}
	callbackArgs := abi.Arguments{
		{Name: "extraData", Type: typBytes},
		{Name: "values", Type: typStrings},
		{Name: "statusCode", Type: typUint},
	}
	pack, err := callbackArgs.Pack(offchainLookup.extraData, values, big.NewInt(int64(statusCode)))
	if err != nil {
		return false, nil, errors.Wrapf(err, "callback args pack error")
	}

	var callbackPayload []byte
	callbackPayload = append(callbackPayload, offchainLookup.callbackFunction[:]...)
	callbackPayload = append(callbackPayload, pack...)

	checkUpkeepGasLimit := uint32(200000) + uint32(6500000) + uint32(300000) + upkeepInfo.ExecuteGas
	callbackMsg := ethereum.CallMsg{
		From: r.addr,             // registry addr
		To:   &upkeepInfo.Target, // upkeep addr
		Gas:  uint64(checkUpkeepGasLimit),
		Data: hexutil.Bytes(callbackPayload), // function callback(bytes calldata extraData, string[] calldata values, uint256 statusCode)
	}

	callbackResp, err := r.client.CallContract(ctx, callbackMsg, opts.BlockNumber)
	if err != nil {
		return false, nil, errors.Wrapf(err, "call contract callback error")
	}

	boolTyp, err := abi.NewType("bool", "", nil)
	if err != nil {
		return false, nil, errors.Wrapf(err, "abi new bool type error")
	}
	callbackOutput := abi.Arguments{
		{Name: "upkeepNeeded", Type: boolTyp},
		{Name: "performData", Type: typBytes},
	}
	unpack, err := callbackOutput.Unpack(callbackResp)
	if err != nil {
		return false, nil, errors.Wrapf(err, "callback output unpack error")
	}

	upkeepNeeded := *abi.ConvertType(unpack[0], new(bool)).(*bool)
	if !upkeepNeeded {
		return false, nil, nil
	}
	performData := *abi.ConvertType(unpack[1], new([]byte)).(*[]byte)
	return true, performData, nil
}

func (r *EvmRegistry) doRequest(o OffchainLookup, upkeepId *big.Int, offchainConfig []byte) ([]byte, int, error) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	var req *http.Request
	var err error

	apiKeys, err := getAPIKeys(upkeepId, offchainConfig)
	if err != nil {
		r.lggr.Debug("[OffchainLookup] offchain api keys error", err)
	}

	req, err = http.NewRequest("GET", o.url, nil)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "get request error")
	}

	for _, key := range apiKeys.Keys {
		value := key.DecryptVal
		switch strings.ToLower(key.Type) {
		case "header":
			req.Header.Set(key.Name, value)
		case "param":
			newUrlString := strings.ReplaceAll(o.url, fmt.Sprintf("{%s}", key.Name), url.PathEscape(value))
			u, urlErr := url.Parse(newUrlString)
			if urlErr != nil {
				continue
			}
			req.URL = u
		case "hmac":
			// TODO
			continue
		default:
			// not supported
			continue
		}
	}

	// Make an HTTP GET request to the request URL.
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		r.setCachesOnAPIErr(upkeepId)
		return nil, 0, errors.Wrapf(err, "do request error")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.setCachesOnAPIErr(upkeepId)
		return nil, 0, errors.Wrapf(err, "read body error")
	}

	// if http response code is 4xx/5xx then put in cool down
	if resp.StatusCode >= 400 {
		r.setCachesOnAPIErr(upkeepId)
	}

	return body, resp.StatusCode, nil
}

// setCachesOnAPIErr when an off chain look up request fails or gets a 4xx/5xx response code we increment error count and put the upkeep in cooldown state
func (r *EvmRegistry) setCachesOnAPIErr(upkeepId *big.Int) {
	errCount := 1
	cacheKey := upkeepId.String()
	e, ok := r.apiErrCache.Get(cacheKey)
	if ok {
		errCount = e.(int) + 1
	}

	// With a 10m Error Cache Window every error sets the error count and resets the TTL to 10m
	// On every error that hits during this rolling 10m the error count is increased and the cooldown period by associate is increased
	// This means the user will suffer a max cooldown of 17m on the 10th error at which point the error cache will have expired since its window is 10m
	// After that the user will reset to 0 and start over after a combined total of 34m in cooldown state.

	// increment error count and reset expiration to shift window with last seen error
	r.apiErrCache.Set(cacheKey, errCount, cache.DefaultExpiration)
	// put upkeep in cooldown state for 2^errors seconds.
	r.cooldownCache.Set(cacheKey, nil, time.Second*time.Duration(2^errCount))
}

// parseJson expects json bytes as input and offchainLookup.fields to be 'jq' style query strings for extracting data from json
func (o *OffchainLookup) parseJson(body []byte) ([]string, error) {
	var m map[string]any
	err := json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}
	results := make([]string, len(o.fields))
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	for i, field := range o.fields {
		query, parseErr := gojq.Parse(field)
		if parseErr != nil {
			// if we can't parse the jq query use the response we are sending to the user to signal to them there was an error
			fmt.Println(parseErr)
			results[i] = ParseFieldError
			continue
		}

		// always run with context
		iter := query.RunWithContext(ctx, m)
		for {
			fieldValue, ok := iter.Next()
			if !ok {
				break
			}
			if fieldValueError, isError := fieldValue.(error); isError {
				// if we have an issue with getting the value use the result response to signal that to the user
				fmt.Println(fieldValueError)
				results[i] = ParseFieldError
				continue
			}
			if fieldValue == nil {
				results[i] = ""
				continue
			}
			if _, isMap := fieldValue.(map[string]any); isMap {
				// if the return is a map then format it as json
				// so the json subset can be set to the user where they can further parse if they desire
				marshal, marshalErr := gojq.Marshal(fieldValue)
				if marshalErr != nil {
					// if there is an issue marshalling default to setting the field value
					results[i] = fmt.Sprint(fieldValue)
					continue
				}
				results[i] = string(marshal)
				continue
			}
			results[i] = fmt.Sprint(fieldValue)
		}
	}
	return results, nil
}
