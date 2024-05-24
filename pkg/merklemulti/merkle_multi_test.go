package merklemulti

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/stat/combin"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti/fixtures"
)

var (
	hasher           = hashutil.NewKeccak()
	a, b, c, d, e, f = hasher.Hash([]byte{0xa}), hasher.Hash([]byte{0xb}), hasher.Hash([]byte{0xc}), hasher.Hash([]byte{0xd}), hasher.Hash([]byte{0xe}), hasher.Hash([]byte{0xf})
)

func mustDecode(input string) []byte {
	b, err := hex.DecodeString(input[2:])
	if err != nil {
		panic(err)
	}
	return b
}

func hashesFromHexStrings(hexStrs []string) [][32]byte {
	var hashes [][32]byte
	for _, hexStr := range hexStrs {
		var hash [32]byte
		copy(hash[:], mustDecode(fmt.Sprintf("0x%s", hexStr)))
		hashes = append(hashes, hash)
	}
	return hashes
}

func TestReturnErrorForTooLargeInput(t *testing.T) {
	leavesOrProofsToLarge := "leaves or proofs length is beyond the limit 256"

	tests := []struct {
		name                 string
		leavesLen, proofsLen int
		errorMessage         string
	}{
		{"both below maximum, but sum above", MaxNumberTreeLeaves + 1, MaxNumberTreeLeaves + 1, "total hashes length cannot me larger than 256"},
		{"both maximum lengths", MaxNumberTreeLeaves + 2, MaxNumberTreeLeaves + 2, leavesOrProofsToLarge},
		{"leaves are too large", MaxNumberTreeLeaves + 2, 1, leavesOrProofsToLarge},
		{"proofs are too large", 2, MaxNumberTreeLeaves + 2, leavesOrProofsToLarge},
		{"empty", 0, 0, "leaves and proofs are empty"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			leaves := make([][32]byte, test.leavesLen)
			proofs := make([][32]byte, test.proofsLen)

			var flags []bool
			flagsLength := test.leavesLen + test.proofsLen - 1
			if flagsLength > 0 {
				flags = make([]bool, flagsLength)
			}

			_, err := VerifyComputeRoot(hasher, leaves, Proof[[32]byte]{Hashes: proofs, SourceFlags: flags})
			require.Error(t, err)
			require.Equal(t, err.Error(), test.errorMessage)
		})
	}
}

func TestErrorWhenNotAllProofsCanBeUsed(t *testing.T) {
	leaves := [][32]byte{a, b}
	proofs := [][32]byte{c, d}
	sourceFlags := []bool{false, true, true}

	_, err := VerifyComputeRoot(hasher, leaves, Proof[[32]byte]{Hashes: proofs, SourceFlags: sourceFlags})
	require.Error(t, err)
	require.Equal(t, err.Error(), "proof source flags 1 != proof hashes 2")
}

func TestSpecFixtureVerifyProof(t *testing.T) {
	for _, testVector := range fixtures.TestVectors {
		var leafHashes = hashesFromHexStrings(testVector.ProofLeaves)
		var proofHashes = hashesFromHexStrings(testVector.ProofHashes)
		computedRoot, err := VerifyComputeRoot(hasher, leafHashes, Proof[[32]byte]{
			Hashes: proofHashes, SourceFlags: testVector.ProofFlags,
		})
		require.NoError(t, err)
		assert.Equal(t, mustDecode(fmt.Sprintf("0x%s", testVector.ExpectedRoot)), computedRoot[:])
	}
}

func TestSpecFixtureNewTree(t *testing.T) {
	for _, testVector := range fixtures.TestVectors {
		var leafHashes = hashesFromHexStrings(testVector.AllLeafs)
		mctx := hashutil.NewKeccak()
		tree, err := NewTree(mctx, leafHashes)
		assert.NoError(t, err)
		actualRoot := tree.Root()
		assert.Equal(t, testVector.ExpectedRoot, hex.EncodeToString(actualRoot[:]))
	}
}

