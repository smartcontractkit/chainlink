package config

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

// PublicConfig is the configuration disseminated through the smart contract
// It's public, because anybody can read it from the blockchain
type PublicConfig struct {
	DeltaProgress    time.Duration
	DeltaResend      time.Duration
	DeltaRound       time.Duration
	DeltaGrace       time.Duration
	DeltaC           time.Duration
	AlphaPPB         uint64
	DeltaStage       time.Duration
	RMax             uint8
	S                []int
	OracleIdentities []OracleIdentity

	F            int
	ConfigDigest types.ConfigDigest
}

type OracleIdentity struct {
	PeerID                string
	OffchainPublicKey     types.OffchainPublicKey
	OnChainSigningAddress types.OnChainSigningAddress
	TransmitAddress       common.Address
}

// N is the number of oracles participating in the protocol
func (c *PublicConfig) N() int {
	return len(c.OracleIdentities)
}

func (c *PublicConfig) CheckParameterBounds() error {
	if c.F < 0 || c.F > math.MaxUint8 {
		return errors.Errorf("number of potentially faulty oracles must fit in 8 bits.")
	}
	return nil
}

func PublicConfigFromContractConfig(chainID *big.Int, skipChainSpecificChecks bool, change types.ContractConfig) (PublicConfig, error) {
	pubcon, _, err := publicConfigFromContractConfig(chainID, skipChainSpecificChecks, change)
	return pubcon, err
}

func publicConfigFromContractConfig(chainID *big.Int, skipChainSpecificChecks bool, change types.ContractConfig) (PublicConfig, SharedSecretEncryptions, error) {
	oc, err := decodeContractSetConfigEncodedComponents(change.Encoded)
	if err != nil {
		return PublicConfig{}, SharedSecretEncryptions{}, err
	}

	// must check that all lists have the same length, or bad input could crash
	// the following for loop.
	if err := checkIdentityListsHaveTheSameLength(change, oc); err != nil {
		return PublicConfig{}, SharedSecretEncryptions{}, err
	}

	identities := []OracleIdentity{}
	for i := range change.Signers {
		identities = append(identities, OracleIdentity{
			oc.PeerIDs[i],
			oc.OffchainPublicKeys[i],
			types.OnChainSigningAddress(change.Signers[i]),
			change.Transmitters[i],
		})
	}

	cfg := PublicConfig{
		oc.DeltaProgress,
		oc.DeltaResend,
		oc.DeltaRound,
		oc.DeltaGrace,
		oc.DeltaC,
		oc.AlphaPPB,
		oc.DeltaStage,
		oc.RMax,
		oc.S,
		identities,
		int(change.Threshold),
		change.ConfigDigest,
	}

	if err := checkPublicConfigParameters(cfg); err != nil {
		return PublicConfig{}, SharedSecretEncryptions{}, err
	}

	if !skipChainSpecificChecks {
		if err := checkPublicConfigParametersForChain(chainID, cfg); err != nil {
			return PublicConfig{}, SharedSecretEncryptions{}, err
		}
	}

	return cfg, oc.SharedSecretEncryptions, nil
}

