// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {FactoryBurnMintERC20} from "../../../tokenAdminRegistry/TokenPoolFactory/FactoryBurnMintERC20.sol";
import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";

contract FactoryBurnMintERC20_grantRole is BurnMintERC20Setup {
  function test_GrantMintAccess() public {
    assertFalse(s_burnMintERC20.isMinter(STRANGER));

    vm.expectEmit();
    emit FactoryBurnMintERC20.MintAccessGranted(STRANGER);

    s_burnMintERC20.grantMintAndBurnRoles(STRANGER);

    assertTrue(s_burnMintERC20.isMinter(STRANGER));

    vm.expectEmit();
    emit FactoryBurnMintERC20.MintAccessRevoked(STRANGER);

    s_burnMintERC20.revokeMintRole(STRANGER);

    assertFalse(s_burnMintERC20.isMinter(STRANGER));
  }

  function test_GrantBurnAccess() public {
    assertFalse(s_burnMintERC20.isBurner(STRANGER));

    vm.expectEmit();
    emit FactoryBurnMintERC20.BurnAccessGranted(STRANGER);

    s_burnMintERC20.grantBurnRole(STRANGER);

    assertTrue(s_burnMintERC20.isBurner(STRANGER));

    vm.expectEmit();
    emit FactoryBurnMintERC20.BurnAccessRevoked(STRANGER);

    s_burnMintERC20.revokeBurnRole(STRANGER);

    assertFalse(s_burnMintERC20.isBurner(STRANGER));
  }

  function test_GrantMany() public {
    // Since alice was already granted mint and burn roles in the setup, we will revoke them
    // and then grant them again for the purposes of the test
    s_burnMintERC20.revokeMintRole(s_alice);
    s_burnMintERC20.revokeBurnRole(s_alice);

    uint256 numberOfPools = 10;
    address[] memory permissionedAddresses = new address[](numberOfPools + 1);
    permissionedAddresses[0] = s_mockPool;

    for (uint160 i = 0; i < numberOfPools; ++i) {
      permissionedAddresses[i + 1] = address(i);
      s_burnMintERC20.grantMintAndBurnRoles(address(i));
    }

    assertEq(permissionedAddresses, s_burnMintERC20.getBurners());
    assertEq(permissionedAddresses, s_burnMintERC20.getMinters());
  }
}
