package common

// AddressCodec is a struct that holds the chain specific address codecs
type AddressCodec struct {
	EVMAddressCodec    ChainSpecificAddressCodec
	SolanaAddressCodec ChainSpecificAddressCodec
}

// NewAddressCodec is a constructor for NewAddressCodec
func NewAddressCodec(evmAddrCodec, solanaAddrCodec ChainSpecificAddressCodec) AddressCodec {
	return AddressCodec{
		EVMAddressCodec:    evmAddrCodec,
		SolanaAddressCodec: solanaAddrCodec,
	}
}
