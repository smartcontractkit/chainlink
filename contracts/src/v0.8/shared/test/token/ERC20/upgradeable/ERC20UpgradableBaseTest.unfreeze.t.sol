// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IBurnMintERC20Upgradeable} from "../../../../../shared/token/ERC20/upgradeable/IBurnMintERC20Upgradeable.sol";
import {PausableUpgradeable} from "../../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/utils/PausableUpgradeable.sol";
import {IAccessControl} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/access/IAccessControl.sol";
import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC20.sol";

import {ERC20UpgradableBaseTest, IERC20UpgradeableBase} from "./ERC20UpgradableBaseTest.t.sol";

contract ERC20UpgradableBaseTest_unfreeze is ERC20UpgradableBaseTest {
  function should_Unfreeze(address implementation) public {
    changePrank(DEFAULT_FREEZER);
    IERC20UpgradeableBase(implementation).freeze(OWNER);

    vm.expectEmit();
    emit IERC20UpgradeableBase.AccountUnfrozen(OWNER);
    IERC20UpgradeableBase(implementation).unfreeze(OWNER);

    assertFalse(IERC20UpgradeableBase(implementation).isFrozen(OWNER));

    changePrank(DEFAULT_ADMIN);
    IERC20UpgradeableBase(implementation).grantMintAndBurnRoles(DEFAULT_ADMIN);
    IBurnMintERC20Upgradeable(implementation).mint(OWNER, AMOUNT);
    assertEq(IBurnMintERC20Upgradeable(implementation).balanceOf(OWNER), AMOUNT);

    changePrank(OWNER);
    IERC20(implementation).approve(STRANGER, AMOUNT);
  }

  function should_Unfreeze_EvenWhenImplementationIsPaused(address implementation) public {
    changePrank(DEFAULT_FREEZER);
    IERC20UpgradeableBase(implementation).freeze(OWNER);
    assertTrue(IERC20UpgradeableBase(implementation).isFrozen(OWNER));

    changePrank(DEFAULT_PAUSER);
    IERC20UpgradeableBase(implementation).pause();
    assertTrue(PausableUpgradeable(implementation).paused());

    changePrank(DEFAULT_FREEZER);
    IERC20UpgradeableBase(implementation).unfreeze(OWNER);
    assertFalse(IERC20UpgradeableBase(implementation).isFrozen(OWNER));
  }

  function should_Unfreeze_RevertWhen_CallerDoesNotHaveFreezerRole(
    address implementation,
    bytes32 FREEZER_ROLE
  ) public {
    changePrank(DEFAULT_FREEZER);
    IERC20UpgradeableBase(implementation).freeze(OWNER);

    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, STRANGER, FREEZER_ROLE)
    );

    IERC20UpgradeableBase(implementation).unfreeze(OWNER);
  }

  function should_Unfreeze_RevertWhen_AccountIsNotFrozen(address implementation, bytes4 errorSelector) public {
    changePrank(DEFAULT_FREEZER);

    assertFalse(IERC20UpgradeableBase(implementation).isFrozen(STRANGER));

    vm.expectRevert(abi.encodeWithSelector(errorSelector, STRANGER));

    IERC20UpgradeableBase(implementation).unfreeze(STRANGER);
  }
}
