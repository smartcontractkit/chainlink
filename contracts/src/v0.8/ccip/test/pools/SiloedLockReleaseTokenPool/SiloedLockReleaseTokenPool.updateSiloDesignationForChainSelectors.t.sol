// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_updateSiloDesignationForChainSelectors is SiloedLockReleaseTokenPoolSetup {
  function test_updateSiloDesignationForChainSelectors() public {
    uint256 amount = 1e18;

    SiloedLockReleaseTokenPool.ChainSiloConfigUpdate[] memory chainSelectors =
      new SiloedLockReleaseTokenPool.ChainSiloConfigUpdate[](1);

    chainSelectors[0] =
      SiloedLockReleaseTokenPool.ChainSiloConfigUpdate({remoteChainSelector: SILOED_CHAIN_SELECTOR, rebalancer: OWNER});

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.ChainSiloed(SILOED_CHAIN_SELECTOR);

    s_siloedLockReleaseTokenPool.updateSiloDesignationForChainSelectors(new uint64[](0), chainSelectors);

    // Assert that the funds are siloed correctly
    assertTrue(s_siloedLockReleaseTokenPool.chainFundsAreSiloed(SILOED_CHAIN_SELECTOR));
    assertEq(s_siloedLockReleaseTokenPool.getSiloedTokensByChain(SILOED_CHAIN_SELECTOR), 0);
    assertEq(s_siloedLockReleaseTokenPool.getRebalancerByChain(SILOED_CHAIN_SELECTOR), OWNER);

    // Provide some Liquidity so that we can then check that it gets removed.
    s_siloedLockReleaseTokenPool.setSiloedChainRebalancer(SILOED_CHAIN_SELECTOR, OWNER);
    s_siloedLockReleaseTokenPool.provideSiloedLiquidity(SILOED_CHAIN_SELECTOR, amount);
    assertEq(s_siloedLockReleaseTokenPool.getSiloedTokensByChain(SILOED_CHAIN_SELECTOR), amount);

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.ChainUnsiloed(SILOED_CHAIN_SELECTOR, amount);

    assertEq(s_siloedLockReleaseTokenPool.getliquidityForUnsiloedChains(), 0);

    uint64[] memory removableChainSelectors = new uint64[](1);
    removableChainSelectors[0] = SILOED_CHAIN_SELECTOR;

    s_siloedLockReleaseTokenPool.updateSiloDesignationForChainSelectors(
      removableChainSelectors, new SiloedLockReleaseTokenPool.ChainSiloConfigUpdate[](0)
    );

    // Check that the locked funds accounting was cleared when the funds were un-siloed.
    assertFalse(s_siloedLockReleaseTokenPool.chainFundsAreSiloed(SILOED_CHAIN_SELECTOR));
    assertEq(s_siloedLockReleaseTokenPool.getSiloedTokensByChain(SILOED_CHAIN_SELECTOR), amount);

    // Assert that the available liquidity moved from being siloed to unsiloed.
    assertEq(s_siloedLockReleaseTokenPool.getliquidityForUnsiloedChains(), amount);
  }
}
