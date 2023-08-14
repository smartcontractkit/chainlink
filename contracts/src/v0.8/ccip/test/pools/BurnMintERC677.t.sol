// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import "../BaseTest.t.sol";
import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {BurnMintTokenPool} from "../../pools/BurnMintTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";

contract BurnMintERC677Setup is BaseTest {
  event Transfer(address indexed from, address indexed to, uint256 value);

  BurnMintERC677 internal s_burnMintERC677;
  BurnMintTokenPool internal s_pool;
  address internal s_burnMintOffRamp;
  address internal s_burnMintOnRamp;

  function setUp() public virtual override {
    BaseTest.setUp();
    s_burnMintERC677 = new BurnMintERC677("Chainlink Token", "LINK", 18, 0);
    s_pool = new BurnMintTokenPool(s_burnMintERC677, new address[](0), address(s_mockARM));
    s_burnMintERC677.grantMintAndBurnRoles(address(s_pool));
    s_burnMintOffRamp = address(238323465456);
    TokenPool.RampUpdate[] memory offRamps = new TokenPool.RampUpdate[](1);
    offRamps[0] = TokenPool.RampUpdate({
      ramp: s_burnMintOffRamp,
      allowed: true,
      rateLimiterConfig: rateLimiterConfig()
    });
    s_pool.applyRampUpdates(new TokenPool.RampUpdate[](0), offRamps);

    s_burnMintOnRamp = address(238323465457);
    TokenPool.RampUpdate[] memory onRamps = new TokenPool.RampUpdate[](1);
    onRamps[0] = TokenPool.RampUpdate({ramp: s_burnMintOnRamp, allowed: true, rateLimiterConfig: rateLimiterConfig()});
    s_pool.applyRampUpdates(onRamps, new TokenPool.RampUpdate[](0));
  }
}

contract BurnMintERC677_mint is BurnMintERC677Setup {
  function testPoolMintSuccess() public {
    uint256 amount = 1e19;
    changePrank(s_burnMintOffRamp);
    vm.expectEmit();
    emit Transfer(address(0), OWNER, amount);
    s_pool.releaseOrMint(bytes(""), OWNER, amount, 0, bytes(""));
    assertEq(s_burnMintERC677.balanceOf(OWNER), amount);
  }

  function testPoolMintNotHealthyReverts() public {
    // Should not mint tokens if cursed.
    s_mockARM.voteToCurse(bytes32(0));
    uint256 before = s_burnMintERC677.balanceOf(OWNER);
    changePrank(s_burnMintOffRamp);
    vm.expectRevert(EVM2EVMOffRamp.BadARMSignal.selector);
    s_pool.releaseOrMint(bytes(""), OWNER, 1e5, 0, bytes(""));
    assertEq(s_burnMintERC677.balanceOf(OWNER), before);
  }
}

contract BurnMintERC677_burn is BurnMintERC677Setup {
  function testPoolBurnSuccess() public {
    uint256 burnAmount = 1e19;

    deal(address(s_burnMintERC677), address(s_pool), burnAmount);
    changePrank(s_burnMintOnRamp);

    vm.expectEmit();
    emit Transfer(address(s_pool), address(0), burnAmount);

    s_pool.lockOrBurn(OWNER, bytes(""), burnAmount, 0, bytes(""));

    assertEq(s_burnMintERC677.balanceOf(address(s_pool)), 0);
  }

  function testPoolBurnRevertNotHealthyReverts() public {
    // Should not burn tokens if cursed.
    s_mockARM.voteToCurse(bytes32(0));
    uint256 before = s_burnMintERC677.balanceOf(address(s_pool));
    changePrank(s_burnMintOnRamp);
    vm.expectRevert(EVM2EVMOnRamp.BadARMSignal.selector);
    s_pool.lockOrBurn(OWNER, bytes(""), 1e5, 0, bytes(""));
    assertEq(s_burnMintERC677.balanceOf(address(s_pool)), before);
  }
}
