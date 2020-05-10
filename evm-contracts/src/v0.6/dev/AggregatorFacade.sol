pragma solidity 0.6.2;

import "./AggregatorInterface.sol";
import "./HistoricAggregatorInterface.sol";

/**
 * @title A facade for Historic Aggregator versions to conform to the new v0.6
 * Aggregator interface.
 */
contract AggregatorFacade is AggregatorInterface {

  HistoricAggregatorInterface public aggregator;
  uint8 public override decimals;

  constructor(address _aggregator, uint8 _decimals) public {
    aggregator = HistoricAggregatorInterface(_aggregator);
    decimals = _decimals;
  }

  /**
   * @notice get the latest completed round where the answer was updated
   */
  function latestRound()
    external
    view
    virtual
    override
    returns (uint256)
  {
    return aggregator.latestRound();
  }

  /**
   * @notice Reads the current answer from aggregator delegated to.
   */
  function latestAnswer()
    external
    view
    virtual
    override
    returns (int256)
  {
    return aggregator.latestAnswer();
  }

  /**
   * @notice Reads the last updated height from aggregator delegated to.
   */
  function latestTimestamp()
    external
    view
    virtual
    override
    returns (uint256)
  {
    return aggregator.latestTimestamp();
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
    virtual
    override
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    )
  {
    return _getRoundData(aggregator.latestRound());
  }

  /**
   * @notice get past rounds answers
   * @param _roundId the answer number to retrieve the answer for
   */
  function getAnswer(uint256 _roundId)
    external
    view
    virtual
    override
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
    view
    virtual
    override
    returns (uint256)
  {
    return aggregator.getTimestamp(_roundId);
  }

  /**
   * @notice get data about a round. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt and
   * answeredInRound return values.
   * @param _roundId the round ID to retrieve the round data for
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
    virtual
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

  /*
   * Internal
   */

  function _getRoundData(uint256 _roundId)
    internal
    view
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    )
  {
    answer = aggregator.getAnswer(_roundId);
    updatedAt = uint64(aggregator.getTimestamp(_roundId));
    return (_roundId, answer, updatedAt, updatedAt, _roundId);
  }

}
