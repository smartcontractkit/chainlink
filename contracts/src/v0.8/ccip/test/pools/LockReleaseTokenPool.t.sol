// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IPoolV1} from "../../interfaces/IPool.sol";

import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {Router} from "../../Router.sol";
import {Pool} from "../../libraries/Pool.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {IERC165} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";
import {RouterSetup} from "../router/RouterSetup.t.sol";

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
    s_allowedList.push(DUMMY_CONTRACT_ADDRESS);
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

contract LockReleaseTokenPool_setRebalancer is LockReleaseTokenPoolSetup {
  function test_SetRebalancer_Success() public {
    assertEq(address(s_lockReleaseTokenPool.getRebalancer()), OWNER);
    s_lockReleaseTokenPool.setRebalancer(STRANGER);
    assertEq(address(s_lockReleaseTokenPool.getRebalancer()), STRANGER);
  }

  function test_SetRebalancer_Revert() public {
    vm.startPrank(STRANGER);

    vm.expectRevert("Only callable by owner");
    s_lockReleaseTokenPool.setRebalancer(STRANGER);
  }
}

contract LockReleaseTokenPool_lockOrBurn is LockReleaseTokenPoolSetup {
  function test_Fuzz_LockOrBurnNoAllowList_Success(
    uint256 amount
  ) public {
    amount = bound(amount, 1, _getOutboundRateLimiterConfig().capacity);
    vm.startPrank(s_allowedOnRamp);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(amount);
    vm.expectEmit();
    emit TokenPool.Locked(s_allowedOnRamp, amount);

    s_lockReleaseTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: STRANGER,
        receiver: bytes(""),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );
  }

  function test_LockOrBurnWithAllowList_Success() public {
    uint256 amount = 100;
    vm.startPrank(s_allowedOnRamp);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(amount);
    vm.expectEmit();
    emit TokenPool.Locked(s_allowedOnRamp, amount);

    s_lockReleaseTokenPoolWithAllowList.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: s_allowedList[0],
        receiver: bytes(""),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    vm.expectEmit();
    emit TokenPool.Locked(s_allowedOnRamp, amount);

    s_lockReleaseTokenPoolWithAllowList.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: s_allowedList[1],
        receiver: bytes(""),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );
  }

  function test_LockOrBurnWithAllowList_Revert() public {
    vm.startPrank(s_allowedOnRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.SenderNotAllowed.selector, STRANGER));

    s_lockReleaseTokenPoolWithAllowList.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: STRANGER,
        receiver: bytes(""),
        amount: 100,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );
  }

  function test_PoolBurnRevertNotHealthy_Revert() public {
    // Should not burn tokens if cursed.
    s_mockRMN.setGlobalCursed(true);
    uint256 before = s_token.balanceOf(address(s_lockReleaseTokenPoolWithAllowList));

    vm.startPrank(s_allowedOnRamp);
    vm.expectRevert(TokenPool.CursedByRMN.selector);

    s_lockReleaseTokenPoolWithAllowList.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: s_allowedList[0],
        receiver: bytes(""),
        amount: 1e5,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    assertEq(s_token.balanceOf(address(s_lockReleaseTokenPoolWithAllowList)), before);
  }
}

