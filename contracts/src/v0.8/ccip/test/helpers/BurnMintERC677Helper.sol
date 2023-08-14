// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";

contract BurnMintERC677Helper is BurnMintERC677 {
  constructor(string memory name, string memory symbol) BurnMintERC677(name, symbol, 18, 0) {}

  // Gives one full token to any given address.
  function drip(address to) external {
    _mint(to, 1e18);
  }
}
