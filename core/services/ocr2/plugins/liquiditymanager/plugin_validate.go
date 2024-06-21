package liquiditymanager

import (
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func (p *Plugin) ValidateObservation(outctx ocr3types.OutcomeContext, query ocrtypes.Query, ao ocrtypes.AttributedObservation) error {
	// todo: improve logging - including duration of each phase, etc.
	p.lggr.Infow("in validate observation", "seqNr", outctx.SeqNr, "phase", "ValidateObservation")

	obs, err := models.DecodeObservation(ao.Observation)
	if err != nil {
		return fmt.Errorf("invalid observation: %w", err)
	}

	if err := validateItems(dedupKeyNetworkLiquidity, obs.LiquidityPerChain, validateNetworkLiquidity); err != nil {
		return fmt.Errorf("invalid LiquidityPerChain: %w", err)
	}
	if err := validateItems(dedupKeyTransfer, obs.ResolvedTransfers, validateTransfer); err != nil {
		return fmt.Errorf("invalid ResolvedTransfers: %w", err)
	}
	if err := validateItems(dedupKeyPendingTransfer, obs.PendingTransfers, validatePendingTransfer); err != nil {
		return fmt.Errorf("invalid PendingTransfers: %w", err)
	}
	if err := validateItems(dedupKeyTransfer, obs.InflightTransfers, validateTransfer); err != nil {
		return fmt.Errorf("invalid InflightTransfers: %w", err)
	}
	if err := validateItems(dedupKeyEdge, obs.Edges); err != nil {
		return fmt.Errorf("invalid Edges: %w", err)
	}
	if err := validateItems(dedupKeyConfigDigest, obs.ConfigDigests); err != nil {
		return fmt.Errorf("invalid ConfigDigests: %w", err)
	}

	return nil
}

func dedupKeyNetworkLiquidity(liq models.NetworkLiquidity) string {
	return fmt.Sprintf("%d", liq.Network)
}

func validateNetworkLiquidity(liq models.NetworkLiquidity) error {
	if liq.Liquidity == nil {
		return fmt.Errorf("nil Liquidity")
	}
	return nil
}

func dedupKeyPendingTransfer(pt models.PendingTransfer) string {
	return pt.ID // TODO: check if we need to use dedupKeyTransfer
}

func validatePendingTransfer(t models.PendingTransfer) error {
	if t.Transfer.Amount == nil {
		return fmt.Errorf("nil Amount")
	}
	if t.Transfer.NativeBridgeFee == nil {
		return fmt.Errorf("nil NativeBridgeFee")
	}
	return nil
}

func dedupKeyTransfer(t models.Transfer) string {
	return fmt.Sprintf("%d-%d-%s-%s-%d", t.From, t.To, t.LocalTokenAddress.String(), t.RemoteTokenAddress.String(), t.Stage)
}

func validateTransfer(t models.Transfer) error {
	if t.Amount == nil {
		return fmt.Errorf("nil Amount")
	}
	if t.NativeBridgeFee == nil {
		return fmt.Errorf("nil NativeBridgeFee")
	}
	return nil
}

func dedupKeyEdge(e models.Edge) string {
	return fmt.Sprintf("%d-%d", e.Source, e.Dest)
}

func dedupKeyConfigDigest(obs models.ConfigDigestWithMeta) string {
	return fmt.Sprintf("%d", obs.NetworkSel) // we only allow 1 config digest per network
}

// validateItems verifies there are no duplicated items and runs the given validate functions against all items.
func validateItems[T any](keyFn func(T) string, items []T, validateFns ...func(T) error) error {
	existing := map[string]bool{}
	for _, item := range items {
		k := keyFn(item)
		if existing[k] {
			return fmt.Errorf("duplicated item (%s)", k)
		}
		for _, validateFn := range validateFns {
			if err := validateFn(item); err != nil {
				return fmt.Errorf("invalid item (%s): %w", k, err)
			}
		}
		existing[k] = true
	}
	return nil
}
