package integration_tests

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

const (
	CapabilityTypeTrigger   = 0
	CapabilityTypeAction    = 1
	CapabilityTypeConsensus = 2
	CapabilityTypeTarget    = 3
)

type peer struct {
	PeerID string
	Signer string
}

func peerIDToBytes(peerID string) ([32]byte, error) {
	var peerIDB ragetypes.PeerID
	err := peerIDB.UnmarshalText([]byte(peerID))
	if err != nil {
		return [32]byte{}, err
	}

	return peerIDB, nil
}

func peers(ps []peer) ([][32]byte, error) {
	out := [][32]byte{}
	for _, p := range ps {
		b, err := peerIDToBytes(p.PeerID)
		if err != nil {
			return nil, err
		}

		out = append(out, b)
	}

	return out, nil
}

func peerToNode(nopID uint32, p peer) (kcr.CapabilitiesRegistryNodeParams, error) {
	peerIDB, err := peerIDToBytes(p.PeerID)
	if err != nil {
		return kcr.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert peerID: %w", err)
	}

	sig := strings.TrimPrefix(p.Signer, "0x")
	signerB, err := hex.DecodeString(sig)
	if err != nil {
		return kcr.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert signer: %w", err)
	}

	var sigb [32]byte
	copy(sigb[:], signerB)

	return kcr.CapabilitiesRegistryNodeParams{
		NodeOperatorId: nopID,
		P2pId:          peerIDB,
		Signer:         sigb,
	}, nil
}

func setupCapabilitiesRegistryContract(ctx context.Context, t *testing.T, workflowDon donInfo, triggerDon donInfo,
	targetDon donInfo,
	transactOpts *bind.TransactOpts, backend *ethBackend) common.Address {
	addr, _, reg, err := kcr.DeployCapabilitiesRegistry(transactOpts, backend)
	require.NoError(t, err)

	backend.Commit()

	streamsTrigger := kcr.CapabilitiesRegistryCapability{
		LabelledName:   "streams-trigger",
		Version:        "1.0.0",
		CapabilityType: CapabilityTypeTrigger,
	}
	sid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, streamsTrigger.LabelledName, streamsTrigger.Version)
	require.NoError(t, err)

	writeChain := kcr.CapabilitiesRegistryCapability{
		LabelledName: "write_geth-testnet",
		Version:      "1.0.0",

		CapabilityType: CapabilityTypeTarget,
	}
	wid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, writeChain.LabelledName, writeChain.Version)
	if err != nil {
		log.Printf("failed to call GetHashedCapabilityId: %s", err)
	}

	ocr := kcr.CapabilitiesRegistryCapability{
		LabelledName:   "offchain_reporting",
		Version:        "1.0.0",
		CapabilityType: CapabilityTypeConsensus,
	}
	ocrid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, ocr.LabelledName, ocr.Version)
	require.NoError(t, err)

	_, err = reg.AddCapabilities(transactOpts, []kcr.CapabilitiesRegistryCapability{
		streamsTrigger,
		writeChain,
		ocr,
	})
	require.NoError(t, err)
	backend.Commit()

	_, err = reg.AddNodeOperators(transactOpts, []kcr.CapabilitiesRegistryNodeOperator{
		{
			Admin: transactOpts.From,
			Name:  "TEST_NODE_OPERATOR",
		},
	})
	require.NoError(t, err)
	blockHash := backend.Commit()

	logs, err := backend.FilterLogs(ctx, ethereum.FilterQuery{
		BlockHash: &blockHash,
		FromBlock: nil,
		ToBlock:   nil,
		Addresses: nil,
		Topics:    nil,
	})

	require.NoError(t, err)

	recLog, err := reg.ParseNodeOperatorAdded(logs[0])
	require.NoError(t, err)

	nopID := recLog.NodeOperatorId
	nodes := []kcr.CapabilitiesRegistryNodeParams{}
	for _, wfPeer := range workflowDon.peerIDs {
		n, innerErr := peerToNode(nopID, wfPeer)
		require.NoError(t, innerErr)

		n.HashedCapabilityIds = [][32]byte{ocrid}
		nodes = append(nodes, n)
	}

	for _, triggerPeer := range triggerDon.peerIDs {
		n, innerErr := peerToNode(nopID, triggerPeer)
		require.NoError(t, innerErr)

		n.HashedCapabilityIds = [][32]byte{sid}
		nodes = append(nodes, n)
	}

	for _, targetPeer := range targetDon.peerIDs {
		n, innerErr := peerToNode(nopID, targetPeer)
		require.NoError(t, innerErr)

		n.HashedCapabilityIds = [][32]byte{wid}
		nodes = append(nodes, n)
	}

	_, err = reg.AddNodes(transactOpts, nodes)
	require.NoError(t, err)

	// workflow DON
	ps, err := peers(workflowDon.peerIDs)
	require.NoError(t, err)

	cc := newCapabilityConfig()
	ccb, err := proto.Marshal(cc)
	require.NoError(t, err)

	cfgs := []kcr.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: ocrid,
			Config:       ccb,
		},
	}

	_, err = reg.AddDON(transactOpts, ps, cfgs, false, true, workflowDon.F)
	require.NoError(t, err)

	// trigger DON
	ps, err = peers(triggerDon.peerIDs)
	require.NoError(t, err)

	triggerCapabilityConfig := newCapabilityConfig()
	triggerCapabilityConfig.RemoteConfig = &pb.CapabilityConfig_RemoteTriggerConfig{
		RemoteTriggerConfig: &pb.RemoteTriggerConfig{
			RegistrationRefresh: durationpb.New(1000 * time.Millisecond),
			RegistrationExpiry:  durationpb.New(60000 * time.Millisecond),
			// F + 1
			MinResponsesToAggregate: uint32(triggerDon.F) + 1,
		},
	}

	configb, err := proto.Marshal(triggerCapabilityConfig)
	require.NoError(t, err)

	cfgs = []kcr.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: sid,
			Config:       configb,
		},
	}

	_, err = reg.AddDON(transactOpts, ps, cfgs, true, false, triggerDon.F)
	require.NoError(t, err)

	// target DON
	ps, err = peers(targetDon.peerIDs)
	require.NoError(t, err)

	targetCapabilityConfig := newCapabilityConfig()

	configWithLimit, err := values.WrapMap(map[string]any{"gasLimit": 500000})
	require.NoError(t, err)

	targetCapabilityConfig.DefaultConfig = values.Proto(configWithLimit).GetMapValue()

	targetCapabilityConfig.RemoteConfig = &pb.CapabilityConfig_RemoteTargetConfig{
		RemoteTargetConfig: &pb.RemoteTargetConfig{
			RequestHashExcludedAttributes: []string{"signed_report.Signatures"},
		},
	}

	remoteTargetConfigBytes, err := proto.Marshal(targetCapabilityConfig)
	require.NoError(t, err)

	cfgs = []kcr.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: wid,
			Config:       remoteTargetConfigBytes,
		},
	}

	_, err = reg.AddDON(transactOpts, ps, cfgs, true, false, targetDon.F)
	require.NoError(t, err)

	backend.Commit()

	return addr
}

