// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ILiquidityManager} from "../interfaces/ILiquidityManager.sol";
import {IBridgeAdapter} from "../interfaces/IBridge.sol";

import {LiquidityManager} from "../LiquidityManager.sol";
import {MockL1BridgeAdapter} from "./mocks/MockBridgeAdapter.sol";
import {LiquidityManagerBaseTest} from "./LiquidityManagerBaseTest.t.sol";
import {LiquidityManagerHelper} from "./helpers/LiquidityManagerHelper.sol";
import {MockLockReleaseTokenPool} from "./mocks/MockTokenPool.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

// FOUNDRY_PROFILE=liquiditymanager forge test --match-path src/v0.8/liquiditymanager/test/LiquidityManager.t.sol

contract LiquidityManagerSetup is LiquidityManagerBaseTest {
  event FinalizationStepCompleted(
    uint64 indexed ocrSeqNum,
    uint64 indexed remoteChainSelector,
    bytes bridgeSpecificData
  );
  event LiquidityTransferred(
    uint64 indexed ocrSeqNum,
    uint64 indexed fromChainSelector,
    uint64 indexed toChainSelector,
    address to,
    uint256 amount,
    bytes bridgeSpecificPayload,
    bytes bridgeReturnData
  );
  event FinalizationFailed(
    uint64 indexed ocrSeqNum,
    uint64 indexed remoteChainSelector,
    bytes bridgeSpecificData,
    bytes reason
  );
  event FinanceRoleSet(address financeRole);
  event LiquidityAddedToContainer(address indexed provider, uint256 indexed amount);
  event LiquidityRemovedFromContainer(address indexed remover, uint256 indexed amount);
  // Liquidity container event
  event LiquidityAdded(address indexed provider, uint256 indexed amount);
  event LiquidityRemoved(address indexed remover, uint256 indexed amount);

  error NonceAlreadyUsed(uint256 nonce);

  LiquidityManagerHelper internal s_liquidityManager;
  MockLockReleaseTokenPool internal s_lockReleaseTokenPool;
  MockL1BridgeAdapter internal s_bridgeAdapter;

  // LiquidityManager that rebalances weth.
  LiquidityManagerHelper internal s_wethRebalancer;
  MockLockReleaseTokenPool internal s_wethMockLockReleaseTokenPool;
  MockL1BridgeAdapter internal s_wethBridgeAdapter;

  function setUp() public virtual override {
    LiquidityManagerBaseTest.setUp();

    s_bridgeAdapter = new MockL1BridgeAdapter(s_l1Token, false);
    s_lockReleaseTokenPool = new MockLockReleaseTokenPool(s_l1Token);
    s_liquidityManager = new LiquidityManagerHelper(
      s_l1Token,
      i_localChainSelector,
      s_lockReleaseTokenPool,
      0,
      FINANCE
    );

    s_wethBridgeAdapter = new MockL1BridgeAdapter(IERC20(address(s_l1Weth)), true);
    s_wethMockLockReleaseTokenPool = new MockLockReleaseTokenPool(IERC20(address(s_l1Weth)));
    s_wethRebalancer = new LiquidityManagerHelper(
      IERC20(address(s_l1Weth)),
      i_localChainSelector,
      s_wethMockLockReleaseTokenPool,
      0,
      FINANCE
    );
  }
}

contract LiquidityManager_addLiquidity is LiquidityManagerSetup {
  function test_addLiquidity() external {
    address caller = STRANGER;
    changePrank(caller);

    uint256 amount = 12345679;
    deal(address(s_l1Token), caller, amount);

    s_l1Token.approve(address(s_liquidityManager), amount);

    vm.expectEmit();
    emit LiquidityAddedToContainer(caller, amount);

    s_liquidityManager.addLiquidity(amount);

    assertEq(s_l1Token.balanceOf(address(s_lockReleaseTokenPool)), amount);
  }
}

contract LiquidityManager_removeLiquidity is LiquidityManagerSetup {
  function test_removeLiquiditySuccess() external {
    uint256 amount = 12345679;
    deal(address(s_l1Token), address(s_lockReleaseTokenPool), amount);

    vm.expectEmit();
    emit LiquidityRemovedFromContainer(FINANCE, amount);

    changePrank(FINANCE);
    s_liquidityManager.removeLiquidity(amount);

    assertEq(s_l1Token.balanceOf(address(s_liquidityManager)), 0);
  }

  function test_InsufficientLiquidityReverts() external {
    uint256 balance = 923;
    uint256 requested = balance + 1;

    deal(address(s_l1Token), address(s_lockReleaseTokenPool), balance);

    vm.expectRevert(abi.encodeWithSelector(LiquidityManager.InsufficientLiquidity.selector, requested, balance, 0));

    changePrank(FINANCE);
    s_liquidityManager.removeLiquidity(requested);
  }

  function test_removeLiquidity_RevertWhen_NotFinanceRole() external {
    vm.stopPrank();

    vm.expectRevert(LiquidityManager.OnlyFinanceRole.selector);

    s_liquidityManager.removeLiquidity(123);
  }
}

