// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_setRebalancer is SiloedLockReleaseTokenPoolSetup {
  address public constant REBALANCER_ADDRESS = address(0xdeadbeef);

  function test_setSiloRebalancer() public {
    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.SiloRebalancerSet(SILOED_CHAIN_SELECTOR, OWNER, REBALANCER_ADDRESS);

    s_siloedLockReleaseTokenPool.setSiloRebalancer(SILOED_CHAIN_SELECTOR, REBALANCER_ADDRESS);

    assertEq(s_siloedLockReleaseTokenPool.getChainRebalancer(SILOED_CHAIN_SELECTOR), REBALANCER_ADDRESS);
    assertEq(s_siloedLockReleaseTokenPool.getChainRebalancer(DEST_CHAIN_SELECTOR), OWNER);
  }

  function test_setRebalancer_UnsiloedChains() public {
    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.UnsiloedRebalancerSet(OWNER, REBALANCER_ADDRESS);

    s_siloedLockReleaseTokenPool.setRebalancer(REBALANCER_ADDRESS);

    assertEq(s_siloedLockReleaseTokenPool.getChainRebalancer(DEST_CHAIN_SELECTOR), REBALANCER_ADDRESS);
    assertEq(s_siloedLockReleaseTokenPool.getRebalancer(), REBALANCER_ADDRESS);
  }

  // Reverts

  function test_setSiloRebalancer_RevertWhen_ChainNotSiloed() public {
    vm.expectRevert(abi.encodeWithSelector(SiloedLockReleaseTokenPool.ChainNotSiloed.selector, DEST_CHAIN_SELECTOR));

    s_siloedLockReleaseTokenPool.setSiloRebalancer(DEST_CHAIN_SELECTOR, REBALANCER_ADDRESS);
  }

  function test_SetSiloRebalancer_RevertWhen_InvalidZeroAddress() public {
    vm.expectRevert(abi.encodeWithSelector(TokenPool.ZeroAddressNotAllowed.selector));

    s_siloedLockReleaseTokenPool.setSiloRebalancer(SILOED_CHAIN_SELECTOR, address(0));
  }

  function test_SetRebalancer_RevertWhen_InvalidZeroAddress() public {
    vm.expectRevert(abi.encodeWithSelector(TokenPool.ZeroAddressNotAllowed.selector));

    s_siloedLockReleaseTokenPool.setRebalancer(address(0));
  }
}
