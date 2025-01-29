package v1_5

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

type TokenPoolContract interface {
	Address() common.Address
	Owner(opts *bind.CallOpts) (common.Address, error)
	TypeAndVersion(*bind.CallOpts) (string, error)
	GetToken(opts *bind.CallOpts) (common.Address, error)
	GetSupportedChains(opts *bind.CallOpts) ([]uint64, error)
	GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error)
	GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)
}

type RemoteChainConfig struct {
	// RemoteTokenAddress is a hex representation of byte arrays
	RemoteTokenAddress string
	// RemotePoolAddresses are a hex representations of byte arrays
	RemotePoolAddresses []string
}

type TokenPoolView struct {
	types.ContractMetaData
	Token              common.Address               `json:"token"`
	RemoteChainConfigs map[uint64]RemoteChainConfig `json:"remoteChainConfigs"`
}

func GenerateTokenPoolView(pool TokenPoolContract) (TokenPoolView, error) {
	owner, err := pool.Owner(nil)
	if err != nil {
		return TokenPoolView{}, err
	}
	typeAndVersion, err := pool.TypeAndVersion(nil)
	if err != nil {
		return TokenPoolView{}, err
	}
	token, err := pool.GetToken(nil)
	if err != nil {
		return TokenPoolView{}, err
	}
	remoteChains, err := pool.GetSupportedChains(nil)
	if err != nil {
		return TokenPoolView{}, err
	}
	remoteChainConfigs := make(map[uint64]RemoteChainConfig)
	for _, remoteChain := range remoteChains {
		remotePools, err := pool.GetRemotePools(nil, remoteChain)
		if err != nil {
			return TokenPoolView{}, err
		}
		remoteToken, err := pool.GetRemoteToken(nil, remoteChain)
		if err != nil {
			return TokenPoolView{}, err
		}
		remotePoolsHex := make([]string, len(remotePools))
		for i, remotePool := range remotePools {
			remotePoolsHex[i] = hex.EncodeToString(remotePool)
		}
		remoteChainConfigs[remoteChain] = RemoteChainConfig{
			RemoteTokenAddress:  hexutil.Encode(remoteToken),
			RemotePoolAddresses: remotePoolsHex,
		}
	}
	return TokenPoolView{
		ContractMetaData: types.ContractMetaData{
			TypeAndVersion: typeAndVersion,
			Address:        pool.Address(),
			Owner:          owner,
		},
		Token:              token,
		RemoteChainConfigs: remoteChainConfigs,
	}, nil
}