contract LiquidityManager__report is LiquidityManagerSetup {
  function test_report_RevertWhen_EmptyReport() external {
    ILiquidityManager.LiquidityInstructions memory instructions = ILiquidityManager.LiquidityInstructions({
      sendLiquidityParams: new ILiquidityManager.SendLiquidityParams[](0),
      receiveLiquidityParams: new ILiquidityManager.ReceiveLiquidityParams[](0)
    });

    vm.expectRevert(LiquidityManager.EmptyReport.selector);

    s_liquidityManager.report(abi.encode(instructions), 123);
  }
}

contract LiquidityManager_rebalanceLiquidity is LiquidityManagerSetup {
  uint256 internal constant AMOUNT = 12345679;

  function test_rebalanceLiquidity() external {
    deal(address(s_l1Token), address(s_lockReleaseTokenPool), AMOUNT);

    LiquidityManager.CrossChainRebalancerArgs[] memory args = new LiquidityManager.CrossChainRebalancerArgs[](1);
    args[0] = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: address(s_liquidityManager),
      localBridge: s_bridgeAdapter,
      remoteToken: address(s_l2Token),
      remoteChainSelector: i_remoteChainSelector,
      enabled: true
    });
    s_liquidityManager.setCrossChainRebalancers(args);

    vm.expectEmit();
    emit Transfer(address(s_lockReleaseTokenPool), address(s_liquidityManager), AMOUNT);

    vm.expectEmit();
    emit Approval(address(s_liquidityManager), address(s_bridgeAdapter), AMOUNT);

    vm.expectEmit();
    emit Transfer(address(s_liquidityManager), address(s_bridgeAdapter), AMOUNT);

    vm.expectEmit();
    bytes memory encodedNonce = abi.encode(uint256(1));
    emit LiquidityTransferred(
      type(uint64).max,
      i_localChainSelector,
      i_remoteChainSelector,
      address(s_liquidityManager),
      AMOUNT,
      bytes(""),
      encodedNonce
    );

    vm.startPrank(FINANCE);
    s_liquidityManager.rebalanceLiquidity(i_remoteChainSelector, AMOUNT, 0, bytes(""));

    assertEq(s_l1Token.balanceOf(address(s_liquidityManager)), 0);
    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), AMOUNT);
    assertEq(s_l1Token.allowance(address(s_liquidityManager), address(s_bridgeAdapter)), 0);
  }

  /// @notice this test sets up a circular system where the liquidity container of
  /// the local Liquidity manager is the bridge adapter of the remote liquidity manager
  /// and the other way around for the remote liquidity manager. This allows us to
  /// rebalance funds between the two liquidity managers on the same chain.
  function test_rebalanceLiquidit() external {
    uint256 amount = 12345670;

    s_liquidityManager = new LiquidityManagerHelper(s_l1Token, i_localChainSelector, s_bridgeAdapter, 0, FINANCE);

    MockL1BridgeAdapter mockRemoteBridgeAdapter = new MockL1BridgeAdapter(s_l1Token, false);
    LiquidityManager mockRemoteRebalancer = new LiquidityManager(
      s_l1Token,
      i_remoteChainSelector,
      mockRemoteBridgeAdapter,
      0,
      FINANCE
    );

    LiquidityManager.CrossChainRebalancerArgs[] memory args = new LiquidityManager.CrossChainRebalancerArgs[](1);
    args[0] = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: address(mockRemoteRebalancer),
      localBridge: mockRemoteBridgeAdapter,
      remoteToken: address(s_l1Token),
      remoteChainSelector: i_remoteChainSelector,
      enabled: true
    });

    s_liquidityManager.setCrossChainRebalancers(args);

    args[0] = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: address(s_liquidityManager),
      localBridge: s_bridgeAdapter,
      remoteToken: address(s_l1Token),
      remoteChainSelector: i_localChainSelector,
      enabled: true
    });

    mockRemoteRebalancer.setCrossChainRebalancers(args);

    deal(address(s_l1Token), address(s_bridgeAdapter), amount);

    vm.startPrank(FINANCE);
    s_liquidityManager.rebalanceLiquidity(i_remoteChainSelector, amount, 0, bytes(""));

    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), 0);
    assertEq(s_l1Token.balanceOf(address(mockRemoteBridgeAdapter)), amount);
    assertEq(s_l1Token.allowance(address(s_liquidityManager), address(s_bridgeAdapter)), 0);

    // attach a bridge fee and see the relevant adapter's ether balance change.
    // the bridge fee is sent along with the sendERC20 call.
    uint256 bridgeFee = 123;
    vm.deal(address(mockRemoteRebalancer), bridgeFee);
    mockRemoteRebalancer.rebalanceLiquidity(i_localChainSelector, amount, bridgeFee, bytes(""));

    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), amount);
    assertEq(s_l1Token.balanceOf(address(mockRemoteBridgeAdapter)), 0);
    assertEq(address(s_bridgeAdapter).balance, bridgeFee);

    // Assert partial rebalancing works correctly
    s_liquidityManager.rebalanceLiquidity(i_remoteChainSelector, amount / 2, 0, bytes(""));

    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), amount / 2);
    assertEq(s_l1Token.balanceOf(address(mockRemoteBridgeAdapter)), amount / 2);
  }

  function test_rebalanceBetweenPools_AlreadyFinalized() external {
    // set up a rebalancer on another chain, an "L2".
    // note we use the L1 bridge adapter because it has the reverting logic
    // when finalization is already done.
    MockL1BridgeAdapter remoteBridgeAdapter = new MockL1BridgeAdapter(s_l2Token, false);
    MockLockReleaseTokenPool remotePool = new MockLockReleaseTokenPool(s_l2Token);
    LiquidityManager remoteRebalancer = new LiquidityManager(s_l2Token, i_remoteChainSelector, remotePool, 0, FINANCE);

    // set up the cross chain rebalancer on "L1".
    LiquidityManager.CrossChainRebalancerArgs[] memory args = new LiquidityManager.CrossChainRebalancerArgs[](1);
    args[0] = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: address(remoteRebalancer),
      localBridge: s_bridgeAdapter,
      remoteToken: address(s_l2Token),
      remoteChainSelector: i_remoteChainSelector,
      enabled: true
    });

    s_liquidityManager.setCrossChainRebalancers(args);

    // set up the cross chain rebalancer on "L2".
    args[0] = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: address(s_liquidityManager),
      localBridge: remoteBridgeAdapter,
      remoteToken: address(s_l1Token),
      remoteChainSelector: i_localChainSelector,
      enabled: true
    });

    remoteRebalancer.setCrossChainRebalancers(args);

    // deal some L1 tokens to the L1 bridge adapter so that it can send them to the rebalancer
    // when the withdrawal gets finalized.
    deal(address(s_l1Token), address(s_bridgeAdapter), AMOUNT);
    // deal some L2 tokens to the remote token pool so that we can withdraw it when we rebalance.
    deal(address(s_l2Token), address(remotePool), AMOUNT);

    uint256 nonce = 1;
    uint64 maxSeqNum = type(uint64).max;
    bytes memory bridgeSendReturnData = abi.encode(nonce);
    bytes memory bridgeSpecificPayload = bytes("");
    vm.expectEmit();
    emit LiquidityRemoved(address(remoteRebalancer), AMOUNT);
    vm.expectEmit();
    emit LiquidityTransferred(
      maxSeqNum,
      i_remoteChainSelector,
      i_localChainSelector,
      address(s_liquidityManager),
      AMOUNT,
      bridgeSpecificPayload,
      bridgeSendReturnData
    );
    vm.startPrank(FINANCE);
    remoteRebalancer.rebalanceLiquidity(i_localChainSelector, AMOUNT, 0, bridgeSpecificPayload);

    // available liquidity has been moved to the remote bridge adapter from the token pool.
    assertEq(s_l2Token.balanceOf(address(remoteBridgeAdapter)), AMOUNT, "remoteBridgeAdapter balance");
    assertEq(s_l2Token.balanceOf(address(remotePool)), 0, "remotePool balance");

    // prove and finalize manually on the L1 bridge adapter.
    // this should transfer the funds to the rebalancer.
    MockL1BridgeAdapter.ProvePayload memory provePayload = MockL1BridgeAdapter.ProvePayload({nonce: nonce});
    MockL1BridgeAdapter.Payload memory payload = MockL1BridgeAdapter.Payload({
      action: MockL1BridgeAdapter.FinalizationAction.ProveWithdrawal,
      data: abi.encode(provePayload)
    });
    bool fundsAvailable = s_bridgeAdapter.finalizeWithdrawERC20(
      address(0),
      address(s_liquidityManager),
      abi.encode(payload)
    );
    assertFalse(fundsAvailable, "fundsAvailable must be false");
    MockL1BridgeAdapter.FinalizePayload memory finalizePayload = MockL1BridgeAdapter.FinalizePayload({
      nonce: nonce,
      amount: AMOUNT
    });
    payload = MockL1BridgeAdapter.Payload({
      action: MockL1BridgeAdapter.FinalizationAction.FinalizeWithdrawal,
      data: abi.encode(finalizePayload)
    });
    fundsAvailable = s_bridgeAdapter.finalizeWithdrawERC20(
      address(0),
      address(s_liquidityManager),
      abi.encode(payload)
    );
    assertTrue(fundsAvailable, "fundsAvailable must be true");

    // available balance on the L1 bridge adapter has been moved to the rebalancer.
    assertEq(s_l1Token.balanceOf(address(s_liquidityManager)), AMOUNT, "rebalancer balance 1");
    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), 0, "bridgeAdapter balance");

    // try to finalize on L1 again
    // bytes memory revertData = abi.encodeWithSelector(NonceAlreadyUsed.selector, nonce);
    vm.expectEmit();
    emit FinalizationFailed(
      maxSeqNum,
      i_remoteChainSelector,
      abi.encode(payload),
      abi.encodeWithSelector(NonceAlreadyUsed.selector, nonce)
    );
    vm.expectEmit();
    emit LiquidityAdded(address(s_liquidityManager), AMOUNT);
    vm.expectEmit();
    emit LiquidityTransferred(
      maxSeqNum,
      i_remoteChainSelector,
      i_localChainSelector,
      address(s_liquidityManager),
      AMOUNT,
      abi.encode(payload),
      bytes("")
    );
    s_liquidityManager.receiveLiquidity(i_remoteChainSelector, AMOUNT, false, abi.encode(payload));

    // available balance on the rebalancer has been injected into the token pool.
    assertEq(s_l1Token.balanceOf(address(s_liquidityManager)), 0, "rebalancer balance 2");
    assertEq(s_l1Token.balanceOf(address(s_lockReleaseTokenPool)), AMOUNT, "lockReleaseTokenPool balance");
  }

  function test_rebalanceLiquidity_RebalanceBetweenPoolsMultiStageFinalization() external {
    // set up a rebalancer on another chain, an "L2".
    // note we use the L1 bridge adapter because it has the reverting logic
    // when finalization is already done.
    MockL1BridgeAdapter remoteBridgeAdapter = new MockL1BridgeAdapter(s_l2Token, false);
    MockLockReleaseTokenPool remotePool = new MockLockReleaseTokenPool(s_l2Token);
    LiquidityManager remoteRebalancer = new LiquidityManager(s_l2Token, i_remoteChainSelector, remotePool, 0, FINANCE);

    // set up the cross chain rebalancer on "L1".
    LiquidityManager.CrossChainRebalancerArgs[] memory args = new LiquidityManager.CrossChainRebalancerArgs[](1);
    args[0] = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: address(remoteRebalancer),
      localBridge: s_bridgeAdapter,
      remoteToken: address(s_l2Token),
      remoteChainSelector: i_remoteChainSelector,
      enabled: true
    });

    s_liquidityManager.setCrossChainRebalancers(args);

    // set up the cross chain rebalancer on "L2".
    args[0] = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: address(s_liquidityManager),
      localBridge: remoteBridgeAdapter,
      remoteToken: address(s_l1Token),
      remoteChainSelector: i_localChainSelector,
      enabled: true
    });

    remoteRebalancer.setCrossChainRebalancers(args);

    // deal some L1 tokens to the L1 bridge adapter so that it can send them to the rebalancer
    // when the withdrawal gets finalized.
    deal(address(s_l1Token), address(s_bridgeAdapter), AMOUNT);
    // deal some L2 tokens to the remote token pool so that we can withdraw it when we rebalance.
    deal(address(s_l2Token), address(remotePool), AMOUNT);

    // initiate a send from remote rebalancer to s_liquidityManager.
    uint256 nonce = 1;
    uint64 maxSeqNum = type(uint64).max;
    bytes memory bridgeSendReturnData = abi.encode(nonce);
    bytes memory bridgeSpecificPayload = bytes("");
    vm.expectEmit();
    emit LiquidityRemoved(address(remoteRebalancer), AMOUNT);
    vm.expectEmit();
    emit LiquidityTransferred(
      maxSeqNum,
      i_remoteChainSelector,
      i_localChainSelector,
      address(s_liquidityManager),
      AMOUNT,
      bridgeSpecificPayload,
      bridgeSendReturnData
    );
    vm.startPrank(FINANCE);
    remoteRebalancer.rebalanceLiquidity(i_localChainSelector, AMOUNT, 0, bridgeSpecificPayload);

    // available liquidity has been moved to the remote bridge adapter from the token pool.
    assertEq(s_l2Token.balanceOf(address(remoteBridgeAdapter)), AMOUNT, "remoteBridgeAdapter balance");
    assertEq(s_l2Token.balanceOf(address(remotePool)), 0, "remotePool balance");

    // prove withdrawal on the L1 bridge adapter, through the rebalancer.
    uint256 balanceBeforeProve = s_l1Token.balanceOf(address(s_lockReleaseTokenPool));
    MockL1BridgeAdapter.ProvePayload memory provePayload = MockL1BridgeAdapter.ProvePayload({nonce: nonce});
    MockL1BridgeAdapter.Payload memory payload = MockL1BridgeAdapter.Payload({
      action: MockL1BridgeAdapter.FinalizationAction.ProveWithdrawal,
      data: abi.encode(provePayload)
    });
    vm.expectEmit();
    emit FinalizationStepCompleted(maxSeqNum, i_remoteChainSelector, abi.encode(payload));
    s_liquidityManager.receiveLiquidity(i_remoteChainSelector, AMOUNT, false, abi.encode(payload));

    // s_liquidityManager should have no tokens.
    assertEq(s_l1Token.balanceOf(address(s_liquidityManager)), 0, "rebalancer balance 1");
    // balance of s_lockReleaseTokenPool should be unchanged since no liquidity got added yet.
    assertEq(
      s_l1Token.balanceOf(address(s_lockReleaseTokenPool)),
      balanceBeforeProve,
      "s_lockReleaseTokenPool balance should be unchanged"
    );

    // finalize withdrawal on the L1 bridge adapter, through the rebalancer.
    MockL1BridgeAdapter.FinalizePayload memory finalizePayload = MockL1BridgeAdapter.FinalizePayload({
      nonce: nonce,
      amount: AMOUNT
    });
    payload = MockL1BridgeAdapter.Payload({
      action: MockL1BridgeAdapter.FinalizationAction.FinalizeWithdrawal,
      data: abi.encode(finalizePayload)
    });
    vm.expectEmit();
    emit LiquidityAdded(address(s_liquidityManager), AMOUNT);
    vm.expectEmit();
    emit LiquidityTransferred(
      maxSeqNum,
      i_remoteChainSelector,
      i_localChainSelector,
      address(s_liquidityManager),
      AMOUNT,
      abi.encode(payload),
      bytes("")
    );
    s_liquidityManager.receiveLiquidity(i_remoteChainSelector, AMOUNT, false, abi.encode(payload));

    // s_liquidityManager should have no tokens.
    assertEq(s_l1Token.balanceOf(address(s_liquidityManager)), 0, "rebalancer balance 2");
    // balance of s_lockReleaseTokenPool should be updated
    assertEq(
      s_l1Token.balanceOf(address(s_lockReleaseTokenPool)),
      balanceBeforeProve + AMOUNT,
      "s_lockReleaseTokenPool balance should be updated"
    );
  }

  function test_rebalanceLiquidity_NativeRewrap() external {
    // set up a rebalancer similar to the above on another chain, an "L2".
    MockL1BridgeAdapter remoteBridgeAdapter = new MockL1BridgeAdapter(IERC20(address(s_l2Weth)), true);
    MockLockReleaseTokenPool remotePool = new MockLockReleaseTokenPool(IERC20(address(s_l2Weth)));
    LiquidityManager remoteRebalancer = new LiquidityManager(
      IERC20(address(s_l2Weth)),
      i_remoteChainSelector,
      remotePool,
      0,
      FINANCE
    );

    // set up the cross chain rebalancer on "L1".
    LiquidityManager.CrossChainRebalancerArgs[] memory args = new LiquidityManager.CrossChainRebalancerArgs[](1);
    args[0] = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: address(remoteRebalancer),
      localBridge: s_wethBridgeAdapter,
      remoteToken: address(s_l2Weth),
      remoteChainSelector: i_remoteChainSelector,
      enabled: true
    });

    s_wethRebalancer.setCrossChainRebalancers(args);

    // set up the cross chain rebalancer on "L2".
    args[0] = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: address(s_wethRebalancer),
      localBridge: remoteBridgeAdapter,
      remoteToken: address(s_l1Weth),
      remoteChainSelector: i_localChainSelector,
      enabled: true
    });

    remoteRebalancer.setCrossChainRebalancers(args);

    // deal some ether to the L1 bridge adapter so that it can send them to the rebalancer
    // when the withdrawal gets finalized.
    vm.deal(address(s_wethBridgeAdapter), AMOUNT);
    // deal some L2 tokens to the remote token pool so that we can withdraw it when we rebalance.
    deal(address(s_l2Weth), address(remotePool), AMOUNT);
    // deposit some eth to the weth contract on L2 from the remote bridge adapter
    // so that the withdraw() call succeeds.
    vm.deal(address(remoteBridgeAdapter), AMOUNT);
    vm.startPrank(address(remoteBridgeAdapter));
    s_l2Weth.deposit{value: AMOUNT}();
    vm.stopPrank();

    // switch to finance for the rest of the test to avoid reverts.
    vm.startPrank(FINANCE);

    // initiate a send from remote rebalancer to s_wethRebalancer.
    uint256 nonce = 1;
    uint64 maxSeqNum = type(uint64).max;
    bytes memory bridgeSendReturnData = abi.encode(nonce);
    bytes memory bridgeSpecificPayload = bytes("");
    vm.expectEmit();
    emit LiquidityRemoved(address(remoteRebalancer), AMOUNT);
    vm.expectEmit();
    emit LiquidityTransferred(
      maxSeqNum,
      i_remoteChainSelector,
      i_localChainSelector,
      address(s_wethRebalancer),
      AMOUNT,
      bridgeSpecificPayload,
      bridgeSendReturnData
    );
    remoteRebalancer.rebalanceLiquidity(i_localChainSelector, AMOUNT, 0, bridgeSpecificPayload);

    // available liquidity has been moved to the remote bridge adapter from the token pool.
    assertEq(s_l2Weth.balanceOf(address(remoteBridgeAdapter)), AMOUNT, "remoteBridgeAdapter balance");
    assertEq(s_l2Weth.balanceOf(address(remotePool)), 0, "remotePool balance");

    // prove withdrawal on the L1 bridge adapter, through the rebalancer.
    uint256 balanceBeforeProve = s_l1Weth.balanceOf(address(s_wethMockLockReleaseTokenPool));
    MockL1BridgeAdapter.ProvePayload memory provePayload = MockL1BridgeAdapter.ProvePayload({nonce: nonce});
    MockL1BridgeAdapter.Payload memory payload = MockL1BridgeAdapter.Payload({
      action: MockL1BridgeAdapter.FinalizationAction.ProveWithdrawal,
      data: abi.encode(provePayload)
    });
    vm.expectEmit();
    emit FinalizationStepCompleted(maxSeqNum, i_remoteChainSelector, abi.encode(payload));
    s_wethRebalancer.receiveLiquidity(i_remoteChainSelector, AMOUNT, false, abi.encode(payload));

    // s_wethRebalancer should have no tokens.
    assertEq(s_l1Weth.balanceOf(address(s_wethRebalancer)), 0, "rebalancer balance 1");
    // balance of s_wethMockTokenPool should be unchanged since no liquidity got added yet.
    assertEq(
      s_l1Weth.balanceOf(address(s_wethMockLockReleaseTokenPool)),
      balanceBeforeProve,
      "s_wethMockLockReleaseTokenPool balance should be unchanged"
    );

    // finalize withdrawal on the L1 bridge adapter, through the rebalancer.
    MockL1BridgeAdapter.FinalizePayload memory finalizePayload = MockL1BridgeAdapter.FinalizePayload({
      nonce: nonce,
      amount: AMOUNT
    });
    payload = MockL1BridgeAdapter.Payload({
      action: MockL1BridgeAdapter.FinalizationAction.FinalizeWithdrawal,
      data: abi.encode(finalizePayload)
    });
    vm.expectEmit();
    emit LiquidityAdded(address(s_wethRebalancer), AMOUNT);
    vm.expectEmit();
    emit LiquidityTransferred(
      maxSeqNum,
      i_remoteChainSelector,
      i_localChainSelector,
      address(s_wethRebalancer),
      AMOUNT,
      abi.encode(payload),
      bytes("")
    );
    s_wethRebalancer.receiveLiquidity(i_remoteChainSelector, AMOUNT, true, abi.encode(payload));

    // s_wethRebalancer should have no tokens.
    assertEq(s_l1Weth.balanceOf(address(s_wethRebalancer)), 0, "rebalancer balance 2");
    // s_wethRebalancer should have no native tokens.
    assertEq(address(s_wethRebalancer).balance, 0, "rebalancer native balance should be zero");
    // balance of s_wethMockLockReleaseTokenPool should be updated
    assertEq(
      s_l1Weth.balanceOf(address(s_wethMockLockReleaseTokenPool                           )),
      balanceBeforeProve + AMOUNT,
      "s_wethMockLockReleaseTokenPool                            balance should be updated"
    );
  }

  // Reverts

  function test_rebalanceLiquidity_RevertWhen_InsufficientLiquidity() external {
    s_liquidityManager.setMinimumLiquidity(3);
    deal(address(s_l1Token), address(s_lockReleaseTokenPool), AMOUNT);
    vm.expectRevert(abi.encodeWithSelector(LiquidityManager.InsufficientLiquidity.selector, AMOUNT, AMOUNT, 3));

    vm.startPrank(FINANCE);
    s_liquidityManager.rebalanceLiquidity(0, AMOUNT, 0, bytes(""));
  }

  function test_rebalanceLiquidity_RevertWhen_InvalidRemoteChain() external {
    deal(address(s_l1Token), address(s_lockReleaseTokenPool), AMOUNT);

    vm.expectRevert(abi.encodeWithSelector(LiquidityManager.InvalidRemoteChain.selector, i_remoteChainSelector));

    vm.startPrank(FINANCE);
    s_liquidityManager.rebalanceLiquidity(i_remoteChainSelector, AMOUNT, 0, bytes(""));
  }
}

