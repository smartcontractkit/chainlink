// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import "../BaseTest.t.sol";
import {TokenPoolHelper} from "../helpers/TokenPoolHelper.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";

contract TokenPoolSetup is BaseTest {
  IERC20 internal s_token;
  TokenPoolHelper internal s_tokenPool;

  function setUp() public virtual override {
    BaseTest.setUp();
    s_token = new BurnMintERC677("LINK", "LNK", 18, 0);
    deal(address(s_token), OWNER, type(uint256).max);

    s_tokenPool = new TokenPoolHelper(s_token, new address[](0), address(s_mockARM));
  }
}

contract TokenPool_constructor is TokenPoolSetup {
  // Reverts
  function testZeroAddressNotAllowedReverts() public {
    vm.expectRevert(TokenPool.ZeroAddressNotAllowed.selector);

    s_tokenPool = new TokenPoolHelper(IERC20(address(0)), new address[](0), address(s_mockARM));
  }
}

contract TokenPool_applyRampUpdates is TokenPoolSetup {
  event OnRampAdded(address onRamp, RateLimiter.Config rateLimiterConfig);
  event OnRampRemoved(address onRamp);
  event OffRampAdded(address offRamp, RateLimiter.Config rateLimiterConfig);
  event OffRampRemoved(address offRamp);

  function addressesFromUpdates(TokenPool.RampUpdate[] memory updates) public pure returns (address[] memory) {
    address[] memory addresses = new address[](updates.length);
    for (uint256 i = 0; i < updates.length; i++) {
      addresses[i] = updates[i].ramp;
    }
    return addresses;
  }

  function assertState(TokenPool.RampUpdate[] memory onRamps, TokenPool.RampUpdate[] memory offRamps) public {
    assertEq(s_tokenPool.getOnRamps(), addressesFromUpdates(onRamps));
    for (uint256 i = 0; i < onRamps.length; ++i) {
      assertTrue(s_tokenPool.isOnRamp(onRamps[i].ramp));
      RateLimiter.TokenBucket memory bkt = s_tokenPool.currentOnRampRateLimiterState(onRamps[i].ramp);
      assertEq(bkt.capacity, onRamps[i].rateLimiterConfig.capacity);
      assertEq(bkt.rate, onRamps[i].rateLimiterConfig.rate);
      assertEq(bkt.isEnabled, onRamps[i].rateLimiterConfig.isEnabled);
    }
    assertEq(s_tokenPool.getOffRamps(), addressesFromUpdates(offRamps));
    for (uint256 i = 0; i < offRamps.length; ++i) {
      assertTrue(s_tokenPool.isOffRamp(offRamps[i].ramp));
      RateLimiter.TokenBucket memory bkt = s_tokenPool.currentOffRampRateLimiterState(offRamps[i].ramp);
      assertEq(bkt.capacity, offRamps[i].rateLimiterConfig.capacity);
      assertEq(bkt.rate, offRamps[i].rateLimiterConfig.rate);
      assertEq(bkt.isEnabled, offRamps[i].rateLimiterConfig.isEnabled);
    }
  }

  function testApplyRampUpdatesSuccess() public {
    // Create on and offramps.
    RateLimiter.Config memory rateLimit1 = RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e15});
    RateLimiter.Config memory rateLimit2 = RateLimiter.Config({isEnabled: true, capacity: 101e28, rate: 1e14});
    TokenPool.RampUpdate[] memory onRampUpdates1 = new TokenPool.RampUpdate[](2);
    onRampUpdates1[0] = TokenPool.RampUpdate({ramp: address(1), allowed: true, rateLimiterConfig: rateLimit1});
    onRampUpdates1[1] = TokenPool.RampUpdate({ramp: address(2), allowed: true, rateLimiterConfig: rateLimit2});
    TokenPool.RampUpdate[] memory offRampUpdates1 = new TokenPool.RampUpdate[](2);
    offRampUpdates1[0] = TokenPool.RampUpdate({ramp: address(11), allowed: true, rateLimiterConfig: rateLimit2});
    offRampUpdates1[1] = TokenPool.RampUpdate({ramp: address(12), allowed: true, rateLimiterConfig: rateLimit1});

    // Assert configuration is applied
    vm.expectEmit();
    emit OnRampAdded(onRampUpdates1[0].ramp, onRampUpdates1[0].rateLimiterConfig);
    vm.expectEmit();
    emit OnRampAdded(onRampUpdates1[1].ramp, onRampUpdates1[1].rateLimiterConfig);
    vm.expectEmit();
    emit OffRampAdded(offRampUpdates1[0].ramp, offRampUpdates1[0].rateLimiterConfig);
    vm.expectEmit();
    emit OffRampAdded(offRampUpdates1[1].ramp, offRampUpdates1[1].rateLimiterConfig);
    s_tokenPool.applyRampUpdates(onRampUpdates1, offRampUpdates1);
    // on1: rateLimit1, on2: rateLimit2, off1: rateLimit1, off2: rateLimit3
    assertState(onRampUpdates1, offRampUpdates1);

    // Removing an non-existent onRamp should revert
    TokenPool.RampUpdate[] memory onRampRemoves = new TokenPool.RampUpdate[](1);
    address strangerOnRamp = address(120938);
    onRampRemoves[0] = TokenPool.RampUpdate({ramp: strangerOnRamp, allowed: false, rateLimiterConfig: rateLimit1});
    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentRamp.selector, strangerOnRamp));
    s_tokenPool.applyRampUpdates(onRampRemoves, offRampUpdates1);
    // State remains
    assertState(onRampUpdates1, offRampUpdates1);

    // Same with offramps
    TokenPool.RampUpdate[] memory offRampRemoves = new TokenPool.RampUpdate[](1);
    offRampRemoves[0] = TokenPool.RampUpdate({ramp: strangerOnRamp, allowed: false, rateLimiterConfig: rateLimit1});
    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentRamp.selector, strangerOnRamp));
    s_tokenPool.applyRampUpdates(new TokenPool.RampUpdate[](0), offRampRemoves);
    // State remains
    assertState(onRampUpdates1, offRampUpdates1);

    // Can remove an onRamp
    onRampRemoves[0].ramp = address(1);
    vm.expectEmit();
    emit OnRampRemoved(onRampRemoves[0].ramp);
    s_tokenPool.applyRampUpdates(onRampRemoves, new TokenPool.RampUpdate[](0));
    // State updated, only onRamp2 remains
    TokenPool.RampUpdate[] memory onRampUpdates2 = new TokenPool.RampUpdate[](1);
    onRampUpdates2[0] = onRampUpdates1[1];
    assertState(onRampUpdates2, offRampUpdates1);

    // Can remove an offRamp
    offRampRemoves[0].ramp = address(11);
    vm.expectEmit();
    emit OffRampRemoved(offRampRemoves[0].ramp);
    s_tokenPool.applyRampUpdates(new TokenPool.RampUpdate[](0), offRampRemoves);
    TokenPool.RampUpdate[] memory offRampUpdates2 = new TokenPool.RampUpdate[](1);
    offRampUpdates2[0] = offRampUpdates1[1];
    assertState(onRampUpdates2, offRampUpdates2);

    // Cannot reset already configured ramp
    vm.expectRevert(abi.encodeWithSelector(TokenPool.RampAlreadyExists.selector, address(2)));
    s_tokenPool.applyRampUpdates(onRampUpdates2, new TokenPool.RampUpdate[](0));
    vm.expectRevert(abi.encodeWithSelector(TokenPool.RampAlreadyExists.selector, address(12)));
    s_tokenPool.applyRampUpdates(new TokenPool.RampUpdate[](0), offRampUpdates2);
  }

  // Reverts

  function testOnlyCallableByOwnerReverts() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_tokenPool.applyRampUpdates(new TokenPool.RampUpdate[](0), new TokenPool.RampUpdate[](0));
  }
}

