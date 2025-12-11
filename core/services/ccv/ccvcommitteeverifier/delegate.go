package ccvcommitteeverifier

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-ccv/integration/pkg/constructors"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-ccv/protocol/common/hmac"
	"github.com/smartcontractkit/chainlink-ccv/verifier/commit"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
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
	delegateLogger logger.Logger
	// logger used to create new Named loggers for the services
	lggr logger.Logger
	// Houses secrets that are needed by the verifier (e.g. aggregator API keys).
	ccvConfig config.CCV
	// TODO: EVM specific (!)
	chainServices []commontypes.ChainService
	ocrKs         keystore.OCR2
	ds            sqlutil.DataSource

	isNewlyCreatedJob bool
}

func NewDelegate(lggr logger.Logger, ds sqlutil.DataSource, ccvConfig config.CCV, ocrKs keystore.OCR2, chainServices []commontypes.ChainService) *Delegate {
	return &Delegate{
		delegateLogger: lggr.Named("CCVCommitteeVerifierDelegate"),
		lggr:           lggr,
		ds:             ds,
		ccvConfig:      ccvConfig,
		chainServices:  chainServices,
		ocrKs:          ocrKs,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.CCVCommitteeVerifier
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {
	d.isNewlyCreatedJob = true
}

func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) (services []job.ServiceCtx, err error) {
	d.delegateLogger.Infow("Creating services for CCV committee verifier job", "jobID", spec.ID)

	var decodedCfg commit.Config
	err = toml.Unmarshal([]byte(spec.CCVCommitteeVerifierSpec.CommitteeVerifierConfig), &decodedCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal committeeVerifierConfig into the verifier config struct: %w", err)
	}

	d.delegateLogger.Infow("validating committee verifier config", "config", decodedCfg, "raw", spec.CCVCommitteeVerifierSpec.CommitteeVerifierConfig)

	err = decodedCfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate committee verifier config: %w", err)
	}

	err = decodedCfg.Monitoring.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate committee verifier monitoring config: %w", err)
	}

	// Chains in the committee verifier configuration should dictate what we end up verifying for.
	var chainsInConfig = make([]protocol.ChainSelector, 0, len(decodedCfg.CommitteeVerifierAddresses))
	for chainSelStr := range decodedCfg.CommitteeVerifierAddresses {
		parsed, err2 := strconv.ParseUint(chainSelStr, 10, 64)
		if err2 != nil {
			return nil, fmt.Errorf("failed to parse chain selector string from committee verifier config (%s): %w", chainSelStr, err)
		}
		chainsInConfig = append(chainsInConfig, protocol.ChainSelector(parsed))
	}
	legacyChains, err := ccvcommon.GetLegacyChains(ctx, d.lggr, d.chainServices, chainsInConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get legacy chains: %w", err)
	}

	signingKeys, err := d.ocrKs.GetAllOfType(chaintype.EVM)
	if err != nil {
		return nil, fmt.Errorf("failed to get signing key %s from eth keystore: %w", decodedCfg.SignerAddress, err)
	}

	var signingKey ocr2key.KeyBundle
	switch len(signingKeys) {
	case 0:
		return nil, errors.New("no signing key found for EVM in OCR2 keystore")
	case 1:
		signingKey = signingKeys[0]
	default:
		d.delegateLogger.Warnw("multiple signing keys found for EVM, using the first", "keys", signingKeys)
		signingKey = signingKeys[0]
	}

	d.delegateLogger.Infow("using ocr2 onchain key for signing", "publicKey", signingKey.OnChainPublicKey())
	onchainPubKeyBytes, err := hex.DecodeString(signingKey.OnChainPublicKey())
	if err != nil {
		return nil, fmt.Errorf("failed to decode onchain public key: %w", err)
	}
	configPubKeyBytes, err := hexutil.Decode(decodedCfg.SignerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signer address: %w", err)
	}
	if !bytes.Equal(onchainPubKeyBytes, configPubKeyBytes) {
		return nil, fmt.Errorf("onchain public key does not match signer address in config, want %s, got %s", signingKey.OnChainPublicKey(), decodedCfg.SignerAddress)
	}

	apiKey, apiSecret, err := getAggregatorSecrets(d.ccvConfig, decodedCfg.VerifierID)
	if err != nil {
		return nil, fmt.Errorf("failed to get aggregator secrets from secrets toml: %w", err)
	}

	vc, err := constructors.NewVerificationCoordinator(
		d.lggr.
			Named("CCVCommitteeVerificationCoordinator").
			Named(decodedCfg.VerifierID),
		decodedCfg,
		&hmac.ClientConfig{
			APIKey: apiKey,
			Secret: apiSecret,
		},
		common.HexToAddress(decodedCfg.SignerAddress).Bytes(),
		newSignerAdapter(signingKey),
		legacyChains,
		d.ds,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create verification coordinator: %w", err)
	}

	services = append(services, vc)

	return services, nil
}

func getAggregatorSecrets(ccvConfig config.CCV, verifierID string) (string, string, error) {
	verifierIDs := make([]string, 0, len(ccvConfig.AggregatorSecrets()))
	for _, secret := range ccvConfig.AggregatorSecrets() {
		if secret.VerifierID() == verifierID {
			return secret.APIKey(), secret.APISecret(), nil
		}
		verifierIDs = append(verifierIDs, secret.VerifierID())
	}
	return "", "", fmt.Errorf("no aggregator secrets found for verifier ID %s, found %s", verifierID, strings.Join(verifierIDs, ", "))
}

func (d *Delegate) AfterJobCreated(spec job.Job) {}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) OnDeleteJob(ctx context.Context, spec job.Job) error {
	return nil
}

// signerAdapter is an adapter that implements the verifier.MessageSigner interface.
type signerAdapter struct {
	kb ocr2key.KeyBundle
}

func newSignerAdapter(kb ocr2key.KeyBundle) *signerAdapter { return &signerAdapter{kb} }

func (s *signerAdapter) Sign(input []byte) ([]byte, error) {
	return s.kb.SignBlob(input)
}
