// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IPool} from "../../interfaces/pools/IPool.sol";

import {BaseTest} from "../BaseTest.t.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {EVM2EVMOnRamp} from "../../onRamp/EVM2EVMOnRamp.sol";
import {EVM2EVMOffRamp} from "../../offRamp/EVM2EVMOffRamp.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {Router} from "../../Router.sol";

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/IERC165.sol";
import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {RouterSetup} from "../router/RouterSetup.t.sol";

contract LockReleaseTokenPoolSetup is RouterSetup {
  IERC20 internal s_token;
  LockReleaseTokenPool internal s_lockReleaseTokenPool;
  LockReleaseTokenPool internal s_lockReleaseTokenPoolWithAllowList;
  address[] internal s_allowedList;

  address internal s_allowedOnRamp = address(123);
  address internal s_allowedOffRamp = address(234);

  function setUp() public virtual override {
    RouterSetup.setUp();
    s_token = new BurnMintERC677("LINK", "LNK", 18, 0);
    deal(address(s_token), OWNER, type(uint256).max);
    s_lockReleaseTokenPool = new LockReleaseTokenPool(
      s_token,
      new address[](0),
      address(s_mockARM),
      true,
      address(s_sourceRouter)
    );

    s_allowedList.push(USER_1);
    s_allowedList.push(DUMMY_CONTRACT_ADDRESS);
    s_lockReleaseTokenPoolWithAllowList = new LockReleaseTokenPool(
      s_token,
      s_allowedList,
      address(s_mockARM),
      true,
      address(s_sourceRouter)
    );

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
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

contract LockReleaseTokenPool_setRebalancer is LockReleaseTokenPoolSetup {
  function testSetRebalancerSuccess() public {
    assertEq(address(s_lockReleaseTokenPool.getRebalancer()), OWNER);
    s_lockReleaseTokenPool.setRebalancer(STRANGER);
    assertEq(address(s_lockReleaseTokenPool.getRebalancer()), STRANGER);
  }

  function testSetRebalancerReverts() public {
    vm.startPrank(STRANGER);

    vm.expectRevert("Only callable by owner");
    s_lockReleaseTokenPool.setRebalancer(STRANGER);
  }
}

contract LockReleaseTokenPool_lockOrBurn is LockReleaseTokenPoolSetup {
  error SenderNotAllowed(address sender);

  event Locked(address indexed sender, uint256 amount);
  event TokensConsumed(uint256 tokens);

  function testFuzz_LockOrBurnNoAllowListSuccess(uint256 amount) public {
    amount = bound(amount, 1, getOutboundRateLimiterConfig().capacity);
    changePrank(s_allowedOnRamp);

    vm.expectEmit();
    emit TokensConsumed(amount);
    vm.expectEmit();
    emit Locked(s_allowedOnRamp, amount);

    s_lockReleaseTokenPool.lockOrBurn(STRANGER, bytes(""), amount, DEST_CHAIN_SELECTOR, bytes(""));
  }

  function testLockOrBurnWithAllowListSuccess() public {
    uint256 amount = 100;
    changePrank(s_allowedOnRamp);

    vm.expectEmit();
    emit TokensConsumed(amount);
    vm.expectEmit();
    emit Locked(s_allowedOnRamp, amount);

    s_lockReleaseTokenPoolWithAllowList.lockOrBurn(s_allowedList[0], bytes(""), amount, DEST_CHAIN_SELECTOR, bytes(""));

    vm.expectEmit();
    emit Locked(s_allowedOnRamp, amount);

    s_lockReleaseTokenPoolWithAllowList.lockOrBurn(s_allowedList[1], bytes(""), amount, DEST_CHAIN_SELECTOR, bytes(""));
  }

  function testLockOrBurnWithAllowListReverts() public {
    changePrank(s_allowedOnRamp);

    vm.expectRevert(abi.encodeWithSelector(SenderNotAllowed.selector, STRANGER));

    s_lockReleaseTokenPoolWithAllowList.lockOrBurn(STRANGER, bytes(""), 100, DEST_CHAIN_SELECTOR, bytes(""));
  }

  function testPoolBurnRevertNotHealthyReverts() public {
    // Should not burn tokens if cursed.
    s_mockARM.voteToCurse(bytes32(0));
    uint256 before = s_token.balanceOf(address(s_lockReleaseTokenPoolWithAllowList));

    changePrank(s_allowedOnRamp);
    vm.expectRevert(EVM2EVMOnRamp.BadARMSignal.selector);

    s_lockReleaseTokenPoolWithAllowList.lockOrBurn(s_allowedList[0], bytes(""), 1e5, DEST_CHAIN_SELECTOR, bytes(""));

    assertEq(s_token.balanceOf(address(s_lockReleaseTokenPoolWithAllowList)), before);
  }
}

contract LockReleaseTokenPool_releaseOrMint is LockReleaseTokenPoolSetup {
  event TokensConsumed(uint256 tokens);
  event Released(address indexed sender, address indexed recipient, uint256 amount);

  function setUp() public virtual override {
    LockReleaseTokenPoolSetup.setUp();
    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: SOURCE_CHAIN_SELECTOR,
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
    });

    s_lockReleaseTokenPool.applyChainUpdates(chainUpdate);
    s_lockReleaseTokenPoolWithAllowList.applyChainUpdates(chainUpdate);
  }

  function test_ReleaseOrMintSuccess() public {
    vm.startPrank(s_allowedOffRamp);

    uint256 amount = 100;
    deal(address(s_token), address(s_lockReleaseTokenPool), amount);

    vm.expectEmit();
    emit TokensConsumed(amount);
    vm.expectEmit();
    emit Released(s_allowedOffRamp, OWNER, amount);

    s_lockReleaseTokenPool.releaseOrMint(bytes(""), OWNER, amount, SOURCE_CHAIN_SELECTOR, bytes(""));
  }

  function testFuzz_ReleaseOrMintSuccess(address recipient, uint256 amount) public {
    // Since the owner already has tokens this would break the checks
    vm.assume(recipient != OWNER);
    vm.assume(recipient != address(0));
    vm.assume(recipient != address(s_token));

    // Makes sure the pool always has enough funds
    deal(address(s_token), address(s_lockReleaseTokenPool), amount);
    vm.startPrank(s_allowedOffRamp);

    uint256 capacity = getInboundRateLimiterConfig().capacity;
    // Determine if we hit the rate limit or the txs should succeed.
    if (amount > capacity) {
      vm.expectRevert(
        abi.encodeWithSelector(RateLimiter.TokenMaxCapacityExceeded.selector, capacity, amount, address(s_token))
      );
    } else {
      // Only rate limit if the amount is >0
      if (amount > 0) {
        vm.expectEmit();
        emit TokensConsumed(amount);
      }

      vm.expectEmit();
      emit Released(s_allowedOffRamp, recipient, amount);
    }

    s_lockReleaseTokenPool.releaseOrMint(bytes(""), recipient, amount, SOURCE_CHAIN_SELECTOR, bytes(""));
  }

  function testChainNotAllowedReverts() public {
    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: SOURCE_CHAIN_SELECTOR,
      allowed: false,
      outboundRateLimiterConfig: RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0}),
      inboundRateLimiterConfig: RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0})
    });

    s_lockReleaseTokenPool.applyChainUpdates(chainUpdate);

    vm.startPrank(s_allowedOffRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, SOURCE_CHAIN_SELECTOR));
    s_lockReleaseTokenPool.releaseOrMint(bytes(""), OWNER, 1e5, SOURCE_CHAIN_SELECTOR, bytes(""));
  }

  function testPoolMintNotHealthyReverts() public {
    // Should not mint tokens if cursed.
    s_mockARM.voteToCurse(bytes32(0));
    uint256 before = s_token.balanceOf(OWNER);
    vm.startPrank(s_allowedOffRamp);
    vm.expectRevert(EVM2EVMOffRamp.BadARMSignal.selector);
    s_lockReleaseTokenPool.releaseOrMint(bytes(""), OWNER, 1e5, SOURCE_CHAIN_SELECTOR, bytes(""));
    assertEq(s_token.balanceOf(OWNER), before);
  }
}

