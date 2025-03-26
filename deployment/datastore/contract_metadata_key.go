package datastore

type ContractMetadataKey interface {
	Comparable[ContractMetadataKey]

	Chain() uint64
	Address() string
}

type contractMetadataKey struct {
	chain   uint64
	address string
}

func (c contractMetadataKey) Chain() uint64   { return c.chain }
func (c contractMetadataKey) Address() string { return c.address }

func (c contractMetadataKey) Equals(other ContractMetadataKey) bool {
	return c.chain == other.Chain() &&
		c.address == other.Address()
}

func NewContractMetadataKey(chain uint64, address string) ContractMetadataKey {
	return contractMetadataKey{
		chain:   chain,
		address: address,
	}
}
