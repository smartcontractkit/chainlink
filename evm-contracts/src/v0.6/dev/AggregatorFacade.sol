pragma solidity 0.6.2;

import "../interfaces/AggregatorInterface.sol";
import "../interfaces/HistoricAggregatorInterface.sol";

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
    virtual
    override
    returns (uint256)
  {
    return aggregator.latestTimestamp();
  }

  /**
   * @notice get data about the latest round. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt value.
   * @return roundId is the round ID for which data was retrieved
   * @return answer is the answer for the given round
   * @return startedAt is always equal to updatedAt because the underlying
   * Aggregator contract does not expose this information.
   * @return updatedAt is the timestamp when the round last was updated (i.e.
   * answer was last computed)
   * @return answeredInRound is always equal to roundId because the underlying
   * Aggregator contract does not expose this information.
   * @dev Note that for rounds that haven't yet received responses from all
   * oracles, answer and updatedAt may change between queries.
   */
  function latestRoundData()
    external
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
    virtual
    override
    returns (uint256)
  {
    return aggregator.getTimestamp(_roundId);
  }

  /**
   * @notice get data about a round. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt value.
   * @param _roundId the round ID to retrieve the round data for
   * @return roundId is the round ID for which data was retrieved
   * @return answer is the answer for the given round
   * @return startedAt is always equal to updatedAt because the underlying
   * Aggregator contract does not expose this information.
   * @return updatedAt is the timestamp when the round last was updated (i.e.
   * answer was last computed)
   * @return answeredInRound is always equal to roundId because the underlying
   * Aggregator contract does not expose this information.
   * @dev Note that for rounds that haven't yet received responses from all
   * oracles, answer and updatedAt may change between queries.
   */
  function getRoundData(uint256 _roundId)
    external
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
    if (updatedAt == 0) {
      answeredInRound = 0;
    } else {
      answeredInRound = _roundId;
    }
    return (_roundId, answer, updatedAt, updatedAt, answeredInRound);
  }

}
