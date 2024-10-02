// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {AggregatorV2V3Interface} from "../../../shared/interfaces/AggregatorV2V3Interface.sol";

/// @notice IDataFeedsLegacyTranslator
/// Reflects the same interface methods as legacy Data Feeds contracts,
/// but includes payable on relevant methods and forwards msg.sender
/// address in request to the DataFeedsRouter. This allows Data Feeds
/// to properly bill for requests without any changes needed on the
/// users' contracts besides setting a balance address and ERC-20 token
/// choice. More sophisticated logic is up to the user to implement.
interface IDataFeedsLegacyTranslator is AggregatorV2V3Interface {
  /// AggregatorInterface

  function latestAnswer() external view returns (int256);

  function latestTimestamp() external view returns (uint256);

  function latestRound() external view returns (uint256);

  function getAnswer(uint256 roundId) external view returns (int256);

  function getTimestamp(uint256 roundId) external view returns (uint256);

  /// AggregatorV3Interface

  function decimals() external view returns (uint8);

  function description() external view returns (string memory);

  function version() external view returns (uint256);

  function getRoundData(
    uint80 _roundId
  ) external view returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound);

  function latestRoundData()
    external
    view
    returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound);
}
