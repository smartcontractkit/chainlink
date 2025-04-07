package crib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/smartcontractkit/chainlink/deployment"
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

func (r *OutputReader) ReadAddressBook() (*deployment.AddressBookMap, error) {
	var result map[uint64]map[string]deployment.TypeAndVersion
	byteValue, err := r.readCRIBDataFile(AddressBookFileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}

	return deployment.NewMemoryAddressBookFromMap(result), nil
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