contract LockReleaseTokenPool_canAcceptLiquidity is LockReleaseTokenPoolSetup {
  function test_CanAcceptLiquiditySuccess() public {
    assertEq(true, s_lockReleaseTokenPool.canAcceptLiquidity());

    s_lockReleaseTokenPool = new LockReleaseTokenPool(
      s_token,
      new address[](0),
      address(s_mockARM),
      false,
      address(s_sourceRouter)
    );
    assertEq(false, s_lockReleaseTokenPool.canAcceptLiquidity());
  }
}

contract LockReleaseTokenPool_provideLiquidity is LockReleaseTokenPoolSetup {
  function testFuzz_ProvideLiquiditySuccess(uint256 amount) public {
    uint256 balancePre = s_token.balanceOf(OWNER);
    s_token.approve(address(s_lockReleaseTokenPool), amount);

    s_lockReleaseTokenPool.provideLiquidity(amount);

    assertEq(s_token.balanceOf(OWNER), balancePre - amount);
    assertEq(s_token.balanceOf(address(s_lockReleaseTokenPool)), amount);
  }

  // Reverts

  function test_UnauthorizedReverts() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(LockReleaseTokenPool.Unauthorized.selector, STRANGER));

    s_lockReleaseTokenPool.provideLiquidity(1);
  }

  function testFuzz_ExceedsAllowance(uint256 amount) public {
    vm.assume(amount > 0);
    vm.expectRevert("ERC20: insufficient allowance");
    s_lockReleaseTokenPool.provideLiquidity(amount);
  }

  function testLiquidityNotAcceptedReverts() public {
    s_lockReleaseTokenPool = new LockReleaseTokenPool(
      s_token,
      new address[](0),
      address(s_mockARM),
      false,
      address(s_sourceRouter)
    );

    vm.expectRevert(LockReleaseTokenPool.LiquidityNotAccepted.selector);
    s_lockReleaseTokenPool.provideLiquidity(1);
  }
}

