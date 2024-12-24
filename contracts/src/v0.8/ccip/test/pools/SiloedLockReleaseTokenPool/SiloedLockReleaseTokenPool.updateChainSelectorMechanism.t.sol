// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_updateChainSelectorMechanism is SiloedLockReleaseTokenPoolSetup {
  function test_updateChainSelectorMechanism_Success() public {
    uint256 amount = 1e18;
    uint64[] memory chainSelectors = new uint64[](1);

    chainSelectors[0] = SILOED_CHAIN_SELECTOR;

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.ChainSiloeDesignationUpdated(SILOED_CHAIN_SELECTOR, true);

    s_siloedLockReleaseTokenPool.updateChainSelectorMechanisms(new uint64[](0), chainSelectors);

    // Assert that the funds are siloed correctly
    assertTrue(s_siloedLockReleaseTokenPool.chainFundsAreSiloed(SILOED_CHAIN_SELECTOR));

    // Provide some Liquidity so that we can then check that it gets removed.
    s_siloedLockReleaseTokenPool.setRebalancer(SILOED_CHAIN_SELECTOR, OWNER);
    s_siloedLockReleaseTokenPool.provideLiquidity(SILOED_CHAIN_SELECTOR, amount);
    assertEq(s_siloedLockReleaseTokenPool.getLockedTokensByChain(SILOED_CHAIN_SELECTOR), amount);

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.ChainSiloeDesignationUpdated(SILOED_CHAIN_SELECTOR, false);

    s_siloedLockReleaseTokenPool.updateChainSelectorMechanisms(chainSelectors, new uint64[](0));

    // Check that the locked funds accounting was cleared when the funds were un-siloed.
    assertFalse(s_siloedLockReleaseTokenPool.chainFundsAreSiloed(SILOED_CHAIN_SELECTOR));
    assertEq(s_siloedLockReleaseTokenPool.getLockedTokensByChain(SILOED_CHAIN_SELECTOR), 0);
  }
}
