package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/dontpanicdao/caigo/types"
	"github.com/google/go-querystring/query"
)

type StarkResp struct {
	Result []string `json:"result"`
}

type StateUpdate struct {
	BlockHash string `json:"block_hash"`
	NewRoot   string `json:"new_root"`
	OldRoot   string `json:"old_root"`
	StateDiff struct {
		StorageDiffs      map[string]interface{} `json:"storage_diffs"`
		DeployedContracts []struct {
			Address   string `json:"address"`
			ClassHash string `json:"class_hash"`
		} `json:"deployed_contracts"`
	} `json:"state_diff"`
}

func (sg *Gateway) ChainID(context.Context) (string, error) {
	return sg.ChainId, nil
}

type GatewayFunctionCall struct {
	FunctionCall
	Signature []string `json:"signature"`
}

type FunctionCall types.FunctionCall

func (f FunctionCall) MarshalJSON() ([]byte, error) {
	output := map[string]interface{}{}
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
	return json.Marshal(output)
}

/*
'call_contract' wrapper and can accept a blockId in the hash or height format
*/
func (sg *Gateway) Call(ctx context.Context, call types.FunctionCall, blockHashOrTag string) ([]string, error) {
	gc := GatewayFunctionCall{
		FunctionCall: FunctionCall(call),
	}
	gc.EntryPointSelector = types.BigToHex(types.GetSelectorFromName(gc.EntryPointSelector))
	if len(gc.Calldata) == 0 {
		gc.Calldata = []string{}
	}

	if len(gc.Signature) == 0 {
		gc.Signature = []string{"0", "0"} // allows rpc and http clients to implement(has to be a better way)
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/call_contract", gc)
	if err != nil {
		return nil, err
	}

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

	var resp StarkResp
	return resp.Result, sg.do(req, &resp)
}

/*
'add_transaction' wrapper for invokation requests
*/
func (sg *Gateway) Invoke(ctx context.Context, invoke types.FunctionInvoke) (*types.AddInvokeTransactionOutput, error) {
	tx := Transaction{
		Type:            INVOKE,
		ContractAddress: invoke.ContractAddress.Hex(),
		Version:         fmt.Sprintf("0x%d", invoke.Version),
		MaxFee:          fmt.Sprintf("0x%s", invoke.MaxFee.Text(16)),
	}
	if invoke.EntryPointSelector != "" {
		tx.EntryPointSelector = types.BigToHex(types.GetSelectorFromName(invoke.EntryPointSelector))
	}
	if invoke.Nonce != nil {
		tx.Nonce = fmt.Sprintf("0x%s", invoke.Nonce.Text(16))
	}

	calldata := []string{}
	for _, v := range invoke.Calldata {
		bv, _ := big.NewInt(0).SetString(v, 0)
		calldata = append(calldata, bv.Text(10))
	}
	tx.Calldata = calldata

	if len(invoke.Signature) == 0 {
		tx.Signature = []string{}
	} else {
		// stop-gap before full types.Felt cutover
		tx.Signature = []string{invoke.Signature[0].String(), invoke.Signature[1].String()}
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/add_transaction", tx)
	if err != nil {
		return nil, err
	}
	var resp types.AddInvokeTransactionOutput
	return &resp, sg.do(req, &resp)
}

/*
'add_transaction' wrapper for compressing and deploying a compiled StarkNet contract
*/
func (sg *Gateway) Deploy(ctx context.Context, contract types.ContractClass, deployRequest types.DeployRequest) (resp types.AddDeployResponse, err error) {
	d := DeployRequest(deployRequest)
	d.Type = DEPLOY
	if len(d.ConstructorCalldata) == 0 {
		d.ConstructorCalldata = []string{}
	}
	if d.ContractAddressSalt == "" {
		d.ContractAddressSalt = "0x0"
	}

	d.ContractDefinition = contract
	if err != nil {
		return resp, err
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/add_transaction", d)
	if err != nil {
		return resp, err
	}

	return resp, sg.do(req, &resp)
}

type DeployAccountRequest types.DeployAccountRequest

func (d DeployAccountRequest) MarshalJSON() ([]byte, error) {
	if d.Type != "DEPLOY_ACCOUNT" {
		return nil, errors.New("wrong type")
	}
	output := map[string]interface{}{}
	constructorCalldata := []string{}
	for _, value := range d.ConstructorCalldata {
		constructorCalldata = append(constructorCalldata, types.SNValToBN(value).Text(10))
	}
	output["constructor_calldata"] = constructorCalldata
	output["max_fee"] = fmt.Sprintf("0x%s", d.MaxFee.Text(16))
	output["version"] = fmt.Sprintf("0x%s", big.NewInt(int64(d.Version)).Text(16))
	signature := []string{}
	for _, value := range d.Signature {
		signature = append(signature, value.Text(10))
	}
	output["signature"] = signature
	nonce := "0x0"
	if d.Nonce != nil {
		output["version"] = fmt.Sprintf("0x%s", d.Nonce.Text(16))
	}
	output["nonce"] = nonce
	output["type"] = "DEPLOY_ACCOUNT"
	if d.ContractAddressSalt == "" {
		d.ContractAddressSalt = "0x0"
	}
	contractAddressSalt := fmt.Sprintf("0x%s", types.SNValToBN(d.ContractAddressSalt).Text(16))
	output["contract_address_salt"] = contractAddressSalt
	classHash := fmt.Sprintf("0x%s", types.SNValToBN(d.ClassHash).Text(16))
	output["class_hash"] = classHash
	return json.Marshal(output)
}

/*
'add_transaction' wrapper for deploying a compiled StarkNet account
*/
func (sg *Gateway) DeployAccount(ctx context.Context, deployAccountRequest types.DeployAccountRequest) (resp types.AddDeployResponse, err error) {
	d := DeployAccountRequest(deployAccountRequest)
	d.Type = DEPLOY_ACCOUNT

	req, err := sg.newRequest(ctx, http.MethodPost, "/add_transaction", d)
	if err != nil {
		return resp, err
	}

	return resp, sg.do(req, &resp)
}

/*
'add_transaction' wrapper for compressing and declaring a contract class
*/
func (sg *Gateway) Declare(ctx context.Context, contract types.ContractClass, declareRequest DeclareRequest) (resp types.AddDeclareResponse, err error) {
	declareRequest.Type = DECLARE
	declareRequest.SenderAddress = "0x1"
	declareRequest.MaxFee = "0x0"
	declareRequest.Nonce = "0x0"
	declareRequest.Signature = []string{}
	declareRequest.ContractClass = contract
	if err != nil {
		return resp, err
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/add_transaction", declareRequest)
	if err != nil {
		return resp, err
	}

	return resp, sg.do(req, &resp)
}

type DeployRequest types.DeployRequest

func (d DeployRequest) MarshalJSON() ([]byte, error) {
	calldata := []string{}
	for _, value := range d.ConstructorCalldata {
		calldata = append(calldata, types.SNValToBN(value).Text(10))
	}
	d.ConstructorCalldata = calldata
	return json.Marshal(types.DeployRequest(d))
}

type DeclareRequest struct {
	Type          string              `json:"type"`
	SenderAddress string              `json:"sender_address"`
	MaxFee        string              `json:"max_fee"`
	Nonce         string              `json:"nonce"`
	Signature     []string            `json:"signature"`
	ContractClass types.ContractClass `json:"contract_class"`
}

func (sg *Gateway) StateUpdate(ctx context.Context, opts *BlockOptions) (*StateUpdate, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_state_update", nil)
	if err != nil {
		return nil, err
	}

	if opts != nil {
		vs, err := query.Values(opts)
		if err != nil {
			return nil, err
		}
		appendQueryValues(req, vs)
	}

	var resp StateUpdate
	return &resp, sg.do(req, &resp)
}

func (sg *Gateway) ContractAddresses(ctx context.Context) (*ContractAddresses, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_contract_addresses", nil)
	if err != nil {
		return nil, err
	}

	var resp ContractAddresses
	return &resp, sg.do(req, &resp)
}

type ContractAddresses struct {
	Starknet             string `json:"Starknet"`
	GpsStatementVerifier string `json:"GpsStatementVerifier"`
}
