// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {CallWithExactGas} from "../../call/CallWithExactGas.sol";
import {CallWithExactGasHelper} from "./CallWithExactGasHelper.sol";
import {BaseTest} from "../BaseTest.t.sol";
import {GenericReceiver} from "../testhelpers/GenericReceiver.sol";

contract CallWithExactGasSetup is BaseTest {
  GenericReceiver internal s_receiver;
  CallWithExactGasHelper internal s_caller;
  uint256 internal constant DEFAULT_GAS_LIMIT = 20_000;
  uint16 internal constant DEFAULT_GAS_FOR_CALL_EXACT_CHECK = 5000;
  uint256 internal constant EXTCODESIZE_GAS_COST = 2600;

  function setUp() public virtual override {
    BaseTest.setUp();

    s_receiver = new GenericReceiver(false);
    s_caller = new CallWithExactGasHelper();
  }
}

contract CallWithExactGas__callWithExactGasSafeReturnData is CallWithExactGasSetup {
  function testFuzz_CallWithExactGasSafeReturnDataSuccess(bytes memory payload, bytes4 funcSelector) public {
    vm.pauseGasMetering();
    bytes memory data = abi.encodeWithSelector(funcSelector, payload);
    vm.assume(
      funcSelector != GenericReceiver.setRevert.selector &&
        funcSelector != GenericReceiver.setErr.selector &&
        funcSelector != 0x5100fc21 // s_toRevert(), which is public and therefore has a function selector
    );

    uint16 maxRetBytes = 0;

    vm.expectCall(address(s_receiver), data);
    vm.resumeGasMetering();

    (bool success, bytes memory retData, uint256 gasUsed) = s_caller.callWithExactGasSafeReturnData(
      data,
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      maxRetBytes,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    assertTrue(success);
    assertEq(retData.length, 0);
    assertGt(gasUsed, 500);
  }

  function test_CallWithExactGasSafeReturnDataExactGas() public {
    // The gas cost for `extcodesize`
    uint256 extcodesizeGas = EXTCODESIZE_GAS_COST;
    // The calculated overhead for retData initialization
    uint256 overheadForRetDataInit = 114;
    // The calculated overhead for otherwise unaccounted for gas usage
    uint256 overheadForCallWithExactGas = 416;

    bytes memory payload = abi.encodeWithSelector(
      s_caller.callWithExactGasSafeReturnData.selector,
      "",
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK,
      0
    );

    // Since only 63/64th of the gas gets passed, we compensate
    uint256 allowedGas = (DEFAULT_GAS_LIMIT + (DEFAULT_GAS_LIMIT / 64));

    allowedGas +=
      extcodesizeGas +
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK +
      overheadForRetDataInit +
      overheadForCallWithExactGas;

    // Due to EIP-150 we expect to lose 1/64, so we compensate for this
    allowedGas = (allowedGas * 64) / 63;

    vm.expectCall(address(s_receiver), "");
    (bool success, bytes memory retData) = address(s_caller).call{gas: allowedGas}(payload);

    assertTrue(success);
    (bool innerSuccess, bytes memory innerRetData, uint256 gasUsed) = abi.decode(retData, (bool, bytes, uint256));

    assertTrue(innerSuccess);
    assertEq(innerRetData.length, 0);
    assertGt(gasUsed, 500);
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

    (bool success, bytes memory retData, uint256 gasUsed) = s_caller.callWithExactGasSafeReturnData(
      data,
      address(s_receiver),
      DEFAULT_GAS_LIMIT * 10,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK,
      maxReturnBytes
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
    assertGt(gasUsed, 500);
  }

  function test_NoContractReverts() public {
    address addressWithoutContract = address(1337);

    vm.expectRevert(CallWithExactGas.NoContract.selector);

    s_caller.callWithExactGasSafeReturnData(
      "", // empty payload as it will revert well before needing it
      addressWithoutContract,
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK,
      0
    );
  }

  function test_NoGasForCallExactCheckReverts() public {
    bytes memory payload = abi.encodeWithSelector(
      s_caller.callWithExactGasSafeReturnData.selector,
      "", // empty payload as it will revert well before needing it
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK,
      0
    );

    (bool success, bytes memory retData) = address(s_caller).call{gas: DEFAULT_GAS_FOR_CALL_EXACT_CHECK - 1}(payload);
    assertFalse(success);
    assertEq(retData.length, CallWithExactGas.NoGasForCallExactCheck.selector.length);
    assertEq(abi.encodeWithSelector(CallWithExactGas.NoGasForCallExactCheck.selector), retData);
  }

  function test_NotEnoughGasForCallReverts() public {
    bytes memory payload = abi.encodeWithSelector(
      s_caller.callWithExactGasSafeReturnData.selector,
      "", // empty payload as it will revert well before needing it
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK,
      0
    );

    // Supply enough gas for the final call, the DEFAULT_GAS_FOR_CALL_EXACT_CHECK,
    // the extcodesize and account for EIP-150. This doesn't account for any other gas
    // usage, and will therefore fail because the checks and memory stored/loads
    // also cost gas.
    uint256 allowedGas = (DEFAULT_GAS_LIMIT + (DEFAULT_GAS_LIMIT / 64)) + DEFAULT_GAS_FOR_CALL_EXACT_CHECK;
    // extcodesize gas cost
    allowedGas += EXTCODESIZE_GAS_COST;
    // EIP-150
    allowedGas = (allowedGas * 64) / 63;

    // Expect this call to fail due to not having enough gas for the final call
    (bool success, bytes memory retData) = address(s_caller).call{gas: allowedGas}(payload);

    assertFalse(success);
    assertEq(retData.length, CallWithExactGas.NotEnoughGasForCall.selector.length);
    assertEq(abi.encodeWithSelector(CallWithExactGas.NotEnoughGasForCall.selector), retData);
  }
}
