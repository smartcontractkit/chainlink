// This is free and unencumbered software released into the public domain.
//
// Anyone is free to copy, modify, publish, use, compile, sell, or
// distribute this software, either in source code form or as a compiled
// binary, for any purpose, commercial or non-commercial, and by any
// means.
//
// In jurisdictions that recognize copyright laws, the author or authors
// of this software dedicate any and all copyright interest in the
// software to the public domain. We make this dedication for the benefit
// of the public at large and to the detriment of our heirs and
// successors. We intend this dedication to be an overt act of
// relinquishment in perpetuity of all present and future rights to this
// software under copyright law.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// For more information, please refer to <https://unlicense.org>

package verkle

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"sync"

	"github.com/crate-crypto/go-ipa/banderwagon"
)

type (
	NodeFlushFn    func([]byte, VerkleNode)
	NodeResolverFn func([]byte) ([]byte, error)
)

type keylist [][]byte

func (kl keylist) Len() int {
	return len(kl)
}

func (kl keylist) Less(i, j int) bool {
	return bytes.Compare(kl[i], kl[j]) == -1
}

func (kl keylist) Swap(i, j int) {
	kl[i], kl[j] = kl[j], kl[i]
}

type VerkleNode interface {
	// Insert or Update value into the tree
	Insert([]byte, []byte, NodeResolverFn) error

	// Delete a leaf with the given key
	Delete([]byte, NodeResolverFn) (bool, error)

	// Get value at a given key
	Get([]byte, NodeResolverFn) ([]byte, error)

	// Commit computes the commitment of the node. The
	// result (the curve point) is cached.
	Commit() *Point

	// Commitment is a getter for the cached commitment
	// to this node.
	Commitment() *Point

	// Hash returns the field representation of the commitment.
	Hash() *Fr

	// GetProofItems collects the various proof elements, and
	// returns them breadth-first. On top of that, it returns
	// one "extension status" per stem, and an alternate stem
	// if the key is missing but another stem has been found.
	GetProofItems(keylist, NodeResolverFn) (*ProofElements, []byte, [][]byte, error)

	// Serialize encodes the node to RLP.
	Serialize() ([]byte, error)

	// Copy a node and its children
	Copy() VerkleNode

	// toDot returns a string representing this subtree in DOT language
	toDot(string, string) string

	setDepth(depth byte)
}

// ProofElements gathers the elements needed to build a proof.
type ProofElements struct {
	Cis    []*Point
	Zis    []byte
	Yis    []*Fr
	Fis    [][]Fr
	ByPath map[string]*Point // Gather commitments by path
	Vals   [][]byte          // list of values read from the tree

	// dedups flags the presence of each (Ci,zi) tuple
	dedups map[*Point]map[byte]struct{}
}

// Merge merges the elements of two proofs and removes duplicates.
func (pe *ProofElements) Merge(other *ProofElements) {
	// Build the local map if it's missing
	if pe.dedups == nil {
		pe.dedups = make(map[*Point]map[byte]struct{})

		for i, ci := range pe.Cis {
			if _, ok := pe.dedups[ci]; !ok {
				pe.dedups[ci] = make(map[byte]struct{})
			}

			pe.dedups[ci][pe.Zis[i]] = struct{}{}
		}
	}

	for i, ci := range other.Cis {
		if _, ok := pe.dedups[ci]; !ok {
			// First time this commitment has been seen, create
			// the map and flatten the zi.
			pe.dedups[ci] = make(map[byte]struct{})
		}

		if _, ok := pe.dedups[ci][other.Zis[i]]; ok {
			// duplicate, skip
			continue
		}

		pe.dedups[ci][other.Zis[i]] = struct{}{}

		pe.Cis = append(pe.Cis, ci)
		pe.Zis = append(pe.Zis, other.Zis[i])
		pe.Yis = append(pe.Yis, other.Yis[i])
		if pe.Fis != nil {
			pe.Fis = append(pe.Fis, other.Fis[i])
		}
	}

	for path, C := range other.ByPath {
		if _, ok := pe.ByPath[path]; !ok {
			pe.ByPath[path] = C
		}
	}

	pe.Vals = append(pe.Vals, other.Vals...)
}

const (
	// These types will distinguish internal
	// and leaf nodes when decoding from RLP.
	internalRLPType byte = 1
	leafRLPType     byte = 2
)

type (
	// Represents an internal node at any level
	InternalNode struct {
		// List of child nodes of this internal node.
		children []VerkleNode

		// node depth in the tree, in bits
		depth byte

		// Cache the commitment value
		commitment *Point

		cow map[byte]*Point
	}

	LeafNode struct {
		stem   []byte
		values [][]byte

		commitment *Point
		c1, c2     *Point

		depth byte

		// IsPOAStub indicates if this LeafNode is a proof of absence
		// for a steam that isn't present in the tree. This flag is only
		// true in the context of a stateless tree.
		isPOAStub bool
	}
)

func (n *InternalNode) toExportable() *ExportableInternalNode {
	comm := n.commitment.Bytes()
	exportable := &ExportableInternalNode{
		Children:   make([]interface{}, NodeWidth),
		Commitment: comm[:],
	}

	for i := range exportable.Children {
		switch child := n.children[i].(type) {
		case Empty:
			exportable.Children[i] = nil
		case HashedNode:
			exportable.Children[i] = n.commitment.Bytes()
		case *InternalNode:
			exportable.Children[i] = child.toExportable()
		case *LeafNode:
			exportable.Children[i] = &ExportableLeafNode{
				Stem:   child.stem,
				Values: child.values,
				C:      child.commitment.Bytes(),
				C1:     child.c1.Bytes(),
			}
		default:
			panic("unexportable type")
		}
	}
	return exportable
}

// Turn an internal node into a JSON string
func (n *InternalNode) ToJSON() ([]byte, error) {
	return json.Marshal(n.toExportable())
}

func newInternalNode(depth byte) VerkleNode {
	node := new(InternalNode)
	node.children = make([]VerkleNode, NodeWidth)
	for idx := range node.children {
		node.children[idx] = Empty(struct{}{})
	}
	node.depth = depth
	node.commitment = new(Point).SetIdentity()
	return node
}

// New creates a new tree root
func New() VerkleNode {
	return newInternalNode(0)
}

func NewStatelessInternal(depth byte, comm *Point) VerkleNode {
	node := &InternalNode{
		children:   make([]VerkleNode, NodeWidth),
		depth:      depth,
		commitment: comm,
	}
	for idx := range node.children {
		node.children[idx] = UnknownNode(struct{}{})
	}
	return node
}

