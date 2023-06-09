package threshold

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	decryptionPlugin "github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type ThresholdServicesConfig struct {
	DecryptionQueue    decryptionPlugin.DecryptionQueuingService
	KeyshareWithPubKey []byte
	ConfigParser       decryptionPlugin.ConfigParser
}

type KeyshareWithPubKey struct {
	PublicKey       string `json:"PublicKey"`  //*tdh2easy.PublicKey
	PrivateKeyShare string `json:"PrivateKey"` //*tdh2easy.PrivateShare
}

func NewThresholdService(sharedOracleArgs *libocr2.OCR2OracleArgs, conf *ThresholdServicesConfig) (job.ServiceCtx, error) {
	kswpk, err := UnmarshalKeyshareWithPubKey(conf.KeyshareWithPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal threshold key share with public key")
	}

	// TODO: This multi-step un-marshaling may not be necessary
	var publicKey *tdh2easy.PublicKey
	err = publicKey.Unmarshal([]byte(kswpk.PublicKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal threshold encryption public key")
	}

	var privKeyShare *tdh2easy.PrivateShare
	err = privKeyShare.Unmarshal([]byte(kswpk.PrivateKeyShare))
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

func UnmarshalKeyshareWithPubKey(raw []byte) (*KeyshareWithPubKey, error) {
	var kwpk KeyshareWithPubKey
	err := json.Unmarshal(raw, &kwpk)
	if err != nil {
		return nil, err
	}
	return &kwpk, nil
}
