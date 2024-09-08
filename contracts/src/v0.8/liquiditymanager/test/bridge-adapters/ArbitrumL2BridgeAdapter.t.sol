// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IWrappedNative} from "../../../ccip/interfaces/IWrappedNative.sol";

import {ArbitrumL2BridgeAdapter, IL2GatewayRouter} from "../../bridge-adapters/ArbitrumL2BridgeAdapter.sol";
import "forge-std/Test.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

//contract ArbitrumL2BridgeAdapterSetup is Test {
//  uint256 internal arbitrumFork;
//
//  string internal constant ARBITRUM_RPC_URL = "<url>";
//
//  address internal constant L2_GATEWAY_ROUTER = 0x5288c571Fd7aD117beA99bF60FE0846C4E84F933;
//  address internal constant L2_ETH_WITHDRAWAL_PRECOMPILE = 0x0000000000000000000000000000000000000064;
//
//  IERC20 internal constant L1_LINK = IERC20(0x514910771AF9Ca656af840dff83E8264EcF986CA);
//  IERC20 internal constant L2_LINK = IERC20(0xf97f4df75117a78c1A5a0DBb814Af92458539FB4);
//  IWrappedNative internal constant L2_WRAPPED_NATIVE = IWrappedNative(0x82aF49447D8a07e3bd95BD0d56f35241523fBab1);
//
//  uint256 internal constant TOKEN_BALANCE = 10e18;
//  address internal constant OWNER = address(0xdead);
//
//  ArbitrumL2BridgeAdapter internal s_l2BridgeAdapter;
//
//  function setUp() public {
//    vm.startPrank(OWNER);
//
//    arbitrumFork = vm.createFork(ARBITRUM_RPC_URL);
//
//    vm.selectFork(arbitrumFork);
//    s_l2BridgeAdapter = new ArbitrumL2BridgeAdapter(IL2GatewayRouter(L2_GATEWAY_ROUTER));
//    deal(address(L2_LINK), OWNER, TOKEN_BALANCE);
//    deal(address(L2_WRAPPED_NATIVE), OWNER, TOKEN_BALANCE);
//
//    vm.label(OWNER, "Owner");
//    vm.label(L2_GATEWAY_ROUTER, "L2GatewayRouterProxy");
//    vm.label(0xe80eb0238029333e368e0bDDB7acDf1b9cb28278, "L2GatewayRouter");
//    vm.label(L2_ETH_WITHDRAWAL_PRECOMPILE, "Precompile: ArbSys");
//  }
//}
//
//contract ArbitrumL2BridgeAdapter_sendERC20 is ArbitrumL2BridgeAdapterSetup {
//  function test_sendERC20Success() public {
//    L2_LINK.approve(address(s_l2BridgeAdapter), TOKEN_BALANCE);
//
//    s_l2BridgeAdapter.sendERC20(address(L1_LINK), address(L2_LINK), OWNER, TOKEN_BALANCE);
//  }
//}
