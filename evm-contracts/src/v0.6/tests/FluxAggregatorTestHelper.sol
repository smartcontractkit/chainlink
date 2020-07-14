pragma solidity ^0.6.0;

import "../dev/FluxAggregator.sol";

contract FluxAggregatorTestHelper {

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

  function readOracleRoundState(address _aggregator, address _oracle)
    external
  {
    FluxAggregator(_aggregator).oracleRoundState(_oracle, 0);
    emit Here();
  }

}
