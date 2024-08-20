// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {IPoolV1} from "../../interfaces/IPool.sol";
import {IPoolPriorTo1_5} from "../../interfaces/IPoolPriorTo1_5.sol";

import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {Router} from "../../Router.sol";
import {Client} from "../../libraries/Client.sol";
import {Pool} from "../../libraries/Pool.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {BurnMintTokenPoolAndProxy} from "../../pools/BurnMintTokenPoolAndProxy.sol";
import {BurnWithFromMintTokenPoolAndProxy} from "../../pools/BurnWithFromMintTokenPoolAndProxy.sol";
import {LockReleaseTokenPoolAndProxy} from "../../pools/LockReleaseTokenPoolAndProxy.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {EVM2EVMOnRampSetup} from "../onRamp/EVM2EVMOnRampSetup.t.sol";
import {RouterSetup} from "../router/RouterSetup.t.sol";
import {BurnMintTokenPool1_2, TokenPool1_2} from "./BurnMintTokenPool1_2.sol";
import {BurnMintTokenPool1_4, TokenPool1_4} from "./BurnMintTokenPool1_4.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {IERC165} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";

contract TokenPoolAndProxyMigration is EVM2EVMOnRampSetup {
  BurnMintTokenPoolAndProxy internal s_newPool;
  IPoolPriorTo1_5 internal s_legacyPool;
  BurnMintERC677 internal s_token;

  address internal s_offRamp;
  address internal s_sourcePool = makeAddr("source_pool");
  address internal s_sourceToken = makeAddr("source_token");
  uint256 internal constant AMOUNT = 1;

  function setUp() public virtual override {
    super.setUp();
    // Create a system with a token and a legacy pool
    s_token = new BurnMintERC677("Test", "TEST", 18, type(uint256).max);
    // dealing doesn't update the total supply, meaning the first time we burn a token we underflow, which isn't
    // guarded against. Then, when we mint a token, we overflow, which is guarded against and will revert.
    s_token.grantMintAndBurnRoles(OWNER);
    s_token.mint(OWNER, 1e18);

    s_offRamp = s_offRamps[0];
    // Approve enough for a few calls
    s_token.approve(address(s_sourceRouter), AMOUNT * 100);

    // Approve infinite fee tokens
    IERC20(s_sourceFeeToken).approve(address(s_sourceRouter), type(uint256).max);
  }

  /// @notice This test covers the entire migration plan for 1.0-1.2 pools to 1.5 pools. For simplicity
  /// we will refer to the 1.0/1.2 pools as 1.2 pools, as they are functionally the same.
  function test_tokenPoolMigration_Success_1_2() public {
    // ================================================================
    // |          1          1.2 prior to upgrade                     |
    // ================================================================
    _deployPool1_2();

    // Ensure everything works on the 1.2 pool
    _ccipSend_OLD();
    _fakeReleaseOrMintFromOffRamp_OLD();

    // ================================================================
    // |          2           Deploy self serve                       |
    // ================================================================
    _deploySelfServe();

    // This doesn't impact the 1.2 pool, so it should still be functional
    _ccipSend_OLD();
    _fakeReleaseOrMintFromOffRamp_OLD();

    // ================================================================
    // |          3     Configure new pool on old pool                |
    // ================================================================
    // In the 1.2 case, everything keeps working on both the 1.2 and 1.5 pools. This config can be
    // done in advance of the actual swap to 1.5 lanes.
    vm.startPrank(OWNER);
    TokenPool1_2.RampUpdate[] memory rampUpdates = new TokenPool1_2.RampUpdate[](1);
    rampUpdates[0] = TokenPool1_2.RampUpdate({
      ramp: address(s_newPool),
      allowed: true,
      // The rate limits should be turned off for this fake ramp, as the 1.5 pool will handle all the
      // rate limiting for us.
      rateLimiterConfig: RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0})
    });
    // Since this call doesn't impact the usability of the old pool, we can do it whenever we want
    BurnMintTokenPool1_2(address(s_legacyPool)).applyRampUpdates(rampUpdates, rampUpdates);

    // Assert the 1.2 lanes still work
    _ccipSend_OLD();
    _fakeReleaseOrMintFromOffRamp_OLD();

    // ================================================================
    // |          4     Update the router with to 1.5                 |
    // ================================================================

    // This will stop any new messages entering the old lanes, and will direct all traffic to the
    // new 1.5 lanes, and therefore to the 1.5 pools. Note that the old pools will still receive
    // inflight messages, and will need to continue functioning until all of those are processed.
    _fakeReleaseOrMintFromOffRamp_OLD();

    // Everything is configured, we can now send a ccip tx to the new pool
    _ccipSend1_5();
    _fakeReleaseOrMintFromOffRamp1_5();

    // ================================================================
    // |          5      Migrate to using 1.5 the pool                |
    // ================================================================
    // Turn off the legacy pool, this enabled the 1.5 pool logic. This should be done AFTER the new pool
    // has gotten permissions to mint/burn. We see the case where that isn't done below.
    vm.startPrank(OWNER);
    s_newPool.setPreviousPool(IPoolPriorTo1_5(address(0)));

    // The new pool is now active, but is has not been given permissions to burn/mint yet
    vm.expectRevert(abi.encodeWithSelector(BurnMintERC677.SenderNotBurner.selector, address(s_newPool)));
    _ccipSend1_5();
    vm.expectRevert(abi.encodeWithSelector(BurnMintERC677.SenderNotMinter.selector, address(s_newPool)));
    _fakeReleaseOrMintFromOffRamp1_5();

    // When we do give burn/mint, the new pool is fully active
    vm.startPrank(OWNER);
    s_token.grantMintAndBurnRoles(address(s_newPool));
    _ccipSend1_5();
    _fakeReleaseOrMintFromOffRamp1_5();

    // Even after the pool has taken over as primary, the old pool can still process messages from the old lane
    _fakeReleaseOrMintFromOffRamp_OLD();
  }

  function test_tokenPoolMigration_Success_1_4() public {
    // ================================================================
    // |          1          1.4 prior to upgrade                     |
    // ================================================================
    _deployPool1_4();

    // Ensure everything works on the 1.4 pool
    _ccipSend_OLD();
    _fakeReleaseOrMintFromOffRamp_OLD();

    // ================================================================
    // |          2           Deploy self serve                       |
    // ================================================================
    _deploySelfServe();

    // This doesn't impact the 1.4 pool, so it should still be functional
    _ccipSend_OLD();
    _fakeReleaseOrMintFromOffRamp_OLD();

    // ================================================================
    // |          3     Configure new pool on old pool                |
    // |                           AND                                |
    // |                Update the router with to 1.5                 |
    // ================================================================
    // NOTE: when this call is made, the SENDING SIDE of old lanes stop working.
    vm.startPrank(OWNER);
    BurnMintTokenPool1_4(address(s_legacyPool)).setRouter(address(s_newPool));

    // This will stop any new messages entering the old lanes, and will direct all traffic to the
    // new 1.5 lanes, and therefore to the 1.5 pools. Note that the old pools will still receive
    // inflight messages, and will need to continue functioning until all of those are processed.
    _fakeReleaseOrMintFromOffRamp_OLD();

    // Sending to the old 1.4 pool no longer works
    _ccipSend_OLD_Reverts();

    // Everything is configured, we can now send a ccip tx
    _ccipSend1_5();
    _fakeReleaseOrMintFromOffRamp1_5();

    // ================================================================
    // |          4      Migrate to using 1.5 the pool                |
    // ================================================================
    // Turn off the legacy pool, this enabled the 1.5 pool logic. This should be done AFTER the new pool
    // has gotten permissions to mint/burn. We see the case where that isn't done below.
    vm.startPrank(OWNER);
    s_newPool.setPreviousPool(IPoolPriorTo1_5(address(0)));

    // The new pool is now active, but is has not been given permissions to burn/mint yet
    vm.expectRevert(abi.encodeWithSelector(BurnMintERC677.SenderNotBurner.selector, address(s_newPool)));
    _ccipSend1_5();
    vm.expectRevert(abi.encodeWithSelector(BurnMintERC677.SenderNotMinter.selector, address(s_newPool)));
    _fakeReleaseOrMintFromOffRamp1_5();

    // When we do give burn/mint, the new pool is fully active
    vm.startPrank(OWNER);
    s_token.grantMintAndBurnRoles(address(s_newPool));
    _ccipSend1_5();
    _fakeReleaseOrMintFromOffRamp1_5();

    // Even after the pool has taken over as primary, the old pool can still process messages from the old lane
    _fakeReleaseOrMintFromOffRamp_OLD();
  }

  function _ccipSend_OLD() internal {
    // We send the funds to the pool manually, as the ramp normally does that
    deal(address(s_token), address(s_legacyPool), AMOUNT);
    vm.startPrank(address(s_onRamp));
    s_legacyPool.lockOrBurn(OWNER, abi.encode(OWNER), AMOUNT, DEST_CHAIN_SELECTOR, "");
  }

  function _ccipSend_OLD_Reverts() internal {
    // We send the funds to the pool manually, as the ramp normally does that
    deal(address(s_token), address(s_legacyPool), AMOUNT);
    vm.startPrank(address(s_onRamp));

    vm.expectRevert(abi.encodeWithSelector(TokenPool1_4.CallerIsNotARampOnRouter.selector, address(s_onRamp)));

    s_legacyPool.lockOrBurn(OWNER, abi.encode(OWNER), AMOUNT, DEST_CHAIN_SELECTOR, "");
  }

  function _ccipSend1_5() internal {
    vm.startPrank(address(OWNER));
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({token: address(s_token), amount: AMOUNT});

    s_sourceRouter.ccipSend(
      DEST_CHAIN_SELECTOR,
      Client.EVM2AnyMessage({
        receiver: abi.encode(OWNER),
        data: "",
        tokenAmounts: tokenAmounts,
        feeToken: s_sourceFeeToken,
        extraArgs: ""
      })
    );
  }

  function _fakeReleaseOrMintFromOffRamp1_5() internal {
    // This is a fake call to simulate the release or mint from the "offRamp"
    vm.startPrank(s_offRamp);
    s_newPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        receiver: OWNER,
        amount: AMOUNT,
        localToken: address(s_token),
        sourcePoolAddress: abi.encode(s_sourcePool),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );
  }

  function _fakeReleaseOrMintFromOffRamp_OLD() internal {
    // This is a fake call to simulate the release or mint from the "offRamp"
    vm.startPrank(s_offRamp);
    s_legacyPool.releaseOrMint(abi.encode(OWNER), OWNER, AMOUNT, SOURCE_CHAIN_SELECTOR, "");
  }

  function _deployPool1_2() internal {
    vm.startPrank(OWNER);
    s_legacyPool = new BurnMintTokenPool1_2(s_token, new address[](0), address(s_mockRMN));
    s_token.grantMintAndBurnRoles(address(s_legacyPool));

    TokenPool1_2.RampUpdate[] memory onRampUpdates = new TokenPool1_2.RampUpdate[](1);
    onRampUpdates[0] = TokenPool1_2.RampUpdate({
      ramp: address(s_onRamp),
      allowed: true,
      rateLimiterConfig: _getInboundRateLimiterConfig()
    });
    TokenPool1_2.RampUpdate[] memory offRampUpdates = new TokenPool1_2.RampUpdate[](1);
    offRampUpdates[0] = TokenPool1_2.RampUpdate({
      ramp: address(s_offRamp),
      allowed: true,
      rateLimiterConfig: _getInboundRateLimiterConfig()
    });
    BurnMintTokenPool1_2(address(s_legacyPool)).applyRampUpdates(onRampUpdates, offRampUpdates);
  }

  function _deployPool1_4() internal {
    vm.startPrank(OWNER);
    s_legacyPool = new BurnMintTokenPool1_4(s_token, new address[](0), address(s_mockRMN), address(s_sourceRouter));
    s_token.grantMintAndBurnRoles(address(s_legacyPool));

    TokenPool1_4.ChainUpdate[] memory legacyChainUpdates = new TokenPool1_4.ChainUpdate[](2);
    legacyChainUpdates[0] = TokenPool1_4.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    legacyChainUpdates[1] = TokenPool1_4.ChainUpdate({
      remoteChainSelector: SOURCE_CHAIN_SELECTOR,
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    BurnMintTokenPool1_4(address(s_legacyPool)).applyChainUpdates(legacyChainUpdates);
  }

  function _deploySelfServe() internal {
    vm.startPrank(OWNER);
    // Deploy the new pool
    s_newPool = new BurnMintTokenPoolAndProxy(s_token, new address[](0), address(s_mockRMN), address(s_sourceRouter));
    // Set the previous pool on the new pool
    s_newPool.setPreviousPool(s_legacyPool);

    // Configure the lanes just like the legacy pool
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](2);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(s_destTokenPool),
      remoteTokenAddress: abi.encode(s_destToken),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    chainUpdates[1] = TokenPool.ChainUpdate({
      remoteChainSelector: SOURCE_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(s_sourcePool),
      remoteTokenAddress: abi.encode(s_sourceToken),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_newPool.applyChainUpdates(chainUpdates);

    // Register the token on the token admin registry
    s_tokenAdminRegistry.proposeAdministrator(address(s_token), OWNER);
    // Accept ownership of the token
    s_tokenAdminRegistry.acceptAdminRole(address(s_token));
    // Set the pool on the admin registry
    s_tokenAdminRegistry.setPool(address(s_token), address(s_newPool));
  }
}

contract TokenPoolAndProxy is EVM2EVMOnRampSetup {
  event Burned(address indexed sender, uint256 amount);
  event Minted(address indexed sender, address indexed recipient, uint256 amount);

  IPoolV1 internal s_pool;
  BurnMintERC677 internal s_token;
  IPoolPriorTo1_5 internal s_legacyPool;
  address internal s_fakeOffRamp = makeAddr("off_ramp");

  address internal s_destPool = makeAddr("dest_pool");

  function setUp() public virtual override {
    super.setUp();
    s_token = BurnMintERC677(s_sourceFeeToken);

    Router.OffRamp[] memory fakeOffRamps = new Router.OffRamp[](1);
    fakeOffRamps[0] = Router.OffRamp({sourceChainSelector: DEST_CHAIN_SELECTOR, offRamp: s_fakeOffRamp});
    s_sourceRouter.applyRampUpdates(new Router.OnRamp[](0), new Router.OffRamp[](0), fakeOffRamps);

    s_token.grantMintAndBurnRoles(OWNER);
    s_token.mint(OWNER, 1e18);
  }

  function test_lockOrBurn_burnMint_Success() public {
    s_pool = new BurnMintTokenPoolAndProxy(s_token, new address[](0), address(s_mockRMN), address(s_sourceRouter));
    _configurePool();
    _deployOldPool();
    _assertLockOrBurnCorrect();

    vm.startPrank(OWNER);
    BurnMintTokenPoolAndProxy(address(s_pool)).setPreviousPool(IPoolPriorTo1_5(address(0)));

    _assertReleaseOrMintCorrect();
  }

  function test_lockOrBurn_burnWithFromMint_Success() public {
    s_pool =
      new BurnWithFromMintTokenPoolAndProxy(s_token, new address[](0), address(s_mockRMN), address(s_sourceRouter));
    _configurePool();
    _deployOldPool();
    _assertLockOrBurnCorrect();

    vm.startPrank(OWNER);
    BurnMintTokenPoolAndProxy(address(s_pool)).setPreviousPool(IPoolPriorTo1_5(address(0)));

    _assertReleaseOrMintCorrect();
  }

  function test_lockOrBurn_lockRelease_Success() public {
    s_pool =
      new LockReleaseTokenPoolAndProxy(s_token, new address[](0), address(s_mockRMN), false, address(s_sourceRouter));
    _configurePool();
    _deployOldPool();
    _assertLockOrBurnCorrect();

    vm.startPrank(OWNER);
    BurnMintTokenPoolAndProxy(address(s_pool)).setPreviousPool(IPoolPriorTo1_5(address(0)));

    _assertReleaseOrMintCorrect();
  }

  function _deployOldPool() internal {
    s_legacyPool = new BurnMintTokenPool1_2(s_token, new address[](0), address(s_mockRMN));
    s_token.grantMintAndBurnRoles(address(s_legacyPool));

    TokenPool1_2.RampUpdate[] memory onRampUpdates = new TokenPool1_2.RampUpdate[](1);
    onRampUpdates[0] =
      TokenPool1_2.RampUpdate({ramp: address(s_pool), allowed: true, rateLimiterConfig: _getInboundRateLimiterConfig()});
    TokenPool1_2.RampUpdate[] memory offRampUpdates = new TokenPool1_2.RampUpdate[](1);
    offRampUpdates[0] =
      TokenPool1_2.RampUpdate({ramp: address(s_pool), allowed: true, rateLimiterConfig: _getInboundRateLimiterConfig()});
    BurnMintTokenPool1_2(address(s_legacyPool)).applyRampUpdates(onRampUpdates, offRampUpdates);
  }

  function _configurePool() internal {
    TokenPool.ChainUpdate[] memory chains = new TokenPool.ChainUpdate[](1);
    chains[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(s_destPool),
      remoteTokenAddress: abi.encode(s_destToken),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });

    BurnMintTokenPoolAndProxy(address(s_pool)).applyChainUpdates(chains);

    // CCIP Token Admin has already been registered from TokenSetup
    s_tokenAdminRegistry.setPool(address(s_token), address(s_pool));

    s_token.grantMintAndBurnRoles(address(s_pool));
  }

  function _assertLockOrBurnCorrect() internal {
    uint256 amount = 1234;
    vm.startPrank(address(s_onRamp));

    // lockOrBurn, assert normal path is taken
    deal(address(s_token), address(s_pool), amount);

    s_pool.lockOrBurn(
      Pool.LockOrBurnInV1({
        receiver: abi.encode(OWNER),
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        originalSender: OWNER,
        amount: amount,
        localToken: address(s_token)
      })
    );

    // set legacy pool

    vm.startPrank(OWNER);
    BurnMintTokenPoolAndProxy(address(s_pool)).setPreviousPool(s_legacyPool);

    // lockOrBurn, assert legacy pool is called

    vm.startPrank(address(s_onRamp));
    deal(address(s_token), address(s_pool), amount);

    vm.expectEmit(address(s_legacyPool));
    emit Burned(address(s_pool), amount);

    s_pool.lockOrBurn(
      Pool.LockOrBurnInV1({
        receiver: abi.encode(OWNER),
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        originalSender: OWNER,
        amount: amount,
        localToken: address(s_token)
      })
    );
  }

  function _assertReleaseOrMintCorrect() internal {
    uint256 amount = 1234;
    vm.startPrank(s_fakeOffRamp);

    // releaseOrMint, assert normal path is taken
    deal(address(s_token), address(s_pool), amount);

    s_pool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        receiver: OWNER,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        originalSender: abi.encode(OWNER),
        amount: amount,
        localToken: address(s_token),
        sourcePoolAddress: abi.encode(s_destPool),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );

    // set legacy pool

    vm.startPrank(OWNER);
    BurnMintTokenPoolAndProxy(address(s_pool)).setPreviousPool(s_legacyPool);

    // releaseOrMint, assert legacy pool is called

    vm.startPrank(address(s_fakeOffRamp));

    vm.expectEmit(address(s_legacyPool));
    emit Minted(address(s_pool), address(OWNER), amount);

    s_pool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        receiver: OWNER,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        originalSender: abi.encode(OWNER),
        amount: amount,
        localToken: address(s_token),
        sourcePoolAddress: abi.encode(s_destPool),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );
  }

  function test_setPreviousPool_Success() public {
    LockReleaseTokenPoolAndProxy pool =
      new LockReleaseTokenPoolAndProxy(s_token, new address[](0), address(s_mockRMN), true, address(s_sourceRouter));

    assertEq(pool.getPreviousPool(), address(0));

    address newLegacyPool = makeAddr("new_legacy_pool");

    vm.startPrank(OWNER);
    pool.setPreviousPool(IPoolPriorTo1_5(newLegacyPool));
    assertEq(pool.getPreviousPool(), address(newLegacyPool));
  }
}

