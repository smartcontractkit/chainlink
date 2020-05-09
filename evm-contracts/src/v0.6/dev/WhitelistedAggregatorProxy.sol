pragma solidity 0.6.2;

import "./AggregatorProxy.sol";
import "./Whitelisted.sol";

/**
 * @title A trusted proxy for updating where current answers are read from
 * @notice This contract provides a consistent address for the
 * AggregatorInterface but delegates where it reads from to the owner, who is
 * trusted to update it.
 * @notice Only whitelisted addresses are allowed to access getters for
 * aggregated answers and round information.
 */
contract WhitelistedAggregatorProxy is AggregatorProxy, Whitelisted {

  constructor(address _aggregator) public AggregatorProxy(_aggregator) {
  }

  /**
   * @notice Reads the current answer from aggregator delegated to.
   * @dev overridden function to add the isWhitelisted() modifier
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
   * @dev overridden function to add the isWhitelisted() modifier
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
   * @dev overridden function to add the isWhitelisted() modifier
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
   * @dev overridden function to add the isWhitelisted() modifier
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
   * @dev overridden function to add the isWhitelisted() modifier
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

  /**
   * @notice get all details about a round. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt and
   * answeredInRound return values.
   * @param _roundId the round ID to retrieve the details for. If _roundId
   * has the special value UINT256_MAX (2**256-1), the contract will retrieve
   * the latest round's details.
   * @return roundId is the round ID for which details were retrieved
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
    view
    isWhitelisted()
    override
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
   * @notice get all details about the latest round.
   * @return roundId is the round ID for which details were retrieved
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
    view
    isWhitelisted()
    override
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
