// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {BurnMintERC677} from "../../../../shared/token/ERC677/BurnMintERC677.sol";
import {Router} from "../../../Router.sol";
import {LockReleaseTokenPool} from "../../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {RouterSetup} from "../../router/Router/RouterSetup.t.sol";

contract LockReleaseTokenPoolSetup is RouterSetup {
  IERC20 internal s_token;
  LockReleaseTokenPool internal s_lockReleaseTokenPool;
  LockReleaseTokenPool internal s_lockReleaseTokenPoolWithAllowList;
  address[] internal s_allowedList;

  address internal s_allowedOnRamp = address(123);
  address internal s_allowedOffRamp = address(234);

  address internal s_destPoolAddress = address(2736782345);
  address internal s_sourcePoolAddress = address(53852352095);

  function setUp() public virtual override {
    RouterSetup.setUp();
    s_token = new BurnMintERC677("LINK", "LNK", 18, 0);
    deal(address(s_token), OWNER, type(uint256).max);
    s_lockReleaseTokenPool =
      new LockReleaseTokenPool(s_token, new address[](0), address(s_mockRMN), true, address(s_sourceRouter));

    s_allowedList.push(USER_1);
    s_allowedList.push(OWNER);
    s_lockReleaseTokenPoolWithAllowList =
      new LockReleaseTokenPool(s_token, s_allowedList, address(s_mockRMN), true, address(s_sourceRouter));

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(s_destPoolAddress),
      remoteTokenAddress: abi.encode(address(2)),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });

    s_lockReleaseTokenPool.applyChainUpdates(chainUpdate);
    s_lockReleaseTokenPoolWithAllowList.applyChainUpdates(chainUpdate);
    s_lockReleaseTokenPool.setRebalancer(OWNER);

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: s_allowedOnRamp});
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: s_allowedOffRamp});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
  }
}
