package verkle

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"sort"

	"github.com/crate-crypto/go-ipa/banderwagon"
	"golang.org/x/sync/errgroup"
)

// BatchNewLeafNodeData is a struct that contains the data needed to create a new leaf node.
type BatchNewLeafNodeData struct {
	Stem   []byte
	Values map[byte][]byte
}

// BatchNewLeafNode creates a new leaf node from the given data. It optimizes LeafNode creation
// by batching expensive cryptography operations. It returns the LeafNodes sorted by stem.
func BatchNewLeafNode(nodesValues []BatchNewLeafNodeData) ([]LeafNode, error) {
	cfg := GetConfig()
	ret := make([]LeafNode, len(nodesValues))

	numBatches := runtime.NumCPU()
	batchSize := len(nodesValues) / numBatches

	group, _ := errgroup.WithContext(context.Background())
	for i := 0; i < numBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if i == numBatches-1 {
			end = len(nodesValues)
		}

		work := func(ret []LeafNode, nodesValues []BatchNewLeafNodeData) func() error {
			return func() error {
				c1c2points := make([]*Point, 2*len(nodesValues))
				c1c2frs := make([]*Fr, 2*len(nodesValues))
				for i, nv := range nodesValues {
					valsslice := make([][]byte, NodeWidth)
					for idx := range nv.Values {
						valsslice[idx] = nv.Values[idx]
					}

					var leaf *LeafNode
					leaf, err := NewLeafNode(nv.Stem, valsslice)
					if err != nil {
						return err
					}
					ret[i] = *leaf

					c1c2points[2*i], c1c2points[2*i+1] = ret[i].c1, ret[i].c2
					c1c2frs[2*i], c1c2frs[2*i+1] = new(Fr), new(Fr)
				}

				if err := banderwagon.BatchMapToScalarField(c1c2frs, c1c2points); err != nil {
					return fmt.Errorf("mapping to scalar field: %s", err)
				}

				var poly [NodeWidth]Fr
				poly[0].SetUint64(1)
				for i, nv := range nodesValues {
					if err := StemFromBytes(&poly[1], nv.Stem); err != nil {
						return err
					}
					poly[2] = *c1c2frs[2*i]
					poly[3] = *c1c2frs[2*i+1]

					ret[i].commitment = cfg.CommitToPoly(poly[:], 252)
				}
				return nil
			}
		}
		group.Go(work(ret[start:end], nodesValues[start:end]))
	}
	if err := group.Wait(); err != nil {
		return nil, fmt.Errorf("creating leaf node: %s", err)
	}

	sort.Slice(ret, func(i, j int) bool {
		return bytes.Compare(ret[i].stem, ret[j].stem) < 0
	})

	return ret, nil
}

// firstDiffByteIdx will return the first index in which the two stems differ.
// Both stems *must* be different.
func firstDiffByteIdx(stem1 []byte, stem2 []byte) int {
	for i := range stem1 {
		if stem1[i] != stem2[i] {
			return i
		}
	}
	panic("stems are equal")
}

func (n *InternalNode) InsertMigratedLeaves(leaves []LeafNode, resolver NodeResolverFn) error {
	sort.Slice(leaves, func(i, j int) bool {
		return bytes.Compare(leaves[i].stem, leaves[j].stem) < 0
	})

	// We first mark all children of the subtreess that we'll update in parallel,
	// so the subtree updating doesn't produce a concurrent access to n.cowChild(...).
	var lastChildrenIdx = -1
	for i := range leaves {
		if int(leaves[i].stem[0]) != lastChildrenIdx {
			lastChildrenIdx = int(leaves[i].stem[0])
			if _, ok := n.children[lastChildrenIdx].(HashedNode); ok {
				serialized, err := resolver([]byte{byte(lastChildrenIdx)})
				if err != nil {
					return fmt.Errorf("resolving node: %s", err)
				}
				resolved, err := ParseNode(serialized, 1)
				if err != nil {
					return fmt.Errorf("parsing node %x: %w", serialized, err)
				}
				n.children[lastChildrenIdx] = resolved
			}
			n.cowChild(byte(lastChildrenIdx))
		}
	}

	// We insert the migrated leaves for each subtree of the root node.
	group, _ := errgroup.WithContext(context.Background())
	group.SetLimit(runtime.NumCPU())
	currStemFirstByte := 0
	for i := range leaves {
		if leaves[currStemFirstByte].stem[0] != leaves[i].stem[0] {
			start := currStemFirstByte
			end := i
			group.Go(func() error {
				return n.insertMigratedLeavesSubtree(leaves[start:end], resolver)
			})
			currStemFirstByte = i
		}
	}
	group.Go(func() error {
		return n.insertMigratedLeavesSubtree(leaves[currStemFirstByte:], resolver)
	})
	if err := group.Wait(); err != nil {
		return fmt.Errorf("inserting migrated leaves: %w", err)
	}

	return nil
}