// New creates a new leaf node
func NewLeafNode(stem []byte, values [][]byte) (*LeafNode, error) {
	cfg := GetConfig()

	// C1.
	var c1poly [NodeWidth]Fr
	var c1 *Point
	count, err := fillSuffixTreePoly(c1poly[:], values[:NodeWidth/2])
	if err != nil {
		return nil, err
	}
	containsEmptyCodeHash := len(c1poly) >= EmptyCodeHashSecondHalfIdx &&
		c1poly[EmptyCodeHashFirstHalfIdx].Equal(&EmptyCodeHashFirstHalfValue) &&
		c1poly[EmptyCodeHashSecondHalfIdx].Equal(&EmptyCodeHashSecondHalfValue)
	if containsEmptyCodeHash {
		// Clear out values of the cached point.
		c1poly[EmptyCodeHashFirstHalfIdx] = FrZero
		c1poly[EmptyCodeHashSecondHalfIdx] = FrZero
		// Calculate the remaining part of c1 and add to the base value.
		partialc1 := cfg.CommitToPoly(c1poly[:], NodeWidth-count-2)
		c1 = new(Point)
		c1.Add(&EmptyCodeHashPoint, partialc1)
	} else {
		c1 = cfg.CommitToPoly(c1poly[:], NodeWidth-count)
	}

	// C2.
	var c2poly [NodeWidth]Fr
	count, err = fillSuffixTreePoly(c2poly[:], values[NodeWidth/2:])
	if err != nil {
		return nil, err
	}
	c2 := cfg.CommitToPoly(c2poly[:], NodeWidth-count)

	// Root commitment preparation for calculation.
	stem = stem[:StemSize] // enforce a 31-byte length
	var poly [NodeWidth]Fr
	poly[0].SetUint64(1)
	if err := StemFromBytes(&poly[1], stem); err != nil {
		return nil, err
	}
	if err := banderwagon.BatchMapToScalarField([]*Fr{&poly[2], &poly[3]}, []*Point{c1, c2}); err != nil {
		return nil, fmt.Errorf("batch mapping to scalar fields: %s", err)
	}

	return &LeafNode{
		// depth will be 0, but the commitment calculation
		// does not need it, and so it won't be free.
		values:     values,
		stem:       stem,
		commitment: cfg.CommitToPoly(poly[:], NodeWidth-4),
		c1:         c1,
		c2:         c2,
	}, nil
}

// NewLeafNodeWithNoComms create a leaf node but does compute its
// commitments. The created node's commitments are intended to be
// initialized with `SetTrustedBytes` in a deserialization context.
func NewLeafNodeWithNoComms(stem []byte, values [][]byte) *LeafNode {
	return &LeafNode{
		// depth will be 0, but the commitment calculation
		// does not need it, and so it won't be free.
		values: values,
		stem:   stem,
	}
}

// Children return the children of the node. The returned slice is
// internal to the tree, so callers *must* consider it readonly.
func (n *InternalNode) Children() []VerkleNode {
	return n.children
}

// SetChild *replaces* the child at the given index with the given node.
func (n *InternalNode) SetChild(i int, c VerkleNode) error {
	if i >= NodeWidth {
		return errors.New("child index higher than node width")
	}
	n.children[i] = c
	return nil
}

func (n *InternalNode) cowChild(index byte) {
	if n.cow == nil {
		n.cow = make(map[byte]*Point)
	}

	if n.cow[index] == nil {
		n.cow[index] = new(Point)
		n.cow[index].Set(n.children[index].Commitment())
	}
}

func (n *InternalNode) Insert(key []byte, value []byte, resolver NodeResolverFn) error {
	values := make([][]byte, NodeWidth)
	values[key[31]] = value
	return n.InsertValuesAtStem(key[:31], values, resolver)
}

func (n *InternalNode) InsertValuesAtStem(stem []byte, values [][]byte, resolver NodeResolverFn) error {
	nChild := offset2key(stem, n.depth) // index of the child pointed by the next byte in the key

	switch child := n.children[nChild].(type) {
	case UnknownNode:
		return errMissingNodeInStateless
	case Empty:
		n.cowChild(nChild)
		var err error
		n.children[nChild], err = NewLeafNode(stem, values)
		if err != nil {
			return err
		}
		n.children[nChild].setDepth(n.depth + 1)
	case HashedNode:
		if resolver == nil {
			return errInsertIntoHash
		}
		serialized, err := resolver(stem[:n.depth+1])
		if err != nil {
			return fmt.Errorf("verkle tree: error resolving node %x at depth %d: %w", stem, n.depth, err)
		}
		resolved, err := ParseNode(serialized, n.depth+1)
		if err != nil {
			return fmt.Errorf("verkle tree: error parsing resolved node %x: %w", stem, err)
		}
		n.children[nChild] = resolved
		n.cowChild(nChild)
		// recurse to handle the case of a LeafNode child that
		// splits.
		return n.InsertValuesAtStem(stem, values, resolver)
	case *LeafNode:
		if equalPaths(child.stem, stem) {
			// We can't insert any values into a POA leaf node.
			if child.isPOAStub {
				return errIsPOAStub
			}
			n.cowChild(nChild)
			return child.insertMultiple(stem, values)
		}
		n.cowChild(nChild)

		// A new branch node has to be inserted. Depending
		// on the next word in both keys, a recursion into
		// the moved leaf node can occur.
		nextWordInExistingKey := offset2key(child.stem, n.depth+1)
		newBranch := newInternalNode(n.depth + 1).(*InternalNode)
		newBranch.cowChild(nextWordInExistingKey)
		n.children[nChild] = newBranch
		newBranch.children[nextWordInExistingKey] = child
		child.depth += 1

		nextWordInInsertedKey := offset2key(stem, n.depth+1)
		if nextWordInInsertedKey == nextWordInExistingKey {
			return newBranch.InsertValuesAtStem(stem, values, resolver)
		}

		// Next word differs, so this was the last level.
		// Insert it directly into its final slot.
		leaf, err := NewLeafNode(stem, values)
		if err != nil {
			return err
		}
		leaf.setDepth(n.depth + 2)
		newBranch.cowChild(nextWordInInsertedKey)
		newBranch.children[nextWordInInsertedKey] = leaf
	case *InternalNode:
		n.cowChild(nChild)
		return child.InsertValuesAtStem(stem, values, resolver)
	default: // It should be an UknownNode.
		return errUnknownNodeType
	}

	return nil
}

// CreatePath inserts a given stem in the tree, placing it as
// described by stemInfo. Its third parameters is the list of
// commitments that have not been assigned a node. It returns
// the same list, save the commitments that were consumed
// during this call.
func (n *InternalNode) CreatePath(path []byte, stemInfo stemInfo, comms []*Point, values [][]byte) ([]*Point, error) {
	if len(path) == 0 {
		return comms, errors.New("invalid path")
	}

	// path is 1 byte long, the leaf node must be created
	if len(path) == 1 {
		switch stemInfo.stemType & 3 {
		case extStatusAbsentEmpty:
			// Set child to Empty so that, in a stateless context,
			// a node known to be absent is differentiated from an
			// unknown node.
			n.children[path[0]] = Empty{}
		case extStatusAbsentOther:
			// insert poa stem
			newchild := &LeafNode{
				commitment: comms[0],
				stem:       stemInfo.stem,
				values:     nil,
				depth:      n.depth + 1,
				isPOAStub:  true,
			}
			n.children[path[0]] = newchild
			comms = comms[1:]
		case extStatusPresent:
			// insert stem
			newchild := &LeafNode{
				commitment: comms[0],
				stem:       stemInfo.stem,
				values:     values,
				depth:      n.depth + 1,
			}
			n.children[path[0]] = newchild
			comms = comms[1:]
			if stemInfo.has_c1 {
				newchild.c1 = comms[0]
				comms = comms[1:]
			} else {
				newchild.c1 = new(Point)
			}
			if stemInfo.has_c2 {
				newchild.c2 = comms[0]
				comms = comms[1:]
			} else {
				newchild.c2 = new(Point)
			}
			for b, value := range stemInfo.values {
				newchild.values[b] = value
			}
		}
		return comms, nil
	}

	switch child := n.children[path[0]].(type) {
	case UnknownNode:
		// create the child node if missing
		n.children[path[0]] = NewStatelessInternal(n.depth+1, comms[0])
		comms = comms[1:]
	case *InternalNode:
	// nothing else to do
	case *LeafNode:
		return comms, fmt.Errorf("error rebuilding the tree from a proof: stem %x leads to an already-existing leaf node at depth %x", stemInfo.stem, n.depth)
	default:
		return comms, fmt.Errorf("error rebuilding the tree from a proof: stem %x leads to an unsupported node type %v", stemInfo.stem, child)
	}

	// This should only be used in the context of
	// stateless nodes, so panic if another node
	// type is found.
	child := n.children[path[0]].(*InternalNode)

	// recurse
	return child.CreatePath(path[1:], stemInfo, comms, values)
}

