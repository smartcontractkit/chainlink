// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";
import {IERC165} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {EnumerableSet} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";
import {ICapabilityConfiguration} from "./interfaces/ICapabilityConfiguration.sol";

/// @notice CapabilityRegistry is used to manage Nodes (including their links to Node
/// Operators), Capabilities, and DONs (Decentralized Oracle Networks) which are
/// sets of nodes that support those Capabilities.
contract CapabilityRegistry is OwnerIsCreator, TypeAndVersionInterface {
  // Add the library methods
  using EnumerableSet for EnumerableSet.Bytes32Set;

  struct NodeOperator {
    /// @notice The address of the admin that can manage a node
    /// operator
    address admin;
    /// @notice Human readable name of a Node Operator managing the node
    string name;
  }

  struct NodeParams {
    /// @notice The id of the node operator that manages this node
    uint32 nodeOperatorId;
    /// @notice The signer address for application-layer message verification.
    address signer;
    /// @notice This is an Ed25519 public key that is used to identify a node.
    /// This key is guaranteed to be unique in the CapabilityRegistry. It is
    /// used to identify a node in the the P2P network.
    bytes32 p2pId;
    /// @notice The list of hashed capability IDs supported by the node
    bytes32[] hashedCapabilityIds;
  }

  struct Node {
    /// @notice The node's parameters
    /// @notice The id of the node operator that manages this node
    uint32 nodeOperatorId;
    /// @notice The number of times the node's capability has been updated
    uint32 configCount;
    /// @notice The signer address for application-layer message verification.
    address signer;
    /// @notice This is an Ed25519 public key that is used to identify a node.
    /// This key is guaranteed to be unique in the CapabilityRegistry. It is
    /// used to identify a node in the the P2P network.
    bytes32 p2pId;
    /// @notice The node's supported capabilities
    /// @dev This is stored as a map so that we can easily update to a set of
    /// new capabilities by incrementing the configCount and creating a
    /// new set of supported capability IDs
    mapping(uint32 configCount => EnumerableSet.Bytes32Set capabilityId) supportedCapabilityIds;
  }

  // CapabilityResponseType indicates whether remote response requires
  // aggregation or is an already aggregated report. There are multiple
  // possible ways to aggregate.
  enum CapabilityResponseType {
    // No additional aggregation is needed on the remote response.
    REPORT,
    // A number of identical observations need to be aggregated.
    OBSERVATION_IDENTICAL
  }

  struct Capability {
    // The `labelledName` is a partially qualified ID for the capability.
    //
    // Given the following capability ID: {name}:{label1_key}_{label1_value}:{label2_key}_{label2_value}@{version}
    // Then we denote the `labelledName` as the `{name}:{label1_key}_{label1_value}:{label2_key}_{label2_value}` portion of the ID.
    //
    // Ex. id = "data-streams-reports:chain:ethereum@1.0.0"
    //     labelledName = "data-streams-reports:chain:ethereum"
    //
    // bytes32(string); validation regex: ^[a-z0-9_\-:]{1,32}$
    bytes32 labelledName;
    // Semver, e.g., "1.2.3"
    // bytes32(string); must be valid Semver + max 32 characters.
    bytes32 version;
    // responseType indicates whether remote response requires
    // aggregation or is an OCR report. There are multiple possible
    // ways to aggregate.
    CapabilityResponseType responseType;
    // An address to the capability configuration contract. Having this defined
    // on a capability enforces consistent configuration across DON instances
    // serving the same capability. Configuration contract MUST implement
    // CapabilityConfigurationContractInterface.
    //
    // The main use cases are:
    // 1) Sharing capability configuration across DON instances
    // 2) Inspect and modify on-chain configuration without off-chain
    // capability code.
    //
    // It is not recommended to store configuration which requires knowledge of
    // the DON membership.
    address configurationContract;
  }

  /// @notice CapabilityConfiguration is a struct that holds the capability configuration
  /// for a specific DON
  struct CapabilityConfiguration {
    /// @notice The capability Id
    bytes32 capabilityId;
    /// @notice The capability config specific to a DON.  This will be decoded
    /// offchain
    bytes config;
  }

  /// @notice DON (Decentralized Oracle Network) is a grouping of nodes that support
  // the same capabilities.
  struct DON {
    /// @notice Computed. Auto-increment.
    uint32 id;
    /// @notice True if the DON is public.  A public DON means that it accepts
    /// external capability requests
    bool isPublic;
    /// @notice The set of p2pIds of nodes that belong to this DON. A node (the same
    // p2pId) can belong to multiple DONs.
    EnumerableSet.Bytes32Set nodes;
    /// @notice The set of capabilityIds
    EnumerableSet.Bytes32Set capabilityIds;
    /// @notice Mapping from hashed capability IDs to configs
    mapping(bytes32 capabilityId => bytes config) capabilityConfigs;
  }

  /// @notice This error is thrown when a caller is not allowed
  /// to execute the transaction
  error AccessForbidden();

  /// @notice This error is thrown when there is a mismatch between
  /// array arguments
  /// @param lengthOne The length of the first array argument
  /// @param lengthTwo The length of the second array argument
  error LengthMismatch(uint256 lengthOne, uint256 lengthTwo);

  /// @notice This error is thrown when trying to set a node operator's
  /// admin address to the zero address
  error InvalidNodeOperatorAdmin();

  /// @notice This error is thrown when trying to add a node with P2P ID that
  /// is empty bytes or a duplicate.
  /// @param p2pId The provided P2P ID
  error InvalidNodeP2PId(bytes32 p2pId);

  /// @notice This error is thrown when trying to add a node without
  /// capabilities or with capabilities that do not exist.
  /// @param hashedCapabilityIds The IDs of the capabilities that are being added.
  error InvalidNodeCapabilities(bytes32[] hashedCapabilityIds);

  /// @notice This event is emitted when a new node is added
  /// @param p2pId The P2P ID of the node
  /// @param nodeOperatorId The ID of the node operator that manages this node
  event NodeAdded(bytes32 p2pId, uint256 nodeOperatorId);

  /// @notice This event is emitted when a node is removed
  /// @param p2pId The P2P ID of the node that was removed
  event NodeRemoved(bytes32 p2pId);

  /// @notice This event is emitted when a node is updated
  /// @param p2pId The P2P ID of the node
  /// @param nodeOperatorId The ID of the node operator that manages this node
  /// @param signer The node's signer address
  event NodeUpdated(bytes32 p2pId, uint256 nodeOperatorId, address signer);

  /// @notice This event is emitted when a new DON is created
  /// @param donId The ID of the newly created DON
  /// @param isPublic True if the newly created DON is public
  event DONAdded(uint256 donId, bool isPublic);

  /// @notice This error is thrown when trying to set the node's
  /// signer address to zero
  error InvalidNodeSigner();

  /// @notice This error is thrown when trying add a capability that already
  /// exists.
  error CapabilityAlreadyExists();

  /// @notice This error is thrown when trying to add a node to a DON where
  /// the node does not support the capability
  /// @param nodeP2PId The P2P ID of the node
  /// @param capabilityId The ID of the capability
  error NodeDoesNotSupportCapability(bytes32 nodeP2PId, bytes32 capabilityId);

  /// @notice This error is thrown when trying to add a capability configuration
  /// for a capability that was already configured on a DON
  /// @param donId The ID of the DON that the capability was configured for
  /// @param capabilityId The ID of the capability that was configured
  error DuplicateDONCapability(uint32 donId, bytes32 capabilityId);

  /// @notice This error is thrown when trying to add a duplicate node to a DON
  /// @param donId The ID of the DON that the node was added for
  /// @param nodeP2PId The P2P ID of the node
  error DuplicateDONNode(uint32 donId, bytes32 nodeP2PId);

  /// @notice This error is thrown when a capability with the provided hashed ID is
  /// not found.
  /// @param hashedCapabilityId The hashed ID used for the lookup.
  error CapabilityDoesNotExist(bytes32 hashedCapabilityId);

  /// @notice This error is thrown when trying to deprecate a capability that
  /// is deprecated.
  /// @param hashedCapabilityId The hashed ID of the capability that is deprecated.
  error CapabilityIsDeprecated(bytes32 hashedCapabilityId);

  /// @notice This error is thrown when trying to add a capability with a
  /// configuration contract that does not implement the required interface.
  /// @param proposedConfigurationContract The address of the proposed
  /// configuration contract.
  error InvalidCapabilityConfigurationContractInterface(address proposedConfigurationContract);

  /// @notice This event is emitted when a new node operator is added
  /// @param nodeOperatorId The ID of the newly added node operator
  /// @param admin The address of the admin that can manage the node
  /// operator
  /// @param name The human readable name of the node operator
  event NodeOperatorAdded(uint256 nodeOperatorId, address indexed admin, string name);

  /// @notice This event is emitted when a node operator is removed
  /// @param nodeOperatorId The ID of the node operator that was removed
  event NodeOperatorRemoved(uint256 nodeOperatorId);

  /// @notice This event is emitted when a node operator is updated
  /// @param nodeOperatorId The ID of the node operator that was updated
  /// @param admin The address of the node operator's admin
  /// @param name The node operator's human readable name
  event NodeOperatorUpdated(uint256 nodeOperatorId, address indexed admin, string name);

  /// @notice This event is emitted when a new capability is added
  /// @param hashedCapabilityId The hashed ID of the newly added capability
  event CapabilityAdded(bytes32 indexed hashedCapabilityId);

  /// @notice This event is emitted when a capability is deprecated
  /// @param hashedCapabilityId The hashed ID of the deprecated capability
  event CapabilityDeprecated(bytes32 indexed hashedCapabilityId);

  mapping(bytes32 => Capability) private s_capabilities;
  /// @notice Set of hashed capability IDs.
  /// A hashed ID is created by the function `getHashedCapabilityId`.
  EnumerableSet.Bytes32Set private s_hashedCapabilityIds;
  /// @notice Set of deprecated hashed capability IDs,
  /// A hashed ID is created by the function `getHashedCapabilityId`.
  ///
  /// Deprecated capabilities are skipped by the `getCapabilities` function.
  EnumerableSet.Bytes32Set private s_deprecatedHashedCapabilityIds;

  /// @notice Mapping of node operators
  mapping(uint256 nodeOperatorId => NodeOperator nodeOperator) private s_nodeOperators;

  /// @notice Mapping of nodes
  mapping(bytes32 p2pId => Node node) private s_nodes;

  /// @notice Mapping of DON IDs to DONs
  mapping(uint32 donId => DON don) private s_dons;

  /// @notice The latest node operator ID
  /// @dev No getter for this as this is an implementation detail
  uint256 private s_nodeOperatorId;

  /// @notice Starting with 1 to avoid confusion with the zero value.
  /// @dev No getter for this as this is an implementation detail
  uint32 private s_donId = 1;

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

  /// @notice Updates a node operator
  /// @param nodeOperatorIds The ID of the node operator being updated
  function updateNodeOperators(uint256[] calldata nodeOperatorIds, NodeOperator[] calldata nodeOperators) external {
    if (nodeOperatorIds.length != nodeOperators.length)
      revert LengthMismatch(nodeOperatorIds.length, nodeOperators.length);

    address owner = owner();
    for (uint256 i; i < nodeOperatorIds.length; ++i) {
      uint256 nodeOperatorId = nodeOperatorIds[i];
      NodeOperator memory nodeOperator = nodeOperators[i];
      if (nodeOperator.admin == address(0)) revert InvalidNodeOperatorAdmin();
      if (msg.sender != nodeOperator.admin && msg.sender != owner) revert AccessForbidden();

      if (
        s_nodeOperators[nodeOperatorId].admin != nodeOperator.admin ||
        keccak256(abi.encode(s_nodeOperators[nodeOperatorId].name)) != keccak256(abi.encode(nodeOperator.name))
      ) {
        s_nodeOperators[nodeOperatorId].admin = nodeOperator.admin;
        s_nodeOperators[nodeOperatorId].name = nodeOperator.name;
        emit NodeOperatorUpdated(nodeOperatorId, nodeOperator.admin, nodeOperator.name);
      }
    }
  }

  /// @notice Gets a node operator's data
  /// @param nodeOperatorId The ID of the node operator to query for
  /// @return NodeOperator The node operator data
  function getNodeOperator(uint256 nodeOperatorId) external view returns (NodeOperator memory) {
    return s_nodeOperators[nodeOperatorId];
  }

  /// @notice Adds nodes. Nodes can be added with deprecated capabilities to
  /// avoid breaking changes when deprecating capabilities.
  /// @param nodes The nodes to add
  function addNodes(NodeParams[] calldata nodes) external {
    for (uint256 i; i < nodes.length; ++i) {
      NodeParams memory node = nodes[i];

      bool isOwner = msg.sender == owner();

      NodeOperator memory nodeOperator = s_nodeOperators[node.nodeOperatorId];
      if (!isOwner && msg.sender != nodeOperator.admin) revert AccessForbidden();

      bool nodeExists = s_nodes[node.p2pId].signer != address(0);
      if (nodeExists || bytes32(node.p2pId) == bytes32("")) revert InvalidNodeP2PId(node.p2pId);

      if (node.signer == address(0)) revert InvalidNodeSigner();

      bytes32[] memory capabilityIds = node.hashedCapabilityIds;
      if (capabilityIds.length == 0) revert InvalidNodeCapabilities(capabilityIds);

      ++s_nodes[node.p2pId].configCount;

      uint32 capabilityConfigCount = s_nodes[node.p2pId].configCount;
      for (uint256 j; j < capabilityIds.length; ++j) {
        if (!s_hashedCapabilityIds.contains(capabilityIds[j])) revert InvalidNodeCapabilities(capabilityIds);
        s_nodes[node.p2pId].supportedCapabilityIds[capabilityConfigCount].add(capabilityIds[j]);
      }

      s_nodes[node.p2pId].nodeOperatorId = node.nodeOperatorId;
      s_nodes[node.p2pId].p2pId = node.p2pId;
      s_nodes[node.p2pId].signer = node.signer;
      emit NodeAdded(node.p2pId, node.nodeOperatorId);
    }
  }

  /// @notice Removes nodes.  The node operator admin or contract owner
  /// can remove nodes
  /// @param removedNodeP2PIds The P2P Ids of the nodes to remove
  function removeNodes(bytes32[] calldata removedNodeP2PIds) external {
    bool isOwner = msg.sender == owner();
    for (uint256 i; i < removedNodeP2PIds.length; ++i) {
      bytes32 p2pId = removedNodeP2PIds[i];

      bool nodeExists = s_nodes[p2pId].signer != address(0);
      if (!nodeExists) revert InvalidNodeP2PId(p2pId);

      NodeOperator memory nodeOperator = s_nodeOperators[s_nodes[p2pId].nodeOperatorId];

      if (!isOwner && msg.sender != nodeOperator.admin) revert AccessForbidden();
      delete s_nodes[p2pId];
      emit NodeRemoved(p2pId);
    }
  }

  /// @notice Updates nodes.  The node admin can update the node's signer address
  /// and reconfigure its supported capabilities
  /// @param nodes The nodes to update
  function updateNodes(NodeParams[] calldata nodes) external {
    for (uint256 i; i < nodes.length; ++i) {
      NodeParams memory node = nodes[i];

      bool isOwner = msg.sender == owner();

      NodeOperator memory nodeOperator = s_nodeOperators[node.nodeOperatorId];
      if (!isOwner && msg.sender != nodeOperator.admin) revert AccessForbidden();

      bool nodeExists = s_nodes[node.p2pId].signer != address(0);
      if (!nodeExists) revert InvalidNodeP2PId(node.p2pId);

      if (node.signer == address(0)) revert InvalidNodeSigner();

      bytes32[] memory supportedCapabilityIds = node.hashedCapabilityIds;
      if (supportedCapabilityIds.length == 0) revert InvalidNodeCapabilities(supportedCapabilityIds);

      s_nodes[node.p2pId].configCount++;
      uint32 capabilityConfigCount = s_nodes[node.p2pId].configCount;
      for (uint256 j; j < supportedCapabilityIds.length; ++j) {
        if (!s_hashedCapabilityIds.contains(supportedCapabilityIds[j]))
          revert InvalidNodeCapabilities(supportedCapabilityIds);
        s_nodes[node.p2pId].supportedCapabilityIds[capabilityConfigCount].add(supportedCapabilityIds[j]);
      }

      s_nodes[node.p2pId].nodeOperatorId = node.nodeOperatorId;
      s_nodes[node.p2pId].p2pId = node.p2pId;
      s_nodes[node.p2pId].signer = node.signer;
      emit NodeUpdated(node.p2pId, node.nodeOperatorId, node.signer);
    }
  }

  /// @notice Gets a node's data
  /// @param p2pId The P2P ID of the node to query for
  /// @return NodeParams The node data
  /// @return configCount The number of times the node has been configured
  function getNode(bytes32 p2pId) external view returns (NodeParams memory, uint32 configCount) {
    return (
      NodeParams({
        nodeOperatorId: s_nodes[p2pId].nodeOperatorId,
        p2pId: s_nodes[p2pId].p2pId,
        signer: s_nodes[p2pId].signer,
        hashedCapabilityIds: s_nodes[p2pId].supportedCapabilityIds[s_nodes[p2pId].configCount].values()
      }),
      s_nodes[p2pId].configCount
    );
  }

  function addCapability(Capability calldata capability) external onlyOwner {
    bytes32 hashedId = getHashedCapabilityId(capability.labelledName, capability.version);
    if (s_hashedCapabilityIds.contains(hashedId)) revert CapabilityAlreadyExists();

    if (capability.configurationContract != address(0)) {
      if (
        capability.configurationContract.code.length == 0 ||
        !IERC165(capability.configurationContract).supportsInterface(
          ICapabilityConfiguration.getCapabilityConfiguration.selector
        )
      ) revert InvalidCapabilityConfigurationContractInterface(capability.configurationContract);
    }

    s_hashedCapabilityIds.add(hashedId);
    s_capabilities[hashedId] = capability;

    emit CapabilityAdded(hashedId);
  }

  /// @notice Deprecates a capability by adding it to the deprecated list
  /// @param hashedCapabilityId The ID of the capability to deprecate
  function deprecateCapability(bytes32 hashedCapabilityId) external onlyOwner {
    if (!s_hashedCapabilityIds.contains(hashedCapabilityId)) revert CapabilityDoesNotExist(hashedCapabilityId);
    if (s_deprecatedHashedCapabilityIds.contains(hashedCapabilityId)) revert CapabilityIsDeprecated(hashedCapabilityId);

    s_deprecatedHashedCapabilityIds.add(hashedCapabilityId);
    emit CapabilityDeprecated(hashedCapabilityId);
  }

  /// @notice This function returns a Capability by its hashed ID. Use `getHashedCapabilityId` to get the hashed ID.
  function getCapability(bytes32 hashedId) public view returns (Capability memory) {
    return s_capabilities[hashedId];
  }

  /// @notice Returns all capabilities. This operation will copy capabilities
  /// to memory, which can be quite expensive. This is designed to mostly be
  /// used by view accessors that are queried without any gas fees.
  /// @return Capability[] An array of capabilities
  function getCapabilities() external view returns (Capability[] memory) {
    bytes32[] memory hashedCapabilityIds = s_hashedCapabilityIds.values();

    // Solidity does not support dynamic arrays in memory, so we create a
    // fixed-size array and copy the capabilities into it.
    Capability[] memory capabilities = new Capability[](
      hashedCapabilityIds.length - s_deprecatedHashedCapabilityIds.length()
    );

    // We need to keep track of the new index because we are skipping
    // deprecated capabilities.
    uint256 newIndex;

    for (uint256 i; i < hashedCapabilityIds.length; ++i) {
      bytes32 hashedCapabilityId = hashedCapabilityIds[i];

      if (!s_deprecatedHashedCapabilityIds.contains(hashedCapabilityId)) {
        capabilities[newIndex] = getCapability(hashedCapabilityId);
        newIndex++;
      }
    }

    return capabilities;
  }

  /// @notice This functions returns a capability id that has been hashed to fit into a bytes32 for cheaper access
  /// @return bytes32 A unique identifier for the capability
  function getHashedCapabilityId(bytes32 labelledName, bytes32 version) public pure returns (bytes32) {
    return keccak256(abi.encodePacked(labelledName, version));
  }

  /// @notice Returns whether a capability is deprecated
  /// @param hashedCapabilityId The hashed ID of the capability to check
  /// @return bool True if the capability is deprecated, false otherwise
  function isCapabilityDeprecated(bytes32 hashedCapabilityId) external view returns (bool) {
    return s_deprecatedHashedCapabilityIds.contains(hashedCapabilityId);
  }

  /// @notice Adds a DON made up by a group of nodes that support a list
  /// of capability configurations
  /// @param nodes The nodes making up the DON
  /// @param capabilityConfigurations The list of configurations for the
  /// capabilities supported by the DON
  /// @param isPublic True if the DON is public
  function addDON(
    bytes32[] calldata nodes,
    CapabilityConfiguration[] calldata capabilityConfigurations,
    bool isPublic
  ) external onlyOwner {
    uint32 id = s_donId;

    s_dons[id].id = id;
    s_dons[id].isPublic = isPublic;

    for (uint256 i; i < nodes.length; ++i) {
      bytes32 nodeP2PId = nodes[i];
      if (s_dons[id].nodes.contains(nodeP2PId)) revert DuplicateDONNode(id, nodeP2PId);

      s_dons[id].nodes.add(nodeP2PId);
    }

    for (uint256 i; i < capabilityConfigurations.length; ++i) {
      CapabilityConfiguration calldata configuration = capabilityConfigurations[i];
      bytes32 capabilityId = configuration.capabilityId;

      if (!s_hashedCapabilityIds.contains(capabilityId)) revert CapabilityDoesNotExist(capabilityId);
      if (s_deprecatedHashedCapabilityIds.contains(capabilityId)) revert CapabilityIsDeprecated(capabilityId);

      if (s_dons[id].capabilityIds.contains(capabilityId)) revert DuplicateDONCapability(id, capabilityId);

      for (uint256 j; j < nodes.length; ++j) {
        bytes32 nodeP2PId = nodes[j];
        if (!s_nodes[nodeP2PId].supportedCapabilityIds[s_nodes[nodeP2PId].configCount].contains(capabilityId))
          revert NodeDoesNotSupportCapability(nodeP2PId, capabilityId);
      }

      s_dons[id].capabilityIds.add(capabilityId);
      s_dons[id].capabilityConfigs[capabilityId] = configuration.config;
    }

    ++s_donId;

    emit DONAdded(id, isPublic);
  }

  /// @notice Gets DON's data
  /// @param donId The DON ID
  /// @return uint32 The DON ID
  /// @return bool True if the DON is public
  /// @return bytes32[] The list of node P2P IDs that are in the DON
  /// @return CapabilityConfiguration[] The list of capability configurations supported by the DON
  function getDON(
    uint32 donId
  ) external view returns (uint32, bool, bytes32[] memory, CapabilityConfiguration[] memory) {
    bytes32[] memory capabilityIds = s_dons[donId].capabilityIds.values();
    CapabilityConfiguration[] memory capabilityConfigurations = new CapabilityConfiguration[](capabilityIds.length);

    for (uint256 i; i < capabilityConfigurations.length; ++i) {
      capabilityConfigurations[i] = CapabilityConfiguration({
        capabilityId: capabilityIds[i],
        config: s_dons[donId].capabilityConfigs[capabilityIds[i]]
      });
    }

    return (s_dons[donId].id, s_dons[donId].isPublic, s_dons[donId].nodes.values(), capabilityConfigurations);
  }

  /// @notice Returns the DON specific configuration for a capability
  /// @param donId The DON's ID
  /// @param capabilityId The Capability ID
  /// @return bytes The DON specific configuration for the capability
  function getDONCapabilityConfig(uint32 donId, bytes32 capabilityId) external view returns (bytes memory) {
    return s_dons[donId].capabilityConfigs[capabilityId];
  }
}
