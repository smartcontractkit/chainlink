package por

import (
	"context"
	"encoding/json"
	"fmt"

	"slices"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/quorumhelper"
)

type PorReportingPluginFactory struct {
	Logger          commontypes.Logger
	ExternalAdapter ExternalAdapter
	ContractReader  ContractReader
	ReportMarshaler ReportMarshaler
}

var _ ocr3types.ReportingPluginFactory[ChainSelector] = &PorReportingPluginFactory{}

func (f *PorReportingPluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[ChainSelector], ocr3types.ReportingPluginInfo, error) {
	porOffchainConfig, err := DeserializePorOffchainConfig(config.OffchainConfig)
	if err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, fmt.Errorf("could not deserialize offchain config: %w", err)
	}

	maxChains := porOffchainConfig.MaxChains

	mintablesMapLength := maxChains * (8 + 32)   // 8 bytes for BlockNumber, 32 bytes for big.Int (assuming max 256-bit big.Int size)
	honestBlocksMapLength := maxChains * (8 + 8) // 8 bytes for BlockNumber, 8 bytes for ChainSelector (64-bit integers)

	porObservationLength := mintablesMapLength + honestBlocksMapLength
	porPluginOutcomeLength := mintablesMapLength + honestBlocksMapLength + 1

	maxObservationLength := int(max((3*porObservationLength)/2, 128)) // estimate the size of the JSON-encoded observation
	maxOutcomeLength := int(max((3*porPluginOutcomeLength)/2, 128))   // estimate the size of the JSON-encoded outcome

	rm := f.ReportMarshaler
	if rm == nil {
		rm = NewMockReportMarshaler()
	}
	maxReportLength := rm.MaxReportSize(ctx)

	limits := ocr3types.ReportingPluginLimits{
		0,
		maxObservationLength,
		maxOutcomeLength,
		maxReportLength,
		int(maxChains),
	}

	ea := f.ExternalAdapter
	if ea == nil {
		ea = NewMockExternalAdapterImpl()
	}

	cr := f.ContractReader
	if cr == nil {
		cr = NewMockContractReader(config.ConfigDigest)
	}

	return &porReportingPlugin{
			maxChains,
			ea,
			cr,
			rm,
			config,
			f.Logger,
		}, ocr3types.ReportingPluginInfo{
			"PorReportingPluginV1",
			limits,
		},
		nil
}

type porReportingPlugin struct {
	maxChains       uint32
	externalAdapter ExternalAdapter
	contractReader  ContractReader
	reportMarshaler ReportMarshaler
	config          ocr3types.ReportingPluginConfig
	logger          commontypes.Logger
}

type porPluginObservation struct {
	QueryResponse Mintables
	LatestBlocks  Blocks
}

type porPluginOutcome struct {
	ChangedMintables        map[ChainSelector]bool // Indicates per chain if the mintable changed from the previous round
	LatestAcquiredMintables Mintables
	HonestBlocks            Blocks
}

var _ ocr3types.ReportingPlugin[ChainSelector] = &porReportingPlugin{}

// We do not use the Query method in this plugin. The query is implied by the honest blocks from the previous round.
func (p *porReportingPlugin) Query(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Query, error) {
	return nil, nil
}

