package ccip

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"io"
	"os"
	"time"
)

const (
	transmitted = iota
	committed
	executed
	LokiLoadLabel = "ccip_load_test"
	ErrLokiClient = "failed to create Loki client for monitoring"
	ErrLokiPush   = "failed to push metrics to Loki"
)

// todo: Have a different struct for commit/exec?
type LokiMetric struct {
	EventType      int       `json:"event_type"`
	Timestamp      time.Time `json:"timestamp"`
	GasUsed        uint64    `json:"gas_used"`
	SequenceNumber uint64    `json:"sequence_number"`
}

func GetAddressFromTypeAndVersion(ab deployment.AddressBook, cs uint64, tv string) (common.Address, error) {
	allAddr, err := ab.AddressesForChain(cs)
	if err != nil {
		return common.Address{}, err
	}
	for addr, tv := range allAddr {
		if tv.Type == tv.Type && tv.Version == tv.Version {
			return common.HexToAddress(addr), nil
		}
	}

	return common.Address{}, fmt.Errorf("address not found for chain selector %d and typeAndVersion %s", cs, tv)
}

func SendMetricsToLoki(l zerolog.Logger, lc *wasp.LokiClient, updatedLabels map[string]string, metrics *LokiMetric) {
	if err := lc.HandleStruct(wasp.LabelsMapToModel(updatedLabels), time.Now(), metrics); err != nil {
		l.Error().Err(err).Msg(ErrLokiPush)
	}
}

func setLokiLabels(src, dst uint64) map[string]string {
	return map[string]string{
		"sourceSelector":      fmt.Sprintf("%d", src),
		"destinationSelector": fmt.Sprintf("%d", dst),
		"testType":            LokiLoadLabel,
	}
}

func readFile(inputDir string, fileName string) []byte {
	file, err := os.Open(fmt.Sprintf("%s/%s", inputDir, fileName))
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

func ReadAddressBook(inputDir string) *deployment.AddressBookMap {
	byteValue := readFile(inputDir, AddressBookFileName)

	var result map[uint64]map[string]deployment.TypeAndVersion

	// Unmarshal the JSON into the map
	err := json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic(err)
	}

	// Print the deserialized map
	fmt.Println(result)
	return deployment.NewMemoryAddressBookFromMap(result)
}

func ReadNodesDetails(inputDir string) NodesDetails {
	byteValue := readFile(inputDir, NodesDetailsFileName)

	var result NodesDetails

	// Unmarshal the JSON into the map
	err := json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic(err)
	}

	// Print the deserialized map
	fmt.Println(result)
	return result
}

func NewDeployEnvironmentFromCribOutput(lggr logger.Logger, ab deployment.AddressBook, nodeIDs []string) (*deployment.Environment, error) {
	chains, err := devenv.NewChains(lggr, output.Chains)
	if err != nil {
		return nil, err
	}
	return deployment.NewEnvironment(
		"Crib Environment",
		lggr,
		ab,
		chains,
		nodeIDs,
		nil, // todo: populate the offchain client using output.DON
	), nil
}
