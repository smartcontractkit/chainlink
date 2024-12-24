// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";

contract SiloedLockReleaseTokenPool_transferLiquidity is SiloedLockReleaseTokenPoolSetup {
  function test_transferLiquidity_Success() public {
    uint256 amount = 10e24;

    // Set rebalancer and provide some liquidity that can be transferred
    s_siloedLockReleaseTokenPool.setRebalancer(SILOED_CHAIN_SELECTOR, OWNER);
    s_siloedLockReleaseTokenPool.provideLiquidity(SILOED_CHAIN_SELECTOR, amount);

    SiloedLockReleaseTokenPool newTokenPool = new SiloedLockReleaseTokenPool(
      s_token, DEFAULT_TOKEN_DECIMALS, new address[](0), address(s_mockRMNRemote), address(s_sourceRouter)
    );

    // Mark the chain on the new token pool as siloed.
    uint64[] memory chainSelectors = new uint64[](1);
    chainSelectors[0] = SILOED_CHAIN_SELECTOR;
    newTokenPool.updateChainSelectorMechanisms(new uint64[](0), chainSelectors);

    // Begin transferring ownership
    s_siloedLockReleaseTokenPool.transferOwnership(address(newTokenPool));

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.LiquidityTransferred(
      SILOED_CHAIN_SELECTOR, address(s_siloedLockReleaseTokenPool), amount
    );

    newTokenPool.transferLiquidity(SILOED_CHAIN_SELECTOR, address(s_siloedLockReleaseTokenPool), amount);

    assertEq(s_siloedLockReleaseTokenPool.owner(), address(newTokenPool));
    assertEq(newTokenPool.getLockedTokensByChain(SILOED_CHAIN_SELECTOR), amount);
    assertEq(s_siloedLockReleaseTokenPool.getLockedTokensByChain(SILOED_CHAIN_SELECTOR), 0);
    assertEq(s_token.balanceOf(address(newTokenPool)), amount);
    assertEq(s_token.balanceOf(address(s_siloedLockReleaseTokenPool)), 0);
  }

  function test_transferLiquidity_MultipleTransfers_Success() public {
    uint256 amount = 10e24;

    // Set rebalancer and provide some liquidity that can be transferred
    s_siloedLockReleaseTokenPool.setRebalancer(SILOED_CHAIN_SELECTOR, OWNER);
    s_siloedLockReleaseTokenPool.provideLiquidity(SILOED_CHAIN_SELECTOR, amount);

    SiloedLockReleaseTokenPool newTokenPool = new SiloedLockReleaseTokenPool(
      s_token, DEFAULT_TOKEN_DECIMALS, new address[](0), address(s_mockRMNRemote), address(s_sourceRouter)
    );

    // Mark the chain on the new token pool as siloed.
    uint64[] memory chainSelectors = new uint64[](1);
    chainSelectors[0] = SILOED_CHAIN_SELECTOR;
    newTokenPool.updateChainSelectorMechanisms(new uint64[](0), chainSelectors);

    // Begin transferring ownership
    s_siloedLockReleaseTokenPool.transferOwnership(address(newTokenPool));

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.LiquidityTransferred(
      SILOED_CHAIN_SELECTOR, address(s_siloedLockReleaseTokenPool), amount / 2
    );

    // Do the first transfer which should accept ownership
    newTokenPool.transferLiquidity(SILOED_CHAIN_SELECTOR, address(s_siloedLockReleaseTokenPool), amount / 2);

    assertEq(s_siloedLockReleaseTokenPool.owner(), address(newTokenPool));
    assertEq(newTokenPool.getLockedTokensByChain(SILOED_CHAIN_SELECTOR), amount / 2);
    assertEq(s_siloedLockReleaseTokenPool.getLockedTokensByChain(SILOED_CHAIN_SELECTOR), amount / 2);
    assertEq(s_token.balanceOf(address(newTokenPool)), amount / 2);
    assertEq(s_token.balanceOf(address(s_siloedLockReleaseTokenPool)), amount / 2);

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.LiquidityTransferred(
      SILOED_CHAIN_SELECTOR, address(s_siloedLockReleaseTokenPool), amount / 2
    );

    newTokenPool.transferLiquidity(SILOED_CHAIN_SELECTOR, address(s_siloedLockReleaseTokenPool), amount / 2);

    assertEq(s_siloedLockReleaseTokenPool.owner(), address(newTokenPool));
    assertEq(newTokenPool.getLockedTokensByChain(SILOED_CHAIN_SELECTOR), amount);
    assertEq(s_siloedLockReleaseTokenPool.getLockedTokensByChain(SILOED_CHAIN_SELECTOR), 0);
    assertEq(s_token.balanceOf(address(newTokenPool)), amount);
    assertEq(s_token.balanceOf(address(s_siloedLockReleaseTokenPool)), 0);
  }
}
