// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "../FluxAggregator.sol";

contract FluxAggregatorTestHelper {

  uint80 public requestedRoundId;

  function readOracleRoundState(address _aggregator, address _oracle)
    external
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    FluxAggregator(_aggregator).oracleRoundState(_oracle, 0);
  }

  function readGetRoundData(address _aggregator, uint80 _roundID)
    external
  {
    FluxAggregator(_aggregator).getRoundData(_roundID);
  }

  function readLatestRoundData(address _aggregator)
    external
  {
    FluxAggregator(_aggregator).latestRoundData();
  }

  function readLatestAnswer(address _aggregator)
    external
  {
    FluxAggregator(_aggregator).latestAnswer();
  }

  function readLatestTimestamp(address _aggregator)
    external
  {
    FluxAggregator(_aggregator).latestTimestamp();
  }

  function readLatestRound(address _aggregator)
    external
  {
    FluxAggregator(_aggregator).latestRound();
  }

  function requestNewRound(address _aggregator)
    external
  {
    requestedRoundId = FluxAggregator(_aggregator).requestNewRound();
  }

  function readGetAnswer(address _aggregator, uint256 _roundID)
    external
  {
    FluxAggregator(_aggregator).getAnswer(_roundID);
  }

  function readGetTimestamp(address _aggregator, uint256 _roundID)
    external
  {
    FluxAggregator(_aggregator).getTimestamp(_roundID);
  }

}
