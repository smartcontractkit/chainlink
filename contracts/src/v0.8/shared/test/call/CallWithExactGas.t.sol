// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {CallWithExactGas} from "../../call/CallWithExactGas.sol";
import {CallWithExactGasHelper} from "./CallWithExactGasHelper.sol";
import {BaseTest} from "../BaseTest.t.sol";
import {GenericReceiver} from "../testhelpers/GenericReceiver.sol";

contract CallWithExactGas_callWithExactGas is BaseTest {
  GenericReceiver internal s_receiver;
  CallWithExactGasHelper internal s_caller;
  uint256 internal constant DEFAULT_GAS_LIMIT = 20_000;
  uint16 internal constant DEFAULT_GAS_FOR_CALL_EXACT_CHECK = 5000;

  function setUp() public virtual override {
    BaseTest.setUp();

    s_receiver = new GenericReceiver(false);

    s_caller = new CallWithExactGasHelper();
  }

  /// forge-config: shared.fuzz.runs = 3200
  function test_CallWithExactGasSafeReturnDataSuccess(bytes memory payload, bytes4 funcSelector) public {
    bytes memory data = abi.encodeWithSelector(funcSelector, payload);
    vm.assume(funcSelector != GenericReceiver.setRevert.selector && funcSelector != GenericReceiver.setErr.selector);

    uint16 maxRetBytes = 0;

    vm.expectCall(address(s_receiver), data);

    (bool success, ) = s_caller.callWithExactGasSafeReturnData(
      data,
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      maxRetBytes,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    assertTrue(success);
  }

  function test_CallWithExactGasSafeReturnDataExactGas() public {
    uint256 gasLimit = 10_000;
    uint16 gasForCallExactCheck = 5_000;

    bytes memory payload = abi.encodeWithSelector(
      s_caller.callWithExactGasSafeReturnData.selector,
      "", // empty payload as it will revert well before needing it
      address(s_receiver),
      gasLimit,
      0,
      gasForCallExactCheck
    );

    // Since only 63/64th of the gas gets passed, we compensate
    uint256 allowedGas = (gasLimit + (gasLimit / 64)); // 10,156
    // We call `extcodesize` which costs 2600 gas
    allowedGas += 2600; //  10,156 + 2,600 = 12,756

    // Add DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    allowedGas += gasForCallExactCheck; // 12,756 + 5,000 = 17,756

    // Add gas to init the retData field, calculated to be 114 gas for 0 length
    allowedGas += 114; // 17,756 + 114 = 17,870

    // Add some margin for 5 mstore's, 1 gas() call, a function call, 3 slots of func args
    // and some basic arithmetic. Rough estimate of ~100 total.
    allowedGas += 100; // 17,870 + 100 = 17,970

    // Extra padding to handle e.g. calldata cost
    allowedGas += 559; // Magic padding required = 18,529

    // Due to EIP-150 we expect to lose 1/64, so we compensate for this
    allowedGas = (allowedGas * 64) / 63; // 18,529   * 64 / 63 = 18,823

    (bool success, ) = address(s_caller).call{gas: allowedGas}(payload);

    assertTrue(success);
  }

  function testFuzz_CallWithExactGasReceiverErrorSuccess(uint16 testRetBytes) public {
    uint16 maxReturnBytes = 500;
    // Bound with upper limit, otherwise the test runs out of gas.
    testRetBytes = uint16(bound(testRetBytes, 0, maxReturnBytes * 10));

    bytes memory data = abi.encode("0x52656E73");

    bytes memory errorData = new bytes(testRetBytes);
    for (uint256 i = 0; i < errorData.length; ++i) {
      errorData[i] = 0x01;
    }
    s_receiver.setErr(errorData);
    s_receiver.setRevert(true);

    vm.expectCall(address(s_receiver), data);

    (bool success, bytes memory retData) = s_caller.callWithExactGasSafeReturnData(
      data,
      address(s_receiver),
      DEFAULT_GAS_LIMIT * 10,
      maxReturnBytes,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    assertFalse(success);

    bytes memory expectedReturnData = errorData;

    // If expected return data is longer than MAX_RET_BYTES, truncate it to MAX_RET_BYTES
    if (expectedReturnData.length > maxReturnBytes) {
      expectedReturnData = new bytes(maxReturnBytes);
      for (uint256 i = 0; i < maxReturnBytes; ++i) {
        expectedReturnData[i] = errorData[i];
      }
    }
    assertEq(expectedReturnData, retData);
  }

  function test_NoContractReverts() public {
    address addressWithoutContract = address(1337);

    vm.expectRevert(CallWithExactGas.NoContract.selector);

    s_caller.callWithExactGasSafeReturnData(
      "", // empty payload as it will revert well before needing it
      addressWithoutContract,
      DEFAULT_GAS_LIMIT,
      0,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );
  }

  function test_NoGasForCallExactCheckReverts() public {
    bytes memory payload = abi.encodeWithSelector(
      s_caller.callWithExactGasSafeReturnData.selector,
      "", // empty payload as it will revert well before needing it
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      0,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    (bool success, bytes memory retData) = address(s_caller).call{gas: DEFAULT_GAS_FOR_CALL_EXACT_CHECK - 1}(payload);
    assertFalse(success);
    assertEq(retData.length, CallWithExactGas.NoGasForCallExactCheck.selector.length);
    assertEq(abi.encodeWithSelector(CallWithExactGas.NoGasForCallExactCheck.selector), retData);
  }

  function test_NotEnoughGasForCallReverts() public {
    //    vm.expectRevert(CallWithExactGas.NotEnoughGasForCall.selector);
    //
    //    s_caller.callWithExactGasSafeReturnData(
    //      "",
    //      address(s_receiver),
    //      type(uint256).max,
    //      0,
    //      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    //    );

    bytes memory payload = abi.encodeWithSelector(
      s_caller.callWithExactGasSafeReturnData.selector,
      "", // empty payload as it will revert well before needing it
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      0,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    uint256 allowedGas = (DEFAULT_GAS_LIMIT + (DEFAULT_GAS_LIMIT / 64)) + DEFAULT_GAS_FOR_CALL_EXACT_CHECK;

    (bool success, bytes memory retData) = address(s_caller).call{gas: allowedGas}(payload);

    assertFalse(success);
    assertEq(retData.length, CallWithExactGas.NotEnoughGasForCall.selector.length);
    assertEq(abi.encodeWithSelector(CallWithExactGas.NotEnoughGasForCall.selector), retData);
  }
}
