package mercury

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
)

// FeedScopedConfigSet ConfigSet with FeedID for use with mercury (and multi-config DON)
var FeedScopedConfigSet common.Hash

var verifierABI abi.ABI

const (
	configSetEventName = "ConfigSet"
	feedIdTopicIndex   = 1
)

func init() {
	var err error
	verifierABI, err = abi.JSON(strings.NewReader(verifier.VerifierABI))
	if err != nil {
		panic(err)
	}
	FeedScopedConfigSet = verifierABI.Events[configSetEventName].ID
}

// FullConfigFromLog defines the contract config with the feedID
type FullConfigFromLog struct {
	ocrtypes.ContractConfig
	feedID utils.FeedID
}

func unpackLogData(d []byte) (*verifier.VerifierConfigSet, error) {
	unpacked := new(verifier.VerifierConfigSet)

	err := verifierABI.UnpackIntoInterface(unpacked, configSetEventName, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack log data")
	}

	return unpacked, nil
}

func configFromLog(logData []byte) (FullConfigFromLog, error) {
	unpacked, err := unpackLogData(logData)
	if err != nil {
		return FullConfigFromLog{}, err
	}

	var transmitAccounts []ocrtypes.Account
	for _, addr := range unpacked.OffchainTransmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(fmt.Sprintf("%x", addr)))
	}
	var signers []ocrtypes.OnchainPublicKey
	for _, addr := range unpacked.Signers {
		addr := addr
		signers = append(signers, addr[:])
	}

	return FullConfigFromLog{
		feedID: unpacked.FeedId,
		ContractConfig: ocrtypes.ContractConfig{
			ConfigDigest:          unpacked.ConfigDigest,
			ConfigCount:           unpacked.ConfigCount,
			Signers:               signers,
			Transmitters:          transmitAccounts,
			F:                     unpacked.F,
			OnchainConfig:         unpacked.OnchainConfig,
			OffchainConfigVersion: unpacked.OffchainConfigVersion,
			OffchainConfig:        unpacked.OffchainConfig,
		},
	}, nil
}

// ConfigPoller defines the Mercury Config Poller
type ConfigPoller struct {
	lggr               logger.Logger
	destChainLogPoller logpoller.LogPoller
	addr               common.Address
	feedId             common.Hash
	notifyCh           chan struct{}
	shutdown           chan struct{}
}

func FilterName(addr common.Address, feedID common.Hash) string {
	return logpoller.FilterName("OCR3 Mercury ConfigPoller", addr.String(), feedID.Hex())
}

// NewConfigPoller creates a new Mercury ConfigPoller
func NewConfigPoller(lggr logger.Logger, destChainPoller logpoller.LogPoller, addr common.Address, feedId common.Hash) (*ConfigPoller, error) {
	err := destChainPoller.RegisterFilter(logpoller.Filter{Name: FilterName(addr, feedId), EventSigs: []common.Hash{FeedScopedConfigSet}, Addresses: []common.Address{addr}})
	if err != nil {
		return nil, err
	}

	cp := &ConfigPoller{
		lggr:               lggr,
		destChainLogPoller: destChainPoller,
		addr:               addr,
		feedId:             feedId,
		notifyCh:           make(chan struct{}, 1),
		shutdown:           make(chan struct{}),
	}

	return cp, nil
}

// Start the subscription to Postgres' notify events.
func (cp *ConfigPoller) Start() {
	// todo: why doesn't start have a ctx?
	go cp.trackConfig(context.Background())
}

// Close the subscription to Postgres' notify events.
func (cp *ConfigPoller) Close() error {
	cp.shutdown <- struct{}{}
	return nil
}

// Notify abstracts the logpoller.LogPoller Notify() implementation
func (cp *ConfigPoller) Notify() <-chan struct{} {
	return cp.notifyCh
}

// Replay abstracts the logpoller.LogPoller Replay() implementation
func (cp *ConfigPoller) Replay(ctx context.Context, fromBlock int64) error {
	return cp.destChainLogPoller.Replay(ctx, fromBlock)
}

// LatestConfigDetails returns the latest config details from the logs
func (cp *ConfigPoller) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	cp.lggr.Debugw("LatestConfigDetails", "eventSig", FeedScopedConfigSet, "addr", cp.addr, "topicIndex", feedIdTopicIndex, "feedID", cp.feedId)
	logs, err := cp.destChainLogPoller.IndexedLogs(FeedScopedConfigSet, cp.addr, feedIdTopicIndex, []common.Hash{cp.feedId}, 1, pg.WithParentCtx(ctx))
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}
	if len(logs) == 0 {
		return 0, ocrtypes.ConfigDigest{}, nil
	}
	latest := logs[len(logs)-1]
	latestConfigSet, err := configFromLog(latest.Data)
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}
	return uint64(latest.BlockNumber), latestConfigSet.ConfigDigest, nil
}

// LatestConfig returns the latest config from the logs on a certain block
func (cp *ConfigPoller) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	lgs, err := cp.destChainLogPoller.IndexedLogsByBlockRange(int64(changedInBlock), int64(changedInBlock), FeedScopedConfigSet, cp.addr, feedIdTopicIndex, []common.Hash{cp.feedId}, pg.WithParentCtx(ctx))
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	if len(lgs) == 0 {
		return ocrtypes.ContractConfig{}, nil
	}
	latestConfigSet, err := configFromLog(lgs[len(lgs)-1].Data)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	cp.lggr.Infow("LatestConfig", "latestConfig", latestConfigSet)
	return latestConfigSet.ContractConfig, nil
}

// LatestBlockHeight returns the latest block height from the logs
func (cp *ConfigPoller) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	latest, err := cp.destChainLogPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return uint64(latest.BlockNumber), nil
}

func (cp *ConfigPoller) trackConfig(ctx context.Context) {

	lastConfig := ocrtypes.ConfigDigest{}

	// poll at half the logpoller frequency to avoid aliasing
	ticker := time.NewTicker(time.Duration(cp.destChainLogPoller.PollInterval().Milliseconds() / 2))

	for {
		select {
		case <-ctx.Done():
			return
		case <-cp.shutdown:
			return
		case <-ticker.C:
			// todo: under the cover this is calling LogPoller. logpoller must be able to satisfy tight polling loops like this for mercury and therefore may need some caching layer to
			// prevent hammering the database. this will be a concern when there are many mercury concurrent providers i.e. many concurrent mercury jobs on a single node
			block, config, err := cp.LatestConfigDetails(ctx)
			if err != nil {
				continue
			}
			if block > 0 && bytes.Compare(lastConfig[:], config[:]) != 0 {
				lastConfig = config
				cp.notifyCh <- struct{}{}
			}
		}
	}
}
