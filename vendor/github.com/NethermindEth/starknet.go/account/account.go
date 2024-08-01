package account

import (
	"context"
	"errors"
	"time"

	"github.com/NethermindEth/juno/core/crypto"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

var (
	ErrNotAllParametersSet   = errors.New("Not all neccessary parameters have been set")
	ErrTxnTypeUnSupported    = errors.New("Unsupported transction type")
	ErrTxnVersionUnSupported = errors.New("Unsupported transction version")
	ErrFeltToBigInt          = errors.New("Felt to BigInt error")
)

var (
	PREFIX_TRANSACTION      = new(felt.Felt).SetBytes([]byte("invoke"))
	PREFIX_DECLARE          = new(felt.Felt).SetBytes([]byte("declare"))
	PREFIX_CONTRACT_ADDRESS = new(felt.Felt).SetBytes([]byte("STARKNET_CONTRACT_ADDRESS"))
	PREFIX_DEPLOY_ACCOUNT   = new(felt.Felt).SetBytes([]byte("deploy_account"))
)

//go:generate mockgen -destination=../mocks/mock_account.go -package=mocks -source=account.go AccountInterface
type AccountInterface interface {
	Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error)
	TransactionHashInvoke(invokeTxn rpc.InvokeTxnType) (*felt.Felt, error)
	TransactionHashDeployAccount(tx rpc.DeployAccountType, contractAddress *felt.Felt) (*felt.Felt, error)
	TransactionHashDeclare(tx rpc.DeclareTxnType) (*felt.Felt, error)
	SignInvokeTransaction(ctx context.Context, tx *rpc.InvokeTxnV1) error
	SignDeployAccountTransaction(ctx context.Context, tx *rpc.DeployAccountTxn, precomputeAddress *felt.Felt) error
	SignDeclareTransaction(ctx context.Context, tx *rpc.DeclareTxnV2) error
	PrecomputeAddress(deployerAddress *felt.Felt, salt *felt.Felt, classHash *felt.Felt, constructorCalldata []*felt.Felt) (*felt.Felt, error)
	WaitForTransactionReceipt(ctx context.Context, transactionHash *felt.Felt, pollInterval time.Duration) (*rpc.TransactionReceiptWithBlockInfo, error)
}

var _ AccountInterface = &Account{}
var _ rpc.RpcProvider = &Account{}

type Account struct {
	provider       rpc.RpcProvider
	ChainId        *felt.Felt
	AccountAddress *felt.Felt
	publicKey      string
	CairoVersion   int
	ks             Keystore
}

// NewAccount creates a new Account instance.
//
// Parameters:
// - provider: is the provider of type rpc.RpcProvider
// - accountAddress: is the account address of type *felt.Felt
// - publicKey: is the public key of type string
// - keystore: is the keystore of type Keystore
// It returns:
// - *Account: a pointer to newly created Account
// - error: an error if any
func NewAccount(provider rpc.RpcProvider, accountAddress *felt.Felt, publicKey string, keystore Keystore, cairoVersion int) (*Account, error) {
	account := &Account{
		provider:       provider,
		AccountAddress: accountAddress,
		publicKey:      publicKey,
		ks:             keystore,
		CairoVersion:   cairoVersion,
	}

	chainID, err := provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	account.ChainId = new(felt.Felt).SetBytes([]byte(chainID))

	return account, nil
}

// Sign signs the given felt message using the account's private key.
//
// Parameters:
// - ctx: is the context used for the signing operation
// - msg: is the felt message to be signed
// Returns:
// - []*felt.Felt: an array of signed felt messages
// - error: an error, if any
func (account *Account) Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error) {

	msgBig := utils.FeltToBigInt(msg)

	s1, s2, err := account.ks.Sign(ctx, account.publicKey, msgBig)
	if err != nil {
		return nil, err
	}
	s1Felt := utils.BigIntToFelt(s1)
	s2Felt := utils.BigIntToFelt(s2)

	return []*felt.Felt{s1Felt, s2Felt}, nil
}

// SignInvokeTransaction signs and invokes a transaction.
//
// Parameters:
// - ctx: the context.Context for the function execution.
// - invokeTx: the InvokeTxnV1 struct representing the transaction to be invoked.
// Returns:
// - error: an error if there was an error in the signing or invoking process
func (account *Account) SignInvokeTransaction(ctx context.Context, invokeTx *rpc.InvokeTxnV1) error {

	txHash, err := account.TransactionHashInvoke(*invokeTx)
	if err != nil {
		return err
	}
	signature, err := account.Sign(ctx, txHash)
	if err != nil {
		return err
	}
	invokeTx.Signature = signature
	return nil
}

