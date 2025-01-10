package v1_5_0

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

type CommitStore struct {
	*v1_2_0.CommitStore
	commitStore *commit_store.CommitStore
}

func (c *CommitStore) GetCommitStoreStaticConfig(ctx context.Context) (cciptypes.CommitStoreStaticConfig, error) {
	staticConfig, err := c.commitStore.GetStaticConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return cciptypes.CommitStoreStaticConfig{}, err
	}
	return cciptypes.CommitStoreStaticConfig{
		ChainSelector:       staticConfig.ChainSelector,
		SourceChainSelector: staticConfig.SourceChainSelector,
		OnRamp:              cciptypes.Address(staticConfig.OnRamp.String()),
		ArmProxy:            cciptypes.Address(staticConfig.RmnProxy.String()),
	}, nil
}

func (c *CommitStore) IsDown(ctx context.Context) (bool, error) {
	unPausedAndNotCursed, err := c.commitStore.IsUnpausedAndNotCursed(&bind.CallOpts{Context: ctx})
	if err != nil {
		return true, err
	}
	return !unPausedAndNotCursed, nil
}

func NewCommitStore(
	lggr logger.Logger,
	addr common.Address,
	ec client.Client,
	lp logpoller.LogPoller,
	feeEstimatorConfig ccipdata.FeeEstimatorConfigReader,
) (*CommitStore, error) {
	v120, err := v1_2_0.NewCommitStore(lggr, addr, ec, lp, feeEstimatorConfig)
	if err != nil {
		return nil, err
	}

	commitStore, err := commit_store.NewCommitStore(addr, ec)
	if err != nil {
		return nil, err
	}

	return &CommitStore{
		commitStore: commitStore,
		CommitStore: v120,
	}, nil
}
