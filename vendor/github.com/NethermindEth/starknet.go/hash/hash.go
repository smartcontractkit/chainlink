package hash

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// ComputeHashOnElementsFelt computes the hash on elements of a Felt array.
//
// Parameters:
// - feltArr: A pointer to an array of Felt objects.
// Returns:
// - *felt.Felt: a pointer to a Felt object
// - error: an error if any
func ComputeHashOnElementsFelt(feltArr []*felt.Felt) (*felt.Felt, error) {
	bigIntArr := utils.FeltArrToBigIntArr(feltArr)

	hash, err := curve.Curve.ComputeHashOnElements(bigIntArr)
	if err != nil {
		return nil, err
	}
	return utils.BigIntToFelt(hash), nil
}

// CalculateTransactionHashCommon calculates the transaction hash common to be used in the StarkNet network - a unique identifier of the transaction.
// [specification]: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/core/os/transaction_hash/transaction_hash.py#L27C5-L27C38
//
// Parameters:
// - txHashPrefix: The prefix of the transaction hash
// - version: The version of the transaction
// - contractAddress: The address of the contract
// - entryPointSelector: The selector of the entry point
// - calldata: The data of the transaction
// - maxFee: The maximum fee for the transaction
// - chainId: The ID of the blockchain
// - additionalData: Additional data to be included in the hash
// Returns:
// - *felt.Felt: the calculated transaction hash
// - error: an error if any
func CalculateTransactionHashCommon(
	txHashPrefix *felt.Felt,
	version *felt.Felt,
	contractAddress *felt.Felt,
	entryPointSelector *felt.Felt,
	calldata *felt.Felt,
	maxFee *felt.Felt,
	chainId *felt.Felt,
	additionalData []*felt.Felt) (*felt.Felt, error) {

	dataToHash := []*felt.Felt{
		txHashPrefix,
		version,
		contractAddress,
		entryPointSelector,
		calldata,
		maxFee,
		chainId,
	}
	dataToHash = append(dataToHash, additionalData...)
	return ComputeHashOnElementsFelt(dataToHash)
}

// ClassHash calculates the hash of a contract class.
//
// It takes a contract class as input and calculates the hash by combining various elements of the class.
// The hash is calculated using the PoseidonArray function from the Curve package.
// The elements used in the hash calculation include the contract class version, constructor entry point, external entry point, L1 handler entry point, ABI, and Sierra program.
// The ABI is converted to bytes and then hashed using the StarknetKeccak function from the Curve package.
// Finally, the ContractClassVersionHash, ExternalHash, L1HandleHash, ConstructorHash, ABIHash, and SierraProgamHash are combined using the PoseidonArray function from the Curve package.
//
// Parameters:
// - contract: A contract class object of type rpc.ContractClass.
// Returns:
// - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
// - error: an error object if there was an error during the hash calculation.
func ClassHash(contract rpc.ContractClass) (*felt.Felt, error) {
	// https://docs.starknet.io/documentation/architecture_and_concepts/Smart_Contracts/class-hash/

	Version := "CONTRACT_CLASS_V" + contract.ContractClassVersion
	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte(Version))
	ConstructorHash := hashEntryPointByType(contract.EntryPointsByType.Constructor)
	ExternalHash := hashEntryPointByType(contract.EntryPointsByType.External)
	L1HandleHash := hashEntryPointByType(contract.EntryPointsByType.L1Handler)
	SierraProgamHash := curve.Curve.PoseidonArray(contract.SierraProgram...)
	ABIHash, err := curve.Curve.StarknetKeccak([]byte(contract.ABI))
	if err != nil {
		return nil, err
	}

	// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#deploy_account_hash_calculation
	return curve.Curve.PoseidonArray(ContractClassVersionHash, ExternalHash, L1HandleHash, ConstructorHash, ABIHash, SierraProgamHash), nil
}

// hashEntryPointByType calculates the hash of an entry point by type.
//
// Parameters:
// - entryPoint: A slice of rpc.SierraEntryPoint objects
// Returns:
// - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
func hashEntryPointByType(entryPoint []rpc.SierraEntryPoint) *felt.Felt {
	flattened := make([]*felt.Felt, 0, len(entryPoint))
	for _, elt := range entryPoint {
		flattened = append(flattened, elt.Selector, new(felt.Felt).SetUint64(uint64(elt.FunctionIdx)))
	}
	return curve.Curve.PoseidonArray(flattened...)
}

// CompiledClassHash calculates the hash of a compiled class in the Casm format.
//
// Parameters:
// - casmClass: A `contracts.CasmClass` object
// Returns:
// - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
func CompiledClassHash(casmClass contracts.CasmClass) *felt.Felt {
	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte("COMPILED_CLASS_V1"))
	ExternalHash := hashCasmClassEntryPointByType(casmClass.EntryPointByType.External)
	L1HandleHash := hashCasmClassEntryPointByType(casmClass.EntryPointByType.L1Handler)
	ConstructorHash := hashCasmClassEntryPointByType(casmClass.EntryPointByType.Constructor)
	ByteCodeHasH := curve.Curve.PoseidonArray(casmClass.ByteCode...)

	// https://github.com/software-mansion/starknet.py/blob/development/starknet_py/hash/casm_class_hash.py#L10
	return curve.Curve.PoseidonArray(ContractClassVersionHash, ExternalHash, L1HandleHash, ConstructorHash, ByteCodeHasH)
}

// hashCasmClassEntryPointByType calculates the hash of a CasmClassEntryPoint array.
//
// Parameters:
// - entryPoint: An array of CasmClassEntryPoint objects
// Returns:
// - *felt.Felt: a pointer to a Felt type
func hashCasmClassEntryPointByType(entryPoint []contracts.CasmClassEntryPoint) *felt.Felt {
	flattened := make([]*felt.Felt, 0, len(entryPoint))
	for _, elt := range entryPoint {
		builtInFlat := []*felt.Felt{}
		for _, builtIn := range elt.Builtins {
			builtInFlat = append(builtInFlat, new(felt.Felt).SetBytes([]byte(builtIn)))
		}
		builtInHash := curve.Curve.PoseidonArray(builtInFlat...)
		flattened = append(flattened, elt.Selector, new(felt.Felt).SetUint64(uint64(elt.Offset)), builtInHash)
	}
	return curve.Curve.PoseidonArray(flattened...)
}
