package handler

import (
	"context"
	"log"
	"math/big"
)

// UpkeepHistory prints the checkUpkeep status and keeper responsibility for a given upkeep in a set block range.
// Legacy keeper registry 1.1/1.2 support has been removed.
func (k *Keeper) UpkeepHistory(ctx context.Context, upkeepId *big.Int, from, to, gasPrice uint64) {
	_, _, _, _, _ = ctx, upkeepId, from, to, gasPrice
	log.Fatal("upkeep-history was only implemented for keeper registry 1.1/1.2, which are no longer supported")
}