////
/// Duplicated tests from LockReleaseTokenPool.t.sol
///

contract LockReleaseTokenPoolAndProxySetup is RouterSetup {
  IERC20 internal s_token;
  LockReleaseTokenPoolAndProxy internal s_lockReleaseTokenPoolAndProxy;
  LockReleaseTokenPoolAndProxy internal s_lockReleaseTokenPoolAndProxyWithAllowList;
  address[] internal s_allowedList;

  address internal s_allowedOnRamp = address(123);
  address internal s_allowedOffRamp = address(234);

  address internal s_destPoolAddress = address(2736782345);
  address internal s_sourcePoolAddress = address(53852352095);

  function setUp() public virtual override {
    RouterSetup.setUp();
    s_token = new BurnMintERC677("LINK", "LNK", 18, 0);
    deal(address(s_token), OWNER, type(uint256).max);
    s_lockReleaseTokenPoolAndProxy =
      new LockReleaseTokenPoolAndProxy(s_token, new address[](0), address(s_mockRMN), true, address(s_sourceRouter));

    s_allowedList.push(USER_1);
    s_allowedList.push(DUMMY_CONTRACT_ADDRESS);
    s_lockReleaseTokenPoolAndProxyWithAllowList =
      new LockReleaseTokenPoolAndProxy(s_token, s_allowedList, address(s_mockRMN), true, address(s_sourceRouter));

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(s_destPoolAddress),
      remoteTokenAddress: abi.encode(address(s_token)),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });

    s_lockReleaseTokenPoolAndProxy.applyChainUpdates(chainUpdate);
    s_lockReleaseTokenPoolAndProxyWithAllowList.applyChainUpdates(chainUpdate);
    s_lockReleaseTokenPoolAndProxy.setRebalancer(OWNER);

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: s_allowedOnRamp});
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: s_allowedOffRamp});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
  }
}

