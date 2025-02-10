// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {IBundleBaseAggregator} from "./IBundleBaseAggregator.sol";
import {ICommonAggregator} from "./ICommonAggregator.sol";
import {IDecimalAggregator} from "./IDecimalAggregator.sol";

/// @notice IDataFeedsCache
/// Responsible for storing data associated with a given data ID and additional request data.
interface IDataFeedsCache is IDecimalAggregator, IBundleBaseAggregator, ICommonAggregator {
  /// @notice Remove feed configs.
  /// @param dataIds List of data IDs
  function removeFeedConfigs(
    bytes16[] calldata dataIds
  ) external;

  /// @notice Update mappings for AggregatorProxy -> Data ID
  /// @param proxies AggregatorProxy addresses
  /// @param dataIds Data IDs
  function updateDataIdMappingsForProxies(address[] calldata proxies, bytes16[] calldata dataIds) external;

  /// @notice Remove mappings for AggregatorProxy -> Data IDs
  /// @param proxies  AggregatorProxy addresses to remove
  function removeDataIdMappingsForProxies(
    address[] calldata proxies
  ) external;

  /// @notice Get the Data ID mapping for a AggregatorProxy
  /// @param proxy AggregatorProxy addresses which will be reading feed data
  function getDataIdForProxy(
    address proxy
  ) external view returns (bytes16 dataId);
}
