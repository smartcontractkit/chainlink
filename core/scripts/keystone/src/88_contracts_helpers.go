package src

import (
	"context"

	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"
)

var ZeroAddress = common.Address{}

type OnChainMetaSerialized struct {
	OCRContract       common.Address `json:"ocrContract"`
	ForwarderContract common.Address `json:"forwarderContract"`
	// The block number of the transaction that set the config on the OCR3 contract. We use this to replay blocks from this point on
	// when we load the OCR3 job specs on the nodes.
	SetConfigTxBlock uint64 `json:"setConfigTxBlock"`

	CapabilitiesRegistry common.Address `json:"CapabilitiesRegistry"`
}

type onchainMeta struct {
	OCRContract       ocr3_capability.OCR3CapabilityInterface
	ForwarderContract forwarder.KeystoneForwarderInterface
	// The block number of the transaction that set the config on the OCR3 contract. We use this to replay blocks from this point on
	// when we load the OCR3 job specs on the nodes.
	SetConfigTxBlock uint64

	CapabilitiesRegistry capabilities_registry.CapabilitiesRegistryInterface
}

func WriteOnchainMeta(o onchainMeta, artefactsDir string) {
	_, err := os.Stat(artefactsDir)
	if err != nil {
		fmt.Println("Creating artefacts directory")
		err = os.MkdirAll(artefactsDir, 0700)
		PanicErr(err)
	}

	fmt.Println("Writing deployed contract addresses to file...")
	serialzed := OnChainMetaSerialized{
		OCRContract:          o.OCRContract.Address(),
		ForwarderContract:    o.ForwarderContract.Address(),
		SetConfigTxBlock:     o.SetConfigTxBlock,
		CapabilitiesRegistry: o.CapabilitiesRegistry.Address(),
	}

	jsonBytes, err := json.Marshal(serialzed)
	PanicErr(err)

	err = os.WriteFile(deployedContractsFilePath(artefactsDir), jsonBytes, 0600)
	PanicErr(err)
}

func LoadOnchainMeta(artefactsDir string, env helpers.Environment) onchainMeta {
	if !ContractsAlreadyDeployed(artefactsDir) {
		fmt.Printf("No deployed contracts file found at %s\n", deployedContractsFilePath(artefactsDir))
		return onchainMeta{}
	}

	jsonBytes, err := os.ReadFile(deployedContractsFilePath(artefactsDir))
	if err != nil {
		fmt.Printf("Error reading deployed contracts file: %s\n", err)
		return onchainMeta{}
	}

	var s OnChainMetaSerialized
	err = json.Unmarshal(jsonBytes, &s)
	if err != nil {
		return onchainMeta{}
	}

	hydrated := onchainMeta{
		SetConfigTxBlock: s.SetConfigTxBlock,
	}

	if s.OCRContract != ZeroAddress {
		if !contractExists(s.OCRContract, env) {
			fmt.Printf("OCR contract at %s does not exist, setting to zero address\n", s.OCRContract.Hex())
			s.OCRContract = ZeroAddress
		}

		ocr3, err := ocr3_capability.NewOCR3Capability(s.OCRContract, env.Ec)
		PanicErr(err)
		hydrated.OCRContract = ocr3
	}

	if s.ForwarderContract != ZeroAddress {
		if !contractExists(s.ForwarderContract, env) {
			fmt.Printf("Forwarder contract at %s does not exist, setting to zero address\n", s.ForwarderContract.Hex())
			s.ForwarderContract = ZeroAddress
		}

		fwdr, err := forwarder.NewKeystoneForwarder(s.ForwarderContract, env.Ec)
		PanicErr(err)
		hydrated.ForwarderContract = fwdr
	}

	if s.CapabilitiesRegistry != ZeroAddress {
		if !contractExists(s.CapabilitiesRegistry, env) {
			fmt.Printf("CapabilityRegistry contract at %s does not exist, setting to zero address\n", s.CapabilitiesRegistry.Hex())
			s.CapabilitiesRegistry = ZeroAddress
		}

		cr, err := capabilities_registry.NewCapabilitiesRegistry(s.CapabilitiesRegistry, env.Ec)
		PanicErr(err)
		hydrated.CapabilitiesRegistry = cr
	}

	blkNum, err := env.Ec.BlockNumber(context.Background())
	PanicErr(err)

	if s.SetConfigTxBlock > blkNum {
		fmt.Printf("Stale SetConfigTxBlock: %d, current block number: %d\n", s.SetConfigTxBlock, blkNum)

		return onchainMeta{}
	}

	return hydrated
}

func ContractsAlreadyDeployed(artefactsDir string) bool {
	_, err := os.Stat(artefactsDir)
	if err != nil {
		return false
	}

	_, err = os.Stat(deployedContractsFilePath(artefactsDir))
	if err != nil {
		return false
	}

	return true
}

func deployedContractsFilePath(artefactsDir string) string {
	return filepath.Join(artefactsDir, deployedContractsJSON)
}

func contractExists(address common.Address, env helpers.Environment) bool {
	byteCode, err := env.Ec.CodeAt(context.Background(), address, nil)
	if err != nil {
		return false
	}
	return len(byteCode) != 0
}