// The Observation method is called by the OCR node to obtain an observation for this oracle for the current round.
// It queries the external adapter for mintable amounts based on the honest blocks from the previous round.
// (See the Outcome method for how the honest blocks are determined.)
//
// If the external adapter is up to date on all chains corresponding to the honest blocks query, the mintables will be
// non-nil and contain the mintable amounts for each chain that is tracked by the plugin.
//
// If the external adapter is not tracking the same set of chains as the honest blocks, or if it is not up to date
// with the latest blocks, it will return a nil mintables vector. This is expected behavior during an upgrade (adding a new
// chain), where the external adapter might not yet be aware of the new chain, or it might not have reached the latest
// block for the chains it is tracking.
//
// The latest blocks are also returned by the external adapter, which are used to track the latest blocks for
// each chain. These are used to calculate honest blocks for the next round's query to the external adapter (pipelining).
// The latest blocks should include all chains that are globally tracked by the plugin, i.e., in the honestBlocks
// vector. If it internally has not started tracking a new chain yet, it should return the latest block number as 0 for that chain.
// It is valid for the latest blocks to include more chains than the honest blocks, as long as it does
// not exceed the maximum number of chains. Including a chain not in honestBlocks suggests that the external adapter has
// started tracking a new chain and wants to extend the honestBlocks vector for the next round.
func (p *porReportingPlugin) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query) (types.Observation, error) {
	honestBlockQuery, err := p.getOrInitializeHonestBlocks(outctx)
	if err != nil {
		return nil, fmt.Errorf("could not get 'honest blocks': %w", err)
	}

	payload, err := p.externalAdapter.GetPayload(ctx, honestBlockQuery)
	if err != nil {
		return nil, fmt.Errorf("could not get payload from the external adapter: %w", err)
	}

	if len(payload.Mintables) == 0 && len(honestBlockQuery) > 0 {
		p.logger.Warn("PorReportingPlugin: external adapter could not respond to 'honest blocks' query. This could be due to a new chain being added, or the external adapter being outdated in relation to the query.", commontypes.LogFields{
			"Query": honestBlockQuery,
		})
	}

	if err := p.checkValidLatestBlocks(payload.LatestBlocks, honestBlockQuery); err != nil {
		return nil, fmt.Errorf("invalid 'latest blocks' in payload: %w", err)
	}

	if err = checkValidMintables(payload.Mintables, honestBlockQuery); err != nil {
		return nil, fmt.Errorf("invalid 'mintables' in payload: %w", err)
	}

	ppo := porPluginObservation{
		payload.Mintables,
		payload.LatestBlocks,
	}

	observation, err := serializePorPluginObservation(ppo)
	if err != nil {
		return nil, fmt.Errorf("could not serialize observation: %w", err)
	}

	return observation, nil
}

// ValidateObservation checks if the observation is valid.
// We only reject observations which are *guaranteed* to come from dishonest oracles.
func (p *porReportingPlugin) ValidateObservation(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query, ao types.AttributedObservation) error {
	honestBlocks, err := p.getOrInitializeHonestBlocks(outctx)
	if err != nil {
		return fmt.Errorf("could not get 'honest blocks': %w", err)
	}

	// Deserialize the observation
	ppo, err := deserializePorPluginObservation(ao.Observation)
	if err != nil {
		return fmt.Errorf("could not deserialize observation: %w", err)
	}

	if err = p.checkValidLatestBlocks(ppo.LatestBlocks, honestBlocks); err != nil {
		return fmt.Errorf("invalid observation 'latest blocks': %w", err)
	}

	if err = checkValidMintables(ppo.QueryResponse, honestBlocks); err != nil {
		return fmt.Errorf("invalid observation 'mintables': %w", err)
	}

	return nil
}

// ObservationQuorum checks if the number of observations reaches the minimum necessary quorum to attempt
// to acquire new mintables for the current round AND calculate honest blocks for the next round.
func (p *porReportingPlugin) ObservationQuorum(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query, aos []types.AttributedObservation) (bool, error) {
	return quorumhelper.ObservationCountReachesObservationQuorum(quorumhelper.QuorumTwoFPlusOne, p.config.N, p.config.F, aos), nil
}

