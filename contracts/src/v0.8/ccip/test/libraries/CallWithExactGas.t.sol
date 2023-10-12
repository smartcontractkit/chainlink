// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IAny2EVMMessageReceiver} from "../../interfaces/IAny2EVMMessageReceiver.sol";

import {Internal} from "../../libraries/Internal.sol";
import {Client} from "../../libraries/Client.sol";
import {CallWithExactGas} from "../../libraries/CallWithExactGas.sol";
import {CallWithExactGasHelper} from "../helpers/CallWithExactGasHelper.sol";
import {BaseTest} from "../BaseTest.t.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";

contract CallWithExactGas_callWithExactGas is BaseTest {
  event MessageReceived();

  IAny2EVMMessageReceiver internal s_receiver;
  MaybeRevertMessageReceiver internal s_reverting_receiver;
  CallWithExactGasHelper internal s_caller;

  function setUp() public virtual override {
    BaseTest.setUp();

    s_receiver = new MaybeRevertMessageReceiver(false);
    s_reverting_receiver = new MaybeRevertMessageReceiver(true);

    s_caller = new CallWithExactGasHelper();
  }

  function test_CallWithExactGasSuccess() public {
    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: "1",
      sourceChainSelector: 1,
      sender: "",
      data: "",
      destTokenAmounts: new Client.EVMTokenAmount[](0)
    });
    bytes memory data = abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message);

    vm.expectEmit();
    emit MessageReceived();

    vm.expectCall(address(s_receiver), data);

    (bool success, ) = s_caller.callWithExactGas(
      data,
      address(s_receiver),
      100_000,
      Internal.MAX_RET_BYTES,
      Internal.GAS_FOR_CALL_EXACT_CHECK
    );

    assertTrue(success);
  }

  function testFuzz_CallWithExactGasReceiverErrorSuccess(uint16 testRetBytes) public {
    // Bound with upper limit, otherwise the test runs out of gas.
    testRetBytes = uint16(bound(testRetBytes, 0, Internal.MAX_RET_BYTES * 10));

    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: "1",
      sourceChainSelector: 1,
      sender: "",
      data: "",
      destTokenAmounts: new Client.EVMTokenAmount[](0)
    });
    bytes memory data = abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message);

    bytes memory errorData = new bytes(testRetBytes);
    for (uint256 i = 0; i < errorData.length; ++i) {
      errorData[i] = 0x01;
    }
    s_reverting_receiver.setErr(errorData);

    vm.expectCall(address(s_reverting_receiver), data);

    (bool success, bytes memory retData) = s_caller.callWithExactGas(
      data,
      address(s_reverting_receiver),
      100_000,
      Internal.MAX_RET_BYTES,
      Internal.GAS_FOR_CALL_EXACT_CHECK
    );

    assertFalse(success);

    bytes memory totalReturnData = abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, errorData);
    bytes memory expectedReturnData = totalReturnData;

    // If expected return data is longer than MAX_RET_BYTES, truncate it to MAX_RET_BYTES
    if (expectedReturnData.length > Internal.MAX_RET_BYTES) {
      expectedReturnData = new bytes(Internal.MAX_RET_BYTES);
      for (uint256 i = 0; i < Internal.MAX_RET_BYTES; ++i) {
        expectedReturnData[i] = totalReturnData[i];
      }
    }
    assertEq(expectedReturnData, retData);
  }

  function test_NoContractReverts() public {
    bytes memory data = abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, "");

    vm.expectRevert(CallWithExactGas.NoContract.selector);

    s_caller.callWithExactGas(data, address(1), 100_000, Internal.MAX_RET_BYTES, Internal.GAS_FOR_CALL_EXACT_CHECK);
  }

  function test_NoGasForCallExactCheckReverts() public {
    bytes memory data = abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, "");

    bytes memory payload = abi.encodeWithSelector(
      s_caller.callWithExactGas.selector,
      data,
      address(s_receiver),
      100_000,
      Internal.MAX_RET_BYTES,
      Internal.GAS_FOR_CALL_EXACT_CHECK
    );

    (bool success, bytes memory retData) = address(s_caller).call{gas: GAS_FOR_CALL_EXACT_CHECK - 1}(payload);
    assertFalse(success);
    assertEq(abi.encodeWithSelector(CallWithExactGas.NoGasForCallExactCheck.selector), retData);
  }

  function test_NotEnoughGasForCallReverts() public {
    bytes memory data = abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, "");

    vm.expectRevert(CallWithExactGas.NotEnoughGasForCall.selector);

    s_caller.callWithExactGas(
      data,
      address(s_receiver),
      type(uint256).max,
      Internal.MAX_RET_BYTES,
      Internal.GAS_FOR_CALL_EXACT_CHECK
    );
  }
}
