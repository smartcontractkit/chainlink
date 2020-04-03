pragma solidity 0.6.2;

import "../dev/SignedSafeMath.sol";
import "./AggregatorInterface.sol";
import "../Owned.sol";

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
  function latestAnswer() external view virtual override returns (int256) {
    return _latestAnswer();
  }

  /**
   * @notice Calls the `latestTimestamp()` function of the `from`
   * aggregator
   * @return The value of latestTimestamp for the `from` aggregator
   */
  function latestTimestamp() external view virtual override returns (uint256) {
    return _latestTimestamp();
  }

  /**
   * @notice Calls the `latestRound()` function of the `from`
   * aggregator
   * @return The value of latestRound for the `from` aggregator
   */
  function latestRound() external view virtual override returns (uint256) {
    return _latestRound();
  }

  /**
   * @notice Converts the specified answer for `_roundId` of the
   * `from` aggregator to the latestAnswer of the `to` aggregator
   * @return The converted answer for `_roundId` of the `from`
   * aggregator with the amount of precision as defined by `decimals`
   * of the `to` aggregator
   */
  function getAnswer(uint256 _roundId) external view virtual override returns (int256) {
    return _getAnswer(_roundId);
  }

  /**
   * @notice Calls the `getTimestamp(_roundId)` function of the `from`
   * aggregator for the specified `_roundId`
   * @return The timestamp of the `from` aggregator for the specified
   * `_roundId`
   */
  function getTimestamp(uint256 _roundId) external view virtual override returns (uint256) {
    return _getTimestamp(_roundId);
  }

  /**
   * @notice Calls the `decimals()` function of the `to` aggregator
   * @return The amount of precision the converted answer will contain
   */
  function decimals() external view override returns (uint8) {
    return to.decimals();
  }

  function _latestAnswer() internal view returns (int256) {
    return convertAnswer(from.latestAnswer(), to.latestAnswer());
  }

  function _latestTimestamp() internal view returns (uint256) {
    return from.latestTimestamp();
  }

  function _latestRound() internal view returns (uint256) {
    return from.latestRound();
  }

  function _getAnswer(uint256 _roundId) internal view returns (int256) {
    return convertAnswer(from.getAnswer(_roundId), to.latestAnswer());
  }

  function _getTimestamp(uint256 _roundId) internal view returns (uint256) {
    return from.getTimestamp(_roundId);
  }

  /**
   * @notice Converts the answer of the `from` aggregator to the rate
   * of the `to` aggregator at the precision of `decimals` of the `to`
   * aggregator
   * @return The converted answer
   */
  function convertAnswer(
    int256 _answerFrom,
    int256 _answerTo
  ) internal view returns (int256) {
    return _answerFrom.mul(_answerTo).div(int256(10 ** uint256(to.decimals())));
  }
}
