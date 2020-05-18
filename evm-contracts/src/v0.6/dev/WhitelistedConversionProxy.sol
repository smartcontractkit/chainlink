pragma solidity 0.6.2;

import "./ConversionProxy.sol";
import "./Whitelisted.sol";

/**
 * @title A trusted proxy for updating where current answers are read from
 * @notice This contract provides a consistent address for the
 * AggregatorInterface but delegates where it reads from to the owner, who is
 * trusted to update it.
 * @notice Only whitelisted addresses are allowed to access getters for
 * aggregated answers and round information.
 */
contract WhitelistedConversionProxy is ConversionProxy, Whitelisted {

  /**
   * @notice Deploys the WhitelistedConversionProxy contract
   * @param _from The address of the aggregator contract which
   * needs to be converted
   * @param _to The address of the aggregator contract which stores
   * the rate to convert to
   */
  constructor(
    address _from,
    address _to
  ) public ConversionProxy(
    _from,
    _to
  ) {}

  /**
   * @notice Converts the latest answer of the `from` aggregator
   * to the rate of the `to` aggregator
   * @dev Overridden function to add the `isWhitelisted()` modifier
   * @return The converted answer with amount of precision as defined
   * by `decimals` of the `to` aggregator
   */
  function latestAnswer()
    external
    override
    isWhitelisted()
    returns (int256)
  {
    return _latestAnswer();
  }

  /**
   * @notice Calls the `latestTimestamp()` function of the `from`
   * aggregator
   * @dev Overridden function to add the `isWhitelisted()` modifier
   * @return The value of latestTimestamp for the `from` aggregator
   */
  function latestTimestamp()
    external
    override
    isWhitelisted()
    returns (uint256)
  {
    return _latestTimestamp();
  }

  /**
   * @notice Calls the `latestRound()` function of the `from`
   * aggregator
   * @dev Overridden function to add the `isWhitelisted()` modifier
   * @return The value of latestRound for the `from` aggregator
   */
  function latestRound()
    external
    override
    isWhitelisted()
    returns (uint256)
  {
    return _latestRound();
  }

  /**
   * @notice Converts the specified answer for `_roundId` of the
   * `from` aggregator to the latestAnswer of the `to` aggregator
   * @dev Overridden function to add the `isWhitelisted()` modifier
   * @return The converted answer for `_roundId` of the `from`
   * aggregator with the amount of precision as defined by `decimals`
   * of the `to` aggregator
   */
  function getAnswer(uint256 _roundId)
    external
    override
    isWhitelisted()
    returns (int256)
  {
    return _getAnswer(_roundId);
  }

  /**
   * @notice Calls the `getTimestamp(_roundId)` function of the `from`
   * aggregator for the specified `_roundId`
   * @dev Overridden function to add the `isWhitelisted()` modifier
   * @return The timestamp of the `from` aggregator for the specified
   * `_roundId`
   */
  function getTimestamp(uint256 _roundId)
    external
    override
    isWhitelisted()
    returns (uint256)
  {
    return _getTimestamp(_roundId);
  }

  /**
   * @notice get data about a round. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt and
   * answeredInRound return values.
   * @param _roundId the round ID to retrieve the round data for
   * @return roundId is the round ID for which data was retrieved
   * @return answer is the answer for the given round
   * @return startedAt is the timestamp when the round was started. This is 0
   * if the round hasn't been started yet.
   * @return updatedAt is the timestamp when the round last was updated (i.e.
   * answer was last computed)
   * @return answeredInRound is the round ID of the round in which the answer
   * was computed. answeredInRound may be smaller than roundId when the round
   * timed out. answerInRound is equal to roundId when the round didn't time out
   * and was completed regularly.
   * @dev Note that for in-progress rounds (i.e. rounds that haven't yet received
   * maxSubmissions) answer and updatedAt may change between queries.
   */
  function getRoundData(uint256 _roundId)
    external
    override
    isWhitelisted()
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    )
  {
    return _getRoundData(_roundId);
  }

  /**
   * @notice get data about the latest round. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt and
   * answeredInRound return values.
   * @return roundId is the round ID for which data was retrieved
   * @return answer is the answer for the given round
   * @return startedAt is the timestamp when the round was started. This is 0
   * if the round hasn't been started yet.
   * @return updatedAt is the timestamp when the round last was updated (i.e.
   * answer was last computed)
   * @return answeredInRound is the round ID of the round in which the answer
   * was computed. answeredInRound may be smaller than roundId when the round
   * timed out. answerInRound is equal to roundId when the round didn't time out
   * and was completed regularly.
   * @dev Note that for in-progress rounds (i.e. rounds that haven't yet received
   * maxSubmissions) answer and updatedAt may change between queries.
   */
  function latestRoundData()
    external
    override
    isWhitelisted()
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    )
  {
    return _latestRoundData();
  }
}
