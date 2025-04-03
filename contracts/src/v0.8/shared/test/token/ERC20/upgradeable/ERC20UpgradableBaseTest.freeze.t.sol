// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IBurnMintERC20Upgradeable} from "../../../../../shared/token/ERC20/upgradeable/IBurnMintERC20Upgradeable.sol";
import {PausableUpgradeable} from "../../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/utils/PausableUpgradeable.sol";
import {IAccessControl} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/access/IAccessControl.sol";
import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC20.sol";

import {ERC20UpgradableBaseTest, IERC20UpgradeableBase} from "./ERC20UpgradableBaseTest.t.sol";

contract ERC20UpgradableBaseTest_freeze is ERC20UpgradableBaseTest {
  function should_Freeze(address implementation) public {
    changePrank(DEFAULT_FREEZER);

    vm.expectEmit();
    emit IERC20UpgradeableBase.AccountFrozen(OWNER);
    IERC20UpgradeableBase(implementation).freeze(OWNER);

    assertTrue(IERC20UpgradeableBase(implementation).isFrozen(OWNER));
  }

  function should_Freeze_EvenWhenImplementationIsPaused(address implementation) public {
    changePrank(DEFAULT_PAUSER);
    IERC20UpgradeableBase(implementation).pause();
    assertTrue(PausableUpgradeable(implementation).paused());

    should_Freeze(implementation);
  }

  function should_Freeze_RevertWhen_CallerDoesNotHaveFreezerRole(address implementation, bytes32 FREEZER_ROLE) public {
    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, STRANGER, FREEZER_ROLE)
    );

    IERC20UpgradeableBase(implementation).freeze(OWNER);
  }

  function should_Freeze_RevertWhen_RecipientIsAddressZero(address implementation, bytes4 errorSelector) public {
    changePrank(DEFAULT_FREEZER);

    vm.expectRevert(abi.encodeWithSelector(errorSelector, address(0)));

    IERC20UpgradeableBase(implementation).freeze(address(0));
  }

  function should_Freeze_RevertWhen_RecipientIsImplementationItself(
    address implementation,
    bytes4 errorSelector
  ) public {
    changePrank(DEFAULT_FREEZER);

    vm.expectRevert(abi.encodeWithSelector(errorSelector, implementation));

    IERC20UpgradeableBase(implementation).freeze(implementation);
  }

  function should_Freeze_RevertWhen_AccountIsAlreadyFrozen(address implementation, bytes4 errorSelector) public {
    changePrank(DEFAULT_FREEZER);
    IERC20UpgradeableBase(implementation).freeze(OWNER);

    vm.expectRevert(abi.encodeWithSelector(errorSelector, OWNER));

    IERC20UpgradeableBase(implementation).freeze(OWNER);
  }

  function should_Mint_RevertWhen_AccountIsFrozen(address implementation, bytes4 errorSelector) public {
    changePrank(DEFAULT_FREEZER);
    IERC20UpgradeableBase(implementation).freeze(OWNER);

    changePrank(DEFAULT_ADMIN);
    IERC20UpgradeableBase(implementation).grantMintAndBurnRoles(DEFAULT_ADMIN);

    vm.expectRevert(abi.encodeWithSelector(errorSelector, OWNER));
    IBurnMintERC20Upgradeable(implementation).mint(OWNER, AMOUNT);
  }

  function should_Transfer_RevertWhen_SenderOrRecipientAreFrozen(
    address implementation,
    bytes4 accountFrozenErrorSelector
  ) public {
    changePrank(DEFAULT_FREEZER);
    IERC20UpgradeableBase(implementation).freeze(OWNER);

    changePrank(OWNER);

    vm.expectRevert(abi.encodeWithSelector(accountFrozenErrorSelector, OWNER));
    IERC20(implementation).transfer(STRANGER, AMOUNT);

    changePrank(DEFAULT_FREEZER);
    IERC20UpgradeableBase(implementation).unfreeze(OWNER);
    IERC20UpgradeableBase(implementation).freeze(STRANGER);

    changePrank(OWNER);
    vm.expectRevert(abi.encodeWithSelector(accountFrozenErrorSelector, STRANGER));
    IERC20(implementation).transfer(STRANGER, AMOUNT);
  }

  function should_Approve_RevertWhen_OwnerOrSpenderAreFrozen(
    address implementation,
    bytes4 accountFrozenErrorSelector
  ) public {
    changePrank(DEFAULT_ADMIN);
    IERC20UpgradeableBase(implementation).grantMintAndBurnRoles(DEFAULT_ADMIN);
    IBurnMintERC20Upgradeable(implementation).mint(OWNER, AMOUNT);

    changePrank(DEFAULT_FREEZER);
    IERC20UpgradeableBase(implementation).freeze(OWNER);

    changePrank(OWNER);

    vm.expectRevert(abi.encodeWithSelector(accountFrozenErrorSelector, OWNER));
    IERC20(implementation).approve(STRANGER, AMOUNT);

    changePrank(DEFAULT_FREEZER);
    IERC20UpgradeableBase(implementation).unfreeze(OWNER);
    IERC20UpgradeableBase(implementation).freeze(STRANGER);

    changePrank(OWNER);

    vm.expectRevert(abi.encodeWithSelector(accountFrozenErrorSelector, STRANGER));
    IERC20(implementation).approve(STRANGER, AMOUNT);
  }
}