// The Outcome method processes the observations and generates a new outcome, from which reports will be obtained.
//
// In the common case, the Outcome tries to simultaneously do two things (pipeline):
//  1. Deduce the honest mintables, which are mintable amounts for each chain that are reported by at least F+1 oracles.
//     The mintable amounts are calculated based on the honest blocks determined in the previous round (for this round).
//  2. Deduce the latest honest blocks from the observations, which are used to track the latest blocks for each chain.
//     These are used in the next round to query the external adapter for mintable amounts (see above).
//
// Exceptions:
//   - If the latest honest blocks are extended with new chains, this indicates that the oracles are undergoing an upgrade,
//     i.e., they are starting to track a new chain (all oracles will eventually do so). In this situation, we ALWAYS update the
//     honest blocks, even if no new mintables were acquired, because the upgrade is expected to be completed eventually and is
//     only a temporary, infrequent event.
//   - If no new mintables were acquired (and there is no new chain) we keep the previous round's mintables. This ensures that
//     we (eventually) acquire mintables, since all honest oracles will (eventually) reach the fixed mintables (as we stop
//     updating it until we succeed).
//
// The outcome also keeps track of which chains had their mintable amounts or block numbers changed compared to the previous round.
// This is used to determine which reports to generate in the Reports method and avoid generating reports for chains that did not change.
func (p *porReportingPlugin) Outcome(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query, aos []types.AttributedObservation) (ocr3types.Outcome, error) {
	// Deserialize the previous outcome in outcome context (outctx)
	prevOutcome, err := previousPorPluginOutcome(outctx)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize previous outcome: %w", err)
	}

	// Deserialize the observations
	ppos := make([]porPluginObservation, 0, len(aos))
	for i, ao := range aos {
		po, err := deserializePorPluginObservation(ao.Observation)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize %d-th observation from sender %v: %w", i, ao.Observer, err)
		}
		ppos = append(ppos, po)
	}

	newHonestBlocks := p.deduceLatestHonestBlocks(prevOutcome.HonestBlocks, ppos)

	honestMintables, err := p.deduceHonestMintables(outctx, ppos)
	if err != nil {
		return nil, fmt.Errorf("error deducing honest mintables: %w", err)
	}

	newOutcome := porPluginOutcome{
		ChangedMintables:        make(map[ChainSelector]bool),
		LatestAcquiredMintables: honestMintables,
		HonestBlocks:            newHonestBlocks,
	}

	// If we did not acquire new mintables, we keep the previous round's mintables
	if honestMintables == nil {
		p.logger.Info("💥💥💥💥💥 PorReportingPlugin: could not acquire mintables", commontypes.LogFields{
			"HonestBlocks":            newOutcome.HonestBlocks,
			"LatestAcquiredMintables": newOutcome.LatestAcquiredMintables,
		})

		newOutcome.LatestAcquiredMintables = prevOutcome.LatestAcquiredMintables
	} else {
		p.logger.Info("🚀🚀🚀🚀🚀 PorReportingPlugin: acquired mintables, generating a new outcome", commontypes.LogFields{
			"HonestBlocks":            newOutcome.HonestBlocks,
			"LatestAcquiredMintables": newOutcome.LatestAcquiredMintables,
		})
	}

	// If we did not acquire new mintables and there is no new chain to track, we keep the previous round's honest blocks.
	if honestMintables == nil && len(newHonestBlocks) == len(prevOutcome.HonestBlocks) {
		newOutcome.HonestBlocks = prevOutcome.HonestBlocks
	}

	if len(newHonestBlocks) > len(prevOutcome.HonestBlocks) {
		p.logger.Info("🆕🆕🆕🆕🆕 PorReportingPlugin: extended honest blocks with new chains", commontypes.LogFields{
			"NewHonestBlocks":      newHonestBlocks,
			"PreviousHonestBlocks": prevOutcome.HonestBlocks,
		})
	}

	newOutcome.ChangedMintables = getChanged(prevOutcome.LatestAcquiredMintables, honestMintables)

	outcome, err := serializePorPluginOutcome(newOutcome)
	if err != nil {
		return nil, fmt.Errorf("could not serialize outcome: %w", err)
	}

	return outcome, nil
}

// The Reports method generates reports based on the outcome of the previous round.
// It creates a report for each chain that had its mintable amount or block number changed compared to the previous round.
// The reports are created using the reportMarshaler, which serializes the report data into a format suitable for transmission.
func (p *porReportingPlugin) Reports(ctx context.Context, seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportPlus[ChainSelector], error) {
	// Deserialize the outcome
	ppo, err := deserializePorPluginOutcome(outcome)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize outcome: %w", err)
	}

	// Create a report for each chain
	reports := make([]ocr3types.ReportPlus[ChainSelector], 0, len(ppo.LatestAcquiredMintables))
	for chain, changed := range ppo.ChangedMintables {
		if !changed {
			continue // Skip chains that did not change
		}

		pair := ppo.LatestAcquiredMintables[chain]
		report := PorReport{
			p.config.ConfigDigest,
			seqNr,
			pair.Block,
			pair.Mintable,
		}

		encodedReport, err := p.reportMarshaler.Serialize(ctx, chain, report)
		if err != nil {
			return nil, fmt.Errorf("could not encode block-mintable pair (%v): %w", pair, err)
		}

		// Create a report for the chain
		reports = append(reports, ocr3types.ReportPlus[ChainSelector]{
			ReportWithInfo: ocr3types.ReportWithInfo[ChainSelector]{
				encodedReport,
				chain,
			},
			TransmissionScheduleOverride: nil,
		})
	}

	return reports, nil
}

