// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IAccessControl} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/access/IAccessControl.sol";

import {ERC20UpgradableBaseTest, IERC20UpgradeableBase} from "./ERC20UpgradableBaseTest.t.sol";

contract ERC20UpgradableBaseTest_initialize is ERC20UpgradableBaseTest {
  function should_Initialize(address implementation, bytes32 DEFAULT_ADMIN_ROLE) public view {
    assertEq(IERC20UpgradeableBase(implementation).name(), NAME);
    assertEq(IERC20UpgradeableBase(implementation).symbol(), SYMBOL);
    assertEq(IERC20UpgradeableBase(implementation).decimals(), DECIMALS);
    assertEq(IERC20UpgradeableBase(implementation).maxSupply(), MAX_SUPPLY);
    assertEq(IERC20UpgradeableBase(implementation).totalSupply(), PRE_MINT);

    assertTrue(IAccessControl(implementation).hasRole(DEFAULT_ADMIN_ROLE, DEFAULT_ADMIN));
  }

  function should_Initialize_WithPreMint(address implementation, uint256 newPreMint) public view {
    assertEq(IERC20UpgradeableBase(implementation).totalSupply(), newPreMint);
    assertEq(IERC20UpgradeableBase(implementation).balanceOf(DEFAULT_ADMIN), newPreMint);
  }
}
