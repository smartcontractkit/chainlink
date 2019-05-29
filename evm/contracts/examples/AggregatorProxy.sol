pragma solidity 0.4.24;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "../interfaces/CurrentAnswerInterface.sol";

/**
 * @title A trusted proxy for updating where current answers are read from
 * @notice This contract provides a consistent address for the
 * CurrentAnwerInterface but delegates where it reads from to the owner, who is
 * trusted to update it.
 */
contract AggregatorProxy is Ownable, CurrentAnswerInterface {

  CurrentAnswerInterface public aggregator;

  constructor(address _aggregator)
    public
    Ownable()
  {
    setAggregator(_aggregator);
  }

  /**
   * @notice Reads the current answer from aggregator delegated to.
   */
  function currentAnswer()
    external
    returns (int256)
  {
    return aggregator.currentAnswer();
  }

  /**
   * @notice Allows the owner to update the aggregator address.
   * @param _aggregator The new address for the aggregator contract
   */
  function setAggregator(address _aggregator)
    public
    onlyOwner()
  {
    aggregator = CurrentAnswerInterface(_aggregator);
  }

  /**
   * @notice Allows the owner to destroy the contract if it is not intended to
   * be used any longer.
   */
  function destroy()
    external
    onlyOwner()
  {
    selfdestruct(owner);
  }

}