// ShouldAcceptAttestedReport is called by the OCR node to determine if the attested report should be accepted.
// In this plugin, we always accept the attested report, as we do not have any specific conditions for acceptance.
func (p *porReportingPlugin) ShouldAcceptAttestedReport(context.Context, uint64, ocr3types.ReportWithInfo[ChainSelector]) (bool, error) {
	return true, nil
}

// ShouldTransmitAcceptedReport is called by the OCR node to determine if the accepted report should be transmitted on-chain.
// In the future, this responsibility might be entirely moved to the transmission infrastucture, but for now,
// we implement it in the plugin to ensure that the report is only transmitted if it is valid and has not been
// transmitted already.
func (p *porReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, report ocr3types.ReportWithInfo[ChainSelector]) (bool, error) {
	chain := report.Info
	details, err := p.contractReader.GetLatestTransmittedReportDetails(ctx, chain)
	if err != nil {
		return false, fmt.Errorf("could not get latest transmission details from chain: %w", err)
	}
	if details.ConfigDigest.Hex() != p.config.ConfigDigest.Hex() {
		return false, fmt.Errorf("config digest mismatch; expected: %v, got: %v", p.config.ConfigDigest.Hex(), details.ConfigDigest.Hex())
	}
	// This or a later report is already posted on-chain
	if details.SeqNr >= seqNr {
		return false, nil
	}

	return true, nil
}

func (p *porReportingPlugin) Close() error {
	return nil
}

// checkValidLatestBlocks checks that the latestBlocks are a possible response by a correct EA for a given
// honestBlocks query.
//
// latestBlocks should contain the latest blocks for all chains that are globally tracked by the plugin, i.e.,
// in the honestBlocks vector.
//
// latestBlocks including more chains than honestBlocks is valid behavior, as long as it does not exceed the
// maximum number of chains. Including a chain not in honestBlocks suggests that the EA has started tracking a new
// chain and wants to extend the honestBlocks vector for the next round.
func (p *porReportingPlugin) checkValidLatestBlocks(latestBlocks Blocks, honestBlocks Blocks) error {
	if len(latestBlocks) > int(p.maxChains) {
		return fmt.Errorf("'latest blocks' contains too many chains (%d), max allowed is %d", len(latestBlocks), p.maxChains)
	}

	for chain := range honestBlocks {
		if _, ok := latestBlocks[chain]; !ok {
			return fmt.Errorf("'latest blocks' missing info on tracked chain with ID: %v", chain)
		}
	}
	return nil
}

// checkValidMintables checks that the mintables are a possible response by a correct EA for a given honestBlocks“ query.
//
// nil mintables are considered valid, as they can happen during normal behavior or during an upgrade (adding a new chain):
// - during normal behavior, the EA *should* return nil mintables if, for at least one (chain -> block) pair in
// honestBlocks, it has not reached that block yet (so it cannot accurately calculate mintables for it)
// - during an upgrade (a new chain being added to the system), the EA *should* return nil if there is a mismatch between
// the chains the EA is tracking and the honest blocks, which might happen temporarily until the upgrade is complete.
//
// If the mintables are not nil, then they should contain exactly the same chains as honestBlocks, with block numbers not exceeding
// the corresponding honest blocks.
func checkValidMintables(mintables Mintables, honestBlocks Blocks) error {
	if len(mintables) == 0 {
		return nil
	}

	// Check that the mintables are for the honest blocks
	for chain, pair := range mintables {
		if block, ok := honestBlocks[chain]; !ok {
			return fmt.Errorf("'mintables' includes info on untracked chain with ID: %v", chain)
		} else if pair.Block > block {
			return fmt.Errorf("'mintables' bad block number; expected no more than: %v, got: %v", block, pair.Block)
		}
	}
	for chain := range honestBlocks {
		if _, ok := mintables[chain]; !ok {
			return fmt.Errorf("'mintables' missing info on tracked chain with ID: %v", chain)
		}
	}

	return nil
}

