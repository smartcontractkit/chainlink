// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IPool} from "../../interfaces/pools/IPool.sol";

import "../BaseTest.t.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/IERC165.sol";

contract LockReleaseTokenPoolSetup is BaseTest {
  IERC20 internal s_token;
  LockReleaseTokenPool internal s_lockReleaseTokenPool;
  LockReleaseTokenPool internal s_lockReleaseTokenPoolWithAllowList;
  address[] internal s_allowedList;

  address internal s_allowedOnRamp = address(123);
  address internal s_allowedOffRamp = address(234);

  function setUp() public virtual override {
    BaseTest.setUp();
    s_token = new BurnMintERC677("LINK", "LNK", 18, 0);
    deal(address(s_token), OWNER, type(uint256).max);
    s_lockReleaseTokenPool = new LockReleaseTokenPool(s_token, new address[](0), address(s_mockARM), true);

    s_allowedList.push(USER_1);
    s_allowedList.push(DUMMY_CONTRACT_ADDRESS);
    s_lockReleaseTokenPoolWithAllowList = new LockReleaseTokenPool(s_token, s_allowedList, address(s_mockARM), true);

    TokenPool.RampUpdate[] memory onRamps = new TokenPool.RampUpdate[](1);
    onRamps[0] = TokenPool.RampUpdate({ramp: s_allowedOnRamp, allowed: true, rateLimiterConfig: rateLimiterConfig()});
    TokenPool.RampUpdate[] memory offRamps = new TokenPool.RampUpdate[](1);
    offRamps[0] = TokenPool.RampUpdate({ramp: s_allowedOffRamp, allowed: true, rateLimiterConfig: rateLimiterConfig()});

    s_lockReleaseTokenPool.applyRampUpdates(onRamps, offRamps);
    s_lockReleaseTokenPoolWithAllowList.applyRampUpdates(onRamps, offRamps);
    s_lockReleaseTokenPool.setRebalancer(OWNER);
  }
}

contract LockReleaseTokenPool_lockOrBurn is LockReleaseTokenPoolSetup {
  error SenderNotAllowed(address sender);

  event Locked(address indexed sender, uint256 amount);
  event TokensConsumed(uint256 tokens);

  function testFuzz_LockOrBurnNoAllowListSuccess(uint256 amount) public {
    amount = bound(amount, 1, rateLimiterConfig().capacity);
    changePrank(s_allowedOnRamp);

    vm.expectEmit();
    emit TokensConsumed(amount);
    vm.expectEmit();
    emit Locked(s_allowedOnRamp, amount);

    s_lockReleaseTokenPool.lockOrBurn(STRANGER, bytes(""), amount, DEST_CHAIN_ID, bytes(""));
  }

  function testLockOrBurnWithAllowListSuccess() public {
    uint256 amount = 100;
    changePrank(s_allowedOnRamp);

    vm.expectEmit();
    emit TokensConsumed(amount);
    vm.expectEmit();
    emit Locked(s_allowedOnRamp, amount);

    s_lockReleaseTokenPoolWithAllowList.lockOrBurn(s_allowedList[0], bytes(""), amount, DEST_CHAIN_ID, bytes(""));

    vm.expectEmit();
    emit Locked(s_allowedOnRamp, amount);

    s_lockReleaseTokenPoolWithAllowList.lockOrBurn(s_allowedList[1], bytes(""), amount, DEST_CHAIN_ID, bytes(""));
  }

  function testLockOrBurnWithAllowListReverts() public {
    changePrank(s_allowedOnRamp);

    vm.expectRevert(abi.encodeWithSelector(SenderNotAllowed.selector, STRANGER));

    s_lockReleaseTokenPoolWithAllowList.lockOrBurn(STRANGER, bytes(""), 100, DEST_CHAIN_ID, bytes(""));
  }

  function testPoolBurnRevertNotHealthyReverts() public {
    // Should not burn tokens if cursed.
    s_mockARM.voteToCurse(bytes32(0));
    uint256 before = s_token.balanceOf(address(s_lockReleaseTokenPoolWithAllowList));
    changePrank(s_allowedOnRamp);
    vm.expectRevert(EVM2EVMOnRamp.BadARMSignal.selector);
    s_lockReleaseTokenPoolWithAllowList.lockOrBurn(s_allowedList[0], bytes(""), 1e5, 0, bytes(""));
    assertEq(s_token.balanceOf(address(s_lockReleaseTokenPoolWithAllowList)), before);
  }
}

contract LockReleaseTokenPool_releaseOrMint is LockReleaseTokenPoolSetup {
  event TokensConsumed(uint256 tokens);
  event Released(address indexed sender, address indexed recipient, uint256 amount);

  function testFuzz_ReleaseOrMintSuccess(address recipient, uint256 amount) public {
    // Since the owner already has tokens this would break the checks
    vm.assume(recipient != OWNER);
    vm.assume(recipient != address(0));
    vm.assume(recipient != address(s_token));

    // Makes sure the pool always has enough funds
    deal(address(s_token), address(s_lockReleaseTokenPool), amount);
    changePrank(s_allowedOffRamp);

    uint256 capacity = rateLimiterConfig().capacity;
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

    s_lockReleaseTokenPool.releaseOrMint(bytes(""), recipient, amount, SOURCE_CHAIN_ID, bytes(""));
  }

  function testPoolMintNotHealthyReverts() public {
    // Should not mint tokens if cursed.
    s_mockARM.voteToCurse(bytes32(0));
    uint256 before = s_token.balanceOf(OWNER);
    changePrank(s_allowedOffRamp);
    vm.expectRevert(EVM2EVMOffRamp.BadARMSignal.selector);
    s_lockReleaseTokenPool.releaseOrMint(bytes(""), OWNER, 1e5, 0, bytes(""));
    assertEq(s_token.balanceOf(OWNER), before);
  }
}

contract LockReleaseTokenPool_canAcceptLiquidity is LockReleaseTokenPoolSetup {
  function test_CanAcceptLiquiditySuccess() public {
    assertEq(true, s_lockReleaseTokenPool.canAcceptLiquidity());

    s_lockReleaseTokenPool = new LockReleaseTokenPool(s_token, new address[](0), address(s_mockARM), false);
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

  function testFuzz_ExceedsAllowance(uint256 amount) public {
    vm.assume(amount > 0);
    vm.expectRevert("ERC20: insufficient allowance");
    s_lockReleaseTokenPool.provideLiquidity(amount);
  }

  function testLiquidityNotAcceptedReverts() public {
    s_lockReleaseTokenPool = new LockReleaseTokenPool(s_token, new address[](0), address(s_mockARM), false);

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
