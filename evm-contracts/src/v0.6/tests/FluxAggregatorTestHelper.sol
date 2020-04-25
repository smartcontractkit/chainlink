pragma solidity ^0.6.0;

import "../dev/FluxAggregator.sol";

contract FluxAggregatorTestHelper is Owned {

  event Here();

  function readOracleRoundState(address _aggregator, address _oracle)
    external
  {
    FluxAggregator(_aggregator).oracleRoundState(_oracle);
    emit Here();
  }

}
