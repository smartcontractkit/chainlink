package mercury

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mercury_verifier"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
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
	verifierABI, err = abi.JSON(strings.NewReader(mercury_verifier.MercuryVerifierABI))
	if err != nil {
		panic(err)
	}
	FeedScopedConfigSet = verifierABI.Events[configSetEventName].ID
}

// FullConfigFromLog defines the contract config with the feedID
type FullConfigFromLog struct {
	ocrtypes.ContractConfig
	feedID [32]byte
}

func unpackLogData(d []byte) (*mercury_verifier.MercuryVerifierConfigSet, error) {
	unpacked := new(mercury_verifier.MercuryVerifierConfigSet)

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
}

func FilterName(addr common.Address) string {
	return logpoller.FilterName("OCR2 Mercury ConfigPoller", addr.String())
}

// NewConfigPoller creates a new Mercury ConfigPoller
func NewConfigPoller(lggr logger.Logger, destChainPoller logpoller.LogPoller, addr common.Address, feedId common.Hash) (*ConfigPoller, error) {
	err := destChainPoller.RegisterFilter(logpoller.Filter{Name: FilterName(addr), EventSigs: []common.Hash{FeedScopedConfigSet}, Addresses: []common.Address{addr}})
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

// Notify noop method
// TODO: implement this, see: https://smartcontract-it.atlassian.net/browse/MERC-302
func (lp *ConfigPoller) Notify() <-chan struct{} {
	return nil
}

// Replay abstracts the logpoller.LogPoller Replay() implementation
func (lp *ConfigPoller) Replay(ctx context.Context, fromBlock int64) error {
	return lp.destChainLogPoller.Replay(ctx, fromBlock)
}

// LatestConfigDetails returns the latest config details from the logs
func (lp *ConfigPoller) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	logs, err := lp.destChainLogPoller.IndexedLogs(FeedScopedConfigSet, lp.addr, feedIdTopicIndex, []common.Hash{lp.feedId}, 1, pg.WithParentCtx(ctx))
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
func (lp *ConfigPoller) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	lgs, err := lp.destChainLogPoller.IndexedLogsByBlockRange(int64(changedInBlock), int64(changedInBlock), FeedScopedConfigSet, lp.addr, feedIdTopicIndex, []common.Hash{lp.feedId}, pg.WithParentCtx(ctx))
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
	lp.lggr.Infow("LatestConfig", "latestConfig", latestConfigSet)
	return latestConfigSet.ContractConfig, nil
}

// LatestBlockHeight returns the latest block height from the logs
func (lp *ConfigPoller) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	latest, err := lp.destChainLogPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return uint64(latest), nil
}
