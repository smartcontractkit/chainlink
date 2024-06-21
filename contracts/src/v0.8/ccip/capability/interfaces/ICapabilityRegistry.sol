// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

interface ICapabilityRegistry {
  struct NodeInfo {
    /// @notice The id of the node operator that manages this node
    uint32 nodeOperatorId;
    /// @notice The signer address for application-layer message verification.
    bytes32 signer;
    /// @notice This is an Ed25519 public key that is used to identify a node.
    /// This key is guaranteed to be unique in the CapabilityRegistry. It is
    /// used to identify a node in the the P2P network.
    bytes32 p2pId;
    /// @notice The list of hashed capability IDs supported by the node
    bytes32[] hashedCapabilityIds;
  }

  /// @notice Gets a node's data
  /// @param p2pId The P2P ID of the node to query for
  /// @return NodeInfo The node data
  /// @return configCount The number of times the node has been configured
  function getNode(bytes32 p2pId) external view returns (NodeInfo memory, uint32 configCount);
}