// SignDeployAccountTransaction signs a deploy account transaction.
//
// Parameters:
// - ctx: the context.Context for the function execution
// - tx: the *rpc.DeployAccountTxn struct representing the transaction to be signed
// - precomputeAddress: the precomputed address for the transaction
// Returns:
// - error: an error if any
func (account *Account) SignDeployAccountTransaction(ctx context.Context, tx *rpc.DeployAccountTxn, precomputeAddress *felt.Felt) error {

	hash, err := account.TransactionHashDeployAccount(*tx, precomputeAddress)
	if err != nil {
		return err
	}
	signature, err := account.Sign(ctx, hash)
	if err != nil {
		return err
	}
	tx.Signature = signature
	return nil
}

// SignDeclareTransaction signs a DeclareTxnV2 transaction using the provided Account.
//
// Parameters:
// - ctx: the context.Context
// - tx: the *rpc.DeclareTxnV2
// Returns:
// - error: an error if any
func (account *Account) SignDeclareTransaction(ctx context.Context, tx *rpc.DeclareTxnV2) error {

	hash, err := account.TransactionHashDeclare(*tx)
	if err != nil {
		return err
	}
	signature, err := account.Sign(ctx, hash)
	if err != nil {
		return err
	}
	tx.Signature = signature
	return nil
}

// TransactionHashDeployAccount calculates the transaction hash for a deploy account transaction.
//
// Parameters:
// - tx: The deploy account transaction to calculate the hash for
// - contractAddress: The contract address as parameters as a *felt.Felt
// Returns:
// - *felt.Felt: the calculated transaction hash
// - error: an error if any
func (account *Account) TransactionHashDeployAccount(tx rpc.DeployAccountType, contractAddress *felt.Felt) (*felt.Felt, error) {

	// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#deploy_account_transaction
	switch txn := tx.(type) {
	case rpc.DeployAccountTxn:
		calldata := []*felt.Felt{txn.ClassHash, txn.ContractAddressSalt}
		calldata = append(calldata, txn.ConstructorCalldata...)
		calldataHash, err := hash.ComputeHashOnElementsFelt(calldata)
		if err != nil {
			return nil, err
		}

		versionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}

		// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#deploy_account_hash_calculation
		return hash.CalculateTransactionHashCommon(
			PREFIX_DEPLOY_ACCOUNT,
			versionFelt,
			contractAddress,
			&felt.Zero,
			calldataHash,
			txn.MaxFee,
			account.ChainId,
			[]*felt.Felt{txn.Nonce},
		)
	case rpc.DeployAccountTxnV3:
		if txn.Version == "" || txn.ResourceBounds == (rpc.ResourceBoundsMapping{}) || txn.Nonce == nil || txn.PayMasterData == nil {
			return nil, ErrNotAllParametersSet
		}
		calldata := []*felt.Felt{txn.ClassHash, txn.ContractAddressSalt}
		calldata = append(calldata, txn.ConstructorCalldata...)

		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		DAUint64, err := dataAvailabilityMode(txn.FeeMode, txn.NonceDataMode)
		if err != nil {
			return nil, err
		}
		tipUint64, err := txn.Tip.ToUint64()
		if err != nil {
			return nil, err
		}
		tipAndResourceHash, err := tipAndResourcesHash(tipUint64, txn.ResourceBounds)
		if err != nil {
			return nil, err
		}
		// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#deploy_account_hash_calculation
		return crypto.PoseidonArray(
			PREFIX_DEPLOY_ACCOUNT,
			txnVersionFelt,
			contractAddress,
			tipAndResourceHash,
			crypto.PoseidonArray(txn.PayMasterData...),
			account.ChainId,
			txn.Nonce,
			new(felt.Felt).SetUint64(DAUint64),
			crypto.PoseidonArray(txn.ConstructorCalldata...),
			txn.ClassHash,
			txn.ContractAddressSalt,
		), nil
	}
	return nil, ErrTxnTypeUnSupported
}

// TransactionHashInvoke calculates the transaction hash for the given invoke transaction.
//
// Parameters:
// - tx: The invoke transaction to calculate the hash for.
//     The transaction can be of type InvokeTxnV0 or InvokeTxnV1.
//     For InvokeTxnV0:
//         the function checks if all the required parameters are set and then computes the transaction hash using the provided data.
//     For InvokeTxnV1:
//         the function performs similar checks and computes the transaction hash using the provided data.
// Returns:
// - *felt.Felt: The calculated transaction hash as a *felt.Felt
// - error: an error, if any

