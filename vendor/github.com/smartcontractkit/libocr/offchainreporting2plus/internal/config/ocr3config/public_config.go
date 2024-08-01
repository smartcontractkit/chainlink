package ocr3config

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/internal/byzquorum"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// PublicConfig is the configuration disseminated through the smart contract
// It's public, because anybody can read it from the blockchain
type PublicConfig struct {
	// If an epoch (driven by a leader) fails to achieve progress (generate a
	// report) after DeltaProgress, we enter a new epoch. This parameter must be
	// chosen carefully. If the duration is too short, we may keep prematurely
	// switching epochs without ever achieving any progress, resulting in a
	// liveness failure!
	DeltaProgress time.Duration
	// DeltaResend determines how often Pacemaker newepoch messages should be
	// resent, allowing oracles that had crashed and are recovering to rejoin
	// the protocol more quickly. ~30s should be a reasonable default under most
	// circumstances.
	DeltaResend time.Duration
	// If no message from the leader has been received after the epoch start plus
	// DeltaInitial, we enter a new epoch. This parameter must be
	// chosen carefully. If the duration is too short, we may keep prematurely
	// switching epochs without ever achieving any progress, resulting in a
	// liveness failure!
	DeltaInitial time.Duration
	// DeltaRound determines the minimal amount of time that should pass between
	// the start of outcome generation rounds. With OCR3 (not OCR1!) you can
	// set this value very aggressively. Note that this only provides a lower
	// bound on the round interval; actual rounds might take longer.
	DeltaRound time.Duration
	// Once the leader of a outcome generation round has collected sufficiently
	// many observations, it will wait for DeltaGrace to pass to allow slower
	// oracles to still contribute an observation before moving on to generating
	// the report. Consequently, rounds driven by correct leaders will always
	// take at least DeltaGrace.
	DeltaGrace time.Duration
	// DeltaCertifiedCommitRequest determines the duration between requests for
	// a certified commit after we have received f+1 signatures in the report
	// attestation protocol but are still missing the certified commit/outcome
	// required for validating the report signatures.
	DeltaCertifiedCommitRequest time.Duration
	// DeltaStage determines the duration between stages of the transmission
	// protocol. In each stage, a certain number of oracles (determined by S)
	// will attempt to transmit, assuming that no other oracle has yet
	// successfully transmitted a report.
	DeltaStage time.Duration
	// The maximum number of rounds during an epoch.
	RMax uint64
	// S is the transmission schedule. For example, S = [1,2,3] indicates that
	// in the first stage of transmission one oracle will attempt to transmit,
	// in the second stage two more will attempt to transmit (if in their view
	// the first stage didn't succeed), and in the third stage three more will
	// attempt to transmit (if in their view the first and second stage didn't
	// succeed).
	//
	// sum(S) should equal n.
	S []int
	// Identities (i.e. public keys) of the oracles participating in this
	// protocol instance.
	OracleIdentities []config.OracleIdentity

	// Binary blob containing configuration passed through to the
	// ReportingPlugin.
	ReportingPluginConfig []byte

	// MaxDurationX is the maximum duration a ReportingPlugin should spend
	// performing X. Reasonable values for these will be specific to each
	// ReportingPlugin. Be sure to not set these too short, or the corresponding
	// ReportingPlugin function may always time out.
	MaxDurationQuery                        time.Duration
	MaxDurationObservation                  time.Duration
	MaxDurationShouldAcceptAttestedReport   time.Duration
	MaxDurationShouldTransmitAcceptedReport time.Duration

	// The maximum number of oracles that are assumed to be faulty while the
	// protocol can retain liveness and safety. Unless you really know what
	// you’re doing, be sure to set this to floor((n-1)/3) where n is the total
	// number of oracles.
	F int

	// Binary blob containing configuration passed through to the
	// ReportingPlugin, and also available to the contract. (Unlike
	// ReportingPluginConfig which is only available offchain.)
	OnchainConfig []byte

	ConfigDigest types.ConfigDigest
}

