// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IAccessControl} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/access/IAccessControl.sol";
import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC20.sol";

import {ERC20UpgradableBaseTest, IERC20UpgradeableBase} from "./ERC20UpgradableBaseTest.t.sol";

contract ERC20UpgradableBaseTest_mint is ERC20UpgradableBaseTest {
  function should_Mint(address implementation) public {
    changePrank(i_mockPool);

    uint256 balanceBefore = IERC20UpgradeableBase(implementation).balanceOf(STRANGER);

    vm.expectEmit();
    emit IERC20.Transfer(address(0), STRANGER, AMOUNT);

    IERC20UpgradeableBase(implementation).mint(STRANGER, AMOUNT);

    assertEq(IERC20UpgradeableBase(implementation).balanceOf(STRANGER), balanceBefore + AMOUNT);
    assertEq(IERC20UpgradeableBase(implementation).totalSupply(), PRE_MINT + AMOUNT);
  }

  function should_Mint_RevertWhen_AmountExceedsMaxSupply(address implementation, bytes4 errorSelector) public {
    changePrank(i_mockPool);

    uint256 amountToMint = IERC20UpgradeableBase(implementation).maxSupply() + AMOUNT;

    vm.expectRevert(abi.encodeWithSelector(errorSelector, amountToMint));

    IERC20UpgradeableBase(implementation).mint(STRANGER, amountToMint);
  }

  function should_Mint_RevertWhen_CallerDoesNotHaveMinterRole(address implementation, bytes32 MINTER_ROLE) public {
    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, STRANGER, MINTER_ROLE)
    );

    IERC20UpgradeableBase(implementation).mint(STRANGER, AMOUNT);
  }

  function should_Mint_RevertWhen_RecipientIsImplementationItself(address implementation, bytes4 errorSelector) public {
    changePrank(i_mockPool);

    vm.expectRevert(abi.encodeWithSelector(errorSelector, implementation));

    IERC20UpgradeableBase(implementation).mint(implementation, AMOUNT);
  }
}
