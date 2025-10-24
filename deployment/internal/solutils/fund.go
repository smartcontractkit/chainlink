package solutils

import (
	"context"
	"errors"
	"time"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
)

// FundAccounts funds the given accounts with the given amount of SOL and waits for confirmation.
func FundAccounts(
	ctx context.Context, solanaGoClient *solRpc.Client, accounts []solana.PublicKey, solAmount uint64,
) error {
	sigs := make([]solana.Signature, 0, len(accounts))
	for _, account := range accounts {
		sig, err := solanaGoClient.RequestAirdrop(
			ctx, account, solAmount*solana.LAMPORTS_PER_SOL, solRpc.CommitmentFinalized,
		)
		if err != nil {
			return err
		}

		sigs = append(sigs, sig)
	}

	const timeout = 100 * time.Second
	const pollInterval = 500 * time.Millisecond

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	remaining := len(sigs)
	for remaining > 0 {
		select {
		case <-timeoutCtx.Done():
			return errors.New("unable to fund transaction within timeout")
		case <-ticker.C:
			statusRes, sigErr := solanaGoClient.GetSignatureStatuses(ctx, true, sigs...)
			if sigErr != nil {
				return sigErr
			}
			if statusRes == nil {
				return errors.New("status response is nil")
			}
			if statusRes.Value == nil {
				return errors.New("status response value is nil")
			}

			unfinalizedCount := 0
			for _, res := range statusRes.Value {
				if res == nil || res.ConfirmationStatus == solRpc.ConfirmationStatusFinalized {
					unfinalizedCount++
				}
			}
			remaining = unfinalizedCount
		}
	}

	return nil
}
