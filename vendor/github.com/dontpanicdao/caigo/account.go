package caigo

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/types"
)

var (
	ErrUnsupportedAccount = errors.New("unsupported account implementation")
	MAX_FEE, _            = big.NewInt(0).SetString("0x20000000000", 0)
)

const (
	TRANSACTION_PREFIX      = "invoke"
	EXECUTE_SELECTOR        = "__execute__"
	CONTRACT_ADDRESS_PREFIX = "STARKNET_CONTRACT_ADDRESS"
)

type account interface {
	Sign(msgHash *big.Int) (*big.Int, *big.Int, error)
	TransactionHash(calls []types.FunctionCall, details types.ExecuteDetails) (*big.Int, error)
	Call(ctx context.Context, call types.FunctionCall) ([]string, error)
	Nonce(ctx context.Context) (*big.Int, error)
	EstimateFee(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FeeEstimate, error)
	Execute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.AddInvokeTransactionOutput, error)
}

var _ account = &Account{}

type AccountPlugin interface {
	PluginCall(calls []types.FunctionCall) (types.FunctionCall, error)
}

type ProviderType string

const (
	ProviderRPCv01  ProviderType = "rpcv01"
	ProviderGateway ProviderType = "gateway"
)

type Account struct {
	rpcv01         *rpcv01.Provider
	sequencer      *gateway.GatewayProvider
	provider       ProviderType
	chainId        string
	AccountAddress string
	private        *big.Int
	version        uint64
	plugin         AccountPlugin
}

type AccountOption struct {
	AccountPlugin AccountPlugin
	version       uint64
}

type AccountOptionFunc func(string, string) (AccountOption, error)

func AccountVersion0(string, string) (AccountOption, error) {
	return AccountOption{
		version: uint64(0),
	}, nil
}

func AccountVersion1(string, string) (AccountOption, error) {
	return AccountOption{
		version: uint64(1),
	}, nil
}

func newAccount(private, address string, options ...AccountOptionFunc) (*Account, error) {
	var accountPlugin AccountPlugin
	version := uint64(0)
	for _, o := range options {
		opt, err := o(private, address)
		if err != nil {
			return nil, err
		}
		if opt.version != 0 {
			version = opt.version
		}
		if opt.AccountPlugin != nil {
			if accountPlugin != nil {
				return nil, errors.New("multiple plugins not supported")
			}
			accountPlugin = opt.AccountPlugin
		}
	}
	priv := types.SNValToBN(private)
	return &Account{
		AccountAddress: address,
		private:        priv,
		version:        version,
		plugin:         accountPlugin,
	}, nil
}

func NewRPCAccount(private, address string, provider *rpcv01.Provider, options ...AccountOptionFunc) (*Account, error) {
	account, err := newAccount(private, address, options...)
	if err != nil {
		return nil, err
	}
	chainID, err := provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	account.chainId = chainID
	account.provider = ProviderRPCv01
	account.rpcv01 = provider
	return account, nil
}

func NewGatewayAccount(private, address string, provider *gateway.GatewayProvider, options ...AccountOptionFunc) (*Account, error) {
	account, err := newAccount(private, address, options...)
	if err != nil {
		return nil, err
	}
	chainID, err := provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	account.chainId = chainID
	account.provider = ProviderGateway
	account.sequencer = provider
	return account, nil
}

func (account *Account) Call(ctx context.Context, call types.FunctionCall) ([]string, error) {
	switch account.provider {
	case ProviderRPCv01:
		if account.rpcv01 == nil {
			return nil, ErrUnsupportedAccount
		}
		return account.rpcv01.Call(ctx, call, rpcv01.WithBlockTag("latest"))
	case ProviderGateway:
		if account.sequencer == nil {
			return nil, ErrUnsupportedAccount
		}
		return account.sequencer.Call(ctx, call, "latest")
	}
	return nil, ErrUnsupportedAccount
}

func (account *Account) Sign(msgHash *big.Int) (*big.Int, *big.Int, error) {
	return Curve.Sign(msgHash, account.private)
}