// If the transaction type is unsupported, the function returns an error.
func (account *Account) TransactionHashInvoke(tx rpc.InvokeTxnType) (*felt.Felt, error) {

	// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#v0_hash_calculation
	switch txn := tx.(type) {
	case rpc.InvokeTxnV0:
		if txn.Version == "" || len(txn.Calldata) == 0 || txn.MaxFee == nil || txn.EntryPointSelector == nil {
			return nil, ErrNotAllParametersSet
		}

		calldataHash, err := hash.ComputeHashOnElementsFelt(txn.Calldata)
		if err != nil {
			return nil, err
		}

		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		return hash.CalculateTransactionHashCommon(
			PREFIX_TRANSACTION,
			txnVersionFelt,
			txn.ContractAddress,
			txn.EntryPointSelector,
			calldataHash,
			txn.MaxFee,
			account.ChainId,
			[]*felt.Felt{},
		)

	case rpc.InvokeTxnV1:
		if txn.Version == "" || len(txn.Calldata) == 0 || txn.Nonce == nil || txn.MaxFee == nil || txn.SenderAddress == nil {
			return nil, ErrNotAllParametersSet
		}

		calldataHash, err := hash.ComputeHashOnElementsFelt(txn.Calldata)
		if err != nil {
			return nil, err
		}
		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		return hash.CalculateTransactionHashCommon(
			PREFIX_TRANSACTION,
			txnVersionFelt,
			txn.SenderAddress,
			&felt.Zero,
			calldataHash,
			txn.MaxFee,
			account.ChainId,
			[]*felt.Felt{txn.Nonce},
		)
	case rpc.InvokeTxnV3:
		// https://github.com/starknet-io/SNIPs/blob/main/SNIPS/snip-8.md#protocol-changes
		if txn.Version == "" || txn.ResourceBounds == (rpc.ResourceBoundsMapping{}) || len(txn.Calldata) == 0 || txn.Nonce == nil || txn.SenderAddress == nil || txn.PayMasterData == nil || txn.AccountDeploymentData == nil {
			return nil, ErrNotAllParametersSet
		}

		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		DAUint64, err := dataAvailabilityMode(txn.FeeMode, txn.NonceDataMode)
		if err != nil {
			return nil, err
		}
		tipUint64, err := txn.Tip.ToUint64()
		if err != nil {
			return nil, err
		}
		tipAndResourceHash, err := tipAndResourcesHash(tipUint64, txn.ResourceBounds)
		if err != nil {
			return nil, err
		}
		return crypto.PoseidonArray(
			PREFIX_TRANSACTION,
			txnVersionFelt,
			txn.SenderAddress,
			tipAndResourceHash,
			crypto.PoseidonArray(txn.PayMasterData...),
			account.ChainId,
			txn.Nonce,
			new(felt.Felt).SetUint64(DAUint64),
			crypto.PoseidonArray(txn.AccountDeploymentData...),
			crypto.PoseidonArray(txn.Calldata...),
		), nil
	}
	return nil, ErrTxnTypeUnSupported
}

func tipAndResourcesHash(tip uint64, resourceBounds rpc.ResourceBoundsMapping) (*felt.Felt, error) {
	l1Bytes, err := resourceBounds.L1Gas.Bytes(rpc.ResourceL1Gas)
	if err != nil {
		return nil, err
	}
	l2Bytes, err := resourceBounds.L2Gas.Bytes(rpc.ResourceL2Gas)
	if err != nil {
		return nil, err
	}
	l1Bounds := new(felt.Felt).SetBytes(l1Bytes)
	l2Bounds := new(felt.Felt).SetBytes(l2Bytes)
	return crypto.PoseidonArray(new(felt.Felt).SetUint64(tip), l1Bounds, l2Bounds), nil
}

func dataAvailabilityMode(feeDAMode, nonceDAMode rpc.DataAvailabilityMode) (uint64, error) {
	const dataAvailabilityModeBits = 32
	fee64, err := feeDAMode.UInt64()
	if err != nil {
		return 0, err
	}
	nonce64, err := nonceDAMode.UInt64()
	if err != nil {
		return 0, err
	}
	return fee64 + nonce64<<dataAvailabilityModeBits, nil
}

