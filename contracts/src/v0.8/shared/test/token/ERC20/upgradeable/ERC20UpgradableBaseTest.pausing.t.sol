// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IBurnMintERC20Upgradeable} from "../../../../../shared/token/ERC20/upgradeable/IBurnMintERC20Upgradeable.sol";
import {PausableUpgradeable} from "../../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/utils/PausableUpgradeable.sol";
import {IAccessControl} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/access/IAccessControl.sol";
import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC20.sol";

import {ERC20UpgradableBaseTest, IERC20UpgradeableBase} from "./ERC20UpgradableBaseTest.t.sol";

contract ERC20UpgradableBaseTest_pausing is ERC20UpgradableBaseTest {
  function should_Pause(address implementation) public {
    changePrank(DEFAULT_PAUSER);

    vm.expectEmit();
    emit PausableUpgradeable.Paused(DEFAULT_PAUSER);

    IERC20UpgradeableBase(implementation).pause();

    assertTrue(PausableUpgradeable(implementation).paused());
  }

  function should_Unpause(address implementation) public {
    changePrank(DEFAULT_PAUSER);
    IERC20UpgradeableBase(implementation).pause();

    changePrank(DEFAULT_ADMIN);
    vm.expectEmit();
    emit PausableUpgradeable.Unpaused(DEFAULT_ADMIN);
    IERC20UpgradeableBase(implementation).unpause();

    assertFalse(PausableUpgradeable(implementation).paused());

    changePrank(i_mockPool);
    vm.expectEmit();
    emit IERC20.Transfer(address(0), STRANGER, AMOUNT);
    IBurnMintERC20Upgradeable(implementation).mint(STRANGER, AMOUNT);

    changePrank(STRANGER);
    IERC20(implementation).approve(OWNER, AMOUNT);
    assertEq(IERC20(implementation).allowance(STRANGER, OWNER), AMOUNT);
  }

  function should_Pause_RevertWhen_CallerDoesNotHavePauserRole(address implementation, bytes32 PAUSER_ROLE) public {
    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, STRANGER, PAUSER_ROLE)
    );

    IERC20UpgradeableBase(implementation).pause();
  }

  function should_Unpause_RevertWhen_CallerDoesNotHaveDefaultAdminRole(
    address implementation,
    bytes32 DEFAULT_ADMIN_ROLE
  ) public {
    changePrank(DEFAULT_PAUSER);
    IERC20UpgradeableBase(implementation).pause();

    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, STRANGER, DEFAULT_ADMIN_ROLE)
    );

    IERC20UpgradeableBase(implementation).unpause();
  }

  function should_Mint_RevertWhen_ImplementationIsPaused(address implementation) public {
    changePrank(DEFAULT_PAUSER);
    IERC20UpgradeableBase(implementation).pause();

    changePrank(i_mockPool);

    vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));
    IBurnMintERC20Upgradeable(implementation).mint(STRANGER, AMOUNT);
  }

  function should_Transfer_RevertWhen_ImplementationIsPaused(address implementation) public {
    changePrank(i_mockPool);
    IBurnMintERC20Upgradeable(implementation).mint(STRANGER, AMOUNT);

    changePrank(DEFAULT_PAUSER);
    IERC20UpgradeableBase(implementation).pause();

    changePrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));
    IERC20(implementation).transfer(OWNER, AMOUNT);
  }

  function should_Burn_RevertWhen_ImplementationIsPaused(address implementation) public {
    changePrank(DEFAULT_PAUSER);
    IERC20UpgradeableBase(implementation).pause();

    changePrank(i_mockPool);

    vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));

    IBurnMintERC20Upgradeable(implementation).burn(0);
  }

  function should_BurnFrom_RevertWhen_ImplementationIsPaused(address implementation) public {
    changePrank(DEFAULT_PAUSER);
    IERC20UpgradeableBase(implementation).pause();

    changePrank(i_mockPool);

    vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));
    IBurnMintERC20Upgradeable(implementation).burnFrom(STRANGER, 0);
  }

  function should_BurnFrom_alias_RevertWhen_ImplementationIsPaused(address implementation) public {
    changePrank(DEFAULT_PAUSER);
    IERC20UpgradeableBase(implementation).pause();

    changePrank(i_mockPool);

    vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));
    IBurnMintERC20Upgradeable(implementation).burn(STRANGER, 0);
  }

  function should_Approve_RevertWhen_ImplementationIsPaused(address implementation) public {
    changePrank(DEFAULT_PAUSER);
    IERC20UpgradeableBase(implementation).pause();

    changePrank(STRANGER);

    vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));

    IERC20(implementation).approve(i_mockPool, AMOUNT);
  }
}
