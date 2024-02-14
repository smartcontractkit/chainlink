// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {BaseTest} from "../BaseTest.t.sol";
import {TokenPoolHelper} from "../helpers/TokenPoolHelper.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {RouterSetup} from "../router/RouterSetup.t.sol";
import {Router} from "../../Router.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract TokenPoolSetup is RouterSetup {
  IERC20 internal s_token;
  TokenPoolHelper internal s_tokenPool;

  function setUp() public virtual override {
    RouterSetup.setUp();
    s_token = new BurnMintERC677("LINK", "LNK", 18, 0);
    deal(address(s_token), OWNER, type(uint256).max);

    s_tokenPool = new TokenPoolHelper(s_token, new address[](0), address(s_mockARM), address(s_sourceRouter));
  }
}

contract TokenPool_constructor is TokenPoolSetup {
  function test_immutableFieldsSuccess() public {
    assertEq(address(s_token), address(s_tokenPool.getToken()));
    assertEq(address(s_mockARM), s_tokenPool.getArmProxy());
    assertEq(false, s_tokenPool.getAllowListEnabled());
    assertEq(address(s_sourceRouter), s_tokenPool.getRouter());
  }

  // Reverts
  function testZeroAddressNotAllowedReverts() public {
    vm.expectRevert(TokenPool.ZeroAddressNotAllowed.selector);

    s_tokenPool = new TokenPoolHelper(
      IERC20(address(0)),
      new address[](0),
      address(s_mockARM),
      address(s_sourceRouter)
    );
  }
}