func (account *Account) TransactionHash(calls []types.FunctionCall, details types.ExecuteDetails) (*big.Int, error) {

	var callArray []*big.Int
	switch {
	case account.version == 0:
		callArray = fmtV0Calldata(details.Nonce, calls)
	case account.version == 1:
		callArray = fmtCalldata(calls)
	default:
		return nil, fmt.Errorf("version %d unsupported", account.version)
	}
	cdHash, err := Curve.ComputeHashOnElements(callArray)
	if err != nil {
		return nil, err
	}

	var multiHashData []*big.Int
	switch {
	case account.version == 0:
		multiHashData = []*big.Int{
			types.UTF8StrToBig(TRANSACTION_PREFIX),
			big.NewInt(int64(account.version)),
			types.SNValToBN(account.AccountAddress),
			types.GetSelectorFromName(EXECUTE_SELECTOR),
			cdHash,
			details.MaxFee,
			types.UTF8StrToBig(account.chainId),
		}
	case account.version == 1:
		multiHashData = []*big.Int{
			types.UTF8StrToBig(TRANSACTION_PREFIX),
			big.NewInt(int64(account.version)),
			types.SNValToBN(account.AccountAddress),
			big.NewInt(0),
			cdHash,
			details.MaxFee,
			types.UTF8StrToBig(account.chainId),
			details.Nonce,
		}
	default:
		return nil, fmt.Errorf("version %d unsupported", account.version)
	}
	return Curve.ComputeHashOnElements(multiHashData)
}

func (account *Account) estimateFeeHash(calls []types.FunctionCall, details types.ExecuteDetails, version *big.Int) (*big.Int, error) {
	var callArray []*big.Int
	switch {
	case account.version == 0:
		callArray = fmtV0Calldata(details.Nonce, calls)
	case account.version == 1:
		callArray = fmtCalldata(calls)
	default:
		return nil, fmt.Errorf("version %d unsupported", account.version)
	}
	cdHash, err := Curve.ComputeHashOnElements(callArray)
	if err != nil {
		return nil, err
	}
	var multiHashData []*big.Int
	switch {
	case account.version == 0:
		multiHashData = []*big.Int{
			types.UTF8StrToBig(TRANSACTION_PREFIX),
			version,
			types.SNValToBN(account.AccountAddress),
			types.GetSelectorFromName(EXECUTE_SELECTOR),
			cdHash,
			details.MaxFee,
			types.UTF8StrToBig(account.chainId),
		}
	case account.version == 1:
		multiHashData = []*big.Int{
			types.UTF8StrToBig(TRANSACTION_PREFIX),
			version,
			types.SNValToBN(account.AccountAddress),
			big.NewInt(0),
			cdHash,
			details.MaxFee,
			types.UTF8StrToBig(account.chainId),
			details.Nonce,
		}
	default:
		return nil, fmt.Errorf("version %d unsupported", account.version)
	}
	return Curve.ComputeHashOnElements(multiHashData)
}

func (account *Account) Nonce(ctx context.Context) (*big.Int, error) {
	switch account.version {
	case 0:
		switch account.provider {
		case ProviderRPCv01:
			nonce, err := account.rpcv01.Call(
				ctx,
				types.FunctionCall{
					ContractAddress:    types.HexToHash(account.AccountAddress),
					EntryPointSelector: "get_nonce",
					Calldata:           []string{},
				},
				rpcv01.WithBlockTag("latest"),
			)
			if err != nil {
				return nil, err
			}
			if len(nonce) == 0 {
				return nil, errors.New("nonce error")
			}
			n, ok := big.NewInt(0).SetString(nonce[0], 0)
			if !ok {
				return nil, errors.New("nonce error")
			}
			return n, nil
		case ProviderGateway:
			return account.sequencer.AccountNonce(ctx, types.HexToHash(account.AccountAddress))
		}
	case 1:
		switch account.provider {
		case ProviderRPCv01:
			nonce, err := account.rpcv01.Nonce(
				ctx,
				types.HexToHash(account.AccountAddress),
			)
			if err != nil {
				return nil, err
			}
			n, ok := big.NewInt(0).SetString(*nonce, 0)
			if !ok {
				return nil, errors.New("nonce error")
			}
			return n, nil
		case ProviderGateway:
			return account.sequencer.Nonce(ctx, account.AccountAddress, "latest")
		}
	}
	return nil, fmt.Errorf("version %d unsupported", account.version)
}

