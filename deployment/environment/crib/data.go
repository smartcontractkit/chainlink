package crib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

type OutputReader struct {
	outputDir string
}

func NewOutputReader(outputDir string) *OutputReader {
	return &OutputReader{outputDir: outputDir}
}

func (r *OutputReader) ReadNodesDetails() NodesDetails {
	byteValue := r.readFile(NodesDetailsFileName)

	var result NodesDetails

	err := json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic(err)
	}

	return result
}

func (r *OutputReader) ReadRMNNodeConfigs() []RMNNodeConfig {
	byteValue := r.readFile(RMNNodeIdentitiesFileName)

	var result []RMNNodeConfig

	err := json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic(err)
	}

	return result
}

func (r *OutputReader) ReadChainConfigs() []devenv.ChainConfig {
	byteValue := r.readFile(ChainsConfigsFileName)

	var result []devenv.ChainConfig

	err := json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic(err)
	}

	return result
}

func (r *OutputReader) ReadAddressBook() *deployment.AddressBookMap {
	byteValue := r.readFile(AddressBookFileName)

	var result map[uint64]map[string]deployment.TypeAndVersion

	err := json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic(err)
	}

	return deployment.NewMemoryAddressBookFromMap(result)
}

func (r *OutputReader) readFile(fileName string) []byte {
	file, err := os.Open(fmt.Sprintf("%s/%s", r.outputDir, fileName))
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
