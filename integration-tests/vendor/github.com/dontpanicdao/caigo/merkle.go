package caigo

import (
	"fmt"
	"math/big"
)

type FixedSizeMerkleTree struct {
	Leaves   []*big.Int
	Branches [][]*big.Int
	Root     *big.Int
}

func NewFixedSizeMerkleTree(leaves ...*big.Int) (*FixedSizeMerkleTree, error) {
	mt := &FixedSizeMerkleTree{
		Leaves:   leaves,
		Branches: [][]*big.Int{},
	}
	root, err := mt.build(leaves)
	if err != nil {
		return nil, err
	}
	mt.Root = root
	return mt, err
}

func MerkleHash(x, y *big.Int) (*big.Int, error) {
	if x.Cmp(y) <= 0 {
		return Curve.HashElements([]*big.Int{x, y})
	}
	return Curve.HashElements([]*big.Int{y, x})
}

func (mt *FixedSizeMerkleTree) build(leaves []*big.Int) (*big.Int, error) {
	if len(leaves) == 1 {
		return leaves[0], nil
	}
	mt.Branches = append(mt.Branches, leaves)
	newLeaves := []*big.Int{}
	for i := 0; i < len(leaves); i += 2 {
		if i+1 == len(leaves) {
			hash, err := MerkleHash(leaves[i], big.NewInt(0))
			if err != nil {
				return nil, err
			}
			newLeaves = append(newLeaves, hash)
			break
		}
		hash, err := MerkleHash(leaves[i], leaves[i+1])
		if err != nil {
			return nil, err
		}
		newLeaves = append(newLeaves, hash)
	}
	return mt.build(newLeaves)
}

func (mt *FixedSizeMerkleTree) Proof(leaf *big.Int) ([]*big.Int, error) {
	return mt.recursiveProof(leaf, 0, []*big.Int{})
}

func (mt *FixedSizeMerkleTree) recursiveProof(leaf *big.Int, branchIndex int, hashPath []*big.Int) ([]*big.Int, error) {
	if branchIndex >= len(mt.Branches) {
		return hashPath, nil
	}
	branch := mt.Branches[branchIndex]
	index := -1
	for k, v := range branch {
		if v.Cmp(leaf) == 0 {
			index = k
			break
		}
	}
	if index == -1 {
		return nil, fmt.Errorf("key 0x%s not found in branch", leaf.Text(16))
	}
	nextProof := big.NewInt(0)
	if index%2 == 0 && index < len(branch) {
		nextProof = branch[index+1]
	}
	if index%2 != 0 {
		nextProof = branch[index-1]
	}
	newLeaf, err := MerkleHash(leaf, nextProof)
	if err != nil {
		return nil, fmt.Errorf("nextproof error: %v", err)
	}
	newHashPath := append(hashPath, nextProof)
	return mt.recursiveProof(newLeaf, branchIndex+1, newHashPath)
}

func ProofMerklePath(root *big.Int, leaf *big.Int, path []*big.Int) bool {
	if len(path) == 0 {
		return root.Cmp(leaf) == 0
	}
	nexLeaf, err := MerkleHash(leaf, path[0])
	if err != nil {
		return false
	}
	return ProofMerklePath(root, nexLeaf, path[1:])
}
