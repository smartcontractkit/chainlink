pragma solidity 0.6.2;

import "../dev/SignedSafeMath.sol";
import "./AggregatorInterface.sol";
import "../Owned.sol";

contract ConversionProxy is AggregatorInterface, Owned {
  using SignedSafeMath for int256;

  AggregatorInterface public from;
  AggregatorInterface public to;

  event AddressesUpdated(
    address from,
    address to
  );

  constructor(
    address _from,
    address _to
  ) public Owned() {
    setAddresses(
      _from,
      _to
    );
  }

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

  function latestAnswer() external view override returns (int256) {
    return convertAnswer(from.latestAnswer(), to.latestAnswer());
  }

  function latestTimestamp() external view override returns (uint256) {
    return from.latestTimestamp();
  }

  function latestRound() external view override returns (uint256) {
    return from.latestRound();
  }

  function getAnswer(uint256 _roundId) external view override returns (int256) {
    return convertAnswer(from.getAnswer(_roundId), to.latestAnswer());
  }

  function getTimestamp(uint256 _roundId) external view override returns (uint256) {
    return from.getTimestamp(_roundId);
  }

  function decimals() external view override returns (uint8) {
    return to.decimals();
  }

  function convertAnswer(
    int256 _answerFrom,
    int256 _answerTo
  ) internal view returns (int256) {
    return _answerFrom.mul(_answerTo).div(int256(10 ** uint256(to.decimals())));
  }
}
