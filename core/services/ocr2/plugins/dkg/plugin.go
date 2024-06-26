package dkg

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-vrf/altbn_128"
	"github.com/smartcontractkit/chainlink-vrf/dkg"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/dkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/dkg/persistence"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func NewDKGServices(
	jb job.Job,
	ocr2Provider evmrelay.DKGProvider,
	lggr logger.Logger,
	ocrLogger commontypes.Logger,
	dkgSignKs keystore.DKGSign,
	dkgEncryptKs keystore.DKGEncrypt,
	ethClient evmclient.Client,
	oracleArgsNoPlugin libocr2.OCR2OracleArgs,
	ds sqlutil.DataSource,
	chainID *big.Int,
	network string,
) ([]job.ServiceCtx, error) {
	var pluginConfig config.PluginConfig
	err := json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal plugin config")
	}
	err = config.ValidatePluginConfig(pluginConfig, dkgSignKs, dkgEncryptKs)
	if err != nil {
		return nil, errors.Wrap(err, "validate plugin config")
	}
	signKey, err := dkgSignKs.Get(pluginConfig.SigningPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "get dkgsign key")
	}
	encryptKey, err := dkgEncryptKs.Get(pluginConfig.EncryptionPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "get dkgencrypt key")
	}
	onchainDKGClient, err := NewOnchainDKGClient(
		jb.OCR2OracleSpec.ContractID,
		ethClient)
	if err != nil {
		return nil, errors.Wrap(err, "new onchain dkg client")
	}
	onchainContract := dkg.NewOnchainContract(onchainDKGClient, &altbn_128.G2{})
	keyConsumer := newDummyKeyConsumer()
	keyID, err := DecodeKeyID(pluginConfig.KeyID)
	if err != nil {
		return nil, errors.Wrap(err, "decode key ID")
	}
	shareDB := persistence.NewShareDB(ds, lggr.Named("DKGShareDB"), chainID, network)
	oracleArgsNoPlugin.ReportingPluginFactory = dkg.NewReportingPluginFactory(
		encryptKey.KyberScalar(),
		signKey.KyberScalar(),
		keyID,
		onchainContract,
		ocrLogger,
		keyConsumer,
		shareDB,
	)
	oracle, err := libocr2.NewOracle(oracleArgsNoPlugin)
	if err != nil {
		return nil, errors.Wrap(err, "error calling NewOracle")
	}
	return []job.ServiceCtx{ocr2Provider, job.NewServiceAdapter(oracle)}, nil
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
