package deployment

import (
	"context"
	"errors"
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	ctfTestEnv "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/gauntlet"
	"github.com/smartcontractkit/chainlink/integration-tests/solclient"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	types2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	types3 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	csav1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/csa/v1"

	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	nodev1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/node/v1"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type OnchainClient interface {
	// For EVM specifically we can use existing geth interface
	// to abstract chain clients.
	bind.ContractBackend
}

type OffchainClient interface {
	// The job distributor grpc interface can be used to abstract offchain read/writes
	jobv1.JobServiceClient
	nodev1.NodeServiceClient
	csav1.CSAServiceClient
	// sometimes we need to connect to the nodes directly to execute actions that are not supported by JD/FMS
	NodeClients() []*client.ChainlinkK8sClient
}

type SolanaChain struct {
	ChainId  uint64
	Deployer *gauntlet.SolanaGauntlet
	Client   *solclient.Client
}

// TODO: we should rename it EVM chain, as it's not generic at all
type Chain struct {
	// Selectors used as canonical chain identifier.
	Selector uint64
	Client   OnchainClient
	// Note the Sign function can be abstract supporting a variety of key storage mechanisms (e.g. KMS etc).
	// TODO remove it and use only DeployerKeys
	DeployerKey *bind.TransactOpts
	// Rename to something more universal, as these keys are not only for deploying, but also for other actions
	DeployerKeys []*bind.TransactOpts
	// we need that data to set up chainlink nodes (we need the RPC URLs)
	// TODO wonder if I could get rid of this field, I could probably return it next to Chain from specific implementations
	// And pass that data to the nodes instead of storing it here
	EVMNetworkWithRPCs RpcProvider
	// Function to execute if transaction submission fails.
	RetrySubmit func(tx *types.Transaction, err error) (*types.Transaction, error)
	Confirm     func(tx common.Hash) (uint64, error)
}

// a bit unfortunate, but we need to be able to pass "private" urls of chains to chainlink nodes
// when running in Docker; that's unless we create a DNS/proxy that would allow to access private
// chains using the same URL both from inside docker containers and from the machine where test code executes
type RpcProvider interface {
	EVMNetwork() blockchain.EVMNetwork
	PrivateHttpUrls() []string
	PrivateWsUrls() []string
	PublicHttpUrls() []string
	PublicWsUrls() []string
}

func NewEVMNetworkWithRPCs(evmNetwork blockchain.EVMNetwork, rpcProvider ctfTestEnv.RpcProvider) RpcProvider {
	return &EVMNetworkWithRPCs{
		evmNetwork,
		rpcProvider,
	}
}

type EVMNetworkWithRPCs struct {
	evmNetwork blockchain.EVMNetwork
	ctfTestEnv.RpcProvider
}

func (s *EVMNetworkWithRPCs) EVMNetwork() blockchain.EVMNetwork {
	return s.evmNetwork
}

func (s *EVMNetworkWithRPCs) PrivateHttpUrls() []string {
	return s.RpcProvider.PrivateHttpUrls()
}

func (s *EVMNetworkWithRPCs) PrivateWsUrls() []string {
	return s.RpcProvider.PrivateWsUrsl()
}

func (s *EVMNetworkWithRPCs) PublicHttpUrls() []string {
	return s.PublicHttpUrls()
}

func (s *EVMNetworkWithRPCs) PublicWsUrls() []string {
	return s.RpcProvider.PublicWsUrls()
}

// NoOpRetrySubmit is a retry submit function that does nothing.
func NoOpRetrySubmit(_ *types.Transaction, err error) (*types.Transaction, error) {
	return nil, err
}

// we do need mocks for higher-level environments
type Mocks struct {
	// we use Mockserver in k8s
	MockServer *ctfClient.MockserverClient
	// we use Killgrave in Docker
	KillGrave *ctfTestEnv.Killgrave
}

type Environment[GenericChain any] struct {
	Name     string
	Chains   map[uint64]GenericChain
	Offchain OffchainClient
	NodeIDs  []string
	Mocks    Mocks
	Logger   logger.Logger
}

func (e Environment[GenericChain]) AllChainSelectors() []uint64 {
	var selectors []uint64
	for sel := range e.Chains {
		selectors = append(selectors, sel)
	}
	return selectors
}

func ConfirmIfNoError(chain Chain, tx *types.Transaction, err error) (uint64, error) {
	if err != nil {
		//revive:disable
		var d rpc.DataError
		ok := errors.As(err, &d)
		if ok {
			return 0, fmt.Errorf("got Data Error: %s", d.ErrorData())
		}
		return 0, err
	}
	return chain.Confirm(tx.Hash())
}

