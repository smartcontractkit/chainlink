package contracts

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	contractsethereum "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_registry_master_wrapper_2_2"
	iregistry22 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_registry_master_wrapper_2_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
)

// OCRv2Config represents the config for the OCRv2 contract
type OCRv2Config struct {
	Signers               []common.Address
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	TypedOnchainConfig21  i_keeper_registry_master_wrapper_2_1.IAutomationV21PlusCommonOnchainConfigLegacy
	TypedOnchainConfig22  i_automation_registry_master_wrapper_2_2.AutomationRegistryBase22OnchainConfig
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

type EthereumFunctionsLoadStats struct {
	LastRequestID string
	LastResponse  string
	LastError     string
	Total         uint32
	Succeeded     uint32
	Errored       uint32
	Empty         uint32
}

func Bytes32ToSlice(a [32]byte) (r []byte) {
	r = append(r, a[:]...)
	return
}

// DefaultOffChainAggregatorOptions returns some base defaults for deploying an OCR contract
func DefaultOffChainAggregatorOptions() OffchainOptions {
	return OffchainOptions{
		MaximumGasPrice:         uint32(3000),
		ReasonableGasPrice:      uint32(10),
		MicroLinkPerEth:         uint32(500),
		LinkGweiPerObservation:  uint32(500),
		LinkGweiPerTransmission: uint32(500),
		MinimumAnswer:           big.NewInt(1),
		MaximumAnswer:           big.NewInt(50000000000000000),
		Decimals:                8,
		Description:             "Test OCR",
	}
}

// DefaultOffChainAggregatorConfig returns some base defaults for configuring an OCR contract
func DefaultOffChainAggregatorConfig(numberNodes int) OffChainAggregatorConfig {
	if numberNodes <= 4 {
		log.Err(fmt.Errorf("insufficient number of nodes (%d) supplied for OCR, need at least 5", numberNodes)).
			Int("Number Chainlink Nodes", numberNodes).
			Msg("You likely need more chainlink nodes to properly configure OCR, try 5 or more.")
	}
	s := []int{1}
	// First node's stage already inputted as a 1 in line above, so numberNodes-1.
	for i := 0; i < numberNodes-1; i++ {
		s = append(s, 2)
	}
	return OffChainAggregatorConfig{
		AlphaPPB:         1,
		DeltaC:           time.Minute * 60,
		DeltaGrace:       time.Second * 12,
		DeltaProgress:    time.Second * 35,
		DeltaStage:       time.Second * 60,
		DeltaResend:      time.Second * 17,
		DeltaRound:       time.Second * 30,
		RMax:             6,
		S:                s,
		N:                numberNodes,
		F:                1,
		OracleIdentities: []ocrConfigHelper.OracleIdentityExtra{},
	}
}

func ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(k8sNodes []*client.ChainlinkK8sClient) []ChainlinkNodeWithKeysAndAddress {
	var nodesAsInterface = make([]ChainlinkNodeWithKeysAndAddress, len(k8sNodes))
	for i, node := range k8sNodes {
		nodesAsInterface[i] = node
	}

	return nodesAsInterface
}

func ChainlinkClientToChainlinkNodeWithKeysAndAddress(k8sNodes []*client.ChainlinkClient) []ChainlinkNodeWithKeysAndAddress {
	var nodesAsInterface = make([]ChainlinkNodeWithKeysAndAddress, len(k8sNodes))
	for i, node := range k8sNodes {
		nodesAsInterface[i] = node
	}

	return nodesAsInterface
}

func V2OffChainAgrregatorToOffChainAggregatorWithRounds(contracts []OffchainAggregatorV2) []OffChainAggregatorWithRounds {
	var contractsAsInterface = make([]OffChainAggregatorWithRounds, len(contracts))
	for i, contract := range contracts {
		contractsAsInterface[i] = contract
	}

	return contractsAsInterface
}

func V1OffChainAgrregatorToOffChainAggregatorWithRounds(contracts []OffchainAggregator) []OffChainAggregatorWithRounds {
	var contractsAsInterface = make([]OffChainAggregatorWithRounds, len(contracts))
	for i, contract := range contracts {
		contractsAsInterface[i] = contract
	}

	return contractsAsInterface
}

func GetRegistryContractABI(version contractsethereum.KeeperRegistryVersion) (*abi.ABI, error) {
	var (
		contractABI *abi.ABI
		err         error
	)
	switch version {
	case contractsethereum.RegistryVersion_1_0, contractsethereum.RegistryVersion_1_1:
		contractABI, err = keeper_registry_wrapper1_1.KeeperRegistryMetaData.GetAbi()
	case contractsethereum.RegistryVersion_1_2:
		contractABI, err = keeper_registry_wrapper1_2.KeeperRegistryMetaData.GetAbi()
	case contractsethereum.RegistryVersion_1_3:
		contractABI, err = keeper_registry_wrapper1_3.KeeperRegistryMetaData.GetAbi()
	case contractsethereum.RegistryVersion_2_0:
		contractABI, err = keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	case contractsethereum.RegistryVersion_2_1:
		contractABI, err = iregistry21.IKeeperRegistryMasterMetaData.GetAbi()
	case contractsethereum.RegistryVersion_2_2:
		contractABI, err = iregistry22.IAutomationRegistryMasterMetaData.GetAbi()
	default:
		contractABI, err = keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	}

	return contractABI, err
}

// DefaultFluxAggregatorOptions produces some basic defaults for a flux aggregator contract
func DefaultFluxAggregatorOptions() FluxAggregatorOptions {
	return FluxAggregatorOptions{
		PaymentAmount: big.NewInt(1),
		Timeout:       uint32(30),
		MinSubValue:   big.NewInt(0),
		MaxSubValue:   big.NewInt(1000000000000),
		Decimals:      uint8(0),
		Description:   "Test Flux Aggregator",
	}
}
