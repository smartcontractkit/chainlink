package direct

import (
	"encoding/json"

	"github.com/consensys/gnark-crypto/accumulator/merkletree"
)

type SerializedMerkleTree struct {
	Nodes     [][]byte `json:"nodes"`
	Root      []byte   `json:"root"`
	Signature []byte   `json:"signature"`
}

// SerializeMerkleTree serializes the Merkle Tree and signature to JSON
func SerializeMerkleTree(tree *merkletree.Tree, signature []byte) ([]byte, error) {
	tree.
	nodes := tree.GetNodes()
	var serializedNodes [][]byte
	for _, node := range nodes {
		serializedNodes = append(serializedNodes, node.Hash)
	}
	serializedTree := SerializedMerkleTree{
		Nodes:     serializedNodes,
		Root:      tree.MerkleRoot(),
		Signature: signature,
	}
	return json.Marshal(serializedTree)
}

// DeserializeMerkleTree deserializes the JSON back into a Merkle Tree and signature
func DeserializeMerkleTree(data []byte) (*merkletree.Tree, []byte, error) {
	var serializedTree SerializedMerkleTree
	if err := json.Unmarshal(data, &serializedTree); err != nil {
		return nil, nil, err
	}
	var nodes []merkletree.Node
	for _, hash := range serializedTree.Nodes {
		nodes = append(nodes, merkletree.Node{Hash: hash})
	}

	tree := &merkletree.Tree{Nodes: nodes}
	return tree, serializedTree.Signature, nil
}
