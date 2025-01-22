package src

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	ksdeploy "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func ocrConfToContractConfig(ocrConf ksdeploy.OCR3OnchainConfig, configCount uint32) types.ContractConfig {
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

func mercuryOCRConfigToContractConfig(ocrConf MercuryOCR2Config, configCount uint32) types.ContractConfig {
	cc := types.ContractConfig{
		Signers:               convertAddressesToOnchainPublicKeys(ocrConf.Signers),
		Transmitters:          convertBytes32sToAccounts(ocrConf.Transmitters),
		F:                     ocrConf.F,
		OnchainConfig:         ocrConf.OnchainConfig,
		OffchainConfigVersion: ocrConf.OffchainConfigVersion,
		OffchainConfig:        ocrConf.OffchainConfig,
		ConfigCount:           uint64(configCount),
	}

	return cc
}

func convertAddressesToOnchainPublicKeys(addresses []common.Address) []types.OnchainPublicKey {
	keys := make([]types.OnchainPublicKey, len(addresses))
	for i, addr := range addresses {
		keys[i] = types.OnchainPublicKey(addr.Bytes())
	}
	return keys
}

func convertAddressesToAccounts(addresses []common.Address) []types.Account {
	accounts := make([]types.Account, len(addresses))
	for i, addr := range addresses {
		accounts[i] = types.Account(addr.Hex())
	}
	return accounts
}

func convertBytes32sToAccounts(bs [][32]byte) []types.Account {
	accounts := make([]types.Account, len(bs))
	for i, b := range bs {
		accounts[i] = types.Account(hex.EncodeToString(b[:]))
	}
	return accounts
}

func convertByteSliceToOnchainPublicKeys(bs [][]byte) []types.OnchainPublicKey {
	keys := make([]types.OnchainPublicKey, len(bs))
	for i, b := range bs {
		keys[i] = types.OnchainPublicKey(b)
	}
	return keys
}