// TransactionHashDeclare calculates the transaction hash for declaring a transaction type.
//
// Parameters:
// - tx: The `tx` parameter of type `rpc.DeclareTxnType`
// Can be one of the following types:
//   - `rpc.DeclareTxnV0`
//   - `rpc.DeclareTxnV1`
//   - `rpc.DeclareTxnV2`
//
// Returns:
// - *felt.Felt: the calculated transaction hash as `*felt.Felt` value
// - error: an error, if any
//
// If the `tx` parameter is not one of the supported types, the function returns an error `ErrTxnTypeUnSupported`.
func (account *Account) TransactionHashDeclare(tx rpc.DeclareTxnType) (*felt.Felt, error) {

	switch txn := tx.(type) {
	case rpc.DeclareTxnV0:
		// Due to inconsistencies in version 0 hash calculation we don't calculate the hash
		return nil, ErrTxnVersionUnSupported
	case rpc.DeclareTxnV1:
		if txn.SenderAddress == nil || txn.Version == "" || txn.ClassHash == nil || txn.MaxFee == nil || txn.Nonce == nil {
			return nil, ErrNotAllParametersSet
		}

		calldataHash, err := hash.ComputeHashOnElementsFelt([]*felt.Felt{txn.ClassHash})
		if err != nil {
			return nil, err
		}

		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		return hash.CalculateTransactionHashCommon(
			PREFIX_DECLARE,
			txnVersionFelt,
			txn.SenderAddress,
			&felt.Zero,
			calldataHash,
			txn.MaxFee,
			account.ChainId,
			[]*felt.Felt{txn.Nonce},
		)
	case rpc.DeclareTxnV2:
		if txn.CompiledClassHash == nil || txn.SenderAddress == nil || txn.Version == "" || txn.ClassHash == nil || txn.MaxFee == nil || txn.Nonce == nil {
			return nil, ErrNotAllParametersSet
		}

		calldataHash, err := hash.ComputeHashOnElementsFelt([]*felt.Felt{txn.ClassHash})
		if err != nil {
			return nil, err
		}

		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		return hash.CalculateTransactionHashCommon(
			PREFIX_DECLARE,
			txnVersionFelt,
			txn.SenderAddress,
			&felt.Zero,
			calldataHash,
			txn.MaxFee,
			account.ChainId,
			[]*felt.Felt{txn.Nonce, txn.CompiledClassHash},
		)
	case rpc.DeclareTxnV3:
		// https://github.com/starknet-io/SNIPs/blob/main/SNIPS/snip-8.md#protocol-changes
		if txn.Version == "" || txn.ResourceBounds == (rpc.ResourceBoundsMapping{}) || txn.Nonce == nil || txn.SenderAddress == nil || txn.PayMasterData == nil || txn.AccountDeploymentData == nil ||
			txn.ClassHash == nil || txn.CompiledClassHash == nil {
			return nil, ErrNotAllParametersSet
		}

		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		DAUint64, err := dataAvailabilityMode(txn.FeeMode, txn.NonceDataMode)
		if err != nil {
			return nil, err
		}
		tipUint64, err := txn.Tip.ToUint64()
		if err != nil {
			return nil, err
		}

		tipAndResourceHash, err := tipAndResourcesHash(tipUint64, txn.ResourceBounds)
		if err != nil {
			return nil, err
		}
		return crypto.PoseidonArray(
			PREFIX_DECLARE,
			txnVersionFelt,
			txn.SenderAddress,
			tipAndResourceHash,
			crypto.PoseidonArray(txn.PayMasterData...),
			account.ChainId,
			txn.Nonce,
			new(felt.Felt).SetUint64(DAUint64),
			crypto.PoseidonArray(txn.AccountDeploymentData...),
			txn.ClassHash,
			txn.CompiledClassHash,
		), nil
	}

	return nil, ErrTxnTypeUnSupported
}

// PrecomputeAddress calculates the precomputed address for an account.
// ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/core/os/contract_address/contract_address.py
//
// Parameters:
// - deployerAddress: the deployer address
// - salt: the salt
// - classHash: the class hash
// - constructorCalldata: the constructor calldata
// Returns:
// - *felt.Felt: the precomputed address as a *felt.Felt
// - error: an error if any
func (account *Account) PrecomputeAddress(deployerAddress *felt.Felt, salt *felt.Felt, classHash *felt.Felt, constructorCalldata []*felt.Felt) (*felt.Felt, error) {

	bigIntArr := utils.FeltArrToBigIntArr([]*felt.Felt{
		PREFIX_CONTRACT_ADDRESS,
		deployerAddress,
		salt,
		classHash,
	})

	constructorCalldataBigIntArr := utils.FeltArrToBigIntArr(constructorCalldata)
	constructorCallDataHashInt, _ := curve.Curve.ComputeHashOnElements(constructorCalldataBigIntArr)
	bigIntArr = append(bigIntArr, constructorCallDataHashInt)

	preBigInt, err := curve.Curve.ComputeHashOnElements(bigIntArr)
	if err != nil {
		return nil, err
	}
	return utils.BigIntToFelt(preBigInt), nil

}