// GetValuesAtStem returns the all NodeWidth values of the stem.
// The returned slice is internal to the tree, so it *must* be considered readonly
// for callers.
func (n *InternalNode) GetValuesAtStem(stem []byte, resolver NodeResolverFn) ([][]byte, error) {
	nchild := offset2key(stem, n.depth) // index of the child pointed by the next byte in the key
	switch child := n.children[nchild].(type) {
	case UnknownNode:
		return nil, errMissingNodeInStateless
	case Empty:
		return nil, nil
	case HashedNode:
		if resolver == nil {
			return nil, fmt.Errorf("hashed node %x at depth %d along stem %x could not be resolved: %w", child.Commitment().Bytes(), n.depth, stem, errReadFromInvalid)
		}
		serialized, err := resolver(stem[:n.depth+1])
		if err != nil {
			return nil, fmt.Errorf("resolving node %x at depth %d: %w", stem, n.depth, err)
		}
		resolved, err := ParseNode(serialized, n.depth+1)
		if err != nil {
			return nil, fmt.Errorf("verkle tree: error parsing resolved node %x: %w", stem, err)
		}
		n.children[nchild] = resolved
		// recurse to handle the case of a LeafNode child that
		// splits.
		return n.GetValuesAtStem(stem, resolver)
	case *LeafNode:
		if equalPaths(child.stem, stem) {
			// We can't return the values since it's a POA leaf node, so we know nothing
			// about its values.
			if child.isPOAStub {
				return nil, errIsPOAStub
			}
			return child.values, nil
		}
		return nil, nil
	case *InternalNode:
		return child.GetValuesAtStem(stem, resolver)
	default:
		return nil, errUnknownNodeType
	}
}

func (n *InternalNode) Delete(key []byte, resolver NodeResolverFn) (bool, error) {
	nChild := offset2key(key, n.depth)
	switch child := n.children[nChild].(type) {
	case Empty:
		return false, nil
	case HashedNode:
		if resolver == nil {
			return false, errDeleteHash
		}
		payload, err := resolver(key[:n.depth+1])
		if err != nil {
			return false, err
		}
		// deserialize the payload and set it as the child
		c, err := ParseNode(payload, n.depth+1)
		if err != nil {
			return false, err
		}
		n.children[nChild] = c
		return n.Delete(key, resolver)
	default:
		n.cowChild(nChild)
		del, err := child.Delete(key, resolver)
		if err != nil {
			return false, err
		}

		// delete the entire child if instructed to by
		// the recursive algorigthm.
		if del {
			n.children[nChild] = Empty{}

			// Check if all children are gone, if so
			// signal that this node should be deleted
			// as well.
			for _, c := range n.children {
				if _, ok := c.(Empty); !ok {
					break
				}
			}

			return true, nil
		}

		return false, nil
	}
}

// Flush hashes the children of an internal node and replaces them
// with HashedNode. It also sends the current node on the flush channel.
func (n *InternalNode) Flush(flush NodeFlushFn) {
	//
	var (
		path                []byte
		flushAndCapturePath = func(p []byte, vn VerkleNode) {
			// wrap the flush function into a function that will
			// capture the path of the first leaf being flushed,
			// so that it can be used as the path to call `flush`.
			if len(path) == 0 {
				path = p[:n.depth]
			}
			flush(p, vn)
		}
	)

	n.Commit()
	for i, child := range n.children {
		if c, ok := child.(*InternalNode); ok {
			c.Commit()
			c.Flush(flushAndCapturePath)
			n.children[i] = HashedNode{}
		} else if c, ok := child.(*LeafNode); ok {
			c.Commit()
			flushAndCapturePath(c.stem[:n.depth+1], n.children[i])
			n.children[i] = HashedNode{}
		}
	}
	flush(path, n)
}

// FlushAtDepth goes over all internal nodes of a given depth, and
// flushes them to disk. Its purpose it to free up space if memory
// is running scarce.
func (n *InternalNode) FlushAtDepth(depth uint8, flush NodeFlushFn) {
	for i, child := range n.children {
		// Skip non-internal nodes
		c, ok := child.(*InternalNode)
		if !ok {
			if c, ok := child.(*LeafNode); ok {
				c.Commit()
				flush(c.stem[:c.depth], c)
				n.children[i] = HashedNode{}
			}
			continue
		}

		// Not deep enough, recurse
		if n.depth < depth {
			c.FlushAtDepth(depth, flush)
			continue
		}

		child.Commit()
		c.Flush(flush)
		n.children[i] = HashedNode{}
	}
}

func (n *InternalNode) Get(key []byte, resolver NodeResolverFn) ([]byte, error) {
	if len(key) != StemSize+1 {
		return nil, fmt.Errorf("invalid key length, expected %d, got %d", StemSize+1, len(key))
	}
	stemValues, err := n.GetValuesAtStem(key[:StemSize], resolver)
	if err != nil {
		return nil, err
	}

	// If the stem results in an empty node, return nil.
	if stemValues == nil {
		return nil, nil
	}

	// Return nil as a signal that the value isn't
	// present in the tree. This matches the behavior
	// of SecureTrie in Geth.
	return stemValues[key[StemSize]], nil
}

func (n *InternalNode) Hash() *Fr {
	var hash Fr
	n.Commitment().MapToScalarField(&hash)
	return &hash
}

func (n *InternalNode) Commitment() *Point {
	if n.commitment == nil {
		panic("nil commitment")
	}
	return n.commitment
}

func (n *InternalNode) fillLevels(levels [][]*InternalNode) {
	levels[int(n.depth)] = append(levels[int(n.depth)], n)
	for idx := range n.cow {
		child := n.children[idx]
		if childInternalNode, ok := child.(*InternalNode); ok && len(childInternalNode.cow) > 0 {
			childInternalNode.fillLevels(levels)
		}
	}
}

