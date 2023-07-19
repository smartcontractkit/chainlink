package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jpillora/backoff"
	"github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"gopkg.in/guregu/null.v2"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	FailedRPCContractCalls = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ocr2_failed_rpc_contract_calls",
		Help: "Running count of failed RPC contract calls to OCR2 configuration contract",
	},
		[]string{"chainID", "contractAddress", "feedID"},
	)
)

func NewRPCCallBackoff() backoff.Backoff {
	return backoff.Backoff{
		Factor: 2,
		Jitter: true,
		Min:    100 * time.Millisecond,
		Max:    1 * time.Hour,
	}
}

type ContractCaller interface {
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	ConfiguredChainID() *big.Int
}

type RelayConfig struct {
	ChainID                *utils.Big  `json:"chainID"`
	FromBlock              uint64      `json:"fromBlock"`
	EffectiveTransmitterID null.String `json:"effectiveTransmitterID"`

	// Contract-specific
	SendingKeys pq.StringArray `json:"sendingKeys"`

	// Mercury-specific
	FeedID *common.Hash `json:"feedID"`
}

type RelayOpts struct {
	// TODO BCF-2508 -- should anyone ever get the raw config bytes that are embedded in args? if not,
	// make this private and wrap the arg fields with funcs on RelayOpts
	relaytypes.RelayArgs
	c *RelayConfig
}

var ErrBadRelayConfig = errors.New("bad relay config")

func NewRelayOpts(args types.RelayArgs) *RelayOpts {
	return &RelayOpts{
		RelayArgs: args,
		c:         nil, // lazy initialization
	}
}

func (o *RelayOpts) RelayConfig() (RelayConfig, error) {
	var empty RelayConfig
	//TODO this should be done once and the error should be cached
	if o.c == nil {
		var c RelayConfig
		err := json.Unmarshal(o.RelayArgs.RelayConfig, &c)
		if err != nil {
			return empty, fmt.Errorf("%w: failed to deserialize relay config: %w", ErrBadRelayConfig, err)
		}
		o.c = &c
	}
	return *o.c, nil
}

type ConfigPoller interface {
	ocrtypes.ContractConfigTracker

	Start()
	Close() error
	Replay(ctx context.Context, fromBlock int64) error
}

// TODO(FUN-668): Migrate this fully into relaytypes.FunctionsProvider
type FunctionsProvider interface {
	relaytypes.FunctionsProvider
	LogPollerWrapper() LogPollerWrapper
}

type OracleRequest struct {
	RequestId           [32]byte
	RequestingContract  common.Address
	RequestInitiator    common.Address
	SubscriptionId      uint64
	SubscriptionOwner   common.Address
	Data                []byte
	DataVersion         uint16
	Flags               [32]byte
	CallbackGasLimit    uint64
	TxHash              common.Hash
	CoordinatorContract common.Address
	OnchainMetadata     []byte
}

type OracleResponse struct {
	RequestId [32]byte
}

type RouteUpdateSubscriber interface {
	UpdateRoutes(activeCoordinator common.Address, proposedCoordinator common.Address) error
}

// A LogPoller wrapper that understands router proxy contracts
//
//go:generate mockery --quiet --name LogPollerWrapper --output ./mocks/ --case=underscore
type LogPollerWrapper interface {
	relaytypes.Service
	LatestEvents() ([]OracleRequest, []OracleResponse, error)

	// TODO (FUN-668): Remove from the LOOP interface and only use internally within the EVM relayer
	SubscribeToUpdates(name string, subscriber RouteUpdateSubscriber)
}
