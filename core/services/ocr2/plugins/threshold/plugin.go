package threshold

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	decryptionPlugin "github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type ThresholdServicesConfig struct {
	DecryptionQueue decryptionPlugin.DecryptionQueuingService
	PublicKey       []byte //*tdh2easy.PublicKey
	PrivKeyShare    []byte //*tdh2easy.PrivateShare
	ConfigParser    decryptionPlugin.ConfigParser
}

func NewThresholdService(sharedOracleArgs *libocr2.OCR2OracleArgs, conf *ThresholdServicesConfig) (job.ServiceCtx, error) {
	var publicKey *tdh2easy.PublicKey
	err := publicKey.Unmarshal(conf.PublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal threshold encryption public key")
	}

	var privKeyShare *tdh2easy.PrivateShare
	err = privKeyShare.Unmarshal(conf.PrivKeyShare)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal threshold decryption private key share")
	}

	// It can be safely assumed all oracle nodes have been given the correct key share,
	// so we use the key share index as the oracle's index.
	// If an oracle was given the wrong key share, key share decryption would have failed.
	oracleToKeyShare := make(map[commontypes.OracleID]int)
	oracleToKeyShare[commontypes.OracleID(privKeyShare.Index())] = privKeyShare.Index()

	sharedOracleArgs.ReportingPluginFactory = decryptionPlugin.DecryptionReportingPluginFactory{
		DecryptionQueue:  conf.DecryptionQueue,
		ConfigParser:     conf.ConfigParser,
		PublicKey:        publicKey,
		PrivKeyShare:     privKeyShare,
		OracleToKeyShare: oracleToKeyShare,
		Logger:           sharedOracleArgs.Logger,
	}

	thresholdReportingPluginOracle, err := libocr2.NewOracle(*sharedOracleArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call NewOracle to create a Threshold Reporting Plugin")
	}

	return job.NewServiceAdapter(thresholdReportingPluginOracle), nil
}
