// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IBurnMintERC20Upgradeable} from "../../../../../shared/token/ERC20/upgradeable/IBurnMintERC20Upgradeable.sol";
import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC20.sol";

import {ERC20UpgradableBaseTest} from "./ERC20UpgradableBaseTest.t.sol";

contract ERC20UpgradableBaseTest_approve is ERC20UpgradableBaseTest {
  function should_Approve(address implementation) public {
    changePrank(i_mockPool);
    IBurnMintERC20Upgradeable(implementation).mint(STRANGER, AMOUNT);

    changePrank(STRANGER);

    vm.expectEmit();
    emit IERC20.Approval(STRANGER, i_mockPool, AMOUNT);

    IBurnMintERC20Upgradeable(implementation).approve(i_mockPool, AMOUNT);

    assertEq(IBurnMintERC20Upgradeable(implementation).allowance(STRANGER, i_mockPool), AMOUNT);
  }

  function should_Approve_RevertWhen_RecipientIsImplementationItself(
    address implementation,
    bytes4 errorSelector
  ) public {
    changePrank(STRANGER);

    vm.expectRevert(abi.encodeWithSelector(errorSelector, implementation));

    IBurnMintERC20Upgradeable(implementation).approve(implementation, AMOUNT);
  }
}