contract LockReleaseTokenPool_withdrawalLiquidity is LockReleaseTokenPoolSetup {
  function testFuzz_WithdrawalLiquiditySuccess(uint256 amount) public {
    uint256 balancePre = s_token.balanceOf(OWNER);
    s_token.approve(address(s_lockReleaseTokenPool), amount);
    s_lockReleaseTokenPool.provideLiquidity(amount);

    s_lockReleaseTokenPool.withdrawLiquidity(amount);

    assertEq(s_token.balanceOf(OWNER), balancePre);
  }

  // Reverts

  function test_UnauthorizedReverts() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(LockReleaseTokenPool.Unauthorized.selector, STRANGER));

    s_lockReleaseTokenPool.withdrawLiquidity(1);
  }

  function testInsufficientLiquidityReverts() public {
    uint256 maxUint256 = 2 ** 256 - 1;
    s_token.approve(address(s_lockReleaseTokenPool), maxUint256);
    s_lockReleaseTokenPool.provideLiquidity(maxUint256);

    changePrank(address(s_lockReleaseTokenPool));
    s_token.transfer(OWNER, maxUint256);
    changePrank(OWNER);

    vm.expectRevert(LockReleaseTokenPool.InsufficientLiquidity.selector);
    s_lockReleaseTokenPool.withdrawLiquidity(1);
  }
}

contract LockReleaseTokenPool_supportsInterface is LockReleaseTokenPoolSetup {
  function testSupportsInterfaceSuccess() public {
    assertTrue(s_lockReleaseTokenPool.supportsInterface(s_lockReleaseTokenPool.getLockReleaseInterfaceId()));
    assertTrue(s_lockReleaseTokenPool.supportsInterface(type(IPool).interfaceId));
    assertTrue(s_lockReleaseTokenPool.supportsInterface(type(IERC165).interfaceId));
  }
}

