pragma solidity 0.4.24;

import "./vendor/SignedSafeMath.sol";
import "./interfaces/HistoricAggregatorInterface.sol";
import "./vendor/Ownable.sol";

/**
 * @title The ConversionProxy contract for Solidity 4
 * @notice This contract allows for the rate of one aggregator
 * contract to be represented in the currency of another aggregator
 * contract's current rate. Rounds and timestamps are referred to
 * relative to the _from address. Historic answers are provided at
 * the latest rate of _to address.
 */
contract ConversionProxy is HistoricAggregatorInterface, Ownable {
  using SignedSafeMath for int256;

  uint8 public decimals;
  HistoricAggregatorInterface public from;
  HistoricAggregatorInterface public to;

  event AddressesUpdated(
    address from,
    address to,
    uint8 decimals
  );

  /**
   * @notice Deploys the ConversionProxy contract
   * @param _decimals The number of decimals that the result will
   * be returned with (should be specified for the `_to` aggregator)
   * @param _from The address of the aggregator contract which
   * needs to be converted
   * @param _to The address of the aggregator contract which stores
   * the rate to convert to
   */
  constructor(
    uint8 _decimals,
    address _from,
    address _to
  ) public Ownable() {
    setAddresses(
      _decimals,
      _from,
      _to
    );
  }

  /**
   * @dev Only callable by the owner of the contract
   * @param _decimals The number of decimals that the result will
   * be returned with (should be specified for the `_to` aggregator)
   * @param _from The address of the aggregator contract which
   * needs to be converted
   * @param _to The address of the aggregator contract which stores
   * the rate to convert to
   */
  function setAddresses(
    uint8 _decimals,
    address _from,
    address _to
  ) public onlyOwner() {
    require(_decimals > 0, "Decimals must be greater than 0");
    require(_from != _to, "Cannot use same address");
    decimals = _decimals;
    from = HistoricAggregatorInterface(_from);
    to = HistoricAggregatorInterface(_to);
    emit AddressesUpdated(
      _from,
      _to,
      _decimals
    );
  }

  /**
   * @notice Converts the latest answer of the `from` aggregator
   * to the rate of the `to` aggregator
   * @return The converted answer with amount of precision as defined
   * by `decimals`
   */
  function latestAnswer()
    external
    returns (int256)
  {
    return convertAnswer(from.latestAnswer(), to.latestAnswer());
  }

  /**
   * @notice Calls the `latestTimestamp()` function of the `from`
   * aggregator
   * @return The value of latestTimestamp for the `from` aggregator
   */
  function latestTimestamp()
    external
    returns (uint256)
  {
    return from.latestTimestamp();
  }

  /**
   * @notice Calls the `latestRound()` function of the `from`
   * aggregator
   * @return The value of latestRound for the `from` aggregator
   */
  function latestRound()
    external
    returns (uint256)
  {
    return from.latestRound();
  }

  /**
   * @notice Converts the specified answer for `_roundId` of the
   * `from` aggregator to the latestAnswer of the `to` aggregator
   * @return The converted answer for `_roundId` of the `from`
   * aggregator with the amount of precision as defined by `decimals`
   */
  function getAnswer(uint256 _roundId)
    external
    returns (int256)
  {
    return convertAnswer(from.getAnswer(_roundId), to.latestAnswer());
  }

  /**
   * @notice Calls the `getTimestamp(_roundId)` function of the `from`
   * aggregator for the specified `_roundId`
   * @return The timestamp of the `from` aggregator for the specified
   * `_roundId`
   */
  function getTimestamp(uint256 _roundId)
    external
    returns (uint256)
  {
    return from.getTimestamp(_roundId);
  }

  /**
   * @notice Converts the answer of the `from` aggregator to the rate
   * of the `to` aggregator at the precision of `decimals`
   * @return The converted answer
   */
  function convertAnswer(
    int256 _answerFrom,
    int256 _answerTo
  ) internal view returns (int256) {
    return _answerFrom.mul(_answerTo).div(int256(10 ** uint256(decimals)));
  }
}