func (n *InternalNode) insertMigratedLeavesSubtree(leaves []LeafNode, resolver NodeResolverFn) error { // skipcq: GO-R1005
	for i := range leaves {
		ln := leaves[i]
		parent := n

		// Look for the appropriate parent for the leaf node.
		for {
			if _, ok := parent.children[ln.stem[parent.depth]].(HashedNode); ok {
				serialized, err := resolver(ln.stem[:parent.depth+1])
				if err != nil {
					return fmt.Errorf("resolving node path=%x: %w", ln.stem[:parent.depth+1], err)
				}
				resolved, err := ParseNode(serialized, parent.depth+1)
				if err != nil {
					return fmt.Errorf("parsing node %x: %w", serialized, err)
				}
				parent.children[ln.stem[parent.depth]] = resolved
			}

			nextParent, ok := parent.children[ln.stem[parent.depth]].(*InternalNode)
			if !ok {
				break
			}

			parent.cowChild(ln.stem[parent.depth])
			parent = nextParent
		}

		switch node := parent.children[ln.stem[parent.depth]].(type) {
		case Empty:
			parent.cowChild(ln.stem[parent.depth])
			parent.children[ln.stem[parent.depth]] = &ln
			ln.setDepth(parent.depth + 1)
		case *LeafNode:
			if bytes.Equal(node.stem, ln.stem) {
				// In `ln` we have migrated key/values which should be copied to the leaf
				// only if there isn't a value there. If there's a value, we skip it since
				// our migrated value is stale.
				nonPresentValues := make([][]byte, NodeWidth)
				for i := range ln.values {
					if node.values[i] == nil {
						nonPresentValues[i] = ln.values[i]
					}
				}

				if err := node.updateMultipleLeaves(nonPresentValues); err != nil {
					return fmt.Errorf("updating leaves: %s", err)
				}
				continue
			}

			// Otherwise, we need to create the missing internal nodes depending in the fork point in their stems.
			idx := firstDiffByteIdx(node.stem, ln.stem)
			// We do a sanity check to make sure that the fork point is not before the current depth.
			if byte(idx) <= parent.depth {
				return fmt.Errorf("unexpected fork point %d for nodes %x and %x", idx, node.stem, ln.stem)
			}
			// Create the missing internal nodes.
			for i := parent.depth + 1; i <= byte(idx); i++ {
				nextParent := newInternalNode(parent.depth + 1).(*InternalNode)
				parent.cowChild(ln.stem[parent.depth])
				parent.children[ln.stem[parent.depth]] = nextParent
				parent = nextParent
			}
			// Add old and new leaf node to the latest created parent.
			parent.cowChild(node.stem[parent.depth])
			parent.children[node.stem[parent.depth]] = node
			node.setDepth(parent.depth + 1)
			parent.cowChild(ln.stem[parent.depth])
			parent.children[ln.stem[parent.depth]] = &ln
			ln.setDepth(parent.depth + 1)
		default:
			return fmt.Errorf("unexpected node type %T", node)
		}
	}
	return nil
}
