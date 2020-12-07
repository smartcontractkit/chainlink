// SPDX-License-Identifier: MIT
pragma solidity 0.6.6;

import "./FluxAggregator.sol";
import "./SimpleReadAccessController.sol";

/**
 * @title AccessControlled FluxAggregator contract
 * @notice This contract requires addresses to be added to a controller
 * in order to read the answers stored in the FluxAggregator contract
 */
contract AccessControlledAggregator is FluxAggregator, SimpleReadAccessController {

  /**
   * @notice set up the aggregator with initial configuration
   * @param _link The address of the LINK token
   * @param _paymentAmount The amount paid of LINK paid to each oracle per submission, in wei (units of 10⁻¹⁸ LINK)
   * @param _timeout is the number of seconds after the previous round that are
   * allowed to lapse before allowing an oracle to skip an unfinished round
   * @param _validator is an optional contract address for validating
   * external validation of answers
   * @param _minSubmissionValue is an immutable check for a lower bound of what
   * submission values are accepted from an oracle
   * @param _maxSubmissionValue is an immutable check for an upper bound of what
   * submission values are accepted from an oracle
   * @param _decimals represents the number of decimals to offset the answer by
   * @param _description a short description of what is being reported
   */
  constructor(
    address _link,
    uint128 _paymentAmount,
    uint32 _timeout,
    address _validator,
    int256 _minSubmissionValue,
    int256 _maxSubmissionValue,
    uint8 _decimals,
    string memory _description
  ) public FluxAggregator(
    _link,
    _paymentAmount,
    _timeout,
    _validator,
    _minSubmissionValue,
    _maxSubmissionValue,
    _decimals,
    _description
  ){}

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
   * @dev overridden funcion to add the checkAccess() modifier
   * @dev Note that for in-progress rounds (i.e. rounds that haven't yet
   * received maxSubmissions) answer and updatedAt may change between queries.
   */
  function getRoundData(uint80 _roundId)
    public
    view
    override
    checkAccess()
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    return super.getRoundData(_roundId);
  }

  /**
   * @notice get data about the latest round. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt and
   * answeredInRound return values. Consumers are encouraged to
   * use this more fully featured method over the "legacy" latestAnswer
   * functions. Consumers are encouraged to check that they're receiving fresh
   * data by inspecting the updatedAt and answeredInRound return values.
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
   * @dev overridden funcion to add the checkAccess() modifier
   * @dev Note that for in-progress rounds (i.e. rounds that haven't yet
   * received maxSubmissions) answer and updatedAt may change between queries.
   */
  function latestRoundData()
    public
    view
    override
    checkAccess()
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    return super.latestRoundData();
  }

  /**
   * @notice get the most recently reported answer
   * @dev overridden funcion to add the checkAccess() modifier
   *
   * @dev #[deprecated] Use latestRoundData instead. This does not error if no
   * answer has been reached, it will simply return 0. Either wait to point to
   * an already answered Aggregator or use the recommended latestRoundData
   * instead which includes better verification information.
   */
  function latestAnswer()
    public
    view
    override
    checkAccess()
    returns (int256)
  {
    return super.latestAnswer();
  }

  /**
   * @notice get the most recently reported round ID
   * @dev overridden funcion to add the checkAccess() modifier
   *
   * @dev #[deprecated] Use latestRoundData instead. This does not error if no
   * answer has been reached, it will simply return 0. Either wait to point to
   * an already answered Aggregator or use the recommended latestRoundData
   * instead which includes better verification information.
   */
  function latestRound()
    public
    view
    override
    checkAccess()
    returns (uint256)
  {
    return super.latestRound();
  }

  /**
   * @notice get the most recent updated at timestamp
   * @dev overridden funcion to add the checkAccess() modifier
   *
   * @dev #[deprecated] Use latestRoundData instead. This does not error if no
   * answer has been reached, it will simply return 0. Either wait to point to
   * an already answered Aggregator or use the recommended latestRoundData
   * instead which includes better verification information.
   */
  function latestTimestamp()
    public
    view
    override
    checkAccess()
    returns (uint256)
  {
    return super.latestTimestamp();
  }

  /**
   * @notice get past rounds answers
   * @dev overridden funcion to add the checkAccess() modifier
   * @param _roundId the round number to retrieve the answer for
   *
   * @dev #[deprecated] Use getRoundData instead. This does not error if no
   * answer has been reached, it will simply return 0. Either wait to point to
   * an already answered Aggregator or use the recommended getRoundData
   * instead which includes better verification information.
   */
  function getAnswer(uint256 _roundId)
    public
    view
    override
    checkAccess()
    returns (int256)
  {
    return super.getAnswer(_roundId);
  }

  /**
   * @notice get timestamp when an answer was last updated
   * @dev overridden funcion to add the checkAccess() modifier
   * @param _roundId the round number to retrieve the updated timestamp for
   *
   * @dev #[deprecated] Use getRoundData instead. This does not error if no
   * answer has been reached, it will simply return 0. Either wait to point to
   * an already answered Aggregator or use the recommended getRoundData
   * instead which includes better verification information.
   */
  function getTimestamp(uint256 _roundId)
    public
    view
    override
    checkAccess()
    returns (uint256)
  {
    return super.getTimestamp(_roundId);
  }

}
