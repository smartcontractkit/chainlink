package crib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

type OutputReader struct {
	cribEnvStateDirPath string
}

// NewOutputReader creates new instance
func NewOutputReader(cribEnvStateDirPath string) *OutputReader {
	return &OutputReader{cribEnvStateDirPath: cribEnvStateDirPath}
}

func (r *OutputReader) ReadNodesDetails() (NodesDetails, error) {
	var result NodesDetails
	byteValue, err := r.readCRIBDataFile(NodesDetailsFileName)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return result, err
	}

	return result, nil
}

func (r *OutputReader) ReadRMNNodeConfigs() ([]RMNNodeConfig, error) {
	var result []RMNNodeConfig
	byteValue, err := r.readCRIBDataFile(RMNNodeIdentitiesFileName)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return result, err
	}

	return result, nil
}

func (r *OutputReader) ReadChainConfigs() ([]devenv.ChainConfig, error) {
	var result []devenv.ChainConfig
	byteValue, err := r.readCRIBDataFile(ChainsConfigsFileName)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return result, err
	}

	return result, nil
}

func (r *OutputReader) ReadAddressBook() (*cldf.AddressBookMap, error) {
	var result map[uint64]map[string]cldf.TypeAndVersion
	byteValue, err := r.readCRIBDataFile(AddressBookFileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}

	return cldf.NewMemoryAddressBookFromMap(result), nil
}

// ReadDataStore loads the DataStore from address_refs.json and metadata files in the data directory.
// Falls back gracefully if files are missing (returns empty DataStore).
func (r *OutputReader) ReadDataStore() (*datastore.MemoryDataStore, error) {
	ds := datastore.NewMemoryDataStore()

	refs, err := r.readCRIBDataFile(AddressRefsFileName)
	if err != nil {
		// address_refs.json is required for a useful DataStore
		return ds, fmt.Errorf("failed to read %s: %w", AddressRefsFileName, err)
	}
	if len(refs) > 0 {
		if err = json.Unmarshal(refs, &ds.AddressRefStore.Records); err != nil {
			return nil, fmt.Errorf("failed to unmarshal address refs: %w", err)
		}
	}

	// Metadata files are optional — load if present
	if chainMeta, err := r.readCRIBDataFile(ChainMetadataFileName); err == nil && len(chainMeta) > 0 {
		_ = json.Unmarshal(chainMeta, &ds.ChainMetadataStore.Records)
	}
	if ctrMeta, err := r.readCRIBDataFile(ContractMetadataFileName); err == nil && len(ctrMeta) > 0 {
		_ = json.Unmarshal(ctrMeta, &ds.ContractMetadataStore.Records)
	}
	if envMeta, err := r.readCRIBDataFile(EnvMetadataFileName); err == nil && len(envMeta) > 0 {
		_ = json.Unmarshal(envMeta, &ds.EnvMetadataStore.Record)
	}

	return ds, nil
}

func (r *OutputReader) readCRIBDataFile(fileName string) ([]byte, error) {
	dataDirPath := path.Join(r.cribEnvStateDirPath, "data")
	file, err := os.Open(fmt.Sprintf("%s/%s", dataDirPath, fileName))
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	// Read the file's content into a byte slice
	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}
	return byteValue, nil
}
