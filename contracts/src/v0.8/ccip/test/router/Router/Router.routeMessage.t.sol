// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IAny2EVMMessageReceiver} from "../../../interfaces/IAny2EVMMessageReceiver.sol";
import {IRouter} from "../../../interfaces/IRouter.sol";

import {Router} from "../../../Router.sol";
import {Client} from "../../../libraries/Client.sol";
import {MaybeRevertMessageReceiver} from "../../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {OffRampSetup} from "../../offRamp/OffRamp/OffRampSetup.t.sol";

contract Router_routeMessage is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    vm.startPrank(address(s_offRamp));
  }

  function _generateManualGasLimit(
    uint256 callDataLength
  ) internal view returns (uint256) {
    return ((gasleft() - 2 * (16 * callDataLength + GAS_FOR_CALL_EXACT_CHECK)) * 62) / 64;
  }

  function _generateReceiverMessage(
    uint64 chainSelector
  ) internal pure returns (Client.Any2EVMMessage memory) {
    Client.EVMTokenAmount[] memory ta = new Client.EVMTokenAmount[](0);
    return Client.Any2EVMMessage({
      messageId: bytes32("a"),
      sourceChainSelector: chainSelector,
      sender: bytes("a"),
      data: bytes("a"),
      destTokenAmounts: ta
    });
  }

  function test_routeMessage_ManualExec() public {
    Client.Any2EVMMessage memory message = _generateReceiverMessage(SOURCE_CHAIN_SELECTOR);
    // Manuel execution cannot run out of gas

    (bool success, bytes memory retData, uint256 gasUsed) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR),
      GAS_FOR_CALL_EXACT_CHECK,
      _generateManualGasLimit(message.data.length),
      address(s_receiver)
    );
    assertTrue(success);
    assertEq("", retData);
    assertGt(gasUsed, 3_000);
  }

  function test_routeMessage_ExecutionEvent() public {
    Client.Any2EVMMessage memory message = _generateReceiverMessage(SOURCE_CHAIN_SELECTOR);
    // Should revert with reason
    bytes memory realError1 = new bytes(2);
    realError1[0] = 0xbe;
    realError1[1] = 0xef;
    s_reverting_receiver.setErr(realError1);

    vm.expectEmit();
    emit Router.MessageExecuted(
      message.messageId,
      message.sourceChainSelector,
      address(s_offRamp),
      keccak256(abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message))
    );

    (bool success, bytes memory retData, uint256 gasUsed) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR),
      GAS_FOR_CALL_EXACT_CHECK,
      _generateManualGasLimit(message.data.length),
      address(s_reverting_receiver)
    );

    assertFalse(success);
    assertEq(abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, realError1), retData);
    assertGt(gasUsed, 2850);

    // Reason is truncated
    // Over the MAX_RET_BYTES limit (including offset and length word since we have a dynamic values), should be ignored
    bytes memory realError2 = new bytes(32 * 2 + 1);
    realError2[32 * 2 - 1] = 0xAA;
    realError2[32 * 2] = 0xFF;
    s_reverting_receiver.setErr(realError2);

    vm.expectEmit();
    emit Router.MessageExecuted(
      message.messageId,
      message.sourceChainSelector,
      address(s_offRamp),
      keccak256(abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message))
    );

    (success, retData, gasUsed) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR),
      GAS_FOR_CALL_EXACT_CHECK,
      _generateManualGasLimit(message.data.length),
      address(s_reverting_receiver)
    );

    assertFalse(success);
    assertEq(
      abi.encodeWithSelector(
        MaybeRevertMessageReceiver.CustomError.selector,
        uint256(32),
        uint256(realError2.length),
        uint256(0),
        uint256(0xAA)
      ),
      retData
    );
    assertGt(gasUsed, 3_000);

    // Should emit success
    vm.expectEmit();
    emit Router.MessageExecuted(
      message.messageId,
      message.sourceChainSelector,
      address(s_offRamp),
      keccak256(abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message))
    );

    (success, retData, gasUsed) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR),
      GAS_FOR_CALL_EXACT_CHECK,
      _generateManualGasLimit(message.data.length),
      address(s_receiver)
    );

    assertTrue(success);
    assertEq("", retData);
    assertGt(gasUsed, 3_000);
  }

  function testFuzz_routeMessage_ExecutionEvent_Success(
    bytes calldata error
  ) public {
    Client.Any2EVMMessage memory message = _generateReceiverMessage(SOURCE_CHAIN_SELECTOR);
    s_reverting_receiver.setErr(error);

    bytes memory expectedRetData;

    if (error.length >= 33) {
      uint256 cutOff = error.length > 64 ? 64 : error.length;
      vm.expectEmit();
      emit Router.MessageExecuted(
        message.messageId,
        message.sourceChainSelector,
        address(s_offRamp),
        keccak256(abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message))
      );
      expectedRetData = abi.encodeWithSelector(
        MaybeRevertMessageReceiver.CustomError.selector,
        uint256(32),
        uint256(error.length),
        bytes32(error[:32]),
        bytes32(error[32:cutOff])
      );
    } else {
      vm.expectEmit();
      emit Router.MessageExecuted(
        message.messageId,
        message.sourceChainSelector,
        address(s_offRamp),
        keccak256(abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message))
      );
      expectedRetData = abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, error);
    }

    (bool success, bytes memory retData,) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR),
      GAS_FOR_CALL_EXACT_CHECK,
      _generateManualGasLimit(message.data.length),
      address(s_reverting_receiver)
    );

    assertFalse(success);
    assertEq(expectedRetData, retData);
  }

  function test_routeMessage_AutoExec() public {
    (bool success,,) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR), GAS_FOR_CALL_EXACT_CHECK, 100_000, address(s_receiver)
    );

    assertTrue(success);

    (success,,) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR), GAS_FOR_CALL_EXACT_CHECK, 1, address(s_receiver)
    );

    // Can run out of gas, should return false
    assertFalse(success);
  }

  // Reverts
  function test_RevertWhen_routeMessage_OnlyOffRamp() public {
    vm.stopPrank();
    vm.startPrank(STRANGER);

    vm.expectRevert(IRouter.OnlyOffRamp.selector);
    s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR), GAS_FOR_CALL_EXACT_CHECK, 100_000, address(s_receiver)
    );
  }

  function test_RevertWhen_routeMessage_WhenNotHealthy() public {
    vm.mockCall(address(s_mockRMNRemote), abi.encodeWithSignature("isCursed()"), abi.encode(true));
    vm.expectRevert(Router.BadARMSignal.selector);
    s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR), GAS_FOR_CALL_EXACT_CHECK, 100_000, address(s_receiver)
    );
  }
}
