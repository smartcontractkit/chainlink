pragma solidity 0.5.0;

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
    isWhitelisted()
    view
    returns (int256)
  {
    return rounds[latestRoundId].answer;
  }

  /**
   * @notice get the most recent updated at timestamp
   * @dev overridden funcion to add the isWhitelisted() modifier
   */
  function latestTimestamp()
    external
    view
    isWhitelisted()
    returns (uint256)
  {
    return rounds[latestRoundId].updatedAt;
  }

  /**
   * @notice get past rounds answers
   * @dev overridden funcion to add the isWhitelisted() modifier
   * @param _roundId the round number to retrieve the answer for
   */
  function getAnswer(uint256 _roundId)
    external
    view
    isWhitelisted()
    returns (int256)
  {
    return rounds[uint32(_roundId)].answer;
  }

  /**
   * @notice get timestamp when an answer was last updated
   * @dev overridden funcion to add the isWhitelisted() modifier
   * @param _roundId the round number to retrieve the updated timestamp for
   */
  function getTimestamp(uint256 _roundId)
    external
    isWhitelisted()
    view
    returns (uint256)
  {
    return rounds[uint32(_roundId)].updatedAt;
  }
}