func TestPadding(t *testing.T) {
	tr4, err := NewTree(hasher, [][32]byte{a, b, c})
	require.NoError(t, err)
	assert.Equal(t, 4, len(tr4.layers[0]))
	tr8, err := NewTree(hasher, [][32]byte{a, b, c, d, e})
	require.NoError(t, err)
	assert.Equal(t, 6, len(tr8.layers[0]))
	assert.Equal(t, 4, len(tr8.layers[1]))
	p, err := tr8.Prove([]int{0})
	assert.NoError(t, err)
	h, err := VerifyComputeRoot(hasher, [][32]byte{a}, p)
	require.NoError(t, err)
	assert.Equal(t, tr8.Root(), h)
	expected := hasher.HashInternal(hasher.HashInternal(hasher.HashInternal(a, b), hasher.HashInternal(c, d)), hasher.HashInternal(hasher.HashInternal(e, hasher.ZeroHash()), hasher.ZeroHash()))
	assert.Equal(t, expected, tr8.Root())
}

func TestMerkleMultiProofSecondPreimage(t *testing.T) {
	tr, err := NewTree(hasher, [][32]byte{a, b})
	require.NoError(t, err)
	pr, err := tr.Prove([]int{0})
	require.NoError(t, err)
	root, err := VerifyComputeRoot(hasher, [][32]byte{a}, pr)
	require.NoError(t, err)
	assert.Equal(t, root, tr.Root())
	tr2, err := NewTree(hasher, [][32]byte{hasher.Hash(append(a[:], b[:]...))})
	require.NoError(t, err)
	assert.NotEqual(t, tr2.Root(), tr.Root())
}

func TestMerkleMultiProof(t *testing.T) {
	leafHashes := [][32]byte{a, b, c, d, e, f}
	expectedRoots := [][32]byte{
		a,
		hasher.HashInternal(a, b),
		hasher.HashInternal(hasher.HashInternal(a, b), hasher.HashInternal(c, hasher.ZeroHash())),
		hasher.HashInternal(hasher.HashInternal(a, b), hasher.HashInternal(c, d)),
		hasher.HashInternal(hasher.HashInternal(hasher.HashInternal(a, b), hasher.HashInternal(c, d)), hasher.HashInternal(hasher.HashInternal(e, hasher.ZeroHash()), hasher.ZeroHash())),
		hasher.HashInternal(hasher.HashInternal(hasher.HashInternal(a, b), hasher.HashInternal(c, d)), hasher.HashInternal(hasher.HashInternal(e, f), hasher.ZeroHash())),
	}
	// For every size tree from 0..len(leaves)
	for length := 1; length <= len(leafHashes); length++ {
		tr, err := NewTree(hasher, leafHashes[:length])
		require.NoError(t, err)
		expectedRoot := expectedRoots[length-1]
		require.Equal(t, tr.Root(), expectedRoot)
		// Prove every subset of its leaves
		for k := 1; k <= length; k++ {
			gen := combin.NewCombinationGenerator(length, k)
			for gen.Next() {
				leaveIndices := gen.Combination(nil)
				proof, err := tr.Prove(leaveIndices)
				require.NoError(t, err)
				var leavesToProve [][32]byte
				for _, idx := range leaveIndices {
					leavesToProve = append(leavesToProve, leafHashes[idx])
				}
				root, err := VerifyComputeRoot(hasher, leavesToProve, proof)
				require.NoError(t, err)
				assert.Equal(t, expectedRoot, root)
			}
		}
	}

	t.Run("invalid indices should not lead to panic", func(t *testing.T) {
		tr, err := NewTree(hasher, leafHashes[:])
		require.NoError(t, err)
		_, err = tr.Prove([]int{1, 2, 3, 9999})
		require.Error(t, err)
	})
}
