package eth

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// Contract holds the solidity contract's parsed ABI
type Contract struct {
	ABI abi.ABI
}

func getContract(name string, box packr.Box) (*Contract, error) {
	jsonFile, err := box.Find(name + ".json")
	if err != nil {
		return nil, errors.Wrap(err, "unable to read contract JSON")
	}

	abiBytes := gjson.GetBytes(jsonFile, "compilerOutput.abi")
	abiParsed, err := abi.JSON(strings.NewReader(abiBytes.Raw))
	if err != nil {
		return nil, err
	}

	return &Contract{abiParsed}, nil
}

// GetContract loads the contract JSON file from ../evm/dist/artifacts
// and parses the ABI JSON contents into an abi.ABI object
//
// NB: These contracts can be built by running
//    yarn workspace chainlink run setup
// in the base project directory.
func GetContract(name string) (*Contract, error) {
	box := packr.NewBox("../../evm/dist/artifacts")
	return getContract(name, box)
}

// GetAdvancedContract loads the contract JSON file from ../evm/v0.5/dist/artifacts
// and parses the ABI JSON contents into an abi.ABI object
//
// NB: These contracts can be built by running
//    yarn workspace chainlinkv0.5 run setup
// in the base project directory.
func GetV5Contract(name string) (*Contract, error) {
	box := packr.NewBox("../../evm/v0.5/dist/artifacts")
	return getContract(name, box)
}

// EncodeMessageCall encodes method name and arguments into a byte array
// to conform with the contract's ABI
func (contract *Contract) EncodeMessageCall(method string, args ...interface{}) ([]byte, error) {
	return contract.ABI.Pack(method, args...)
}
