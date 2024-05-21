// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {LinkToken} from "../../token/ERC677/LinkToken.sol";

// This contract exists to mirror the functionality of the old token, which
// always deployed with 1b tokens sent to the deployer.
contract LinkTokenTestHelper is LinkToken {
  constructor() {
    _mint(msg.sender, 1e27);
  }
}
