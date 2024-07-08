// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";
import {IERC165} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {EnumerableSet} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";
import {ICapabilityConfiguration} from "./interfaces/ICapabilityConfiguration.sol";

/// @notice CapabilitiesRegistry is used to manage Nodes (including their links to Node
/// Operators), Capabilities, and DONs (Decentralized Oracle Networks) which are
/// sets of nodes that support those Capabilities.
/// @dev The contract currently stores the entire state of Node Operators, Nodes,
/// Capabilities and DONs in the contract and requires a full state migration
/// if an upgrade is ever required.  The team acknowledges this and is fine
/// reconfiguring the upgraded contract in the future so as to not add extra
/// complexity to this current version.
contract CapabilitiesRegistry is OwnerIsCreator, TypeAndVersionInterface {
  // Add the library methods
  using EnumerableSet for EnumerableSet.Bytes32Set;
  using EnumerableSet for EnumerableSet.UintSet;

  struct NodeOperator {
    /// @notice The address of the admin that can manage a node
    /// operator
    address admin;
    /// @notice Human readable name of a Node Operator managing the node
    /// @dev The contract does not validate the length or characters of the
    /// node operator name because a trusted admin will supply these names.
    /// We reduce gas costs by omitting these checks on-chain.
    string name;
  }

  struct NodeParams {
    /// @notice The id of the node operator that manages this node
    uint32 nodeOperatorId;
    /// @notice The signer address for application-layer message verification.
    bytes32 signer;
    /// @notice This is an Ed25519 public key that is used to identify a node.
    /// This key is guaranteed to be unique in the CapabilitiesRegistry. It is
    /// used to identify a node in the the P2P network.
    bytes32 p2pId;
    /// @notice The list of hashed capability IDs supported by the node
    bytes32[] hashedCapabilityIds;
  }

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
    /// @notice The list of hashed capability IDs supported by the node
    bytes32[] hashedCapabilityIds;
    /// @notice The list of capabilities DON Ids supported by the node. A node
    /// can belong to multiple capabilities DONs. This list does not include a
    /// Workflow DON id if the node belongs to one.
    uint256[] capabilitiesDONIds;
  }

  struct Node {
    /// @notice The node's parameters
    /// @notice The id of the node operator that manages this node
    uint32 nodeOperatorId;
    /// @notice The number of times the node's configuration has been updated
    uint32 configCount;
    /// @notice The ID of the Workflow DON that the node belongs to. A node can
    /// only belong to one DON that accepts Workflows.
    uint32 workflowDONId;
    /// @notice The signer address for application-layer message verification.
    /// @dev This key is guaranteed to be unique in the CapabilitiesRegistry
    /// as a signer address can only belong to one node.
    /// @dev This should be the ABI encoded version of the node's address.
    /// I.e 0x0000address.  The Capability Registry does not store it as an address so that
    /// non EVM chains with addresses greater than 20 bytes can be supported
    /// in the future.
    bytes32 signer;
    /// @notice This is an Ed25519 public key that is used to identify a node.
    /// This key is guaranteed to be unique in the CapabilitiesRegistry. It is
    /// used to identify a node in the the P2P network.
    bytes32 p2pId;
    /// @notice The node's supported capabilities
    /// @dev This is stored as a map so that we can easily update to a set of
    /// new capabilities by incrementing the configCount and creating a
    /// new set of supported capability IDs
    mapping(uint32 configCount => EnumerableSet.Bytes32Set capabilityId) supportedHashedCapabilityIds;
    /// @notice The list of capabilities DON Ids supported by the node. A node
    /// can belong to multiple capabilities DONs. This list does not include a
    /// Workflow DON id if the node belongs to one.
    EnumerableSet.UintSet capabilitiesDONIds;
  }

  /// @notice CapabilityResponseType indicates whether remote response requires
  // aggregation or is an already aggregated report. There are multiple
  // possible ways to aggregate.
  /// @dev REPORT response type receives signatures together with the response that
  /// is used to verify the data.  OBSERVATION_IDENTICAL just receives data without
  /// signatures and waits for some number of observations before proceeeding to
  /// the next step
  enum CapabilityResponseType {
    // No additional aggregation is needed on the remote response.
    REPORT,
    // A number of identical observations need to be aggregated.
    OBSERVATION_IDENTICAL
  }

  /// @notice CapabilityType indicates the type of capability which determines
  /// where the capability can be used in a Workflow Spec.
  enum CapabilityType {
    TRIGGER,
    ACTION,
    CONSENSUS,
    TARGET
  }

  struct Capability {
    /// @notice The partially qualified ID for the capability.
    /// @dev Given the following capability ID: {name}:{label1_key}_{label1_value}:{label2_key}_{label2_value}@{version}
    // Then we denote the `labelledName` as the `{name}:{label1_key}_{label1_value}:{label2_key}_{label2_value}` portion of the ID.
    ///
    /// Ex. id = "data-streams-reports:chain:ethereum@1.0.0"
    ///     labelledName = "data-streams-reports:chain:ethereum"
    string labelledName;
    /// @notice Semver, e.g., "1.2.3"
    /// @dev must be valid Semver + max 32 characters.
    string version;
    /// @notice CapabilityType indicates the type of capability which determines
    /// where the capability can be used in a Workflow Spec.
    CapabilityType capabilityType;
    /// @notice CapabilityResponseType indicates whether remote response requires
    // aggregation or is an already aggregated report. There are multiple
    // possible ways to aggregate.
    CapabilityResponseType responseType;
    /// @notice An address to the capability configuration contract. Having this defined
    // on a capability enforces consistent configuration across DON instances
    // serving the same capability. Configuration contract MUST implement
    // CapabilityConfigurationContractInterface.
    //
    /// @dev The main use cases are:
    // 1) Sharing capability configuration across DON instances
    // 2) Inspect and modify on-chain configuration without off-chain
    // capability code.
    //
    // It is not recommended to store configuration which requires knowledge of
    // the DON membership.
    address configurationContract;
  }

  struct CapabilityInfo {
    /// @notice A hashed ID created by the `getHashedCapabilityId` function.
    bytes32 hashedId;
    /// @notice The partially qualified ID for the capability.
    /// @dev Given the following capability ID: {name}:{label1_key}_{label1_value}:{label2_key}_{label2_value}@{version}
    // Then we denote the `labelledName` as the `{name}:{label1_key}_{label1_value}:{label2_key}_{label2_value}` portion of the ID.
    ///
    /// Ex. id = "data-streams-reports:chain:ethereum@1.0.0"
    ///     labelledName = "data-streams-reports:chain:ethereum"
    string labelledName;
    /// @notice Semver, e.g., "1.2.3"
    /// @dev must be valid Semver + max 32 characters.
    string version;
    /// @notice CapabilityType indicates the type of capability which determines
    /// where the capability can be used in a Workflow Spec.
    CapabilityType capabilityType;
    /// @notice CapabilityResponseType indicates whether remote response requires
    // aggregation or is an already aggregated report. There are multiple
    // possible ways to aggregate.
    CapabilityResponseType responseType;
    /// @notice An address to the capability configuration contract. Having this defined
    // on a capability enforces consistent configuration across DON instances
    // serving the same capability. Configuration contract MUST implement
    // CapabilityConfigurationContractInterface.
    //
    /// @dev The main use cases are:
    // 1) Sharing capability configuration across DON instances
    // 2) Inspect and modify on-chain configuration without off-chain
    // capability code.
    //
    // It is not recommended to store configuration which requires knowledge of
    // the DON membership.
    address configurationContract;
    /// @notice True if the capability is deprecated
    bool isDeprecated;
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

  struct DONCapabilityConfig {
    /// @notice The set of p2pIds of nodes that belong to this DON. A node (the same
    // p2pId) can belong to multiple DONs.
    EnumerableSet.Bytes32Set nodes;
    /// @notice The set of capabilityIds
    bytes32[] capabilityIds;
    /// @notice Mapping from hashed capability IDs to configs
    mapping(bytes32 capabilityId => bytes config) capabilityConfigs;
  }

  /// @notice DON (Decentralized Oracle Network) is a grouping of nodes that support
  // the same capabilities.
  struct DON {
    /// @notice Computed. Auto-increment.
    uint32 id;
    /// @notice The number of times the DON was configured
    uint32 configCount;
    /// @notice The f value for the DON.  This is the number of faulty nodes
    /// that the DON can tolerate. This can be different from the f value of
    /// the OCR instances that capabilities spawn.
    uint8 f;
    /// @notice True if the DON is public. A public DON means that it accepts
    /// external capability requests
    bool isPublic;
    /// @notice True if the DON accepts Workflows. A DON that accepts Workflows
    /// is called Workflow DON and it can process Workflow Specs. A Workflow
    /// DON also support one or more capabilities as well.
    bool acceptsWorkflows;
    /// @notice Mapping of config counts to configurations
    mapping(uint32 configCount => DONCapabilityConfig donConfig) config;
  }

  struct DONInfo {
    /// @notice Computed. Auto-increment.
    uint32 id;
    /// @notice The number of times the DON was configured
    uint32 configCount;
    /// @notice The f value for the DON.  This is the number of faulty nodes
    /// that the DON can tolerate. This can be different from the f value of
    /// the OCR instances that capabilities spawn.
    uint8 f;
    /// @notice True if the DON is public.  A public DON means that it accepts
    /// external capability requests
    bool isPublic;
    /// @notice True if the DON accepts Workflows.
    bool acceptsWorkflows;
    /// @notice List of member node P2P Ids
    bytes32[] nodeP2PIds;
    /// @notice List of capability configurations
    CapabilityConfiguration[] capabilityConfigurations;
  }

  /// @notice DONParams is a struct that holds the parameters for a DON.
  /// @dev This is needed to avoid "stack too deep" errors in _setDONConfig.
  struct DONParams {
    uint32 id;
    uint32 configCount;
    bool isPublic;
    bool acceptsWorkflows;
    uint8 f;
  }

  /// @notice This error is thrown when a caller is not allowed
  /// to execute the transaction
  /// @param sender The address that tried to execute the transaction
  error AccessForbidden(address sender);

  /// @notice This error is thrown when there is a mismatch between
  /// array arguments
  /// @param lengthOne The length of the first array argument
  /// @param lengthTwo The length of the second array argument
  error LengthMismatch(uint256 lengthOne, uint256 lengthTwo);

  /// @notice This error is thrown when trying to set a node operator's
  /// admin address to the zero address
  error InvalidNodeOperatorAdmin();

  /// @notice This error is thrown when trying to add a node with P2P ID that
  /// is empty bytes
  /// @param p2pId The provided P2P ID
  error InvalidNodeP2PId(bytes32 p2pId);

  /// @notice This error is thrown when trying to add a node without
  /// capabilities or with capabilities that do not exist.
  /// @param hashedCapabilityIds The IDs of the capabilities that are being added.
  error InvalidNodeCapabilities(bytes32[] hashedCapabilityIds);

  /// @notice This error is emitted when a DON does not exist
  /// @param donId The ID of the nonexistent DON
  error DONDoesNotExist(uint32 donId);

  /// @notice This error is thrown when trying to set the node's
  /// signer address to zero or if the signer address has already
  /// been used by another node
  error InvalidNodeSigner();

  /// @notice This error is thrown when trying to add a capability that already
  /// exists.
  /// @param hashedCapabilityId The hashed capability ID of the capability
  /// that already exists
  error CapabilityAlreadyExists(bytes32 hashedCapabilityId);

  /// @notice This error is thrown when trying to add a node that already
  /// exists.
  /// @param nodeP2PId The P2P ID of the node that already exists
  error NodeAlreadyExists(bytes32 nodeP2PId);

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

  /// @notice This error is thrown when trying to configure a DON with invalid
  /// fault tolerance value.
  /// @param f The proposed fault tolerance value
  /// @param nodeCount The proposed number of nodes in the DON
  error InvalidFaultTolerance(uint8 f, uint256 nodeCount);

  /// @notice This error is thrown when a capability with the provided hashed ID is
  /// not found.
  /// @param hashedCapabilityId The hashed ID used for the lookup.
  error CapabilityDoesNotExist(bytes32 hashedCapabilityId);

  /// @notice This error is thrown when trying to deprecate a capability that
  /// is deprecated.
  /// @param hashedCapabilityId The hashed ID of the capability that is deprecated.
  error CapabilityIsDeprecated(bytes32 hashedCapabilityId);

  /// @notice This error is thrown when a node with the provided P2P ID is
  /// not found.
  /// @param nodeP2PId The node P2P ID used for the lookup.
  error NodeDoesNotExist(bytes32 nodeP2PId);

  /// @notice This error is thrown when a node operator does not exist
  /// @param nodeOperatorId The ID of the node operator that does not exist
  error NodeOperatorDoesNotExist(uint32 nodeOperatorId);

  /// @notice This error is thrown when trying to remove a node that is still
  /// part of a capabitlities DON
  /// @param donId The Id of the DON the node belongs to
  /// @param nodeP2PId The P2P Id of the node being removed
  error NodePartOfCapabilitiesDON(uint32 donId, bytes32 nodeP2PId);

  /// @notice This error is thrown when attempting to add a node to a second
  /// Workflow DON or when trying to remove a node that belongs to a Workflow
  /// DON
  /// @param donId The Id of the DON the node belongs to
  /// @param nodeP2PId The P2P Id of the node
  error NodePartOfWorkflowDON(uint32 donId, bytes32 nodeP2PId);

  /// @notice This error is thrown when removing a capability from the node
  /// when that capability is still required by one of the DONs the node
  /// belongs to.
  /// @param hashedCapabilityId The hashed ID of the capability
  /// @param donId The ID of the DON that requires the capability
  error CapabilityRequiredByDON(bytes32 hashedCapabilityId, uint32 donId);

  /// @notice This error is thrown when trying to add a capability with a
  /// configuration contract that does not implement the required interface.
  /// @param proposedConfigurationContract The address of the proposed
  /// configuration contract.
  error InvalidCapabilityConfigurationContractInterface(address proposedConfigurationContract);

  /// @notice This event is emitted when a new node is added
  /// @param p2pId The P2P ID of the node
  /// @param nodeOperatorId The ID of the node operator that manages this node
  /// @param signer The encoded node's signer address
  event NodeAdded(bytes32 p2pId, uint32 indexed nodeOperatorId, bytes32 signer);

  /// @notice This event is emitted when a node is removed
  /// @param p2pId The P2P ID of the node that was removed
  event NodeRemoved(bytes32 p2pId);

  /// @notice This event is emitted when a node is updated
  /// @param p2pId The P2P ID of the node
  /// @param nodeOperatorId The ID of the node operator that manages this node
  /// @param signer The node's signer address
  event NodeUpdated(bytes32 p2pId, uint32 indexed nodeOperatorId, bytes32 signer);

  /// @notice This event is emitted when a DON's config is set
  /// @param donId The ID of the DON the config was set for
  /// @param configCount The number of times the DON has been
  /// configured
  event ConfigSet(uint32 donId, uint32 configCount);

  /// @notice This event is emitted when a new node operator is added
  /// @param nodeOperatorId The ID of the newly added node operator
  /// @param admin The address of the admin that can manage the node
  /// operator
  /// @param name The human readable name of the node operator
  event NodeOperatorAdded(uint32 indexed nodeOperatorId, address indexed admin, string name);

  /// @notice This event is emitted when a node operator is removed
  /// @param nodeOperatorId The ID of the node operator that was removed
  event NodeOperatorRemoved(uint32 indexed nodeOperatorId);

  /// @notice This event is emitted when a node operator is updated
  /// @param nodeOperatorId The ID of the node operator that was updated
  /// @param admin The address of the node operator's admin
  /// @param name The node operator's human readable name
  event NodeOperatorUpdated(uint32 indexed nodeOperatorId, address indexed admin, string name);

  /// @notice This event is emitted when a new capability is added
  /// @param hashedCapabilityId The hashed ID of the newly added capability
  event CapabilityConfigured(bytes32 indexed hashedCapabilityId);

  /// @notice This event is emitted when a capability is deprecated
  /// @param hashedCapabilityId The hashed ID of the deprecated capability
  event CapabilityDeprecated(bytes32 indexed hashedCapabilityId);

  /// @notice Mapping of capabilities
  mapping(bytes32 hashedCapabilityId => Capability capability) private s_capabilities;

  /// @notice Set of hashed capability IDs.
  /// A hashed ID is created by the function `getHashedCapabilityId`.
  EnumerableSet.Bytes32Set private s_hashedCapabilityIds;

  /// @notice Set of deprecated hashed capability IDs,
  /// A hashed ID is created by the function `getHashedCapabilityId`.
  EnumerableSet.Bytes32Set private s_deprecatedHashedCapabilityIds;

  /// @notice Encoded node signer addresses
  EnumerableSet.Bytes32Set private s_nodeSigners;

  /// @notice Set of node P2P IDs
  EnumerableSet.Bytes32Set private s_nodeP2PIds;

  /// @notice Mapping of node operators
  mapping(uint32 nodeOperatorId => NodeOperator nodeOperator) private s_nodeOperators;

  /// @notice Mapping of nodes
  mapping(bytes32 p2pId => Node node) private s_nodes;

  /// @notice Mapping of DON IDs to DONs
  mapping(uint32 donId => DON don) private s_dons;

  /// @notice The next ID to assign a new node operator to
  /// @dev Starting with 1 to avoid confusion with the zero value
  /// @dev No getter for this as this is an implementation detail
  uint32 private s_nextNodeOperatorId = 1;

  /// @notice The next ID to assign a new DON to
  /// @dev Starting with 1 to avoid confusion with the zero value
  /// @dev No getter for this as this is an implementation detail
  uint32 private s_nextDONId = 1;

  function typeAndVersion() external pure override returns (string memory) {
    return "CapabilitiesRegistry 1.0.0";
  }

  /// @notice Adds a list of node operators
  /// @param nodeOperators List of node operators to add
  function addNodeOperators(NodeOperator[] calldata nodeOperators) external onlyOwner {
    for (uint256 i; i < nodeOperators.length; ++i) {
      NodeOperator memory nodeOperator = nodeOperators[i];
      if (nodeOperator.admin == address(0)) revert InvalidNodeOperatorAdmin();
      uint32 nodeOperatorId = s_nextNodeOperatorId;
      s_nodeOperators[nodeOperatorId] = NodeOperator({admin: nodeOperator.admin, name: nodeOperator.name});
      ++s_nextNodeOperatorId;
      emit NodeOperatorAdded(nodeOperatorId, nodeOperator.admin, nodeOperator.name);
    }
  }

  /// @notice Removes a node operator
  /// @param nodeOperatorIds The IDs of the node operators to remove
  function removeNodeOperators(uint32[] calldata nodeOperatorIds) external onlyOwner {
    for (uint32 i; i < nodeOperatorIds.length; ++i) {
      uint32 nodeOperatorId = nodeOperatorIds[i];
      delete s_nodeOperators[nodeOperatorId];
      emit NodeOperatorRemoved(nodeOperatorId);
    }
  }

  /// @notice Updates a node operator
  /// @param nodeOperatorIds The ID of the node operator being updated
  /// @param nodeOperators The updated node operator params
  function updateNodeOperators(uint32[] calldata nodeOperatorIds, NodeOperator[] calldata nodeOperators) external {
    if (nodeOperatorIds.length != nodeOperators.length)
      revert LengthMismatch(nodeOperatorIds.length, nodeOperators.length);

    address owner = owner();
    for (uint256 i; i < nodeOperatorIds.length; ++i) {
      uint32 nodeOperatorId = nodeOperatorIds[i];

      NodeOperator storage currentNodeOperator = s_nodeOperators[nodeOperatorId];
      if (currentNodeOperator.admin == address(0)) revert NodeOperatorDoesNotExist(nodeOperatorId);

      NodeOperator memory nodeOperator = nodeOperators[i];
      if (nodeOperator.admin == address(0)) revert InvalidNodeOperatorAdmin();
      if (msg.sender != nodeOperator.admin && msg.sender != owner) revert AccessForbidden(msg.sender);

      if (
        currentNodeOperator.admin != nodeOperator.admin ||
        keccak256(abi.encode(currentNodeOperator.name)) != keccak256(abi.encode(nodeOperator.name))
      ) {
        currentNodeOperator.admin = nodeOperator.admin;
        currentNodeOperator.name = nodeOperator.name;
        emit NodeOperatorUpdated(nodeOperatorId, nodeOperator.admin, nodeOperator.name);
      }
    }
  }

  /// @notice Gets a node operator's data
  /// @param nodeOperatorId The ID of the node operator to query for
  /// @return NodeOperator The node operator data
  function getNodeOperator(uint32 nodeOperatorId) external view returns (NodeOperator memory) {
    return s_nodeOperators[nodeOperatorId];
  }

  /// @notice Gets all node operators
  /// @return NodeOperator[] All node operators
  function getNodeOperators() external view returns (NodeOperator[] memory) {
    uint32 nodeOperatorId = s_nextNodeOperatorId;
    /// Minus one to account for s_nextNodeOperatorId starting at index 1
    NodeOperator[] memory nodeOperators = new NodeOperator[](s_nextNodeOperatorId - 1);
    uint256 idx;
    for (uint32 i = 1; i < nodeOperatorId; ++i) {
      if (s_nodeOperators[i].admin != address(0)) {
        nodeOperators[idx] = s_nodeOperators[i];
        ++idx;
      }
    }
    if (idx != s_nextNodeOperatorId - 1) {
      assembly {
        mstore(nodeOperators, idx)
      }
    }
    return nodeOperators;
  }

  /// @notice Adds nodes. Nodes can be added with deprecated capabilities to
  /// avoid breaking changes when deprecating capabilities.
  /// @param nodes The nodes to add
  function addNodes(NodeParams[] calldata nodes) external {
    bool isOwner = msg.sender == owner();
    for (uint256 i; i < nodes.length; ++i) {
      NodeParams memory node = nodes[i];

      NodeOperator memory nodeOperator = s_nodeOperators[node.nodeOperatorId];
      if (nodeOperator.admin == address(0)) revert NodeOperatorDoesNotExist(node.nodeOperatorId);
      if (!isOwner && msg.sender != nodeOperator.admin) revert AccessForbidden(msg.sender);

      Node storage storedNode = s_nodes[node.p2pId];
      if (storedNode.signer != bytes32("")) revert NodeAlreadyExists(node.p2pId);
      if (node.p2pId == bytes32("")) revert InvalidNodeP2PId(node.p2pId);

      if (node.signer == bytes32("") || s_nodeSigners.contains(node.signer)) revert InvalidNodeSigner();

      bytes32[] memory capabilityIds = node.hashedCapabilityIds;
      if (capabilityIds.length == 0) revert InvalidNodeCapabilities(capabilityIds);

      ++storedNode.configCount;

      uint32 capabilityConfigCount = storedNode.configCount;
      for (uint256 j; j < capabilityIds.length; ++j) {
        if (!s_hashedCapabilityIds.contains(capabilityIds[j])) revert InvalidNodeCapabilities(capabilityIds);
        storedNode.supportedHashedCapabilityIds[capabilityConfigCount].add(capabilityIds[j]);
      }

      storedNode.nodeOperatorId = node.nodeOperatorId;
      storedNode.p2pId = node.p2pId;
      storedNode.signer = node.signer;
      s_nodeSigners.add(node.signer);
      s_nodeP2PIds.add(node.p2pId);
      emit NodeAdded(node.p2pId, node.nodeOperatorId, node.signer);
    }
  }

  /// @notice Removes nodes.  The node operator admin or contract owner
  /// can remove nodes
  /// @param removedNodeP2PIds The P2P Ids of the nodes to remove
  function removeNodes(bytes32[] calldata removedNodeP2PIds) external {
    bool isOwner = msg.sender == owner();
    for (uint256 i; i < removedNodeP2PIds.length; ++i) {
      bytes32 p2pId = removedNodeP2PIds[i];

      Node storage node = s_nodes[p2pId];

      if (node.signer == bytes32("")) revert NodeDoesNotExist(p2pId);
      if (node.capabilitiesDONIds.length() > 0)
        revert NodePartOfCapabilitiesDON(uint32(node.capabilitiesDONIds.at(i)), p2pId);
      if (node.workflowDONId != 0) revert NodePartOfWorkflowDON(node.workflowDONId, p2pId);

      if (!isOwner && msg.sender != s_nodeOperators[node.nodeOperatorId].admin) revert AccessForbidden(msg.sender);
      s_nodeSigners.remove(node.signer);
      s_nodeP2PIds.remove(node.p2pId);
      delete s_nodes[p2pId];
      emit NodeRemoved(p2pId);
    }
  }

  /// @notice Updates nodes.  The node admin can update the node's signer address
  /// and reconfigure its supported capabilities
  /// @param nodes The nodes to update
  function updateNodes(NodeParams[] calldata nodes) external {
    bool isOwner = msg.sender == owner();
    for (uint256 i; i < nodes.length; ++i) {
      NodeParams memory node = nodes[i];
      NodeOperator memory nodeOperator = s_nodeOperators[node.nodeOperatorId];
      if (!isOwner && msg.sender != nodeOperator.admin) revert AccessForbidden(msg.sender);

      Node storage storedNode = s_nodes[node.p2pId];
      if (storedNode.signer == bytes32("")) revert NodeDoesNotExist(node.p2pId);
      if (node.signer == bytes32("")) revert InvalidNodeSigner();

      bytes32 previousSigner = storedNode.signer;
      if (previousSigner != node.signer) {
        if (s_nodeSigners.contains(node.signer)) revert InvalidNodeSigner();
        storedNode.signer = node.signer;
        s_nodeSigners.remove(previousSigner);
        s_nodeSigners.add(node.signer);
      }

      bytes32[] memory supportedHashedCapabilityIds = node.hashedCapabilityIds;
      if (supportedHashedCapabilityIds.length == 0) revert InvalidNodeCapabilities(supportedHashedCapabilityIds);

      uint32 capabilityConfigCount = ++storedNode.configCount;
      for (uint256 j; j < supportedHashedCapabilityIds.length; ++j) {
        if (!s_hashedCapabilityIds.contains(supportedHashedCapabilityIds[j]))
          revert InvalidNodeCapabilities(supportedHashedCapabilityIds);
        storedNode.supportedHashedCapabilityIds[capabilityConfigCount].add(supportedHashedCapabilityIds[j]);
      }

      // Validate that capabilities required by a Workflow DON are still supported
      uint32 nodeWorkflowDONId = storedNode.workflowDONId;
      if (nodeWorkflowDONId != 0) {
        bytes32[] memory workflowDonCapabilityIds = s_dons[nodeWorkflowDONId]
          .config[s_dons[nodeWorkflowDONId].configCount]
          .capabilityIds;

        for (uint256 j; j < workflowDonCapabilityIds.length; ++j) {
          if (!storedNode.supportedHashedCapabilityIds[capabilityConfigCount].contains(workflowDonCapabilityIds[j]))
            revert CapabilityRequiredByDON(workflowDonCapabilityIds[j], nodeWorkflowDONId);
        }
      }

      // Validate that capabilities required by capabilities DONs are still supported
      uint256[] memory capabilitiesDONIds = storedNode.capabilitiesDONIds.values();
      for (uint32 j; j < capabilitiesDONIds.length; ++j) {
        uint32 donId = uint32(capabilitiesDONIds[j]);
        bytes32[] memory donCapabilityIds = s_dons[donId].config[s_dons[donId].configCount].capabilityIds;

        for (uint256 k; k < donCapabilityIds.length; ++k) {
          if (!storedNode.supportedHashedCapabilityIds[capabilityConfigCount].contains(donCapabilityIds[k]))
            revert CapabilityRequiredByDON(donCapabilityIds[k], donId);
        }
      }

      storedNode.nodeOperatorId = node.nodeOperatorId;
      storedNode.p2pId = node.p2pId;

      emit NodeUpdated(node.p2pId, node.nodeOperatorId, node.signer);
    }
  }

  /// @notice Gets a node's data
  /// @param p2pId The P2P ID of the node to query for
  /// @return nodeInfo NodeInfo The node data
  function getNode(bytes32 p2pId) public view returns (NodeInfo memory nodeInfo) {
    return (
      NodeInfo({
        nodeOperatorId: s_nodes[p2pId].nodeOperatorId,
        p2pId: s_nodes[p2pId].p2pId,
        signer: s_nodes[p2pId].signer,
        hashedCapabilityIds: s_nodes[p2pId].supportedHashedCapabilityIds[s_nodes[p2pId].configCount].values(),
        configCount: s_nodes[p2pId].configCount,
        workflowDONId: s_nodes[p2pId].workflowDONId,
        capabilitiesDONIds: s_nodes[p2pId].capabilitiesDONIds.values()
      })
    );
  }

  /// @notice Gets all nodes
  /// @return NodeInfo[] All nodes in the capability registry
  function getNodes() external view returns (NodeInfo[] memory) {
    bytes32[] memory p2pIds = s_nodeP2PIds.values();
    NodeInfo[] memory nodesInfo = new NodeInfo[](p2pIds.length);

    for (uint256 i; i < p2pIds.length; ++i) {
      nodesInfo[i] = getNode(p2pIds[i]);
    }
    return nodesInfo;
  }

  /// @notice Adds a new capability to the capability registry
  /// @param capabilities The capabilities being added
  /// @dev There is no function to update capabilities as this would require
  /// nodes to trust that the capabilities they support are not updated by the
  /// admin
  function addCapabilities(Capability[] calldata capabilities) external onlyOwner {
    for (uint256 i; i < capabilities.length; ++i) {
      Capability memory capability = capabilities[i];
      bytes32 hashedCapabilityId = getHashedCapabilityId(capability.labelledName, capability.version);
      if (!s_hashedCapabilityIds.add(hashedCapabilityId)) revert CapabilityAlreadyExists(hashedCapabilityId);
      _setCapability(hashedCapabilityId, capability);
    }
  }

  /// @notice Deprecates a capability
  /// @param hashedCapabilityIds[] The IDs of the capabilities to deprecate
  function deprecateCapabilities(bytes32[] calldata hashedCapabilityIds) external onlyOwner {
    for (uint256 i; i < hashedCapabilityIds.length; ++i) {
      bytes32 hashedCapabilityId = hashedCapabilityIds[i];
      if (!s_hashedCapabilityIds.contains(hashedCapabilityId)) revert CapabilityDoesNotExist(hashedCapabilityId);
      if (!s_deprecatedHashedCapabilityIds.add(hashedCapabilityId)) revert CapabilityIsDeprecated(hashedCapabilityId);

      emit CapabilityDeprecated(hashedCapabilityId);
    }
  }

  /// @notice Returns a Capability by its hashed ID.
  /// @dev Use `getHashedCapabilityId` to get the hashed ID.
  function getCapability(bytes32 hashedId) public view returns (CapabilityInfo memory) {
    return (
      CapabilityInfo({
        hashedId: hashedId,
        labelledName: s_capabilities[hashedId].labelledName,
        version: s_capabilities[hashedId].version,
        capabilityType: s_capabilities[hashedId].capabilityType,
        responseType: s_capabilities[hashedId].responseType,
        configurationContract: s_capabilities[hashedId].configurationContract,
        isDeprecated: s_deprecatedHashedCapabilityIds.contains(hashedId)
      })
    );
  }

  /// @notice Returns all capabilities. This operation will copy capabilities
  /// to memory, which can be quite expensive. This is designed to mostly be
  /// used by view accessors that are queried without any gas fees.
  /// @return CapabilityInfo[] List of capabilities
  function getCapabilities() external view returns (CapabilityInfo[] memory) {
    bytes32[] memory hashedCapabilityIds = s_hashedCapabilityIds.values();
    CapabilityInfo[] memory capabilitiesInfo = new CapabilityInfo[](hashedCapabilityIds.length);

    for (uint256 i; i < hashedCapabilityIds.length; ++i) {
      capabilitiesInfo[i] = getCapability(hashedCapabilityIds[i]);
    }
    return capabilitiesInfo;
  }

  /// @notice This functions returns a capability id that has been hashed to fit into a bytes32 for cheaper access
  /// @param labelledName The name of the capability
  /// @param version The capability's version number
  /// @return bytes32 A unique identifier for the capability
  /// @dev The hash of the encoded labelledName and version
  function getHashedCapabilityId(string memory labelledName, string memory version) public pure returns (bytes32) {
    return keccak256(abi.encode(labelledName, version));
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
    bool isPublic,
    bool acceptsWorkflows,
    uint8 f
  ) external onlyOwner {
    uint32 id = s_nextDONId++;
    s_dons[id].id = id;

    _setDONConfig(
      nodes,
      capabilityConfigurations,
      DONParams({id: id, configCount: 1, isPublic: isPublic, acceptsWorkflows: acceptsWorkflows, f: f})
    );
  }

  /// @notice Updates a DON's configuration.  This allows
  /// the admin to reconfigure the list of capabilities supported
  /// by the DON, the list of nodes that make up the DON as well
  /// as whether or not the DON can accept external workflows
  /// @param nodes The nodes making up the DON
  /// @param capabilityConfigurations The list of configurations for the
  /// capabilities supported by the DON
  /// @param isPublic True if the DON is can accept external workflows
  function updateDON(
    uint32 donId,
    bytes32[] calldata nodes,
    CapabilityConfiguration[] calldata capabilityConfigurations,
    bool isPublic,
    bool acceptsWorkflows,
    uint8 f
  ) external onlyOwner {
    uint32 configCount = s_dons[donId].configCount;
    if (configCount == 0) revert DONDoesNotExist(donId);
    _setDONConfig(
      nodes,
      capabilityConfigurations,
      DONParams({id: donId, configCount: ++configCount, isPublic: isPublic, acceptsWorkflows: acceptsWorkflows, f: f})
    );
  }

  /// @notice Removes DONs from the Capability Registry
  /// @param donIds The IDs of the DON to be removed
  function removeDONs(uint32[] calldata donIds) external onlyOwner {
    for (uint256 i; i < donIds.length; ++i) {
      uint32 donId = donIds[i];
      DON storage don = s_dons[donId];

      uint32 configCount = don.configCount;
      EnumerableSet.Bytes32Set storage nodeP2PIds = don.config[configCount].nodes;

      bool isWorkflowDON = don.acceptsWorkflows;
      for (uint256 j; j < nodeP2PIds.length(); ++j) {
        if (isWorkflowDON) {
          delete s_nodes[nodeP2PIds.at(j)].workflowDONId;
        } else {
          s_nodes[nodeP2PIds.at(j)].capabilitiesDONIds.remove(donId);
        }
      }

      // DON config count starts at index 1
      if (don.configCount == 0) revert DONDoesNotExist(donId);
      delete s_dons[donId];
      emit ConfigSet(donId, 0);
    }
  }

  /// @notice Gets DON's data
  /// @param donId The DON ID
  /// @return DONInfo The DON's parameters
  function getDON(uint32 donId) external view returns (DONInfo memory) {
    return _getDON(donId);
  }

  /// @notice Returns the list of configured DONs
  /// @return DONInfo[] The list of configured DONs
  function getDONs() external view returns (DONInfo[] memory) {
    /// Minus one to account for s_nextDONId starting at index 1
    uint32 donId = s_nextDONId;
    DONInfo[] memory dons = new DONInfo[](donId - 1);
    uint256 idx;
    ///
    for (uint32 i = 1; i < donId; ++i) {
      if (s_dons[i].id != 0) {
        dons[idx] = _getDON(i);
        ++idx;
      }
    }
    if (idx != donId - 1) {
      assembly {
        mstore(dons, idx)
      }
    }
    return dons;
  }

  /// @notice Returns the DON specific configuration for a capability
  /// @param donId The DON's ID
  /// @param capabilityId The Capability ID
  /// @return bytes The DON specific configuration for the capability stored on the capability registry
  /// @return bytes The DON specific configuration stored on the capability's configuration contract
  function getCapabilityConfigs(uint32 donId, bytes32 capabilityId) external view returns (bytes memory, bytes memory) {
    uint32 configCount = s_dons[donId].configCount;

    bytes memory donCapabilityConfig = s_dons[donId].config[configCount].capabilityConfigs[capabilityId];
    bytes memory globalCapabilityConfig;

    if (s_capabilities[capabilityId].configurationContract != address(0)) {
      globalCapabilityConfig = ICapabilityConfiguration(s_capabilities[capabilityId].configurationContract)
        .getCapabilityConfiguration(donId);
    }

    return (donCapabilityConfig, globalCapabilityConfig);
  }

  /// @notice Sets the configuration for a DON
  /// @param nodes The nodes making up the DON
  /// @param capabilityConfigurations The list of configurations for the
  /// capabilities supported by the DON
  /// @param donParams The DON's parameters
  function _setDONConfig(
    bytes32[] calldata nodes,
    CapabilityConfiguration[] calldata capabilityConfigurations,
    DONParams memory donParams
  ) internal {
    DONCapabilityConfig storage donCapabilityConfig = s_dons[donParams.id].config[donParams.configCount];

    // Validate the f value. We are intentionally relaxing the 3f+1 requirement
    // as not all DONs will run OCR instances.
    if (donParams.f == 0 || donParams.f + 1 > nodes.length) revert InvalidFaultTolerance(donParams.f, nodes.length);

    // Skip removing supported DON Ids from previously configured nodes in DON if
    // we are adding the DON for the first time
    if (donParams.configCount > 1) {
      DONCapabilityConfig storage prevDONCapabilityConfig = s_dons[donParams.id].config[donParams.configCount - 1];

      // We acknowledge that this may result in an out of gas error if the number of configured
      // nodes is large.  This is mitigated by ensuring that there will not be a large number
      // of nodes configured to a DON.
      // We also do not remove the nodes from the previous DON capability config.  This is not
      // needed as the previous config will be overwritten by storing the latest config
      // at configCount
      for (uint256 i; i < prevDONCapabilityConfig.nodes.length(); ++i) {
        s_nodes[prevDONCapabilityConfig.nodes.at(i)].capabilitiesDONIds.remove(donParams.id);
        delete s_nodes[prevDONCapabilityConfig.nodes.at(i)].workflowDONId;
      }
    }

    for (uint256 i; i < nodes.length; ++i) {
      if (!donCapabilityConfig.nodes.add(nodes[i])) revert DuplicateDONNode(donParams.id, nodes[i]);

      if (donParams.acceptsWorkflows) {
        if (s_nodes[nodes[i]].workflowDONId != donParams.id && s_nodes[nodes[i]].workflowDONId != 0)
          revert NodePartOfWorkflowDON(donParams.id, nodes[i]);
        s_nodes[nodes[i]].workflowDONId = donParams.id;
      } else {
        /// Fine to add a duplicate DON ID to the set of supported DON IDs again as the set
        /// will only store unique DON IDs
        s_nodes[nodes[i]].capabilitiesDONIds.add(donParams.id);
      }
    }

    for (uint256 i; i < capabilityConfigurations.length; ++i) {
      CapabilityConfiguration calldata configuration = capabilityConfigurations[i];

      if (!s_hashedCapabilityIds.contains(configuration.capabilityId))
        revert CapabilityDoesNotExist(configuration.capabilityId);
      if (s_deprecatedHashedCapabilityIds.contains(configuration.capabilityId))
        revert CapabilityIsDeprecated(configuration.capabilityId);

      if (donCapabilityConfig.capabilityConfigs[configuration.capabilityId].length > 0)
        revert DuplicateDONCapability(donParams.id, configuration.capabilityId);

      for (uint256 j; j < nodes.length; ++j) {
        if (
          !s_nodes[nodes[j]].supportedHashedCapabilityIds[s_nodes[nodes[j]].configCount].contains(
            configuration.capabilityId
          )
        ) revert NodeDoesNotSupportCapability(nodes[j], configuration.capabilityId);
      }

      donCapabilityConfig.capabilityIds.push(configuration.capabilityId);
      donCapabilityConfig.capabilityConfigs[configuration.capabilityId] = configuration.config;

      _setDONCapabilityConfig(
        donParams.id,
        donParams.configCount,
        configuration.capabilityId,
        nodes,
        configuration.config
      );
    }
    s_dons[donParams.id].isPublic = donParams.isPublic;
    s_dons[donParams.id].acceptsWorkflows = donParams.acceptsWorkflows;
    s_dons[donParams.id].f = donParams.f;
    s_dons[donParams.id].configCount = donParams.configCount;
    emit ConfigSet(donParams.id, donParams.configCount);
  }

  /// @notice Sets the capability's config on the config contract
  /// @param donId The ID of the DON the capability is being configured for
  /// @param configCount The number of times the DON has been configured
  /// @param capabilityId The capability's ID
  /// @param nodes The nodes in the DON
  /// @param config The DON's capability config
  /// @dev Helper function used to resolve stack too deep errors in _setDONConfig
  function _setDONCapabilityConfig(
    uint32 donId,
    uint32 configCount,
    bytes32 capabilityId,
    bytes32[] calldata nodes,
    bytes memory config
  ) internal {
    if (s_capabilities[capabilityId].configurationContract != address(0)) {
      ICapabilityConfiguration(s_capabilities[capabilityId].configurationContract).beforeCapabilityConfigSet(
        nodes,
        config,
        configCount,
        donId
      );
    }
  }

  /// @notice Sets a capability's data
  /// @param hashedCapabilityId The ID of the capability being set
  /// @param capability The capability's data
  function _setCapability(bytes32 hashedCapabilityId, Capability memory capability) internal {
    if (capability.configurationContract != address(0)) {
      /// Check that the configuration contract being assigned
      /// correctly supports the ICapabilityConfiguration interface
      /// by implementing both getCapabilityConfiguration and
      /// beforeCapabilityConfigSet
      if (
        capability.configurationContract.code.length == 0 ||
        !IERC165(capability.configurationContract).supportsInterface(type(ICapabilityConfiguration).interfaceId)
      ) revert InvalidCapabilityConfigurationContractInterface(capability.configurationContract);
    }
    s_capabilities[hashedCapabilityId] = capability;
    emit CapabilityConfigured(hashedCapabilityId);
  }

  /// @notice Gets DON's data
  /// @param donId The DON ID
  /// @return DONInfo The DON's parameters
  function _getDON(uint32 donId) internal view returns (DONInfo memory) {
    uint32 configCount = s_dons[donId].configCount;

    DONCapabilityConfig storage donCapabilityConfig = s_dons[donId].config[configCount];

    bytes32[] memory capabilityIds = donCapabilityConfig.capabilityIds;
    CapabilityConfiguration[] memory capabilityConfigurations = new CapabilityConfiguration[](capabilityIds.length);

    for (uint256 i; i < capabilityConfigurations.length; ++i) {
      capabilityConfigurations[i] = CapabilityConfiguration({
        capabilityId: capabilityIds[i],
        config: donCapabilityConfig.capabilityConfigs[capabilityIds[i]]
      });
    }

    return
      DONInfo({
        id: s_dons[donId].id,
        configCount: configCount,
        f: s_dons[donId].f,
        isPublic: s_dons[donId].isPublic,
        acceptsWorkflows: s_dons[donId].acceptsWorkflows,
        nodeP2PIds: donCapabilityConfig.nodes.values(),
        capabilityConfigurations: capabilityConfigurations
      });
  }
}
