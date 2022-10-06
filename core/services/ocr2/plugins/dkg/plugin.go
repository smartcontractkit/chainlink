package dkg

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	"github.com/smartcontractkit/ocr2vrf/dkg"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/dkg/config"
	evmrelay "github.com/smartcontractkit/chainlink/core/services/relay/evm"
)

type container struct {
	jb           job.Job
	logger       logger.Logger
	ocrLogger    commontypes.Logger
	pluginConfig config.PluginConfig
	ocr2Provider evmrelay.DKGProvider
	dkgSignKs    keystore.DKGSign
	dkgEncryptKs keystore.DKGEncrypt
	ethClient    evmclient.Client
}

func (d *container) Name() string {
	return "DKGPlugin"
}

var _ plugins.OraclePlugin = &container{}

func NewDKG(
	jb job.Job,
	ocr2Provider evmrelay.DKGProvider,
	l logger.Logger,
	ocrLogger commontypes.Logger,
	dkgSignKs keystore.DKGSign,
	dkgEncryptKs keystore.DKGEncrypt,
	ethClient evmclient.Client) (plugins.OraclePlugin, error) {
	var pluginConfig config.PluginConfig
	err := json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal plugin config")
	}
	err = config.ValidatePluginConfig(pluginConfig, dkgSignKs, dkgEncryptKs)
	if err != nil {
		return nil, errors.Wrap(err, "validate plugin config")
	}

	return &container{
		logger:       l,
		jb:           jb,
		ocrLogger:    ocrLogger,
		pluginConfig: pluginConfig,
		ocr2Provider: ocr2Provider,
		dkgSignKs:    dkgSignKs,
		dkgEncryptKs: dkgEncryptKs,
		ethClient:    ethClient,
	}, nil
}

func (d *container) GetPluginFactory() (ocr2types.ReportingPluginFactory, error) {
	signKey, err := d.dkgSignKs.Get(d.pluginConfig.SigningPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "get dkgsign key")
	}
	encryptKey, err := d.dkgEncryptKs.Get(d.pluginConfig.EncryptionPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "get dkgencrypt key")
	}
	onchainDKGClient, err := NewOnchainDKGClient(
		d.jb.OCR2OracleSpec.ContractID,
		d.ethClient)
	if err != nil {
		return nil, errors.Wrap(err, "new onchain dkg client")
	}
	onchainContract := dkg.NewOnchainContract(onchainDKGClient, &altbn_128.G2{})
	keyConsumer := newDummyKeyConsumer()
	keyID, err := DecodeKeyID(d.pluginConfig.KeyID)
	if err != nil {
		return nil, errors.Wrap(err, "decode key ID")
	}

	factory := dkg.NewReportingPluginFactory(
		encryptKey.KyberScalar(),
		signKey.KyberScalar(),
		keyID,
		onchainContract,
		d.ocrLogger,
		keyConsumer,
	)
	return factory, nil
}

func (d *container) GetServices() ([]job.ServiceCtx, error) {
	return []job.ServiceCtx{}, nil
}

func DecodeKeyID(val string) (byteArray [32]byte, err error) {
	decoded, err := hex.DecodeString(val)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "hex decode string")
	}
	if len(decoded) != 32 {
		return [32]byte{}, fmt.Errorf("expected value to be 32 bytes but received %d bytes", len(decoded))
	}
	copy(byteArray[:], decoded)
	return
}