contract TokenPool_setOnRampRateLimiterConfig is TokenPoolSetup {
  event ConfigChanged(RateLimiter.Config);
  event OnRampConfigured(address onRamp, RateLimiter.Config);
  address internal s_onRamp;

  function setUp() public virtual override {
    TokenPoolSetup.setUp();
    TokenPool.RampUpdate[] memory onRampUpdates1 = new TokenPool.RampUpdate[](1);
    s_onRamp = address(1);
    onRampUpdates1[0] = TokenPool.RampUpdate({ramp: s_onRamp, allowed: true, rateLimiterConfig: rateLimiterConfig()});
    s_tokenPool.applyRampUpdates(onRampUpdates1, new TokenPool.RampUpdate[](0));
  }

  function testFuzz_SetRateLimiterConfigSuccess(uint128 capacity, uint128 rate, uint32 newTime) public {
    // Bucket updates only work on increasing time
    vm.assume(newTime >= block.timestamp);
    vm.warp(newTime);

    uint256 oldTokens = s_tokenPool.currentOnRampRateLimiterState(s_onRamp).tokens;

    RateLimiter.Config memory newConfig = RateLimiter.Config({isEnabled: true, capacity: capacity, rate: rate});

    vm.expectEmit();
    emit ConfigChanged(newConfig);
    vm.expectEmit();
    emit OnRampConfigured(s_onRamp, newConfig);

    s_tokenPool.setOnRampRateLimiterConfig(s_onRamp, newConfig);

    uint256 expectedTokens = RateLimiter._min(newConfig.capacity, oldTokens);

    RateLimiter.TokenBucket memory bucket = s_tokenPool.currentOnRampRateLimiterState(s_onRamp);
    assertEq(bucket.capacity, newConfig.capacity);
    assertEq(bucket.rate, newConfig.rate);
    assertEq(bucket.tokens, expectedTokens);
    assertEq(bucket.lastUpdated, newTime);
  }

  // Reverts

  function testOnlyOwnerReverts() public {
    changePrank(STRANGER);

    vm.expectRevert("Only callable by owner");
    s_tokenPool.setOnRampRateLimiterConfig(s_onRamp, rateLimiterConfig());
  }

  function testNonExistentRampReverts() public {
    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentRamp.selector, address(120938)));
    s_tokenPool.setOnRampRateLimiterConfig(address(120938), rateLimiterConfig());
  }
}