func (n *InternalNode) Commit() *Point {
	if len(n.cow) == 0 {
		return n.commitment
	}

	internalNodeLevels := make([][]*InternalNode, StemSize)
	n.fillLevels(internalNodeLevels)

	for level := len(internalNodeLevels) - 1; level >= 0; level-- {
		nodes := internalNodeLevels[level]
		if len(nodes) == 0 {
			continue
		}

		minBatchSize := 4
		if len(nodes) <= minBatchSize {
			if err := commitNodesAtLevel(nodes); err != nil {
				// TODO: make Commit() return an error
				panic(err)
			}
		} else {
			var wg sync.WaitGroup
			numBatches := runtime.NumCPU()
			batchSize := (len(nodes) + numBatches - 1) / numBatches
			if batchSize < minBatchSize {
				batchSize = minBatchSize
			}
			for i := 0; i < len(nodes); i += batchSize {
				start := i
				end := i + batchSize
				if end > len(nodes) {
					end = len(nodes)
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := commitNodesAtLevel(nodes[start:end]); err != nil {
						// TODO: make Commit() return an error
						panic(err)
					}

				}()
			}
			wg.Wait()
		}
	}
	return n.commitment
}

func commitNodesAtLevel(nodes []*InternalNode) error {
	points := make([]*Point, 0, 1024)
	cowIndexes := make([]int, 0, 1024)

	// For each internal node, we collect in `points` all the ones we need to map to a field element.
	// That is, for each touched children in a node, we collect the old and new commitment to do the diff updating
	// later.
	for _, node := range nodes {
		for idx, nodeChildComm := range node.cow {
			points = append(points, nodeChildComm)
			points = append(points, node.children[idx].Commitment())
			cowIndexes = append(cowIndexes, int(idx))
		}
	}

	// We generate `frs` which will contain the result for each element in `points`.
	frs := make([]*Fr, len(points))
	for i := range frs {
		frs[i] = &Fr{}
	}

	// Do a single batch calculation for all the points in this level.
	if err := banderwagon.BatchMapToScalarField(frs, points); err != nil {
		return fmt.Errorf("batch mapping to scalar fields: %s", err)
	}

	// We calculate the difference between each (new commitment - old commitment) pair, and store it
	// in the same slice to avoid allocations.
	for i := 0; i < len(frs); i += 2 {
		frs[i/2].Sub(frs[i+1], frs[i])
	}
	// Now `frs` have half of the elements, and these are the Frs differences to update commitments.
	frs = frs[:len(frs)/2]

	// Now we iterate on the nodes, and use this calculated differences to update their commitment.
	var frsIdx int
	var cowIndex int

	for _, node := range nodes {
		poly := make([]Fr, NodeWidth)
		for i := 0; i < len(node.cow); i++ {
			poly[cowIndexes[cowIndex]] = *frs[frsIdx]
			frsIdx++
			cowIndex++
		}
		node.cow = nil
		node.commitment.Add(node.commitment, cfg.CommitToPoly(poly, 0))
	}

	return nil
}

// groupKeys groups a set of keys based on their byte at a given depth.
func groupKeys(keys keylist, depth byte) []keylist {
	// special case: no key
	if len(keys) == 0 {
		return []keylist{}
	}

	// special case: only one key left
	if len(keys) == 1 {
		return []keylist{keys}
	}

	// there are at least two keys left in the list at this depth
	groups := make([]keylist, 0, len(keys))
	firstkey, lastkey := 0, 1
	for ; lastkey < len(keys); lastkey++ {
		key := keys[lastkey]
		keyidx := offset2key(key, depth)
		previdx := offset2key(keys[lastkey-1], depth)

		if keyidx != previdx {
			groups = append(groups, keys[firstkey:lastkey])
			firstkey = lastkey
		}
	}

	groups = append(groups, keys[firstkey:lastkey])

	return groups
}

func (n *InternalNode) GetProofItems(keys keylist, resolver NodeResolverFn) (*ProofElements, []byte, [][]byte, error) {
	var (
		groups = groupKeys(keys, n.depth)
		pe     = &ProofElements{
			Cis:    []*Point{},
			Zis:    []byte{},
			Yis:    []*Fr{}, // Should be 0
			Fis:    [][]Fr{},
			ByPath: map[string]*Point{},
		}

		esses []byte   = nil // list of extension statuses
		poass [][]byte       // list of proof-of-absence stems
	)

	// fill in the polynomial for this node
	var fi [NodeWidth]Fr
	var fiPtrs [NodeWidth]*Fr
	var points [NodeWidth]*Point
	for i, child := range n.children {
		fiPtrs[i] = &fi[i]
		if child != nil {
			var c VerkleNode
			if _, ok := child.(HashedNode); ok {
				childpath := make([]byte, n.depth+1)
				copy(childpath[:n.depth+1], keys[0][:n.depth])
				childpath[n.depth] = byte(i)
				if resolver == nil {
					return nil, nil, nil, fmt.Errorf("no resolver for path %x", childpath)
				}
				serialized, err := resolver(childpath)
				if err != nil {
					return nil, nil, nil, fmt.Errorf("error resolving for path %x: %w", childpath, err)
				}
				c, err = ParseNode(serialized, n.depth+1)
				if err != nil {
					return nil, nil, nil, err
				}
				n.children[i] = c
			} else {
				c = child
			}
			points[i] = c.Commitment()
		} else {
			// TODO: add a test case to cover this scenario.
			points[i] = new(Point)
		}
	}
	if err := banderwagon.BatchMapToScalarField(fiPtrs[:], points[:]); err != nil {
		return nil, nil, nil, fmt.Errorf("batch mapping to scalar fields: %s", err)
	}

	for _, group := range groups {
		childIdx := offset2key(group[0], n.depth)

		// Build the list of elements for this level
		var yi Fr
		yi.Set(&fi[childIdx])

		pe.Cis = append(pe.Cis, n.commitment)
		pe.Zis = append(pe.Zis, childIdx)
		pe.Yis = append(pe.Yis, &yi)
		pe.Fis = append(pe.Fis, fi[:])
		pe.ByPath[string(group[0][:n.depth])] = n.commitment
	}

	// Loop over again, collecting the children's proof elements
	// This is because the order is breadth-first.
	for _, group := range groups {
		childIdx := offset2key(group[0], n.depth)

		if _, isunknown := n.children[childIdx].(UnknownNode); isunknown {
			// TODO: add a test case to cover this scenario.
			return nil, nil, nil, errMissingNodeInStateless
		}

		// Special case of a proof of absence: no children
		// commitment, as the value is 0.
		_, isempty := n.children[childIdx].(Empty)
		if isempty {
			// A question arises here: what if this proof of absence
			// corresponds to several stems? Should the ext status be
			// repeated as many times? It would be wasteful, so the
			// decoding code has to be aware of this corner case.
			esses = append(esses, extStatusAbsentEmpty|((n.depth+1)<<3))
			for i := 0; i < len(group); i++ {
				// Append one nil value per key in this missing stem
				pe.Vals = append(pe.Vals, nil)
			}
			continue
		}

		pec, es, other, err := n.children[childIdx].GetProofItems(group, resolver)
		if err != nil {
			// TODO: add a test case to cover this scenario.
			return nil, nil, nil, err
		}
		pe.Merge(pec)
		poass = append(poass, other...)
		esses = append(esses, es...)
	}

	return pe, esses, poass, nil
}

