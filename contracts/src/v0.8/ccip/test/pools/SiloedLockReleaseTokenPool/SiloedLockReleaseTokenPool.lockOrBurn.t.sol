// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Pool} from "../../../libraries/Pool.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";

import {TokenPool} from "../../../pools/TokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_lockOrBurn is SiloedLockReleaseTokenPoolSetup {
  function test_LockOrBurn_SiloedFunds_Success(
    uint256 amount
  ) public {
    amount = 10e18;

    vm.startPrank(s_allowedOnRamp);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(amount);
    vm.expectEmit();
    emit TokenPool.Locked(s_allowedOnRamp, amount);

    s_siloedLockReleaseTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: STRANGER,
        receiver: bytes(""),
        amount: amount,
        remoteChainSelector: SILOED_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    assertEq(s_siloedLockReleaseTokenPool.getSiloedTokensByChain(SILOED_CHAIN_SELECTOR), amount);
  }

  function testLockOrBurn_NonSiloedFunds(
    uint256 amount
  ) public {
    amount = 10e18;
    vm.startPrank(s_allowedOnRamp);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(amount);
    vm.expectEmit();
    emit TokenPool.Locked(s_allowedOnRamp, amount);

    s_siloedLockReleaseTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: STRANGER,
        receiver: bytes(""),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    assertEq(s_siloedLockReleaseTokenPool.getSiloedTokensByChain(DEST_CHAIN_SELECTOR), amount);
  }
}
