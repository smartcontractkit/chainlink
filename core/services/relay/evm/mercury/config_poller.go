package mercury

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
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

func ConfigFromLog(logData []byte) (FullConfigFromLog, error) {
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
}

func FilterName(addr common.Address, feedID common.Hash) string {
	return logpoller.FilterName("OCR3 Mercury ConfigPoller", addr.String(), feedID.Hex())
}

// NewConfigPoller creates a new Mercury ConfigPoller
func NewConfigPoller(ctx context.Context, lggr logger.Logger, destChainPoller logpoller.LogPoller, addr common.Address, feedId common.Hash) (*ConfigPoller, error) {
	err := destChainPoller.RegisterFilter(ctx, logpoller.Filter{Name: FilterName(addr, feedId), EventSigs: []common.Hash{FeedScopedConfigSet}, Addresses: []common.Address{addr}})
	if err != nil {
		return nil, err
	}

	cp := &ConfigPoller{
		lggr:               lggr,
		destChainLogPoller: destChainPoller,
		addr:               addr,
		feedId:             feedId,
	}

	return cp, nil
}

func (cp *ConfigPoller) Start() {}

func (cp *ConfigPoller) Close() error {
	return nil
}

func (cp *ConfigPoller) Notify() <-chan struct{} {
	return nil // rely on libocr's builtin config polling
}

// Replay abstracts the logpoller.LogPoller Replay() implementation
func (cp *ConfigPoller) Replay(ctx context.Context, fromBlock int64) error {
	return cp.destChainLogPoller.Replay(ctx, fromBlock)
}

// LatestConfigDetails returns the latest config details from the logs
func (cp *ConfigPoller) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	cp.lggr.Debugw("LatestConfigDetails", "eventSig", FeedScopedConfigSet, "addr", cp.addr, "topicIndex", feedIdTopicIndex, "feedID", cp.feedId)
	logs, err := cp.destChainLogPoller.IndexedLogs(ctx, FeedScopedConfigSet, cp.addr, feedIdTopicIndex, []common.Hash{cp.feedId}, 1)
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}
	if len(logs) == 0 {
		return 0, ocrtypes.ConfigDigest{}, nil
	}
	latest := logs[len(logs)-1]
	latestConfigSet, err := ConfigFromLog(latest.Data)
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}
	return uint64(latest.BlockNumber), latestConfigSet.ConfigDigest, nil
}

// LatestConfig returns the latest config from the logs on a certain block
func (cp *ConfigPoller) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	lgs, err := cp.destChainLogPoller.IndexedLogsByBlockRange(ctx, int64(changedInBlock), int64(changedInBlock), FeedScopedConfigSet, cp.addr, feedIdTopicIndex, []common.Hash{cp.feedId})
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	if len(lgs) == 0 {
		return ocrtypes.ContractConfig{}, nil
	}
	latestConfigSet, err := ConfigFromLog(lgs[len(lgs)-1].Data)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	cp.lggr.Infow("LatestConfig", "latestConfig", latestConfigSet)
	return latestConfigSet.ContractConfig, nil
}

// LatestBlockHeight returns the latest block height from the logs
func (cp *ConfigPoller) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	latest, err := cp.destChainLogPoller.LatestBlock(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return uint64(latest.BlockNumber), nil
}
