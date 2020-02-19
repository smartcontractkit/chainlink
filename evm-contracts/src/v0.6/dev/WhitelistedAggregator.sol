pragma solidity 0.6.2;

import "./PrepaidAggregator.sol";
import "./Whitelisted.sol";

/**
 * @title Whitelisted Prepaid Aggregator contract
 * @notice This contract requires addresses to be added to a whitelist
 * in order to read the answers stored in the PrepaidAggregator contract
 */
contract WhitelistedAggregator is PrepaidAggregator, Whitelisted {

  constructor(
    address _link,
    uint128 _paymentAmount,
    uint32 _timeout,
    uint8 _decimals,
    bytes32 _description
  ) public PrepaidAggregator(
    _link,
    _paymentAmount,
    _timeout,
    _decimals,
    _description
  ){}

  /**
   * @notice get the most recently reported answer
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
   * @notice get the most recent updated at timestamp
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
   * @dev overridden funcion to add the isWhitelisted() modifier
   * @param _roundId the round number to retrieve the answer for
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
   * @notice get timestamp when an answer was last updated
   * @dev overridden funcion to add the isWhitelisted() modifier
   * @param _roundId the round number to retrieve the updated timestamp for
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
}
