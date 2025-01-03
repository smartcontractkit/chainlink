// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {CallWithExactGas} from "../../call/CallWithExactGas.sol";
import {CallWithExactGasHelper} from "./CallWithExactGasHelper.sol";
import {BaseTest} from "../BaseTest.t.sol";
import {GenericReceiver} from "../testhelpers/GenericReceiver.sol";
import {GasConsumer} from "../testhelpers/GasConsumer.sol";

contract CallWithExactGasSetup is BaseTest {
  GenericReceiver internal s_receiver;
  CallWithExactGasHelper internal s_caller;
  GasConsumer internal s_gasConsumer;

  uint256 internal constant DEFAULT_GAS_LIMIT = 20_000;
  uint16 internal constant DEFAULT_GAS_FOR_CALL_EXACT_CHECK = 5000;
  uint256 internal constant EXTCODESIZE_GAS_COST = 2600;

  function setUp() public virtual override {
    BaseTest.setUp();

    s_receiver = new GenericReceiver(false);
    s_caller = new CallWithExactGasHelper();
    s_gasConsumer = new GasConsumer();
  }
}

contract CallWithExactGas__callWithExactGas is CallWithExactGasSetup {
  function testFuzz_callWithExactGasSuccess(bytes memory payload, bytes4 funcSelector) public {
    vm.pauseGasMetering();

    bytes memory data = abi.encodeWithSelector(funcSelector, payload);
    vm.assume(
      funcSelector != GenericReceiver.setRevert.selector &&
        funcSelector != GenericReceiver.setErr.selector &&
        funcSelector != 0x5100fc21 // s_toRevert(), which is public and therefore has a function selector
    );

    vm.expectCall(address(s_receiver), data);
    vm.resumeGasMetering();

    bool success = s_caller.callWithExactGas(
      data,
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    assertTrue(success);
  }

  function test_CallWithExactGasSafeReturnDataExactGas() public {
    // The calculated overhead for otherwise unaccounted for gas usage
    uint256 overheadForCallWithExactGas = 364;

    bytes memory payload = abi.encodeWithSelector(
      s_caller.callWithExactGas.selector,
      "",
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    // Since only 63/64th of the gas gets passed, we compensate
    uint256 allowedGas = (DEFAULT_GAS_LIMIT + (DEFAULT_GAS_LIMIT / 64));

    allowedGas += EXTCODESIZE_GAS_COST + DEFAULT_GAS_FOR_CALL_EXACT_CHECK + overheadForCallWithExactGas;

    // Due to EIP-150 we expect to lose 1/64, so we compensate for this
    allowedGas = (allowedGas * 64) / 63;

    (bool success, bytes memory retData) = address(s_caller).call{gas: allowedGas}(payload);

    assertTrue(success);
    assertEq(abi.encode(true), retData);
  }

  function test_CallWithExactGasReceiverErrorSuccess() public {
    bytes memory data = abi.encode("0x52656E73");

    bytes memory errorData = new bytes(20);
    for (uint256 i = 0; i < errorData.length; ++i) {
      errorData[i] = 0x01;
    }
    s_receiver.setErr(errorData);
    s_receiver.setRevert(true);

    vm.expectCall(address(s_receiver), data);

    bool success = s_caller.callWithExactGas(
      data,
      address(s_receiver),
      DEFAULT_GAS_LIMIT * 10,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    assertFalse(success);
  }

  function test_NoContractReverts() public {
    address addressWithoutContract = address(1337);

    vm.expectRevert(CallWithExactGas.NoContract.selector);

    s_caller.callWithExactGas(
      "", // empty payload as it will revert well before needing it
      addressWithoutContract,
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );
  }

  function test_NoGasForCallExactCheckReverts() public {
    bytes memory payload = abi.encodeWithSelector(
      s_caller.callWithExactGas.selector,
      "", // empty payload as it will revert well before needing it
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    (bool success, bytes memory retData) = address(s_caller).call{gas: DEFAULT_GAS_FOR_CALL_EXACT_CHECK - 1}(payload);
    assertFalse(success);
    assertEq(retData.length, CallWithExactGas.NoGasForCallExactCheck.selector.length);
    assertEq(abi.encodeWithSelector(CallWithExactGas.NoGasForCallExactCheck.selector), retData);
  }

  function test_NotEnoughGasForCallReverts() public {
    bytes memory payload = abi.encodeWithSelector(
      s_caller.callWithExactGas.selector,
      "", // empty payload as it will revert well before needing it
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
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
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK,
      maxRetBytes
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
    uint256 overheadForCallWithExactGas = 486;

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

  function testFuzz_CallWithExactGasSafeReturnData_ConsumeAllGas_Success(uint8 gasLimitMultiplier) external {
    vm.assume(gasLimitMultiplier > 0); // Assume not zero to avoid zero gas being passed to s_gasConsumer
    uint16 maxRetBytes = 0;

    vm.expectCall(address(s_gasConsumer), abi.encodeWithSelector(s_gasConsumer.consumeAllGas.selector));

    // Default gas limit of 20k times the amount generated by the fuzzer
    uint256 gasLimit = DEFAULT_GAS_LIMIT * gasLimitMultiplier;

    (bool success, , uint256 gasUsed) = s_caller.callWithExactGasSafeReturnData(
      abi.encodePacked(s_gasConsumer.consumeAllGas.selector),
      address(s_gasConsumer),
      gasLimit,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK,
      maxRetBytes
    );

    // Was determined as the least amount of gas needed for pre-call stack management.
    // when an OOG error is thrown, used gas returned by the library is 122 <= x <= 137 based on the
    // gas limit passed in. This was determined through trial and error and fuzzing.
    uint256 CALL_WITH_EXACT_GAS_SAFE_RETURN_DATA_GAS_OVERHEAD = 122;

    assertTrue(success, "Error: External Call Failed");

    // Assert equal within a margin of error of 1/64 of the gas limit to account for excess gas used by execution library
    assertApproxEqAbs(
      gasUsed - CALL_WITH_EXACT_GAS_SAFE_RETURN_DATA_GAS_OVERHEAD,
      gasLimit,
      gasLimit / 64,
      "Error: All gas not consumed by receiver"
    );
  }

  function test_callWithExactGasSafeReturnData_ThrowOOGError_Revert() external {
    uint16 maxRetBytes = 0;

    vm.expectCall(address(s_gasConsumer), abi.encodeWithSelector(s_gasConsumer.throwOutOfGasError.selector));

    (bool success, , ) = s_caller.callWithExactGasSafeReturnData(
      abi.encodePacked(s_gasConsumer.throwOutOfGasError.selector),
      address(s_gasConsumer),
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK,
      maxRetBytes
    );

    assertFalse(success, "Error: External Call Succeeded where it should not");
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

contract CallWithExactGas__callWithExactGasEvenIfTargetIsNoContract is CallWithExactGasSetup {
  function testFuzz_CallWithExactGasEvenIfTargetIsNoContractSuccess(bytes memory payload, bytes4 funcSelector) public {
    vm.pauseGasMetering();
    bytes memory data = abi.encodeWithSelector(funcSelector, payload);
    vm.assume(
      funcSelector != GenericReceiver.setRevert.selector &&
        funcSelector != GenericReceiver.setErr.selector &&
        funcSelector != 0x5100fc21 // s_toRevert(), which is public and therefore has a function selector
    );
    vm.expectCall(address(s_receiver), data);
    vm.resumeGasMetering();

    (bool success, bool sufficientGas) = s_caller.callWithExactGasEvenIfTargetIsNoContract(
      data,
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    assertTrue(success);
    assertTrue(sufficientGas);
  }

  function test_CallWithExactGasEvenIfTargetIsNoContractExactGasSuccess() public {
    // The calculated overhead for otherwise unaccounted for gas usage
    uint256 overheadForCallWithExactGas = 446;

    bytes memory data = abi.encode("0x52656E73");

    bytes memory payload = abi.encodeWithSelector(
      s_caller.callWithExactGasEvenIfTargetIsNoContract.selector,
      data,
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    // Since only 63/64th of the gas gets passed, we compensate
    uint256 allowedGas = (DEFAULT_GAS_LIMIT + (DEFAULT_GAS_LIMIT / 64));

    allowedGas += DEFAULT_GAS_FOR_CALL_EXACT_CHECK + overheadForCallWithExactGas;

    // Due to EIP-150 we expect to lose 1/64, so we compensate for this
    allowedGas = (allowedGas * 64) / 63;

    vm.expectCall(address(s_receiver), data);
    (bool outerCallSuccess, bytes memory SuccessAndSufficientGas) = address(s_caller).call{gas: allowedGas}(payload);

    // The call succeeds
    assertTrue(outerCallSuccess);

    (bool success, bool sufficientGas) = abi.decode(SuccessAndSufficientGas, (bool, bool));
    assertTrue(success);
    assertTrue(sufficientGas);
  }

  function test_CallWithExactGasEvenIfTargetIsNoContractReceiverErrorSuccess() public {
    bytes memory data = abi.encode("0x52656E73");

    bytes memory errorData = new bytes(20);
    for (uint256 i = 0; i < errorData.length; ++i) {
      errorData[i] = 0x01;
    }
    s_receiver.setErr(errorData);
    s_receiver.setRevert(true);

    vm.expectCall(address(s_receiver), data);

    (bool success, bool sufficientGas) = s_caller.callWithExactGasEvenIfTargetIsNoContract(
      data,
      address(s_receiver),
      DEFAULT_GAS_LIMIT * 10,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    // We don't care if it reverts, we only care if we have enough gas
    assertFalse(success);
    assertTrue(sufficientGas);
  }

  function test_NoContractSuccess() public {
    bytes memory data = abi.encode("0x52656E73");
    address addressWithoutContract = address(1337);

    (bool success, bool sufficientGas) = s_caller.callWithExactGasEvenIfTargetIsNoContract(
      data,
      addressWithoutContract,
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    assertTrue(success);
    assertTrue(sufficientGas);
  }

  function test_NoGasForCallExactCheckReturnFalseSuccess() public {
    bytes memory payload = abi.encodeWithSelector(
      s_caller.callWithExactGasEvenIfTargetIsNoContract.selector,
      "", // empty payload as it will revert well before needing it
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    (bool outerCallSuccess, bytes memory SuccessAndSufficientGas) = address(s_caller).call{
      gas: DEFAULT_GAS_FOR_CALL_EXACT_CHECK - 1
    }(payload);

    // The call succeeds
    assertTrue(outerCallSuccess);

    (bool success, bool sufficientGas) = abi.decode(SuccessAndSufficientGas, (bool, bool));
    assertFalse(success);
    assertFalse(sufficientGas);
  }

  function test_NotEnoughGasForCallReturnsFalseSuccess() public {
    bytes memory payload = abi.encodeWithSelector(
      s_caller.callWithExactGasEvenIfTargetIsNoContract.selector,
      "", // empty payload as it will revert well before needing it
      address(s_receiver),
      DEFAULT_GAS_LIMIT,
      DEFAULT_GAS_FOR_CALL_EXACT_CHECK
    );

    // Supply enough gas for the final call, the DEFAULT_GAS_FOR_CALL_EXACT_CHECK,
    // and account for EIP-150. This doesn't account for any other gas usage, and
    // will therefore fail because the checks and memory stored/loads also cost gas.
    uint256 allowedGas = (DEFAULT_GAS_LIMIT + (DEFAULT_GAS_LIMIT / 64)) + DEFAULT_GAS_FOR_CALL_EXACT_CHECK;
    // EIP-150
    allowedGas = (allowedGas * 64) / 63;

    // Expect this call to fail due to not having enough gas for the final call
    (bool outerCallSuccess, bytes memory SuccessAndSufficientGas) = address(s_caller).call{gas: allowedGas}(payload);

    // The call succeeds
    assertTrue(outerCallSuccess);

    (bool success, bool sufficientGas) = abi.decode(SuccessAndSufficientGas, (bool, bool));
    assertFalse(success);
    assertFalse(sufficientGas);
  }
}
