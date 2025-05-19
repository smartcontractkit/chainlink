// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";

contract FactoryBurnMintERC20_decreaseApproval is BurnMintERC20Setup {
  function test_DecreaseApproval() public {
    s_burnMintERC20.approve(s_mockPool, s_amount);
    uint256 allowance = s_burnMintERC20.allowance(OWNER, s_mockPool);
    assertEq(allowance, s_amount);
    s_burnMintERC20.decreaseApproval(s_mockPool, s_amount);
    assertEq(s_burnMintERC20.allowance(OWNER, s_mockPool), allowance - s_amount);
  }
}
