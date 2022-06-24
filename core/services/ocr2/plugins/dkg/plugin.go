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
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/dkg/config"
	evmrelay "github.com/smartcontractkit/chainlink/core/services/relay/evm"
)

type DKGContainer struct {
	jb           job.Job
	logger       logger.Logger
	ocrLogger    commontypes.Logger
	pluginConfig config.PluginConfig
	ocr2Provider evmrelay.DKGProvider
	dkgSignKs    keystore.DKGSign
	dkgEncryptKs keystore.DKGEncrypt
	ethClient    evmclient.Client
}

var _ plugins.OraclePlugin = &DKGContainer{}

func NewDKG(
	jb job.Job,
	ocr2Provider evmrelay.DKGProvider,
	l logger.Logger,
	ocrLogger commontypes.Logger,
	dkgSignKs keystore.DKGSign,
	dkgEncryptKs keystore.DKGEncrypt,
	ethClient evmclient.Client) (*DKGContainer, error) {
	var pluginConfig config.PluginConfig
	err := json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return &DKGContainer{}, err
	}
	err = config.ValidatePluginConfig(pluginConfig)
	if err != nil {
		return &DKGContainer{}, err
	}

	return &DKGContainer{
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

func (d *DKGContainer) GetPluginFactory() (ocr2types.ReportingPluginFactory, error) {
	signKey, err := d.dkgSignKs.Get(d.pluginConfig.SigningPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "get dkgsign key")
	}
	encryptKey, err := d.dkgEncryptKs.Get(d.pluginConfig.EncryptionPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "get dkgencrypt key")
	}
	onchainDKGClient, err := newOnchainDKGClient(
		d.jb.OCR2OracleSpec.ContractID,
		d.ethClient)
	if err != nil {
		return nil, errors.Wrap(err, "ew onchain dkg client")
	}
	onchainContract := dkg.NewOnchainContract(onchainDKGClient, &altbn_128.G1{})
	keyConsumer := newDummyKeyConsumer()
	factory := dkg.NewReportingPluginFactory(
		encryptKey.KyberScalar(),
		signKey.KyberScalar(),
		dkg.KeyID(decodeKeyID(d.pluginConfig.KeyID)),
		onchainContract,
		d.ocrLogger,
		keyConsumer,
	)
	return factory, nil
}

func (d *DKGContainer) GetServices() ([]job.ServiceCtx, error) {
	return []job.ServiceCtx{}, nil
}

func decodeKeyID(val string) (byteArray [32]byte) {
	decoded, err := hex.DecodeString(val)
	helpers.PanicErr(err)
	if len(decoded) != 32 {
		panic(fmt.Sprintf("expected value to be 32 bytes but received %d bytes", len(decoded)))
	}
	copy(byteArray[:], decoded)
	return
}
