package changeset

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math/big"
	"strings"

	workflowUtils "github.com/smartcontractkit/chainlink-common/pkg/workflows"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/shared"
)

type WorkflowMetadata struct {
	Workflows map[string]string
}

func FeedIDsToBytes16(feedIDs []string) ([][16]byte, error) {
	dataIDs := make([][16]byte, len(feedIDs))
	for i, feedID := range feedIDs {
		err := shared.ValidateFeedID(feedID)
		if err != nil {
			return nil, err
		}
		dataIDs[i], err = ConvertHexToBytes16(feedID)
		if err != nil {
			return nil, err
		}
	}

	return dataIDs, nil
}

func FeedIDsToBytes(feedIDs []string) ([][]byte, error) {
	dataIDs16, err := FeedIDsToBytes16(feedIDs)
	if err != nil {
		return nil, err
	}

	dataSlices := make([][]byte, len(dataIDs16))
	for i, v := range dataIDs16 {
		b := make([]byte, 32)
		copy(b, v[:])
		dataSlices[i] = b
	}

	return dataSlices, nil
}

func ConvertHexToBytes16(hexStr string) ([16]byte, error) {
	if hexStr[:2] == "0x" {
		hexStr = hexStr[2:] // Remove "0x" prefix
	}
	decodedBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return [16]byte{}, fmt.Errorf("failed to decode hex string: %w", err)
	}

	var result [16]byte
	copy(result[:], decodedBytes[:16])

	return result, nil
}

func ExtractTypeAndVersion(hexStr string) (string, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", fmt.Errorf("invalid hex: %w", err)
	}

	if len(data) < 64 {
		return "", errors.New("data too short to be ABI-encoded")
	}

	// Extract the length (32 bytes from offset 32)
	lengthBytes := data[32:64]
	strLen := new(big.Int).SetBytes(lengthBytes).Int64()

	if strLen < 0 {
		return "", errors.New("negative string length")
	}

	if len(data) < 64+int(strLen) {
		return "", errors.New("data too short for expected string length")
	}

	strBytes := data[64 : 64+int(strLen)]
	return string(strBytes), nil
}

// HashedWorkflowName returns first 10 bytes of the sha256(workflow_name)
func HashedWorkflowName(name string) [10]byte {
	nameHash := workflowUtils.HashTruncateName(name)
	var result [10]byte
	copy(result[:], nameHash)
	return result
}

func LoadJSON[T any](pth string, fs fs.ReadFileFS) (T, error) {
	var dflt T
	f, err := fs.ReadFile(pth)
	if err != nil {
		return dflt, fmt.Errorf("failed to read %s: %w", pth, err)
	}
	var v T
	err = json.Unmarshal(f, &v)
	if err != nil {
		return dflt, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return v, nil
}

func GetDecimalsFromFeedID(feedID string) (uint8, error) {
	err := shared.ValidateFeedID(feedID)
	if err != nil {
		return 0, fmt.Errorf("invalid feed ID: %w", err)
	}
	feedIDBytes, err := ConvertHexToBytes16(feedID)
	if err != nil {
		return 0, fmt.Errorf("failed to convert feed ID to bytes: %w", err)
	}

	if feedIDBytes[7] >= 0x20 && feedIDBytes[7] <= 0x60 {
		return feedIDBytes[7] - 32, nil
	}

	return 0, nil
}

func GetDataFeedsCacheAddress(ab cldf.AddressBook, dataStore datastore.AddressRefStore, chainSelector uint64, label *string) string {
	var qualifier string
	if label != nil {
		qualifier = *label
	} else {
		qualifier = "data-feeds"
	}

	// try to find the address in datastore, fallback to addressbook
	record, err := dataStore.Get(
		datastore.NewAddressRefKey(chainSelector, DataFeedsCache, &deployment.Version1_0_0, qualifier),
	)
	if err == nil {
		return record.Address
	}

	// legacy addressbook
	dataFeedsCacheAddress := ""
	cacheTV := cldf.NewTypeAndVersion("DataFeedsCache", deployment.Version1_0_0)
	cacheTV.Labels.Add(qualifier)

	address, err := ab.AddressesForChain(chainSelector)
	if err != nil {
		return ""
	}

	for addr, tv := range address {
		if strings.Contains(tv.String(), cacheTV.String()) {
			dataFeedsCacheAddress = addr
		}
	}

	return dataFeedsCacheAddress
}

func UpdateWorkflowMetadataDS(
	env cldf.Environment,
	ds datastore.MutableDataStore,
	wfName string,
	wfSpec string,
) error {
	// environment metadata is overwritten with every Set(), so we need to read the existing metadata first
	record, err := env.DataStore.EnvMetadata().Get()
	if err != nil {
		// if the datastore is not initialized, we should create a new one
		env.Logger.Errorf("failed to get env datastore: %v", err)
	}

	metadata, err := datastore.As[WorkflowMetadata](record.Metadata)
	if err != nil {
		return fmt.Errorf("failed to cast env metadata: %w", err)
	}

	if metadata.Workflows == nil {
		metadata.Workflows = make(map[string]string)
	}

	// upsert the workflow spec in the metadata
	metadata.Workflows[wfName] = wfSpec

	err = ds.EnvMetadata().Set(
		datastore.EnvMetadata{
			Metadata: metadata,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to set workflow spec in datastore: %w", err)
	}

	return nil
}
