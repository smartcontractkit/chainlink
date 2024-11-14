// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_getRemotePool is TokenPoolSetup {
  function test_getRemotePool_Success() public {
    uint64 chainSelector = 123124;
    address remotePool = makeAddr("remotePool");
    address remoteToken = makeAddr("remoteToken");

    // Zero indicates nothing is set
    assertEq(0, s_tokenPool.getRemotePool(chainSelector).length);

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      remotePoolAddress: abi.encode(remotePool),
      remoteTokenAddress: abi.encode(remoteToken),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(chainUpdates);

    assertEq(remotePool, abi.decode(s_tokenPool.getRemotePool(chainSelector), (address)));
  }
}
