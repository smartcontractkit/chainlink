pragma solidity ^0.6.0;

import "../dev/Whitelisted.sol";

contract WhitelistedTestHelper is Whitelisted {

  int256 private value;

  constructor(int256 _value)
    public
  {
    value = _value;
  }

  function getValue()
    external
    view
    isWhitelisted()
    returns (int256)
  {
    return value;
  }

}