// This function deduces the honest mintables from the plugin observations (ppos).
// An honest mintables vector is one that is reported by at least F+1 oracles.
// It returns nil if no honest mintables could be deduced, which is expected if:
// - It is the first round of the plugin (no honest blocks yet),
// - Many oracles could not respond to the query (e.g., due to a new chain being added or their EA being outdated),
// - Mintable amounts differing, e.g., due to temporary inconsistencies between reserve amounts at different oracles.
// In these cases, the plugin will not be able to deduce honest mintables (we will retry in the next OCR round).
func (p *porReportingPlugin) deduceHonestMintables(outctx ocr3types.OutcomeContext, ppos []porPluginObservation) (Mintables, error) {
	honestBlocks, err := p.getOrInitializeHonestBlocks(outctx)
	if err != nil {
		return nil, fmt.Errorf("could not get honest blocks: %w", err)
	}

	if len(honestBlocks) == 0 {
		return nil, nil // No honest blocks to deduce mintables from, this is expected in the first round
	}

	mintablesFrequencyMap := make(map[string]int, len(ppos))
	for _, po := range ppos {
		if len(po.QueryResponse) != len(honestBlocks) {
			continue
		}

		uniqueEncoding, err := po.QueryResponse.toString()
		if err != nil {
			return nil, fmt.Errorf("could not convert query response to string: %w", err)
		}

		if _, exists := mintablesFrequencyMap[string(uniqueEncoding)]; exists {
			mintablesFrequencyMap[string(uniqueEncoding)]++
		} else {
			mintablesFrequencyMap[string(uniqueEncoding)] = 1
		}
	}

	// Find mintables with enough support (seen by at least F+1 oracles)
	maxMintablesWithEnoughSupport := p.config.N / (p.config.F + 1)
	mintablesWithEnoughSupport := make([]Mintables, 0, maxMintablesWithEnoughSupport)
	for mintablesAsString, count := range mintablesFrequencyMap {
		if count <= p.config.F {
			// Not enough support, skip this mintables vector
			continue
		}

		// If there is a mintable amount that more than F oracles reported, it is guaranteed to be honest.
		mintables, err := mintablesFromString(mintablesAsString)
		if err != nil {
			// This error should be impossible, as the serialization comes from at least one honest oracle.
			return nil, fmt.Errorf("could not deserialize mintables from string: %w", err)
		}
		mintablesWithEnoughSupport = append(mintablesWithEnoughSupport, mintables)
	}

	if len(mintablesWithEnoughSupport) == 0 {
		p.logger.Warn("PorReportingPlugin: no mintable vectors with enough support found", commontypes.LogFields{
			"MintablesFrequencyMap": mintablesFrequencyMap,
		})
		return nil, nil // No mintables with enough support found, this is expected if the EA could not respond to the query
	}

	if len(mintablesWithEnoughSupport) > 1 {
		p.logger.Warn("PorReportingPlugin: multiple mintable vectors with enough support found, picking the first one", commontypes.LogFields{
			"MintableVectors": mintablesWithEnoughSupport,
		})
	}

	return mintablesWithEnoughSupport[0], nil
}

