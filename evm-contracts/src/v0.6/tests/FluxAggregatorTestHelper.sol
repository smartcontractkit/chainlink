pragma solidity ^0.6.0;

import "../dev/FluxAggregator.sol";

contract FluxAggregatorTestHelper is Owned {

  event Here();

  function readOracleRoundState(address _aggregator, address _oracle)
    external
  {
    FluxAggregator(_aggregator).oracleRoundState(_oracle, 0);
    emit Here();
  }

  function readLatestAnswer(address _aggregator)
    external
  {
    FluxAggregator(_aggregator).latestAnswer();
    emit Here();
  }

  function readLatestTimestamp(address _aggregator)
    external
  {
    FluxAggregator(_aggregator).latestTimestamp();
    emit Here();
  }

  function readGetAnswer(address _aggregator, uint256 _roundID)
    external
  {
    FluxAggregator(_aggregator).getAnswer(_roundID);
    emit Here();
  }

  function readGetTimestamp(address _aggregator, uint256 _roundID)
    external
  {
    FluxAggregator(_aggregator).getTimestamp(_roundID);
    emit Here();
  }

}
