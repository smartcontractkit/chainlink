package ton

import "github.com/xssnick/tonutils-go/address"

type Ton2AnyMessage struct {
	Receiver      []byte
	Data          []byte
	TokenAmounts  []Ton2AnyTokenAmount
	FeeToken      address.Address
	FeeTokenStore address.Address
	ExtraArgs     []byte
}

type Ton2AnyTokenAmount struct {
	Token      address.Address
	Amount     uint64
	TokenStore address.Address
}
