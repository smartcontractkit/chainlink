package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

func main() {
	var (
		cmd             = flag.NewFlagSet("verify", flag.ExitOnError)
		chainID         = cmd.Int("chain-id", 1, "chain ID of the chain the contract is deployed on")
		contractName    = cmd.String("contract-name", "", "name of the contract to verify")
		contractAddress = cmd.String("contract-address", "", "address of the contract to verify")
		// TODO: have users provide just a version and we can fill out the rest by requesting https://etherscan.io/solcversions and pattern matching
		// Or perhaps https://raw.githubusercontent.com/ethereum/solc-bin/gh-pages/bin/list.txt
		compilerVersion  = cmd.String("compiler-version", "", "solc compiler version used to compile contract")
		optimizationUsed = cmd.Bool("optimization-used", true, "whether optimization was enabled")
		constructorArgs  = cmd.String("constructor-args", "", "abi-encoded constructor args (if necessary)")
		standardJSONPath = cmd.String("std-json", "", "path to solc StandardJSON output of the contract compilation")
		apiKey           = cmd.String("api-key", "", "etherscan or similar API key")
	)

	helpers.ParseArgs(cmd, os.Args[1:],
		"contract-address", "std-json", "api-key", "contract-name")

	apiURL, ok := helpers.EtherscanVerifyEndpoints[*chainID]
	if !ok {
		panic(
			fmt.Errorf(
				"etherscan verification endpoint not found for chain ID %d - either add one or it doesn't exist?",
				*chainID))
	}

	data := url.Values{}

	// Request metadata
	data.Set("apikey", *apiKey)
	data.Set("module", "contract")
	data.Set("action", "verifysourcecode")
	data.Set("codeformat", "solidity-standard-json-input")

	// Contract information
	// See https://docs.etherscan.io/tutorials/verifying-contracts-programmatically#4.-configuring-source-code-parameters
	// for required parameters.
	data.Set("sourceCode", readStandardJSON(*standardJSONPath))
	data.Set("contractname", *contractName)
	data.Set("contractaddress", *contractAddress)
	data.Set("compilerversion", *compilerVersion)                       // NOTE: not clear if this is necessary if we pass in standardjson
	data.Set("optimizationUsed", strconv.FormatBool(*optimizationUsed)) // NOTE: not clear if this is necessary if we pass in standardjson

	fmt.Println("encoded:", data.Encode())

	if constructorArgs != nil {
		data.Set("constructorArguements", *constructorArgs) // NOTE: typo not ours
	}

	request, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(data.Encode()))
	helpers.PanicErr(err)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println("sending request...")

	client := &http.Client{}
	resp, err := client.Do(request)
	helpers.PanicErr(err)

	// Get GUID from response and poll the checking endpoint
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("expected status %d got %d", http.StatusOK, resp.StatusCode))
	}

	fmt.Println("reading response...")

	respBytes, err := ioutil.ReadAll(resp.Body)
	helpers.PanicErr(err)

	var eResp etherscanResponse
	err = json.Unmarshal(respBytes, &eResp)
	if err != nil {
		fmt.Println("response:", string(respBytes))
		helpers.PanicErr(err)
	}

	fmt.Println("response:", eResp, "now polling for completion...")

	// Follow https://docs.etherscan.io/api-endpoints/contracts#check-source-code-verification-submission-status
	// to check verification status
	data = url.Values{}
	data.Set("apikey", *apiKey)
	data.Set("guid", eResp.Result) // GUID is stored in the "result" field
	data.Set("module", "contract")
	data.Set("action", "checkverifystatus")

	request, err = http.NewRequest("GET", apiURL, strings.NewReader(data.Encode()))
	helpers.PanicErr(err)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for i := 0; i < 5; i++ {
		resp, err := client.Do(request)
		helpers.PanicErr(err)

		fmt.Println("reading polling response...")

		respBytes, err := ioutil.ReadAll(resp.Body)
		helpers.PanicErr(err)

		var eResp etherscanResponse
		helpers.PanicErr(json.Unmarshal(respBytes, &eResp))

		fmt.Println("response:", eResp)
		if eResp.Message != "OK" {
			fmt.Println("contract not verified yet, trying again in 5 seconds. Result:", eResp.Result)
			time.Sleep(5 * time.Second)
		} else {
			fmt.Println("contract verified, done!")
			return
		}
	}

}

type etherscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func readStandardJSON(path string) string {
	fileBytes, err := ioutil.ReadFile(path)
	helpers.PanicErr(err)
	return string(fileBytes)
}
