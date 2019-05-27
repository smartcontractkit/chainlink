pragma solidity 0.4.24;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "../interfaces/CurrentAnswerInterface.sol";

contract AggregatorProxy is Ownable, CurrentAnswerInterface {

  CurrentAnswerInterface public aggregator;

  constructor(address _aggregator)
    public
    Ownable()
  {
    updateAggregator(_aggregator);
  }

  function updateAggregator(address _aggregator)
    public
    onlyOwner()
  {
    aggregator = CurrentAnswerInterface(_aggregator);
  }

  function currentAnswer()
    public
    returns (uint256)
  {
    return aggregator.currentAnswer();
  }

}
