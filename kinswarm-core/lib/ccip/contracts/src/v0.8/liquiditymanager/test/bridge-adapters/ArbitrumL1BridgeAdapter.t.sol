// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IWrappedNative} from "../../../ccip/interfaces/IWrappedNative.sol";

import {ArbitrumL1BridgeAdapter, IOutbox} from "../../bridge-adapters/ArbitrumL1BridgeAdapter.sol";
import "forge-std/Test.sol";

import {IL1GatewayRouter} from "@arbitrum/token-bridge-contracts/contracts/tokenbridge/ethereum/gateway/IL1GatewayRouter.sol";
import {IGatewayRouter} from "@arbitrum/token-bridge-contracts/contracts/tokenbridge/libraries/gateway/IGatewayRouter.sol";
import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

//contract ArbitrumL1BridgeAdapterSetup is Test {
//  uint256 internal mainnetFork;
//  uint256 internal arbitrumFork;
//
//  string internal constant MAINNET_RPC_URL = "<url>";
//
//  address internal constant L1_GATEWAY_ROUTER = 0x72Ce9c846789fdB6fC1f34aC4AD25Dd9ef7031ef;
//  address internal constant L1_ERC20_GATEWAY = 0xa3A7B6F88361F48403514059F1F16C8E78d60EeC;
//  address internal constant L1_INBOX = 0x4Dbd4fc535Ac27206064B68FfCf827b0A60BAB3f;
//  // inbox 0x5aED5f8A1e3607476F1f81c3d8fe126deB0aFE94?
//  address internal constant L1_OUTBOX = 0x0B9857ae2D4A3DBe74ffE1d7DF045bb7F96E4840;
//
//  IERC20 internal constant L1_LINK = IERC20(0x514910771AF9Ca656af840dff83E8264EcF986CA);
//  IWrappedNative internal constant L1_WRAPPED_NATIVE = IWrappedNative(0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2);
//
//  address internal constant L2_GATEWAY_ROUTER = 0x5288c571Fd7aD117beA99bF60FE0846C4E84F933;
//  address internal constant L2_ETH_WITHDRAWAL_PRECOMPILE = 0x0000000000000000000000000000000000000064;
//
//  IERC20 internal constant L2_LINK = IERC20(0xf97f4df75117a78c1A5a0DBb814Af92458539FB4);
//  IWrappedNative internal constant L2_WRAPPED_NATIVE = IWrappedNative(0x82aF49447D8a07e3bd95BD0d56f35241523fBab1);
//
//  ArbitrumL1BridgeAdapter internal s_l1BridgeAdapter;
//
//  uint256 internal constant TOKEN_BALANCE = 10e18;
//  address internal constant OWNER = address(0xdead);
//
//  function setUp() public {
//    vm.startPrank(OWNER);
//
//    mainnetFork = vm.createFork(MAINNET_RPC_URL);
//    vm.selectFork(mainnetFork);
//
//    s_l1BridgeAdapter = new ArbitrumL1BridgeAdapter(
//      IL1GatewayRouter(L1_GATEWAY_ROUTER),
//      IOutbox(L1_OUTBOX),
//      L1_ERC20_GATEWAY
//    );
//
//    deal(address(L1_LINK), OWNER, TOKEN_BALANCE);
//    deal(address(L1_WRAPPED_NATIVE), OWNER, TOKEN_BALANCE);
//
//    vm.label(OWNER, "Owner");
//    vm.label(L1_GATEWAY_ROUTER, "L1GatewayRouter");
//    vm.label(L1_ERC20_GATEWAY, "L1 ERC20 Gateway");
//  }
//}
//
//contract ArbitrumL1BridgeAdapter_sendERC20 is ArbitrumL1BridgeAdapterSetup {
//  event TransferRouted(address indexed token, address indexed _userFrom, address indexed _userTo, address gateway);
//
//  function test_sendERC20Success() public {
//    L1_LINK.approve(address(s_l1BridgeAdapter), TOKEN_BALANCE);
//
//    vm.expectEmit();
//    emit TransferRouted(address(L1_LINK), address(s_l1BridgeAdapter), OWNER, L1_ERC20_GATEWAY);
//
//    uint256 expectedCost = s_l1BridgeAdapter.MAX_GAS() *
//      s_l1BridgeAdapter.GAS_PRICE_BID() +
//      s_l1BridgeAdapter.MAX_SUBMISSION_COST();
//
//    s_l1BridgeAdapter.sendERC20{value: expectedCost}(address(L1_LINK), OWNER, OWNER, TOKEN_BALANCE);
//  }
//
//  function test_BridgeFeeTooLowReverts() public {
//    L1_LINK.approve(address(s_l1BridgeAdapter), TOKEN_BALANCE);
//    uint256 expectedCost = s_l1BridgeAdapter.MAX_GAS() *
//      s_l1BridgeAdapter.GAS_PRICE_BID() +
//      s_l1BridgeAdapter.MAX_SUBMISSION_COST();
//
//    vm.expectRevert(
//      abi.encodeWithSelector(ArbitrumL1BridgeAdapter.InsufficientEthValue.selector, expectedCost, expectedCost - 1)
//    );
//
//    s_l1BridgeAdapter.sendERC20{value: expectedCost - 1}(address(L1_LINK), OWNER, OWNER, TOKEN_BALANCE);
//  }
//
//  function test_noApprovalReverts() public {
//    uint256 expectedCost = s_l1BridgeAdapter.MAX_GAS() *
//      s_l1BridgeAdapter.GAS_PRICE_BID() +
//      s_l1BridgeAdapter.MAX_SUBMISSION_COST();
//
//    vm.expectRevert("SafeERC20: low-level call failed");
//
//    s_l1BridgeAdapter.sendERC20{value: expectedCost}(address(L1_LINK), OWNER, OWNER, TOKEN_BALANCE);
//  }
//}
