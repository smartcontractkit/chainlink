// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "../SimpleReadAccessController.sol";

contract AccessControlTestHelper is SimpleReadAccessController {

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