contract LockReleaseTokenPool_releaseOrMint is LockReleaseTokenPoolSetup {
  function setUp() public virtual override {
    LockReleaseTokenPoolSetup.setUp();
    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: SOURCE_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(s_sourcePoolAddress),
      remoteTokenAddress: abi.encode(address(2)),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });

    s_lockReleaseTokenPool.applyChainUpdates(chainUpdate);
    s_lockReleaseTokenPoolWithAllowList.applyChainUpdates(chainUpdate);
  }

  function test_ReleaseOrMint_Success() public {
    vm.startPrank(s_allowedOffRamp);

    uint256 amount = 100;
    deal(address(s_token), address(s_lockReleaseTokenPool), amount);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(amount);
    vm.expectEmit();
    emit TokenPool.Released(s_allowedOffRamp, OWNER, amount);

    s_lockReleaseTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: OWNER,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_sourcePoolAddress),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );
  }

  function test_Fuzz_ReleaseOrMint_Success(address recipient, uint256 amount) public {
    // Since the owner already has tokens this would break the checks
    vm.assume(recipient != OWNER);
    vm.assume(recipient != address(0));
    vm.assume(recipient != address(s_token));

    // Makes sure the pool always has enough funds
    deal(address(s_token), address(s_lockReleaseTokenPool), amount);
    vm.startPrank(s_allowedOffRamp);

    uint256 capacity = _getInboundRateLimiterConfig().capacity;
    // Determine if we hit the rate limit or the txs should succeed.
    if (amount > capacity) {
      vm.expectRevert(
        abi.encodeWithSelector(RateLimiter.TokenMaxCapacityExceeded.selector, capacity, amount, address(s_token))
      );
    } else {
      // Only rate limit if the amount is >0
      if (amount > 0) {
        vm.expectEmit();
        emit RateLimiter.TokensConsumed(amount);
      }

      vm.expectEmit();
      emit TokenPool.Released(s_allowedOffRamp, recipient, amount);
    }

    s_lockReleaseTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: recipient,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_sourcePoolAddress),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );
  }

  function test_ChainNotAllowed_Revert() public {
    address notAllowedRemotePoolAddress = address(1);

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: SOURCE_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(notAllowedRemotePoolAddress),
      remoteTokenAddress: abi.encode(address(2)),
      allowed: false,
      outboundRateLimiterConfig: RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0}),
      inboundRateLimiterConfig: RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0})
    });

    s_lockReleaseTokenPool.applyChainUpdates(chainUpdate);

    vm.startPrank(s_allowedOffRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, SOURCE_CHAIN_SELECTOR));
    s_lockReleaseTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: OWNER,
        amount: 1e5,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_sourcePoolAddress),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );
  }

  function test_PoolMintNotHealthy_Revert() public {
    // Should not mint tokens if cursed.
    s_mockRMN.setGlobalCursed(true);
    uint256 before = s_token.balanceOf(OWNER);
    vm.startPrank(s_allowedOffRamp);
    vm.expectRevert(TokenPool.CursedByRMN.selector);
    s_lockReleaseTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: OWNER,
        amount: 1e5,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: _generateSourceTokenData().sourcePoolAddress,
        sourcePoolData: _generateSourceTokenData().extraData,
        offchainTokenData: ""
      })
    );

    assertEq(s_token.balanceOf(OWNER), before);
  }
}

contract LockReleaseTokenPool_canAcceptLiquidity is LockReleaseTokenPoolSetup {
  function test_CanAcceptLiquidity_Success() public {
    assertEq(true, s_lockReleaseTokenPool.canAcceptLiquidity());

    s_lockReleaseTokenPool =
      new LockReleaseTokenPool(s_token, new address[](0), address(s_mockRMN), false, address(s_sourceRouter));
    assertEq(false, s_lockReleaseTokenPool.canAcceptLiquidity());
  }
}

contract LockReleaseTokenPool_provideLiquidity is LockReleaseTokenPoolSetup {
  function test_Fuzz_ProvideLiquidity_Success(
    uint256 amount
  ) public {
    uint256 balancePre = s_token.balanceOf(OWNER);
    s_token.approve(address(s_lockReleaseTokenPool), amount);

    s_lockReleaseTokenPool.provideLiquidity(amount);

    assertEq(s_token.balanceOf(OWNER), balancePre - amount);
    assertEq(s_token.balanceOf(address(s_lockReleaseTokenPool)), amount);
  }

  // Reverts

  function test_Unauthorized_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, STRANGER));

    s_lockReleaseTokenPool.provideLiquidity(1);
  }

  function test_Fuzz_ExceedsAllowance(
    uint256 amount
  ) public {
    vm.assume(amount > 0);
    vm.expectRevert("ERC20: insufficient allowance");
    s_lockReleaseTokenPool.provideLiquidity(amount);
  }

  function test_LiquidityNotAccepted_Revert() public {
    s_lockReleaseTokenPool =
      new LockReleaseTokenPool(s_token, new address[](0), address(s_mockRMN), false, address(s_sourceRouter));

    vm.expectRevert(LockReleaseTokenPool.LiquidityNotAccepted.selector);
    s_lockReleaseTokenPool.provideLiquidity(1);
  }
}

