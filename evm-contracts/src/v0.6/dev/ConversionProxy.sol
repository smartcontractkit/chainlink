pragma solidity 0.6.2;

import "./SignedSafeMath.sol";
import "../Owned.sol";
import "../interfaces/AggregatorInterface.sol";

/**
 * @title The ConversionProxy contract for Solidity v0.6
 * @notice This contract allows for the rate of one aggregator
 * contract to be represented in the currency of another aggregator
 * contract's current rate. Rounds and timestamps are referred to
 * relative to the _from address. Historic answers are provided at
 * the latest rate of _to address.
 */
contract ConversionProxy is AggregatorInterface, Owned {
  using SignedSafeMath for int256;

  AggregatorInterface public from;
  AggregatorInterface public to;

  event AddressesUpdated(
    address from,
    address to
  );

  /**
   * @notice Deploys the ConversionProxy contract
   * @param _from The address of the aggregator contract which
   * needs to be converted
   * @param _to The address of the aggregator contract which stores
   * the rate to convert to
   */
  constructor(
    address _from,
    address _to
  ) public Owned() {
    setAddresses(
      _from,
      _to
    );
  }

  /**
   * @dev Only callable by the owner of the contract
   * @param _from The address of the aggregator contract which
   * needs to be converted
   * @param _to The address of the aggregator contract which stores
   * the rate to convert to
   */
  function setAddresses(
    address _from,
    address _to
  ) public onlyOwner() {
    require(_from != _to, "Cannot use same address");
    from = AggregatorInterface(_from);
    to = AggregatorInterface(_to);
    emit AddressesUpdated(
      _from,
      _to
    );
  }

  /**
   * @notice Converts the latest answer of the `from` aggregator
   * to the rate of the `to` aggregator
   * @return The converted answer with amount of precision as defined
   * by `decimals` of the `to` aggregator
   */
  function latestAnswer()
    external
    virtual
    override
    returns (int256)
  {
    return _latestAnswer();
  }

  /**
   * @notice Calls the `latestTimestamp()` function of the `from`
   * aggregator
   * @return The value of latestTimestamp for the `from` aggregator
   */
  function latestTimestamp()
    external
    virtual
    override
    returns (uint256)
  {
    return _latestTimestamp();
  }

  /**
   * @notice Calls the `latestRound()` function of the `from`
   * aggregator
   * @return The value of latestRound for the `from` aggregator
   */
  function latestRound()
    external
    virtual
    override
    returns (uint256)
  {
    return _latestRound();
  }

  /**
   * @notice Converts the specified answer for `_roundId` of the
   * `from` aggregator to the latestAnswer of the `to` aggregator
   * @return The converted answer for `_roundId` of the `from`
   * aggregator with the amount of precision as defined by `decimals`
   * of the `to` aggregator
   */
  function getAnswer(uint256 _roundId)
    external
    virtual
    override
    returns (int256)
  {
    return _getAnswer(_roundId);
  }

  /**
   * @notice Calls the `getTimestamp(_roundId)` function of the `from`
   * aggregator for the specified `_roundId`
   * @return The timestamp of the `from` aggregator for the specified
   * `_roundId`
   */
  function getTimestamp(uint256 _roundId)
    external
    virtual
    override
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
   * @notice Calls the `decimals()` function of the `to` aggregator
   * @return The amount of precision the converted answer will contain
   */
  function decimals()
    external
    override
    returns (uint8)
  {
    return to.decimals();
  }


  function _latestAnswer()
    internal
    returns (int256)
  {
    return convertAnswer(from.latestAnswer());
  }

  function _latestTimestamp()
    internal
    returns (uint256)
  {
    return from.latestTimestamp();
  }

  function _latestRound()
    internal
    returns (uint256)
  {
    return from.latestRound();
  }

  function _getAnswer(uint256 _roundId)
    internal
    returns (int256)
  {
    return convertAnswer(from.getAnswer(_roundId));
  }

  function _getTimestamp(uint256 _roundId)
    internal
    returns (uint256)
  {
    return from.getTimestamp(_roundId);
  }

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
    uint256 roundIdFrom;
    int256 answerFrom;
    uint256 startedAtFrom;
    uint256 updatedAtFrom;
    uint256 answeredInRoundFrom;

    (roundIdFrom, answerFrom, startedAtFrom, updatedAtFrom, answeredInRoundFrom) = from.getRoundData(_roundId);
    return (roundIdFrom, convertAnswer(answerFrom), startedAtFrom, updatedAtFrom, answeredInRoundFrom);
  }

  function _latestRoundData()
    internal
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    )
  {
    uint256 roundIdFrom;
    int256 answerFrom;
    uint256 startedAtFrom;
    uint256 updatedAtFrom;
    uint256 answeredInRoundFrom;

    (roundIdFrom, answerFrom, startedAtFrom, updatedAtFrom, answeredInRoundFrom) = from.getRoundData(_latestRound());
    return (roundIdFrom, convertAnswer(answerFrom), startedAtFrom, updatedAtFrom, answeredInRoundFrom);
  }

  /**
   * @notice Converts the answer of the `from` aggregator to the rate
   * of the `to` aggregator at the precision of `decimals` of the `to`
   * aggregator
   * @param _answerFrom The answer of the `from` aggregator
   * @return The converted answer
   */
  function convertAnswer(int256 _answerFrom)
    internal
    returns (int256)
  {
    return _answerFrom.mul(to.latestAnswer()).div(int256(10 ** uint256(to.decimals())));
  }
}
