// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_setRebalancer is SiloedLockReleaseTokenPoolSetup {
  address public REBALANCER_ADDRESS = address(0xdeadbeef);

  function test_setSiloedChainRebalancer() public {
    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.SiloedChainRebalancerSet(SILOED_CHAIN_SELECTOR, REBALANCER_ADDRESS, OWNER);

    s_siloedLockReleaseTokenPool.setSiloedChainRebalancer(SILOED_CHAIN_SELECTOR, REBALANCER_ADDRESS);

    assertEq(s_siloedLockReleaseTokenPool.getRebalancerByChain(SILOED_CHAIN_SELECTOR), REBALANCER_ADDRESS);
    assertEq(s_siloedLockReleaseTokenPool.getRebalancerByChain(DEST_CHAIN_SELECTOR), OWNER);
  }

  function test_setUnsiloedChainRebalancer() public {
    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.UnsiloedChainRebalancerSet(REBALANCER_ADDRESS, OWNER);

    s_siloedLockReleaseTokenPool.setUnsiloedChainRebalancer(REBALANCER_ADDRESS);

    assertEq(s_siloedLockReleaseTokenPool.getRebalancerByChain(DEST_CHAIN_SELECTOR), REBALANCER_ADDRESS);
  }

  // Reverts

  function test_setSiloedChainRebalancer_RevertWhen_ChainNotSiloed() public {
    vm.expectRevert(abi.encodeWithSelector(SiloedLockReleaseTokenPool.ChainNotSiloed.selector, DEST_CHAIN_SELECTOR));

    s_siloedLockReleaseTokenPool.setSiloedChainRebalancer(DEST_CHAIN_SELECTOR, REBALANCER_ADDRESS);
  }
}
