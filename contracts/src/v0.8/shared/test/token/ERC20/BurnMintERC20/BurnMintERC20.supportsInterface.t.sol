// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IBurnMintERC20} from "../../../../../shared/token/ERC20/IBurnMintERC20.sol";

import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";

import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {IERC165} from "../../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/IERC165.sol";

contract BurnMintERC20_supportsInterface is BurnMintERC20Setup {
  function test_SupportsInterface() public view {
    assertTrue(s_burnMintERC20.supportsInterface(type(IERC20).interfaceId));
    assertTrue(s_burnMintERC20.supportsInterface(type(IBurnMintERC20).interfaceId));
    assertTrue(s_burnMintERC20.supportsInterface(type(IERC165).interfaceId));
  }
}
