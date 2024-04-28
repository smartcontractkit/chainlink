// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";

contract CapabilityRegistry is OwnerIsCreator, TypeAndVersionInterface {
  struct NodeOperator {
    /// @notice The address of the admin that can manage a node
    /// operator
    address admin;
    /// @notice Human readable name of a Node Operator managing the node
    string name;
  }

  struct Node {
    /// @notice The id of the node operator that manages this node
    uint256 nodeOperatorId;
    /// @notice The P2P ID used for OCR
    bytes p2pId;
    /// @notice The list of capability IDs
    /// this node supports
    string[] supportedCapabilityIds;
  }

  struct Capability {
    // Capability type, e.g. "data-streams-reports"
    // bytes32(string); validation regex: ^[a-z0-9_\-:]{1,32}$
    // Not "type" because that's a reserved keyword in Solidity.
    bytes32 capabilityType;
    // Semver, e.g., "1.2.3"
    // bytes32(string); must be valid Semver + max 32 characters.
    bytes32 version;
  }

  /// @notice This error is thrown when a caller tries to call a function
  /// it does not have authorization to call
  error AccessForbidden();

  /// @notice This error is thrown when trying to set a node operator's
  /// admin address to the zero address
  error InvalidNodeOperatorAdmin();

  /// @notice This error is thrown when trying to configure the node's
  /// p2p ID to an empty bytes
  error InvalidNodeP2PId();

  /// @notice This event is emitted when a new node is added
  /// @param nodeId The ID of the newly added node
  /// @param nodeOperatorId The ID of the node operator that manages this node
  /// @param p2pId The P2P ID the node will use for OCR
  event NodeAdded(uint256 nodeId, uint256 nodeOperatorId, bytes p2pId);

  /// @notice This event is emitted when a new node operator is added
  /// @param nodeOperatorId The ID of the newly added node operator
  /// @param admin The address of the admin that can manage the node
  /// operator
  /// @param name The human readable name of the node operator
  event NodeOperatorAdded(uint256 nodeOperatorId, address indexed admin, string name);

  /// @notice This event is emitted when a node operator is removed
  /// @param nodeOperatorId The ID of the node operator that was removed
  event NodeOperatorRemoved(uint256 nodeOperatorId);

  /// @notice This event is emitted when a new capability is added
  /// @param capabilityId The ID of the newly added capability
  event CapabilityAdded(bytes32 indexed capabilityId);

  mapping(bytes32 => Capability) private s_capabilities;

  /// @notice Mapping of node operators
  mapping(uint256 nodeOperatorId => NodeOperator nodeOperator) private s_nodeOperators;

  /// @notice Mapping of nodes
  mapping(uint256 nodeId => Node node) private s_nodes;

  /// @notice The latest node operator ID
  /// @dev No getter for this as this is an implementation detail
  uint256 private s_nodeOperatorId;

  /// @notice The latest node ID
  /// @dev No getter for this as this is an implementation detail
  uint256 private s_nodeId;

  function typeAndVersion() external pure override returns (string memory) {
    return "CapabilityRegistry 1.0.0";
  }

  /// @notice Adds a list of node operators
  /// @param nodeOperators List of node operators to add
  function addNodeOperators(NodeOperator[] calldata nodeOperators) external onlyOwner {
    for (uint256 i; i < nodeOperators.length; ++i) {
      NodeOperator memory nodeOperator = nodeOperators[i];
      if (nodeOperator.admin == address(0)) revert InvalidNodeOperatorAdmin();
      uint256 nodeOperatorId = s_nodeOperatorId;
      s_nodeOperators[nodeOperatorId] = NodeOperator({admin: nodeOperator.admin, name: nodeOperator.name});
      ++s_nodeOperatorId;
      emit NodeOperatorAdded(nodeOperatorId, nodeOperator.admin, nodeOperator.name);
    }
  }

  /// @notice Removes a node operator
  /// @param nodeOperatorIds The IDs of the node operators to remove
  function removeNodeOperators(uint256[] calldata nodeOperatorIds) external onlyOwner {
    for (uint256 i; i < nodeOperatorIds.length; ++i) {
      uint256 nodeOperatorId = nodeOperatorIds[i];
      delete s_nodeOperators[nodeOperatorId];
      emit NodeOperatorRemoved(nodeOperatorId);
    }
  }

  /// @notice Gets a node operator's data
  /// @param nodeOperatorId The ID of the node operator to query for
  /// @return NodeOperator The node operator data
  function getNodeOperator(uint256 nodeOperatorId) external view returns (NodeOperator memory) {
    return s_nodeOperators[nodeOperatorId];
  }

  /// @notice Adds nodes
  /// @param nodes The nodes to add
  function addNodes(Node[] calldata nodes) external {
    for (uint256 i; i < nodes.length; ++i) {
      Node memory node = nodes[i];
      if (bytes32(node.p2pId) == bytes32("")) revert InvalidNodeP2PId();
      NodeOperator memory nodeOperator = s_nodeOperators[node.nodeOperatorId];
      if (msg.sender != nodeOperator.admin) revert AccessForbidden();

      uint256 nodeId = s_nodeId;
      s_nodes[nodeId] = node;
      s_nodeId++;
      emit NodeAdded(nodeId, node.nodeOperatorId, node.p2pId);
    }
  }

  /// @notice Gets a node's data
  /// @param nodeId The ID of the node to query for
  /// @return Node The node data
  function getNode(uint256 nodeId) external view returns (Node memory) {
    return s_nodes[nodeId];
  }

  function addCapability(Capability calldata capability) external onlyOwner {
    bytes32 capabilityId = getCapabilityID(capability.capabilityType, capability.version);
    s_capabilities[capabilityId] = capability;
    emit CapabilityAdded(capabilityId);
  }

  function getCapability(bytes32 capabilityID) public view returns (Capability memory) {
    return s_capabilities[capabilityID];
  }

  /// @notice This functions returns a Capability ID packed into a bytes32 for cheaper access
  /// @return bytes32 A unique identifier for the capability
  function getCapabilityID(bytes32 capabilityType, bytes32 version) public pure returns (bytes32) {
    return keccak256(abi.encodePacked(capabilityType, version));
  }
}
