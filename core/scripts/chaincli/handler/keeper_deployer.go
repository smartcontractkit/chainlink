package handler

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/core/cmd"
	registry11 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry12 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	registry20 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
)

// canceller describes the behavior to cancel upkeeps
type canceller interface {
	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)
	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)
	RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error)
}

// upkeepDeployer contains functions needed to deploy an upkeep
type upkeepDeployer interface {
	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error)
	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)
}

// keepersDeployer contains functions needed to deploy keepers
type keepersDeployer interface {
	canceller
	upkeepDeployer
	SetKeepers(opts *bind.TransactOpts, _ []cmd.HTTPClient, keepers []common.Address, payees []common.Address) (*types.Transaction, error)
}

type v11KeeperDeployer struct {
	registry11.KeeperRegistryInterface
}

func (d *v11KeeperDeployer) SetKeepers(opts *bind.TransactOpts, _ []cmd.HTTPClient, keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return d.KeeperRegistryInterface.SetKeepers(opts, keepers, payees)
}

type v12KeeperDeployer struct {
	registry12.KeeperRegistryInterface
}

func (d *v12KeeperDeployer) SetKeepers(opts *bind.TransactOpts, _ []cmd.HTTPClient, keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return d.KeeperRegistryInterface.SetKeepers(opts, keepers, payees)
}

type v20KeeperDeployer struct {
	registry20.KeeperRegistryInterface
	cfg *config.Config
}

func (d *v20KeeperDeployer) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return d.KeeperRegistryInterface.RegisterUpkeep(opts, target, gasLimit, admin, checkData, nil)
}

type OCR2KeeperNode struct {
	Id             string `json:"id"`
	OnChainPubKey  string `json:"onChainPubKey"`
	OffChainPubKey string `json:"offChainPubKey"`
	ConfigPubKey   string `json:"configPubKey"`
	PeerID         string `json:"peerID"`
	KeeperAddress  string `json:"keeperAddress"`
	PayeeAddress   string `json:"payeeAddress"`
}

func (d *v20KeeperDeployer) SetKeepers(opts *bind.TransactOpts, cls []cmd.HTTPClient, keepers []common.Address, _ []common.Address) (*types.Transaction, error) {
	nodes := make([]OCR2KeeperNode, len(cls))
	f, _ := os.Create("data")

	for i, cl := range cls {
		node := OCR2KeeperNode{}
		ocr2Config, err := getNodeOCR2Config(cl)
		if err != nil {
			panic(err)
		}
		p2pKeyID, err := getP2PKeyID(cl)
		if err != nil {
			panic(err)
		}

		node.Id = ocr2Config.ID
		node.PeerID = p2pKeyID
		node.OffChainPubKey = strings.TrimPrefix(ocr2Config.OffChainPublicKey, "ocr2off_evm_")
		node.ConfigPubKey = strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_")
		node.OnChainPubKey = strings.TrimPrefix(ocr2Config.OnchainPublicKey, "ocr2on_evm_")
		node.KeeperAddress = keepers[i].String()
		node.PayeeAddress = "0x8fA510072009E71CfD447169AB5A84cAc394f58A"
		nodes[i] = node
		//fmt.Println(node)
		//j, err := json.Marshal(node)
		//if err != nil {
		//	continue
		//}

		nodeBytes, err := json.MarshalIndent(node, "", "  ")
		if err != nil {
			continue
		}
		f.Write(nodeBytes)
		f.WriteString(",\n")
		fmt.Println(string(nodeBytes))
	}

	return nil, nil

	//signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err := ocr2config.ContractSetConfigArgsForTests(
	//	10*time.Second,        // deltaProgress time.Duration,
	//	15*time.Second,        // deltaResend time.Duration,
	//	3000*time.Millisecond, // deltaRound time.Duration,
	//	50*time.Millisecond,   // deltaGrace time.Duration,
	//	90*time.Second,        // deltaStage time.Duration,
	//	20,                    // rMax uint8,
	//	S,                     // s []int,
	//	oracleIdentities,      // oracles []OracleIdentityExtra,
	//	ocr2keepers.OffchainConfig{
	//		PerformLockoutWindow: 100 * 12 * 1000, // ~100 block lockout (on goerli)
	//		UniqueReports:        false,           // set quorum requirements
	//		TargetProbability:    "0.99",
	//		TargetInRounds:       4,
	//	}.Encode(), // reportingPluginConfig []byte,
	//	15*time.Millisecond,   // maxDurationQuery time.Duration,
	//	1900*time.Millisecond, // maxDurationObservation time.Duration,
	//	900*time.Millisecond,  // maxDurationReport time.Duration,
	//	15*time.Millisecond,   // maxDurationShouldAcceptFinalizedReport time.Duration,
	//	15*time.Millisecond,   // maxDurationShouldTransmitAcceptedReport time.Duration,
	//	1,                     // f int,
	//	nil,                   // onchainConfig []byte,
	//)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var signers []common.Address
	//for _, signer := range signerOnchainPublicKeys {
	//	if len(signer) != 20 {
	//		return nil, fmt.Errorf("OnChainPublicKey has wrong length for address")
	//	}
	//	signers = append(signers, common.BytesToAddress(signer))
	//}
	//
	//var transmitters []common.Address
	//for _, transmitter := range transmitterAccounts {
	//	if !common.IsHexAddress(string(transmitter)) {
	//		return nil, fmt.Errorf("TransmitAccount is not a valid Ethereum address")
	//	}
	//	transmitters = append(transmitters, common.HexToAddress(string(transmitter)))
	//}
	//
	//configType := abi.MustNewType("tuple(uint32 paymentPremiumPPB,uint32 flatFeeMicroLink,uint32 checkGasLimit,uint24 stalenessSeconds,uint16 gasCeilingMultiplier,uint96 minUpkeepSpend,uint32 maxPerformGas,uint32 maxCheckDataSize,uint32 maxPerformDataSize,uint256 fallbackGasPrice,uint256 fallbackLinkPrice,address transcoder,address registrar)")
	//onchainConfig, err := abi.Encode(map[string]interface{}{
	//	"paymentPremiumPPB":    d.cfg.PaymentPremiumPBB,
	//	"flatFeeMicroLink":     d.cfg.FlatFeeMicroLink,
	//	"checkGasLimit":        d.cfg.CheckGasLimit,
	//	"stalenessSeconds":     d.cfg.StalenessSeconds,
	//	"gasCeilingMultiplier": d.cfg.GasCeilingMultiplier,
	//	"minUpkeepSpend":       d.cfg.MinUpkeepSpend,
	//	"maxPerformGas":        d.cfg.MaxPerformGas,
	//	"maxCheckDataSize":     d.cfg.MaxCheckDataSize,
	//	"maxPerformDataSize":   d.cfg.MaxPerformDataSize,
	//	"fallbackGasPrice":     big.NewInt(d.cfg.FallbackGasPrice),
	//	"fallbackLinkPrice":    big.NewInt(d.cfg.FallbackLinkPrice),
	//	"transcoder":           common.HexToAddress(d.cfg.Transcoder),
	//	"registrar":            common.HexToAddress(d.cfg.Registrar),
	//}, configType)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return d.KeeperRegistryInterface.SetConfig(opts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}
