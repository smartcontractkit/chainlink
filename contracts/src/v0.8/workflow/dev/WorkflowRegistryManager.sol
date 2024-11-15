// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";

import {Ownable2StepMsgSender} from "../../shared/access/Ownable2StepMsgSender.sol";

/// @title WorkflowRegistryManager
/// @notice This contract manages the versions of WorkflowRegistry contracts deployed over time.
/// @dev This contract allows the owner to add, activate, and manage versions of WorkflowRegistry contracts. It tracks
/// deployment information for each version, including deployment timestamp, chain ID, and active status. Only one
/// version can be active at any given time.
contract WorkflowRegistryManager is Ownable2StepMsgSender, ITypeAndVersion {
  string public constant override typeAndVersion = "WorkflowRegistryManager 1.0.0";
  uint8 private constant MAX_PAGINATION_LIMIT = 100;

  struct Version {
    address contractAddress; // ─╮ Address of the WorkflowRegistry contract
    uint64 chainID; //           │ Chain ID of the EVM chain where the WorkflowRegistry is deployed.
    uint32 deployedAt; // ───────╯ Block timestamp of deployment (sufficient until year 2106).
    string contractTypeAndVersion; // WorkflowRegistry's typeAndVersion.
  }

  /// @notice Maps version numbers to their corresponding `Version` details.
  /// @dev This mapping is 1-based, meaning version numbers start from 1. Ensure that all operations account for this
  /// indexing strategy to avoid off-by-one errors.
  mapping(uint32 versionNumber => Version versionInfo) private s_versions;

  /// @notice The version number of the currently active WorkflowRegistry.
  /// @dev Initialized to `type(uint32).max` to indicate no active version. Updated when a version is activated.
  uint32 private s_activeVersionNumber = type(uint32).max;

  /// @notice The latest version number registered in the contract.
  /// @dev Incremented each time a new version is added. Useful for iterating over all registered versions.
  uint32 private s_latestVersionNumber = 0;

  // Errors
  error InvalidContractAddress(address invalidAddress);
  error InvalidContractType(address invalidAddress);
  error NoActiveVersionAvailable();
  error NoVersionsRegistered();
  error VersionNotRegistered(uint32 versionNumber);
  // Events

  event VersionAddedV1(address indexed contractAddress, uint64 chainID, uint32 deployedAt, uint32 version);
  event VersionActivatedV1(address indexed contractAddress, uint64 chainID, uint32 indexed version);
  event VersionDeactivatedV1(address indexed contractAddress, uint64 chainID, uint32 indexed version);

  // ================================================================
  // |                            ADMIN                             |
  // ================================================================
  /// @notice Adds a new WorkflowRegistry version to the version history and optionally activates it.
  /// @dev This function records the deployment details of a new registry version. It deactivates the currently active
  /// version (if any) and activates the newly added version if `autoActivate` is true.
  /// @param contractAddress The address of the deployed WorkflowRegistry contract. Must be a valid contract address.
  /// @param chainID The chain ID of the EVM chain where the WorkflowRegistry is deployed.
  /// @param autoActivate A boolean indicating whether the new version should be activated immediately.
  /// @custom:throws InvalidContractType if the provided contract address is zero or not a WorkflowRegistry.
  function addVersion(address contractAddress, uint64 chainID, bool autoActivate) external onlyOwner {
    string memory typeVer = _getTypeAndVersionForContract(contractAddress);
    uint32 latestVersionNumber = ++s_latestVersionNumber;
    uint32 deployedAt = uint32(block.timestamp);

    s_versions[latestVersionNumber] = Version({
      contractAddress: contractAddress,
      chainID: chainID,
      deployedAt: deployedAt,
      contractTypeAndVersion: typeVer
    });

    if (autoActivate) {
      _activateVersion(latestVersionNumber);
    }

    emit VersionAddedV1(contractAddress, chainID, deployedAt, latestVersionNumber);
  }

  /// @notice Activates a specific WorkflowRegistry version by its version number.
  /// @dev This contract uses a 1-based index, meaning the `versionNumber` parameter must start at 1, with 1 representing the
  /// first version. Setting `versionNumber` to 0 will revert, as 0 is not a valid index in this context. Only one version
  /// can be active at a time; activating a new version automatically deactivates the currently active one (if any).
  /// @param versionNumber The 1-based version number to activate (minimum value is 1).
  /// @custom:throws VersionNotRegistered if the `versionNumber` is not valid or not registered.
  function activateVersion(
    uint32 versionNumber
  ) external onlyOwner {
    _activateVersion(versionNumber);
  }

  // ================================================================
  // |                         GET VERSIONS                         |
  // ================================================================
  /// @notice Returns a paginated list of all WorkflowRegistry versions.
  /// @dev This function retrieves a range of versions based on the provided `start` and `limit` parameters. The contract uses
  /// a 1-based index, so the `start` parameter must be at least 1, representing the first version. If `limit` is set to
  /// 0 or exceeds `MAX_PAGINATION_LIMIT`, it defaults to `MAX_PAGINATION_LIMIT`. If `start` exceeds the total number of
  /// versions, an empty array is returned.
  /// @param start The index at which to start retrieving versions (1-based index, minimum value is 1).
  /// @param limit The maximum number of versions to retrieve (maximum is `MAX_PAGINATION_LIMIT`).
  /// @return versions An array of `Version` structs containing version details, starting from the `start` index up to the
  /// specified `limit`.
  function getAllVersions(uint32 start, uint32 limit) external view returns (Version[] memory versions) {
    uint32 totalVersions = s_latestVersionNumber;

    // Adjust for 1-based index
    if (start == 0 || start > totalVersions) {
      return new Version[](0);
    }

    if (limit > MAX_PAGINATION_LIMIT || limit == 0) {
      limit = MAX_PAGINATION_LIMIT;
    }

    uint32 end = (start + limit - 1 > totalVersions) ? totalVersions : start + limit - 1;
    uint32 resultLength = end - start + 1;

    versions = new Version[](resultLength);
    for (uint32 i = 0; i < resultLength; ++i) {
      versions[i] = s_versions[start + i];
    }

    return versions;
  }

  /// @notice Retrieves the details of a specific WorkflowRegistry version by its version number.
  /// @dev This contract uses a 1-based index, so `versionNumber` must be at least 1. This means the first version is
  /// represented by `versionNumber` of 1, not 0. Attempting to retrieve a version with a `versionNumber` of 0 or exceeding
  /// `s_latestVersionNumber` will revert.
  /// @param versionNumber The 1-based version number of the version to retrieve (minimum value is 1).
  /// @return A `Version` struct containing the details of the specified version.
  /// @custom:throws VersionNotRegistered if the `versionNumber` is not valid or not registered.
  function getVersion(
    uint32 versionNumber
  ) external view returns (Version memory) {
    if (versionNumber == 0 || versionNumber > s_latestVersionNumber) {
      revert VersionNotRegistered(versionNumber);
    }
    return s_versions[versionNumber];
  }

  /// @notice Retrieves the details of the currently active WorkflowRegistry version.
  /// @dev Assumes there is only one active version. Throws if no version is currently active.
  /// @return A `Version` struct containing the details of the active version.
  /// @custom:throws NoActiveVersionAvailable if no version is currently active.
  function getActiveVersion() external view returns (Version memory) {
    uint32 activateVersionNumber = s_activeVersionNumber;
    if (activateVersionNumber == type(uint32).max) revert NoActiveVersionAvailable();
    return s_versions[activateVersionNumber];
  }

  /// @notice Retrieves the details of the latest registered WorkflowRegistry version.
  /// @return A `Version` struct containing the details of the latest version.
  /// @custom:throws NoActiveVersionAvailable if no versions have been registered.
  function getLatestVersion() external view returns (Version memory) {
    uint32 latestVersionNumber = s_latestVersionNumber;
    if (latestVersionNumber == 0) revert NoActiveVersionAvailable();
    return s_versions[latestVersionNumber];
  }

  /// @notice Retrieves the version number of the currently active WorkflowRegistry version.
  /// @return activeVersionNumber The version number of the active version.
  /// @custom:throws NoActiveVersionAvailable if s_activeVersionNumber is `type(uint32).max`.
  function getActiveVersionNumber() external view returns (uint32 activeVersionNumber) {
    activeVersionNumber = s_activeVersionNumber;
    if (activeVersionNumber == type(uint32).max) revert NoActiveVersionAvailable();
    return activeVersionNumber;
  }

  /// @notice Retrieves the version number of the latest registered WorkflowRegistry version.
  /// @return latestVersionNumber The version number of the latest version.
  /// @custom:throws NoVersionsRegistered if s_latestVersionNumber is 0.
  function getLatestVersionNumber() external view returns (uint32 latestVersionNumber) {
    latestVersionNumber = s_latestVersionNumber;
    if (latestVersionNumber == 0) revert NoVersionsRegistered();
    return latestVersionNumber;
  }

  // ================================================================
  // |                          PRIVATE                             |
  // ================================================================
  /// @dev This function deactivates the currently active version (if any) before activating the specified version. It
  /// emits events for both deactivation and activation.
  /// @param versionNumber The version number of the version to activate.
  /// @custom:throws IndexOutOfBounds if the version number does not exist.
  function _activateVersion(
    uint32 versionNumber
  ) private {
    // Cache the current active version number to reduce storage reads
    uint32 currentActiveVersionNumber = s_activeVersionNumber;

    // Check that the provided version number is within a valid range
    if (versionNumber == 0 || versionNumber > s_latestVersionNumber) {
      revert VersionNotRegistered(versionNumber);
    }

    // Emit deactivation event if there is an active version
    if (currentActiveVersionNumber != type(uint32).max) {
      Version memory currentActive = s_versions[currentActiveVersionNumber];
      emit VersionDeactivatedV1(currentActive.contractAddress, currentActive.chainID, currentActiveVersionNumber);
    }

    // Set the new active version (which deactivates the previous one)
    s_activeVersionNumber = versionNumber;
    Version memory newActive = s_versions[versionNumber];
    emit VersionActivatedV1(newActive.contractAddress, newActive.chainID, versionNumber);
  }

  /// @dev Validates that a given contract address is non-zero, contains code, and supports the IWorkflowRegistry interface.
  /// @param contractAddress The address of the contract to validate.
  /// @custom:throws InvalidContractAddress if the address is zero or contains no code.
  /// @custom:throws InvalidContractType if the contract does not implement typeAndVersion().
  function _getTypeAndVersionForContract(
    address contractAddress
  ) internal view returns (string memory) {
    if (!_isNonZeroWithCode(contractAddress)) {
      revert InvalidContractAddress(contractAddress);
    }

    try ITypeAndVersion(contractAddress).typeAndVersion() returns (string memory retrievedVersion) {
      return retrievedVersion;
    } catch {
      revert InvalidContractType(contractAddress);
    }
  }

  function _isNonZeroWithCode(
    address _addr
  ) internal view returns (bool) {
    return _addr != address(0) && _addr.code.length > 0;
  }
}