contract TokenPool_applyChainUpdates is TokenPoolSetup {
  event ChainAdded(
    uint64 chainSelector,
    RateLimiter.Config outboundRateLimiterConfig,
    RateLimiter.Config inboundRateLimiterConfig
  );
  event ChainRemoved(uint64 chainSelector);

  function assertState(TokenPool.ChainUpdate[] memory chainUpdates) public {
    uint64[] memory chainSelectors = s_tokenPool.getSupportedChains();
    for (uint256 i = 0; i < chainUpdates.length; i++) {
      assertEq(chainUpdates[i].remoteChainSelector, chainSelectors[i]);
    }

    for (uint256 i = 0; i < chainUpdates.length; ++i) {
      assertTrue(s_tokenPool.isSupportedChain(chainUpdates[i].remoteChainSelector));
      RateLimiter.TokenBucket memory bkt = s_tokenPool.getCurrentOutboundRateLimiterState(
        chainUpdates[i].remoteChainSelector
      );
      assertEq(bkt.capacity, chainUpdates[i].outboundRateLimiterConfig.capacity);
      assertEq(bkt.rate, chainUpdates[i].outboundRateLimiterConfig.rate);
      assertEq(bkt.isEnabled, chainUpdates[i].outboundRateLimiterConfig.isEnabled);

      bkt = s_tokenPool.getCurrentInboundRateLimiterState(chainUpdates[i].remoteChainSelector);
      assertEq(bkt.capacity, chainUpdates[i].inboundRateLimiterConfig.capacity);
      assertEq(bkt.rate, chainUpdates[i].inboundRateLimiterConfig.rate);
      assertEq(bkt.isEnabled, chainUpdates[i].inboundRateLimiterConfig.isEnabled);
    }
  }

  function testSuccess() public {
    RateLimiter.Config memory outboundRateLimit1 = RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e18});
    RateLimiter.Config memory inboundRateLimit1 = RateLimiter.Config({isEnabled: true, capacity: 100e29, rate: 1e19});
    RateLimiter.Config memory outboundRateLimit2 = RateLimiter.Config({isEnabled: true, capacity: 100e26, rate: 1e16});
    RateLimiter.Config memory inboundRateLimit2 = RateLimiter.Config({isEnabled: true, capacity: 100e27, rate: 1e17});
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](2);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: 1,
      allowed: true,
      outboundRateLimiterConfig: outboundRateLimit1,
      inboundRateLimiterConfig: inboundRateLimit1
    });
    chainUpdates[1] = TokenPool.ChainUpdate({
      remoteChainSelector: 2,
      allowed: true,
      outboundRateLimiterConfig: outboundRateLimit2,
      inboundRateLimiterConfig: inboundRateLimit2
    });

    // Assert configuration is applied
    vm.expectEmit();
    emit ChainAdded(
      chainUpdates[0].remoteChainSelector,
      chainUpdates[0].outboundRateLimiterConfig,
      chainUpdates[0].inboundRateLimiterConfig
    );
    vm.expectEmit();
    emit ChainAdded(
      chainUpdates[1].remoteChainSelector,
      chainUpdates[1].outboundRateLimiterConfig,
      chainUpdates[1].inboundRateLimiterConfig
    );
    s_tokenPool.applyChainUpdates(chainUpdates);
    // on1: rateLimit1, on2: rateLimit2, off1: rateLimit1, off2: rateLimit3
    assertState(chainUpdates);

    // Removing an non-existent chain should revert
    TokenPool.ChainUpdate[] memory chainRemoves = new TokenPool.ChainUpdate[](1);
    uint64 strangerChainSelector = 120938;
    chainRemoves[0] = TokenPool.ChainUpdate({
      remoteChainSelector: strangerChainSelector,
      allowed: false,
      outboundRateLimiterConfig: RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0}),
      inboundRateLimiterConfig: RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0})
    });
    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, strangerChainSelector));
    s_tokenPool.applyChainUpdates(chainRemoves);
    // State remains
    assertState(chainUpdates);

    // Can remove a chain
    chainRemoves[0].remoteChainSelector = 1;

    vm.expectEmit();
    emit ChainRemoved(chainRemoves[0].remoteChainSelector);

    s_tokenPool.applyChainUpdates(chainRemoves);

    // State updated, only chain 2 remains
    TokenPool.ChainUpdate[] memory singleChainConfigured = new TokenPool.ChainUpdate[](1);
    singleChainConfigured[0] = chainUpdates[1];
    assertState(singleChainConfigured);

    // Cannot reset already configured ramp
    vm.expectRevert(
      abi.encodeWithSelector(TokenPool.ChainAlreadyExists.selector, singleChainConfigured[0].remoteChainSelector)
    );
    s_tokenPool.applyChainUpdates(singleChainConfigured);
  }

  // Reverts

  function testOnlyCallableByOwnerReverts() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_tokenPool.applyChainUpdates(new TokenPool.ChainUpdate[](0));
  }

  function testDisabledNonZeroRateLimitReverts() public {
    RateLimiter.Config memory outboundRateLimit = RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e18});
    RateLimiter.Config memory inboundRateLimit = RateLimiter.Config({isEnabled: true, capacity: 100e22, rate: 1e12});
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: 1,
      allowed: true,
      outboundRateLimiterConfig: outboundRateLimit,
      inboundRateLimiterConfig: inboundRateLimit
    });

    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].allowed = false;
    chainUpdates[0].outboundRateLimiterConfig = RateLimiter.Config({isEnabled: false, capacity: 10, rate: 1});
    chainUpdates[0].inboundRateLimiterConfig = RateLimiter.Config({isEnabled: false, capacity: 10, rate: 1});

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.DisabledNonZeroRateLimit.selector, chainUpdates[0].outboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);
  }

  function testNonExistentChainReverts() public {
    RateLimiter.Config memory outboundRateLimit = RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0});
    RateLimiter.Config memory inboundRateLimit = RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0});
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: 1,
      allowed: false,
      outboundRateLimiterConfig: outboundRateLimit,
      inboundRateLimiterConfig: inboundRateLimit
    });

    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, chainUpdates[0].remoteChainSelector));
    s_tokenPool.applyChainUpdates(chainUpdates);
  }

  function testInvalidRatelimitRateReverts() public {
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: 1,
      allowed: true,
      outboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 0, rate: 0}),
      inboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 100e22, rate: 1e12})
    });

    // Outbound

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRatelimitRate.selector, chainUpdates[0].outboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].outboundRateLimiterConfig.rate = 100;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRatelimitRate.selector, chainUpdates[0].outboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].outboundRateLimiterConfig.capacity = 100;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRatelimitRate.selector, chainUpdates[0].outboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].outboundRateLimiterConfig.capacity = 101;

    s_tokenPool.applyChainUpdates(chainUpdates);

    // Change the chain selector as adding the same one would revert
    chainUpdates[0].remoteChainSelector = 2;

    // Inbound

    chainUpdates[0].inboundRateLimiterConfig.capacity = 0;
    chainUpdates[0].inboundRateLimiterConfig.rate = 0;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRatelimitRate.selector, chainUpdates[0].inboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].inboundRateLimiterConfig.rate = 100;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRatelimitRate.selector, chainUpdates[0].inboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].inboundRateLimiterConfig.capacity = 100;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRatelimitRate.selector, chainUpdates[0].inboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].inboundRateLimiterConfig.capacity = 101;

    s_tokenPool.applyChainUpdates(chainUpdates);
  }
}