// WaitForTransactionReceipt waits for the transaction receipt of the given transaction hash to succeed or fail.
//
// Parameters:
// - ctx: The context
// - transactionHash: The hash
// - pollInterval: The poll interval as parameters
// It returns:
// - *rpc.TransactionReceipt: the transaction receipt
// - error: an error
func (account *Account) WaitForTransactionReceipt(ctx context.Context, transactionHash *felt.Felt, pollInterval time.Duration) (*rpc.TransactionReceiptWithBlockInfo, error) {
	t := time.NewTicker(pollInterval)
	for {
		select {
		case <-ctx.Done():
			return nil, rpc.Err(rpc.InternalError, ctx.Err())
		case <-t.C:
			receiptWithBlockInfo, err := account.TransactionReceipt(ctx, transactionHash)
			if err != nil {
				if errors.Is(err, rpc.ErrHashNotFound) {
					continue
				} else {
					return nil, err
				}
			}
			return receiptWithBlockInfo, nil
		}
	}
}

// AddInvokeTransaction generates an invoke transaction and adds it to the account's provider.
//
// Parameters:
// - ctx: the context.Context object for the transaction.
// - invokeTx: the invoke transaction to be added.
// Returns:
// - *rpc.AddInvokeTransactionResponse: The response for the AddInvokeTransactionResponse
// - error: an error if any.
func (account *Account) AddInvokeTransaction(ctx context.Context, invokeTx rpc.BroadcastInvokeTxnType) (*rpc.AddInvokeTransactionResponse, error) {
	return account.provider.AddInvokeTransaction(ctx, invokeTx)
}

// AddDeclareTransaction adds a declare transaction to the account.
//
// Parameters:
// - ctx: The context.Context for the request.
// - declareTransaction: The input for adding a declare transaction.
// Returns:
// - *rpc.AddDeclareTransactionResponse: The response for adding a declare transaction
// - error: an error, if any
func (account *Account) AddDeclareTransaction(ctx context.Context, declareTransaction rpc.BroadcastDeclareTxnType) (*rpc.AddDeclareTransactionResponse, error) {
	return account.provider.AddDeclareTransaction(ctx, declareTransaction)
}

// AddDeployAccountTransaction adds a deploy account transaction to the account.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - deployAccountTransaction: The rpc.DeployAccountTxn object representing the deploy account transaction.
// Returns:
// - *rpc.AddDeployAccountTransactionResponse: a pointer to rpc.AddDeployAccountTransactionResponse
// - error: an error if any
func (account *Account) AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction rpc.BroadcastAddDeployTxnType) (*rpc.AddDeployAccountTransactionResponse, error) {
	return account.provider.AddDeployAccountTransaction(ctx, deployAccountTransaction)
}

// BlockHashAndNumber returns the block hash and number for the account.
//
// Parameters:
// - ctx: The context in which the function is called.
// Returns:
// - rpc.BlockHashAndNumberOutput: the block hash and number as an rpc.BlockHashAndNumberOutput object.
// - error: an error if there was an issue retrieving the block hash and number.
func (account *Account) BlockHashAndNumber(ctx context.Context) (*rpc.BlockHashAndNumberOutput, error) {
	return account.provider.BlockHashAndNumber(ctx)
}

// BlockNumber returns the block number of the account.
//
// Parameters:
// - ctx: The context in which the function is called.
// Returns:
// - uint64: the block number as a uint64
// - error: an error encountered
func (account *Account) BlockNumber(ctx context.Context) (uint64, error) {
	return account.provider.BlockNumber(ctx)
}

// BlockTransactionCount returns the number of transactions in a block.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - blockID: The rpc.BlockID object representing the block.
// Returns:
// - uint64: the number of transactions in the block
//   - error: an error, if any
func (account *Account) BlockTransactionCount(ctx context.Context, blockID rpc.BlockID) (uint64, error) {
	return account.provider.BlockTransactionCount(ctx, blockID)
}

// BlockWithTxHashes retrieves a block with transaction hashes.
//
// Parameters:
// - ctx: the context.Context object for the request.
// - blockID: the rpc.BlockID object specifying the block to retrieve.
// Returns:
// - interface{}: an interface{} representing the retrieved block
// - error: an error if there was any issue retrieving the block
func (account *Account) BlockWithTxHashes(ctx context.Context, blockID rpc.BlockID) (interface{}, error) {
	return account.provider.BlockWithTxHashes(ctx, blockID)
}

