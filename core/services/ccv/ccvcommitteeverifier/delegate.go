package ccvcommitteeverifier

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml"
	"github.com/smartcontractkit/chainlink-ccv/integration/pkg/constructors"
	"github.com/smartcontractkit/chainlink-ccv/verifier"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ccv/ccvcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

type Delegate struct {
	lggr logger.Logger
	// Houses secrets that are needed by the verifier (e.g. aggregator API keys).
	ccvConfig config.CCV
	// TODO: EVM specific (!)
	chainServices []commontypes.ChainService
	ocrKs         keystore.OCR2

	isNewlyCreatedJob bool
}

func NewDelegate(lggr logger.Logger, ccvConfig config.CCV, ocrKs keystore.OCR2, chainServices []commontypes.ChainService) *Delegate {
	return &Delegate{
		lggr:          lggr.Named("CCVCommitteeVerifierDelegate"),
		ccvConfig:     ccvConfig,
		chainServices: chainServices,
		ocrKs:         ocrKs,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.CCVCommitteeVerifier
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {
	d.isNewlyCreatedJob = true
}

func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) (services []job.ServiceCtx, err error) {
	var decodedCfg verifier.Config
	err = toml.Unmarshal([]byte(spec.CCVCommitteeVerifierSpec.CommitteeVerifierConfig), &decodedCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal committeeVerifierConfig into the verifier config struct: %w", err)
	}

	legacyChains := ccvcommon.GetLegacyChains(d.lggr, d.chainServices)

	// TODO: the Get() here expects a OCR key ID, which is not the same as the hex address
	// of the onchain signing key.
	// kb, err := d.ks.Get(decodedCfg.SignerAddress)
	evmKbs, err := d.ocrKs.GetAllOfType(chaintype.EVM)
	if err != nil {
		return nil, fmt.Errorf("failed to get OCR2 key bundles for EVM chains: %w", err)
	}

	if len(evmKbs) == 0 {
		return nil, fmt.Errorf("no OCR2 key bundles found for EVM chains")
	}

	// there should also usually be just one key bundle per chain?
	if len(evmKbs) > 1 {
		return nil, fmt.Errorf("multiple OCR2 key bundles found for EVM chains")
	}

	kb := evmKbs[0]
	onChainPubKeyBytes := common.HexToAddress(kb.OnChainPublicKey()).Bytes()
	configPubKeyBytes := common.HexToAddress(decodedCfg.SignerAddress).Bytes()
	if !bytes.Equal(onChainPubKeyBytes, configPubKeyBytes) {
		return nil, fmt.Errorf("onchain public key in the node's OCR2 key bundle (%s) does not match the signer address in the config (%s)", kb.OnChainPublicKey(), decodedCfg.SignerAddress)
	}

	d.lggr.Infow("using OCR2 key bundle for signing", "id", kb.ID(), "pubKey", kb.OnChainPublicKey())

	apiKey, apiSecret, err := getAggregatorSecrets(d.ccvConfig, decodedCfg.VerifierID)
	if err != nil {
		return nil, fmt.Errorf("failed to get aggregator secrets: %w", err)
	}

	// TODO: pass secrets as a separate param in the constructor.
	vc, err := constructors.NewVerificationCoordinator(
		d.lggr.Named("CCVCommitteeVerificationCoordinator"),
		decodedCfg,
		constructors.AggregatorSecret{
			APIKey:    apiKey,
			SecretKey: apiSecret,
		},
		common.HexToAddress(decodedCfg.SignerAddress).Bytes(),
		newSigner(kb),
		legacyChains,
	)
	if err != nil {
		d.lggr.Errorw("failed to create verification coordinator", "error", err)
		return nil, fmt.Errorf("failed to create verification coordinator: %w", err)
	}

	services = append(services, vc)

	return services, nil
}

func getAggregatorSecrets(ccvConfig config.CCV, verifierID string) (string, string, error) {
	for _, secret := range ccvConfig.AggregatorSecrets() {
		if secret.CommitteeID() == verifierID {
			return secret.APIKey(), secret.APISecret(), nil
		}
	}
	return "", "", fmt.Errorf("no aggregator secrets found for verifier ID: %s", verifierID)
}

func (d *Delegate) AfterJobCreated(spec job.Job) {}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) OnDeleteJob(ctx context.Context, spec job.Job) error {
	// TODO: shut down needed services?
	return nil
}

// signer adapts SignBlob to Sign so that the verifier can use it.
type signer struct {
	kb ocr2key.KeyBundle
}

func newSigner(kb ocr2key.KeyBundle) *signer {
	return &signer{kb}
}

func (s *signer) Sign(msg []byte) ([]byte, error) {
	return s.kb.SignBlob(msg)
}
