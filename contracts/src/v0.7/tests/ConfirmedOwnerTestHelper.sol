// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../ConfirmedOwner.sol";

contract ConfirmedOwnerTestHelper is ConfirmedOwner {

  event Here();

  constructor() 
    ConfirmedOwner(msg.sender) 
  {}

  function modifierOnlyOwner()
    public
    onlyOwner()
  {
    emit Here();
  }

}
