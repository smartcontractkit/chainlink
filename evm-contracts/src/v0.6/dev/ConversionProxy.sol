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

  function latestAnswer() external view virtual override returns (int256) {
    return _latestAnswer();
  }

  function latestTimestamp() external view virtual override returns (uint256) {
    return _latestTimestamp();
  }

  function latestRound() external view virtual override returns (uint256) {
    return _latestRound();
  }

  function getAnswer(uint256 _roundId) external view virtual override returns (int256) {
    return _getAnswer(_roundId);
  }

  function getTimestamp(uint256 _roundId) external view virtual override returns (uint256) {
    return _getTimestamp(_roundId);
  }

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

  function convertAnswer(
    int256 _answerFrom,
    int256 _answerTo
  ) internal view returns (int256) {
    return _answerFrom.mul(_answerTo).div(int256(10 ** uint256(to.decimals())));
  }
}
