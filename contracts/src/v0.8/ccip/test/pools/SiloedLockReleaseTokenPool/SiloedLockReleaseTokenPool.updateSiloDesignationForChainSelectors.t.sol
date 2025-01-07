// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_updateSiloDesignationForChainSelectors is SiloedLockReleaseTokenPoolSetup {
  function test_updateSiloDesignationForChainSelectors_Success() public {
    uint256 amount = 1e18;
    uint64[] memory chainSelectors = new uint64[](1);

    chainSelectors[0] = SILOED_CHAIN_SELECTOR;

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.ChainSiloeDesignationUpdated(SILOED_CHAIN_SELECTOR, true);

    s_siloedLockReleaseTokenPool.updateSiloDesignationForChainSelectors(new uint64[](0), chainSelectors);

    // Assert that the funds are siloed correctly
    assertTrue(s_siloedLockReleaseTokenPool.chainFundsAreSiloed(SILOED_CHAIN_SELECTOR));

    // Provide some Liquidity so that we can then check that it gets removed.
    s_siloedLockReleaseTokenPool.setRebalancer(SILOED_CHAIN_SELECTOR, OWNER);
    s_siloedLockReleaseTokenPool.provideLiquidity(SILOED_CHAIN_SELECTOR, amount);
    assertEq(s_siloedLockReleaseTokenPool.getSiloedTokensByChain(SILOED_CHAIN_SELECTOR), amount);

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.ChainSiloeDesignationUpdated(SILOED_CHAIN_SELECTOR, false);

    assertEq(s_siloedLockReleaseTokenPool.getliquidityForUnsiloedChains(), 0);

    s_siloedLockReleaseTokenPool.updateSiloDesignationForChainSelectors(chainSelectors, new uint64[](0));

    // Check that the locked funds accounting was cleared when the funds were un-siloed.
    assertFalse(s_siloedLockReleaseTokenPool.chainFundsAreSiloed(SILOED_CHAIN_SELECTOR));
    assertEq(s_siloedLockReleaseTokenPool.getSiloedTokensByChain(SILOED_CHAIN_SELECTOR), amount);

    // Assert that the available liquidity moved from being siloed to unsiloed.
    assertEq(s_siloedLockReleaseTokenPool.getliquidityForUnsiloedChains(), amount);
  }
}
