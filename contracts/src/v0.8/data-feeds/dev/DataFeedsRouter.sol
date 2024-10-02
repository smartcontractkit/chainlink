// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IDataFeedsRouter} from "./interfaces/IDataFeedsRouter.sol";
import {IDataFeedsRegistry} from "./interfaces/IDataFeedsRegistry.sol";
import {IDataFeedsFeeManager} from "./interfaces/IDataFeedsFeeManager.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";
import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";

contract DataFeedsRouter is IDataFeedsRouter, ConfirmedOwner, TypeAndVersionInterface {
  using EnumerableSet for EnumerableSet.AddressSet;

  string public constant override typeAndVersion = "DataFeedsRouter 1.0.0";

  /// @notice Mapping from feed id to registry address
  mapping(bytes32 => address) internal s_feedIdToRegistry;

  /// @notice Set of authorized users for zero-fee data access
  EnumerableSet.AddressSet internal s_authorizedNonbillableUsers;

  IDataFeedsFeeManager internal s_feeManager;

  event NonbillableUserAdded(address user);
  event NonbillableUserRemoved(address user);

  error UnauthorizedNonbillableAccess();

  modifier authorizeNonbillableAccess() {
    if (!s_authorizedNonbillableUsers.contains(msg.sender)) revert UnauthorizedNonbillableAccess();

    _;
  }

  constructor() ConfirmedOwner(msg.sender) {
    s_authorizedNonbillableUsers.add(msg.sender);
    s_authorizedNonbillableUsers.add(address(0));
  }

  // Given a list of feedIds and a mapping from feedId to address, return the unique addresses associated with those feedIds
  function _getUniqueAssignedAddresses(
    bytes32[] calldata feedIds,
    mapping(bytes32 => address) storage feedIdToAddress
  ) internal view returns (address[] memory) {
    address[] memory uniqueAddresses = new address[](feedIds.length);
    uint256 lastAddressIndex = 0;

    for (uint256 i; i < feedIds.length; ++i) {
      address assignedAddress = feedIdToAddress[feedIds[i]];

      // Check if assignedAddress is already present in the list of unique addresses
      bool addressPresent = false;
      for (uint256 j = 0; j < uniqueAddresses.length; j++) {
        if (uniqueAddresses[j] == assignedAddress) {
          addressPresent = true;
          break;
        }
      }

      // If assignedAddress is not already present in the list of unique addresses, add it
      if (!addressPresent) {
        uniqueAddresses[lastAddressIndex] = assignedAddress;
        lastAddressIndex++;
      }
    }

    // Trim the size of the uniqueAddresses array down to only include set values
    address[] memory trimmedAddresses = new address[](lastAddressIndex);
    for (uint256 i; i < lastAddressIndex; ++i) {
      trimmedAddresses[i] = uniqueAddresses[i];
    }

    return trimmedAddresses;
  }

  // Given a list of feed IDs, an address and a mapping from feedID to address
  // return the feedIds that are associated with that address, and the index of the feedIds in the original list.
  // This is used to combine results fetched from multiple contracts back into a single list, in the correct order.
  function _getAssignedFeedIdsForAddress(
    bytes32[] calldata feedIds,
    address assignedAddress,
    mapping(bytes32 => address) storage feedIdToAddress
  ) internal view returns (bytes32[] memory, uint256[] memory) {
    bytes32[] memory assignedFeedIds = new bytes32[](feedIds.length);
    uint256[] memory feedsIdsIndexes = new uint256[](feedIds.length);

    uint256 lastFeedIdIndex = 0;
    for (uint256 i; i < feedIds.length; ++i) {
      if (feedIdToAddress[feedIds[i]] == assignedAddress) {
        assignedFeedIds[lastFeedIdIndex] = feedIds[i];
        feedsIdsIndexes[lastFeedIdIndex] = i;
        lastFeedIdIndex++;
      }
    }

    // Trim the size of the arrays down to only include set values
    bytes32[] memory trimmedAssignedFeedIds = new bytes32[](lastFeedIdIndex);
    uint256[] memory trimmedFeedsIdsIndexes = new uint256[](lastFeedIdIndex);
    for (uint256 i; i < lastFeedIdIndex; ++i) {
      trimmedAssignedFeedIds[i] = assignedFeedIds[i];
      trimmedFeedsIdsIndexes[i] = feedsIdsIndexes[i];
    }

    return (trimmedAssignedFeedIds, trimmedFeedsIdsIndexes);
  }

  function _getBenchmarks(
    bytes32[] calldata feedIds
  ) internal view returns (int256[] memory benchmarks, uint256[] memory observationTimestamps) {
    benchmarks = new int256[](feedIds.length);
    observationTimestamps = new uint256[](feedIds.length);

    address[] memory uniqueRegistries = _getUniqueAssignedAddresses(feedIds, s_feedIdToRegistry);

    // For each Registry, fetch the benchmarks and timestamps for the feedIds associated with that Registry
    // and map the results back to the index of the feedIds in the original list
    for (uint256 i; i < uniqueRegistries.length; ++i) {
      address registryAddress = uniqueRegistries[i];

      (bytes32[] memory registryFeedIds, uint256[] memory feedIdsIdxs) = _getAssignedFeedIdsForAddress(
        feedIds,
        registryAddress,
        s_feedIdToRegistry
      );
      IDataFeedsRegistry registry = IDataFeedsRegistry(registryAddress);

      (int256[] memory registryBenchmarks, uint256[] memory registryObservationTimestamps) = registry.getBenchmarks(
        registryFeedIds
      );

      // Map the results back to the index of the feedIds in the original list
      for (uint256 j; j < feedIdsIdxs.length; ++j) {
        uint256 feedIdIdx = feedIdsIdxs[j];
        benchmarks[feedIdIdx] = registryBenchmarks[j];
        observationTimestamps[feedIdIdx] = registryObservationTimestamps[j];
      }
    }

    return (benchmarks, observationTimestamps);
  }

  function _getReports(
    bytes32[] calldata feedIds
  ) internal view returns (bytes[] memory reports, uint256[] memory observationTimestamps) {
    reports = new bytes[](feedIds.length);
    observationTimestamps = new uint256[](feedIds.length);

    address[] memory uniqueRegistries = _getUniqueAssignedAddresses(feedIds, s_feedIdToRegistry);

    // For each Registry, fetch the benchmarks and timestamps for the feedIds associated with that Registry
    // and map the results back to the index of the feedIds in the original list
    for (uint256 i; i < uniqueRegistries.length; ++i) {
      address registryAddress = uniqueRegistries[i];

      (bytes32[] memory registryFeedIds, uint256[] memory feedIdsIdxs) = _getAssignedFeedIdsForAddress(
        feedIds,
        registryAddress,
        s_feedIdToRegistry
      );
      IDataFeedsRegistry registry = IDataFeedsRegistry(registryAddress);

      (bytes[] memory registryReports, uint256[] memory registryObservationTimestamps) = registry.getReports(
        registryFeedIds
      );

      // Map the results back to the index of the feedIds in the original list
      for (uint256 j; j < feedIdsIdxs.length; ++j) {
        uint256 feedIdIdx = feedIdsIdxs[j];
        reports[feedIdIdx] = registryReports[j];
        observationTimestamps[feedIdIdx] = registryObservationTimestamps[j];
      }
    }
    return (reports, observationTimestamps);
  }

  /// @inheritdoc IDataFeedsRouter
  function getBenchmarks(
    bytes32[] calldata feedIds,
    bytes calldata billingData
  ) external payable returns (int256[] memory benchmarks, uint256[] memory observationTimestamps) {
    s_feeManager.processFee{value: msg.value}(
      msg.sender,
      IDataFeedsFeeManager.Service.GetBenchmarks,
      feedIds,
      billingData
    );

    return _getBenchmarks(feedIds);
  }

  /// @inheritdoc IDataFeedsRouter
  function getBenchmarksNonbillable(
    bytes32[] calldata feedIds
  )
    external
    view
    authorizeNonbillableAccess
    returns (int256[] memory benchmarks, uint256[] memory observationTimestamps)
  {
    return _getBenchmarks(feedIds);
  }

  /// @inheritdoc IDataFeedsRouter
  function getReports(
    bytes32[] calldata feedIds,
    bytes calldata billingData
  ) external payable returns (bytes[] memory reports, uint256[] memory observationTimestamps) {
    s_feeManager.processFee{value: msg.value}(
      msg.sender,
      IDataFeedsFeeManager.Service.GetReports,
      feedIds,
      billingData
    );

    return _getReports(feedIds);
  }

  /// @inheritdoc IDataFeedsRouter
  function getReportsNonbillable(
    bytes32[] calldata feedIds
  ) external view authorizeNonbillableAccess returns (bytes[] memory reports, uint256[] memory observationTimestamps) {
    return _getReports(feedIds);
  }

  /// @inheritdoc IDataFeedsRouter
  function getDescriptions(bytes32[] calldata feedIds) external view returns (string[] memory descriptions) {
    descriptions = new string[](feedIds.length);

    address[] memory uniqueRegistries = _getUniqueAssignedAddresses(feedIds, s_feedIdToRegistry);

    // For each Registry, fetch the benchmarks and timestamps for the feedIds associated with that Registry
    // and map the results back to the index of the feedIds in the original list
    for (uint256 i; i < uniqueRegistries.length; ++i) {
      address registryAddress = uniqueRegistries[i];

      (bytes32[] memory registryFeedIds, uint256[] memory feedIdsIdxs) = _getAssignedFeedIdsForAddress(
        feedIds,
        registryAddress,
        s_feedIdToRegistry
      );
      IDataFeedsRegistry registry = IDataFeedsRegistry(registryAddress);

      (string[] memory registryDescriptions, , , , ) = registry.getFeedMetadata(registryFeedIds);

      // Map the results back to the index of the feedIds in the original list
      for (uint256 j; j < feedIdsIdxs.length; ++j) {
        uint256 feedIdIdx = feedIdsIdxs[j];
        descriptions[feedIdIdx] = registryDescriptions[j];
      }
    }

    return descriptions;
  }

  /// @inheritdoc IDataFeedsRouter
  function requestUpkeep(bytes32[] calldata feedIds, bytes calldata billingData) external payable {
    address[] memory uniqueRegistries = _getUniqueAssignedAddresses(feedIds, s_feedIdToRegistry);
    // For each Registry, fetch the benchmarks and timestamps for the feedIds associated with that Registry
    // and map the results back to the index of the feedIds in the original list
    for (uint256 i; i < uniqueRegistries.length; ++i) {
      address registryAddress = uniqueRegistries[i];

      (bytes32[] memory registryFeedIds, ) = _getAssignedFeedIdsForAddress(
        feedIds,
        registryAddress,
        s_feedIdToRegistry
      );

      IDataFeedsRegistry registry = IDataFeedsRegistry(registryAddress);
      registry.requestUpkeep(registryFeedIds);
    }

    s_feeManager.processFee{value: msg.value}(
      msg.sender,
      IDataFeedsFeeManager.Service.RequestUpkeep,
      feedIds,
      billingData
    );
  }

  /// @inheritdoc IDataFeedsRouter
  function configureFeeds(
    bytes32[] calldata feedIds,
    string[] calldata descriptions,
    bytes32 configId,
    address upkeep,
    address registryAddress,
    bytes32 feeConfigId
  ) external onlyOwner {
    IDataFeedsRegistry registry = IDataFeedsRegistry(registryAddress);
    registry.setFeeds(feedIds, descriptions, configId, upkeep);

    setRegistry(feedIds, registryAddress);

    s_feeManager.setFeedServiceFees(feeConfigId, feedIds);
  }

  /**
   * @notice indicates that the Registry address has been updated for a feed
   * @param feedId feed id
   * @param currentRegistry current Registry address
   * @param previousRegistry previous Registry address
   */
  event RegistrySet(bytes32 indexed feedId, address indexed currentRegistry, address indexed previousRegistry);

  /// @inheritdoc IDataFeedsRouter
  function setRegistry(bytes32[] calldata feedIds, address registryAddress) public onlyOwner {
    for (uint256 i; i < feedIds.length; ++i) {
      address previousRegistryAddress = s_feedIdToRegistry[feedIds[i]];
      s_feedIdToRegistry[feedIds[i]] = registryAddress;
      emit RegistrySet(feedIds[i], registryAddress, previousRegistryAddress);
    }
  }

  /// @inheritdoc IDataFeedsRouter
  function getRegistry(bytes32 feedId) external view returns (address registryAddress) {
    return s_feedIdToRegistry[feedId];
  }

  /**
   * @notice indicates that the Fee Manager address has been updated for a feed
   * @param currentFeeManager current Fee Manager address
   * @param previousFeeManager previous Fee Manager address
   */
  event FeeManagerSet(address indexed currentFeeManager, address indexed previousFeeManager);

  /// @inheritdoc IDataFeedsRouter
  function setFeeManager(address feeManagerAddress) public onlyOwner {
    address previousFeeManager = address(s_feeManager);
    s_feeManager = IDataFeedsFeeManager(feeManagerAddress);

    if (previousFeeManager != address(0)) s_authorizedNonbillableUsers.remove(previousFeeManager);
    s_authorizedNonbillableUsers.add(feeManagerAddress);
    emit FeeManagerSet(feeManagerAddress, previousFeeManager);
  }

  /// @inheritdoc IDataFeedsRouter
  function getFeeManager() external view returns (address feeManagerAddress) {
    return address(s_feeManager);
  }

  /// @inheritdoc IDataFeedsRouter
  function addNonbillableUser(address user) public onlyOwner {
    s_authorizedNonbillableUsers.add(user);

    emit NonbillableUserAdded(user);
  }

  /// @inheritdoc IDataFeedsRouter
  function removeNonbillableUser(address user) public onlyOwner {
    s_authorizedNonbillableUsers.remove(user);

    emit NonbillableUserRemoved(user);
  }
}
