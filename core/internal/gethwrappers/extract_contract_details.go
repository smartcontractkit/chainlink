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
		Binary:  binaryData.String(),
		ABI:     contractABI,
		Sources: sources.Sources,
	}, nil
}

// binaryContractSections represents an analysis of the contract binary produced
// by solc. Constituent strings are the hex representations of those sections
type binaryContractSections struct {
	contractBinary string // The actual contract bytecode
	// The contract binary is followed by a hash of the contract metadata Because
	// the metadata can vary from compile to compile, it's necessary to replace
	// this with a constant.
	// https://solidity.readthedocs.io/en/v0.6.2/metadata.html#encoding-of-the-metadata-hash-in-the-bytecode
	metadataHash string
	// If any static data such as often-used strings need to be stored, they come
	// after the metadata hash
	// https://gitter.im/ethereum/solidity?at=5f10a247d60398014655861f
	staticData string
}

// String returns the concatenated contract binary, as hex
func (b *binaryContractSections) String() string {
	return b.contractBinary + b.metadataHash + b.staticData
}

// extractSuffix returns the binary with the metadata replaced by a constant
// value, and any stored code (such as interned strings at the end of the
// deployed contract removed)
//
// Since the binary metadata suffix varies so much, it can't be included in a
// reliable check that the golang wrapper is up-to-date, so we replace it here
func replaceMetadata(beltArtifact []byte) (*binaryContractSections, error) {
	// We want the bytecode here, not the deployedByteCode. The latter does not
	// include the initialization code.
	rawBinary := gjson.Get(string(beltArtifact),
		"compilerOutput.evm.bytecode.object").String()
	if rawBinary == "" {
		return nil, errors.Errorf(
			"could not parse belt contract binary JSON as JSON")
	}
	if _, err := hexutil.Decode(rawBinary); err != nil {
		return nil, errors.Errorf("contract binary from belt artifact is not JSON")
	}
	// Since regexp.Regexp.FindAll returns successive NON-OVERLAPPING matches, and
	// we want the FINAL match, it is necessary to reverse the binary data and
	// regexp, so that we can search from the front.
	//
	// See https://golang.org/pkg/regexp/#Overview for submatch structure, search
	// for "matches and submatches are identified by byte index"
	revMetadataBounds := reversedMetaData.FindAllStringSubmatchIndex(
		utils.ReverseString(rawBinary), 1)
	// We want the last submatch
	revBounds := revMetadataBounds[0][len(revMetadataBounds[0])-2:]
	metadataStart := len(rawBinary) - revBounds[1]
	metadataEnd := len(rawBinary) - revBounds[0]
	metadataHash := rawBinary[metadataStart:metadataEnd]
	if !metadata.MatchString(metadataHash) {
		panic(errors.Errorf("retrieved bad metadata: %s", metadataHash))
	}
	rv := &binaryContractSections{
		contractBinary: rawBinary[:metadataStart],
		metadataHash:   constantBinaryMetadataHash,
		staticData:     rawBinary[metadataEnd:],
	}
	if len(rv.String()) != len(rawBinary) {
		panic("reconstruction with fixed metadata hash failed")
	}
	rawCode, err := hex.DecodeString(rv.staticData)
	if err != nil {
		panic(errors.Wrapf(err, "can't fail, due to earlier hexutil.Decode check!"))
	}
	if !ascii.Match(rawCode) {
		// So far, we have only seen the compiler store ASCII strings after the
		// metadata hash. Fail noisily if we see something else(?)
		panic(errors.Errorf("non-ascii"))
	}
	return rv, nil
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

var reversedMetaData = regexp.MustCompile("(..)*(3300[[:xdigit:]]{6}" +
	"3436c6f63746[[:xdigit:]]{68}228537660796462a)")

// constantBinaryMetadataHash is an arbitrary constant stand-in for the
// metadata suffix which the EVM expects (and in some cases, it seems, requires)
// to find at the end of the binary object representing a contract under
// deployment. See metadata docstring.
var constantBinaryMetadataHash = "a264697066735822" +
	strings.Repeat("beef", 68/4) + "64736f6c6343" + "decafe" + "0033"

func init() {
	if !metadata.MatchString(constantBinaryMetadataHash) {
		panic("constantBinaryMetadataHash is invalid")
	}
	if !reversedMetaData.MatchString(utils.ReverseString(
		constantBinaryMetadataHash)) {
		panic("constantBinaryMetadataHash does not match reversedMetaData regex")
	}
}

// ascii matches any ascii string
var ascii = regexp.MustCompile("^[[:ascii:]]*$")
