// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Router} from "../../../Router.sol";
import {Pool} from "../../../libraries/Pool.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_addRemotePool is TokenPoolSetup {
  function test_addRemotePool() public {
    // Use a longer data type to ensure it also works for non-evm
    bytes memory remotePool = abi.encode(makeAddr("non-evm-1"), makeAddr("non-evm-2"));

    vm.startPrank(OWNER);

    vm.expectEmit();
    emit TokenPool.RemotePoolAdded(DEST_CHAIN_SELECTOR, remotePool);

    s_tokenPool.addRemotePool(DEST_CHAIN_SELECTOR, remotePool);

    bytes[] memory remotePools = s_tokenPool.getRemotePools(DEST_CHAIN_SELECTOR);

    assertEq(remotePools.length, 2);
    assertEq(remotePools[0], abi.encode(s_initialRemotePool));
    assertEq(remotePools[1], remotePool);
  }

  function test_addRemotePool_MultipleActive() public {
    bytes[] memory remotePools = new bytes[](3);
    remotePools[0] = abi.encode(makeAddr("remotePool1"));
    remotePools[1] = abi.encode(makeAddr("remotePool2"));
    remotePools[2] = abi.encode(makeAddr("remotePool3"));

    address fakeOffRamp = makeAddr("fakeOffRamp");

    vm.mockCall(
      address(s_sourceRouter), abi.encodeCall(Router.isOffRamp, (DEST_CHAIN_SELECTOR, fakeOffRamp)), abi.encode(true)
    );

    vm.startPrank(fakeOffRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.InvalidSourcePoolAddress.selector, remotePools[0]));
    s_tokenPool.releaseOrMint(_getReleaseOrMintInV1(remotePools[0]));

    // There's already one pool setup through the test setup
    assertEq(s_tokenPool.getRemotePools(DEST_CHAIN_SELECTOR).length, 1);

    vm.startPrank(OWNER);
    s_tokenPool.addRemotePool(DEST_CHAIN_SELECTOR, remotePools[0]);

    vm.startPrank(fakeOffRamp);
    s_tokenPool.releaseOrMint(_getReleaseOrMintInV1(remotePools[0]));

    // Adding an additional pool does not remove the previous one
    vm.startPrank(OWNER);
    s_tokenPool.addRemotePool(DEST_CHAIN_SELECTOR, remotePools[1]);

    // Both should now work
    assertEq(s_tokenPool.getRemotePools(DEST_CHAIN_SELECTOR).length, 3);
    vm.startPrank(fakeOffRamp);
    s_tokenPool.releaseOrMint(_getReleaseOrMintInV1(remotePools[0]));
    s_tokenPool.releaseOrMint(_getReleaseOrMintInV1(remotePools[1]));

    // Adding a third pool, and removing the first one
    vm.startPrank(OWNER);
    s_tokenPool.addRemotePool(DEST_CHAIN_SELECTOR, remotePools[2]);
    assertEq(s_tokenPool.getRemotePools(DEST_CHAIN_SELECTOR).length, 4);
    s_tokenPool.removeRemotePool(DEST_CHAIN_SELECTOR, remotePools[0]);
    assertEq(s_tokenPool.getRemotePools(DEST_CHAIN_SELECTOR).length, 3);

    // Only the last two should work
    vm.startPrank(fakeOffRamp);
    vm.expectRevert(abi.encodeWithSelector(TokenPool.InvalidSourcePoolAddress.selector, remotePools[0]));
    s_tokenPool.releaseOrMint(_getReleaseOrMintInV1(remotePools[0]));
    s_tokenPool.releaseOrMint(_getReleaseOrMintInV1(remotePools[1]));
    s_tokenPool.releaseOrMint(_getReleaseOrMintInV1(remotePools[2]));

    // Removing the chain removes all associated pool hashes
    vm.startPrank(OWNER);

    uint64[] memory chainSelectorsToRemove = new uint64[](1);
    chainSelectorsToRemove[0] = DEST_CHAIN_SELECTOR;
    s_tokenPool.applyChainUpdates(chainSelectorsToRemove, new TokenPool.ChainUpdate[](0));

    assertEq(s_tokenPool.getRemotePools(DEST_CHAIN_SELECTOR).length, 0);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, DEST_CHAIN_SELECTOR));
    s_tokenPool.releaseOrMint(_getReleaseOrMintInV1(remotePools[0]));
    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, DEST_CHAIN_SELECTOR));
    s_tokenPool.releaseOrMint(_getReleaseOrMintInV1(remotePools[1]));
    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, DEST_CHAIN_SELECTOR));
    s_tokenPool.releaseOrMint(_getReleaseOrMintInV1(remotePools[2]));

    // Adding the chain back should NOT allow the previous pools to work again
    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddresses: new bytes[](0),
      remoteTokenAddress: abi.encode(s_initialRemoteToken),
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });

    vm.startPrank(OWNER);
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdate);

    vm.startPrank(fakeOffRamp);
    vm.expectRevert(abi.encodeWithSelector(TokenPool.InvalidSourcePoolAddress.selector, remotePools[0]));
    s_tokenPool.releaseOrMint(_getReleaseOrMintInV1(remotePools[0]));
    vm.expectRevert(abi.encodeWithSelector(TokenPool.InvalidSourcePoolAddress.selector, remotePools[1]));
    s_tokenPool.releaseOrMint(_getReleaseOrMintInV1(remotePools[1]));
    vm.expectRevert(abi.encodeWithSelector(TokenPool.InvalidSourcePoolAddress.selector, remotePools[2]));
    s_tokenPool.releaseOrMint(_getReleaseOrMintInV1(remotePools[2]));
  }

  function _getReleaseOrMintInV1(
    bytes memory sourcePoolAddress
  ) internal view returns (Pool.ReleaseOrMintInV1 memory) {
    return Pool.ReleaseOrMintInV1({
      originalSender: abi.encode(OWNER),
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      receiver: OWNER,
      amount: 1000,
      localToken: address(s_token),
      sourcePoolAddress: sourcePoolAddress,
      sourcePoolData: "",
      offchainTokenData: ""
    });
  }

  // Reverts

  function test_RevertWhen_NonExistentChain() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR + 1;
    bytes memory remotePool = abi.encode(type(uint256).max);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, chainSelector));

    s_tokenPool.addRemotePool(chainSelector, remotePool);
  }

  function test_RevertWhen_ZeroLengthAddressNotAllowed() public {
    bytes memory remotePool = "";

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ZeroAddressNotAllowed.selector));

    s_tokenPool.addRemotePool(DEST_CHAIN_SELECTOR, remotePool);
  }

  function test_RevertWhen_PoolAlreadyAdded() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;

    bytes memory remotePool = abi.encode(type(uint256).max);

    vm.expectEmit();
    emit TokenPool.RemotePoolAdded(chainSelector, remotePool);

    s_tokenPool.addRemotePool(chainSelector, remotePool);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.PoolAlreadyAdded.selector, chainSelector, remotePool));

    // Attempt to add the same pool again but revert
    s_tokenPool.addRemotePool(chainSelector, remotePool);
  }
}
