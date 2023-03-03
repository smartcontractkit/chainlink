package types

import (
	"encoding/json"
)

type Address Felt

func BytesToAddress(b []byte) Address {
	return Address(BytesToFelt(b))
}

func HexToAddress(s string) Address {
	return Address(HexToFelt(s))
}

func (a *Address) Felt() Felt {
	return Felt(*a)
}

func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Felt())
}

func (a *Address) UnmarshalJSON(data []byte) error {
	felt := a.Felt()
	err := felt.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	*a = Address(felt)
	return nil
}

func (a *Address) Bytes() []byte {
	return a.Felt().Bytes()
}

func (a Address) Hex() string {
	return a.Felt().Hex()
}