// N is the number of oracles participating in the protocol
func (c *PublicConfig) N() int {
	return len(c.OracleIdentities)
}

func (c *PublicConfig) ByzQuorumSize() int {
	return byzquorum.Size(c.N(), c.F)
}

func (c *PublicConfig) MinRoundInterval() time.Duration {
	if c.DeltaRound > c.DeltaGrace {
		return c.DeltaRound
	}
	return c.DeltaGrace
}

func (c *PublicConfig) CheckParameterBounds() error {
	if c.F < 0 || c.F > math.MaxUint8 {
		return errors.Errorf("number of potentially faulty oracles must fit in 8 bits.")
	}
	return nil
}

func PublicConfigFromContractConfig(skipResourceExhaustionChecks bool, change types.ContractConfig) (PublicConfig, error) {
	pubcon, _, err := publicConfigFromContractConfig(skipResourceExhaustionChecks, change)
	return pubcon, err
}

func publicConfigFromContractConfig(skipResourceExhaustionChecks bool, change types.ContractConfig) (PublicConfig, config.SharedSecretEncryptions, error) {
	if change.OffchainConfigVersion != config.OCR3OffchainConfigVersion {
		return PublicConfig{}, config.SharedSecretEncryptions{}, fmt.Errorf("unsuppported OffchainConfigVersion %v, supported OffchainConfigVersion is %v", change.OffchainConfigVersion, config.OCR3OffchainConfigVersion)
	}

	oc, err := deserializeOffchainConfig(change.OffchainConfig)
	if err != nil {
		return PublicConfig{}, config.SharedSecretEncryptions{}, err
	}

	if err := checkIdentityListsHaveNoDuplicates(change, oc); err != nil {
		return PublicConfig{}, config.SharedSecretEncryptions{}, err
	}

	// must check that all lists have the same length, or bad input could crash
	// the following for loop.
	if err := checkIdentityListsHaveTheSameLength(change, oc); err != nil {
		return PublicConfig{}, config.SharedSecretEncryptions{}, err
	}

	identities := []config.OracleIdentity{}
	for i := range change.Signers {
		identities = append(identities, config.OracleIdentity{
			oc.OffchainPublicKeys[i],
			types.OnchainPublicKey(change.Signers[i][:]),
			oc.PeerIDs[i],
			change.Transmitters[i],
		})
	}

	cfg := PublicConfig{
		oc.DeltaProgress,
		oc.DeltaResend,
		oc.DeltaInitial,
		oc.DeltaRound,
		oc.DeltaGrace,
		oc.DeltaCertifiedCommitRequest,
		oc.DeltaStage,
		oc.RMax,
		oc.S,
		identities,
		oc.ReportingPluginConfig,
		oc.MaxDurationQuery,
		oc.MaxDurationObservation,
		oc.MaxDurationShouldAcceptAttestedReport,
		oc.MaxDurationShouldTransmitAcceptedReport,

		int(change.F),
		change.OnchainConfig,
		change.ConfigDigest,
	}

	if err := checkPublicConfigParameters(cfg); err != nil {
		return PublicConfig{}, config.SharedSecretEncryptions{}, err
	}

	if !skipResourceExhaustionChecks {
		if err := checkResourceExhaustion(cfg); err != nil {
			return PublicConfig{}, config.SharedSecretEncryptions{}, err
		}
	}

	return cfg, oc.SharedSecretEncryptions, nil
}

