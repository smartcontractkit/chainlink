pragma solidity ^0.6.0;

import "../dev/FluxAggregator.sol";

contract FluxAggregatorTestHelper {

  uint80 public requestedRoundId;

  function readLatestRoundData(address _aggregator)
    external
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    return FluxAggregator(_aggregator).latestRoundData();
  }

  function readGetRoundData(address _aggregator, uint80 _roundID)
    external
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    return FluxAggregator(_aggregator).getRoundData(_roundID);
  }

  function readLatestAnswer(address _aggregator)
    external
    returns(
      int256 answer
    )
  {
    return FluxAggregator(_aggregator).latestAnswer();
  }

  function readOracleRoundState(address _aggregator, address _oracle)
    external
  {
    FluxAggregator(_aggregator).oracleRoundState(_oracle, 0);
  }

  function requestNewRound(address _aggregator)
    external
  {
    requestedRoundId = FluxAggregator(_aggregator).requestNewRound();
  }

}
