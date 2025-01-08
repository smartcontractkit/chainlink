// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_setRebalancer is SiloedLockReleaseTokenPoolSetup {
  address public REBALANCER_ADDRESS = address(0xdeadbeef);

  function test_setRebalancer_Success() public {
    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.RebalancerSet(SILOED_CHAIN_SELECTOR, REBALANCER_ADDRESS, OWNER);

    s_siloedLockReleaseTokenPool.setRebalancer(SILOED_CHAIN_SELECTOR, REBALANCER_ADDRESS);

    assertEq(s_siloedLockReleaseTokenPool.getRebalancerByChain(SILOED_CHAIN_SELECTOR), REBALANCER_ADDRESS);
    assertEq(s_siloedLockReleaseTokenPool.getRebalancerByChain(DEST_CHAIN_SELECTOR), OWNER);
  }

  // Reverts

  function test_setRebalancer_RevertWhen_ChainNotSiloed() public {
    vm.expectRevert(abi.encodeWithSelector(SiloedLockReleaseTokenPool.ChainNotSiloed.selector, DEST_CHAIN_SELECTOR));

    s_siloedLockReleaseTokenPool.setRebalancer(DEST_CHAIN_SELECTOR, REBALANCER_ADDRESS);
  }
}
