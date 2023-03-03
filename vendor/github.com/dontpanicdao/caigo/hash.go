package caigo

import (
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo/types"
)

func fmtCalldataStrings(calls []types.FunctionCall) (calldataStrings []string) {
	callArray := fmtCalldata(calls)
	for _, data := range callArray {
		calldataStrings = append(calldataStrings, fmt.Sprintf("0x%s", data.Text(16)))
	}
	return calldataStrings
}

func fmtV0CalldataStrings(nonce *big.Int, calls []types.FunctionCall) (calldataStrings []string) {
	calldataStrings = fmtCalldataStrings(calls)
	calldataStrings = append(calldataStrings, fmt.Sprintf("0x%s", nonce.Text(16)))
	return calldataStrings
}

/*
Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func fmtCalldata(calls []types.FunctionCall) (calldataArray []*big.Int) {
	callArray := []*big.Int{big.NewInt(int64(len(calls)))}

	for _, tx := range calls {
		address, _ := big.NewInt(0).SetString(tx.ContractAddress.Hex(), 0)
		callArray = append(callArray, address, types.GetSelectorFromName(tx.EntryPointSelector))

		if len(tx.Calldata) == 0 {
			callArray = append(callArray, big.NewInt(0), big.NewInt(0))

			continue
		}

		callArray = append(callArray, big.NewInt(int64(len(calldataArray))), big.NewInt(int64(len(tx.Calldata))))
		for _, cd := range tx.Calldata {
			calldataArray = append(calldataArray, types.SNValToBN(cd))
		}
	}

	callArray = append(callArray, big.NewInt(int64(len(calldataArray))))
	callArray = append(callArray, calldataArray...)
	return callArray
}

/*
Formats the multicall transactions with v0 of OpenZeppelin contract
*/
func fmtV0Calldata(nonce *big.Int, calls []types.FunctionCall) (calldataArray []*big.Int) {
	callArray := fmtCalldata(calls)
	callArray = append(callArray, nonce)
	return callArray
}
