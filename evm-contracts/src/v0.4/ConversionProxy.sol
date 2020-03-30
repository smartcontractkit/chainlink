pragma solidity 0.4.24;

import "./vendor/SignedSafeMath.sol";
import "./interfaces/AggregatorInterface.sol";
import "./vendor/Ownable.sol";

contract ConversionProxy is AggregatorInterface, Ownable {
  using SignedSafeMath for int256;

  uint8 public decimals;
  AggregatorInterface public from;
  AggregatorInterface public to;

  event AddressesUpdated(
    uint8 decimals,
    address from,
    address to
  );

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

  function setAddresses(
    uint8 _decimals,
    address _from,
    address _to
  ) public onlyOwner() {
    require(_decimals > 0, "Decimals must be greater than 0");
    require(_from != _to, "Cannot use same address");
    decimals = _decimals;
    from = AggregatorInterface(_from);
    to = AggregatorInterface(_to);
    emit AddressesUpdated(
      _decimals,
      _from,
      _to
    );
  }

  function latestAnswer() external view returns (int256) {
    return convertAnswer(from.latestAnswer(), to.latestAnswer());
  }

  function latestTimestamp() external view returns (uint256) {
    return from.latestTimestamp();
  }

  function latestRound() external view returns (uint256) {
    return from.latestRound();
  }

  function getAnswer(uint256 _roundId) external view returns (int256) {
    return convertAnswer(from.getAnswer(_roundId), to.latestAnswer());
  }

  function getTimestamp(uint256 _roundId) external view returns (uint256) {
    return from.getTimestamp(_roundId);
  }

  function convertAnswer(
    int256 _answerFrom,
    int256 _answerTo
  ) internal view returns (int256) {
    return _answerFrom.mul(_answerTo).div(int256(10 ** uint256(decimals)));
  }
}