// deduceLatestHonestBlocks deduces the latest honest blocks from the plugin observations (ppos).
// This is done in two steps:
//  1. For each chain that is already tracked, we deduce the latest honest block by taking the maximum of the latest blocks
//     reported by the oracles, ignoring the top F (which could be overestimates), and ensuring that it is at least as high
//     as the original honest block.
//  2. For each chain that is not currently tracked, we check if at least F+1 oracles want to start tracking it (=> at least one
//     honest oracle has started tracking, so the "upgrade" is true).
//     If so, we take the maximum of the latest blocks reported by the oracles, ignoring the top F (which could be overestimates).
//     (A possible more conservative alternative is to initially set a new chain's latest block to 0)
//
// Note: during an upgrade (adding a new chain), honest blocks vector will be extended with new chains, which might happen
// before all oracles start tracking those new chains in their EA. For oracles missing the new chain, it is expected behavior
// to return a nil mintables vector in this case. For the latestBlocks, the new chains must be included with
// the latest block number set to 0.
func (p *porReportingPlugin) deduceLatestHonestBlocks(honestBlocks Blocks, ppos []porPluginObservation) Blocks {
	newHonestBlocks := make(Blocks)

	// First, we deduce the latest honest blocks for each chain that is already tracked.
	for chain := range honestBlocks {
		numbers := []BlockNumber{}
		for _, ppo := range ppos {
			number := ppo.LatestBlocks[chain]
			numbers = append(numbers, number)
		}

		slices.Sort(numbers)

		// Among |pos| >= 2F+1 observations, it is safe to pick any block index where N-F > X >= F.
		// In this case, we are choosing the highest block number among the safe options.
		// (In practice, the number of faults will often be zero, to in the next round the top F+1
		// fastest oracles should be able to answer the query corresponding the F+1th fastest oracle.)
		safeBlock := numbers[p.config.N-p.config.F-1]
		originalHonestBlock := honestBlocks[chain]

		// Ensure monotonicity of the honest blocks
		newHonestBlocks[chain] = max(safeBlock, originalHonestBlock)
	}

	// Second, we check if enough oracles want to start tracking a new chain. For each chain not in honestBlocks,
	// we first count how many oracles want to start tracking it and what block number they suggest.
	newChains := make(map[ChainSelector][]BlockNumber)
	for _, ppo := range ppos {
		for chain, blockNumber := range ppo.LatestBlocks {
			if _, ok := honestBlocks[chain]; !ok {
				// This chain is not currently tracked, so we count how many oracles want to start tracking it
				// and what block number they suggest.
				newChains[chain] = append(newChains[chain], blockNumber)
			}
		}
	}

	// Then we check if enough oracles (> F) want to start tracking each new chain to avoid malicious suggestions.
	for chain, blockNumbers := range newChains {
		if len(blockNumbers) > int(p.config.F) {
			// It is safe to underestimate the block number. In theory, we could even just start at 0.
			// However, it is *not* safe to overestimate.
			// Here, we use highest *assuredly honest* block number suggested by the oracles, by eliminating the top F.
			slices.Sort(blockNumbers)
			newHonestBlocks[chain] = blockNumbers[len(blockNumbers)-p.config.F-1]
		}
	}

	return newHonestBlocks
}

// getChanged compares the old and new Mintables and returns a map indicating which chains have changed.
// A chain is considered changed if its mintable amount or block number has changed compared to the old Mintables.
func getChanged(old, new Mintables) map[ChainSelector]bool {
	changed := make(map[ChainSelector]bool, len(new))

	for chain, newPair := range new {
		oldPair, ok := old[chain]
		if !ok || oldPair.Mintable != newPair.Mintable || oldPair.Block != newPair.Block {
			changed[chain] = true
		}
	}

	return changed
}

func deserializePorPluginOutcome(outcome []byte) (porPluginOutcome, error) {
	var ppo porPluginOutcome
	err := json.Unmarshal(outcome, &ppo)
	if err != nil {
		return porPluginOutcome{}, err
	}
	return ppo, err
}

func serializePorPluginOutcome(ppo porPluginOutcome) (ocr3types.Outcome, error) {
	return json.Marshal(ppo)
}

func previousPorPluginOutcome(outctx ocr3types.OutcomeContext) (porPluginOutcome, error) {
	if outctx.SeqNr == 1 {
		return porPluginOutcome{
			ChangedMintables:        make(map[ChainSelector]bool),
			LatestAcquiredMintables: make(Mintables),
			HonestBlocks:            make(Blocks),
		}, nil
	} else {
		return deserializePorPluginOutcome(outctx.PreviousOutcome)
	}
}

// Deserialize the previous outcome in outcome context (outctx)
// If we are in the first round, this will be empty, which is intentional:
// We will not have any honest blocks to obtain mintables for in the first round,
// but we will extend the honestBlocks with new chains to track at the end for the next round.
// (This is done in the Outcome method.)
func (p *porReportingPlugin) getOrInitializeHonestBlocks(outctx ocr3types.OutcomeContext) (Blocks, error) {
	prevOutcome, err := previousPorPluginOutcome(outctx)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize previous outcome: %w", err)
	}

	return prevOutcome.HonestBlocks, nil
}

func deserializePorPluginObservation(raw []byte) (porPluginObservation, error) {
	var ppo porPluginObservation
	err := json.Unmarshal(raw, &ppo)
	if err != nil {
		return porPluginObservation{}, err
	}
	return ppo, nil
}

func serializePorPluginObservation(rq porPluginObservation) (types.Observation, error) {
	return json.Marshal(rq)
}
