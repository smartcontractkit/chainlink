pragma solidity ^0.6.0;

import "../dev/SimpleAccessControl.sol";

contract AccessControlTestHelper is SimpleAccessControl {

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
