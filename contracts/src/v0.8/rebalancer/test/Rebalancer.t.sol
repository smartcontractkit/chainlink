// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IBridgeAdapter, IL1BridgeAdapter} from "../interfaces/IBridge.sol";
import {IRebalancer} from "../interfaces/IRebalancer.sol";

import {LockReleaseTokenPool} from "../../ccip/pools/LockReleaseTokenPool.sol";
import {Rebalancer} from "../Rebalancer.sol";
import {MockL1BridgeAdapter} from "./mocks/MockBridgeAdapter.sol";
import {RebalancerBaseTest} from "./RebalancerBaseTest.sol";

import {ERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract RebalancerSetup is RebalancerBaseTest {
  event LiquidityTransferred(
    uint64 indexed ocrSeqNum,
    uint64 indexed fromChainSelector,
    uint64 indexed toChainSelector,
    address to,
    uint256 amount
  );

  Rebalancer internal s_rebalancer;
  LockReleaseTokenPool internal s_lockReleaseTokenPool;
  MockL1BridgeAdapter internal s_bridgeAdapter;

  function setUp() public override {
    RebalancerBaseTest.setUp();

    s_bridgeAdapter = new MockL1BridgeAdapter(s_l1Token);
    s_lockReleaseTokenPool = new LockReleaseTokenPool(s_l1Token, new address[](0), address(1), true);
    s_rebalancer = new Rebalancer(s_l1Token, i_localChainSelector, s_lockReleaseTokenPool);

    s_lockReleaseTokenPool.setRebalancer(address(s_rebalancer));
  }
}

contract Rebalancer_rebalanceLiquidity is RebalancerSetup {
  function test_rebalanceLiquiditySuccess() external {
    uint256 amount = 12345679;
    deal(address(s_l1Token), address(s_lockReleaseTokenPool), amount);

    Rebalancer.CrossChainRebalancerArgs[] memory args = new Rebalancer.CrossChainRebalancerArgs[](1);
    args[0] = IRebalancer.CrossChainRebalancerArgs({
      remoteRebalancer: address(s_rebalancer),
      localBridge: s_bridgeAdapter,
      remoteToken: address(s_l2Token),
      remoteChainSelector: i_remoteChainSelector,
      enabled: true
    });
    s_rebalancer.setCrossChainRebalancer(args);

    vm.expectEmit();
    emit Transfer(address(s_lockReleaseTokenPool), address(s_rebalancer), amount);

    vm.expectEmit();
    emit Approval(address(s_rebalancer), address(s_bridgeAdapter), amount);

    vm.expectEmit();
    emit Transfer(address(s_rebalancer), address(s_bridgeAdapter), amount);

    vm.expectEmit();
    emit LiquidityTransferred(
      type(uint64).max,
      i_localChainSelector,
      i_remoteChainSelector,
      address(s_rebalancer),
      amount
    );

    s_rebalancer.rebalanceLiquidity(i_remoteChainSelector, amount);

    assertEq(s_l1Token.balanceOf(address(s_rebalancer)), 0);
    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), amount);
    assertEq(s_l1Token.allowance(address(s_rebalancer), address(s_bridgeAdapter)), 0);
  }

  /// @notice this test sets up a circular system where the liquidity container of
  /// the local Liquidity manager is the bridge adapter of the remote liquidity manager
  /// and the other way around for the remote liquidity manager. This allows us to
  /// rebalance funds between the two liquidity managers on the same chain.
  function test_rebalanceBetweenPoolsSuccess() external {
    uint256 amount = 12345670;

    s_rebalancer = new Rebalancer(s_l1Token, i_localChainSelector, s_bridgeAdapter);

    MockL1BridgeAdapter mockRemoteBridgeAdapter = new MockL1BridgeAdapter(s_l1Token);
    Rebalancer mockRemoteRebalancer = new Rebalancer(s_l1Token, i_remoteChainSelector, mockRemoteBridgeAdapter);

    Rebalancer.CrossChainRebalancerArgs[] memory args = new Rebalancer.CrossChainRebalancerArgs[](1);
    args[0] = IRebalancer.CrossChainRebalancerArgs({
      remoteRebalancer: address(mockRemoteRebalancer),
      localBridge: mockRemoteBridgeAdapter,
      remoteToken: address(s_l1Token),
      remoteChainSelector: i_remoteChainSelector,
      enabled: true
    });

    s_rebalancer.setCrossChainRebalancer(args);

    args[0] = IRebalancer.CrossChainRebalancerArgs({
      remoteRebalancer: address(s_rebalancer),
      localBridge: s_bridgeAdapter,
      remoteToken: address(s_l1Token),
      remoteChainSelector: i_localChainSelector,
      enabled: true
    });

    mockRemoteRebalancer.setCrossChainRebalancer(args);

    deal(address(s_l1Token), address(s_bridgeAdapter), amount);

    s_rebalancer.rebalanceLiquidity(i_remoteChainSelector, amount);

    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), 0);
    assertEq(s_l1Token.balanceOf(address(mockRemoteBridgeAdapter)), amount);
    assertEq(s_l1Token.allowance(address(s_rebalancer), address(s_bridgeAdapter)), 0);

    mockRemoteRebalancer.rebalanceLiquidity(i_localChainSelector, amount);

    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), amount);
    assertEq(s_l1Token.balanceOf(address(mockRemoteBridgeAdapter)), 0);

    // Assert partial rebalancing works correctly
    s_rebalancer.rebalanceLiquidity(i_remoteChainSelector, amount / 2);

    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), amount / 2);
    assertEq(s_l1Token.balanceOf(address(mockRemoteBridgeAdapter)), amount / 2);
  }

  // Reverts

  function test_InsufficientLiquidityReverts() external {
    uint256 amount = 1245;

    vm.expectRevert(abi.encodeWithSelector(Rebalancer.InsufficientLiquidity.selector, amount, 0));

    s_rebalancer.rebalanceLiquidity(0, amount);
  }

  function test_InvalidRemoteChainReverts() external {
    uint256 amount = 12345679;
    deal(address(s_l1Token), address(s_lockReleaseTokenPool), amount);

    vm.expectRevert(abi.encodeWithSelector(Rebalancer.InvalidRemoteChain.selector, i_remoteChainSelector));

    s_rebalancer.rebalanceLiquidity(i_remoteChainSelector, amount);
  }
}
