pragma solidity ^0.6.0;

import "./HistoricAggregatorInterface.sol";

interface AggregatorInterface is HistoricAggregatorInterface {
  function decimals() external view returns (uint8);
  function getRoundData(uint256 _roundId)
    external
    view
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    );
  function latestRoundData()
    external
    view
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    );
}
