package handler

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ocr2config "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/types"
	"github.com/umbracle/ethgo/abi"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	registry11 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry12 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	registry20 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
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

func (d *v20KeeperDeployer) SetKeepers(opts *bind.TransactOpts, cls []cmd.HTTPClient, keepers []common.Address, _ []common.Address) (*types.Transaction, error) {
	S := make([]int, len(cls))
	oracleIdentities := make([]ocr2config.OracleIdentityExtra, len(cls))
	sharedSecretEncryptionPublicKeys := make([]ocr2types.ConfigEncryptionPublicKey, len(cls))
	var wg sync.WaitGroup
	for i, cl := range cls {
		wg.Add(1)
		go func(i int, cl cmd.HTTPClient) {
			defer wg.Done()

			ocr2Config, err := getNodeOCR2Config(cl)
			if err != nil {
				panic(err)
			}

			p2pKeyID, err := getP2PKeyID(cl)
			if err != nil {
				panic(err)
			}

			offchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OffChainPublicKey, "ocr2off_evm_"))
			if err != nil {
				panic(fmt.Errorf("failed to decode %s: %v", ocr2Config.OffChainPublicKey, err))
			}

			offchainPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n := copy(offchainPkBytesFixed[:], offchainPkBytes)
			if n != ed25519.PublicKeySize {
				panic(fmt.Errorf("wrong num elements copied"))
			}

			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_"))
			if err != nil {
				panic(fmt.Errorf("failed to decode %s: %v", ocr2Config.ConfigPublicKey, err))
			}

			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n = copy(configPkBytesFixed[:], configPkBytes)
			if n != ed25519.PublicKeySize {
				panic(fmt.Errorf("wrong num elements copied"))
			}

			onchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OnchainPublicKey, "ocr2on_evm_"))
			if err != nil {
				panic(fmt.Errorf("failed to decode %s: %v", ocr2Config.OnchainPublicKey, err))
			}

			sharedSecretEncryptionPublicKeys[i] = configPkBytesFixed
			oracleIdentities[i] = ocr2config.OracleIdentityExtra{
				OracleIdentity: ocr2config.OracleIdentity{
					OnchainPublicKey:  onchainPkBytes,
					OffchainPublicKey: offchainPkBytesFixed,
					PeerID:            p2pKeyID,
					TransmitAccount:   ocr2types.Account(keepers[i].String()),
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}
			S[i] = 1
		}(i, cl)
	}
	wg.Wait()

	signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err := ocr2config.ContractSetConfigArgsForTests(
		5*time.Second,         // deltaProgress time.Duration,
		10*time.Second,        // deltaResend time.Duration,
		2500*time.Millisecond, // deltaRound time.Duration,
		40*time.Millisecond,   // deltaGrace time.Duration,
		30*time.Second,        // deltaStage time.Duration,
		50,                    // rMax uint8,
		S,                     // s []int,
		oracleIdentities,      // oracles []OracleIdentityExtra,
		ocr2keepers.OffchainConfig{
			PerformLockoutWindow: 100 * 3 * 1000, // ~100 block lockout (on mumbai)
			MinConfirmations:     1,
		}.Encode(), // reportingPluginConfig []byte,
		20*time.Millisecond,   // maxDurationQuery time.Duration,
		1600*time.Millisecond, // maxDurationObservation time.Duration,
		800*time.Millisecond,  // maxDurationReport time.Duration,
		20*time.Millisecond,   // maxDurationShouldAcceptFinalizedReport time.Duration,
		20*time.Millisecond,   // maxDurationShouldTransmitAcceptedReport time.Duration,
		1,                     // f int,
		nil,                   // onchainConfig []byte,
	)
	if err != nil {
		return nil, err
	}

	var signers []common.Address
	for _, signer := range signerOnchainPublicKeys {
		if len(signer) != 20 {
			return nil, fmt.Errorf("OnChainPublicKey has wrong length for address")
		}
		signers = append(signers, common.BytesToAddress(signer))
	}

	var transmitters []common.Address
	for _, transmitter := range transmitterAccounts {
		if !common.IsHexAddress(string(transmitter)) {
			return nil, fmt.Errorf("TransmitAccount is not a valid Ethereum address")
		}
		transmitters = append(transmitters, common.HexToAddress(string(transmitter)))
	}

	configType := abi.MustNewType("tuple(uint32 paymentPremiumPPB,uint32 flatFeeMicroLink,uint32 checkGasLimit,uint24 stalenessSeconds,uint16 gasCeilingMultiplier,uint96 minUpkeepSpend,uint32 maxPerformGas,uint32 maxCheckDataSize,uint32 maxPerformDataSize,uint256 fallbackGasPrice,uint256 fallbackLinkPrice,address transcoder,address registrar)")
	onchainConfig, err := abi.Encode(map[string]interface{}{
		"paymentPremiumPPB":    d.cfg.PaymentPremiumPBB,
		"flatFeeMicroLink":     d.cfg.FlatFeeMicroLink,
		"checkGasLimit":        d.cfg.CheckGasLimit,
		"stalenessSeconds":     d.cfg.StalenessSeconds,
		"gasCeilingMultiplier": d.cfg.GasCeilingMultiplier,
		"minUpkeepSpend":       d.cfg.MinUpkeepSpend,
		"maxPerformGas":        d.cfg.MaxPerformGas,
		"maxCheckDataSize":     d.cfg.MaxCheckDataSize,
		"maxPerformDataSize":   d.cfg.MaxPerformDataSize,
		"fallbackGasPrice":     big.NewInt(d.cfg.FallbackGasPrice),
		"fallbackLinkPrice":    big.NewInt(d.cfg.FallbackLinkPrice),
		"transcoder":           common.HexToAddress(d.cfg.Transcoder),
		"registrar":            common.HexToAddress(d.cfg.Registrar),
	}, configType)
	if err != nil {
		return nil, err
	}

	return d.KeeperRegistryInterface.SetConfig(opts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}
