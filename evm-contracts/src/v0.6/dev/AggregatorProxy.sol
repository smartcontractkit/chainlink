pragma solidity 0.6.2;

import "./AggregatorInterface.sol";
import "../Owned.sol";

/**
 * @title A trusted proxy for updating where current answers are read from
 * @notice This contract provides a consistent address for the
 * CurrentAnwerInterface but delegates where it reads from to the owner, who is
 * trusted to update it.
 */
contract AggregatorProxy is AggregatorInterface, Owned {

  AggregatorInterface public aggregator;

  constructor(address _aggregator) public Owned() {
    setAggregator(_aggregator);
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
    return _latestAnswer();
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
    return _latestTimestamp();
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
    return _getAnswer(_roundId);
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
    return _getTimestamp(_roundId);
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
    return _latestRound();
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
    return _latestRoundData();
  }

  /**
   * @notice represents the number of decimals the aggregator responses represent.
   */
  function decimals()
    external
    view
    override
    returns (uint8)
  {
    return aggregator.decimals();
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

  /*
   * Internal
   */

  function _latestAnswer()
    internal
    view
    returns (int256)
  {
    return aggregator.latestAnswer();
  }

  function _latestTimestamp()
    internal
    view
    returns (uint256)
  {
    return aggregator.latestTimestamp();
  }

  function _getAnswer(uint256 _roundId)
    internal
    view
    returns (int256)
  {
    return aggregator.getAnswer(_roundId);
  }

  function _getTimestamp(uint256 _roundId)
    internal
    view
    returns (uint256)
  {
    return aggregator.getTimestamp(_roundId);
  }

  function _latestRound()
    internal
    view
    returns (uint256)
  {
    return aggregator.latestRound();
  }

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
    return aggregator.getRoundData(_roundId);
  }

  function _latestRoundData()
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
    return aggregator.latestRoundData();
  }
}
