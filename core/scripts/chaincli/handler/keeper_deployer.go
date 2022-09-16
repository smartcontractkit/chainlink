package handler

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ocr2config "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/cmd"
	registry11 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry12 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	registry20 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"
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
}

func (d *v20KeeperDeployer) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return d.KeeperRegistryInterface.RegisterUpkeep(opts, target, gasLimit, admin, false, checkData)
}

func (d *v20KeeperDeployer) SetKeepers(opts *bind.TransactOpts, cls []cmd.HTTPClient, keepers []common.Address, _ []common.Address) (*types.Transaction, error) {
	oracleIdentities := make([]ocr2config.OracleIdentityExtra, len(cls))
	var wg sync.WaitGroup
	for i, cl := range cls {
		wg.Add(1)
		go func(i int, cl cmd.HTTPClient) {
			defer wg.Done()

			ocr2Config, err := getNodeOCR2Config(cl)
			if err != nil {
				panic(err)
				return
			}

			p2pKeyID, err := getP2PKeyID(cl)
			if err != nil {
				panic(err)
				return
			}

			offchainPkBytes, err := hex.DecodeString(ocr2Config.OffChainPublicKey)
			if err != nil {
				panic(err)
				return
			}

			offchainPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n := copy(offchainPkBytesFixed[:], offchainPkBytes)
			if n != ed25519.PublicKeySize {
				panic(fmt.Errorf("wrong num elements copied"))
				return
			}

			configPkBytes, err := hex.DecodeString(ocr2Config.ConfigPublicKey)
			if err != nil {
				panic(err)
				return
			}

			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n = copy(configPkBytesFixed[:], configPkBytes)
			if n != ed25519.PublicKeySize {
				panic(fmt.Errorf("wrong num elements copied"))
				return
			}

			oracleIdentities[i] = ocr2config.OracleIdentityExtra{
				OracleIdentity: ocr2config.OracleIdentity{
					OnchainPublicKey:  common.HexToAddress(ocr2Config.OnchainPublicKey).Bytes()[:],
					OffchainPublicKey: offchainPkBytesFixed,
					PeerID:            p2pKeyID,
					TransmitAccount:   ocr2types.Account(keepers[i].String()),
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}
		}(i, cl)
	}
	wg.Wait()

	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr2config.ContractSetConfigArgsForEthereumIntegrationTest(oracleIdentities, 1, uint64(1000))
	if err != nil {
		return nil, err
	}

	return d.KeeperRegistryInterface.SetConfig(opts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}
