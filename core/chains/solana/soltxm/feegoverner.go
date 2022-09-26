package soltxm

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// TODO: constants should be configurable
const (
	feeGovernerPoll = 5 * time.Second

	// fees in microLamports
	minFee  uint64 = 0
	maxFee  uint64 = 1_000_000
	feeStep uint64 = 10
)

func (txm *Txm) feeGovernor(ctx context.Context) {
	defer txm.done.Done()

	tick := time.After(0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			// check a metric
			increase := false
			decrease := false

			// determine if fees can be decrease or should be raised
			if increase {
				// set new fee if fee <= maxFee or new fee > old fee (prevent overflow)
				if fee := txm.GetFee() + feeStep; fee <= maxFee && fee > txm.GetFee() {
					txm.SetFee(fee)
					txm.lggr.Infow("solana fee governer increased fee", "fee", txm.GetFee())
				} else {
					txm.lggr.Warnw("solana fee governer cannot bump fee higher than maxFee", "maxfee", maxFee)
				}
			}
			if decrease {
				// set new fee if fee >= minFee, and new fee is less than old fee (prevent overflow)
				if fee := txm.GetFee() - feeStep; fee >= minFee && fee < txm.GetFee() {
					txm.SetFee(fee)
					txm.lggr.Infow("solana fee governer decreased fee", "fee", txm.GetFee())
				}
			}

		}
		tick = time.After(utils.WithJitter(feeGovernerPoll))

	}
}