// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import "../BaseTest.t.sol";
import {ThirdPartyBurnMintTokenPool} from "../../pools/ThirdPartyBurnMintTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {Router} from "../../Router.sol";
import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";

contract ThirdPartyBurnMintTokenPoolSetup is BaseTest {
  IERC20 internal s_token;
  address internal s_routerAllowedOffRamp = address(234);
  Router internal s_router;

  ThirdPartyBurnMintTokenPool internal s_thirdPartyPool;
  ThirdPartyBurnMintTokenPool internal s_thirdPartyPoolWithAllowList;
  address[] internal s_allowedList;
  address internal s_allowedOnRamp = address(123456789);

  function setUp() public virtual override {
    BaseTest.setUp();
    s_token = new BurnMintERC677("LINK", "LNK", 18, 0);
    deal(address(s_token), OWNER, type(uint256).max);

    s_router = new Router(address(s_token), address(s_mockARM));

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](0);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    address[] memory offRamps = new address[](1);
    offRamps[0] = s_routerAllowedOffRamp;
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_ID, offRamp: s_routerAllowedOffRamp});

    s_router.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);

    s_thirdPartyPool = new ThirdPartyBurnMintTokenPool(
      IBurnMintERC20(address(s_token)),
      new address[](0),
      address(s_router),
      address(s_mockARM)
    );

    s_allowedList.push(USER_1);
    s_allowedList.push(DUMMY_CONTRACT_ADDRESS);
    s_thirdPartyPoolWithAllowList = new ThirdPartyBurnMintTokenPool(
      IBurnMintERC20(address(s_token)),
      s_allowedList,
      address(s_router),
      address(s_mockARM)
    );

    BurnMintERC677(address(s_token)).grantMintAndBurnRoles(address(s_thirdPartyPool));
    BurnMintERC677(address(s_token)).grantMintAndBurnRoles(address(s_thirdPartyPoolWithAllowList));

    TokenPool.RampUpdate[] memory poolOnRampsUpdates = new TokenPool.RampUpdate[](1);
    poolOnRampsUpdates[0] = TokenPool.RampUpdate({
      ramp: s_allowedOnRamp,
      allowed: true,
      rateLimiterConfig: rateLimiterConfig()
    });
    TokenPool.RampUpdate[] memory poolOffRampUpdates = new TokenPool.RampUpdate[](0);

    s_thirdPartyPool.applyRampUpdates(poolOnRampsUpdates, poolOffRampUpdates);
    s_thirdPartyPoolWithAllowList.applyRampUpdates(poolOnRampsUpdates, poolOffRampUpdates);
  }
}

contract ThirdPartyBurnMintTokenPool_lockOrBurn is ThirdPartyBurnMintTokenPoolSetup {
  error SenderNotAllowed(address sender);
  event Burned(address indexed sender, uint256 amount);
  event TokensConsumed(uint256 amount);

  function setUp() public virtual override {
    ThirdPartyBurnMintTokenPoolSetup.setUp();

    deal(address(s_token), address(s_thirdPartyPool), type(uint256).max);
    deal(address(s_token), address(s_thirdPartyPoolWithAllowList), type(uint256).max);
  }

  function testFuzz_LockOrBurnNoAllowListSuccess(uint256 amount) public {
    amount = bound(amount, 1, rateLimiterConfig().capacity);
    changePrank(s_allowedOnRamp);

    vm.expectEmit();
    emit TokensConsumed(amount);
    vm.expectEmit();
    emit Burned(s_allowedOnRamp, amount);

    s_thirdPartyPool.lockOrBurn(STRANGER, bytes(""), amount, DEST_CHAIN_ID, bytes(""));
  }

  function testLockOrBurnWithAllowListSuccess() public {
    uint256 amount = 100;
    changePrank(s_allowedOnRamp);

    vm.expectEmit();
    emit TokensConsumed(amount);
    vm.expectEmit();
    emit Burned(s_allowedOnRamp, amount);

    s_thirdPartyPoolWithAllowList.lockOrBurn(s_allowedList[0], bytes(""), amount, DEST_CHAIN_ID, bytes(""));

    vm.expectEmit();
    emit TokensConsumed(amount);
    vm.expectEmit();
    emit Burned(s_allowedOnRamp, amount);

    s_thirdPartyPoolWithAllowList.lockOrBurn(s_allowedList[1], bytes(""), amount, DEST_CHAIN_ID, bytes(""));
  }

  function testLockOrBurnWithAllowListReverts() public {
    changePrank(s_allowedOnRamp);

    vm.expectRevert(abi.encodeWithSelector(SenderNotAllowed.selector, STRANGER));

    s_thirdPartyPoolWithAllowList.lockOrBurn(STRANGER, bytes(""), 100, DEST_CHAIN_ID, bytes(""));
  }
}

contract ThirdPartyBurnMintTokenPool_applyRampUpdates is ThirdPartyBurnMintTokenPoolSetup {
  // Note applyRampUpdates inherits from TokenPool so we only need to test the new functionality.
  // Reverts
  function testInvalidOffRampReverts() public {
    address invalidOffRamp = address(23456787654321);
    TokenPool.RampUpdate[] memory offRamps = new TokenPool.RampUpdate[](1);
    offRamps[0] = TokenPool.RampUpdate({ramp: invalidOffRamp, allowed: true, rateLimiterConfig: rateLimiterConfig()});

    vm.expectRevert(abi.encodeWithSelector(ThirdPartyBurnMintTokenPool.InvalidOffRamp.selector, invalidOffRamp));

    s_thirdPartyPool.applyRampUpdates(new TokenPool.RampUpdate[](0), offRamps);
  }
}
