package deployment

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	types2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	types3 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/grpc"

	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type OnchainClient interface {
	// For EVM specifically we can use existing geth interface
	// to abstract chain clients.
	bind.ContractBackend
	bind.DeployBackend
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
}

type OffchainClient interface {
	// The job distributor grpc interface can be used to abstract offchain read/writes
	jobv1.JobServiceClient
	nodev1.NodeServiceClient
	csav1.CSAServiceClient
}

type Chain struct {
	// Selectors used as canonical chain identifier.
	Selector uint64
	Client   OnchainClient
	// Note the Sign function can be abstract supporting a variety of key storage mechanisms (e.g. KMS etc).
	DeployerKey *bind.TransactOpts
	Confirm     func(tx *types.Transaction) (uint64, error)
}

type Environment struct {
	Name     string
	Chains   map[uint64]Chain
	Offchain OffchainClient
	NodeIDs  []string
	Logger   logger.Logger
}

func (e Environment) AllChainSelectors() []uint64 {
	var selectors []uint64
	for sel := range e.Chains {
		selectors = append(selectors, sel)
	}
	sort.Slice(selectors, func(i, j int) bool {
		return selectors[i] < selectors[j]
	})
	return selectors
}

func (e Environment) AllChainSelectorsExcluding(excluding []uint64) []uint64 {
	var selectors []uint64
	for sel := range e.Chains {
		excluded := false
		for _, toExclude := range excluding {
			if sel == toExclude {
				excluded = true
			}
		}
		if excluded {
			continue
		}
		selectors = append(selectors, sel)
	}
	sort.Slice(selectors, func(i, j int) bool {
		return selectors[i] < selectors[j]
	})
	return selectors
}

func ConfirmIfNoError(chain Chain, tx *types.Transaction, err error) (uint64, error) {
	if err != nil {
		//revive:disable
		var d rpc.DataError
		ok := errors.As(err, &d)
		if ok {
			return 0, fmt.Errorf("transaction reverted: Error %s ErrorData %v", d.Error(), d.ErrorData())
		}
		return 0, err
	}
	return chain.Confirm(tx)
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
	KeyBundleID               string
}

// Nodes includes is a group CL nodes.
type Nodes []Node

// PeerIDs returns peerIDs in a sorted list
func (n Nodes) PeerIDs() [][32]byte {
	var peerIDs [][32]byte
	for _, node := range n {
		peerIDs = append(peerIDs, node.PeerID)
	}
	sort.Slice(peerIDs, func(i, j int) bool {
		return bytes.Compare(peerIDs[i][:], peerIDs[j][:]) < 0
	})
	return peerIDs
}

func (n Nodes) NonBootstraps() Nodes {
	var nonBootstraps Nodes
	for _, node := range n {
		if node.IsBootstrap {
			continue
		}
		nonBootstraps = append(nonBootstraps, node)
	}
	return nonBootstraps
}

func (n Nodes) DefaultF() uint8 {
	return uint8(len(n) / 3)
}

func (n Nodes) BootstrapLocators() []string {
	bootstrapMp := make(map[string]struct{})
	for _, node := range n {
		if node.IsBootstrap {
			bootstrapMp[fmt.Sprintf("%s@%s",
				// p2p_12D3... -> 12D3...
				node.PeerID.String()[4:], node.MultiAddr)] = struct{}{}
		}
	}
	var locators []string
	for b := range bootstrapMp {
		locators = append(locators, b)
	}
	return locators
}

type Node struct {
	NodeID         string
	SelToOCRConfig map[uint64]OCRConfig
	PeerID         p2pkey.PeerID
	IsBootstrap    bool
	MultiAddr      string
	AdminAddr      string
}

func (n Node) FirstOCRKeybundle() OCRConfig {
	for _, ocrConfig := range n.SelToOCRConfig {
		return ocrConfig
	}
	return OCRConfig{}
}

func MustPeerIDFromString(s string) p2pkey.PeerID {
	p := p2pkey.PeerID{}
	if err := p.UnmarshalString(s); err != nil {
		panic(err)
	}
	return p
}

type NodeChainConfigsLister interface {
	ListNodeChainConfigs(ctx context.Context, in *nodev1.ListNodeChainConfigsRequest, opts ...grpc.CallOption) (*nodev1.ListNodeChainConfigsResponse, error)
}

// Gathers all the node info through JD required to be able to set
// OCR config for example.
func NodeInfo(nodeIDs []string, oc NodeChainConfigsLister) (Nodes, error) {
	var nodes []Node
	for _, nodeID := range nodeIDs {
		// TODO: Filter should accept multiple nodes
		nodeChainConfigs, err := oc.ListNodeChainConfigs(context.Background(), &nodev1.ListNodeChainConfigsRequest{Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
			NodeIds: []string{nodeID},
		}})
		if err != nil {
			return nil, err
		}
		selToOCRConfig := make(map[uint64]OCRConfig)
		bootstrap := false
		var peerID p2pkey.PeerID
		var multiAddr string
		var adminAddr string
		for _, chainConfig := range nodeChainConfigs.ChainConfigs {
			if chainConfig.Chain.Type == nodev1.ChainType_CHAIN_TYPE_SOLANA {
				// Note supported for CCIP yet.
				continue
			}
			// NOTE: Assume same peerID/multiAddr for all chains.
			// Might make sense to change proto as peerID/multiAddr is 1-1 with nodeID?
			peerID = MustPeerIDFromString(chainConfig.Ocr2Config.P2PKeyBundle.PeerId)
			multiAddr = chainConfig.Ocr2Config.Multiaddr
			adminAddr = chainConfig.AdminAddress
			if chainConfig.Ocr2Config.IsBootstrap {
				// NOTE: Assume same peerID for all chains.
				// Might make sense to change proto as peerID is 1-1 with nodeID?
				bootstrap = true
				break
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
				KeyBundleID:               chainConfig.Ocr2Config.OcrKeyBundle.BundleId,
			}
		}
		nodes = append(nodes, Node{
			NodeID:         nodeID,
			SelToOCRConfig: selToOCRConfig,
			IsBootstrap:    bootstrap,
			PeerID:         peerID,
			MultiAddr:      multiAddr,
			AdminAddr:      adminAddr,
		})
	}

	return nodes, nil
}

type CapabilityRegistryConfig struct {
	EVMChainID uint64         // chain id of the chain the CR is deployed on
	Contract   common.Address // address of the CR contract
}
