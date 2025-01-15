// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Router} from "../../../Router.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_onlyOnRamp is TokenPoolSetup {
  function test_onlyOnRamp() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;
    address onRamp = makeAddr("onRamp");

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: chainSelector, onRamp: onRamp});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), new Router.OffRamp[](0));

    vm.startPrank(onRamp);

    s_tokenPool.onlyOnRampModifier(chainSelector);
  }

  function test_RevertWhen_ChainNotAllowed() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR + 1;
    address onRamp = makeAddr("onRamp");

    vm.startPrank(onRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, chainSelector));
    s_tokenPool.onlyOnRampModifier(chainSelector);

    vm.startPrank(OWNER);

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      remotePoolAddresses: new bytes[](0),
      remoteTokenAddress: abi.encode(address(2)),
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdate);

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: chainSelector, onRamp: onRamp});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), new Router.OffRamp[](0));

    vm.startPrank(onRamp);
    // Should succeed now that we've added the chain
    s_tokenPool.onlyOnRampModifier(chainSelector);

    uint64[] memory chainsToRemove = new uint64[](1);
    chainsToRemove[0] = chainSelector;

    vm.startPrank(OWNER);
    s_tokenPool.applyChainUpdates(chainsToRemove, new TokenPool.ChainUpdate[](0));

    vm.startPrank(onRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, chainSelector));
    s_tokenPool.onlyOffRampModifier(chainSelector);
  }

  function test_RevertWhen_CallerIsNotARampOnRouter() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;
    address onRamp = makeAddr("onRamp");

    vm.startPrank(onRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.CallerIsNotARampOnRouter.selector, onRamp));

    s_tokenPool.onlyOnRampModifier(chainSelector);
  }
}