contract LiquidityManager_setCrossChainRebalancer is LiquidityManagerSetup {
  event CrossChainRebalancerSet(
    uint64 indexed remoteChainSelector,
    IBridgeAdapter localBridge,
    address remoteToken,
    address remoteRebalancer,
    bool enabled
  );

  function test_setCrossChainRebalancer() external {
    address newRebalancer = address(23892423);
    uint64 remoteChainSelector = 12301293;

    uint64[] memory supportedChains = s_liquidityManager.getSupportedDestChains();
    assertEq(supportedChains.length, 0);

    LiquidityManager.CrossChainRebalancerArgs[] memory args = new LiquidityManager.CrossChainRebalancerArgs[](1);
    args[0] = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: newRebalancer,
      localBridge: s_bridgeAdapter,
      remoteToken: address(190490124908),
      remoteChainSelector: remoteChainSelector,
      enabled: true
    });

    vm.expectEmit();
    emit CrossChainRebalancerSet(
      remoteChainSelector,
      args[0].localBridge,
      args[0].remoteToken,
      newRebalancer,
      args[0].enabled
    );

    s_liquidityManager.setCrossChainRebalancers(args);

    assertEq(s_liquidityManager.getCrossChainRebalancer(remoteChainSelector).remoteRebalancer, newRebalancer);

    LiquidityManager.CrossChainRebalancerArgs[] memory got = s_liquidityManager.getAllCrossChainRebalancers();
    assertEq(got.length, 1);
    assertEq(got[0].remoteRebalancer, args[0].remoteRebalancer);
    assertEq(address(got[0].localBridge), address(args[0].localBridge));
    assertEq(got[0].remoteToken, args[0].remoteToken);
    assertEq(got[0].remoteChainSelector, args[0].remoteChainSelector);
    assertEq(got[0].enabled, args[0].enabled);

    supportedChains = s_liquidityManager.getSupportedDestChains();
    assertEq(supportedChains.length, 1);
    assertEq(supportedChains[0], remoteChainSelector);

    address anotherRebalancer = address(123);
    args[0].remoteRebalancer = anotherRebalancer;

    vm.expectEmit();
    emit CrossChainRebalancerSet(
      remoteChainSelector,
      args[0].localBridge,
      args[0].remoteToken,
      anotherRebalancer,
      args[0].enabled
    );

    s_liquidityManager.setCrossChainRebalancer(args[0]);

    assertEq(s_liquidityManager.getCrossChainRebalancer(remoteChainSelector).remoteRebalancer, anotherRebalancer);

    supportedChains = s_liquidityManager.getSupportedDestChains();
    assertEq(supportedChains.length, 1);
    assertEq(supportedChains[0], remoteChainSelector);
  }

  function test_ZeroChainSelectorReverts() external {
    LiquidityManager.CrossChainRebalancerArgs memory arg = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: address(9),
      localBridge: s_bridgeAdapter,
      remoteToken: address(190490124908),
      remoteChainSelector: 0,
      enabled: true
    });

    vm.expectRevert(LiquidityManager.ZeroChainSelector.selector);

    s_liquidityManager.setCrossChainRebalancer(arg);
  }

  function test_setCrossChainRebalancer_RevertsWhen_ZeroAddressRemoteRebalancer() external {
    LiquidityManager.CrossChainRebalancerArgs memory arg = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: address(0),
      localBridge: s_bridgeAdapter,
      remoteToken: address(190490124908),
      remoteChainSelector: 123,
      enabled: true
    });

    vm.expectRevert(LiquidityManager.ZeroAddress.selector);

    s_liquidityManager.setCrossChainRebalancer(arg);

    arg.remoteRebalancer = address(9);
    arg.localBridge = IBridgeAdapter(address(0));

    vm.expectRevert(LiquidityManager.ZeroAddress.selector);

    s_liquidityManager.setCrossChainRebalancer(arg);

    arg.localBridge = s_bridgeAdapter;
    arg.remoteToken = address(0);

    vm.expectRevert(LiquidityManager.ZeroAddress.selector);

    s_liquidityManager.setCrossChainRebalancer(arg);
  }

  function test_setCrossChainRebalancer_RevertsWhen_CallerNotOwner() external {
    vm.stopPrank();

    vm.expectRevert("Only callable by owner");

    // Test the entrypoint that takes a list
    s_liquidityManager.setCrossChainRebalancers(new LiquidityManager.CrossChainRebalancerArgs[](0));

    vm.expectRevert("Only callable by owner");

    // Test the entrypoint that takes a single item
    s_liquidityManager.setCrossChainRebalancer(
      ILiquidityManager.CrossChainRebalancerArgs({
        remoteRebalancer: address(9),
        localBridge: s_bridgeAdapter,
        remoteToken: address(190490124908),
        remoteChainSelector: 124,
        enabled: true
      })
    );
  }
}

