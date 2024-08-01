package rpc

import "github.com/NethermindEth/juno/core/felt"

type BroadcastTxn interface{}

var (
	_ BroadcastTxn = BroadcastInvokev0Txn{}
	_ BroadcastTxn = BroadcastInvokev1Txn{}
	_ BroadcastTxn = BroadcastDeclareTxnV1{}
	_ BroadcastTxn = BroadcastDeclareTxnV2{}
	_ BroadcastTxn = BroadcastDeclareTxnV3{}
	_ BroadcastTxn = BroadcastDeployAccountTxn{}
)

type BroadcastInvokeTxnType interface{}

var (
	_ BroadcastInvokeTxnType = BroadcastInvokev0Txn{}
	_ BroadcastInvokeTxnType = BroadcastInvokev1Txn{}
	_ BroadcastInvokeTxnType = BroadcastInvokev3Txn{}
)

type BroadcastDeclareTxnType interface{}

var (
	_ BroadcastDeclareTxnType = BroadcastDeclareTxnV1{}
	_ BroadcastDeclareTxnType = BroadcastDeclareTxnV2{}
	_ BroadcastDeclareTxnType = BroadcastDeclareTxnV3{}
)

type BroadcastAddDeployTxnType interface{}

var (
	_ BroadcastAddDeployTxnType = BroadcastDeployAccountTxn{}
	_ BroadcastAddDeployTxnType = BroadcastDeployAccountTxnV3{}
)

type BroadcastInvokev0Txn struct {
	InvokeTxnV0
}

type BroadcastInvokev1Txn struct {
	InvokeTxnV1
}

type BroadcastInvokev3Txn struct {
	InvokeTxnV3
}

type BroadcastDeclareTxnV1 struct {
	Type TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt              `json:"sender_address"`
	MaxFee        *felt.Felt              `json:"max_fee"`
	Version       NumAsHex                `json:"version"`
	Signature     []*felt.Felt            `json:"signature"`
	Nonce         *felt.Felt              `json:"nonce"`
	ContractClass DeprecatedContractClass `json:"contract_class"`
}
type BroadcastDeclareTxnV2 struct {
	Type TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress     *felt.Felt    `json:"sender_address"`
	CompiledClassHash *felt.Felt    `json:"compiled_class_hash"`
	MaxFee            *felt.Felt    `json:"max_fee"`
	Version           NumAsHex      `json:"version"`
	Signature         []*felt.Felt  `json:"signature"`
	Nonce             *felt.Felt    `json:"nonce"`
	ContractClass     ContractClass `json:"contract_class"`
}

type BroadcastDeclareTxnV3 struct {
	Type              TransactionType       `json:"type"`
	SenderAddress     *felt.Felt            `json:"sender_address"`
	CompiledClassHash *felt.Felt            `json:"compiled_class_hash"`
	Version           NumAsHex              `json:"version"`
	Signature         []*felt.Felt          `json:"signature"`
	Nonce             *felt.Felt            `json:"nonce"`
	ContractClass     *ContractClass        `json:"contract_class"`
	ResourceBounds    ResourceBoundsMapping `json:"resource_bounds"`
	Tip               U64                   `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The data needed to deploy the account contract from which this tx will be initiated
	AccountDeploymentData *felt.Felt `json:"account_deployment_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode DataAvailabilityMode `json:"fee_data_availability_mode"`
}

type BroadcastDeployAccountTxn struct {
	DeployAccountTxn
}
type BroadcastDeployAccountTxnV3 struct {
	DeployAccountTxnV3
}
