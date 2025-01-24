// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ILiquidityManager} from "../interfaces/ILiquidityManager.sol";
import {ILiquidityContainer} from "../interfaces/ILiquidityContainer.sol";

import {LiquidityManagerSetup} from "./LiquidityManagerSetup.t.sol";
import {LiquidityManager} from "../LiquidityManager.sol";
import {LiquidityManagerHelper} from "./helpers/LiquidityManagerHelper.sol";
import {MockL1BridgeAdapter} from "./mocks/MockBridgeAdapter.sol";
import {MockLockReleaseTokenPool} from "./mocks/MockTokenPool.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

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
    emit IERC20.Transfer(address(s_lockReleaseTokenPool), address(s_liquidityManager), AMOUNT);

    vm.expectEmit();
    emit IERC20.Approval(address(s_liquidityManager), address(s_bridgeAdapter), AMOUNT);

    vm.expectEmit();
    emit IERC20.Transfer(address(s_liquidityManager), address(s_bridgeAdapter), AMOUNT);

    vm.expectEmit();
    bytes memory encodedNonce = abi.encode(uint256(1));
    emit LiquidityManager.LiquidityTransferred(
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
  function test_rebalanceLiquidity_BetweenPools() external {
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

  function test_rebalanceLiquidity_BetweenPools_AlreadyFinalized() external {
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
    emit ILiquidityContainer.LiquidityRemoved(address(remoteRebalancer), AMOUNT);
    vm.expectEmit();
    emit LiquidityManager.LiquidityTransferred(
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
    emit LiquidityManager.FinalizationFailed(
      maxSeqNum,
      i_remoteChainSelector,
      abi.encode(payload),
      abi.encodeWithSelector(NonceAlreadyUsed.selector, nonce)
    );

    vm.expectEmit();
    emit ILiquidityContainer.LiquidityAdded(address(s_liquidityManager), AMOUNT);

    vm.expectEmit();
    emit LiquidityManager.LiquidityTransferred(
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
    emit ILiquidityContainer.LiquidityRemoved(address(remoteRebalancer), AMOUNT);

    vm.expectEmit();
    emit LiquidityManager.LiquidityTransferred(
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
    emit LiquidityManager.FinalizationStepCompleted(maxSeqNum, i_remoteChainSelector, abi.encode(payload));

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
    emit ILiquidityContainer.LiquidityAdded(address(s_liquidityManager), AMOUNT);

    vm.expectEmit();
    emit LiquidityManager.LiquidityTransferred(
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
    emit ILiquidityContainer.LiquidityRemoved(address(remoteRebalancer), AMOUNT);

    vm.expectEmit();
    emit LiquidityManager.LiquidityTransferred(
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
    emit LiquidityManager.FinalizationStepCompleted(maxSeqNum, i_remoteChainSelector, abi.encode(payload));

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
    emit ILiquidityContainer.LiquidityAdded(address(s_wethRebalancer), AMOUNT);

    vm.expectEmit();
    emit LiquidityManager.LiquidityTransferred(
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
      s_l1Weth.balanceOf(address(s_wethMockLockReleaseTokenPool)),
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
