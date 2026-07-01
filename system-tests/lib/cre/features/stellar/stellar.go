package stellar

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// Stellar is the Local CRE feature for the Stellar chain capability.
//
// SCAFFOLD (read milestone): this is a registered no-op stub so DONs that do
// NOT request the Stellar capability are unaffected. The read-path body should
// mirror system-tests/lib/cre/features/solana/v2/solana.go (and features/aptos/*),
// implementing:
//   - PreEnvStartup: detect enabled Stellar chain IDs on the DON, build the
//     capability registrations (LabelledName "stellar:ChainSelector:<sel>",
//     version 1.0.0) and per-label OCR3 configs. No forwarder for reads.
//   - PostEnvStartup: propose + approve the Stellar worker job specs.
//
// See STELLAR_LOCAL_CRE_PLAN.md items A11–A16. The write path (forwarder
// deploy/config, WriteReport, secp256k1/EVM-encoder reports) is Milestone B.
type Stellar struct{}

var _ cre.Feature = (*Stellar)(nil)

const (
	CapabilityVersion     = "1.0.0"
	CapabilityLabelPrefix = "stellar:ChainSelector:"
)

func (s *Stellar) Flag() cre.CapabilityFlag { return cre.StellarCapability }

func (s *Stellar) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	// TODO(A11): build capability registrations + per-label OCR3 configs for any
	// enabled Stellar chains on this DON. No-op until implemented.
	return &cre.PreEnvStartupOutput{}, nil
}

func (s *Stellar) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	// TODO(A11): propose Stellar worker job specs and approve them on the DON.
	return nil
}