// BlockWithTxs retrieves the specified block along with its transactions.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - blockID: The rpc.BlockID parameter for the function.
// Returns:
// - interface{}: An interface{}
// - error: An error
func (account *Account) BlockWithTxs(ctx context.Context, blockID rpc.BlockID) (interface{}, error) {
	return account.provider.BlockWithTxs(ctx, blockID)
}

func (account *Account) BlockWithReceipts(ctx context.Context, blockID rpc.BlockID) (interface{}, error) {
	return account.provider.BlockWithReceipts(ctx, blockID)
}

// Call is a function that performs a function call on an Account.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - call: The rpc.FunctionCall object representing the function call.
// - blockID: The rpc.BlockID object representing the block ID.
// Returns:
// - []*felt.Felt: a slice of *felt.Felt
// - error: an error object.
func (account *Account) Call(ctx context.Context, call rpc.FunctionCall, blockId rpc.BlockID) ([]*felt.Felt, error) {
	return account.provider.Call(ctx, call, blockId)
}

// ChainID returns the chain ID associated with the account.
//
// Parameters:
// - ctx: the context.Context object for the function.
// Returns:
//   - string: the chain ID.
//   - error: any error encountered while retrieving the chain ID.
func (account *Account) ChainID(ctx context.Context) (string, error) {
	return account.provider.ChainID(ctx)
}

// Class is a method that calls the `Class` method of the `provider` field of the `account` struct.
//
// Parameters:
// - ctx: The context.Context
// - blockID: The rpc.BlockID
// - classHash: The `*felt.Felt`
// Returns:
//   - *rpc.ClassOutput: The rpc.ClassOutput (the class output could be a DeprecatedContractClass
//     or just a Contract class depending on the contract version)
//   - error: An error if any occurred.
func (account *Account) Class(ctx context.Context, blockID rpc.BlockID, classHash *felt.Felt) (rpc.ClassOutput, error) {
	return account.provider.Class(ctx, blockID, classHash)
}

// ClassAt retrieves the class at the specified block ID and contract address.
// Parameters:
// - ctx: The context.Context object for the function.
// - blockID: The rpc.BlockID object representing the block ID.
// - contractAddress: The felt.Felt object representing the contract address.
// Returns:
//   - *rpc.ClassOutput: The rpc.ClassOutput object (the class output could be a DeprecatedContractClass
//     or just a Contract class depending on the contract version)
//   - error: An error if any occurred.
func (account *Account) ClassAt(ctx context.Context, blockID rpc.BlockID, contractAddress *felt.Felt) (rpc.ClassOutput, error) {
	return account.provider.ClassAt(ctx, blockID, contractAddress)
}

// ClassHashAt returns the class hash at the given block ID for the specified contract address.
//
// Parameters:
// - ctx: The context to use for the function call.
// - blockID: The ID of the block.
// contractAddress - The address of the contract to get the class hash for.
// Returns:
// - *felt.Felt: the class hash as a *felt.Felt
// - error: an error if any occurred.
func (account *Account) ClassHashAt(ctx context.Context, blockID rpc.BlockID, contractAddress *felt.Felt) (*felt.Felt, error) {
	return account.provider.ClassHashAt(ctx, blockID, contractAddress)
}

// EstimateFee estimates the fee for a set of requests in the given block ID.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - requests: An array of rpc.BroadcastTxn objects representing the requests to estimate the fee for.
// - blockID: The rpc.BlockID object representing the block ID for which to estimate the fee.
// Returns:
// - []rpc.FeeEstimate: An array of rpc.FeeEstimate objects representing the estimated fees.
// - error: An error object if any error occurred during the estimation process.
func (account *Account) EstimateFee(ctx context.Context, requests []rpc.BroadcastTxn, simulationFlags []rpc.SimulationFlag, blockID rpc.BlockID) ([]rpc.FeeEstimate, error) {
	return account.provider.EstimateFee(ctx, requests, simulationFlags, blockID)
}

// EstimateMessageFee estimates the fee for a given message in the context of an account.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - msg: The rpc.MsgFromL1 object representing the message.
// - blockID: The rpc.BlockID object representing the block ID.
// Returns:
// - *rpc.FeeEstimate: a pointer to rpc.FeeEstimate
// - error: an error if any.
func (account *Account) EstimateMessageFee(ctx context.Context, msg rpc.MsgFromL1, blockID rpc.BlockID) (*rpc.FeeEstimate, error) {
	return account.provider.EstimateMessageFee(ctx, msg, blockID)
}

// Events retrieves events for the account.
//
// Parameters:
// - ctx: the context.Context to use for the request.
// - input: the input parameters for retrieving events.
// Returns:
// - *rpc.EventChunk: the chunk of events retrieved.
// - error: an error if the retrieval fails.
func (account *Account) Events(ctx context.Context, input rpc.EventsInput) (*rpc.EventChunk, error) {
	return account.provider.Events(ctx, input)
}

