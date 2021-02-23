package gethwrappers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

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

	binaryData, err := replaceMetadata(beltArtifact)
	if err != nil {
		return nil, err
	}
	var sources struct {
		Sources map[string]string `json:"sourceCodes"`
	}
	if err := json.Unmarshal(beltArtifact, &sources); err != nil {
		return nil, errors.Wrapf(err, "could not read source code from compiler artifact")
	}
	return &ContractDetails{
		Binary:  binaryData,
		ABI:     contractABI,
		Sources: sources.Sources,
	}, nil
}

// extractSuffix returns the binary with the metadata replaced by a constant
// value, and any stored code (such as interned strings at the end of the
// deployed contract removed)
//
// Since the binary metadata suffix varies so much, it can't be included in a
// reliable check that the golang wrapper is up-to-date, so we replace it here
// https://solidity.readthedocs.io/en/v0.6.2/metadata.html#encoding-of-the-metadata-hash-in-the-bytecode
//
// Note that if the contract being wrapped uses the `new` keyword to deploy some
// other contract, the binary for wrapped contract must contain the contract to
// be deployed as well, and solc will *also* add metadata to that contract.
// Thus, there can be multiple metadata hashes which need to be replaced.
func replaceMetadata(beltArtifact []byte) (string, error) {
	// We want the bytecode here, not the deployedByteCode. The latter does not
	// include the initialization code.
	rawBinary := gjson.Get(string(beltArtifact),
		"compilerOutput.evm.bytecode.object").String()
	if rawBinary == "" {
		return "", errors.Errorf(
			"could not parse belt contract binary JSON as JSON")
	}
	if _, err := hexutil.Decode(rawBinary); err != nil {
		return "", errors.Errorf("contract binary from belt artifact is not JSON")
	}
	binLen := len(rawBinary)
	binRawBin := []byte(rawBinary)
	// See https://golang.org/pkg/regexp/#Overview for submatch structure, search
	// for "matches and submatches are identified by byte index"
	for _, metadataBounds := range metadata.FindAllSubmatchIndex(binRawBin, binLen) {
		for i := range constantBinaryMetadataHash {
			binRawBin[metadataBounds[0]+i] = constantBinaryMetadataHash[i]
		}
	}
	rawBinary = string(binRawBin)
	_, err := hex.DecodeString(rawBinary[2:])
	if err != nil {
		panic(errors.Wrapf(err, "can't fail, due to earlier hexutil.Decode check!"))
	}
	return rawBinary, nil
}

// binarySuffixRegexp checks that the hex representation of the trailing bytes
// of a contract wrapper follow the expected metadata format.
//
// Modern solc objects have metadata suffixes which vary depending on
// incidental compilation context like absolute paths to source files. See
// https://solidity.readthedocs.io/en/v0.6.6/metadata.html#encoding-of-the-metadata-hash-in-the-bytecode
//
// Note that the metadata hash is not necessarily the last part of the contract
// binary, as a naive reading of that documentation would imply. It can be
// followed by static code sections.
// https://gitter.im/ethereum/solidity?at=5f10a247d60398014655861f
var metadata = regexp.MustCompile("a264697066735822[[:xdigit:]]{68}" +
	"64736f6c6343[[:xdigit:]]{6}0033")

// constantBinaryMetadataHash is an arbitrary constant stand-in for the
// metadata suffix which the EVM expects (and in some cases, it seems, requires)
// to find at the end of the binary object representing a contract under
// deployment. See metadata docstring.
var constantBinaryMetadataHash = []byte("a264697066735822" +
	strings.Repeat("beef", 68/4) + "64736f6c6343" + "decafe" + "0033")

func init() {
	if !metadata.Match(constantBinaryMetadataHash) {
		panic("constantBinaryMetadataHash is invalid")
	}
}
