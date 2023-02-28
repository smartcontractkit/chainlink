package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum/mercury/exchanger"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum/mercury/verifier"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum/mercury/verifier_proxy"
)

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
	txOpts.GasLimit = 1000000
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
	SetConfig(OCRConfig) error
	LatestConfigDetails() (struct {
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

func (v *EthereumVerifier) SetConfig(ocrConfig OCRConfig) error {
	txOpts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.verifier.SetConfig(
		txOpts,
		ocrConfig.Signers,
		ocrConfig.Transmitters,
		ocrConfig.F,
		ocrConfig.OnchainConfig,
		ocrConfig.OffchainConfigVersion,
		ocrConfig.OffchainConfig,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVerifier) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	}
	return v.verifier.LatestConfigDetails(opts)
}

func (e *EthereumContractDeployer) DeployVerifier(feedId [32]byte, verifierProxyAddr string) (Verifier, error) {
	address, _, instance, err := e.client.DeployContract("Verifier", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return verifier.DeployVerifier(auth, backend, feedId, common.HexToAddress(verifierProxyAddr))
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
