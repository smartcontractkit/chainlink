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

	if err := validateDedupedItems(dedupKeyNetworkLiquidity, obs.LiquidityPerChain...); err != nil {
		return fmt.Errorf("invalid LiquidityPerChain: %w", err)
	}
	if err := validateDedupedItems(dedupKeyTransfer, obs.ResolvedTransfers...); err != nil {
		return fmt.Errorf("invalid ResolvedTransfers: %w", err)
	}
	if err := validateDedupedItems(dedupKeyPendingTransfer, obs.PendingTransfers...); err != nil {
		return fmt.Errorf("invalid PendingTransfers: %w", err)
	}
	if err := validateDedupedItems(dedupKeyTransfer, obs.InflightTransfers...); err != nil {
		return fmt.Errorf("invalid InflightTransfers: %w", err)
	}
	if err := validateDedupedItems(dedupKeyEdge, obs.Edges...); err != nil {
		return fmt.Errorf("invalid Edges: %w", err)
	}
	if err := validateDedupedItems(dedupKeyConfigDigest, obs.ConfigDigests...); err != nil {
		return fmt.Errorf("invalid ConfigDigests: %w", err)
	}

	return nil
}

func dedupKeyNetworkLiquidity(liq models.NetworkLiquidity) string {
	return fmt.Sprintf("%d", liq.Network)
}

func dedupKeyPendingTransfer(pt models.PendingTransfer) string {
	return pt.ID // TODO: check if we need to use dedupKeyTransfer
}

func dedupKeyTransfer(t models.Transfer) string {
	return fmt.Sprintf("%d-%d-%s-%s-%d", t.From, t.To, t.LocalTokenAddress.String(), t.RemoteTokenAddress.String(), t.Stage)
}

func dedupKeyEdge(e models.Edge) string {
	return fmt.Sprintf("%d-%d", e.Source, e.Dest)
}

func dedupKeyConfigDigest(obs models.ConfigDigestWithMeta) string {
	return fmt.Sprintf("%d", obs.NetworkSel) // we only allow 1 config digest per network
}

// validateDedupedItems checks if there are any duplicated items in the provided slice.
func validateDedupedItems[T any](keyFn func(T) string, items ...T) error {
	existing := map[string]bool{}
	for _, item := range items {
		k := keyFn(item)
		if existing[k] {
			return fmt.Errorf("duplicated item")
		}
		existing[k] = true
	}
	return nil
}