func (account *Account) prepFunctionInvoke(ctx context.Context, messageType string, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FunctionInvoke, error) {
	if messageType != "invoke" && messageType != "estimate" {
		return nil, errors.New("unsupported message type")
	}
	nonce := details.Nonce
	var err error
	if details.Nonce == nil {
		nonce, err = account.Nonce(ctx)
		if err != nil {
			return nil, err
		}
	}
	maxFee := MAX_FEE
	if details.MaxFee != nil {
		maxFee = details.MaxFee
	}
	if account.plugin != nil {
		call, err := account.plugin.PluginCall(calls)
		if err != nil {
			return nil, err
		}
		calls = append([]types.FunctionCall{call}, calls...)
	}
	// version, _ := big.NewInt(0).SetString("0x100000000000000000000000000000000", 0)
	version, _ := big.NewInt(0).SetString("0x0", 0)
	var txHash *big.Int
	switch messageType {
	case "invoke":
		version = big.NewInt(int64(account.version))
		txHash, err = account.TransactionHash(
			calls,
			types.ExecuteDetails{
				Nonce:  nonce,
				MaxFee: maxFee,
			},
		)
		if err != nil {
			return nil, err
		}
	case "estimate":
		if account.version == 1 {
			// version, _ = big.NewInt(0).SetString("0x100000000000000000000000000000001", 0)
			version, _ = big.NewInt(0).SetString("0x1", 0)
		}
		txHash, err = account.estimateFeeHash(
			calls,
			types.ExecuteDetails{
				Nonce:  nonce,
				MaxFee: maxFee,
			},
			version,
		)
		if err != nil {
			return nil, err
		}
	}
	s1, s2, err := account.Sign(txHash)
	if err != nil {
		return nil, err
	}
	switch account.version {
	case 0:
		calldata := fmtV0CalldataStrings(nonce, calls)
		return &types.FunctionInvoke{
			MaxFee:    maxFee,
			Version:   version,
			Signature: types.Signature{s1, s2},
			FunctionCall: types.FunctionCall{
				ContractAddress:    types.HexToHash(account.AccountAddress),
				EntryPointSelector: EXECUTE_SELECTOR,
				Calldata:           calldata,
			},
		}, nil
	case 1:
		calldata := fmtCalldataStrings(calls)
		return &types.FunctionInvoke{
			MaxFee:    maxFee,
			Version:   version,
			Signature: types.Signature{s1, s2},
			FunctionCall: types.FunctionCall{
				ContractAddress: types.HexToHash(account.AccountAddress),
				Calldata:        calldata,
			},
			Nonce: nonce,
		}, nil
	}
	return nil, ErrUnsupportedAccount
}

func (account *Account) EstimateFee(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FeeEstimate, error) {
	call, err := account.prepFunctionInvoke(ctx, "estimate", calls, details)
	if err != nil {
		return nil, err
	}
	switch account.provider {
	case ProviderRPCv01:
		return account.rpcv01.EstimateFee(ctx, *call, rpcv01.WithBlockTag("latest"))
	case ProviderGateway:
		return account.sequencer.EstimateFee(ctx, *call, "latest")
	}
	return nil, ErrUnsupportedAccount
}

func (account *Account) Execute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.AddInvokeTransactionOutput, error) {
	maxFee := details.MaxFee
	if maxFee == nil {
		estimate, err := account.EstimateFee(ctx, calls, details)
		if err != nil {
			return nil, err
		}
		fmt.Printf("fee %+v\n", estimate)
		v, ok := big.NewInt(0).SetString(string(estimate.OverallFee), 0)
		if !ok {
			return nil, errors.New("could not match OverallFee to big.Int")
		}
		maxFee = v.Mul(v, big.NewInt(2))
	}
	details.MaxFee = maxFee
	call, err := account.prepFunctionInvoke(ctx, "invoke", calls, details)
	if err != nil {
		return nil, err
	}
	switch account.provider {
	case ProviderRPCv01:
		signature := []string{}
		for _, k := range call.Signature {
			signature = append(signature, fmt.Sprintf("0x%s", k.Text(16)))
		}
		return account.rpcv01.AddInvokeTransaction(
			context.Background(),
			call.FunctionCall,
			signature,
			fmt.Sprintf("0x%s", maxFee.Text(16)),
			fmt.Sprintf("0x%d", account.version),
			call.Nonce,
		)
	case ProviderGateway:
		return account.sequencer.Invoke(
			context.Background(),
			*call,
		)
	}
	return nil, ErrUnsupportedAccount
}
