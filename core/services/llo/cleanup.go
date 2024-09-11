package llo

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/mercurytransmitter"
)

func Cleanup(ctx context.Context, lp LogPoller, addr common.Address, donID uint32, ds sqlutil.DataSource, chainSelector uint64) error {
	if (addr != common.Address{} && donID > 0) {
		if err := lp.UnregisterFilter(ctx, filterName(addr, donID)); err != nil {
			return fmt.Errorf("failed to unregister filter: %w", err)
		}
		orm := NewORM(ds, chainSelector)
		if err := orm.CleanupChannelDefinitions(ctx, addr, donID); err != nil {
			return fmt.Errorf("failed to cleanup channel definitions: %w", err)
		}
	}
	torm := mercurytransmitter.NewORM(ds, donID)
	if err := torm.Cleanup(ctx); err != nil {
		return fmt.Errorf("failed to cleanup transmitter: %w", err)
	}
	return nil
}