// Nonce retrieves the nonce for a given block ID and contract address.
//
// Parameters:
// - ctx: is the context.Context for the function call
// - blockID: is the ID of the block
// - contractAddress: is the address of the contract
// Returns:
// - *felt.Felt: the contract's nonce at the requested state
// - error: an error if any
func (account *Account) Nonce(ctx context.Context, blockID rpc.BlockID, contractAddress *felt.Felt) (*felt.Felt, error) {
	return account.provider.Nonce(ctx, blockID, contractAddress)
}

// SimulateTransactions simulates transactions using the provided context
// Parameters:
// - ctx: The context.Context object
// - blockID: The rpc.BlockID object for the block referencing the state or call the transactions are on
// - txns: The slice of rpc.Transaction objects representing the transactions to simulate
// - simulationFlags: The slice of rpc.simulationFlags
// Returns:
// - []rpc.SimulatedTransaction: a list of simulated transactions
// - error: an error, if any.
func (account *Account) SimulateTransactions(ctx context.Context, blockID rpc.BlockID, txns []rpc.Transaction, simulationFlags []rpc.SimulationFlag) ([]rpc.SimulatedTransaction, error) {
	return account.provider.SimulateTransactions(ctx, blockID, txns, simulationFlags)
}

// StorageAt is a function that retrieves the storage value at the given key for a contract address.
//
// Parameters:
// - ctx: The context.Context object for the function
// - contractAddress: The contract address for which to retrieve the storage value
// - key: The key of the storage value to retrieve
// - blockID: The block ID at which to retrieve the storage value
// Returns:
// - string: The storage value at the given key.
// - error: An error if the retrieval fails.
func (account *Account) StorageAt(ctx context.Context, contractAddress *felt.Felt, key string, blockID rpc.BlockID) (string, error) {
	return account.provider.StorageAt(ctx, contractAddress, key, blockID)
}

// StateUpdate updates the state of the Account.
//
// Parameters:
// - context.Context: The context.Context object.
// - blockID: The rpc.BlockID object representing the block to update the state of.
// Returns:
// - *rpc.StateUpdateOutput: a *rpc.StateUpdateOutput
// - error: an error
func (account *Account) StateUpdate(ctx context.Context, blockID rpc.BlockID) (*rpc.StateUpdateOutput, error) {
	return account.provider.StateUpdate(ctx, blockID)
}

// SpecVersion returns the spec version of the account.
// It takes a context as a parameter and returns a string and an error
//
// Parameters:
// - context.Context: The context.Context object
// Returns:
// - string: The spec version
// - error: An error if any
func (account *Account) SpecVersion(ctx context.Context) (string, error) {
	return account.provider.SpecVersion(ctx)
}

// Syncing returns the sync status of the account.
//
// Parameters:
// - ctx: The context.Context object
// Returns:
// - *rpc.SyncStatus: *rpc.SyncStatus
// - error: an error.
func (account *Account) Syncing(ctx context.Context) (*rpc.SyncStatus, error) {
	return account.provider.Syncing(ctx)
}

// TraceBlockTransactions retrieves a list of trace transactions for a given block hash.
//
// Parameters:
// - ctx: The context.Context object.
// - blockID: The hash of the block to retrieve trace transactions for.
// Returns
// - []rpc.Trace: The list of trace transactions for the given block.
// - error: An error if there was a problem retrieving the trace transactions.
func (account *Account) TraceBlockTransactions(ctx context.Context, blockID rpc.BlockID) ([]rpc.Trace, error) {
	return account.provider.TraceBlockTransactions(ctx, blockID)
}

// TransactionReceipt retrieves the transaction receipt for the given transaction hash.
//
// Parameters:
// - ctx: The context to use for the request.
// - transactionHash: The hash of the transaction.
// Returns:
// - rpc.Transactiontype: rpc.TransactionReceipt, error.
func (account *Account) TransactionReceipt(ctx context.Context, transactionHash *felt.Felt) (*rpc.TransactionReceiptWithBlockInfo, error) {
	return account.provider.TransactionReceipt(ctx, transactionHash)
}

// TransactionTrace returns the transaction trace for a given transaction hash.
//
// Parameters:
// - ctx: The context.Context object for the request.
// - transactionHash: The transaction hash for which the transaction trace is to be retrieved.
// Returns:
// - rpc.TxnTrace: The rpc.TxnTrace object representing the transaction trace, and an error if any.
func (account *Account) TraceTransaction(ctx context.Context, transactionHash *felt.Felt) (rpc.TxnTrace, error) {
	return account.provider.TraceTransaction(ctx, transactionHash)
}

