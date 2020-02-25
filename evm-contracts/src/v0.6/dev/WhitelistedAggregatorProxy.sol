pragma solidity 0.6.2;

import "./AggregatorProxy.sol";
import "./Whitelisted.sol";

/**
 * @title A trusted proxy for updating where current answers are read from
 * @notice This contract provides a consistent address for the
 * CurrentAnwerInterface but delegates where it reads from to the owner, who is
 * trusted to update it.
 * @notice Only whitelisted addresses are allowed to access getters for
 * aggregated answers and round information.
 */
contract WhitelistedAggregatorProxy is AggregatorProxy, Whitelisted {

  constructor(address _aggregator) public AggregatorProxy(_aggregator) {
  }

  /**
   * @notice Reads the current answer from aggregator delegated to.
   * @dev overridden funcion to add the isWhitelisted() modifier
   */
  function latestAnswer()
    external
    view
    override
    isWhitelisted()
    returns (int256)
  {
    return _latestAnswer();
  }

  /**
   * @notice Reads the last updated height from aggregator delegated to.
   * @dev overridden funcion to add the isWhitelisted() modifier
   */
  function latestTimestamp()
    external
    view
    override
    isWhitelisted()
    returns (uint256)
  {
    return _latestTimestamp();
  }

  /**
   * @notice get past rounds answers
   * @param _roundId the answer number to retrieve the answer for
   * @dev overridden funcion to add the isWhitelisted() modifier
   */
  function getAnswer(uint256 _roundId)
    external
    view
    override
    isWhitelisted()
    returns (int256)
  {
    return _getAnswer(_roundId);
  }

  /**
   * @notice get block timestamp when an answer was last updated
   * @param _roundId the answer number to retrieve the updated timestamp for
   * @dev overridden funcion to add the isWhitelisted() modifier
   */
  function getTimestamp(uint256 _roundId)
    external
    view
    override
    isWhitelisted()
    returns (uint256)
  {
    return _getTimestamp(_roundId);
  }

  /**
   * @notice get the latest completed round where the answer was updated
   * @dev overridden funcion to add the isWhitelisted() modifier
   */
  function latestRound()
    external
    view
    override
    isWhitelisted()
    returns (uint256)
  {
    return _latestRound();
  }
}