contract LiquidityManager_setLocalLiquidityContainer is LiquidityManagerSetup {
  event LiquidityContainerSet(address indexed newLiquidityContainer);

  function test_setLocalLiquidityContainer() external {
    MockLockReleaseTokenPool newPool = new MockLockReleaseTokenPool(s_l1Token);

    vm.expectEmit();
    emit LiquidityContainerSet(address(newPool));

    s_liquidityManager.setLocalLiquidityContainer(newPool);

    assertEq(s_liquidityManager.getLocalLiquidityContainer(), address(newPool));
  }

  function test_setLocalLiquidityContainer_RevertsWhen_CallerNotOwner() external {
    vm.stopPrank();

    vm.expectRevert("Only callable by owner");

    s_liquidityManager.setLocalLiquidityContainer(MockLockReleaseTokenPool(address(1)));
  }

  function test_setLocalLiquidityContainer_RevertsWhen_CalledWithTheZeroAddress() external {
    vm.expectRevert(LiquidityManager.ZeroAddress.selector);
    s_liquidityManager.setLocalLiquidityContainer(MockLockReleaseTokenPool(address(0)));
  }
}

contract LiquidityManager_setMinimumLiquidity is LiquidityManagerSetup {
  event MinimumLiquiditySet(uint256 oldBalance, uint256 newBalance);

  function test_setMinimumLiquidity() external {
    vm.expectEmit();
    emit MinimumLiquiditySet(uint256(0), uint256(1000));
    s_liquidityManager.setMinimumLiquidity(1000);
    assertEq(s_liquidityManager.getMinimumLiquidity(), uint256(1000));
  }

  function test_setMinimumLiquidity_RevertWhen_CallerNotOwner() external {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    s_liquidityManager.setMinimumLiquidity(uint256(1000));
  }
}

