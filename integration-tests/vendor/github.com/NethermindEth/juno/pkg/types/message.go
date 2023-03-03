package types

type MessageL2ToL1 struct {
	ToAddress EthAddress
	Payload   []Felt
}

type MessageL1ToL2 struct {
	FromAddress EthAddress
	Payload     []Felt
}
