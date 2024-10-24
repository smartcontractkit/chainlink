/*
Package evm provides the EVM relay service along with the necessary utilities to interact with the EVM chain.

Using the ChainReaderService for EMV:

Initialization requires a ChainReaderConfig. The config is expected to be a static json file (though toml is
partially supported) that provides ABI definitions, 'method' and 'event' mappings, event definitions, input and
and output modifications, and various other options.

The general flow of initialization is as follows:

  - For each 'contract' in the config, parse the provided ABI in the config
  - For each 'method' (generically a ReadType), optionally create a method or event reader
  - Set filters in LogPoller for all bindings.

Each method or event binding uses the parsed ABI to create codec instances for encoding and decoding a variety of
data both as incoming parameters and results from contract outputs (methods and logs).

example (config):

	{
		"contracts": {
			"NamedContract": {
				"contractABI": "[encoded_abi]", // string value containing the entire encoded ABI
				"configs": {
					"contractReadName1": {
						"chainSpecificName": "contractRead", // can be different from read name or the same
						"readType": "method", // "method" for contract methods, "event" for contract events
						"confidenceConfirmations": {
							"unconfirmed": 1, // modify confidence levels as needed
							"finalized": 12 // if the concept of 'finalized' needs to be adjusted for a chain
						}
					}
				}
			}
		}
	}

example (method):

	// solidity method on NamedContract
	function contractRead(address token, quantity uint16) external view returns (uint224) {
		// ... implementation
	}

	// ContractReader usage
	type ContractReadParameters struct {
		Address []byte // use chain agnostic address type
		Quantity uint16
	}

	parameters := ContractReadParameters{
		Address: []byte("0x2142"),
		Qunatity: 100,
	}

	var contractReadResult *big.Int // another chain agnostic type

	const readName = "contractReadName1"

	// the following BoundContract should be provided to Bind before calling GetLatestValue
	// this only needs to be done once per address
	identifier := types.BoundContract{
		Address: "0x4221",
		Name: "NamedContract",
	}.ReadIdentifier(readName)

	_ = reader.GetLatestValue(ctx, identifier, primitives.Finalized, parameters, contractReadResult)
*/
package evm
