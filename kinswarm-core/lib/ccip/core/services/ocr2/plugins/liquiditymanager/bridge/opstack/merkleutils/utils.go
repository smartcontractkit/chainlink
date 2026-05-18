package merkleutils

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/opstack/rlphelpers"
)

// MaybeAddProofNode is a fix for the case where the final proof element is less than 32 bytes and the element exists
// inside of a branch node. Current implementation of the onchain MPT contract can't handle this
// natively so we instead append an extra proof element to handle it instead.
// Implementation is ported from: https://github.com/ethereum-optimism/optimism/blob/53573e0ea6a807a125784cc5c7df07cbb4dbe3bc/packages/sdk/src/utils/merkle-utils.ts#L57.
func MaybeAddProofNode(
	key [32]byte,
	proof [][]byte,
) ([][]byte, error) {
	keyHex := hexutil.Encode(key[:])

	var modifiedProof [][]byte
	modifiedProof = append(modifiedProof, proof...)

	finalProofEl := modifiedProof[len(modifiedProof)-1]
	rlpBuffers := rlphelpers.NewRLPBuffers()

	err := rlp.DecodeBytes(finalProofEl, rlpBuffers)
	if err != nil {
		return nil, fmt.Errorf("unable to RLP decode proof element: %w", err)
	}

	if len(rlpBuffers.Children) == 17 {
		for _, item := range rlpBuffers.Children {
			// Find any nodes located inside of the branch node.
			if len(item.Children) > 0 {
				// Check if the key inside the node matches the key we're looking for. We remove the first
				// two characters (0x) and then we remove one more character (the first nibble) since this
				// is the identifier for the type of node we're looking at. In this case we don't actually
				// care what type of node it is because a branch node would only ever be the final proof
				// element if (1) it includes the leaf node we're looking for or (2) it stores the value
				// within itself. If (1) then this logic will work, if (2) then this won't find anything
				// and we won't append any proof elements, which is exactly what we would want.
				itemSuffix := hexutil.Encode(item.Children[0].Data)[3:]
				if strings.HasSuffix(keyHex, itemSuffix) {
					var itemTree [][]byte
					for _, child := range item.Children {
						itemTree = append(itemTree, child.Data)
					}
					encoded, err := rlp.EncodeToBytes(itemTree)
					if err != nil {
						return nil, fmt.Errorf("unable to RLP encode proof element: %w", err)
					}
					modifiedProof = append(modifiedProof, encoded)
				}
			}
		}
	}
	return modifiedProof, nil
}
