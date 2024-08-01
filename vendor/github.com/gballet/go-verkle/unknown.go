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

import "errors"

type UnknownNode struct{}

func (UnknownNode) Insert([]byte, []byte, NodeResolverFn) error {
	return errMissingNodeInStateless
}

func (UnknownNode) Delete([]byte, NodeResolverFn) (bool, error) {
	return false, errors.New("cant delete in a subtree missing form a stateless view")
}

func (UnknownNode) Get([]byte, NodeResolverFn) ([]byte, error) {
	return nil, nil
}

func (n UnknownNode) Commit() *Point {
	return n.Commitment()
}

func (UnknownNode) Commitment() *Point {
	var id Point
	id.SetIdentity()
	return &id
}

func (UnknownNode) GetProofItems(keylist, NodeResolverFn) (*ProofElements, []byte, [][]byte, error) {
	return nil, nil, nil, errors.New("can't generate proof items for unknown node")
}

func (UnknownNode) Serialize() ([]byte, error) {
	return nil, errors.New("trying to serialize a subtree missing from the statless view")
}

func (UnknownNode) Copy() VerkleNode {
	return UnknownNode(struct{}{})
}

func (UnknownNode) toDot(string, string) string {
	return ""
}

func (UnknownNode) setDepth(_ byte) {
	panic("should not be try to set the depth of an UnknownNode node")
}

func (UnknownNode) Hash() *Fr {
	return &FrZero
}
