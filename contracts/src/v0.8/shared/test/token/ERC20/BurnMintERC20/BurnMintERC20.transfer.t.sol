// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";
import {BurnMintERC20} from "../../../../token/ERC20/BurnMintERC20.sol";

contract BurnMintERC20_transfer is BurnMintERC20Setup {
  function test_transfer() public {
    uint256 balancePre = s_burnMintERC20.balanceOf(STRANGER);
    uint256 sendingAmount = s_amount / 2;

    s_burnMintERC20.transfer(STRANGER, sendingAmount);

    assertEq(sendingAmount + balancePre, s_burnMintERC20.balanceOf(STRANGER));
  }

  // Reverts

  function test_transfer_RevertWhen_InvalidAddress() public {
    vm.expectRevert(abi.encodeWithSelector(BurnMintERC20.InvalidRecipient.selector, address(s_burnMintERC20)));

    s_burnMintERC20.transfer(address(s_burnMintERC20), s_amount);
  }
}