func checkIdentityListsHaveTheSameLength(
	change types.ContractConfig, oc setConfigEncodedComponents,
) error {
	expectedLength := len(change.Signers)
	errorMsg := "%s list must have same length as onchain signers list: %d ≠ " +
		strconv.Itoa(expectedLength)
	for _, identityList := range []struct {
		length int
		name   string
	}{
		{len(oc.PeerIDs) /*                       */, "peer ID"},
		{len(oc.OffchainPublicKeys) /*            */, "offchain public keys"},
		{len(change.Transmitters) /*              */, "transmitter address"},
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

	if !(0 <= cfg.DeltaC) {
		return fmt.Errorf("DeltaC (%v) must be non-negative", cfg.DeltaC)
	}

	if !(0 <= cfg.DeltaStage) {
		return fmt.Errorf("DeltaStage (%v) must be non-negative", cfg.DeltaStage)
	}

	if !(0 <= cfg.DeltaRound) {
		return fmt.Errorf("DeltaRound (%v) must be non-negative", cfg.DeltaRound)
	}

	if !(0 <= cfg.DeltaProgress) {
		return fmt.Errorf("DeltaProgress (%v) must be non-negative", cfg.DeltaProgress)
	}

	if !(0 <= cfg.DeltaResend) {
		return fmt.Errorf("DeltaResend (%v) must be non-negative", cfg.DeltaResend)
	}

	if !(0 <= cfg.F && cfg.F*3 < cfg.N()) {
		return fmt.Errorf("F (%v) must be non-negative and less than N/3 (N = %v)",
			cfg.F, cfg.N())
	}

	if !(cfg.N() <= types.MaxOracles) {
		return fmt.Errorf("N (%v) must be less than or equal MaxOracles (%v)",
			cfg.N(), types.MaxOracles)
	}

	if !(0 <= cfg.DeltaGrace) {
		return fmt.Errorf("DeltaGrace (%v) must be non-negative",
			cfg.DeltaGrace)
	}

	if !(cfg.DeltaGrace < cfg.DeltaRound) {
		return fmt.Errorf("DeltaGrace (%v) must be less than DeltaRound (%v)",
			cfg.DeltaGrace, cfg.DeltaRound)
	}

	if !(cfg.DeltaRound < cfg.DeltaProgress) {
		return fmt.Errorf("DeltaRound (%v) must be less than DeltaProgress (%v)",
			cfg.DeltaRound, cfg.DeltaProgress)
	}

	// *less* than 255 is intentional!
	// In report_generation_leader.go, we add 1 to a round number that can equal RMax.
	if !(0 < cfg.RMax && cfg.RMax < 255) {
		return fmt.Errorf("RMax (%v) must be greater than zero and less than 255", cfg.RMax)
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

func checkPublicConfigParametersForChain(chainID *big.Int, cfg PublicConfig) error {
	/////////////////////////////////////////////////////////////////
	// Be sure to think about changes to other tooling that need to
	// be made when you change this function!
	/////////////////////////////////////////////////////////////////

	type chainType int
	const (
		_ chainType = iota
		chainTypeSlowUpdates
		chainTypeModerateUpdates
		chainTypeFastUpdates
		chainTypePublicTestnet
		chainTypePrivateTestnet
	)

	type chainInfo struct {
		Name      string
		ChainType chainType
	}

	type chainLimits struct {
		MinDeltaC        time.Duration
		MinDeltaStage    time.Duration
		MinDeltaRound    time.Duration
		MinDeltaProgress time.Duration
		MinDeltaResend   time.Duration
	}

	if chainID == nil {
		return fmt.Errorf("chainID is nil, cannot perform chain-specific checks")
	}

	info, ok := map[uint64]chainInfo{
		1337:            {"SimulatedBackend", chainTypePrivateTestnet},
		7418:            {"Geth Local Testnet", chainTypePrivateTestnet},
		42161:           {"Arbitrum One", chainTypeFastUpdates},
		42170:           {"Arbitrum Nova", chainTypeFastUpdates},
		144545313136048: {"Arbitrum Testnet Kovan", chainTypePublicTestnet},
		421611:          {"Arbitrum Testnet Rinkeby", chainTypePublicTestnet},
		421613:          {"Arbitrum Testnet Goerli", chainTypePublicTestnet},
		421614:          {"Arbitrum Testnet Sepolia", chainTypePublicTestnet},
		43114:           {"Avalanche", chainTypeFastUpdates},
		43113:           {"Avalanche Testnet Fuji", chainTypePublicTestnet},
		8453:            {"Base", chainTypeFastUpdates},
		84531:           {"Base Testnet Goerli", chainTypePublicTestnet},
		84532:           {"Base Testnet Sepolia", chainTypePublicTestnet},
		56:              {"BSC", chainTypeFastUpdates},
		97:              {"BSC Testnet", chainTypePublicTestnet},
		42220:           {"Celo", chainTypeModerateUpdates},
		44787:           {"Celo Testnet", chainTypePublicTestnet},
		65:              {"Cosmos Testnet Okex", chainTypePublicTestnet},
		128:             {"HECO", chainTypeModerateUpdates},
		256:             {"HECO Testnet", chainTypePublicTestnet},
		1:               {"Ethereum", chainTypeSlowUpdates},
		5:               {"Ethereum Testnet Goerli", chainTypePublicTestnet},
		42:              {"Ethereum Testnet Kovan", chainTypePublicTestnet},
		4:               {"Ethereum Testnet Rinkeby", chainTypePublicTestnet},
		11155111:        {"Ethereum Testnet Sepolia", chainTypePublicTestnet},
		250:             {"Fantom", chainTypeFastUpdates},
		4002:            {"Fantom Testnet", chainTypePublicTestnet},
		1666600000:      {"Harmony Shard 0", chainTypeModerateUpdates},
		1666600001:      {"Harmony Shard 1", chainTypeModerateUpdates},
		1666600002:      {"Harmony Shard 2", chainTypeModerateUpdates},
		1666600003:      {"Harmony Shard 3", chainTypeModerateUpdates},
		1666700000:      {"Harmony Testnet Shard 0", chainTypePublicTestnet},
		1666700001:      {"Harmony Testnet Shard 1", chainTypePublicTestnet},
		1666700002:      {"Harmony Testnet Shard 2", chainTypePublicTestnet},
		1666700003:      {"Harmony Testnet Shard 3", chainTypePublicTestnet},
		8217:            {"Klaytn", chainTypeFastUpdates},
		1001:            {"Klaytn Testnet Baobab", chainTypePublicTestnet},
		59144:           {"Linea", chainTypeFastUpdates},
		59140:           {"Linea Testnet Goerli", chainTypePublicTestnet},
		59141:           {"Linea Testnet Sepolia", chainTypePublicTestnet},
		5000:            {"Mantle", chainTypeModerateUpdates},
		5001:            {"Mantle Testnet", chainTypePublicTestnet},
		1088:            {"Metis", chainTypeFastUpdates},
		588:             {"Metis Testnet Rinkeby", chainTypePublicTestnet},
		1284:            {"Moonbeam", chainTypeFastUpdates},
		1285:            {"Moonriver", chainTypeFastUpdates},
		1287:            {"Moonbeam Testnet Moonbase Alpha", chainTypePublicTestnet},
		10:              {"Optimism", chainTypeFastUpdates},
		69:              {"Optimism Testnet Kovan", chainTypePublicTestnet},
		420:             {"Optimism Testnet Goerli", chainTypePublicTestnet},
		11155420:        {"Optimism Testnet Sepolia", chainTypePublicTestnet},
		137:             {"Polygon", chainTypeFastUpdates},
		80002:           {"Polygon Testnet Amoy", chainTypePublicTestnet},
		80001:           {"Polygon Testnet Mumbai", chainTypePublicTestnet},
		1101:            {"Polygon zkEVM", chainTypeSlowUpdates},
		1442:            {"Polygon zkEVM Testnet", chainTypePublicTestnet},
		2442:            {"Polygon zkEVM Testnet Cardona", chainTypePublicTestnet},
		30:              {"RSK", chainTypeModerateUpdates},
		31:              {"RSK Testnet", chainTypePublicTestnet},
		534352:          {"Scroll Mainnet", chainTypeFastUpdates},
		534351:          {"Scroll Testnet Sepolia", chainTypePublicTestnet},
		196:             {"X Layer Mainnet", chainTypeFastUpdates},
		195:             {"X Layer Testnet Sepolia", chainTypePublicTestnet},
		100:             {"xDai", chainTypeModerateUpdates},
		10200:           {"xDai Testnet Chiado", chainTypePublicTestnet},
		324:             {"zkSync Mainnet", chainTypeFastUpdates},
		280:             {"zkSync Testnet Goerli", chainTypePublicTestnet},
		300:             {"zkSync Testnet Sepolia", chainTypePublicTestnet},
	}[chainID.Uint64()]
	if !ok {
		// "fail-closed" design. If we don't know the chain, we assume that
		// we shouldn't be updating it quickly
		info = chainInfo{"UNKNOWN", chainTypeSlowUpdates}
	}

	limits, ok := map[chainType]chainLimits{
		chainTypeSlowUpdates: {
			10 * time.Minute,
			10 * time.Second,
			20 * time.Second,
			23 * time.Second,
			10 * time.Second,
		},
		chainTypeModerateUpdates: {
			1 * time.Minute,
			5 * time.Second,
			20 * time.Second,
			23 * time.Second,
			10 * time.Second,
		},
		chainTypeFastUpdates: {
			10 * time.Second,
			5 * time.Second,
			5 * time.Second,
			8 * time.Second,
			5 * time.Second,
		},
		chainTypePublicTestnet: {
			1 * time.Second,
			5 * time.Second,
			1 * time.Second,
			2 * time.Second,
			2 * time.Second,
		},
		chainTypePrivateTestnet: {}, // do whatever you want on private testnet
	}[info.ChainType]
	if !ok {
		return fmt.Errorf("unknown chainType (%v) for chainID %v, cannot check config parameters", info.ChainType, chainID)
	}

	if !(limits.MinDeltaC <= cfg.DeltaC) {
		return fmt.Errorf("DeltaC (%v) must be greater or equal %v on chain %v (chainID: %v)",
			cfg.DeltaC, limits.MinDeltaC, info.Name, chainID)
	}

	if !(limits.MinDeltaStage <= cfg.DeltaStage) {
		return fmt.Errorf("DeltaStage (%v) must be greater or equal %v on chain %v (chainID: %v)",
			cfg.DeltaStage, limits.MinDeltaStage, info.Name, chainID)
	}

	if !(limits.MinDeltaRound <= cfg.DeltaRound) {
		return fmt.Errorf("DeltaRound (%v) must be greater or equal %v on chain %v (chainID: %v)",
			cfg.DeltaRound, limits.MinDeltaRound, info.Name, chainID)
	}

	if !(limits.MinDeltaProgress <= cfg.DeltaProgress) {
		return fmt.Errorf("DeltaProgress (%v) must be greater or equal %v on chain %v (chainID: %v)",
			cfg.DeltaProgress, limits.MinDeltaProgress, info.Name, chainID)
	}

	if !(limits.MinDeltaResend <= cfg.DeltaResend) {
		return fmt.Errorf("DeltaResend (%v) must be greater or equal %v on chain %v (chainID: %v)",
			cfg.DeltaResend, limits.MinDeltaResend, info.Name, chainID)
	}

	return nil
}