// Serialize returns the serialized form of the internal node.
// The format is: <nodeType><bitlist><commitment>
func (n *InternalNode) Serialize() ([]byte, error) {
	ret := make([]byte, nodeTypeSize+bitlistSize+banderwagon.UncompressedSize)

	// Write the <bitlist>.
	bitlist := ret[internalBitlistOffset:internalCommitmentOffset]
	for i, c := range n.children {
		if _, ok := c.(Empty); !ok {
			setBit(bitlist, i)
		}
	}

	// Write the <node-type>
	ret[nodeTypeOffset] = internalRLPType

	// Write the <commitment>
	comm := n.commitment.BytesUncompressed()
	copy(ret[internalCommitmentOffset:], comm[:])

	return ret, nil
}

func (n *InternalNode) Copy() VerkleNode {
	ret := &InternalNode{
		children:   make([]VerkleNode, len(n.children)),
		commitment: new(Point),
		depth:      n.depth,
	}

	for i, child := range n.children {
		ret.children[i] = child.Copy()
	}

	if n.commitment != nil {
		ret.commitment.Set(n.commitment)
	}

	if n.cow != nil {
		ret.cow = make(map[byte]*Point)
		for k, v := range n.cow {
			ret.cow[k] = new(Point)
			ret.cow[k].Set(v)
		}
	}

	return ret
}

func (n *InternalNode) toDot(parent, path string) string {
	me := fmt.Sprintf("internal%s", path)
	var hash Fr
	n.commitment.MapToScalarField(&hash)
	ret := fmt.Sprintf("%s [label=\"I: %x\"]\n", me, hash.BytesLE())
	if len(parent) > 0 {
		ret = fmt.Sprintf("%s %s -> %s\n", ret, parent, me)
	}

	for i, child := range n.children {
		if child == nil {
			continue
		}
		ret = fmt.Sprintf("%s%s", ret, child.toDot(me, fmt.Sprintf("%s%02x", path, i)))
	}

	return ret
}

func (n *InternalNode) setDepth(d byte) {
	n.depth = d
}

// MergeTrees takes a series of subtrees that got filled following
// a command-and-conquer method, and merges them into a single tree.
// This method is deprecated, use with caution.
func MergeTrees(subroots []*InternalNode) VerkleNode {
	root := New().(*InternalNode)
	for _, subroot := range subroots {
		for i := 0; i < NodeWidth; i++ {
			if _, ok := subroot.children[i].(Empty); ok {
				continue
			}
			root.touchCoW(byte(i))
			root.children[i] = subroot.children[i]
		}
	}

	return root
}

// TouchCoW is a helper function that will mark a child as
// "inserted into". It is used by the conversion code to
// mark reconstructed subtrees as 'written to', so that its
// root commitment can be computed.
func (n *InternalNode) touchCoW(index byte) {
	n.cowChild(index)
}

func (n *LeafNode) Insert(key []byte, value []byte, _ NodeResolverFn) error {
	if n.isPOAStub {
		return errIsPOAStub
	}

	if len(key) != StemSize+1 {
		return fmt.Errorf("invalid key size: %d", len(key))
	}
	if !bytes.Equal(key[:StemSize], n.stem) {
		return fmt.Errorf("stems doesn't match: %x != %x", key[:StemSize], n.stem)
	}
	values := make([][]byte, NodeWidth)
	values[key[StemSize]] = value
	return n.insertMultiple(key[:StemSize], values)
}

func (n *LeafNode) insertMultiple(stem []byte, values [][]byte) error {
	// Sanity check: ensure the stems are the same.
	if !equalPaths(stem, n.stem) {
		return errInsertIntoOtherStem
	}

	return n.updateMultipleLeaves(values)
}

func (n *LeafNode) updateC(cxIndex int, newC Fr, oldC Fr) {
	// Calculate the Fr-delta.
	var deltaC Fr
	deltaC.Sub(&newC, &oldC)

	// Calculate the Point-delta.
	var poly [NodeWidth]Fr
	poly[cxIndex] = deltaC

	// Add delta to the current commitment.
	n.commitment.Add(n.commitment, cfg.CommitToPoly(poly[:], 0))
}

func (n *LeafNode) updateCn(index byte, value []byte, c *Point) error {
	var (
		old, newH [2]Fr
		diff      Point
		poly      [NodeWidth]Fr
	)

	// Optimization idea:
	// If the value is created (i.e. not overwritten), the leaf marker
	// is already present in the commitment. In order to save computations,
	// do not include it. The result should be the same,
	// but the computation time should be faster as one doesn't need to
	// compute 1 - 1 mod N.
	err := leafToComms(old[:], n.values[index])
	if err != nil {
		return err
	}
	if err = leafToComms(newH[:], value); err != nil {
		return err
	}

	newH[0].Sub(&newH[0], &old[0])
	poly[2*(index%128)] = newH[0]
	diff = cfg.conf.Commit(poly[:])
	poly[2*(index%128)].SetZero()
	c.Add(c, &diff)

	newH[1].Sub(&newH[1], &old[1])
	poly[2*(index%128)+1] = newH[1]
	diff = cfg.conf.Commit(poly[:])
	c.Add(c, &diff)

	return nil
}

func (n *LeafNode) updateLeaf(index byte, value []byte) error {
	// Update the corresponding C1 or C2 commitment.
	var c *Point
	var oldC Point
	if index < NodeWidth/2 {
		c = n.c1
		oldC = *n.c1
	} else {
		c = n.c2
		oldC = *n.c2
	}
	if err := n.updateCn(index, value, c); err != nil {
		return err
	}

	// Batch the Fr transformation of the new and old CX.
	var frs [2]Fr
	if err := banderwagon.BatchMapToScalarField([]*Fr{&frs[0], &frs[1]}, []*Point{c, &oldC}); err != nil {
		return fmt.Errorf("batch mapping to scalar fields: %s", err)
	}

	// If index is in the first NodeWidth/2 elements, we need to update C1. Otherwise, C2.
	cxIndex := 2 + int(index)/(NodeWidth/2) // [1, stem, -> C1, C2 <-]
	n.updateC(cxIndex, frs[0], frs[1])

	n.values[index] = value
	return nil
}

