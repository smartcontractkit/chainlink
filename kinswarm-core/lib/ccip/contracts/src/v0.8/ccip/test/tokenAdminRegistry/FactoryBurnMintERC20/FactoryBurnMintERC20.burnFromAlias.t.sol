// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {FactoryBurnMintERC20} from "../../../tokenAdminRegistry/TokenPoolFactory/FactoryBurnMintERC20.sol";
import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";

contract FactoryBurnMintERC20_burnFromAlias is BurnMintERC20Setup {
  function setUp() public virtual override {
    BurnMintERC20Setup.setUp();
  }

  function test_BurnFrom_Success() public {
    s_burnMintERC20.approve(s_mockPool, s_amount);

    changePrank(s_mockPool);

    s_burnMintERC20.burn(OWNER, s_amount);

    assertEq(0, s_burnMintERC20.balanceOf(OWNER));
  }

  // Reverts

  function test_SenderNotBurner_Reverts() public {
    vm.expectRevert(abi.encodeWithSelector(FactoryBurnMintERC20.SenderNotBurner.selector, OWNER));

    s_burnMintERC20.burn(OWNER, s_amount);
  }

  function test_InsufficientAllowance_Reverts() public {
    changePrank(s_mockPool);

    vm.expectRevert("ERC20: insufficient allowance");

    s_burnMintERC20.burn(OWNER, s_amount);
  }

  function test_ExceedsBalance_Reverts() public {
    s_burnMintERC20.approve(s_mockPool, s_amount * 2);

    changePrank(s_mockPool);

    vm.expectRevert("ERC20: burn amount exceeds balance");

    s_burnMintERC20.burn(OWNER, s_amount * 2);
  }
}
