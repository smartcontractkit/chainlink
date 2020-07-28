pragma solidity ^0.6.0;

import "../dev/FluxAggregator.sol";

contract FluxAggregatorTestHelper {

  uint80 public requestedRoundId;

  event Here();

  function readLatestRoundData(address _aggregator)
    external
  {
    FluxAggregator(_aggregator).latestRoundData();
    emit Here();
  }

  function readGetRoundData(address _aggregator, uint80 _roundID)
    external
  {
    FluxAggregator(_aggregator).getRoundData(_roundID);
    emit Here();
  }

  function readLatestAnswer(address _aggregator)
    external
  {
    FluxAggregator(_aggregator).latestAnswer();
    emit Here();
  }

  function readOracleRoundState(address _aggregator, address _oracle)
    external
  {
    FluxAggregator(_aggregator).oracleRoundState(_oracle, 0);
    emit Here();
  }

  function requestNewRound(address _aggregator)
    external
  {
    requestedRoundId = FluxAggregator(_aggregator).requestNewRound();
    emit Here();
  }

}