func checkIdentityListsHaveNoDuplicates(change types.ContractConfig, oc offchainConfig) error {
	// inefficient, but it doesn't matter
	for i := range change.Signers {
		for j := range change.Signers {
			if i != j && bytes.Equal(change.Signers[i], change.Signers[j]) {
				return fmt.Errorf("%v-th and %v-th signer are identical: %x", i, j, change.Signers[i])
			}
		}
	}

	{
		uniquePeerIDs := map[string]struct{}{}
		for _, peerID := range oc.PeerIDs {
			if _, ok := uniquePeerIDs[peerID]; ok {
				return fmt.Errorf("duplicate PeerID '%v'", peerID)
			}
			uniquePeerIDs[peerID] = struct{}{}
		}
	}

	{
		uniqueOffchainPublicKeys := map[types.OffchainPublicKey]struct{}{}
		for _, ocpk := range oc.OffchainPublicKeys {
			if _, ok := uniqueOffchainPublicKeys[ocpk]; ok {
				return fmt.Errorf("duplicate OffchainPublicKey %x", ocpk)
			}
			uniqueOffchainPublicKeys[ocpk] = struct{}{}
		}
	}

	{
		// this isn't strictly necessary, but since we don't intend to run
		// with duplicate transmitters at this time, we might as well check
		uniqueTransmitters := map[types.Account]struct{}{}
		for _, transmitter := range change.Transmitters {
			if _, ok := uniqueTransmitters[transmitter]; ok {
				return fmt.Errorf("duplicate transmitter '%v'", transmitter)
			}
			uniqueTransmitters[transmitter] = struct{}{}
		}
	}

	// no point in checking SharedSecretEncryptions for uniqueness

	return nil
}

func checkIdentityListsHaveTheSameLength(
	change types.ContractConfig, oc offchainConfig,
) error {
	expectedLength := len(change.Signers)
	errorMsg := "%s list must have same length as onchain signers list: %d ≠ " +
		strconv.Itoa(expectedLength)
	for _, identityList := range []struct {
		length int
		name   string
	}{
		{len(oc.PeerIDs) /*                       */, "peer ids"},
		{len(oc.OffchainPublicKeys) /*            */, "offchain public keys"},
		{len(change.Transmitters) /*              */, "transmitters"},
		{len(oc.SharedSecretEncryptions.Encryptions), "shared-secret encryptions"},
	} {
		if identityList.length != expectedLength {
			return errors.Errorf(errorMsg, identityList.name, identityList.length)
		}
	}
	return nil
}

// Sanity check on parameters:
// (1) violations of fundamental constraints like 3*f<n;
// (2) configurations that would trivially exhaust all of a node's resources;
// (3) (some) simple mistakes

