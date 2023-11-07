// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import "../BaseTest.t.sol";
import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {BurnMintTokenPool} from "../../pools/BurnMintTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";

contract BurnMintSetup is BaseTest {
  event Transfer(address indexed from, address indexed to, uint256 value);
  event TokensConsumed(uint256 tokens);
  event Burned(address indexed sender, uint256 amount);

  BurnMintERC677 internal s_burnMintERC677;
  address internal s_burnMintOffRamp = makeAddr("burn_mint_offRamp");
  address internal s_burnMintOnRamp = makeAddr("burn_mint_onRamp");

  function setUp() public virtual override {
    BaseTest.setUp();

    s_burnMintERC677 = new BurnMintERC677("Chainlink Token", "LINK", 18, 0);
  }

  function applyRampUpdates(address pool) internal {
    TokenPool.RampUpdate[] memory offRamps = new TokenPool.RampUpdate[](1);
    offRamps[0] = TokenPool.RampUpdate({
      ramp: s_burnMintOffRamp,
      allowed: true,
      rateLimiterConfig: rateLimiterConfig()
    });

    TokenPool.RampUpdate[] memory onRamps = new TokenPool.RampUpdate[](1);
    onRamps[0] = TokenPool.RampUpdate({ramp: s_burnMintOnRamp, allowed: true, rateLimiterConfig: rateLimiterConfig()});
    BurnMintTokenPool(pool).applyRampUpdates(onRamps, offRamps);
  }
}