func (n *LeafNode) updateMultipleLeaves(values [][]byte) error { // skipcq: GO-R1005
	var oldC1, oldC2 *Point

	// We iterate the values, and we update the C1 and/or C2 commitments depending on the index.
	// If any of them is touched, we save the original point so we can update the LeafNode root
	// commitment. We copy the original point in oldC1 and oldC2, so we can batch their Fr transformation
	// after this loop.
	for i, v := range values {
		if len(v) != 0 && !bytes.Equal(v, n.values[i]) {
			if i < NodeWidth/2 {
				// First time we touch C1? Save the original point for later.
				if oldC1 == nil {
					oldC1 = &Point{}
					oldC1.Set(n.c1)
				}
				// We update C1 directly in `n`. We have our original copy in oldC1.
				if err := n.updateCn(byte(i), v, n.c1); err != nil {
					return err
				}
			} else {
				// First time we touch C2? Save the original point for later.
				if oldC2 == nil {
					oldC2 = &Point{}
					oldC2.Set(n.c2)
				}
				// We update C2 directly in `n`. We have our original copy in oldC2.
				if err := n.updateCn(byte(i), v, n.c2); err != nil {
					return err
				}
			}
			n.values[i] = v
		}
	}

	// We have three potential cases here:
	// 1. We have touched C1 and C2: we Fr-batch old1, old2 and newC1, newC2. (4x gain ratio)
	// 2. We have touched only one CX: we Fr-batch oldX and newCX. (2x gain ratio)
	// 3. No C1 or C2 was touched, this is a noop.
	var frs [4]Fr
	const c1Idx = 2 // [1, stem, ->C1<-, C2]
	const c2Idx = 3 // [1, stem, C1, ->C2<-]

	if oldC1 != nil && oldC2 != nil { // Case 1.
		if err := banderwagon.BatchMapToScalarField([]*Fr{&frs[0], &frs[1], &frs[2], &frs[3]}, []*Point{n.c1, oldC1, n.c2, oldC2}); err != nil {
			return fmt.Errorf("batch mapping to scalar fields: %s", err)
		}
		n.updateC(c1Idx, frs[0], frs[1])
		n.updateC(c2Idx, frs[2], frs[3])
	} else if oldC1 != nil { // Case 2. (C1 touched)
		if err := banderwagon.BatchMapToScalarField([]*Fr{&frs[0], &frs[1]}, []*Point{n.c1, oldC1}); err != nil {
			return fmt.Errorf("batch mapping to scalar fields: %s", err)
		}
		n.updateC(c1Idx, frs[0], frs[1])
	} else if oldC2 != nil { // Case 2. (C2 touched)
		if err := banderwagon.BatchMapToScalarField([]*Fr{&frs[0], &frs[1]}, []*Point{n.c2, oldC2}); err != nil {
			return fmt.Errorf("batch mapping to scalar fields: %s", err)
		}
		n.updateC(c2Idx, frs[0], frs[1])
	}

	return nil
}

// Delete deletes a value from the leaf, return `true` as a second
// return value, if the parent should entirely delete the child.
func (n *LeafNode) Delete(k []byte, _ NodeResolverFn) (bool, error) {
	// Sanity check: ensure the key header is the same:
	if !equalPaths(k, n.stem) {
		return false, nil
	}

	// Erase the value it used to contain
	original := n.values[k[31]] // save original value
	n.values[k[31]] = nil

	// Check if a Cn subtree is entirely empty, or if
	// the entire subtree is empty.
	var (
		isCnempty = true
		isCempty  = true
	)
	for i := 0; i < NodeWidth; i++ {
		if len(n.values[i]) > 0 {
			// if i and k[31] are in the same subtree,
			// set both values and return.
			if byte(i/128) == k[31]/128 {
				isCnempty = false
				isCempty = false
				break
			}

			// i and k[31] were in a different subtree,
			// so all we can say at this stage, is that
			// the whole tree isn't empty.
			// TODO if i < 128, then k[31] >= 128 and
			// we could skip to 128, but that's an
			// optimization for later.
			isCempty = false
		}
	}

	// if the whole subtree is empty, then the
	// entire node should be deleted.
	if isCempty {
		return true, nil
	}

	// if a Cn branch becomes empty as a result
	// of removing the last value, update C by
	// adding -Cn to it and exit.
	if isCnempty {
		var (
			cn           *Point
			subtreeindex = 2 + k[31]/128
		)

		if k[31] < 128 {
			cn = n.c1
		} else {
			cn = n.c2
		}

		// Update C by subtracting the old value for Cn
		// Note: this isn't done in one swoop, which would make sense
		// since presumably a lot of values would be deleted at the same
		// time when reorging. Nonetheless, a reorg is an already complex
		// operation which is slow no matter what, so ensuring correctness
		// is more important than
		var poly [4]Fr
		cn.MapToScalarField(&poly[subtreeindex])
		n.commitment.Sub(n.commitment, cfg.CommitToPoly(poly[:], 0))

		// Clear the corresponding commitment
		if k[31] < 128 {
			n.c1 = nil
		} else {
			n.c2 = nil
		}

		return false, nil
	}

	// Recompute the updated C & Cn
	//
	// This is done by setting the leaf value
	// to `nil` at this point, and all the
	// diff computation will be performed by
	// updateLeaf since leafToComms supports
	// nil values.
	// Note that the value is set to nil by
	// the method, as it needs the original
	// value to compute the commitment diffs.
	n.values[k[31]] = original
	return false, n.updateLeaf(k[31], nil)
}

func (n *LeafNode) Get(k []byte, _ NodeResolverFn) ([]byte, error) {
	if n.isPOAStub {
		return nil, errIsPOAStub
	}

	if !equalPaths(k, n.stem) {
		// If keys differ, return nil in order to
		// signal that the key isn't present in the
		// tree. Do not return an error, thus matching
		// the behavior of Geth's SecureTrie.
		return nil, nil
	}
	// value can be nil, as expected by geth
	return n.values[k[StemSize]], nil
}

func (n *LeafNode) Hash() *Fr {
	// TODO cache this in a subsequent PR, not done here
	// to reduce complexity.
	// TODO use n.commitment once all Insert* are diff-inserts
	var hash Fr
	n.Commitment().MapToScalarField(&hash)
	return &hash
}

func (n *LeafNode) Commitment() *Point {
	if n.commitment == nil {
		panic("nil commitment")
	}
	return n.commitment
}

func (n *LeafNode) Commit() *Point {
	return n.commitment
}

// fillSuffixTreePoly takes one of the two suffix tree and
// builds the associated polynomial, to be used to compute
// the corresponding C{1,2} commitment.
func fillSuffixTreePoly(poly []Fr, values [][]byte) (int, error) {
	count := 0
	for idx, val := range values {
		if val == nil {
			continue
		}
		count++

		if err := leafToComms(poly[(idx<<1)&0xFF:], val); err != nil {
			return 0, err
		}
	}
	return count, nil
}

// leafToComms turns a leaf into two commitments of the suffix
// and extension tree.
func leafToComms(poly []Fr, val []byte) error {
	if len(val) == 0 {
		return nil
	}
	if len(val) > 32 {
		return fmt.Errorf("invalid leaf length %d, %v", len(val), val)
	}
	var (
		valLoWithMarker [17]byte
		loEnd           = 16
		err             error
	)
	if len(val) < loEnd {
		loEnd = len(val)
	}
	copy(valLoWithMarker[:loEnd], val[:loEnd])
	valLoWithMarker[16] = 1 // 2**128
	if err = FromLEBytes(&poly[0], valLoWithMarker[:]); err != nil {
		return err
	}
	if len(val) >= 16 {
		if err = FromLEBytes(&poly[1], val[16:]); err != nil {
			return err
		}
	}
	return nil
}

