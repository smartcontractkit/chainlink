package solana

import (
	"context"
	"errors"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
)

var (
	configVersion uint8 = 1
)

type StateCache struct {
	*client.Cache[State]
}

func NewStateCache(stateID solana.PublicKey, chainID string, cfg config.Config, reader client.Reader, lggr logger.Logger) *StateCache {
	name := "ocr2_median_state"
	getter := func(ctx context.Context) (State, uint64, error) {
		return GetState(ctx, reader, stateID, cfg.Commitment())
	}
	return &StateCache{client.NewCache(name, stateID, chainID, cfg, getter, logger.With(lggr, "cache", name))}
}

func GetState(ctx context.Context, reader client.AccountReader, account solana.PublicKey, commitment rpc.CommitmentType) (State, uint64, error) {
	res, err := reader.GetAccountInfoWithOpts(ctx, account, &rpc.GetAccountInfoOpts{
		Commitment: commitment,
		Encoding:   "base64",
	})
	if err != nil {
		return State{}, 0, fmt.Errorf("failed to fetch state account at address '%s': %w", account.String(), err)
	}

	// check for nil pointers
	if res == nil || res.Value == nil || res.Value.Data == nil {
		return State{}, 0, errors.New("nil pointer returned in GetState.GetAccountInfoWithOpts")
	}

	var state State
	if err := bin.NewBinDecoder(res.Value.Data.GetBinary()).Decode(&state); err != nil {
		return State{}, 0, fmt.Errorf("failed to decode state account data: %w", err)
	}

	// validation for config version
	if configVersion != state.Version {
		return State{}, 0, fmt.Errorf("decoded config version (%d) does not match expected config version (%d)", state.Version, configVersion)
	}

	blockNum := res.RPCContext.Context.Slot
	return state, blockNum, nil
}
