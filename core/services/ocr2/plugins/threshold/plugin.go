package threshold

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	decryptionPlugin "github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin"
	decryptionPluginConfig "github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin/config"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type ThresholdServicesConfig struct {
	DecryptionQueue    decryptionPlugin.DecryptionQueuingService
	KeyshareWithPubKey []byte
	ConfigParser       decryptionPluginConfig.ConfigParser
}

func NewThresholdService(sharedOracleArgs *libocr2.OCR2OracleArgs, conf *ThresholdServicesConfig) (job.ServiceCtx, error) {
	publicKey, privKeyShare, err := UnmarshalKeys(conf.KeyshareWithPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal threshold key share with public key")
	}

	// The key generation tooling ensures that key IDs correspond to the oracle's index,
	// therefore an identity mapping is used when creating the threshold reporting plugin.
	// maxNumNodes is selected such that it will always be larger than the number of nodes in the DON.
	oracleToKeyShare := make(map[commontypes.OracleID]int)
	maxNumNodes := 100
	for i := 0; i <= maxNumNodes; i++ {
		oracleToKeyShare[commontypes.OracleID(i)] = i
	}

	sharedOracleArgs.ReportingPluginFactory = decryptionPlugin.DecryptionReportingPluginFactory{
		DecryptionQueue:  conf.DecryptionQueue,
		ConfigParser:     conf.ConfigParser,
		PublicKey:        &publicKey,
		PrivKeyShare:     &privKeyShare,
		OracleToKeyShare: oracleToKeyShare,
		Logger:           sharedOracleArgs.Logger,
	}

	thresholdReportingPluginOracle, err := libocr2.NewOracle(*sharedOracleArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call NewOracle to create a Threshold Reporting Plugin")
	}

	return job.NewServiceAdapter(thresholdReportingPluginOracle), nil
}

type KeyshareWithPubKey struct {
	PublicKey       json.RawMessage //tdh2easy.PublicKey
	PrivateKeyShare json.RawMessage //tdh2easy.PrivateShare
}

func UnmarshalKeys(raw []byte) (publicKey tdh2easy.PublicKey, privateShare tdh2easy.PrivateShare, err error) {
	var kwpk KeyshareWithPubKey
	err = json.Unmarshal(raw, &kwpk)
	if err != nil {
		return publicKey, privateShare, err
	}

	err = publicKey.Unmarshal(kwpk.PublicKey)
	if err != nil {
		return publicKey, privateShare, err
	}

	err = privateShare.Unmarshal(kwpk.PrivateKeyShare)
	if err != nil {
		return publicKey, privateShare, err
	}

	return publicKey, privateShare, nil
}
