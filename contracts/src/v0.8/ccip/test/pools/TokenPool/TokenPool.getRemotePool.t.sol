// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_getRemotePool is TokenPoolSetup {
  function test_getRemotePools() public {
    uint64 chainSelector = DEST_CHAIN_SELECTOR + 1;
    bytes memory remotePool = abi.encode(makeAddr("remotePool"));
    bytes memory remoteToken = abi.encode(makeAddr("remoteToken"));

    // Zero indicates nothing is set
    assertEq(0, s_tokenPool.getRemotePools(chainSelector).length);

    bytes[] memory remotePoolAddresses = new bytes[](1);
    remotePoolAddresses[0] = remotePool;

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      remotePoolAddresses: remotePoolAddresses,
      remoteTokenAddress: remoteToken,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    bytes[] memory remotePools = s_tokenPool.getRemotePools(chainSelector);
    assertEq(1, remotePools.length);
    assertEq(remotePool, remotePools[0]);
  }
}
