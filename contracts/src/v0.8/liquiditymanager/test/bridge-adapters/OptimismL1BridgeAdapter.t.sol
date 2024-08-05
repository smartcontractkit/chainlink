// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import "forge-std/Test.sol";

import {IWrappedNative} from "../../../ccip/interfaces/IWrappedNative.sol";
import {WETH9} from "../../../ccip/test/WETH9.sol";
import {OptimismL1BridgeAdapter} from "../../bridge-adapters/OptimismL1BridgeAdapter.sol";
import {Types} from "../../interfaces/optimism/Types.sol";
import {IOptimismPortal} from "../../interfaces/optimism/IOptimismPortal.sol";

import {IL1StandardBridge} from "@eth-optimism/contracts/L1/messaging/IL1StandardBridge.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract OptimismL1BridgeAdapterSetup is Test {
  // addresses below are fake
  address internal constant L1_STANDARD_BRIDGE = address(1234);
  address internal constant OP_PORTAL = address(4567);
  address internal constant OWNER = address(0xdead);

  OptimismL1BridgeAdapter internal s_adapter;

  function setUp() public {
    vm.startPrank(OWNER);

    // deploy wrapped native
    WETH9 weth = new WETH9();

    // deploy bridge adapter
    s_adapter = new OptimismL1BridgeAdapter(
      IL1StandardBridge(L1_STANDARD_BRIDGE),
      IWrappedNative(address(weth)),
      IOptimismPortal(OP_PORTAL)
    );
  }
}

contract OptimismL1BridgeAdapter_finalizeWithdrawERC20 is OptimismL1BridgeAdapterSetup {
  function testfinalizeWithdrawERC20proveWithdrawalSuccess() public {
    // prepare payload
    OptimismL1BridgeAdapter.OptimismProveWithdrawalPayload memory provePayload = OptimismL1BridgeAdapter
      .OptimismProveWithdrawalPayload({
        withdrawalTransaction: Types.WithdrawalTransaction({
          nonce: 1,
          sender: address(0xdead),
          target: address(0xbeef),
          value: 1234,
          gasLimit: 4567,
          data: hex"deadbeef"
        }),
        l2OutputIndex: 1234,
        outputRootProof: Types.OutputRootProof({
          version: bytes32(0),
          stateRoot: bytes32(uint256(500)),
          messagePasserStorageRoot: bytes32(uint256(600)),
          latestBlockhash: bytes32(uint256(700))
        }),
        withdrawalProof: new bytes[](0)
      });
    OptimismL1BridgeAdapter.FinalizeWithdrawERC20Payload memory payload;
    payload.action = OptimismL1BridgeAdapter.FinalizationAction.ProveWithdrawal;
    payload.data = abi.encode(provePayload);

    bytes memory encodedPayload = abi.encode(payload);

    // mock out call to optimism portal
    vm.mockCall(
      OP_PORTAL,
      abi.encodeWithSelector(
        IOptimismPortal.proveWithdrawalTransaction.selector,
        provePayload.withdrawalTransaction,
        provePayload.l2OutputIndex,
        provePayload.outputRootProof,
        provePayload.withdrawalProof
      ),
      ""
    );

    // call finalizeWithdrawERC20
    s_adapter.finalizeWithdrawERC20(address(0), address(0), encodedPayload);
  }

  function testfinalizeWithdrawERC20FinalizeSuccess() public {
    // prepare payload
    OptimismL1BridgeAdapter.OptimismFinalizationPayload memory finalizePayload = OptimismL1BridgeAdapter
      .OptimismFinalizationPayload({
        withdrawalTransaction: Types.WithdrawalTransaction({
          nonce: 1,
          sender: address(0xdead),
          target: address(0xbeef),
          value: 1234,
          gasLimit: 4567,
          data: hex"deadbeef"
        })
      });
    OptimismL1BridgeAdapter.FinalizeWithdrawERC20Payload memory payload;
    payload.action = OptimismL1BridgeAdapter.FinalizationAction.FinalizeWithdrawal;
    payload.data = abi.encode(finalizePayload);

    bytes memory encodedPayload = abi.encode(payload);

    // mock out call to optimism portal
    vm.mockCall(
      OP_PORTAL,
      abi.encodeWithSelector(
        IOptimismPortal.finalizeWithdrawalTransaction.selector,
        finalizePayload.withdrawalTransaction
      ),
      ""
    );

    // call finalizeWithdrawERC20
    s_adapter.finalizeWithdrawERC20(address(0), address(0), encodedPayload);
  }

  function testFinalizeWithdrawERC20Reverts() public {
    // case 1: badly encoded payload
    bytes memory payload = abi.encode(1, 2, 3);
    vm.expectRevert();
    s_adapter.finalizeWithdrawERC20(address(0), address(0), payload);

    // case 2: invalid action
    // can't prepare the payload in solidity
    payload = hex"0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004deadbeef00000000000000000000000000000000000000000000000000000000";
    vm.expectRevert();
    s_adapter.finalizeWithdrawERC20(address(0), address(0), payload);
  }
}
