package ton

import "github.com/xssnick/tonutils-go/address"

type TVM2AnyMessage struct {
	Receiver      []byte
	Data          []byte
	TokenAmounts  []TVM2AnyTokenAmount
	FeeToken      address.Address
	FeeTokenStore address.Address
	ExtraArgs     []byte
}

type TVM2AnyTokenAmount struct {
	Token      address.Address
	Amount     uint64
	TokenStore address.Address
}
