package ccipocr3

type CommitPluginReport struct {
	MerkleRoots  []MerkleRootChain `json:"merkleRoots"`
	PriceUpdates PriceUpdates      `json:"priceUpdates"`
}

func NewCommitPluginReport(merkleRoots []MerkleRootChain, tokenPriceUpdates []TokenPrice, gasPriceUpdate []GasPriceChain) CommitPluginReport {
	return CommitPluginReport{
		MerkleRoots:  merkleRoots,
		PriceUpdates: PriceUpdates{TokenPriceUpdates: tokenPriceUpdates, GasPriceUpdates: gasPriceUpdate},
	}
}

// IsEmpty returns true if the CommitPluginReport is empty
func (r CommitPluginReport) IsEmpty() bool {
	return len(r.MerkleRoots) == 0 &&
		len(r.PriceUpdates.TokenPriceUpdates) == 0 &&
		len(r.PriceUpdates.GasPriceUpdates) == 0
}

type MerkleRootChain struct {
	ChainSel     ChainSelector `json:"chain"`
	SeqNumsRange SeqNumRange   `json:"seqNumsRange"`
	MerkleRoot   Bytes32       `json:"merkleRoot"`
}

func NewMerkleRootChain(
	chainSel ChainSelector,
	seqNumsRange SeqNumRange,
	merkleRoot Bytes32,
) MerkleRootChain {
	return MerkleRootChain{
		ChainSel:     chainSel,
		SeqNumsRange: seqNumsRange,
		MerkleRoot:   merkleRoot,
	}
}

type PriceUpdates struct {
	TokenPriceUpdates []TokenPrice    `json:"tokenPriceUpdates"`
	GasPriceUpdates   []GasPriceChain `json:"gasPriceUpdates"`
}