func (n *LeafNode) GetProofItems(keys keylist, _ NodeResolverFn) (*ProofElements, []byte, [][]byte, error) { // skipcq: GO-R1005
	var (
		poly [NodeWidth]Fr // top-level polynomial
		pe                 = &ProofElements{
			Cis:    []*Point{n.commitment, n.commitment},
			Zis:    []byte{0, 1},
			Yis:    []*Fr{&poly[0], &poly[1]}, // Should be 0
			Fis:    [][]Fr{poly[:], poly[:]},
			Vals:   make([][]byte, 0, len(keys)),
			ByPath: map[string]*Point{},
		}

		esses []byte   = nil // list of extension statuses
		poass [][]byte       // list of proof-of-absence stems
	)

	// Initialize the top-level polynomial with 1 + stem + C1 + C2
	poly[0].SetUint64(1)
	if err := StemFromBytes(&poly[1], n.stem); err != nil {
		return nil, nil, nil, fmt.Errorf("error serializing stem '%x': %w", n.stem, err)
	}

	// First pass: add top-level elements first
	var hasC1, hasC2 bool
	for _, key := range keys {
		// Note that keys might contain keys that don't correspond to this leaf node.
		// We should only analize the inclusion of C1/C2 for keys corresponding to this
		// leaf node stem.
		if equalPaths(n.stem, key) {
			hasC1 = hasC1 || (key[31] < 128)
			hasC2 = hasC2 || (key[31] >= 128)
			if hasC2 {
				break
			}
		}
	}

	// If this tree is a full tree (i.e: not a stateless tree), we know we have c1 and c2 values.
	// Also, we _need_ them independently of hasC1 or hasC2 since the prover needs `Fis`.
	if !n.isPOAStub {
		if err := banderwagon.BatchMapToScalarField([]*Fr{&poly[2], &poly[3]}, []*Point{n.c1, n.c2}); err != nil {
			return nil, nil, nil, fmt.Errorf("batch mapping to scalar fields: %s", err)
		}
	} else if hasC1 || hasC2 || n.c1 != nil || n.c2 != nil {
		// This LeafNode is a proof of absence stub. It must be true that
		// both c1 and c2 are nil, and that hasC1 and hasC2 are false.
		// Let's just check that to be sure, since the code below can't use
		// poly[2] or poly[3].
		return nil, nil, nil, fmt.Errorf("invalid proof of absence stub")
	}

	if hasC1 {
		pe.Cis = append(pe.Cis, n.commitment)
		pe.Zis = append(pe.Zis, 2)
		pe.Yis = append(pe.Yis, &poly[2])
		pe.Fis = append(pe.Fis, poly[:])
	}
	if hasC2 {
		pe.Cis = append(pe.Cis, n.commitment)
		pe.Zis = append(pe.Zis, 3)
		pe.Yis = append(pe.Yis, &poly[3])
		pe.Fis = append(pe.Fis, poly[:])
	}

	// Second pass: add the cn-level elements
	for _, key := range keys {
		pe.ByPath[string(key[:n.depth])] = n.commitment

		// Proof of absence: case of a differing stem.
		// Add an unopened stem-level node.
		if !equalPaths(n.stem, key) {
			// Corner case: don't add the poa stem if it's
			// already present as a proof-of-absence for a
			// different key, or for the same key (case of
			// multiple missing keys being absent).
			// The list of extension statuses has to be of
			// length 1 at this level, so skip otherwise.
			if len(esses) == 0 {
				esses = append(esses, extStatusAbsentOther|(n.depth<<3))
				poass = append(poass, n.stem)
			}
			pe.Vals = append(pe.Vals, nil)
			continue
		}

		// corner case (see previous corner case): if a proof-of-absence
		// stem was found, and it now turns out the same stem is used as
		// a proof of presence, clear the proof-of-absence list to avoid
		// redundancy.
		if len(poass) > 0 {
			poass = nil
			esses = nil
		}

		var (
			suffix   = key[31]
			suffPoly [NodeWidth]Fr // suffix-level polynomial
			count    int
			err      error
		)
		if suffix >= 128 {
			count, err = fillSuffixTreePoly(suffPoly[:], n.values[128:])
			if err != nil {
				return nil, nil, nil, err
			}
		} else {
			count, err = fillSuffixTreePoly(suffPoly[:], n.values[:128])
			if err != nil {
				return nil, nil, nil, err
			}
		}

		// Proof of absence: case of a missing suffix tree.
		//
		// The suffix tree for this value is missing, i.e. all
		// values in the extension-and-suffix tree are grouped
		// in the other suffix tree (e.g. C2 if we are looking
		// at C1).
		if count == 0 {
			// TODO(gballet) maintain a count variable at LeafNode level
			// so that we know not to build the polynomials in this case,
			// as all the information is available before fillSuffixTreePoly
			// has to be called, save the count.
			esses = append(esses, extStatusAbsentEmpty|(n.depth<<3))
			pe.Vals = append(pe.Vals, nil)
			continue
		}

		var scomm *Point
		if suffix < 128 {
			scomm = n.c1
		} else {
			scomm = n.c2
		}
		var leaves [2]Fr
		if n.values[suffix] == nil {
			// Proof of absence: case of a missing value.
			//
			// Suffix tree is present as a child of the extension,
			// but does not contain the requested suffix. This can
			// only happen when the leaf has never been written to
			// since after deletion the value would be set to zero
			// but still contain the leaf marker 2^128.
			leaves[0], leaves[1] = FrZero, FrZero
		} else {
			// suffix tree is present and contains the key
			leaves[0], leaves[1] = suffPoly[2*suffix], suffPoly[2*suffix+1]
		}
		pe.Cis = append(pe.Cis, scomm, scomm)
		pe.Zis = append(pe.Zis, 2*suffix, 2*suffix+1)
		pe.Yis = append(pe.Yis, &leaves[0], &leaves[1])
		pe.Fis = append(pe.Fis, suffPoly[:], suffPoly[:])
		pe.Vals = append(pe.Vals, n.values[key[31]])
		if len(esses) == 0 || esses[len(esses)-1] != extStatusPresent|(n.depth<<3) {
			esses = append(esses, extStatusPresent|(n.depth<<3))
		}
		slotPath := string(key[:n.depth]) + string([]byte{2 + suffix/128})
		pe.ByPath[slotPath] = scomm
	}

	return pe, esses, poass, nil
}

// Serialize serializes a LeafNode.
// The format is: <nodeType><stem><bitlist><comm><c1comm><c2comm><children...>
func (n *LeafNode) Serialize() ([]byte, error) {
	cBytes := banderwagon.BatchToBytesUncompressed(n.commitment, n.c1, n.c2)
	return n.serializeLeafWithUncompressedCommitments(cBytes[0], cBytes[1], cBytes[2]), nil
}

func (n *LeafNode) Copy() VerkleNode {
	l := &LeafNode{}
	l.stem = make([]byte, len(n.stem))
	l.values = make([][]byte, len(n.values))
	l.depth = n.depth
	copy(l.stem, n.stem)
	for i, v := range n.values {
		l.values[i] = make([]byte, len(v))
		copy(l.values[i], v)
	}
	if n.commitment != nil {
		l.commitment = new(Point)
		l.commitment.Set(n.commitment)
	}
	if n.c1 != nil {
		l.c1 = new(Point)
		l.c1.Set(n.c1)
	}
	if n.c2 != nil {
		l.c2 = new(Point)
		l.c2.Set(n.c2)
	}
	l.isPOAStub = n.isPOAStub

	return l
}

