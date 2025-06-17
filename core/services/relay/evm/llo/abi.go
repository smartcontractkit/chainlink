package llo

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
)

var (
	ProductionConfigSet  common.Hash
	StagingConfigSet     common.Hash
	PromoteStagingConfig common.Hash

	configuratorABI abi.ABI
)

func init() {
	var err error
	configuratorABI, err = abi.JSON(strings.NewReader(configurator.ConfiguratorABI))
	if err != nil {
		panic(err)
	}
	ProductionConfigSet = configuratorABI.Events["ProductionConfigSet"].ID
	StagingConfigSet = configuratorABI.Events["StagingConfigSet"].ID
	PromoteStagingConfig = configuratorABI.Events["PromoteStagingConfig"].ID
}
