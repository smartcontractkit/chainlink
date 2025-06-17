package por

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
)

func CreateConfigFile(feedsConsumerAddress common.Address, feedID, dataURL, writeTargetName string, authKeySecretName *string) (*os.File, error) {
	configFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create workflow config file")
	}

	cleanFeedID := strings.TrimPrefix(feedID, "0x")
	feedLength := len(cleanFeedID)

	if feedLength < 32 {
		return nil, errors.Errorf("feed ID must be at least 32 characters long, but was %d", feedLength)
	}

	if feedLength > 32 {
		cleanFeedID = cleanFeedID[:32]
	}

	feedIDToUse := "0x" + cleanFeedID

	workflowConfig := libcrecli.PoRWorkflowConfig{
		FeedID:            feedIDToUse,
		URL:               dataURL,
		ConsumerAddress:   feedsConsumerAddress.Hex(),
		WriteTargetName:   writeTargetName,
		AuthKeySecretName: authKeySecretName,
	}

	configMarshalled, err := json.Marshal(workflowConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal workflow config")
	}

	_, err = configFile.Write(configMarshalled)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write workflow config file")
	}

	return configFile, nil
}