contract LockReleaseTokenPoolAndProxy_setRebalancer is LockReleaseTokenPoolAndProxySetup {
  function test_SetRebalancer_Success() public {
    assertEq(address(s_lockReleaseTokenPoolAndProxy.getRebalancer()), OWNER);
    s_lockReleaseTokenPoolAndProxy.setRebalancer(STRANGER);
    assertEq(address(s_lockReleaseTokenPoolAndProxy.getRebalancer()), STRANGER);
  }

  function test_SetRebalancer_Revert() public {
    vm.startPrank(STRANGER);

    vm.expectRevert("Only callable by owner");
    s_lockReleaseTokenPoolAndProxy.setRebalancer(STRANGER);
  }
}

contract LockReleaseTokenPoolPoolAndProxy_canAcceptLiquidity is LockReleaseTokenPoolAndProxySetup {
  function test_CanAcceptLiquidity_Success() public {
    assertEq(true, s_lockReleaseTokenPoolAndProxy.canAcceptLiquidity());

    s_lockReleaseTokenPoolAndProxy =
      new LockReleaseTokenPoolAndProxy(s_token, new address[](0), address(s_mockRMN), false, address(s_sourceRouter));
    assertEq(false, s_lockReleaseTokenPoolAndProxy.canAcceptLiquidity());
  }
}

