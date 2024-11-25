// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";

contract FactoryBurnMintERC20_approve is BurnMintERC20Setup {
  function test_Approve() public {
    uint256 balancePre = s_burnMintERC20.balanceOf(STRANGER);
    uint256 sendingAmount = s_amount / 2;

    s_burnMintERC20.approve(STRANGER, sendingAmount);

    changePrank(STRANGER);

    s_burnMintERC20.transferFrom(OWNER, STRANGER, sendingAmount);

    assertEq(sendingAmount + balancePre, s_burnMintERC20.balanceOf(STRANGER));
  }

  // Reverts

  function test_RevertWhen_InvalidAddresss() public {
    vm.expectRevert();

    s_burnMintERC20.approve(address(s_burnMintERC20), s_amount);
  }
}
