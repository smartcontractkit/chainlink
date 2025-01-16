// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_setRebalancer is SiloedLockReleaseTokenPoolSetup {
  address public constant REBALANCER_ADDRESS = address(0xdeadbeef);

  function test_setSiloRebalancer() public {
    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.SiloRebalancerSet(SILOED_CHAIN_SELECTOR, REBALANCER_ADDRESS, OWNER);

    s_siloedLockReleaseTokenPool.setSiloRebalancer(SILOED_CHAIN_SELECTOR, REBALANCER_ADDRESS);

    assertEq(s_siloedLockReleaseTokenPool.getSiloRebalancer(SILOED_CHAIN_SELECTOR), REBALANCER_ADDRESS);
    assertEq(s_siloedLockReleaseTokenPool.getSiloRebalancer(DEST_CHAIN_SELECTOR), OWNER);
  }

  function test_setRebalancer_UnsiloedChains() public {
    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.UnsiloedRebalancerSet(REBALANCER_ADDRESS, OWNER);

    s_siloedLockReleaseTokenPool.setRebalancer(REBALANCER_ADDRESS);

    assertEq(s_siloedLockReleaseTokenPool.getSiloRebalancer(DEST_CHAIN_SELECTOR), REBALANCER_ADDRESS);
  }

  // Reverts

  function test_setSiloRebalancer_RevertWhen_ChainNotSiloed() public {
    vm.expectRevert(abi.encodeWithSelector(SiloedLockReleaseTokenPool.ChainNotSiloed.selector, DEST_CHAIN_SELECTOR));

    s_siloedLockReleaseTokenPool.setSiloRebalancer(DEST_CHAIN_SELECTOR, REBALANCER_ADDRESS);
  }
}