contract LiquidityManager_setFinanceRole is LiquidityManagerSetup {
  event MinimumLiquiditySet(uint256 oldBalance, uint256 newBalance);

  function test_setFinanceRole() external {
    vm.expectEmit();
    address newFinanceRole = makeAddr("newFinanceRole");
    assertEq(s_liquidityManager.getFinanceRole(), FINANCE);
    emit FinanceRoleSet(newFinanceRole);
    s_liquidityManager.setFinanceRole(newFinanceRole);
    assertEq(s_liquidityManager.getFinanceRole(), newFinanceRole);
  }

  function test_setFinanceRole_RevertWhen_CallerNotOwner() external {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    s_liquidityManager.setFinanceRole(address(1));
  }
}

contract LiquidityManager_withdrawNative is LiquidityManagerSetup {
  event NativeWithdrawn(uint256 amount, address destination);

  address private receiver = makeAddr("receiver");

  function setUp() public override {
    super.setUp();
    vm.deal(address(s_liquidityManager), 1);
  }

  function test_withdrawNative() external {
    assertEq(receiver.balance, 0);
    vm.expectEmit();
    emit NativeWithdrawn(1, receiver);
    changePrank(FINANCE);
    s_liquidityManager.withdrawNative(1, payable(receiver));
    assertEq(receiver.balance, 1);
  }

  function test_withdrawNative_FinanceRoleReverts() external {
    vm.stopPrank();
    vm.expectRevert(LiquidityManager.OnlyFinanceRole.selector);
    s_liquidityManager.withdrawNative(1, payable(receiver));
  }
}

