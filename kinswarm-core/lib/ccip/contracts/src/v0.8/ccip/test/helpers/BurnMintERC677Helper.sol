// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {IGetCCIPAdmin} from "../../interfaces/IGetCCIPAdmin.sol";

contract BurnMintERC677Helper is BurnMintERC677, IGetCCIPAdmin {
  constructor(string memory name, string memory symbol) BurnMintERC677(name, symbol, 18, 0) {}

  // Gives one full token to any given address.
  function drip(
    address to
  ) external {
    _mint(to, 1e18);
  }

  function getCCIPAdmin() external view override returns (address) {
    return owner();
  }
}
