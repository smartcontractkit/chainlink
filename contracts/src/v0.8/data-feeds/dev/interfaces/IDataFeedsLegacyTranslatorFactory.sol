// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @notice IDataFeedsLegacyTranslatorFactory
/// Produces DataFeedsLegacyTranslator clone contracts that
/// call into DataFeedsLegacyTranslator with their mapped feedId.
/// Legacy AggregatorProxy contracts can point to the respective
/// DataFeedsLegacyTranslator to retrieve benchmark values using
/// the legacy interface.

interface IDataFeedsLegacyTranslatorFactory {
  function createTranslators(
    bytes32[] memory feedIds,
    string[] memory descriptions,
    uint256 decimals,
    uint256 version
  ) external;

  function destructTranslators(address[] memory translatorAddresses) external;

  function getTranslators(bytes32[] memory feedIds) external view returns (address[] memory translators);
}
