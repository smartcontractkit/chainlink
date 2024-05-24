// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {ConfirmedOwner} from "../../access/ConfirmedOwner.sol";

contract ConfirmedOwnerTestHelper is ConfirmedOwner {
  event Here();

  constructor() ConfirmedOwner(msg.sender) {}

  function modifierOnlyOwner() public onlyOwner {
    emit Here();
  }
}
