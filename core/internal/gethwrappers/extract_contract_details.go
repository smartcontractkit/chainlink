package gethwrappers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// ContractDetails contains the contract data needed to make a geth contract
// wrapper for a solidity contract.
type ContractDetails struct {
	Binary  string // Hex representation of the contract's  raw bytes
	ABI     string
	Sources map[string]string // contractName -> source code
}

// VersionHash is the hash used to detect changes in the underlying contract
func (c *ContractDetails) VersionHash() (hash string) {
	hashMsg := c.ABI + c.Binary + "\n"
	return fmt.Sprintf("%x", sha256.Sum256([]byte(hashMsg)))
}

// ExtractContractDetails returns the data in the belt artifact needed to make a
// geth contract wrapper for the corresponding EVM contract
//
// This should be run on artifacts compiled with no metadata hash
// (
//
//   --metadata-hash=none
//
// in the solc command, or
//
//   "metadata": {"bytecodeHash": "none"}
//
// in the sol-compiler config JSON [evm-contracts/app.config.json in the
// Chainlink source code.])
//
// The main drawback to including the metadata hash is that it varies from build
// to build, which causes a lot of false positives in the
// TestCheckContractHashesFromLastGoGenerate and
// TestArtifactCompilerVersionMatchesConfig checks that the golang wrappers are
// up to date with the contract source codes. It also increases contract size,
// slightly, which increases deployment costs.
func ExtractContractDetails(beltArtifactPath string) (*ContractDetails, error) {
	beltArtifact, err := ioutil.ReadFile(beltArtifactPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read belt artifact; you may need "+
			" to run`yarn && yarn workspace @chainlink/contracts belt compile solc`")
	}
	rawABI := gjson.Get(string(beltArtifact), "compilerOutput.abi")
	contractABI, err := utils.NormalizedJSON([]byte(rawABI.String()))
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse belt ABI JSON as JSON")
	}

	details := &ContractDetails{
		ABI: contractABI,
	}

	// We want the bytecode here, not the deployedByteCode. The latter does not
	// include the initialization code.
	binaryData := gjson.Get(string(beltArtifact),
		"compilerOutput.evm.bytecode.object").String()
	if binaryData != "" {
		if _, errr := hexutil.Decode(binaryData); errr != nil {
			return nil, errors.Wrapf(errr, "contract binary from belt artifact is not hex data")
		}
		details.Binary = binaryData
	}

	var sources struct {
		Sources map[string]string `json:"sourceCodes"`
	}
	err = json.Unmarshal(beltArtifact, &sources)
	if err == nil {
		details.Sources = sources.Sources
	}
	return details, nil
}
