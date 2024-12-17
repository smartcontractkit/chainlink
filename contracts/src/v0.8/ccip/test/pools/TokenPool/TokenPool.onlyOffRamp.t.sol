// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Router} from "../../../Router.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_onlyOffRamp is TokenPoolSetup {
  function test_onlyOffRamp() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;
    address offRamp = makeAddr("onRamp");

    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: chainSelector, offRamp: offRamp});
    s_sourceRouter.applyRampUpdates(new Router.OnRamp[](0), new Router.OffRamp[](0), offRampUpdates);

    vm.startPrank(offRamp);

    s_tokenPool.onlyOffRampModifier(chainSelector);
  }

  function test_RevertWhen_ChainNotAllowed() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR + 1;
    address offRamp = makeAddr("onRamp");

    vm.startPrank(offRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, chainSelector));
    s_tokenPool.onlyOffRampModifier(chainSelector);

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

    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: chainSelector, offRamp: offRamp});
    s_sourceRouter.applyRampUpdates(new Router.OnRamp[](0), new Router.OffRamp[](0), offRampUpdates);

    vm.startPrank(offRamp);
    // Should succeed now that we've added the chain
    s_tokenPool.onlyOffRampModifier(chainSelector);

    uint64[] memory chainsToRemove = new uint64[](1);
    chainsToRemove[0] = chainSelector;

    vm.startPrank(OWNER);
    s_tokenPool.applyChainUpdates(chainsToRemove, new TokenPool.ChainUpdate[](0));

    vm.startPrank(offRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, chainSelector));
    s_tokenPool.onlyOffRampModifier(chainSelector);
  }

  function test_RevertWhen_CallerIsNotARampOnRouter() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR;
    address offRamp = makeAddr("offRamp");

    vm.startPrank(offRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.CallerIsNotARampOnRouter.selector, offRamp));

    s_tokenPool.onlyOffRampModifier(chainSelector);
  }
}