contract TokenPool_setOffRampRateLimiterConfig is TokenPoolSetup {
  event ConfigChanged(RateLimiter.Config);
  event OffRampConfigured(address onRamp, RateLimiter.Config);
  address internal s_offRamp;

  function setUp() public virtual override {
    TokenPoolSetup.setUp();
    TokenPool.RampUpdate[] memory offRampUpdates1 = new TokenPool.RampUpdate[](1);
    s_offRamp = address(1);
    offRampUpdates1[0] = TokenPool.RampUpdate({ramp: s_offRamp, allowed: true, rateLimiterConfig: rateLimiterConfig()});
    s_tokenPool.applyRampUpdates(new TokenPool.RampUpdate[](0), offRampUpdates1);
  }

  function testFuzz_SetRateLimiterConfigSuccess(uint128 capacity, uint128 rate, uint32 newTime) public {
    // Bucket updates only work on increasing time
    vm.assume(newTime >= block.timestamp);
    vm.warp(newTime);

    uint256 oldTokens = s_tokenPool.currentOffRampRateLimiterState(s_offRamp).tokens;

    RateLimiter.Config memory newConfig = RateLimiter.Config({isEnabled: true, capacity: capacity, rate: rate});

    vm.expectEmit();
    emit ConfigChanged(newConfig);
    vm.expectEmit();
    emit OffRampConfigured(s_offRamp, newConfig);

    s_tokenPool.setOffRampRateLimiterConfig(s_offRamp, newConfig);

    uint256 expectedTokens = RateLimiter._min(newConfig.capacity, oldTokens);

    RateLimiter.TokenBucket memory bucket = s_tokenPool.currentOffRampRateLimiterState(s_offRamp);
    assertEq(bucket.capacity, newConfig.capacity);
    assertEq(bucket.rate, newConfig.rate);
    assertEq(bucket.tokens, expectedTokens);
    assertEq(bucket.lastUpdated, newTime);
  }

  // Reverts

  function testOnlyOwnerReverts() public {
    changePrank(STRANGER);

    vm.expectRevert("Only callable by owner");
    s_tokenPool.setOffRampRateLimiterConfig(s_offRamp, rateLimiterConfig());
  }

  function testNonExistentRampReverts() public {
    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentRamp.selector, address(120938)));
    s_tokenPool.setOffRampRateLimiterConfig(address(120938), rateLimiterConfig());
  }
}

contract TokenPoolWithAllowListSetup is TokenPoolSetup {
  address[] internal s_allowedSenders;

  function setUp() public virtual override {
    TokenPoolSetup.setUp();

    s_allowedSenders.push(STRANGER);
    s_allowedSenders.push(DUMMY_CONTRACT_ADDRESS);

    s_tokenPool = new TokenPoolHelper(s_token, s_allowedSenders, address(s_mockARM));
  }
}

/// @notice #getAllowListEnabled
contract TokenPoolWithAllowList_getAllowListEnabled is TokenPoolWithAllowListSetup {
  function testGetAllowListEnabledSuccess() public {
    assertTrue(s_tokenPool.getAllowListEnabled());
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
}
