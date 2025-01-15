// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {BurnToAddressMintTokenPool} from "../../../pools/BurnToAddressMintTokenPool.sol";
import {BurnToAddressMintTokenPoolSetup} from "./BurnToAddressMintTokenPoolSetup.t.sol";

contract BurnToAddressMintTokenPool_setLockedTokens is BurnToAddressMintTokenPoolSetup {
  function test_setLockedTokens() public {
    uint256 amount = 1e18;

    assertEq(s_pool.getMintedTokens(), 0);

    vm.expectEmit();
    emit BurnToAddressMintTokenPool.MintedTokensSet(amount, 0);

    s_pool.setMintedTokens(amount);

    assertEq(s_pool.getMintedTokens(), amount);
  }
}