contract TokenPool_setChainRateLimiterConfig is TokenPoolSetup {
  event ConfigChanged(RateLimiter.Config);
  event ChainConfigured(
    uint64 chainSelector,
    RateLimiter.Config outboundRateLimiterConfig,
    RateLimiter.Config inboundRateLimiterConfig
  );

  uint64 internal s_remoteChainSelector;

  function setUp() public virtual override {
    TokenPoolSetup.setUp();
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    s_remoteChainSelector = 123124;
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: s_remoteChainSelector,
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(chainUpdates);
  }

  function testFuzz_SetChainRateLimiterConfigSuccess(uint128 capacity, uint128 rate, uint32 newTime) public {
    // Cap the lower bound to 4 so 4/2 is still >= 2
    vm.assume(capacity >= 4);
    // Cap the lower bound to 2 so 2/2 is still >= 1
    rate = uint128(bound(rate, 2, capacity - 2));
    // Bucket updates only work on increasing time
    newTime = uint32(bound(newTime, block.timestamp + 1, type(uint32).max));
    vm.warp(newTime);

    uint256 oldOutboundTokens = s_tokenPool.getCurrentOutboundRateLimiterState(s_remoteChainSelector).tokens;
    uint256 oldInboundTokens = s_tokenPool.getCurrentInboundRateLimiterState(s_remoteChainSelector).tokens;

    RateLimiter.Config memory newOutboundConfig = RateLimiter.Config({isEnabled: true, capacity: capacity, rate: rate});
    RateLimiter.Config memory newInboundConfig = RateLimiter.Config({
      isEnabled: true,
      capacity: capacity / 2,
      rate: rate / 2
    });

    vm.expectEmit();
    emit ConfigChanged(newOutboundConfig);
    vm.expectEmit();
    emit ConfigChanged(newInboundConfig);
    vm.expectEmit();
    emit ChainConfigured(s_remoteChainSelector, newOutboundConfig, newInboundConfig);

    s_tokenPool.setChainRateLimiterConfig(s_remoteChainSelector, newOutboundConfig, newInboundConfig);

    uint256 expectedTokens = RateLimiter._min(newOutboundConfig.capacity, oldOutboundTokens);

    RateLimiter.TokenBucket memory bucket = s_tokenPool.getCurrentOutboundRateLimiterState(s_remoteChainSelector);
    assertEq(bucket.capacity, newOutboundConfig.capacity);
    assertEq(bucket.rate, newOutboundConfig.rate);
    assertEq(bucket.tokens, expectedTokens);
    assertEq(bucket.lastUpdated, newTime);

    expectedTokens = RateLimiter._min(newInboundConfig.capacity, oldInboundTokens);

    bucket = s_tokenPool.getCurrentInboundRateLimiterState(s_remoteChainSelector);
    assertEq(bucket.capacity, newInboundConfig.capacity);
    assertEq(bucket.rate, newInboundConfig.rate);
    assertEq(bucket.tokens, expectedTokens);
    assertEq(bucket.lastUpdated, newTime);
  }

  // Reverts

  function testOnlyOwnerReverts() public {
    changePrank(STRANGER);

    vm.expectRevert("Only callable by owner");
    s_tokenPool.setChainRateLimiterConfig(
      s_remoteChainSelector,
      getOutboundRateLimiterConfig(),
      getInboundRateLimiterConfig()
    );
  }

  function testNonExistentChainReverts() public {
    uint64 wrongChainSelector = 9084102894;

    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, wrongChainSelector));
    s_tokenPool.setChainRateLimiterConfig(
      wrongChainSelector,
      getOutboundRateLimiterConfig(),
      getInboundRateLimiterConfig()
    );
  }
}

