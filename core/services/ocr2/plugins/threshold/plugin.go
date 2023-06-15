package threshold

import (
	"encoding/json"
	"fmt"

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
	fmt.Println("Private Key Share With Public Key: ", string(conf.KeyshareWithPubKey))

	publicKey, privKeyShare, err := UnmarshalKeys(conf.KeyshareWithPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal threshold key share with public key")
	}

	// It can be safely assumed all oracle nodes have been given the correct key share,
	// so we use the key share index as the oracle's index.
	// If an oracle was given the wrong key share, key share decryption would have failed.
	// oracleToKeyShare := make(map[commontypes.OracleID]int)
	// oracleToKeyShare[commontypes.OracleID(privKeyShare.Index())] = privKeyShare.Index()
	oracleToKeyShare := map[commontypes.OracleID]int{0: 0, 1: 1, 2: 2, 3: 3}

	fmt.Println("THRESHOLD ORACLE INDEX:", privKeyShare.Index())
	fmt.Println("THRESHOLD ORACLETOKEYSHARE:", oracleToKeyShare)

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
	fmt.Println("Public Key:", publicKey)

	err = privateShare.Unmarshal(kwpk.PrivateKeyShare)
	if err != nil {
		return publicKey, privateShare, err
	}

	checkPubKey, err := publicKey.Marshal()
	if err != nil {
		return publicKey, privateShare, err
	}
	checkPrivKey, err := privateShare.Marshal()
	if err != nil {
		return publicKey, privateShare, err
	}
	fmt.Printf("Threshold Public Key: %+v\nThreshold Private Share: %+v\n", string(checkPubKey), string(checkPrivKey))

	return publicKey, privateShare, nil
}
