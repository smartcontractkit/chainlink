pragma solidity 0.4.24;

import "./interfaces/AggregatorInterface.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

/**
 * @title A trusted proxy for updating where current answers are read from
 * @notice This contract provides a consistent address for the
 * CurrentAnwerInterface but delegates where it reads from to the owner, who is
 * trusted to update it.
 */
contract AggregatorProxy is AggregatorInterface, Ownable {

  AggregatorInterface public aggregator;

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
   * @notice Reads the last updated height from aggregator delegated to.
   */
  function updatedHeight()
    external
    returns (uint256)
  {
    return aggregator.updatedHeight();
  }

  /**
   * @notice Allows the owner to update the aggregator address.
   * @param _aggregator The new address for the aggregator contract
   */
  function setAggregator(address _aggregator)
    public
    onlyOwner()
  {
    aggregator = AggregatorInterface(_aggregator);
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
