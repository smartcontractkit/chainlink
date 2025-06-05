package ton

import tonaddress "github.com/xssnick/tonutils-go/address"

type Ton2AnyMessage struct {
	Receiver      []byte
	Data          []byte
	TokenAmounts  []Ton2AnyTokenAmount
	FeeToken      tonaddress.Address
	FeeTokenStore tonaddress.Address
	ExtraArgs     []byte
}

type Ton2AnyTokenAmount struct {
	Token      tonaddress.Address
	Amount     uint64
	TokenStore tonaddress.Address
}
