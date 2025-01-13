// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";

contract SiloedLockReleaseTokenPool_provideLiqudity is SiloedLockReleaseTokenPoolSetup {
  address public UNAUTHORIZED_ADDRESS = address(0xdeadbeef);

  function setUp() public override {
    super.setUp();

    s_siloedLockReleaseTokenPool.setSiloedChainRebalancer(SILOED_CHAIN_SELECTOR, OWNER);
  }

  function test_ProvideLiquidity_ChainNotSiloed() public {
    uint256 amount = 1e24;

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.LiquidityAdded(DEST_CHAIN_SELECTOR, OWNER, amount);

    s_siloedLockReleaseTokenPool.provideLiquidity(DEST_CHAIN_SELECTOR, amount);

    assertEq(s_token.balanceOf(address(s_siloedLockReleaseTokenPool)), amount);

    // Since the funds for the destination chain are not siloed,
    // the locked token amount should not be increased
    assertEq(s_siloedLockReleaseTokenPool.getSiloedTokensByChain(DEST_CHAIN_SELECTOR), amount);
  }

  function test_ProvideLiquidity_ChainSiloed() public {
    uint256 amount = 1e24;

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.LiquidityAdded(SILOED_CHAIN_SELECTOR, OWNER, amount);

    s_siloedLockReleaseTokenPool.provideLiquidity(SILOED_CHAIN_SELECTOR, amount);

    assertEq(s_token.balanceOf(address(s_siloedLockReleaseTokenPool)), amount);

    // Since the funds for the destination chain are not siloed,
    // the locked token amount should not be increased
    assertEq(s_siloedLockReleaseTokenPool.getSiloedTokensByChain(SILOED_CHAIN_SELECTOR), amount);
  }

  // Reverts

  function test_ProvideLiquidity_RevertWhen_UnauthorizedForSiloedChain() public {
    vm.startPrank(UNAUTHORIZED_ADDRESS);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, UNAUTHORIZED_ADDRESS));

    s_siloedLockReleaseTokenPool.provideLiquidity(SILOED_CHAIN_SELECTOR, 1);
  }

  function test_ProvideLiquidity_RevertWhen_UnauthorizedForNotSiloedChain() public {
    vm.startPrank(UNAUTHORIZED_ADDRESS);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, UNAUTHORIZED_ADDRESS));

    s_siloedLockReleaseTokenPool.provideLiquidity(DEST_CHAIN_SELECTOR, 1);
  }
}
