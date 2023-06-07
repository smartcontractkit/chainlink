package threshold

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"

	decryptionPlugin "github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type ThresholdServicesConfig struct {
	DecryptionQueue  decryptionPlugin.DecryptionQueuingService
	PublicKey        []byte //*tdh2easy.PublicKey
	PrivKeyShare     []byte //*tdh2easy.PrivateShare
	OracleToKeyShare map[commontypes.OracleID]int
}

func NewThresholdServices(sharedOracleArgs *libocr2.OracleArgs, conf *ThresholdServicesConfig) (job.ServiceCtx, error) {
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

	sharedOracleArgs.ReportingPluginFactory = decryptionPlugin.DecryptionReportingPluginFactory{
		DecryptionQueue:  conf.DecryptionQueue,
		PublicKey:        publicKey,
		PrivKeyShare:     privKeyShare,
		OracleToKeyShare: conf.OracleToKeyShare,
		Logger:           sharedOracleArgs.Logger,
	}

	thresholdReportingPluginOracle, err := libocr2.NewOracle(*sharedOracleArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call NewOracle to create a Threshold Reporting Plugin")
	}

	return job.NewServiceAdapter(thresholdReportingPluginOracle), nil
}
