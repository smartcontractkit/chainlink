// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

/// @title INodeInfoProvider
/// @notice Interface for retrieving node information.
interface INodeInfoProvider {
  /// @notice This error is thrown when a node with the provided P2P ID is
  /// not found.
  /// @param nodeP2PId The node P2P ID used for the lookup.
  error NodeDoesNotExist(bytes32 nodeP2PId);

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

  /// @notice Retrieves node information by its P2P ID.
  /// @param p2pId The P2P ID of the node to query for.
  /// @return nodeInfo The node data.
  function getNode(bytes32 p2pId) external view returns (NodeInfo memory nodeInfo);

  /// @notice Retrieves all node information.
  /// @return NodeInfo[] Array of all nodes in the registry.
  function getNodes() external view returns (NodeInfo[] memory);

  /// @notice Retrieves nodes by their P2P IDs.
  /// @param p2pIds Array of P2P IDs to query for.
  /// @return NodeInfo[] Array of node data corresponding to the provided P2P IDs.
  function getNodesByP2PIds(bytes32[] calldata p2pIds) external view returns (NodeInfo[] memory);
}