func checkPublicConfigParameters(cfg PublicConfig) error {
	/////////////////////////////////////////////////////////////////
	// Be sure to think about changes to other tooling that need to
	// be made when you change this function!
	/////////////////////////////////////////////////////////////////

	if !(0 <= cfg.F && cfg.F*3 < cfg.N()) {
		return fmt.Errorf("F (%v) must be non-negative and less than N/3 (N = %v)",
			cfg.F, cfg.N())
	}

	if !(cfg.N() <= types.MaxOracles) {
		return fmt.Errorf("N (%v) must be less than or equal MaxOracles (%v)",
			cfg.N(), types.MaxOracles)
	}

	if !(0 <= cfg.DeltaProgress) {
		return fmt.Errorf("DeltaProgress (%v) must be non-negative", cfg.DeltaProgress)
	}

	if !(0 <= cfg.DeltaResend) {
		return fmt.Errorf("DeltaResend (%v) must be non-negative", cfg.DeltaResend)
	}

	if !(0 <= cfg.DeltaInitial) {
		return fmt.Errorf("DeltaInitial (%v) must be non-negative", cfg.DeltaInitial)
	}

	if !(0 <= cfg.DeltaRound) {
		return fmt.Errorf("DeltaRound (%v) must be non-negative", cfg.DeltaRound)
	}

	if !(0 <= cfg.DeltaGrace) {
		return fmt.Errorf("DeltaGrace (%v) must be non-negative",
			cfg.DeltaGrace)
	}

	if !(0 <= cfg.DeltaCertifiedCommitRequest) {
		return fmt.Errorf("DeltaCertifiedCommitRequest (%v) must be non-negative", cfg.DeltaCertifiedCommitRequest)
	}

	if !(0 <= cfg.DeltaStage) {
		return fmt.Errorf("DeltaStage (%v) must be non-negative", cfg.DeltaStage)
	}

	if !(0 <= cfg.MaxDurationQuery) {
		return fmt.Errorf("MaxDurationQuery (%v) must be non-negative", cfg.MaxDurationQuery)
	}

	if !(0 <= cfg.MaxDurationObservation) {
		return fmt.Errorf("MaxDurationObservation (%v) must be non-negative", cfg.MaxDurationObservation)
	}

	if !(0 <= cfg.MaxDurationShouldAcceptAttestedReport) {
		return fmt.Errorf("MaxDurationShouldAcceptAttestedReport (%v) must be non-negative", cfg.MaxDurationShouldAcceptAttestedReport)
	}

	if !(0 <= cfg.MaxDurationShouldTransmitAcceptedReport) {
		return fmt.Errorf("MaxDurationShouldTransmitAcceptedReport (%v) must be non-negative", cfg.MaxDurationShouldTransmitAcceptedReport)
	}

	if !(cfg.DeltaRound < cfg.DeltaProgress) {
		return fmt.Errorf("DeltaRound (%v) must be less than DeltaProgress (%v)",
			cfg.DeltaRound, cfg.DeltaProgress)
	}

	sumMaxDurationsOutcomeGeneration := cfg.MaxDurationQuery + cfg.MaxDurationObservation + cfg.DeltaGrace
	if !(sumMaxDurationsOutcomeGeneration < cfg.DeltaProgress) {
		return fmt.Errorf("sum of MaxDurationQuery/MaxDurationObservation/DeltaGrace (%v) must be less than DeltaProgress (%v)",
			sumMaxDurationsOutcomeGeneration, cfg.DeltaProgress)
	}

	// We cannot easily add a similar check for the MaxDuration variables used
	// in the transmission protocol (MaxDurationShouldAcceptAttestedReport,
	// MaxDurationShouldTransmitAcceptedReport), because we don't know how often
	// they will be triggered. But if we assume that there is one transmission
	// for each round, we should have MaxDurationShouldAcceptAttestedReport +
	// MaxDurationShouldTransmitAcceptedReport < round duration.

	if !(0 < cfg.RMax) {
		return fmt.Errorf("RMax (%v) must be greater than zero", cfg.RMax)
	}

	// This prevents possible overflows adding up the elements of S. We should never
	// hit this.
	if !(len(cfg.S) < 1000) {
		return fmt.Errorf("len(S) (%v) must be less than 1000", len(cfg.S))
	}

	for i, s := range cfg.S {
		if !(0 <= s && s <= types.MaxOracles) {
			return fmt.Errorf("S[%v] (%v) must be between 0 and types.MaxOracles (%v)", i, s, types.MaxOracles)
		}
	}

	return nil
}

func checkResourceExhaustion(cfg PublicConfig) error {
	// Sending messages related to epoch changes and missing certified commits
	// shouldn't be necessary in any realistic WAN deployment and could cause
	// resource exhaustion
	const safeInterval = 100 * time.Millisecond
	if cfg.DeltaProgress < safeInterval {
		return fmt.Errorf("DeltaProgress (%v) is set below the resource exhaustion safe interval (%v)", cfg.DeltaProgress, safeInterval)
	}
	if cfg.DeltaResend < safeInterval {
		return fmt.Errorf("DeltaResend (%v) is set below the resource exhaustion safe interval (%v)", cfg.DeltaResend, safeInterval)
	}
	if cfg.DeltaInitial < safeInterval {
		return fmt.Errorf("DeltaInitial (%v) is set below the resource exhaustion safe interval (%v)", cfg.DeltaInitial, safeInterval)
	}
	if cfg.DeltaCertifiedCommitRequest < safeInterval {
		return fmt.Errorf("DeltaCertifiedCommitRequest (%v) is set below the resource exhaustion safe interval (%v)", cfg.DeltaCertifiedCommitRequest, safeInterval)
	}
	// We don't check DeltaGrace, DeltaRound, DeltaStage since none of them
	// would exhaust the oracle's resources even if they are all set to 0.
	return nil
}
