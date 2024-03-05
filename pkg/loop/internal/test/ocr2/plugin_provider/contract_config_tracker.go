package pluginprovider

import (
	"context"
	"fmt"
	"reflect"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
)

type staticConfigTrackerConfig struct {
	contractConfig libocr.ContractConfig
	configDigest   libocr.ConfigDigest
	changedInBlock uint64
	blockHeight    uint64
}

type staticContractConfigTracker struct {
	staticConfigTrackerConfig
}

var _ testtypes.ContractConfigTrackerEvaluator = staticContractConfigTracker{}

func (s staticContractConfigTracker) Notify() <-chan struct{} { return nil }

func (s staticContractConfigTracker) LatestConfigDetails(ctx context.Context) (uint64, libocr.ConfigDigest, error) {
	return s.changedInBlock, s.configDigest, nil
}

func (s staticContractConfigTracker) LatestConfig(ctx context.Context, cib uint64) (libocr.ContractConfig, error) {
	if s.changedInBlock != cib {
		return libocr.ContractConfig{}, fmt.Errorf("expected changed in block %d but got %d", s.changedInBlock, cib)
	}
	return s.contractConfig, nil
}

func (s staticContractConfigTracker) LatestBlockHeight(ctx context.Context) (uint64, error) {
	return s.blockHeight, nil
}

func (s staticContractConfigTracker) Evaluate(ctx context.Context, cct libocr.ContractConfigTracker) error {
	gotCIB, gotDigest, err := cct.LatestConfigDetails(ctx)
	if err != nil {
		return fmt.Errorf("failed to get LatestConfigDetails: %w", err)
	}
	if gotCIB != s.changedInBlock {
		return fmt.Errorf("expected changed in block %d but got %d", s.changedInBlock, gotCIB)
	}
	if gotDigest != s.configDigest {
		return fmt.Errorf("expected config digest %x but got %x", s.configDigest, gotDigest)
	}
	gotBlockHeight, err := cct.LatestBlockHeight(ctx)
	if err != nil {
		return fmt.Errorf("failed to get LatestBlockHeight: %w", err)
	}
	if gotBlockHeight != s.blockHeight {
		return fmt.Errorf("expected block height %d but got %d", s.blockHeight, gotBlockHeight)
	}
	gotConfig, err := cct.LatestConfig(ctx, gotCIB)
	if err != nil {
		return fmt.Errorf("failed to get LatestConfig: %w", err)
	}
	if !reflect.DeepEqual(gotConfig, s.contractConfig) {
		return fmt.Errorf("expected ContractConfig %v but got %v", s.contractConfig, gotConfig)
	}
	return nil
}
