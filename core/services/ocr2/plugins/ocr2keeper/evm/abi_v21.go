package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
)

type evmRegistryPackerV21 struct {
	abi abi.ABI
}

func NewEvmRegistryPackerV21(abi abi.ABI) *evmRegistryPackerV21 {
	return &evmRegistryPackerV21{abi: abi}
}

// TODO: implement other methods as needed

// UnpackLogTriggerConfig unpacks the log trigger config from the given raw data
// TODO: tests
func (rp *evmRegistryPackerV21) UnpackLogTriggerConfig(raw string) (i_keeper_registry_master_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig, error) {
	var cfg i_keeper_registry_master_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig
	b, err := hexutil.Decode(raw)
	if err != nil {
		return cfg, err
	}

	out, err := rp.abi.Methods["getLogTriggerConfig"].Outputs.UnpackValues(b)
	if err != nil {
		return cfg, fmt.Errorf("%w: unpack getLogTriggerConfig return: %s", err, raw)
	}

	converted, ok := abi.ConvertType(out[0], new(i_keeper_registry_master_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig)).(*i_keeper_registry_master_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig)
	if !ok {
		return cfg, fmt.Errorf("failed to convert type")
	}
	return *converted, nil
}