contract TokenPool_onlyOnRamp is TokenPoolSetup {
  function test_onlyOnRampSuccess() public {
    uint64 chainSelector = 13377;
    address onRamp = makeAddr("onRamp");

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(chainUpdate);

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: chainSelector, onRamp: onRamp});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), new Router.OffRamp[](0));

    vm.startPrank(onRamp);

    s_tokenPool.onlyOnRampModifier(chainSelector);
  }

  function test_ChainNotAllowedReverts() public {
    uint64 chainSelector = 13377;
    address onRamp = makeAddr("onRamp");

    vm.startPrank(onRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, chainSelector));
    s_tokenPool.onlyOnRampModifier(chainSelector);

    vm.startPrank(OWNER);

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(chainUpdate);

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: chainSelector, onRamp: onRamp});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), new Router.OffRamp[](0));

    vm.startPrank(onRamp);
    // Should succeed now that we've added the chain
    s_tokenPool.onlyOnRampModifier(chainSelector);

    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      allowed: false,
      outboundRateLimiterConfig: RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0}),
      inboundRateLimiterConfig: RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0})
    });

    vm.startPrank(OWNER);
    s_tokenPool.applyChainUpdates(chainUpdate);

    vm.startPrank(onRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, chainSelector));
    s_tokenPool.onlyOffRampModifier(chainSelector);
  }

  function test_CallerIsNotARampOnRouterReverts() public {
    uint64 chainSelector = 13377;
    address onRamp = makeAddr("onRamp");

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(chainUpdate);

    vm.startPrank(onRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.CallerIsNotARampOnRouter.selector, onRamp));

    s_tokenPool.onlyOnRampModifier(chainSelector);
  }
}

contract TokenPool_onlyOffRamp is TokenPoolSetup {
  function test_onlyOffRampSuccess() public {
    uint64 chainSelector = 13377;
    address offRamp = makeAddr("onRamp");

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(chainUpdate);

    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: chainSelector, offRamp: offRamp});
    s_sourceRouter.applyRampUpdates(new Router.OnRamp[](0), new Router.OffRamp[](0), offRampUpdates);

    vm.startPrank(offRamp);

    s_tokenPool.onlyOffRampModifier(chainSelector);
  }

  function test_ChainNotAllowedReverts() public {
    uint64 chainSelector = 13377;
    address offRamp = makeAddr("onRamp");

    vm.startPrank(offRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, chainSelector));
    s_tokenPool.onlyOffRampModifier(chainSelector);

    vm.startPrank(OWNER);

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
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

  function test_CallerIsNotARampOnRouterReverts() public {
    uint64 chainSelector = 13377;
    address offRamp = makeAddr("offRamp");

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: chainSelector,
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(chainUpdate);

    vm.startPrank(offRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.CallerIsNotARampOnRouter.selector, offRamp));

    s_tokenPool.onlyOffRampModifier(chainSelector);
  }
}

contract TokenPoolWithAllowListSetup is TokenPoolSetup {
  address[] internal s_allowedSenders;

  function setUp() public virtual override {
    TokenPoolSetup.setUp();

    s_allowedSenders.push(STRANGER);
    s_allowedSenders.push(DUMMY_CONTRACT_ADDRESS);

    s_tokenPool = new TokenPoolHelper(s_token, s_allowedSenders, address(s_mockARM), address(s_sourceRouter));
  }
}

/// @notice #getAllowListEnabled
contract TokenPoolWithAllowList_getAllowListEnabled is TokenPoolWithAllowListSetup {
  function testGetAllowListEnabledSuccess() public {
    assertTrue(s_tokenPool.getAllowListEnabled());
  }
}

