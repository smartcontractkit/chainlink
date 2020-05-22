pragma solidity ^0.6.0;

import "../dev/AccessControl.sol";

contract AccessControlTestHelper is AccessControl {

  int256 private value;

  constructor(int256 _value)
    public
  {
    value = _value;
  }

  function getValue()
    external
    view
    checkAccess()
    returns (int256)
  {
    return value;
  }

}
