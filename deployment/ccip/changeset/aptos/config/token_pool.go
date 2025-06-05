package config

import (
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/common"

	fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type AddTokenPoolConfig struct {
	// DeployAptosTokenConfig
	ChainSelector                       uint64
	TokenAddress                        aptos.AccountAddress // if empty, token will be deployed
	TokenObjAddress                     aptos.AccountAddress // if empty, token will be deployed
	TokenPoolAddress                    aptos.AccountAddress // if empty, token pool will be deployed
	PoolType                            cldf.ContractType
	TokenTransferFeeByRemoteChainConfig map[uint64]fee_quoter.TokenTransferFeeConfig
	EVMRemoteConfigs                    map[uint64]EVMRemoteConfig
	TokenParams                         TokenParams
	MCMSConfig                          *proposalutils.TimelockConfig
}

type EVMRemoteConfig struct {
	TokenAddress common.Address
	// TODO: EVM has a way of picking up Pool by token address and type, use this instead of passing PoolAddress
	TokenPoolAddress common.Address
	RateLimiterConfig
}

type RateLimiterConfig struct {
	RemoteChainSelector uint64
	OutboundIsEnabled   bool
	OutboundCapacity    uint64
	OutboundRate        uint64
	InboundIsEnabled    bool
	InboundCapacity     uint64
	InboundRate         uint64
}
