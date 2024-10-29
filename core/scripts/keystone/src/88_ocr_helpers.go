package src

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ksdeploy "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

func ocrConfToContractConfig(ocrConf ksdeploy.Orc2drOracleConfig, configCount uint32) types.ContractConfig {
	cc := types.ContractConfig{
		Signers:               convertByteSliceToOnchainPublicKeys(ocrConf.Signers),
		Transmitters:          convertAddressesToAccounts(ocrConf.Transmitters),
		F:                     ocrConf.F,
		OnchainConfig:         ocrConf.OnchainConfig,
		OffchainConfigVersion: ocrConf.OffchainConfigVersion,
		OffchainConfig:        ocrConf.OffchainConfig,
		ConfigCount:           uint64(configCount),
	}
	return cc
}

func convertAddressesToAccounts(addresses []common.Address) []types.Account {
	accounts := make([]types.Account, len(addresses))
	for i, addr := range addresses {
		accounts[i] = types.Account(addr.Hex())
	}
	return accounts
}

func convertByteSliceToOnchainPublicKeys(bs [][]byte) []types.OnchainPublicKey {
	keys := make([]types.OnchainPublicKey, len(bs))
	for i, b := range bs {
		keys[i] = types.OnchainPublicKey(hexutil.Encode(b))
	}
	return keys
}
