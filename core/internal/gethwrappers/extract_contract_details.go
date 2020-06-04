package gethwrappers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

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

// c.VersionHash() is the hash used to detect changes in the underlying contract
func (c *ContractDetails) VersionHash() (hash string) {
	hashMsg := c.ABI + c.Binary + "\n"
	return fmt.Sprintf("%x", sha256.Sum256([]byte(hashMsg)))
}

// ExtractContractDetails returns the data in the belt artifact needed to make a
// geth contract wrapper for the corresponding EVM contract
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

	// We want the bytecode here, not the deployedByteCode. The latter does not
	// include the initialization code.
	rawBinary := gjson.Get(string(beltArtifact),
		"compilerOutput.evm.bytecode.object").String()
	if rawBinary == "" {
		return nil, errors.Errorf(
			"could not parse belt contract binary JSON as JSON")
	}

	// Since the binary metadata suffix varies so much, it can't be included in a
	// reliable check that the golang wrapper is up-to-date, so remove it from the
	// message hash.
	truncLen := len(rawBinary) - 106
	truncatedBinary := rawBinary[:truncLen]
	suffix := rawBinary[truncLen:]
	// Verify that the suffix follows the pattern outlined in the above link, to
	// ensure that we're actually truncating what we think we are.
	if !binarySuffixRegexp(suffix) {
		return nil, errors.Errorf(
			"binary suffix has unexpected format; giving up: "+suffix, nil)
	}
	var sources struct {
		Sources map[string]string `json:"sourceCodes"`
	}
	if err := json.Unmarshal(beltArtifact, &sources); err != nil {
		return nil, errors.Wrapf(err, "could not read source code from compiler artifact")
	}
	return &ContractDetails{
		Binary:  truncatedBinary + constantBinaryMetadataSuffix,
		ABI:     contractABI,
		Sources: sources.Sources,
	}, nil
}

// binarySuffixRegexp checks that the hex representation of the trailing bytes
// of a contract wrapper follow the expected metadata format.
//
// Modern solc objects have metadata suffixes which vary depending on
// incidental compilation context like absolute paths to source files. See
// https://solidity.readthedocs.io/en/v0.6.2/metadata.html#encoding-of-the-metadata-hash-in-the-bytecode
var binarySuffixRegexp = regexp.MustCompile(
	"^a264697066735822[[:xdigit:]]{68}64736f6c6343[[:xdigit:]]{6}0033$",
).MatchString

// constantBinaryMetadataSuffix is a constant stand-in for the metadata suffix
// which the EVM expects (and in some cases, it seems, requires) to find at the
// end of the binary object representing a contract under deployment.
var constantBinaryMetadataSuffix = "a264697066735822" +
	strings.Repeat("beef", 68/4) + "64736f6c6343" + "decafe" + "0033"
