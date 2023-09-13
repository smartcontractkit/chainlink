package merklemulti

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
)

type singleLayerProof[H hashlib.Hash] struct {
	nextIndices []int
	subProof    []H
	sourceFlags []bool
}

type Proof[H hashlib.Hash] struct {
	Hashes      []H    `json:"hashes"`
	SourceFlags []bool `json:"source_flags"`
}

func (p Proof[H]) countSourceFlags(b bool) (count int) {
	for _, flag := range p.SourceFlags {
		if flag == b {
			count++
		}
	}
	return
}

const (
	SourceFromHashes = true
	SourceFromProof  = false
	// Maximum number of leaves in a Merkle tree. This is a limitation of the contract.
	MaxNumberTreeLeaves = 256
)

func parentIndex(idx int) int {
	return idx / 2
}

func siblingIndex(idx int) int {
	return idx ^ 1
}

func proveSingleLayer[H hashlib.Hash](layer []H, indices []int) (singleLayerProof[H], error) {
	var (
		authIndices []int
		nextIndices []int
		sourceFlags []bool
	)
	j := 0
	for j < len(indices) {
		x := indices[j]
		nextIndices = append(nextIndices, parentIndex(x))
		if j+1 < len(indices) && indices[j+1] == siblingIndex(x) {
			j++
			sourceFlags = append(sourceFlags, SourceFromHashes)
		} else {
			authIndices = append(authIndices, siblingIndex(x))
			sourceFlags = append(sourceFlags, SourceFromProof)
		}
		j++
	}
	var subProof []H
	for _, i := range authIndices {
		if i < 0 || i >= len(layer) {
			return singleLayerProof[H]{}, fmt.Errorf("auth index %d is out of bounds", i)
		}
		subProof = append(subProof, layer[i])
	}
	return singleLayerProof[H]{
		nextIndices: nextIndices,
		subProof:    subProof,
		sourceFlags: sourceFlags,
	}, nil
}

type Tree[H hashlib.Hash] struct {
	layers [][]H
}

func NewTree[H hashlib.Hash](ctx hashlib.Ctx[H], leafHashes []H) (*Tree[H], error) {
	if len(leafHashes) == 0 {
		return nil, errors.New("Cannot construct a tree without leaves")
	}
	var layer = make([]H, len(leafHashes))
	copy(layer, leafHashes)
	var layers = [][]H{layer}
	var curr int
	for len(layer) > 1 {
		paddedLayer, nextLayer := computeNextLayer(ctx, layer)
		layers[curr] = paddedLayer
		curr++
		layers = append(layers, nextLayer)
		layer = nextLayer
	}
	return &Tree[H]{
		layers: layers,
	}, nil
}

// Revive appears confused with the generics "receiver name t should be consistent with previous receiver name p for invalid-type"
//
//revive:disable:receiver-naming
func (t *Tree[H]) String() string {
	b := strings.Builder{}
	for _, layer := range t.layers {
		b.WriteString(fmt.Sprintf("%v", layer))
	}
	return b.String()
}

func (t *Tree[H]) Root() H {
	return t.layers[len(t.layers)-1][0]
}

func (t *Tree[H]) Prove(indices []int) (Proof[H], error) {
	var proof Proof[H]
	for _, layer := range t.layers[:len(t.layers)-1] {
		res, err := proveSingleLayer(layer, indices)
		if err != nil {
			return Proof[H]{}, err
		}

		indices = res.nextIndices
		proof.Hashes = append(proof.Hashes, res.subProof...)
		proof.SourceFlags = append(proof.SourceFlags, res.sourceFlags...)
	}
	return proof, nil
}

func computeNextLayer[H hashlib.Hash](ctx hashlib.Ctx[H], layer []H) ([]H, []H) {
	if len(layer) == 1 {
		return layer, layer
	}
	if len(layer)%2 != 0 {
		layer = append(layer, ctx.ZeroHash())
	}
	var nextLayer []H
	for i := 0; i < len(layer); i += 2 {
		nextLayer = append(nextLayer, ctx.HashInternal(layer[i], layer[i+1]))
	}
	return layer, nextLayer
}

func VerifyComputeRoot[H hashlib.Hash](ctx hashlib.Ctx[H], leafHashes []H, proof Proof[H]) (H, error) {
	leavesLength := len(leafHashes)
	proofsLength := len(proof.Hashes)
	if leavesLength == 0 && proofsLength == 0 {
		return ctx.ZeroHash(), errors.Errorf("leaves and proofs are empty")
	}
	if leavesLength > MaxNumberTreeLeaves+1 || proofsLength > MaxNumberTreeLeaves+1 {
		return ctx.ZeroHash(), errors.Errorf("leaves or proofs length is beyond the limit %d", MaxNumberTreeLeaves)
	}
	totalHashes := leavesLength + proofsLength - 1
	if totalHashes > MaxNumberTreeLeaves {
		return ctx.ZeroHash(), errors.Errorf("total hashes length cannot me larger than %d", MaxNumberTreeLeaves)
	}
	if totalHashes != len(proof.SourceFlags) {
		return ctx.ZeroHash(), errors.Errorf("hashes %d != sourceFlags %d", totalHashes, len(proof.SourceFlags))
	}
	if totalHashes == 0 {
		return leafHashes[0], nil
	}
	sourceProofCount := proof.countSourceFlags(SourceFromProof)
	if sourceProofCount != proofsLength {
		return ctx.ZeroHash(), errors.Errorf("proof source flags %d != proof hashes %d", sourceProofCount, proofsLength)
	}
	hashes := make([]H, totalHashes)
	for i := 0; i < totalHashes; i++ {
		hashes = append(hashes, leafHashes[0])
	}
	var (
		leafPos  int
		hashPos  int
		proofPos int
	)
	for i := 0; i < totalHashes; i++ {
		var a, b H
		//nolint:gosimple
		if proof.SourceFlags[i] == SourceFromHashes {
			if leafPos < leavesLength {
				a = leafHashes[leafPos]
				leafPos++
			} else {
				a = hashes[hashPos]
				hashPos++
			}
			//nolint:gosimple
		} else if proof.SourceFlags[i] == SourceFromProof {
			a = proof.Hashes[proofPos]
			proofPos++
		}
		if leafPos < leavesLength {
			b = leafHashes[leafPos]
			leafPos++
		} else {
			b = hashes[hashPos]
			hashPos++
		}
		hashes[i] = ctx.HashInternal(a, b)
	}
	if hashPos != totalHashes-1 ||
		leafPos != leavesLength ||
		proofPos != proofsLength {
		return ctx.ZeroHash(), errors.Errorf("not all proofs used during processing")
	}
	return hashes[totalHashes-1], nil
}
