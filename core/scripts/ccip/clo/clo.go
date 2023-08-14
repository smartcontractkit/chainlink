package clo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea"
)

var url = os.Getenv("CLO_QUERY_URL")
var authToken = os.Getenv("CLO_AUTH_TOKEN")

// SetChainConfig overrides any existing config values for given chain
func SetChainConfig(contracts []ContractC, deploymentConfig *rhea.EvmDeploymentConfig) {
	for _, c := range contracts {
		switch c.Name {
		case "MockARMContract":
			deploymentConfig.ChainConfig.ARM = common.HexToAddress(c.Address)
		case "Router":
			deploymentConfig.ChainConfig.Router = common.HexToAddress(c.Address)
		case "PriceRegistry":
			deploymentConfig.ChainConfig.PriceRegistry = common.HexToAddress(c.Address)
		}
	}
}

// SetLaneConfig overrides any existing config values for given leg
func SetLaneConfig(leg LegL, source *rhea.EvmDeploymentConfig, dest *rhea.EvmDeploymentConfig) {
	for _, c := range leg.Source.Contracts {
		switch c.Name {
		case "PingPongDemo":
			source.LaneConfig.PingPongDapp = common.HexToAddress(c.Address)
		case "EVM2EVMOnRamp":
			source.LaneConfig.OnRamp = common.HexToAddress(c.Address)
		}
	}
	for _, c := range leg.Destination.Contracts {
		switch c.Name {
		case "CommitStore":
			dest.LaneConfig.CommitStore = common.HexToAddress(c.Address)
		case "EVM2EVMOffRamp":
			dest.LaneConfig.OffRamp = common.HexToAddress(c.Address)
		}
	}
}

func GetTargetChainsContracts(t *testing.T, sourceChainId uint64, destChainId uint64) ([]ContractC, []ContractC) {
	sourceChainIdStr := strconv.FormatUint(sourceChainId, 10)
	destChainIdStr := strconv.FormatUint(destChainId, 10)
	var sourceChain ChainC
	var destChain ChainC

	var responseBody GenericGraphQLResponseBody[ListChainsResponseData]
	responseBody = requestCLO(t, "POST", url, listChains, responseBody)
	for _, chain := range responseBody.Data.Ccip.Chains {
		if chain.Network.ChainID == sourceChainIdStr {
			sourceChain = ChainC{
				ID:        chain.ID,
				Network:   chain.Network,
				Contracts: chain.Contracts,
			}
		} else if chain.Network.ChainID == destChainIdStr {
			destChain = ChainC{
				ID:        chain.ID,
				Network:   chain.Network,
				Contracts: chain.Contracts,
			}
		}
	}

	return sourceChain.Contracts, destChain.Contracts
}

func GetTargetLaneConfig(t *testing.T, sourceChainId uint64, destChainId uint64, laneID string) (LegL, LegL) {
	var lane LaneL
	var LegA LegL
	var LegB LegL
	sourceChainIdStr := strconv.FormatUint(sourceChainId, 10)
	destChainIdStr := strconv.FormatUint(destChainId, 10)

	var responseBody GenericGraphQLResponseBody[ListLanesResponse]
	responseBody = requestCLO(t, "POST", url, listLanes, responseBody)
	for _, lane = range responseBody.Data.Ccip.Lanes {
		if lane.ID == laneID {
			if lane.LegA.Source.Chain.Network.ChainID == sourceChainIdStr && lane.LegB.Source.Chain.Network.ChainID == destChainIdStr {
				LegA = LegL{
					Source:      lane.LegA.Source,
					Destination: lane.LegA.Destination,
				}
				LegB = LegL{
					Source:      lane.LegB.Source,
					Destination: lane.LegB.Destination,
				}
			} else {
				LegA = LegL{
					Source:      lane.LegB.Source,
					Destination: lane.LegB.Destination,
				}
				LegB = LegL{
					Source:      lane.LegA.Source,
					Destination: lane.LegA.Destination,
				}
			}
		}
	}
	return LegA, LegB
}

func requestCLO[T any](t *testing.T, method string, url string, query string, responseBody T) T {
	if url == "" {
		t.Error("CLO_QUERY_URL must be set")
	}
	if authToken == "" {
		t.Error("CLO_AUTH_TOKEN must be set")
	}
	requestBody := map[string]interface{}{
		"query": query,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Session-Token", authToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	defer response.Body.Close()

	responseBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	err = json.Unmarshal(responseBodyBytes, &responseBody)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	//// Write the response body to a file
	// timestamp := time.Now().Format("2006-01-02T15:04:05.000-07:00")
	// filename := "clo/response_" + timestamp + ".json"
	// err = os.WriteFile(filename, responseBodyBytes, 0644)
	// if err != nil {
	//	 t.Errorf("%s", err.Error())
	// }

	return responseBody
}
