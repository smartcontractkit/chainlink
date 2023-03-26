package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum/mercury/exchanger"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum/mercury/verifier"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum/mercury/verifier_proxy"
)

type MercuryOCRConfig struct {
	Signers               []common.Address
	Transmitters          [][32]byte
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

type Exchanger interface {
	Address() string
	CommitTrade(commitment [32]byte) error
	ResolveTrade(encodedCommitment []byte) (string, error)
	ResolveTradeWithReport(chainlinkBlob []byte, encodedCommitment []byte) (*types.Receipt, error)
}

type EthereumExchanger struct {
	address   *common.Address
	client    blockchain.EVMClient
	exchanger *exchanger.Exchanger
}

func (v *EthereumExchanger) Address() string {
	return v.address.Hex()
}

func (e *EthereumExchanger) CommitTrade(commitment [32]byte) error {
	txOpts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := e.exchanger.CommitTrade(txOpts, commitment)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

func (e *EthereumExchanger) ResolveTrade(encodedCommitment []byte) (string, error) {
	callOpts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	}
	data, err := e.exchanger.ResolveTrade(callOpts, encodedCommitment)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (e *EthereumExchanger) ResolveTradeWithReport(chainlinkBlob []byte, encodedCommitment []byte) (*types.Receipt, error) {
	txOpts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	txOpts.GasLimit = 8000000
	tx, err := e.exchanger.ResolveTradeWithReport(txOpts, chainlinkBlob, encodedCommitment)
	if err != nil {
		// blockchain.LogRevertReason(err, exchanger.ExchangerABI)
		return nil, err
	}
	err = e.client.ProcessTransaction(tx)
	if err != nil {
		return nil, err
	}
	err = e.client.WaitForEvents()
	if err != nil {
		return nil, err
	}
	return e.client.GetTxReceipt(tx.Hash())
}

type VerifierProxy interface {
	Address() string
	InitializeVerifier(configDigest [32]byte, verifierAddress string) error
	Verify(signedReport []byte) error
}

type EthereumVerifierProxy struct {
	address       *common.Address
	client        blockchain.EVMClient
	verifierProxy *verifier_proxy.VerifierProxy
}

func (v *EthereumVerifierProxy) Address() string {
	return v.address.Hex()
}

func (v *EthereumVerifierProxy) InitializeVerifier(configDigest [32]byte, verifierAddr string) error {
	txOpts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.verifierProxy.InitializeVerifier(txOpts, configDigest, common.HexToAddress(verifierAddr))
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVerifierProxy) Verify(signedReport []byte) error {
	txOpts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.verifierProxy.Verify(txOpts, signedReport)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

type Verifier interface {
	Address() string
	SetConfig([32]byte, MercuryOCRConfig) error
	LatestConfigDetails(feedId [32]byte) (struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	}, error)
}

type EthereumVerifier struct {
	address  *common.Address
	client   blockchain.EVMClient
	verifier *verifier.Verifier
}

func (v *EthereumVerifier) Address() string {
	return v.address.Hex()
}

func (v *EthereumVerifier) SetConfig(feedId [32]byte, config MercuryOCRConfig) error {
	txOpts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().Msgf("Setting config, feedId: %s, config: %v", feedId, config)
	for i, s := range config.Signers {
		log.Info().Msgf("Signer %d: %x", i, s)
	}
	for i, s := range config.Transmitters {
		log.Info().Msgf("Transmitter %d: %x", i, s)
	}
	// log.Info().Msgf("Transmitters: %x", config.Transmitters)
	// log.Info().Msgf("OnchainConfig: %x", config.OnchainConfig)
	log.Info().Msgf("OffchainConfig: %x", config.OffchainConfig)

	tx, err := v.verifier.SetConfig(
		txOpts,
		feedId,
		config.Signers,
		config.Transmitters,
		config.F,
		config.OnchainConfig,
		config.OffchainConfigVersion,
		config.OffchainConfig,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVerifier) LatestConfigDetails(feedId [32]byte) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	}
	return v.verifier.LatestConfigDetails(opts, feedId)
}

func (e *EthereumContractDeployer) LoadVerifier(address common.Address) (Verifier, error) {
	instance, err := e.client.LoadContract("Verifier", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return verifier.NewVerifier(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVerifier{
		client:   e.client,
		address:  &address,
		verifier: instance.(*verifier.Verifier),
	}, err
}

func (e *EthereumContractDeployer) DeployVerifier(verifierProxyAddr string) (Verifier, error) {
	address, _, instance, err := e.client.DeployContract("Verifier", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return verifier.DeployVerifier(auth, backend, common.HexToAddress(verifierProxyAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVerifier{
		client:   e.client,
		address:  address,
		verifier: instance.(*verifier.Verifier),
	}, err
}

func (e *EthereumContractDeployer) LoadVerifierProxy(address common.Address) (VerifierProxy, error) {
	instance, err := e.client.LoadContract("VerifierProxy", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return verifier_proxy.NewVerifierProxy(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVerifierProxy{
		client:        e.client,
		address:       &address,
		verifierProxy: instance.(*verifier_proxy.VerifierProxy),
	}, err
}

func (e *EthereumContractDeployer) DeployVerifierProxy(accessControllerAddr string) (VerifierProxy, error) {
	address, _, instance, err := e.client.DeployContract("VerifierProxy", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return verifier_proxy.DeployVerifierProxy(auth, backend, common.HexToAddress(accessControllerAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVerifierProxy{
		client:        e.client,
		address:       address,
		verifierProxy: instance.(*verifier_proxy.VerifierProxy),
	}, err
}

func (e *EthereumContractDeployer) LoadExchanger(address common.Address) (Exchanger, error) {
	instance, err := e.client.LoadContract("Exchanger", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return exchanger.NewExchanger(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumExchanger{
		client:    e.client,
		address:   &address,
		exchanger: instance.(*exchanger.Exchanger),
	}, err
}

func (e *EthereumContractDeployer) DeployExchanger(verifierProxyAddr string, lookupURL string, maxDelay uint8) (Exchanger, error) {
	address, _, instance, err := e.client.DeployContract("Exchanger", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return exchanger.DeployExchanger(auth, backend,
			common.HexToAddress(verifierProxyAddr), lookupURL, maxDelay)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumExchanger{
		client:    e.client,
		address:   address,
		exchanger: instance.(*exchanger.Exchanger),
	}, err
}
