// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {BurnMintERC677} from "../../../../shared/token/ERC677/BurnMintERC677.sol";
import {Router} from "../../../Router.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolHelper} from "../../helpers/TokenPoolHelper.sol";
import {RouterSetup} from "../../router/RouterSetup.t.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

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
