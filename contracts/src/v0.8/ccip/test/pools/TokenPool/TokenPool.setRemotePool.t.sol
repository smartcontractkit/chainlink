// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_setRemotePool is TokenPoolSetup {
  function test_setRemotePool_Success() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;
    address initialPool = makeAddr("remotePool");
    address remoteToken = makeAddr("remoteToken");
    // The new pool is a non-evm pool, as it doesn't fit in the normal 160 bits
    bytes memory newPool = abi.encode(type(uint256).max);

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      remotePoolAddress: abi.encode(initialPool),
      remoteTokenAddress: abi.encode(remoteToken),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(chainUpdates);

    vm.expectEmit();
    emit TokenPool.RemotePoolSet(chainSelector, abi.encode(initialPool), newPool);

    s_tokenPool.setRemotePool(chainSelector, newPool);

    assertEq(keccak256(newPool), keccak256(s_tokenPool.getRemotePool(chainSelector)));
  }

  // Reverts

  function test_setRemotePool_NonExistentChain_Reverts() public {
    uint64 chainSelector = 123124;
    bytes memory remotePool = abi.encode(makeAddr("remotePool"));

    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, chainSelector));
    s_tokenPool.setRemotePool(chainSelector, remotePool);
  }

  function test_setRemotePool_OnlyOwner_Reverts() public {
    vm.startPrank(STRANGER);

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_tokenPool.setRemotePool(123124, abi.encode(makeAddr("remotePool")));
  }
}
