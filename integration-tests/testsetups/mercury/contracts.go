package mercury

import (
	"math/big"

	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type Order struct {
	FeedID       [32]byte
	CurrencySrc  [32]byte
	CurrencyDst  [32]byte
	AmountSrc    *big.Int
	MinAmountDst *big.Int
	Sender       common.Address
	Receiver     common.Address
}

func CreateEncodedCommitment(order Order) ([]byte, error) {
	// bytes32 feedID, bytes32 currencySrc, bytes32 currencyDst, uint256 amountSrc, uint256 minAmountDst, address sender, address receiver
	orderType, _ := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "feedID", Type: "bytes32"},
		{Name: "currencySrc", Type: "bytes32"},
		{Name: "currencyDst", Type: "bytes32"},
		{Name: "amountSrc", Type: "uint256"},
		{Name: "minAmountDst", Type: "uint256"},
		{Name: "sender", Type: "address"},
		{Name: "receiver", Type: "address"},
	})
	var args abi.Arguments = []abi.Argument{{Type: orderType}}
	return args.Pack(order)
}

func CreateCommitmentHash(order Order) common.Hash {
	uint256Ty, _ := abi.NewType("uint256", "", nil)
	bytes32Ty, _ := abi.NewType("bytes32", "", nil)
	addressTy, _ := abi.NewType("address", "", nil)

	arguments := abi.Arguments{
		{
			Type: bytes32Ty,
		},
		{
			Type: bytes32Ty,
		},
		{
			Type: bytes32Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: addressTy,
		},
		{
			Type: addressTy,
		},
	}

	bytes, _ := arguments.Pack(
		order.FeedID,
		order.CurrencySrc,
		order.CurrencyDst,
		order.AmountSrc,
		order.MinAmountDst,
		order.Sender,
		order.Receiver,
	)

	return crypto.Keccak256Hash(bytes)
}

func LoadMercuryContracts(evmClient blockchain.EVMClient,
	verifierAddr string, verifierProxyAddr string, exchangerAddr string) (
	contracts.Verifier, contracts.VerifierProxy, contracts.Exchanger, error) {

	contractDeployer, err := contracts.NewContractDeployer(evmClient)
	if err != nil {
		return nil, nil, nil, err
	}
	verifier, err := contractDeployer.LoadVerifier(common.HexToAddress(verifierAddr))
	if err != nil {
		return nil, nil, nil, err
	}
	verifierProxy, err := contractDeployer.LoadVerifierProxy(common.HexToAddress(verifierProxyAddr))
	if err != nil {
		return verifier, nil, nil, err
	}
	exchanger, err := contractDeployer.LoadExchanger(common.HexToAddress(exchangerAddr))
	if err != nil {
		return verifier, verifierProxy, nil, err
	}

	return verifier, verifierProxy, exchanger, nil
}

func DeployMercuryContracts(evmClient blockchain.EVMClient, lookupUrl string, ocrConfig contracts.MercuryOCRConfig) (
	contracts.Verifier, contracts.VerifierProxy, contracts.Exchanger, contracts.ReadAccessController, error) {
	contractDeployer, err := contracts.NewContractDeployer(evmClient)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	accessController, err := contractDeployer.DeployReadAccessController()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// verifierProxy, err := contractDeployer.DeployVerifierProxy(accessController.Address())
	// Use zero address for access controller disables access control
	verifierProxy, err := contractDeployer.DeployVerifierProxy("0x0")
	if err != nil {
		return nil, nil, nil, accessController, err
	}

	verifier, err := contractDeployer.DeployVerifier(verifierProxy.Address())
	if err != nil {
		return nil, verifierProxy, nil, accessController, err
	}

	exchanger, err := contractDeployer.DeployExchanger(verifierProxy.Address(), lookupUrl, 255)
	if err != nil {
		return verifier, verifierProxy, nil, accessController, err
	}

	err = accessController.AddAccess(exchanger.Address())
	if err != nil {
		return verifier, verifierProxy, exchanger, accessController, err
	}

	return verifier, verifierProxy, exchanger, accessController, nil
}
