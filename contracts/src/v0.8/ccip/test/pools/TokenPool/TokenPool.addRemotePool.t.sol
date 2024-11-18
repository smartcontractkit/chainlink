// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_addRemotePool is TokenPoolSetup {
  function test_addRemotePool_Success() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;

    address initialPool = makeAddr("remotePool");
    address remoteToken = makeAddr("remoteToken");
    bytes memory remotePool = abi.encode(type(uint256).max);

    bytes32 remotePairHash = keccak256(abi.encode(chainSelector, remotePool));

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      remotePoolAddress: abi.encode(initialPool),
      remoteTokenAddress: abi.encode(remoteToken),
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    vm.expectEmit();
    emit TokenPool.RemotePoolSet(chainSelector, remotePool, remotePairHash);

    vm.startPrank(OWNER);

    s_tokenPool.addRemotePool(chainSelector, remotePool);

    bytes[] memory remotePools = s_tokenPool.getRemotePools(chainSelector);

    assertEq(remotePools.length, 2);
    assertEq(remotePools[0], abi.encode(initialPool));
    assertEq(remotePools[1], remotePool);
  }

  // Reverts

  function test_NonExistentChain_Revert() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;
    bytes memory remotePool = abi.encode(type(uint256).max);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, chainSelector));

    s_tokenPool.addRemotePool(chainSelector, remotePool);
  }

  function test_ZeroAddressNotAllowed_Revert() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;
    bytes memory remotePool = abi.encode(address(0));

    address initialPool = makeAddr("remotePool");
    address remoteToken = makeAddr("remoteToken");

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      remotePoolAddress: abi.encode(initialPool),
      remoteTokenAddress: abi.encode(remoteToken),
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ZeroAddressNotAllowed.selector));

    s_tokenPool.addRemotePool(chainSelector, remotePool);
  }

  function test_ZeroLengthAddressNotAllowed_Revert() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;
    bytes memory remotePool = "";
    address initialPool = makeAddr("remotePool");
    address remoteToken = makeAddr("remoteToken");

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      remotePoolAddress: abi.encode(initialPool),
      remoteTokenAddress: abi.encode(remoteToken),
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ZeroAddressNotAllowed.selector));

    s_tokenPool.addRemotePool(chainSelector, remotePool);
  }

  function test_PoolAlreadyAdded_Revert() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;

    address initialPool = makeAddr("remotePool");
    address remoteToken = makeAddr("remoteToken");
    bytes memory remotePool = abi.encode(type(uint256).max);

    bytes32 remotePairHash = keccak256(abi.encode(chainSelector, remotePool));

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      remotePoolAddress: abi.encode(initialPool),
      remoteTokenAddress: abi.encode(remoteToken),
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    vm.expectEmit();
    emit TokenPool.RemotePoolSet(chainSelector, remotePool, remotePairHash);

    s_tokenPool.addRemotePool(chainSelector, remotePool);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.PoolAlreadyAdded.selector, chainSelector, remotePool));

    // Attempt to add the same pool again but revert
    s_tokenPool.addRemotePool(chainSelector, remotePool);
  }
}
