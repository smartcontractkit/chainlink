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

func (r *OutputReader) ReadNodesDetails() NodesDetails {
	byteValue := r.readCRIBDataFile(NodesDetailsFileName)

	var result NodesDetails

	err := json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic(err)
	}

	return result
}

func (r *OutputReader) ReadRMNNodeConfigs() []RMNNodeConfig {
	byteValue := r.readCRIBDataFile(RMNNodeIdentitiesFileName)

	var result []RMNNodeConfig

	err := json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic(err)
	}

	return result
}

func (r *OutputReader) ReadChainConfigs() []devenv.ChainConfig {
	byteValue := r.readCRIBDataFile(ChainsConfigsFileName)

	var result []devenv.ChainConfig

	err := json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic(err)
	}

	return result
}

func (r *OutputReader) ReadAddressBook() *deployment.AddressBookMap {
	byteValue := r.readCRIBDataFile(AddressBookFileName)

	var result map[uint64]map[string]deployment.TypeAndVersion

	err := json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic(err)
	}

	return deployment.NewMemoryAddressBookFromMap(result)
}

func (r *OutputReader) readCRIBDataFile(fileName string) []byte {
	dataDirPath := path.Join(r.cribEnvStateDirPath, "data")
	file, err := os.Open(fmt.Sprintf("%s/%s", dataDirPath, fileName))
	if err != nil {
		fmt.Println("Error opening file:", err)
		panic(err)
	}
	defer file.Close()

	// Read the file's content into a byte slice
	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		panic(err)
	}
	return byteValue
}
