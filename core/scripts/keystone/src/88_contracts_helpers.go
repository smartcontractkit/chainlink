package src

import (
	"context"

	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	capabilities_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"
	verifierContract "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy"
)

var ZeroAddress = common.Address{}

type OnChainMetaSerialized struct {
	OCR       common.Address `json:"ocrContract"`
	Forwarder common.Address `json:"forwarderContract"`
	// The block number of the transaction that set the config on the OCR3 contract. We use this to replay blocks from this point on
	// when we load the OCR3 job specs on the nodes.
	SetConfigTxBlock uint64 `json:"setConfigTxBlock"`

	CapabilitiesRegistry common.Address `json:"CapabilitiesRegistry"`
	Verifier             common.Address `json:"VerifierContract"`
	VerifierProxy        common.Address `json:"VerifierProxy"`
	// Stores the address that has been initialized by the proxy, if any
	InitializedVerifierAddress common.Address `json:"InitializedVerifierAddress"`
}

type onchainMeta struct {
	OCR3      ocr3_capability.OCR3CapabilityInterface
	Forwarder forwarder.KeystoneForwarderInterface
	// The block number of the transaction that set the config on the OCR3 contract. We use this to replay blocks from this point on
	// when we load the OCR3 job specs on the nodes.
	SetConfigTxBlock uint64

	CapabilitiesRegistry       capabilities_registry.CapabilitiesRegistryInterface
	Verifier                   verifierContract.VerifierInterface
	VerifierProxy              verifier_proxy.VerifierProxyInterface
	InitializedVerifierAddress common.Address `json:"InitializedVerifierAddress"`
}

func WriteOnchainMeta(o *onchainMeta, artefactsDir string) {
	ensureArtefactsDir(artefactsDir)

	fmt.Println("Writing deployed contract addresses to file...")
	serialzed := OnChainMetaSerialized{}

	if o.OCR3 != nil {
		serialzed.OCR = o.OCR3.Address()
	}

	if o.Forwarder != nil {
		serialzed.Forwarder = o.Forwarder.Address()
	}

	serialzed.SetConfigTxBlock = o.SetConfigTxBlock
	serialzed.InitializedVerifierAddress = o.InitializedVerifierAddress

	if o.CapabilitiesRegistry != nil {
		serialzed.CapabilitiesRegistry = o.CapabilitiesRegistry.Address()
	}

	if o.Verifier != nil {
		serialzed.Verifier = o.Verifier.Address()
	}

	if o.VerifierProxy != nil {
		serialzed.VerifierProxy = o.VerifierProxy.Address()
	}

	jsonBytes, err := json.Marshal(serialzed)
	PanicErr(err)

	err = os.WriteFile(deployedContractsFilePath(artefactsDir), jsonBytes, 0600)
	PanicErr(err)
}

func LoadOnchainMeta(artefactsDir string, env helpers.Environment) *onchainMeta {
	hydrated := &onchainMeta{}
	if !ContractsAlreadyDeployed(artefactsDir) {
		fmt.Printf("No deployed contracts file found at %s\n", deployedContractsFilePath(artefactsDir))
		return hydrated
	}

	jsonBytes, err := os.ReadFile(deployedContractsFilePath(artefactsDir))
	if err != nil {
		fmt.Printf("Error reading deployed contracts file: %s\n", err)
		return hydrated
	}

	var s OnChainMetaSerialized
	err = json.Unmarshal(jsonBytes, &s)
	if err != nil {
		return hydrated
	}

	hydrated.SetConfigTxBlock = s.SetConfigTxBlock
	if s.OCR != ZeroAddress {
		if !contractExists(s.OCR, env) {
			fmt.Printf("OCR contract at %s does not exist\n", s.OCR.Hex())
		} else {
			ocr3, e := ocr3_capability.NewOCR3Capability(s.OCR, env.Ec)
			PanicErr(e)
			hydrated.OCR3 = ocr3
		}
	}

	if s.Forwarder != ZeroAddress {
		if !contractExists(s.Forwarder, env) {
			fmt.Printf("Forwarder contract at %s does not exist\n", s.Forwarder.Hex())
		} else {
			fwdr, e := forwarder.NewKeystoneForwarder(s.Forwarder, env.Ec)
			PanicErr(e)
			hydrated.Forwarder = fwdr
		}
	}

	if s.CapabilitiesRegistry != ZeroAddress {
		if !contractExists(s.CapabilitiesRegistry, env) {
			fmt.Printf("CapabilityRegistry contract at %s does not exist\n", s.CapabilitiesRegistry.Hex())
		} else {
			cr, e := capabilities_registry.NewCapabilitiesRegistry(s.CapabilitiesRegistry, env.Ec)
			PanicErr(e)
			hydrated.CapabilitiesRegistry = cr
		}
	}

	hydrated.InitializedVerifierAddress = s.InitializedVerifierAddress

	if s.Verifier != ZeroAddress {
		if !contractExists(s.Verifier, env) {
			fmt.Printf("Verifier contract at %s does not exist\n", s.Verifier.Hex())
			hydrated.InitializedVerifierAddress = ZeroAddress
		} else {
			verifier, e := verifierContract.NewVerifier(s.Verifier, env.Ec)
			PanicErr(e)
			hydrated.Verifier = verifier
		}
	}

	if s.VerifierProxy != ZeroAddress {
		if !contractExists(s.VerifierProxy, env) {
			fmt.Printf("VerifierProxy contract at %s does not exist\n", s.VerifierProxy.Hex())
			hydrated.InitializedVerifierAddress = ZeroAddress
		} else {
			verifierProxy, e := verifier_proxy.NewVerifierProxy(s.VerifierProxy, env.Ec)
			PanicErr(e)
			hydrated.VerifierProxy = verifierProxy
		}
	}

	blkNum, err := env.Ec.BlockNumber(context.Background())
	PanicErr(err)

	if s.SetConfigTxBlock > blkNum {
		fmt.Printf("Stale SetConfigTxBlock: %d, current block number: %d\n", s.SetConfigTxBlock, blkNum)
		hydrated.SetConfigTxBlock = 0
	}

	return hydrated
}

func ContractsAlreadyDeployed(artefactsDir string) bool {
	_, err := os.Stat(artefactsDir)
	if err != nil {
		return false
	}

	_, err = os.Stat(deployedContractsFilePath(artefactsDir))

	return err == nil
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
