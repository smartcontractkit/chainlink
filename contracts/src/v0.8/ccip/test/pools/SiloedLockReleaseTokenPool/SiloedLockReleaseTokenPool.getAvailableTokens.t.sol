// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_getAvailableTokens is SiloedLockReleaseTokenPoolSetup {
  function test_getAvailableTokens_UnsiloedChain() public {
    uint256 amount = 1e24;

    s_siloedLockReleaseTokenPool.provideLiquidity(amount);

    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(DEST_CHAIN_SELECTOR), amount);
  }

  function test_getAvailableTokens_SiloedChain() public {
    uint256 amount = 1e24;

    s_siloedLockReleaseTokenPool.provideSiloedLiquidity(SILOED_CHAIN_SELECTOR, amount);

    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(SILOED_CHAIN_SELECTOR), amount);
  }

  function test_getAvailableTokens_RevertWhen_UnsupportedChain() public {
    vm.expectRevert(abi.encodeWithSelector(SiloedLockReleaseTokenPool.InvalidChainSelector.selector, type(uint64).max));

    s_siloedLockReleaseTokenPool.getAvailableTokens(type(uint64).max);
  }
}
