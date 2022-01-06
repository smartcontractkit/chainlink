package blockhashstorefeeder

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	blockhash_store "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/log"
)

type Config interface {
	MinIncomingConfirmations() uint32
	EvmGasLimitDefault() uint64
	KeySpecificMaxGasPriceWei(addr common.Address) *big.Int
	MinRequiredOutgoingConfirmations() uint64
}

// Delegate is the Blockhash Store Feeder delegate.
type Delegate struct {
	chainSet    evm.ChainSet
	lggr        logger.Logger
	ethKeystore keystore.Eth
}

// NewDelegate creates a new Blockhash Store Feeder delegate.
func NewDelegate(chainSet evm.ChainSet, lggr logger.Logger, ks keystore.Master) *Delegate {
	return &Delegate{
		chainSet:    chainSet,
		lggr:        lggr.With("jobType", job.BlockhashStoreFeeder),
		ethKeystore: ks.Eth(),
	}
}

func (d *Delegate) JobType() job.Type {
	return job.BlockhashStoreFeeder
}

func (d *Delegate) AfterJobCreated(spec job.Job) {}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.Service, error) {
	if jb.BlockhashStoreFeederSpec == nil {
		return nil, fmt.Errorf("blockhashstorefeeder.Delegate expects a BlockhashStoreFeederSpec, got %+v", jb)
	}

	chain, err := d.chainSet.Get(jb.BlockhashStoreFeederSpec.EVMChainID.ToInt())
	if err != nil {
		return nil, err
	}

	blockhashStore, err := blockhash_store.NewBlockhashStore(
		jb.BlockhashStoreFeederSpec.BlockhashStoreAddress.Address(),
		chain.Client(),
	)
	if err != nil {
		return nil, err
	}

	// Get parse log function depending on which VRF version we're using
	parseLogFunc, err := getParseLogFunc(
		jb.BlockhashStoreFeederSpec.VRFVersion,
		jb.BlockhashStoreFeederSpec.CoordinatorAddress.Address(),
		chain.Client(),
	)
	if err != nil {
		return nil, err
	}

	// Ditto with log topics
	requestTopic, responseTopic, err := getLogTopics(jb.BlockhashStoreFeederSpec.VRFVersion)
	if err != nil {
		return nil, err
	}

	lggr := d.lggr.With(
		"jobID", jb.ID,
		"externalJobID", jb.ExternalJobID,
		"coordinatorAddress", jb.BlockhashStoreFeederSpec.CoordinatorAddress,
		"blockhashStoreAddress", jb.BlockhashStoreFeederSpec.BlockhashStoreAddress,
		"vrfVersion", jb.BlockhashStoreFeederSpec.VRFVersion,
	)

	abi, err := blockhash_store.BlockhashStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	return []job.Service{
		&feeder{
			config:            chain.Config(),
			lggr:              lggr,
			txm:               chain.TxManager(),
			job:               jb,
			blockhashStore:    blockhashStore,
			parseLogFunc:      parseLogFunc,
			wg:                &sync.WaitGroup{},
			ethClient:         chain.Client(),
			requestTopic:      requestTopic,
			responseTopic:     responseTopic,
			blockhashStoreABI: abi,
			ethKeystore:       d.ethKeystore,
		},
	}, nil
}

func getParseLogFunc(vrfVersion int32, coordinatorAddress common.Address, backend bind.ContractBackend) (log.ParseLogFunc, error) {
	if vrfVersion == 1 {
		coordinator, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(
			coordinatorAddress,
			backend,
		)
		if err != nil {
			return nil, err
		}
		return coordinator.ParseLog, nil
	} else if vrfVersion == 2 {
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(
			coordinatorAddress,
			backend,
		)
		if err != nil {
			return nil, err
		}
		return coordinator.ParseLog, nil
	}
	return nil, fmt.Errorf("invalid vrf version: %d. expected 1 or 2", vrfVersion)
}

func getLogTopics(vrfVersion int32) (requestTopic common.Hash, responseTopic common.Hash, err error) {
	if vrfVersion == 1 {
		return solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{}.Topic(),
			solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled{}.Topic(), nil
	} else if vrfVersion == 2 {
		return vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{}.Topic(),
			vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled{}.Topic(), nil
	}
	return common.Hash{}, common.Hash{}, fmt.Errorf("invalid vrf version: %d. expected 1 or 2", vrfVersion)
}
