package ccip

import (
	"context"
	"encoding/json"
	"fmt"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
	"time"
)

const (
	transmitted = iota
	committed
	executed
	LokiLoadLabel   = "ccip_load_test"
	ErrLokiClient   = "failed to create Loki client for monitoring"
	ErrLokiPush     = "failed to push metrics to Loki"
	abPath          = "/Users/austin.wang/ccip-core/repos/chainlink/integration-tests/load/ccip/testfiles/ccip-v2-scripts-address-book.json"
	nodeIdsPath     = "/Users/austin.wang/ccip-core/repos/chainlink/integration-tests/load/ccip/testfiles/ccip-v2-scripts-node-details.json"
	chainConfigPath = "/Users/austin.wang/ccip-core/repos/chainlink/integration-tests/load/ccip/testfiles/ccip-v2-scripts-chains-details.json"
)

// todo: Have a different struct for commit/exec?
type LokiMetric struct {
	EventType      int       `json:"event_type"`
	Timestamp      time.Time `json:"timestamp"`
	GasUsed        uint64    `json:"gas_used"`
	SequenceNumber uint64    `json:"sequence_number"`
}

func SendMetricsToLoki(l logger.Logger, lc *wasp.LokiClient, updatedLabels map[string]string, metrics *LokiMetric) {
	if err := lc.HandleStruct(wasp.LabelsMapToModel(updatedLabels), time.Now(), metrics); err != nil {
		l.Error(ErrLokiPush)
	}
}

func setLokiLabels(src, dst uint64) (map[string]string, error) {
	srcChainId, err := chainselectors.GetChainIDFromSelector(src)
	if err != nil {
		return nil, err
	}
	dstChainId, err := chainselectors.GetChainIDFromSelector(dst)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"sourceEvmChainId":    fmt.Sprintf("%s", srcChainId),
		"destEvmChainId":      fmt.Sprintf("%s", dstChainId),
		"destinationSelector": fmt.Sprintf("%d", dst),
		"testType":            LokiLoadLabel,
	}, nil
}

func readFile(filePath string) []byte {
	file, err := os.Open(filePath)
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

func readFromFile[T []string | *deployment.AddressBookMap | []devenv.ChainConfig](t *testing.T, inputDir string) T {
	byteValue := readFile(inputDir)

	var result T
	// Unmarshal the JSON into the map
	err := json.Unmarshal(byteValue, &result)
	require.NoError(t, err)

	// Print the deserialized map
	fmt.Println(result)
	return result
}

func CreateEnvironmentFromCribOutput(t *testing.T, lggr logger.Logger) (*deployment.Environment, error) {
	ab := readFromFile[*deployment.AddressBookMap](t, abPath)
	nodeIds := readFromFile[[]string](t, nodeIdsPath)
	chainDetails := readFromFile[[]devenv.ChainConfig](t, chainConfigPath)

	// todo: make sure to call chainDetails.SetDeployerKey() for each chain
	// where private keys should be stored in env vars or .toml

	chains, err := devenv.NewChains(lggr, chainDetails)
	if err != nil {
		return nil, err
	}
	return deployment.NewEnvironment(
		"Crib Environment",
		lggr,
		ab,
		chains,
		nil,
		nodeIds,
		nil,
		func() context.Context { return context.Background() },
		deployment.XXXGenerateTestOCRSecrets(),
	), nil
}
