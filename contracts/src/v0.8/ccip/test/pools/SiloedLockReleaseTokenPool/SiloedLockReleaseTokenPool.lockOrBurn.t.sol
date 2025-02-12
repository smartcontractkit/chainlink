// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Pool} from "../../../libraries/Pool.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_lockOrBurn is SiloedLockReleaseTokenPoolSetup {
  uint256 public constant AMOUNT = 10e18;

  function test_lockOrBurn_SiloedFunds() public {
    assertTrue(s_siloedLockReleaseTokenPool.isSiloed(SILOED_CHAIN_SELECTOR));

    vm.startPrank(s_allowedOnRamp);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(AMOUNT);
    vm.expectEmit();
    emit TokenPool.Locked(s_allowedOnRamp, AMOUNT);

    s_siloedLockReleaseTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: STRANGER,
        receiver: bytes(""),
        amount: AMOUNT,
        remoteChainSelector: SILOED_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(SILOED_CHAIN_SELECTOR), AMOUNT);
  }

  function test_lockOrBurn_UnsiloedFunds() public {
    vm.startPrank(s_allowedOnRamp);

    assertFalse(s_siloedLockReleaseTokenPool.isSiloed(DEST_CHAIN_SELECTOR));

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(AMOUNT);
    vm.expectEmit();
    emit TokenPool.Locked(s_allowedOnRamp, AMOUNT);

    s_siloedLockReleaseTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: STRANGER,
        receiver: bytes(""),
        amount: AMOUNT,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(DEST_CHAIN_SELECTOR), AMOUNT);
  }

  // Reverts
}