contract LiquidityManager_receive is LiquidityManagerSetup {
  event NativeDeposited(uint256 amount, address depositor);

  address private depositor = makeAddr("depositor");

  function test_receive() external {
    vm.deal(depositor, 100);
    uint256 before = address(s_liquidityManager).balance;
    vm.expectEmit();
    emit NativeDeposited(100, depositor);
    changePrank(depositor);
    payable(address(s_liquidityManager)).transfer(100);
    assertEq(address(s_liquidityManager).balance, before + 100);
  }
}

contract LiquidityManager_withdrawERC20 is LiquidityManagerSetup {
  function test_withdrawERC20() external {
    uint256 amount = 100;
    deal(address(s_otherToken), address(s_liquidityManager), amount);
    assertEq(s_otherToken.balanceOf(address(1)), 0);
    assertEq(s_otherToken.balanceOf(address(s_liquidityManager)), amount);
    vm.startPrank(FINANCE);
    s_liquidityManager.withdrawERC20(address(s_otherToken), amount, address(1));
    assertEq(s_otherToken.balanceOf(address(1)), amount);
    assertEq(s_otherToken.balanceOf(address(s_liquidityManager)), 0);
  }

  function test_withdrawERC20_RevertWhen_InvalidCondition() external {
    uint256 amount = 100;
    deal(address(s_otherToken), address(s_liquidityManager), amount);
    vm.startPrank(STRANGER);
    vm.expectRevert(LiquidityManager.OnlyFinanceRole.selector);
    s_liquidityManager.withdrawERC20(address(s_otherToken), amount, address(1));
  }
}