contract TokenPoolWithAllowList_setRouter is TokenPoolWithAllowListSetup {
  event RouterUpdated(address oldRouter, address newRouter);

  function testSetRouterSuccess() public {
    assertEq(address(s_sourceRouter), s_tokenPool.getRouter());

    address newRouter = makeAddr("newRouter");

    vm.expectEmit();
    emit RouterUpdated(address(s_sourceRouter), newRouter);

    s_tokenPool.setRouter(newRouter);

    assertEq(newRouter, s_tokenPool.getRouter());
  }
}

/// @notice #getAllowList
contract TokenPoolWithAllowList_getAllowList is TokenPoolWithAllowListSetup {
  function testGetAllowListSuccess() public {
    address[] memory setAddresses = s_tokenPool.getAllowList();
    assertEq(2, setAddresses.length);
    assertEq(s_allowedSenders[0], setAddresses[0]);
    assertEq(s_allowedSenders[1], setAddresses[1]);
  }
}

/// @notice #setAllowList
contract TokenPoolWithAllowList_applyAllowListUpdates is TokenPoolWithAllowListSetup {
  event AllowListAdd(address sender);
  event AllowListRemove(address sender);

  function testSetAllowListSuccess() public {
    address[] memory newAddresses = new address[](2);
    newAddresses[0] = address(1);
    newAddresses[1] = address(2);

    for (uint256 i = 0; i < 2; ++i) {
      vm.expectEmit();
      emit AllowListAdd(newAddresses[i]);
    }

    s_tokenPool.applyAllowListUpdates(new address[](0), newAddresses);
    address[] memory setAddresses = s_tokenPool.getAllowList();

    assertEq(s_allowedSenders[0], setAddresses[0]);
    assertEq(s_allowedSenders[1], setAddresses[1]);
    assertEq(address(1), setAddresses[2]);
    assertEq(address(2), setAddresses[3]);

    // address(2) exists noop, add address(3), remove address(1)
    newAddresses = new address[](2);
    newAddresses[0] = address(2);
    newAddresses[1] = address(3);

    address[] memory removeAddresses = new address[](1);
    removeAddresses[0] = address(1);

    vm.expectEmit();
    emit AllowListRemove(address(1));

    vm.expectEmit();
    emit AllowListAdd(address(3));

    s_tokenPool.applyAllowListUpdates(removeAddresses, newAddresses);
    setAddresses = s_tokenPool.getAllowList();

    assertEq(s_allowedSenders[0], setAddresses[0]);
    assertEq(s_allowedSenders[1], setAddresses[1]);
    assertEq(address(2), setAddresses[2]);
    assertEq(address(3), setAddresses[3]);

    // remove all from allowList
    for (uint256 i = 0; i < setAddresses.length; ++i) {
      vm.expectEmit();
      emit AllowListRemove(setAddresses[i]);
    }

    s_tokenPool.applyAllowListUpdates(setAddresses, new address[](0));
    setAddresses = s_tokenPool.getAllowList();

    assertEq(0, setAddresses.length);
  }

  function testSetAllowListSkipsZeroSuccess() public {
    uint256 setAddressesLength = s_tokenPool.getAllowList().length;

    address[] memory newAddresses = new address[](1);
    newAddresses[0] = address(0);

    s_tokenPool.applyAllowListUpdates(new address[](0), newAddresses);
    address[] memory setAddresses = s_tokenPool.getAllowList();

    assertEq(setAddresses.length, setAddressesLength);
  }

  // Reverts

  function testOnlyOwnerReverts() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    address[] memory newAddresses = new address[](2);
    s_tokenPool.applyAllowListUpdates(new address[](0), newAddresses);
  }

  function test_AllowListNotEnabledReverts() public {
    s_tokenPool = new TokenPoolHelper(s_token, new address[](0), address(s_mockARM), address(s_sourceRouter));

    vm.expectRevert(TokenPool.AllowListNotEnabled.selector);

    s_tokenPool.applyAllowListUpdates(new address[](0), new address[](2));
  }
}
