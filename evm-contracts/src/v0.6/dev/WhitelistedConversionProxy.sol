pragma solidity 0.6.2;

import "./ConversionProxy.sol";
import "./Whitelisted.sol";

contract WhitelistedConversionProxy is ConversionProxy, Whitelisted {
  constructor(
    address _from,
    address _to
  ) public ConversionProxy(
    _from,
    _to
  ) {}

  function latestAnswer()
    external
    view
    override
    isWhitelisted()
    returns (int256)
  {
    return _latestAnswer();
  }

  function latestTimestamp()
    external
    view
    override
    isWhitelisted()
    returns (uint256)
  {
    return _latestTimestamp();
  }

  function latestRound()
    external
    view
    override
    isWhitelisted()
    returns (uint256)
  {
    return _latestRound();
  }

  function getAnswer(uint256 _roundId)
    external
    view
    override
    isWhitelisted()
    returns (int256)
  {
    return _getAnswer(_roundId);
  }

  function getTimestamp(uint256 _roundId)
    external
    view
    override
    isWhitelisted()
    returns (uint256)
  {
    return _getTimestamp(_roundId);
  }
}
