package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"
)

type OffchainLookup struct {
	url              string
	extraData        []byte
	fields           []string
	callbackFunction [4]byte
}

type OffchainLookupBody struct {
	sender string
	data   string
}

// offchainLookup looks through check upkeep results looking for any that need off chain lookup
func (r *EvmRegistry) offchainLookup(ctx context.Context, upkeepResults []types.UpkeepResult) error {
	for i, _ := range upkeepResults {
		// if its another reason just keep going
		if upkeepResults[i].FailureReason != UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED {
			continue
		}

		var offchainLookup OffchainLookup
		offchainLookup, err := decodeOffchainLookup(upkeepResults[i].PerformData)
		if err != nil {
			r.lggr.Error("[OffchainLookup] decodeOffchainLookup=", err)
			continue
		}
		r.lggr.Infof("[OffchainLookup]= %+v\n", offchainLookup)

		// 	do the http request
		body, statusCode, err := offchainLookup.doRequest()
		if err != nil {
			r.lggr.Error("[OffchainLookup] doRequest=", err)
			continue
		}
		r.lggr.Infof("[OffchainLookup] StatusCode: %d\n", statusCode)
		r.lggr.Infof("[OffchainLookup] Body: %s\n", string(body))

		values, err := offchainLookup.parseJson(body)
		if err != nil {
			r.lggr.Error("[OffchainLookup] parseJson=", err)
			continue
		}
		r.lggr.Infof("[OffchainLookup] Parsed values= %v\n", values)

		// get block for call opts
		block, upkeepId, err := blockAndIdFromKey(upkeepResults[i].Key)
		if err != nil {
			r.lggr.Error("[OffchainLookup] BlockAndIdFromKey=", err)
			continue
		}
		opts, err := r.buildCallOpts(ctx, block)
		if err != nil {
			r.lggr.Error("[OffchainLookup] buildCallOpts=", err)
			continue
		}
		// need upkeep info to hit callback
		upkeepInfo, err := r.registry.GetUpkeep(opts, upkeepId)
		if err != nil {
			r.lggr.Error("[OffchainLookup] GetUpkeep=", err)
			continue
		}
		r.lggr.Infof("[OffchainLookup] GetUpkeep: %v\n", upkeepInfo)

		needed, performData, err := r.offchainLookupCallback(ctx, offchainLookup, values, statusCode, upkeepInfo, opts)
		if err != nil || !needed {
			r.lggr.Info("[OffchainLookup] offchainLookupCallback=", needed, performData, err)
			continue
		}

		// success! set failure and state to eligible and off chain perform data
		upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_NONE
		upkeepResults[i].State = types.Eligible
		upkeepResults[i].PerformData = performData
	}
	return nil
}

// decodeOffchainLookup decodes the revert error ChainlinkAPIFetch(string url, bytes extraData, string[] jsonFields, bytes4 callbackSelector)
func decodeOffchainLookup(data []byte) (OffchainLookup, error) {
	abiUpkeepAPIFetch, err := abi.JSON(strings.NewReader("[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"jsonFields\",\"type\":\"string[]\"},{\"internalType\":\"bytes4\",\"name\":\"callbackSelector\",\"type\":\"bytes4\"}],\"name\":\"ChainlinkAPIFetch\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"fact\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"values\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"statusCode\",\"type\":\"uint256\"}],\"name\":\"callback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"fields\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"input\",\"type\":\"string\"}],\"name\":\"setURLs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"s\",\"type\":\"string\"}],\"name\":\"stringToUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"url\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"))
	if err != nil {
		return OffchainLookup{}, err
	}
	offchainLookup := OffchainLookup{}
	e := abiUpkeepAPIFetch.Errors["ChainlinkAPIFetch"]
	unpack, err := e.Unpack(data)
	if err != nil {
		return OffchainLookup{}, errors.Wrapf(err, "unpack error:")
	}
	errorParameters := unpack.([]interface{})

	offchainLookup.url = *abi.ConvertType(errorParameters[0], new(string)).(*string)
	offchainLookup.extraData = *abi.ConvertType(errorParameters[1], new([]byte)).(*[]byte)
	offchainLookup.fields = *abi.ConvertType(errorParameters[2], new([]string)).(*[]string)
	offchainLookup.callbackFunction = *abi.ConvertType(errorParameters[3], new([4]byte)).(*[4]byte)
	return offchainLookup, nil
}

// offchainLookupCallback calls the callback(bytes calldata extraData, string[] calldata values, uint256 statusCode)
// the return will match check telling us if the upkeep is needed and what the perform data is
func (r *EvmRegistry) offchainLookupCallback(ctx context.Context, offchainLookup OffchainLookup, values []string, statusCode int, upkeepInfo keeper_registry_wrapper2_0.UpkeepInfo, opts *bind.CallOpts) (bool, []byte, error) {
	// call to the contract function specified by the 4-byte selector callbackFunction
	typBytes, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return false, nil, errors.Wrapf(err, "abi new bytes type error:")
	}
	typStrings, err := abi.NewType("string[]", "", nil)
	if err != nil {
		return false, nil, errors.Wrapf(err, "abi new strings type error:")
	}
	typUint, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return false, nil, errors.Wrapf(err, "abi new uint256 type error:")
	}
	callbackArgs := abi.Arguments{
		{Name: "extraData", Type: typBytes},
		{Name: "values", Type: typStrings},
		{Name: "statusCode", Type: typUint},
	}
	pack, err := callbackArgs.Pack(offchainLookup.extraData, values, big.NewInt(int64(statusCode)))
	if err != nil {
		return false, nil, errors.Wrapf(err, "callback args pack error:")
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
		return false, nil, errors.Wrapf(err, "call contract callback error:")
	}

	boolTyp, err := abi.NewType("bool", "", nil)
	callbackOutput := abi.Arguments{
		{Name: "upkeepNeeded", Type: boolTyp},
		{Name: "performData", Type: typBytes},
	}
	unpack, err := callbackOutput.Unpack(callbackResp)
	if err != nil {
		return false, nil, errors.Wrapf(err, "callback ouput unpack error:")
	}

	upkeepNeeded := *abi.ConvertType(unpack[0], new(bool)).(*bool)
	if !upkeepNeeded {
		return false, nil, nil
	}
	performData := *abi.ConvertType(unpack[1], new([]byte)).(*[]byte)
	return true, performData, nil
}

func (o *OffchainLookup) doRequest() ([]byte, int, error) {
	client := http.Client{}
	var req *http.Request
	var err error

	req, err = http.NewRequest("GET", o.url, nil)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "get request error:")
	}

	// Make an HTTP GET request to the request URL.
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "do request error:")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "read body error:")
	}

	return body, resp.StatusCode, nil
}

func (o *OffchainLookup) parseJson(body []byte) ([]string, error) {
	var m map[string]interface{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(o.fields), len(o.fields))
	for i, field := range o.fields {
		for key, val := range m {
			if key == field {
				result[i] = fmt.Sprint(val)
			}
		}
	}
	return result, nil
}