contract LockReleaseTokenPool_setChainRateLimiterConfig is LockReleaseTokenPoolSetup {
  event ConfigChanged(RateLimiter.Config);
  event ChainConfigured(
    uint64 chainSelector,
    RateLimiter.Config outboundRateLimiterConfig,
    RateLimiter.Config inboundRateLimiterConfig
  );

  uint64 internal s_remoteChainSelector;

  function setUp() public virtual override {
    LockReleaseTokenPoolSetup.setUp();
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    s_remoteChainSelector = 123124;
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: s_remoteChainSelector,
      allowed: true,
      outboundRateLimiterConfig: getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: getInboundRateLimiterConfig()
    });
    s_lockReleaseTokenPool.applyChainUpdates(chainUpdates);
  }

  function testFuzz_SetChainRateLimiterConfigSuccess(uint128 capacity, uint128 rate, uint32 newTime) public {
    // Cap the lower bound to 4 so 4/2 is still >= 2
    vm.assume(capacity >= 4);
    // Cap the lower bound to 2 so 2/2 is still >= 1
    rate = uint128(bound(rate, 2, capacity - 2));
    // Bucket updates only work on increasing time
    newTime = uint32(bound(newTime, block.timestamp + 1, type(uint32).max));
    vm.warp(newTime);

    uint256 oldOutboundTokens = s_lockReleaseTokenPool.getCurrentOutboundRateLimiterState(s_remoteChainSelector).tokens;
    uint256 oldInboundTokens = s_lockReleaseTokenPool.getCurrentInboundRateLimiterState(s_remoteChainSelector).tokens;

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

    s_lockReleaseTokenPool.setChainRateLimiterConfig(s_remoteChainSelector, newOutboundConfig, newInboundConfig);

    uint256 expectedTokens = RateLimiter._min(newOutboundConfig.capacity, oldOutboundTokens);

    RateLimiter.TokenBucket memory bucket = s_lockReleaseTokenPool.getCurrentOutboundRateLimiterState(
      s_remoteChainSelector
    );
    assertEq(bucket.capacity, newOutboundConfig.capacity);
    assertEq(bucket.rate, newOutboundConfig.rate);
    assertEq(bucket.tokens, expectedTokens);
    assertEq(bucket.lastUpdated, newTime);

    expectedTokens = RateLimiter._min(newInboundConfig.capacity, oldInboundTokens);

    bucket = s_lockReleaseTokenPool.getCurrentInboundRateLimiterState(s_remoteChainSelector);
    assertEq(bucket.capacity, newInboundConfig.capacity);
    assertEq(bucket.rate, newInboundConfig.rate);
    assertEq(bucket.tokens, expectedTokens);
    assertEq(bucket.lastUpdated, newTime);
  }

  function testOnlyOwnerOrRateLimitAdminReverts() public {
    address rateLimiterAdmin = address(28973509103597907);

    s_lockReleaseTokenPool.setRateLimitAdmin(rateLimiterAdmin);

    changePrank(rateLimiterAdmin);

    s_lockReleaseTokenPool.setChainRateLimiterConfig(
      s_remoteChainSelector,
      getOutboundRateLimiterConfig(),
      getInboundRateLimiterConfig()
    );

    changePrank(OWNER);

    s_lockReleaseTokenPool.setChainRateLimiterConfig(
      s_remoteChainSelector,
      getOutboundRateLimiterConfig(),
      getInboundRateLimiterConfig()
    );
  }

  // Reverts

  function testOnlyOwnerReverts() public {
    changePrank(STRANGER);

    vm.expectRevert(abi.encodeWithSelector(LockReleaseTokenPool.Unauthorized.selector, STRANGER));
    s_lockReleaseTokenPool.setChainRateLimiterConfig(
      s_remoteChainSelector,
      getOutboundRateLimiterConfig(),
      getInboundRateLimiterConfig()
    );
  }

  function testNonExistentChainReverts() public {
    uint64 wrongChainSelector = 9084102894;

    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, wrongChainSelector));
    s_lockReleaseTokenPool.setChainRateLimiterConfig(
      wrongChainSelector,
      getOutboundRateLimiterConfig(),
      getInboundRateLimiterConfig()
    );
  }
}

contract LockReleaseTokenPool_setRateLimitAdmin is LockReleaseTokenPoolSetup {
  function testSetRateLimitAdminSuccess() public {
    assertEq(address(0), s_lockReleaseTokenPool.getRateLimitAdmin());
    s_lockReleaseTokenPool.setRateLimitAdmin(OWNER);
    assertEq(OWNER, s_lockReleaseTokenPool.getRateLimitAdmin());
  }

  // Reverts

  function testSetRateLimitAdminReverts() public {
    vm.startPrank(STRANGER);

    vm.expectRevert("Only callable by owner");
    s_lockReleaseTokenPool.setRateLimitAdmin(STRANGER);
  }
}
