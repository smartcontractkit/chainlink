// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";
import {BurnMintERC20} from "../../../../token/ERC20/BurnMintERC20.sol";

contract BurnMintERC20grantRole is BurnMintERC20Setup {
  function test_GrantMintAccess_Success() public {
    assertFalse(s_burnMintERC20.isMinter(STRANGER));

    vm.expectEmit();
    emit BurnMintERC20.MintAccessGranted(STRANGER);

    s_burnMintERC20.grantMintAndBurnRoles(STRANGER);

    assertTrue(s_burnMintERC20.isMinter(STRANGER));

    vm.expectEmit();
    emit BurnMintERC20.MintAccessRevoked(STRANGER);

    s_burnMintERC20.revokeMintRole(STRANGER);

    assertFalse(s_burnMintERC20.isMinter(STRANGER));
  }

  function test_GrantBurnAccess_Success() public {
    assertFalse(s_burnMintERC20.isBurner(STRANGER));

    vm.expectEmit();
    emit BurnMintERC20.BurnAccessGranted(STRANGER);

    s_burnMintERC20.grantBurnRole(STRANGER);

    assertTrue(s_burnMintERC20.isBurner(STRANGER));

    vm.expectEmit();
    emit BurnMintERC20.BurnAccessRevoked(STRANGER);

    s_burnMintERC20.revokeBurnRole(STRANGER);

    assertFalse(s_burnMintERC20.isBurner(STRANGER));
  }

  function test_GrantMany_Success() public {
    // Since the owner was already granted mint and burn roles in the constructor, we will revoke them
    // and then grant them again for the purposes of the test
    s_burnMintERC20.revokeMintRole(OWNER);
    s_burnMintERC20.revokeBurnRole(OWNER);

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
