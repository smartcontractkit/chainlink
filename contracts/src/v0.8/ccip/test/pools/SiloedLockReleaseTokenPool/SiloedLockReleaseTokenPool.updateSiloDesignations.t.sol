// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_updateSiloDesignations is SiloedLockReleaseTokenPoolSetup {
  function test_updateSiloDesignations() public {
    uint256 amount = 1e18;

    // Provide some Liquidity so that we can then check that it gets removed.
    s_siloedLockReleaseTokenPool.provideSiloedLiquidity(SILOED_CHAIN_SELECTOR, amount);
    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(SILOED_CHAIN_SELECTOR), amount);

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.ChainUnsiloed(SILOED_CHAIN_SELECTOR, amount);

    assertEq(s_siloedLockReleaseTokenPool.getUnsiloedLiquidity(), 0);

    uint64[] memory removableChainSelectors = new uint64[](1);
    removableChainSelectors[0] = SILOED_CHAIN_SELECTOR;

    s_siloedLockReleaseTokenPool.updateSiloDesignations(
      removableChainSelectors, new SiloedLockReleaseTokenPool.SiloConfigUpdate[](0)
    );

    // Check that the locked funds accounting was cleared when the funds were un-siloed.
    assertFalse(s_siloedLockReleaseTokenPool.isSiloed(SILOED_CHAIN_SELECTOR));
    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(SILOED_CHAIN_SELECTOR), amount);

    // Assert that the available liquidity moved from being siloed to unsiloed.
    assertEq(s_siloedLockReleaseTokenPool.getUnsiloedLiquidity(), amount);

    // Now we re-silo the chain
    SiloedLockReleaseTokenPool.SiloConfigUpdate[] memory chainSelectors =
      new SiloedLockReleaseTokenPool.SiloConfigUpdate[](1);

    chainSelectors[0] =
      SiloedLockReleaseTokenPool.SiloConfigUpdate({remoteChainSelector: SILOED_CHAIN_SELECTOR, rebalancer: OWNER});

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.ChainSiloed(SILOED_CHAIN_SELECTOR, OWNER);

    s_siloedLockReleaseTokenPool.updateSiloDesignations(new uint64[](0), chainSelectors);

    // Assert that the funds are siloed correctly
    assertTrue(s_siloedLockReleaseTokenPool.isSiloed(SILOED_CHAIN_SELECTOR));
    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(SILOED_CHAIN_SELECTOR), 0);
    assertEq(s_siloedLockReleaseTokenPool.getChainRebalancer(SILOED_CHAIN_SELECTOR), OWNER);

    // Provide some Liquidity so that we can then check that it gets removed.
    s_siloedLockReleaseTokenPool.provideSiloedLiquidity(SILOED_CHAIN_SELECTOR, amount);
    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(SILOED_CHAIN_SELECTOR), amount);
  }

  // Reverts

  function test_updateSiloDesignations_RevertWhen_ChainNotSiloed() public {
    uint64[] memory removableChainSelectors = new uint64[](1);
    removableChainSelectors[0] = DEST_CHAIN_SELECTOR;

    vm.expectRevert(abi.encodeWithSelector(SiloedLockReleaseTokenPool.ChainNotSiloed.selector, DEST_CHAIN_SELECTOR));

    s_siloedLockReleaseTokenPool.updateSiloDesignations(
      removableChainSelectors, new SiloedLockReleaseTokenPool.SiloConfigUpdate[](0)
    );
  }

  function test_updateSiloDesignations_RevertWhen_InvalidChainSelector_ChainSelectorZero() public {
    SiloedLockReleaseTokenPool.SiloConfigUpdate[] memory adds = new SiloedLockReleaseTokenPool.SiloConfigUpdate[](1);
    adds[0] = SiloedLockReleaseTokenPool.SiloConfigUpdate({remoteChainSelector: 0, rebalancer: OWNER});

    // Chain selector cannot be zero
    vm.expectRevert(abi.encodeWithSelector(SiloedLockReleaseTokenPool.InvalidChainSelector.selector, 0));

    s_siloedLockReleaseTokenPool.updateSiloDesignations(new uint64[](0), adds);
  }

  function test_updateSiloDesignations_RevertWhen_InvalidChainSelector_ChainAlreadySiloed() public {
    SiloedLockReleaseTokenPool.SiloConfigUpdate[] memory adds = new SiloedLockReleaseTokenPool.SiloConfigUpdate[](1);
    adds[0] =
      SiloedLockReleaseTokenPool.SiloConfigUpdate({remoteChainSelector: SILOED_CHAIN_SELECTOR, rebalancer: OWNER});

    // Since the chain is already siloed you cannot re-silo it.
    vm.expectRevert(
      abi.encodeWithSelector(SiloedLockReleaseTokenPool.InvalidChainSelector.selector, SILOED_CHAIN_SELECTOR)
    );

    s_siloedLockReleaseTokenPool.updateSiloDesignations(new uint64[](0), adds);
  }

  function test_updateSiloDesignations_RevertWhen_InvalidZeroRebalancerAddress() public {
    SiloedLockReleaseTokenPool.SiloConfigUpdate[] memory adds = new SiloedLockReleaseTokenPool.SiloConfigUpdate[](1);
    adds[0] =
      SiloedLockReleaseTokenPool.SiloConfigUpdate({remoteChainSelector: DEST_CHAIN_SELECTOR, rebalancer: address(0)});

    // Rebalancer address cannot be zero
    vm.expectRevert(abi.encodeWithSelector(TokenPool.ZeroAddressNotAllowed.selector));

    s_siloedLockReleaseTokenPool.updateSiloDesignations(new uint64[](0), adds);
  }
}
