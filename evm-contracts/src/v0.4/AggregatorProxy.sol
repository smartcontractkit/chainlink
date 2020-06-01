pragma solidity 0.4.24;

import "./interfaces/HistoricAggregatorInterface.sol";
import "./vendor/Ownable.sol";

/**
 * @title A trusted proxy for updating where current answers are read from
 * @notice This contract provides a consistent address for the
 * CurrentAnwerInterface but delegates where it reads from to the owner, who is
 * trusted to update it.
 */
contract AggregatorProxy is HistoricAggregatorInterface, Ownable {

  HistoricAggregatorInterface public aggregator;

  constructor(address _aggregator) public Ownable() {
    setAggregator(_aggregator);
  }

  /**
   * @notice Reads the current answer from aggregator delegated to.
   */
  function latestAnswer()
    external
    returns (int256)
  {
    return aggregator.latestAnswer();
  }

  /**
   * @notice Reads the last updated height from aggregator delegated to.
   */
  function latestTimestamp()
    external
    returns (uint256)
  {
    return aggregator.latestTimestamp();
  }

  /**
   * @notice get past rounds answers
   * @param _roundId the answer number to retrieve the answer for
   */
  function getAnswer(uint256 _roundId)
    external
    returns (int256)
  {
    return aggregator.getAnswer(_roundId);
  }

  /**
   * @notice get block timestamp when an answer was last updated
   * @param _roundId the answer number to retrieve the updated timestamp for
   */
  function getTimestamp(uint256 _roundId)
    external
    returns (uint256)
  {
    return aggregator.getTimestamp(_roundId);
  }

  /**
   * @notice get the latest completed round where the answer was updated
   */
  function latestRound()
    external
    returns (uint256)
  {
    return aggregator.latestRound();
  }

  /**
   * @notice Allows the owner to update the aggregator address.
   * @param _aggregator The new address for the aggregator contract
   */
  function setAggregator(address _aggregator)
    public
    onlyOwner()
  {
    aggregator = HistoricAggregatorInterface(_aggregator);
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