func (n *LeafNode) Key(i int) []byte {
	var ret [32]byte
	copy(ret[:], n.stem)
	ret[31] = byte(i)
	return ret[:]
}

func (n *LeafNode) Value(i int) []byte {
	if i >= NodeWidth {
		panic("leaf node index out of range")
	}
	return n.values[byte(i)]
}

func (n *LeafNode) toDot(parent, path string) string {
	var hash Fr
	n.Commitment().MapToScalarField(&hash)
	ret := fmt.Sprintf("leaf%s [label=\"L: %x\nC: %x\nC₁: %x\nC₂:%x\"]\n%s -> leaf%s\n", path, hash.Bytes(), n.commitment.Bytes(), n.c1.Bytes(), n.c2.Bytes(), parent, path)
	for i, v := range n.values {
		if len(v) != 0 {
			ret = fmt.Sprintf("%sval%s%02x [label=\"%x\"]\nleaf%s -> val%s%02x\n", ret, path, i, v, path, path, i)
		}
	}
	return ret
}

func (n *LeafNode) setDepth(d byte) {
	n.depth = d
}

func (n *LeafNode) Values() [][]byte {
	return n.values
}

func setBit(bitlist []byte, index int) {
	bitlist[index/8] |= mask[index%8]
}

func ToDot(root VerkleNode) string {
	root.Commit()
	return fmt.Sprintf("digraph D {\n%s}", root.toDot("", ""))
}

// SerializedNode contains a serialization of a tree node.
// It provides everything that the client needs to save the node to the database.
// For example, CommitmentBytes is usually use as key and SerializedBytes as value.
// Providing both allows this library to do more optimizations.
type SerializedNode struct {
	Node            VerkleNode
	CommitmentBytes [banderwagon.UncompressedSize]byte
	Path            []byte
	SerializedBytes []byte
}

// BatchSerialize is an optimized serialization API when multiple VerkleNodes serializations are required, and all are
// available in memory.
func (n *InternalNode) BatchSerialize() ([]SerializedNode, error) {
	// Commit to the node to update all the nodes commitments.
	n.Commit()

	// Collect all nodes that we need to serialize.
	nodes := make([]VerkleNode, 0, 1024)
	paths := make([][]byte, 0, 1024)
	nodes, paths = n.collectNonHashedNodes(nodes, paths, nil)

	// We collect all the *Point, so we can batch all projective->affine transformations.
	pointsToCompress := make([]*Point, 0, 3*len(nodes))
	// Contains a map between VerkleNode and the index in the serializedPoints containing the commitment below.
	serializedPointsIdxs := make(map[VerkleNode]int, len(nodes))
	for i := range nodes {
		switch n := nodes[i].(type) {
		case *InternalNode:
			pointsToCompress = append(pointsToCompress, n.commitment)
			serializedPointsIdxs[n] = len(pointsToCompress) - 1
		case *LeafNode:
			pointsToCompress = append(pointsToCompress, n.commitment, n.c1, n.c2)
		}
	}

	// Now we do the all transformations in a single-shot.
	serializedPoints := banderwagon.BatchToBytesUncompressed(pointsToCompress...)

	// Now we that we did the heavy CPU work, we have to do the rest of `nodes` serialization
	// taking the compressed points from this single list.
	ret := make([]SerializedNode, 0, len(nodes))
	idx := 0
	for i := range nodes {
		switch n := nodes[i].(type) {
		case *InternalNode:
			serialized, err := n.serializeInternalWithUncompressedCommitment(serializedPointsIdxs, serializedPoints)
			if err != nil {
				return nil, err
			}
			sn := SerializedNode{
				Node:            n,
				Path:            paths[i],
				CommitmentBytes: serializedPoints[idx],
				SerializedBytes: serialized,
			}
			ret = append(ret, sn)
			idx++
		case *LeafNode:
			cBytes := serializedPoints[idx]
			c1Bytes := serializedPoints[idx+1]
			c2Bytes := serializedPoints[idx+2]
			sn := SerializedNode{
				Node:            n,
				Path:            paths[i],
				CommitmentBytes: serializedPoints[idx],
				SerializedBytes: n.serializeLeafWithUncompressedCommitments(cBytes, c1Bytes, c2Bytes),
			}
			ret = append(ret, sn)
			idx += 3
		}
	}

	return ret, nil
}

func (n *InternalNode) collectNonHashedNodes(list []VerkleNode, paths [][]byte, path []byte) ([]VerkleNode, [][]byte) {
	list = append(list, n)
	paths = append(paths, path)
	for i, child := range n.children {
		switch childNode := child.(type) {
		case *LeafNode:
			list = append(list, childNode)
			paths = append(paths, childNode.stem[:len(path)+1])
		case *InternalNode:
			childpath := make([]byte, len(path)+1)
			copy(childpath, path)
			childpath[len(path)] = byte(i)
			list, paths = childNode.collectNonHashedNodes(list, paths, childpath)
		}
	}
	return list, paths
}

// unpack one compressed commitment from the list of batch-compressed commitments
func (n *InternalNode) serializeInternalWithUncompressedCommitment(pointsIdx map[VerkleNode]int, serializedPoints [][banderwagon.UncompressedSize]byte) ([]byte, error) {
	serialized := make([]byte, nodeTypeSize+bitlistSize+banderwagon.UncompressedSize)
	bitlist := serialized[internalBitlistOffset:internalCommitmentOffset]
	for i, c := range n.children {
		if _, ok := c.(Empty); !ok {
			setBit(bitlist, i)
		}
	}
	serialized[nodeTypeOffset] = internalRLPType
	pointidx, ok := pointsIdx[n]
	if !ok {
		return nil, fmt.Errorf("child node not found in cache")
	}
	copy(serialized[internalCommitmentOffset:], serializedPoints[pointidx][:])

	return serialized, nil
}

func (n *LeafNode) serializeLeafWithUncompressedCommitments(cBytes, c1Bytes, c2Bytes [banderwagon.UncompressedSize]byte) []byte {
	// Empty value in LeafNode used for padding.
	var emptyValue [LeafValueSize]byte

	// Create bitlist and store in children LeafValueSize (padded) values.
	children := make([]byte, 0, NodeWidth*LeafValueSize)
	var bitlist [bitlistSize]byte
	for i, v := range n.values {
		if v != nil {
			setBit(bitlist[:], i)
			children = append(children, v...)
			if padding := emptyValue[:LeafValueSize-len(v)]; len(padding) != 0 {
				children = append(children, padding...)
			}
		}
	}

	// Create the serialization.
	result := make([]byte, nodeTypeSize+StemSize+bitlistSize+3*banderwagon.UncompressedSize+len(children))
	result[0] = leafRLPType
	copy(result[leafSteamOffset:], n.stem[:StemSize])
	copy(result[leafBitlistOffset:], bitlist[:])
	copy(result[leafCommitmentOffset:], cBytes[:])
	copy(result[leafC1CommitmentOffset:], c1Bytes[:])
	copy(result[leafC2CommitmentOffset:], c2Bytes[:])
	copy(result[leafChildrenOffset:], children)

	return result
}
