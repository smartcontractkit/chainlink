// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_removeRemotePool is TokenPoolSetup {
  function test_removeRemotePool() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;
    // Use a longer data type to ensure it also works for non-evm
    bytes memory remotePool = abi.encode(makeAddr("non-evm-1"), makeAddr("non-evm-2"));

    vm.expectEmit();
    emit TokenPool.RemotePoolAdded(chainSelector, remotePool);

    // Add the remote pool properly so that it can be removed
    s_tokenPool.addRemotePool(chainSelector, remotePool);

    bytes[] memory remotePools = s_tokenPool.getRemotePools(chainSelector);
    assertEq(remotePools.length, 2);
    assertEq(remotePools[0], abi.encode(s_initialRemotePool));
    assertEq(remotePools[1], remotePool);

    vm.expectEmit();
    emit TokenPool.RemotePoolRemoved(chainSelector, remotePool);

    s_tokenPool.removeRemotePool(chainSelector, remotePool);

    remotePools = s_tokenPool.getRemotePools(chainSelector);
    assertEq(remotePools.length, 1);
    assertEq(remotePools[0], abi.encode(s_initialRemotePool));

    // Assert that it can be added after it has been removed
    s_tokenPool.addRemotePool(chainSelector, remotePool);

    remotePools = s_tokenPool.getRemotePools(chainSelector);
    assertEq(remotePools.length, 2);
    assertEq(remotePools[0], abi.encode(s_initialRemotePool));
    assertEq(remotePools[1], remotePool);
  }

  // Reverts

  function test_RevertWhen_NonExistentChain() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR + 1;
    bytes memory remotePool = abi.encode(type(uint256).max);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, chainSelector));

    s_tokenPool.removeRemotePool(chainSelector, remotePool);
  }

  function test_RevertWhen_InvalidRemotePoolForChain() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;
    bytes memory remotePool = abi.encode(type(uint256).max);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.InvalidRemotePoolForChain.selector, chainSelector, remotePool));

    s_tokenPool.removeRemotePool(chainSelector, remotePool);
  }
}
