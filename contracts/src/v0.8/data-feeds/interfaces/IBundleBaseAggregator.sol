// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

interface IBundleBaseAggregator {
  function latestBundle() external view returns (bytes memory bundle);

  function bundleDecimals() external view returns (uint8[] memory);

  function latestBundleTimestamp() external view returns (uint256);
}