contract LockReleaseTokenPoolPoolAndProxy_provideLiquidity is LockReleaseTokenPoolAndProxySetup {
  function test_Fuzz_ProvideLiquidity_Success(uint256 amount) public {
    uint256 balancePre = s_token.balanceOf(OWNER);
    s_token.approve(address(s_lockReleaseTokenPoolAndProxy), amount);

    s_lockReleaseTokenPoolAndProxy.provideLiquidity(amount);

    assertEq(s_token.balanceOf(OWNER), balancePre - amount);
    assertEq(s_token.balanceOf(address(s_lockReleaseTokenPoolAndProxy)), amount);
  }

  // Reverts

  function test_Unauthorized_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, STRANGER));

    s_lockReleaseTokenPoolAndProxy.provideLiquidity(1);
  }

  function test_Fuzz_ExceedsAllowance(uint256 amount) public {
    vm.assume(amount > 0);
    vm.expectRevert("ERC20: insufficient allowance");
    s_lockReleaseTokenPoolAndProxy.provideLiquidity(amount);
  }

  function test_LiquidityNotAccepted_Revert() public {
    s_lockReleaseTokenPoolAndProxy =
      new LockReleaseTokenPoolAndProxy(s_token, new address[](0), address(s_mockRMN), false, address(s_sourceRouter));

    vm.expectRevert(LockReleaseTokenPoolAndProxy.LiquidityNotAccepted.selector);
    s_lockReleaseTokenPoolAndProxy.provideLiquidity(1);
  }
}

