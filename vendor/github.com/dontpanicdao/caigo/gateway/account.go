package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dontpanicdao/caigo/types"
)

func (sg *Gateway) AccountNonce(ctx context.Context, address types.Hash) (*big.Int, error) {
	resp, err := sg.Call(ctx, types.FunctionCall{
		ContractAddress:    address,
		EntryPointSelector: "get_nonce",
	}, "")
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("no resp in contract call 'get_nonce' %v", address)
	}

	return types.HexToBN(resp[0]), nil
}

func (sg *Gateway) Nonce(ctx context.Context, contractAddress, blockHashOrTag string) (*big.Int, error) {

	req, err := sg.newRequest(ctx, http.MethodGet, "/get_nonce", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"contractAddress": []string{contractAddress},
	})
	switch {
	case strings.HasPrefix(blockHashOrTag, "0x"):
		appendQueryValues(req, url.Values{
			"blockHash": []string{blockHashOrTag},
		})
	case blockHashOrTag == "":
		appendQueryValues(req, url.Values{
			"blockNumber": []string{"pending"},
		})
	default:
		appendQueryValues(req, url.Values{
			"blockNumber": []string{blockHashOrTag},
		})
	}

	var resp string
	err = sg.do(req, &resp)
	if err != nil {
		return nil, err
	}
	nonce, ok := big.NewInt(0).SetString(resp, 0)
	if !ok {
		return nil, errors.New("nonce not found")
	}
	return nonce, nil
}

type functionInvoke types.FunctionInvoke

func (f functionInvoke) MarshalJSON() ([]byte, error) {
	output := map[string]interface{}{}
	sigs := []string{}
	for _, sig := range f.Signature {
		sigs = append(sigs, sig.Text(10))
	}
	output["signature"] = sigs
	output["contract_address"] = f.ContractAddress.Hex()
	if f.EntryPointSelector != "" {
		output["entry_point_selector"] = f.EntryPointSelector
	}
	calldata := []string{}
	for _, v := range f.Calldata {
		data, _ := big.NewInt(0).SetString(v, 0)
		calldata = append(calldata, data.Text(10))
	}
	output["calldata"] = calldata
	if f.Nonce != nil {
		output["nonce"] = json.RawMessage(
			strconv.Quote(fmt.Sprintf("0x%s", f.Nonce.Text(16))),
		)
	}
	if f.MaxFee != nil {
		output["max_fee"] = json.RawMessage(
			strconv.Quote(fmt.Sprintf("0x%s", f.MaxFee.Text(16))),
		)
	}
	output["version"] = json.RawMessage(strconv.Quote(fmt.Sprintf("0x%s", f.Version.Text(16))))
	output["type"] = "INVOKE_FUNCTION"
	return json.Marshal(output)
}

func (sg *Gateway) EstimateFee(ctx context.Context, call types.FunctionInvoke, hash string) (*types.FeeEstimate, error) {
	if call.EntryPointSelector != "" {
		call.EntryPointSelector = types.BigToHex(types.GetSelectorFromName(call.EntryPointSelector))
	}
	c := functionInvoke(call)
	req, err := sg.newRequest(ctx, http.MethodPost, "/estimate_fee", c)
	if err != nil {
		return nil, err
	}

	if hash != "" {
		appendQueryValues(req, url.Values{
			"blockNumber": []string{hash},
		})
	}
	output := map[string]json.RawMessage{}
	err = sg.do(req, &output)
	if err != nil {
		return nil, err
	}
	gasPrice := output["gas_price"]
	gasConsumed := output["gas_usage"]
	overallFee := output["overall_fee"]
	overallFeeInt, _ := big.NewInt(0).SetString(string(overallFee), 0)
	gasPriceInt, _ := big.NewInt(0).SetString(string(gasPrice), 0)
	gasConsumedInt, _ := big.NewInt(0).SetString(string(gasConsumed), 0)
	resp := types.FeeEstimate{
		GasConsumed: types.NumAsHex("0x" + gasConsumedInt.Text(16)),
		GasPrice:    types.NumAsHex("0x" + gasPriceInt.Text(16)),
		OverallFee:  types.NumAsHex("0x" + overallFeeInt.Text(16)),
	}
	return &resp, nil
}
