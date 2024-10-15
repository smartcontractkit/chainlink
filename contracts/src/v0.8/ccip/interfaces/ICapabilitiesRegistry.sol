// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

interface ICapabilitiesRegistry {
  struct NodeInfo {
    /// @notice The id of the node operator that manages this node
    uint32 nodeOperatorId;
    /// @notice The number of times the node's configuration has been updated
    uint32 configCount;
    /// @notice The ID of the Workflow DON that the node belongs to. A node can
    /// only belong to one DON that accepts Workflows.
    uint32 workflowDONId;
    /// @notice The signer address for application-layer message verification.
    bytes32 signer;
    /// @notice This is an Ed25519 public key that is used to identify a node.
    /// This key is guaranteed to be unique in the CapabilitiesRegistry. It is
    /// used to identify a node in the the P2P network.
    bytes32 p2pId;
    /// @notice Public key used to encrypt secrets for this node
    bytes32 encryptionPublicKey;
    /// @notice The list of hashed capability IDs supported by the node
    bytes32[] hashedCapabilityIds;
    /// @notice The list of capabilities DON Ids supported by the node. A node
    /// can belong to multiple capabilities DONs. This list does not include a
    /// Workflow DON id if the node belongs to one.
    uint256[] capabilitiesDONIds;
  }

  /// @notice Gets a node's data
  /// @param p2pId The P2P ID of the node to query for
  /// @return NodeInfo The node data
  function getNode(
    bytes32 p2pId
  ) external view returns (NodeInfo memory);
}