contract LockReleaseTokenPoolPoolAndProxy_withdrawalLiquidity is LockReleaseTokenPoolAndProxySetup {
  function test_Fuzz_WithdrawalLiquidity_Success(uint256 amount) public {
    uint256 balancePre = s_token.balanceOf(OWNER);
    s_token.approve(address(s_lockReleaseTokenPoolAndProxy), amount);
    s_lockReleaseTokenPoolAndProxy.provideLiquidity(amount);

    s_lockReleaseTokenPoolAndProxy.withdrawLiquidity(amount);

    assertEq(s_token.balanceOf(OWNER), balancePre);
  }

  // Reverts

  function test_Unauthorized_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, STRANGER));

    s_lockReleaseTokenPoolAndProxy.withdrawLiquidity(1);
  }

  function test_InsufficientLiquidity_Revert() public {
    uint256 maxUint256 = 2 ** 256 - 1;
    s_token.approve(address(s_lockReleaseTokenPoolAndProxy), maxUint256);
    s_lockReleaseTokenPoolAndProxy.provideLiquidity(maxUint256);

    vm.startPrank(address(s_lockReleaseTokenPoolAndProxy));
    s_token.transfer(OWNER, maxUint256);
    vm.startPrank(OWNER);

    vm.expectRevert(LockReleaseTokenPoolAndProxy.InsufficientLiquidity.selector);
    s_lockReleaseTokenPoolAndProxy.withdrawLiquidity(1);
  }
}

contract LockReleaseTokenPoolPoolAndProxy_supportsInterface is LockReleaseTokenPoolAndProxySetup {
  function test_SupportsInterface_Success() public view {
    assertTrue(s_lockReleaseTokenPoolAndProxy.supportsInterface(type(IPoolV1).interfaceId));
    assertTrue(s_lockReleaseTokenPoolAndProxy.supportsInterface(type(IERC165).interfaceId));
  }
}
