// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Router} from "../../../Router.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_onlyOffRamp is TokenPoolSetup {
  function test_onlyOffRamp_Success() public {
    uint64 chainSelector = 13377;
    address offRamp = makeAddr("onRamp");

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      remotePoolAddress: abi.encode(address(1)),
      remoteTokenAddress: abi.encode(address(2)),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(chainUpdate);

    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: chainSelector, offRamp: offRamp});
    s_sourceRouter.applyRampUpdates(new Router.OnRamp[](0), new Router.OffRamp[](0), offRampUpdates);

    vm.startPrank(offRamp);

    s_tokenPool.onlyOffRampModifier(chainSelector);
  }

  function test_ChainNotAllowed_Revert() public {
    uint64 chainSelector = 13377;
    address offRamp = makeAddr("onRamp");

    vm.startPrank(offRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, chainSelector));
    s_tokenPool.onlyOffRampModifier(chainSelector);

    vm.startPrank(OWNER);

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      remotePoolAddress: abi.encode(address(1)),
      remoteTokenAddress: abi.encode(address(2)),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(chainUpdate);

    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: chainSelector, offRamp: offRamp});
    s_sourceRouter.applyRampUpdates(new Router.OnRamp[](0), new Router.OffRamp[](0), offRampUpdates);

    vm.startPrank(offRamp);
    // Should succeed now that we've added the chain
    s_tokenPool.onlyOffRampModifier(chainSelector);

    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      remotePoolAddress: abi.encode(address(1)),
      remoteTokenAddress: abi.encode(address(2)),
      allowed: false,
      outboundRateLimiterConfig: RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0}),
      inboundRateLimiterConfig: RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0})
    });

    vm.startPrank(OWNER);
    s_tokenPool.applyChainUpdates(chainUpdate);

    vm.startPrank(offRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, chainSelector));
    s_tokenPool.onlyOffRampModifier(chainSelector);
  }

  function test_CallerIsNotARampOnRouter_Revert() public {
    uint64 chainSelector = 13377;
    address offRamp = makeAddr("offRamp");

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      remotePoolAddress: abi.encode(address(1)),
      remoteTokenAddress: abi.encode(address(2)),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(chainUpdate);

    vm.startPrank(offRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.CallerIsNotARampOnRouter.selector, offRamp));

    s_tokenPool.onlyOffRampModifier(chainSelector);
  }
}