func MaybeDataErr(err error) error {
	//revive:disable
	var d rpc.DataError
	ok := errors.As(err, &d)
	if ok {
		return d
	}
	return err
}

func UBigInt(i uint64) *big.Int {
	return new(big.Int).SetUint64(i)
}

func E18Mult(amount uint64) *big.Int {
	return new(big.Int).Mul(UBigInt(amount), UBigInt(1e18))
}

type OCRConfig struct {
	OffchainPublicKey types2.OffchainPublicKey
	// For EVM-chains, this an *address*.
	OnchainPublicKey          types2.OnchainPublicKey
	PeerID                    p2pkey.PeerID
	TransmitAccount           types2.Account
	ConfigEncryptionPublicKey types3.ConfigEncryptionPublicKey
	IsBootstrap               bool
	MultiAddr                 string // TODO: type
}

type Nodes []Node

func (n Nodes) PeerIDs(chainSel uint64) [][32]byte {
	var peerIDs [][32]byte
	for _, node := range n {
		cfg := node.SelToOCRConfig[chainSel]
		// NOTE: Assume same peerID for all chains.
		// Might make sense to change proto as peerID is 1-1 with node?
		peerIDs = append(peerIDs, cfg.PeerID)
	}
	return peerIDs
}

func (n Nodes) BootstrapPeerIDs(chainSel uint64) [][32]byte {
	var peerIDs [][32]byte
	for _, node := range n {
		cfg := node.SelToOCRConfig[chainSel]
		if !cfg.IsBootstrap {
			continue
		}
		peerIDs = append(peerIDs, cfg.PeerID)
	}
	return peerIDs
}

// OffchainPublicKey types.OffchainPublicKey
// // For EVM-chains, this an *address*.
// OnchainPublicKey types.OnchainPublicKey
// PeerID           string
// TransmitAccount  types.Account
type Node struct {
	SelToOCRConfig map[uint64]OCRConfig
}

func MustPeerIDFromString(s string) p2pkey.PeerID {
	p := p2pkey.PeerID{}
	if err := p.UnmarshalString(s); err != nil {
		panic(err)
	}
	return p
}

// Gathers all the node info through JD required to be able to set
// OCR config for example.
func NodeInfo(nodeIDs []string, oc OffchainClient) (Nodes, error) {
	var nodes []Node
	for _, node := range nodeIDs {
		// TODO: Filter should accept multiple nodes
		nodeChainConfigs, err := oc.ListNodeChainConfigs(context.Background(), &nodev1.ListNodeChainConfigsRequest{Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
			NodeIds: []string{node},
		}})
		if err != nil {
			return nil, err
		}
		selToOCRConfig := make(map[uint64]OCRConfig)
		for _, chainConfig := range nodeChainConfigs.ChainConfigs {
			if chainConfig.Chain.Type == nodev1.ChainType_CHAIN_TYPE_SOLANA {
				// Note supported for CCIP yet.
				continue
			}
			evmChainID, err := strconv.Atoi(chainConfig.Chain.Id)
			if err != nil {
				return nil, err
			}
			sel, err := chain_selectors.SelectorFromChainId(uint64(evmChainID))
			if err != nil {
				return nil, err
			}
			b := common.Hex2Bytes(chainConfig.Ocr2Config.OcrKeyBundle.OffchainPublicKey)
			var opk types2.OffchainPublicKey
			copy(opk[:], b)

			b = common.Hex2Bytes(chainConfig.Ocr2Config.OcrKeyBundle.ConfigPublicKey)
			var cpk types3.ConfigEncryptionPublicKey
			copy(cpk[:], b)

			selToOCRConfig[sel] = OCRConfig{
				OffchainPublicKey:         opk,
				OnchainPublicKey:          common.HexToAddress(chainConfig.Ocr2Config.OcrKeyBundle.OnchainSigningAddress).Bytes(),
				PeerID:                    MustPeerIDFromString(chainConfig.Ocr2Config.P2PKeyBundle.PeerId),
				TransmitAccount:           types2.Account(chainConfig.AccountAddress),
				ConfigEncryptionPublicKey: cpk,
				IsBootstrap:               chainConfig.Ocr2Config.IsBootstrap,
				MultiAddr:                 chainConfig.Ocr2Config.Multiaddr,
			}
		}
		nodes = append(nodes, Node{
			SelToOCRConfig: selToOCRConfig,
		})
	}

	return nodes, nil
}
