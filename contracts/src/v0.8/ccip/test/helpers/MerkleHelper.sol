// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {MerkleMultiProof} from "../../libraries/MerkleMultiProof.sol";

library MerkleHelper {
  /// @notice Generate a Merkle Root from a full set of leaves. When a tree is unbalanced
  /// the value is brought up in the tree. For example consider (a,b,c) as leaves. This would
  /// result in the following tree with d being computed from hash(a,c) and the root r from
  /// hash(d,c). Notice c is not being rehashed when it is brought up in the tree, so the
  /// root is NOT hash(d,hash(c)) but instead hash(d,c) == hash(hash(a,b),c).
  ///       r
  ///     /   \
  ///    d     c
  ///   / \
  ///  a   b
  function getMerkleRoot(
    bytes32[] memory hashedLeaves
  ) public pure returns (bytes32) {
    // solhint-disable-next-line reason-string,gas-custom-errors
    require(hashedLeaves.length <= 256);
    while (hashedLeaves.length > 1) {
      hashedLeaves = computeNextLayer(hashedLeaves);
    }
    return hashedLeaves[0];
  }

  /// @notice Computes a single layer of a merkle proof by hashing each pair (i, i+1) for
  /// each i, i+2, i+4.. n. When an uneven number of leaves is supplied the last item
  /// is simply included as the last element in the result set and not hashed.
  function computeNextLayer(
    bytes32[] memory layer
  ) public pure returns (bytes32[] memory) {
    uint256 leavesLen = layer.length;
    if (leavesLen == 1) return layer;

    unchecked {
      bytes32[] memory nextLayer = new bytes32[]((leavesLen + 1) / 2);
      for (uint256 i = 0; i < leavesLen; i += 2) {
        if (i == leavesLen - 1) {
          nextLayer[i / 2] = layer[i];
        } else {
          nextLayer[i / 2] = hashPair(layer[i], layer[i + 1]);
        }
      }
      return nextLayer;
    }
  }

  function hashPair(bytes32 a, bytes32 b) public pure returns (bytes32) {
    return a < b ? hashInternalNode(a, b) : hashInternalNode(b, a);
  }

  function hashInternalNode(bytes32 left, bytes32 right) public pure returns (bytes32 hash) {
    return keccak256(abi.encode(MerkleMultiProof.INTERNAL_DOMAIN_SEPARATOR, left, right));
  }
}
