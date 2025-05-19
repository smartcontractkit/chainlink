// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";

import {Ownable2StepMsgSender} from "../shared/access/Ownable2StepMsgSender.sol";

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

  /// @notice Maps a combination of address and chain ID to the version number.
  /// @dev This mapping allows for lookup of the version number for a given address and chain ID.
  mapping(bytes32 => uint32) private s_versionNumberByAddressAndChainID;

  /// @notice The version number of the currently active WorkflowRegistry.
  /// @dev Initialized to 0 to indicate no active version. Updated when a version is activated.
  uint32 private s_activeVersionNumber = 0;

  /// @notice The latest version number registered in the contract.
  /// @dev Incremented each time a new version is added. Useful for iterating over all registered versions.
  uint32 private s_latestVersionNumber = 0;

  // Errors
  error ContractAlreadyRegistered(address contractAddress, uint64 chainID);
  error InvalidContractAddress(address invalidAddress);
  error InvalidContractType(address invalidAddress);
  error NoActiveVersionAvailable();
  error NoVersionsRegistered();
  error VersionAlreadyActive(uint32 versionNumber);
  error VersionNotRegistered(uint32 versionNumber);
  // Events

  event VersionAdded(address indexed contractAddress, uint64 chainID, uint32 deployedAt, uint32 version);
  event VersionActivated(address indexed contractAddress, uint64 chainID, uint32 version);
  event VersionDeactivated(address indexed contractAddress, uint64 chainID, uint32 version);

  // ================================================================
  // |                      Manage Versions                         |
  // ================================================================

  /// @notice Adds a new WorkflowRegistry version to the version history and optionally activates it.
  /// @dev This function records the deployment details of a new registry version. It deactivates the currently active
  /// version (if any) and activates the newly added version if `autoActivate` is true.
  /// @param contractAddress The address of the deployed WorkflowRegistry contract. Must be a valid contract address.
  /// @param chainID The chain ID of the EVM chain where the WorkflowRegistry is deployed.
  /// @param autoActivate A boolean indicating whether the new version should be activated immediately.
  /// @custom:throws InvalidContractType if the provided contract address is zero or not a WorkflowRegistry.
  function addVersion(address contractAddress, uint64 chainID, uint32 deployedAt, bool autoActivate) external onlyOwner {
    // Check if the contract is already registered. If it is, you can just activate that existing version.
    bytes32 key = keccak256(abi.encodePacked(contractAddress, chainID));
    if (s_versionNumberByAddressAndChainID[key] != 0) {
      revert ContractAlreadyRegistered(contractAddress, chainID);
    }

    string memory typeVer = _getTypeAndVersionForContract(contractAddress);
    uint32 latestVersionNumber = ++s_latestVersionNumber;

    s_versions[latestVersionNumber] = Version({
      contractAddress: contractAddress,
      chainID: chainID,
      deployedAt: deployedAt,
      contractTypeAndVersion: typeVer
    });

    // Store the version number associated with the hash of contract address and chainID
    s_versionNumberByAddressAndChainID[key] = latestVersionNumber;

    if (autoActivate) {
      _activateVersion(latestVersionNumber);
    }

    emit VersionAdded(contractAddress, chainID, deployedAt, latestVersionNumber);
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

  /// @dev This private function deactivates the currently active version (if any) before activating the specified version. It
  /// emits events for both deactivation and activation.
  /// @param versionNumber The version number of the version to activate.
  /// @custom:throws IndexOutOfBounds if the version number does not exist.
  function _activateVersion(
    uint32 versionNumber
  ) private {
    // Check that the provided version number is within a valid range
    if (versionNumber == 0 || versionNumber > s_latestVersionNumber) {
      revert VersionNotRegistered(versionNumber);
    }

    // Cache the current active version number to reduce storage reads
    uint32 currentActiveVersionNumber = s_activeVersionNumber;

    // Check that the version number is not the same as the current active version number
    if (currentActiveVersionNumber == versionNumber) {
      revert VersionAlreadyActive(versionNumber);
    }

    // Emit deactivation event if there is an active version
    if (currentActiveVersionNumber != 0) {
      Version memory currentActive = s_versions[currentActiveVersionNumber];
      emit VersionDeactivated(currentActive.contractAddress, currentActive.chainID, currentActiveVersionNumber);
    }

    // Set the new active version (which deactivates the previous one)
    s_activeVersionNumber = versionNumber;
    Version memory newActive = s_versions[versionNumber];
    emit VersionActivated(newActive.contractAddress, newActive.chainID, versionNumber);
  }

  // ================================================================
  // |                        Query Versions                        |
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

  /// @notice Retrieves the version number for a specific WorkflowRegistry by its contract address and chain ID.
  /// @param contractAddress The address of the WorkflowRegistry contract.
  /// @param chainID The chain ID of the network where the WorkflowRegistry is deployed.
  /// @return versionNumber The version number associated with the given contract address and chain ID.
  function getVersionNumberByContractAddressAndChainID(
    address contractAddress,
    uint64 chainID
  ) external view returns (uint32 versionNumber) {
    _validateContractAddress(contractAddress);

    bytes32 key = keccak256(abi.encodePacked(contractAddress, chainID));
    versionNumber = s_versionNumberByAddressAndChainID[key];
    if (versionNumber == 0) {
      revert NoVersionsRegistered();
    }
    return versionNumber;
  }

  /// @notice Retrieves the details of the currently active WorkflowRegistry version.
  /// @dev Assumes there is only one active version. Throws if no version is currently active.
  /// @return A `Version` struct containing the details of the active version.
  /// @custom:throws NoActiveVersionAvailable if no version is currently active.
  function getActiveVersion() external view returns (Version memory) {
    uint32 activeVersionNumber = s_activeVersionNumber;
    if (activeVersionNumber == 0) revert NoActiveVersionAvailable();
    return s_versions[activeVersionNumber];
  }

  /// @notice Retrieves the details of the latest registered WorkflowRegistry version.
  /// @return A `Version` struct containing the details of the latest version.
  /// @custom:throws NoVersionsRegistered if no versions have been registered.
  function getLatestVersion() external view returns (Version memory) {
    uint32 latestVersionNumber = s_latestVersionNumber;
    if (latestVersionNumber == 0) revert NoVersionsRegistered();
    return s_versions[latestVersionNumber];
  }

  /// @notice Retrieves the version number of the currently active WorkflowRegistry version.
  /// @return activeVersionNumber The version number of the active version.
  /// @custom:throws NoActiveVersionAvailable if s_activeVersionNumber is `type(uint32).max`.
  function getActiveVersionNumber() external view returns (uint32 activeVersionNumber) {
    activeVersionNumber = s_activeVersionNumber;
    if (activeVersionNumber == 0) revert NoActiveVersionAvailable();
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
  // |                          Validation                          |
  // ================================================================

  /// @dev Validates that a given contract address is non-zero, contains code, and implements typeAndVersion().
  /// @param contractAddress The address of the contract to validate.
  /// @custom:throws InvalidContractAddress if the address is zero or contains no code.
  /// @custom:throws InvalidContractType if the contract does not implement typeAndVersion().
  function _getTypeAndVersionForContract(
    address contractAddress
  ) internal view returns (string memory) {
    _validateContractAddress(contractAddress);

    try ITypeAndVersion(contractAddress).typeAndVersion() returns (string memory retrievedVersion) {
      return retrievedVersion;
    } catch {
      revert InvalidContractType(contractAddress);
    }
  }

  /// @dev Validates that a given contract address is non-zero and contains code.
  /// @param _addr The address of the contract to validate.
  /// @custom:throws InvalidContractAddress if the address is zero or contains no code.
  function _validateContractAddress(
    address _addr
  ) internal view {
    if (_addr == address(0) || _addr.code.length == 0) {
      revert InvalidContractAddress(_addr);
    }
  }
}
