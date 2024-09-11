// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {Router} from "../../Router.sol";
import {BurnMintTokenPool} from "../../pools/BurnMintTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {RouterSetup} from "../router/RouterSetup.t.sol";

contract BurnMintSetup is RouterSetup {
  BurnMintERC677 internal s_burnMintERC677;
  address internal s_burnMintOffRamp = makeAddr("burn_mint_offRamp");
  address internal s_burnMintOnRamp = makeAddr("burn_mint_onRamp");

  address internal s_remoteBurnMintPool = makeAddr("remote_burn_mint_pool");
  address internal s_remoteToken = makeAddr("remote_token");

  function setUp() public virtual override {
    RouterSetup.setUp();

    s_burnMintERC677 = new BurnMintERC677("Chainlink Token", "LINK", 18, 0);
  }

  function _applyChainUpdates(address pool) internal {
    TokenPool.ChainUpdate[] memory chains = new TokenPool.ChainUpdate[](1);
    chains[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(s_remoteBurnMintPool),
      remoteTokenAddress: abi.encode(s_remoteToken),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });

    BurnMintTokenPool(pool).applyChainUpdates(chains);

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: s_burnMintOnRamp});
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: DEST_CHAIN_SELECTOR, offRamp: s_burnMintOffRamp});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
  }
}
