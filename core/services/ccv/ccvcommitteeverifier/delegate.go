package ccvcommitteeverifier

import (
	"context"
	"fmt"
	"strconv"

	burntsushitoml "github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccv/integration/pkg/constructors"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-ccv/verifier"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ccv/ccvcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

type Delegate struct {
	lggr logger.Logger
	// Houses secrets that are needed by the verifier (e.g. aggregator API keys).
	ccvConfig config.CCV
	// TODO: EVM specific (!)
	chainServices []commontypes.ChainService
	// TODO: this is temporary, need to switch to the OCR2 keystore or another
	// custom one.
	ethKs keystore.Eth

	isNewlyCreatedJob bool
}

func NewDelegate(lggr logger.Logger, ccvConfig config.CCV, ocrKs keystore.Eth, chainServices []commontypes.ChainService) *Delegate {
	return &Delegate{
		lggr:          lggr.Named("CCVCommitteeVerifierDelegate"),
		ccvConfig:     ccvConfig,
		chainServices: chainServices,
		ethKs:         ocrKs,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.CCVCommitteeVerifier
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {
	d.isNewlyCreatedJob = true
}

func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) (services []job.ServiceCtx, err error) {
	d.lggr.Infow("Creating services for CCV committee verifier job", "jobID", spec.ID)

	// note that go-toml doesn't correctly parse nested TOMLs, at least from this struct,
	// so burntsushi/toml is needed.
	var decodedCfg verifier.Config
	_, err = burntsushitoml.Decode(spec.CCVCommitteeVerifierSpec.CommitteeVerifierConfig, &decodedCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal committeeVerifierConfig into the verifier config struct: %w", err)
	}

	d.lggr.Infow("validating committee verifier config", "config", decodedCfg, "raw", spec.CCVCommitteeVerifierSpec.CommitteeVerifierConfig)

	err = decodedCfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate committee verifier config: %w", err)
	}

	err = decodedCfg.Monitoring.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate committee verifier monitoring config: %w", err)
	}

	// Chains in the committee verifier configuration should dictate what we end up verifying for.
	var chainsInConfig []protocol.ChainSelector
	for chainSelStr := range decodedCfg.CommitteeVerifierAddresses {
		parsed, err := strconv.ParseUint(chainSelStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse chain selector string from committee verifier config (%s): %w", chainSelStr, err)
		}
		chainsInConfig = append(chainsInConfig, protocol.ChainSelector(parsed))
	}
	legacyChains, err := ccvcommon.GetLegacyChains(d.lggr, d.chainServices, chainsInConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get legacy chains: %w", err)
	}

	signingKey, err := d.ethKs.Get(ctx, decodedCfg.SignerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get signing key %s from eth keystore: %w", decodedCfg.SignerAddress, err)
	}

	d.lggr.Infow("using eth key for signing", "key", signingKey.Address.Hex())

	apiKey, apiSecret := getAggregatorSecrets(d.ccvConfig, decodedCfg.VerifierID)
	if apiKey == "" || apiSecret == "" {
		// fall back to the keys current set in the TOML config
		// TODO: this is a temporary solution to allow the node to run the verifier job but needs
		// to be fixed.
		apiKey = decodedCfg.AggregatorAPIKey
		apiSecret = decodedCfg.AggregatorSecretKey
		d.lggr.Warnw("no aggregator secrets found for verifier ID, using keys current set in the TOML config",
			"verifierID", decodedCfg.VerifierID)
	}

	vc, err := constructors.NewVerificationCoordinator(
		d.lggr.Named("CCVCommitteeVerificationCoordinator"),
		decodedCfg,
		constructors.AggregatorSecret{
			APIKey:    apiKey,
			SecretKey: apiSecret,
		},
		common.HexToAddress(decodedCfg.SignerAddress).Bytes(),
		signingKey,
		legacyChains,
	)
	if err != nil {
		d.lggr.Errorw("failed to create verification coordinator", "error", err)
		return nil, fmt.Errorf("failed to create verification coordinator: %w", err)
	}

	services = append(services, vc)

	return services, nil
}

func getAggregatorSecrets(ccvConfig config.CCV, verifierID string) (string, string) {
	for _, secret := range ccvConfig.AggregatorSecrets() {
		if secret.CommitteeID() == verifierID {
			return secret.APIKey(), secret.APISecret()
		}
	}
	return "", ""
}

func (d *Delegate) AfterJobCreated(spec job.Job) {}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) OnDeleteJob(ctx context.Context, spec job.Job) error {
	// TODO: shut down needed services?
	return nil
}
