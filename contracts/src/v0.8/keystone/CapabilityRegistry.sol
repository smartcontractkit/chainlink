// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";
import {IERC165} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {EnumerableSet} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";
import {ICapabilityConfiguration} from "./interfaces/ICapabilityConfiguration.sol";

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
    // Capability type, e.g. "data-streams-reports"
    // bytes32(string); validation regex: ^[a-z0-9_\-:]{1,32}$
    // Not "type" because that's a reserved keyword in Solidity.
    bytes32 capabilityType;
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

  /// @notice This error is thrown when trying add a capability that already
  /// exists.
  error CapabilityAlreadyExists();

  /// @notice This error is thrown when a capability with the provided ID is
  /// not found.
  /// @param capabilityId The ID used for the lookup.
  error CapabilityDoesNotExist(bytes32 capabilityId);

  /// @notice This error is thrown when trying to deprecate a capability that
  /// is already deprecated.
  /// @param capabilityId The ID of the capability that is already deprecated.
  error CapabilityAlreadyDeprecated(bytes32 capabilityId);

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
  /// @param capabilityId The ID of the newly added capability
  event CapabilityAdded(bytes32 indexed capabilityId);

  /// @notice This event is emitted when a capability is deprecated
  /// @param capabilityId The ID of the deprecated capability
  event CapabilityDeprecated(bytes32 indexed capabilityId);

  mapping(bytes32 => Capability) private s_capabilities;
  EnumerableSet.Bytes32Set private s_capabilityIds;
  EnumerableSet.Bytes32Set private s_deprecatedCapabilityIds;

  /// @notice Mapping of node operators
  mapping(uint256 nodeOperatorId => NodeOperator) private s_nodeOperators;

  /// @notice The latest node operator ID
  /// @dev No getter for this as this is an implementation detail
  uint256 private s_nodeOperatorId;

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

  function addCapability(Capability calldata capability) external onlyOwner {
    bytes32 capabilityId = getCapabilityID(capability.capabilityType, capability.version);

    if (s_capabilityIds.contains(capabilityId)) revert CapabilityAlreadyExists();

    if (capability.configurationContract != address(0)) {
      if (
        capability.configurationContract.code.length == 0 ||
        !IERC165(capability.configurationContract).supportsInterface(
          ICapabilityConfiguration.getCapabilityConfiguration.selector
        )
      ) revert InvalidCapabilityConfigurationContractInterface(capability.configurationContract);
    }

    s_capabilityIds.add(capabilityId);
    s_capabilities[capabilityId] = capability;

    emit CapabilityAdded(capabilityId);
  }

  /// @notice Deprecates a capability by adding it to the deprecated list
  /// @param capabilityId The ID of the capability to deprecate
  function deprecateCapability(bytes32 capabilityId) external onlyOwner {
    if (!s_capabilityIds.contains(capabilityId)) revert CapabilityDoesNotExist(capabilityId);
    if (s_deprecatedCapabilityIds.contains(capabilityId)) revert CapabilityAlreadyDeprecated(capabilityId);

    s_deprecatedCapabilityIds.add(capabilityId);
    emit CapabilityDeprecated(capabilityId);
  }

  function getCapability(bytes32 capabilityID) public view returns (Capability memory) {
    return s_capabilities[capabilityID];
  }

  /// @notice Returns all capabilities. This operation will copy capabilities
  /// to memory, which can be quite expensive. This is designed to mostly be
  /// used by view accessors that are queried without any gas fees.
  /// @return Capability[] An array of capabilities
  function getCapabilities() external view returns (Capability[] memory) {
    bytes32[] memory capabilityIds = s_capabilityIds.values();

    // Solidity does not support dynamic arrays in memory, so we create a
    // fixed-size array and copy the capabilities into it.
    Capability[] memory capabilities = new Capability[](capabilityIds.length - s_deprecatedCapabilityIds.length());

    // We need to keep track of the new index because we are skipping
    // deprecated capabilities.
    uint256 newIndex;

    for (uint256 i; i < capabilityIds.length; ++i) {
      bytes32 capabilityId = capabilityIds[i];

      if (!s_deprecatedCapabilityIds.contains(capabilityId)) {
        capabilities[newIndex] = getCapability(capabilityId);
        newIndex++;
      }
    }

    return capabilities;
  }

  /// @notice This functions returns a Capability ID packed into a bytes32 for cheaper access
  /// @return bytes32 A unique identifier for the capability
  function getCapabilityID(bytes32 capabilityType, bytes32 version) public pure returns (bytes32) {
    return keccak256(abi.encodePacked(capabilityType, version));
  }

  /// @notice Returns whether a capability is deprecated
  /// @param capabilityId The ID of the capability to check
  /// @return bool True if the capability is deprecated, false otherwise
  function isCapabilityDeprecated(bytes32 capabilityId) external view returns (bool) {
    return s_deprecatedCapabilityIds.contains(capabilityId);
  }
}