func newCapabilityConfig() *pb.CapabilityConfig {
	return &pb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
	}
}

func setupForwarderContract(t *testing.T, workflowDon donInfo,
	transactOpts *bind.TransactOpts, backend *ethBackend) (common.Address, *forwarder.KeystoneForwarder) {
	addr, _, fwd, err := forwarder.DeployKeystoneForwarder(transactOpts, backend)
	require.NoError(t, err)
	backend.Commit()

	var signers []common.Address
	for _, p := range workflowDon.peerIDs {
		signers = append(signers, common.HexToAddress(p.Signer))
	}

	_, err = fwd.SetConfig(transactOpts, workflowDon.ID, workflowDon.ConfigVersion, workflowDon.F, signers)
	require.NoError(t, err)
	backend.Commit()

	return addr, fwd
}

func setupConsumerContract(t *testing.T, transactOpts *bind.TransactOpts, backend *ethBackend,
	forwarderAddress common.Address, workflowOwner string, workflowName string) (common.Address, *feeds_consumer.KeystoneFeedsConsumer) {
	addr, _, consumer, err := feeds_consumer.DeployKeystoneFeedsConsumer(transactOpts, backend)
	require.NoError(t, err)
	backend.Commit()

	var nameBytes [10]byte
	copy(nameBytes[:], workflowName)

	ownerAddr := common.HexToAddress(workflowOwner)

	_, err = consumer.SetConfig(transactOpts, []common.Address{forwarderAddress}, []common.Address{ownerAddr}, [][10]byte{nameBytes})
	require.NoError(t, err)

	backend.Commit()

	return addr, consumer
}

type ethBackend struct {
	services.StateMachine
	*backends.SimulatedBackend

	blockTimeProcessingTime time.Duration

	stopCh services.StopChan
	wg     sync.WaitGroup
}

func setupBlockchain(t *testing.T, initialEth int, blockTimeProcessingTime time.Duration) (*ethBackend, *bind.TransactOpts) {
	transactOpts := testutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := core.GenesisAlloc{transactOpts.From: {Balance: assets.Ether(initialEth).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	gethlog.SetDefault(gethlog.NewLogger(gethlog.NewTerminalHandlerWithLevel(os.Stderr, gethlog.LevelWarn, true)))
	backend.Commit()

	return &ethBackend{SimulatedBackend: backend, stopCh: make(services.StopChan),
		blockTimeProcessingTime: blockTimeProcessingTime}, transactOpts
}

func (b *ethBackend) Start(ctx context.Context) error {
	return b.StartOnce("ethBackend", func() error {
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			ticker := time.NewTicker(b.blockTimeProcessingTime)
			defer ticker.Stop()

			for {
				select {
				case <-b.stopCh:
					return
				case <-ctx.Done():
					return
				case <-ticker.C:
					b.SimulatedBackend.Commit()
				}
			}
		}()

		return nil
	})
}

func (b *ethBackend) Close() error {
	return b.StopOnce("ethBackend", func() error {
		close(b.stopCh)
		b.wg.Wait()
		return nil
	})
}