contract LockReleaseTokenPool_withdrawalLiquidity is LockReleaseTokenPoolSetup {
  function test_Fuzz_WithdrawalLiquidity_Success(
    uint256 amount
  ) public {
    uint256 balancePre = s_token.balanceOf(OWNER);
    s_token.approve(address(s_lockReleaseTokenPool), amount);
    s_lockReleaseTokenPool.provideLiquidity(amount);

    s_lockReleaseTokenPool.withdrawLiquidity(amount);

    assertEq(s_token.balanceOf(OWNER), balancePre);
  }

  // Reverts

  function test_Unauthorized_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, STRANGER));

    s_lockReleaseTokenPool.withdrawLiquidity(1);
  }

  function test_InsufficientLiquidity_Revert() public {
    uint256 maxUint256 = 2 ** 256 - 1;
    s_token.approve(address(s_lockReleaseTokenPool), maxUint256);
    s_lockReleaseTokenPool.provideLiquidity(maxUint256);

    vm.startPrank(address(s_lockReleaseTokenPool));
    s_token.transfer(OWNER, maxUint256);
    vm.startPrank(OWNER);

    vm.expectRevert(LockReleaseTokenPool.InsufficientLiquidity.selector);
    s_lockReleaseTokenPool.withdrawLiquidity(1);
  }
}

contract LockReleaseTokenPool_transferLiquidity is LockReleaseTokenPoolSetup {
  LockReleaseTokenPool internal s_oldLockReleaseTokenPool;
  uint256 internal s_amount = 100000;

  function setUp() public virtual override {
    super.setUp();

    s_oldLockReleaseTokenPool =
      new LockReleaseTokenPool(s_token, new address[](0), address(s_mockRMN), true, address(s_sourceRouter));

    deal(address(s_token), address(s_oldLockReleaseTokenPool), s_amount);
  }

  function test_transferLiquidity_Success() public {
    uint256 balancePre = s_token.balanceOf(address(s_lockReleaseTokenPool));

    s_oldLockReleaseTokenPool.setRebalancer(address(s_lockReleaseTokenPool));

    vm.expectEmit();
    emit LockReleaseTokenPool.LiquidityTransferred(address(s_oldLockReleaseTokenPool), s_amount);

    s_lockReleaseTokenPool.transferLiquidity(address(s_oldLockReleaseTokenPool), s_amount);

    assertEq(s_token.balanceOf(address(s_lockReleaseTokenPool)), balancePre + s_amount);
  }

  function test_transferLiquidity_transferTooMuch_Revert() public {
    uint256 balancePre = s_token.balanceOf(address(s_lockReleaseTokenPool));

    s_oldLockReleaseTokenPool.setRebalancer(address(s_lockReleaseTokenPool));

    vm.expectRevert(LockReleaseTokenPool.InsufficientLiquidity.selector);
    s_lockReleaseTokenPool.transferLiquidity(address(s_oldLockReleaseTokenPool), s_amount + 1);

    assertEq(s_token.balanceOf(address(s_lockReleaseTokenPool)), balancePre);
  }
}

contract LockReleaseTokenPool_supportsInterface is LockReleaseTokenPoolSetup {
  function test_SupportsInterface_Success() public view {
    assertTrue(s_lockReleaseTokenPool.supportsInterface(type(IPoolV1).interfaceId));
    assertTrue(s_lockReleaseTokenPool.supportsInterface(type(IERC165).interfaceId));
  }
}
