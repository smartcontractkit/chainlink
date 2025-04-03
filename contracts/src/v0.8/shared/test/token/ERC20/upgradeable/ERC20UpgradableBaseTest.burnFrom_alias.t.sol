// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IBurnMintERC20Upgradeable} from "../../../../../shared/token/ERC20/upgradeable/IBurnMintERC20Upgradeable.sol";
import {IAccessControl} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/access/IAccessControl.sol";
import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC20.sol";

import {ERC20UpgradableBaseTest} from "./ERC20UpgradableBaseTest.t.sol";

contract ERC20UpgradableBaseTest_burnFrom_alias is ERC20UpgradableBaseTest {
  function should_BurnFrom_alias(address implementation) public {
    changePrank(i_mockPool);
    IBurnMintERC20Upgradeable(implementation).mint(STRANGER, AMOUNT);

    uint256 balanceBefore = IBurnMintERC20Upgradeable(implementation).balanceOf(STRANGER);
    uint256 totalSupplyBefore = IBurnMintERC20Upgradeable(implementation).totalSupply();
    uint256 amountToBurn = AMOUNT / 2;

    changePrank(STRANGER);
    IBurnMintERC20Upgradeable(implementation).approve(i_mockPool, amountToBurn);

    changePrank(i_mockPool);

    vm.expectEmit();
    emit IERC20.Transfer(STRANGER, address(0), amountToBurn);

    // burn(account, amount) is alias for burnFrom(account, amount)
    IBurnMintERC20Upgradeable(implementation).burn(STRANGER, amountToBurn);

    assertEq(IBurnMintERC20Upgradeable(implementation).balanceOf(STRANGER), balanceBefore - amountToBurn);
    assertEq(IBurnMintERC20Upgradeable(implementation).totalSupply(), totalSupplyBefore - amountToBurn);
  }

  function should_BurnFrom_alias_RevertWhen_CallerDoesNotHaveBurnerRole(
    address implementation,
    bytes32 BURNER_ROLE
  ) public {
    changePrank(i_mockPool);
    IBurnMintERC20Upgradeable(implementation).mint(STRANGER, AMOUNT);

    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, STRANGER, BURNER_ROLE)
    );

    IBurnMintERC20Upgradeable(implementation).burn(STRANGER, AMOUNT);
  }
}
