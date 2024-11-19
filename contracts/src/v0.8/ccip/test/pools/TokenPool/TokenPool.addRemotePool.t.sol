// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_addRemotePool is TokenPoolSetup {
  function test_addRemotePool_Success() public {
    // Use a longer data type to ensure it also works for non-evm
    bytes memory remotePool = abi.encode(makeAddr("non-evm-1"), makeAddr("non-evm-2"));

    bytes32 remotePairHash = keccak256(abi.encode(DEST_CHAIN_SELECTOR, remotePool));

    vm.startPrank(OWNER);

    vm.expectEmit();
    emit TokenPool.RemotePoolSet(DEST_CHAIN_SELECTOR, remotePool, remotePairHash);

    s_tokenPool.addRemotePool(DEST_CHAIN_SELECTOR, remotePool);

    bytes[] memory remotePools = s_tokenPool.getRemotePools(DEST_CHAIN_SELECTOR);

    assertEq(remotePools.length, 2);
    assertEq(remotePools[0], abi.encode(s_initialRemotePool));
    assertEq(remotePools[1], remotePool);
  }

  function test_addRemotePool_MultipleActive() public {}

  // Reverts

  function test_NonExistentChain_Revert() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR + 1;
    bytes memory remotePool = abi.encode(type(uint256).max);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, chainSelector));

    s_tokenPool.addRemotePool(chainSelector, remotePool);
  }

  function test_ZeroAddressNotAllowed_Revert() public {
    bytes memory remotePool = abi.encode(address(0));

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ZeroAddressNotAllowed.selector));

    s_tokenPool.addRemotePool(DEST_CHAIN_SELECTOR, remotePool);
  }

  function test_ZeroLengthAddressNotAllowed_Revert() public {
    bytes memory remotePool = "";

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ZeroAddressNotAllowed.selector));

    s_tokenPool.addRemotePool(DEST_CHAIN_SELECTOR, remotePool);
  }

  function test_PoolAlreadyAdded_Revert() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;

    bytes memory remotePool = abi.encode(type(uint256).max);

    bytes32 remotePairHash = keccak256(abi.encode(chainSelector, remotePool));

    vm.expectEmit();
    emit TokenPool.RemotePoolSet(chainSelector, remotePool, remotePairHash);

    s_tokenPool.addRemotePool(chainSelector, remotePool);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.PoolAlreadyAdded.selector, chainSelector, remotePool));

    // Attempt to add the same pool again but revert
    s_tokenPool.addRemotePool(chainSelector, remotePool);
  }
}
