package v1_6

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_0_0/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type SetRMNRemoteConfig struct {
	ChainSelector   uint64                        `json:"chainSelector"`
	RMNRemoteConfig RMNRemoteConfig               `json:"rmnRemoteConfigs"`
	MCMSConfig      *proposalutils.TimelockConfig `json:"mcmsConfig,omitempty"`
}

type RMNRemoteConfig struct {
	Signers []rmn_remote.RMNRemoteSigner `json:"signers"`
	F       uint64                       `json:"f"`
}

type DeployRMNRemoteInput struct {
	RMNLegacyAddr common.Address `json:"rmnLegacyAddr"`
	ChainSelector uint64         `json:"chainSelector"`
}

var (
	DeployRMNRemoteOp = opsutil.NewEVMDeployOperation(
		"DeployRMNRemote",
		semver.MustParse("1.0.0"),
		"Deploys RMNRemote 1.6 contract on the specified evm chain",
		shared.RMNRemote,
		rmn_remote.RMNRemoteMetaData,
		&opsutil.ContractOpts{
			Version:          &deployment.Version1_6_0,
			EVMBytecode:      common.FromHex(rmn_remote.RMNRemoteBin),
			ZkSyncVMBytecode: rmn_remote.ZkBytecode,
		},
		func(input DeployRMNRemoteInput) []any {
			return []any{
				input.ChainSelector,
				input.RMNLegacyAddr,
			}
		},
	)

	SetRMNRemoteConfigOp = opsutil.NewEVMCallOperation(
		"SetRMNRemoteConfigOp",
		semver.MustParse("1.0.0"),
		"Setting RMNRemoteConfig based on ActiveDigest from RMNHome",
		rmn_remote.RMNRemoteABI,
		shared.RMNRemote,
		rmn_remote.NewRMNRemote,
		func(rmnRemote *rmn_remote.RMNRemote, opts *bind.TransactOpts, input rmn_remote.RMNRemoteConfig) (*types.Transaction, error) {
			return rmnRemote.SetConfig(opts, input)
		})

	SetRMNRemoteOnRMNProxyOp = opsutil.NewEVMCallOperation(
		"SetRMNRemoteOnRMNProxyOp",
		semver.MustParse("1.0.0"),
		"Sets SetRMNRemote on RMNProxy contract on the specified evm chain",
		rmn_proxy_contract.RMNProxyABI,
		shared.ARMProxy,
		rmn_proxy_contract.NewRMNProxy,
		func(rmnProxy *rmn_proxy_contract.RMNProxy, opts *bind.TransactOpts, input common.Address) (*types.Transaction, error) {
			return rmnProxy.SetARM(opts, input)
		},
	)
)