// TransactionByBlockIdAndIndex returns a transaction by block ID and index.
//
// Parameters:
// - ctx: The context for the function.
// - blockID: The ID of the block.
// - index: The index of the transaction in the block.
// Returns:
// - rpc.Transaction: The transaction and an error, if any.
func (account *Account) TransactionByBlockIdAndIndex(ctx context.Context, blockID rpc.BlockID, index uint64) (rpc.Transaction, error) {
	return account.provider.TransactionByBlockIdAndIndex(ctx, blockID, index)
}

// TransactionByHash returns the transaction with the given hash.
//
// Parameters:
// - ctx: The context.Context
// - hash: The *felt.Felt hash as parameters.
// Returns:
// - rpc.Transaction
// - error
func (account *Account) TransactionByHash(ctx context.Context, hash *felt.Felt) (rpc.Transaction, error) {
	return account.provider.TransactionByHash(ctx, hash)
}

// GetTransactionStatus returns the transaction status.
//
// Parameters:
// - ctx: The context.Context
// - Txnhash: The *felt.Felt Txn hash.
// Returns:
// - *rpc.TxnStatusResp: the transaction status
// - error: anerror if any
func (account *Account) GetTransactionStatus(ctx context.Context, Txnhash *felt.Felt) (*rpc.TxnStatusResp, error) {
	return account.provider.GetTransactionStatus(ctx, Txnhash)
}

// FmtCalldata generates the formatted calldata for the given function calls and Cairo version.
//
// Parameters:
// - fnCalls: a slice of rpc.FunctionCall representing the function calls.
// - cairoVersion: an integer representing the Cairo version.
// Returns:
// - a slice of *felt.Felt representing the formatted calldata.
// - an error if Cairo version is not supported.
func (account *Account) FmtCalldata(fnCalls []rpc.FunctionCall) ([]*felt.Felt, error) {
	switch account.CairoVersion {
	case 0:
		return FmtCallDataCairo0(fnCalls), nil
	case 2:
		return FmtCallDataCairo2(fnCalls), nil
	default:
		return nil, errors.New("Cairo version not supported")
	}
}

// FmtCallDataCairo0 generates a slice of *felt.Felt that represents the calldata for the given function calls in Cairo 0 format.
//
// Parameters:
// - fnCalls: a slice of rpc.FunctionCall containing the function calls.
//
// Returns:
// - a slice of *felt.Felt representing the generated calldata.
// https://github.com/project3fusion/StarkSharp/blob/main/StarkSharp/StarkSharp.Rpc/Modules/Transactions/Hash/TransactionHash.cs#L27
func FmtCallDataCairo0(callArray []rpc.FunctionCall) []*felt.Felt {
	var calldata []*felt.Felt
	var calls []*felt.Felt

	calldata = append(calldata, new(felt.Felt).SetUint64(uint64(len(callArray))))

	offset := uint64(0)
	for _, call := range callArray {
		calldata = append(calldata, call.ContractAddress)
		calldata = append(calldata, call.EntryPointSelector)
		calldata = append(calldata, new(felt.Felt).SetUint64(uint64(offset)))
		callDataLen := uint64(len(call.Calldata))
		calldata = append(calldata, new(felt.Felt).SetUint64(callDataLen))
		offset += callDataLen

		for _, data := range call.Calldata {
			calls = append(calls, data)
		}
	}

	calldata = append(calldata, new(felt.Felt).SetUint64(offset))
	calldata = append(calldata, calls...)

	return calldata
}

// FmtCallDataCairo2 generates the calldata for the given function calls for Cairo 2 contracs.
//
// Parameters:
// - fnCalls: a slice of rpc.FunctionCall containing the function calls.
// Returns:
// - a slice of *felt.Felt representing the generated calldata.
// https://github.com/project3fusion/StarkSharp/blob/main/StarkSharp/StarkSharp.Rpc/Modules/Transactions/Hash/TransactionHash.cs#L22
func FmtCallDataCairo2(callArray []rpc.FunctionCall) []*felt.Felt {
	var result []*felt.Felt

	result = append(result, new(felt.Felt).SetUint64(uint64(len(callArray))))

	for _, call := range callArray {
		result = append(result, call.ContractAddress)
		result = append(result, call.EntryPointSelector)

		callDataLen := uint64(len(call.Calldata))
		result = append(result, new(felt.Felt).SetUint64(callDataLen))

		result = append(result, call.Calldata...)
	}

	return result
}
