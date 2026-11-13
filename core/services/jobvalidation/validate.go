// Package jobvalidation turns a job spec in TOML into a job.Job, the way the
// node's HTTP API does. It is shared by the API and by the Dynamic Config
// reconciler so both accept exactly the same specs.
package jobvalidation

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	ccip "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockheaderfeeder"
	"github.com/smartcontractkit/chainlink/v2/core/services/ccv/ccvcommitteeverifier"
	"github.com/smartcontractkit/chainlink/v2/core/services/ccv/ccvexecutor"
	"github.com/smartcontractkit/chainlink/v2/core/services/cresettings"
	"github.com/smartcontractkit/chainlink/v2/core/services/cron"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

// ErrUnsupportedJobType is returned for a job type the node cannot run.
var ErrUnsupportedJobType = errors.New("unsupported job type")

// Config is the slice of the node's configuration the validators need.
type Config interface {
	OCR() coreconfig.OCR
	OCR2() coreconfig.OCR2
	Insecure() coreconfig.Insecure
}

// Validator validates job specs against the node's configuration.
type Validator struct {
	cfg       Config
	legacyEVM legacyevm.LegacyChainContainer
	registrar plugins.RegistrarConfig
}

// New returns a Validator. legacyEVM may be nil when no EVM chains are
// configured; only OCR v1 specs need it.
func New(cfg Config, legacyEVM legacyevm.LegacyChainContainer, registrar plugins.RegistrarConfig) *Validator {
	return &Validator{cfg: cfg, legacyEVM: legacyEVM, registrar: registrar}
}

// Validate parses and validates tomlString, returning the job to create.
func (v *Validator) Validate(ctx context.Context, tomlString string) (job.Job, error) {
	jobType, err := job.ValidateSpec(tomlString)
	if err != nil {
		return job.Job{}, fmt.Errorf("parse TOML: %w", err)
	}

	switch jobType {
	case job.OffchainReporting:
		if !v.cfg.OCR().Enabled() {
			return job.Job{}, fmt.Errorf("%w: the Offchain Reporting feature is disabled by configuration", ErrUnsupportedJobType)
		}

		return ocr.ValidatedOracleSpecToml(v.cfg, v.legacyEVM, tomlString) //nolint:staticcheck // LegacyEVMChains is deprecated but refactoring to the new relayer interface is out of scope here
	case job.OffchainReporting2:
		if !v.cfg.OCR2().Enabled() {
			return job.Job{}, fmt.Errorf("%w: the Offchain Reporting 2 feature is disabled by configuration", ErrUnsupportedJobType)
		}

		return validate.ValidatedOracleSpecToml(ctx, v.cfg.OCR2(), v.cfg.Insecure(), tomlString, v.registrar)
	case job.DirectRequest:
		return job.Job{}, fmt.Errorf("%w: job type directrequest has been removed and is no longer supported", ErrUnsupportedJobType)
	case job.FluxMonitor:
		return job.Job{}, fmt.Errorf("%w: job type fluxmonitor has been removed and is no longer supported", ErrUnsupportedJobType)
	case job.Webhook:
		return job.Job{}, fmt.Errorf("%w: job type webhook has been removed and is no longer supported", ErrUnsupportedJobType)
	case job.Workflow:
		return job.Job{}, fmt.Errorf("%w: job type workflow has been removed and is no longer supported", ErrUnsupportedJobType)
	case job.CRESettings:
		return cresettings.ValidatedCRESettingsSpec(tomlString)
	case job.Cron:
		return cron.ValidatedCronSpec(tomlString)
	case job.VRF:
		return vrfcommon.ValidatedVRFSpec(tomlString)
	case job.BlockhashStore:
		return blockhashstore.ValidatedSpec(tomlString)
	case job.BlockHeaderFeeder:
		return blockheaderfeeder.ValidatedSpec(tomlString)
	case job.Bootstrap:
		return ocrbootstrap.ValidatedBootstrapSpecToml(tomlString)
	case job.Gateway:
		return gateway.ValidatedGatewaySpec(tomlString)
	case job.Stream:
		return streams.ValidatedStreamSpec(tomlString)
	case job.StandardCapabilities:
		return standardcapabilities.ValidatedStandardCapabilitiesSpec(tomlString)
	case job.CCIP:
		return ccip.ValidatedCCIPSpec(tomlString)
	case job.CCVCommitteeVerifier:
		return ccvcommitteeverifier.ValidatedCCVCommitteeVerifierSpec(tomlString)
	case job.CCVExecutor:
		return ccvexecutor.ValidatedCCVExecutorSpec(tomlString)
	default:
		return job.Job{}, fmt.Errorf("%w: %s", ErrUnsupportedJobType, jobType)
	}
}
