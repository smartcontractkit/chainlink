// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {CallWithExactGasZKSync} from "../../call/CallWithExactGasZKSync.sol";
import {CallWithExactGasZKSyncHelper} from "./CallWithExactGasZKSyncHelper.sol";
import {BaseTest} from "../BaseTest.t.sol";

import {MockSystemContext} from "../mocks/MockSystemContext.sol";
import {TestTarget} from "../testhelpers/TestTarget.sol";

contract CallWithExactGasZKSyncSetup is BaseTest {
  CallWithExactGasZKSyncHelper internal s_helper;
  MockSystemContext internal s_mockSystemContext;
  TestTarget internal s_target;

  // Import the errors from the library (for vm.expectRevert checks)
  error NoContract();
  error NotEnoughGasForPubdata();
  error NotEnoughGasForCall();

  function setUp() public virtual override {
    s_mockSystemContext = new MockSystemContext();
    // Write mock's code to 0x800b so library calls see it
    vm.etch(address(0x800b), address(s_mockSystemContext).code);

    s_helper = new CallWithExactGasZKSyncHelper();
    s_target = new TestTarget();
  }

  function _limitedGasCallWithExactGas(
    uint256 allowedGas,
    address _to,
    uint256 _maxTotalGas,
    bytes memory _data,
    uint16 _maxReturnBytes
  ) internal returns (bool success, bytes memory retData) {
    // Encode the call to the helper function:
    bytes memory payload = abi.encodeWithSelector(
      CallWithExactGasZKSyncHelper.callWithExactGasSafeReturnData.selector,
      _data, // bytes
      _to, // address
      _maxTotalGas, // uint256
      _maxReturnBytes // uint16
    );

    // Constrain the subcall to `allowedGas`
    (success, retData) = address(s_helper).call{gas: allowedGas}(payload);

    return (success, retData);
  }

  function _decodeResult(
    bytes memory retData
  ) internal pure returns (bool callSuccess, bytes memory callRetData, uint256 pubdataGasSpent) {
    // The helper returns (bool, bytes, uint256)
    return abi.decode(retData, (bool, bytes, uint256));
  }
}

contract CallWithExactGasZKSync__callWithExactGasSafeReturnData is CallWithExactGasZKSyncSetup {
  function test__callWithExactGasSafeReturnData_RevertWhen_NoContract() public {
    // Expect custom error NoContract()
    vm.expectRevert(NoContract.selector);

    // We'll allow ~2e6 gas for the subcall, which is plenty so we
    // don't trigger NotEnoughGasForCall by accident.
    _limitedGasCallWithExactGas(
      2_000_000,
      address(1234), // no code
      1_000_000, // _maxTotalGas
      abi.encodeWithSelector(TestTarget.returnData.selector),
      100
    );
  }

  function test__callWithExactGasSafeReturnData_RevertWhen_NotEnoughGasForCall() public {
    // If subcall has 500k gas, but we pass _maxTotalGas=100k, that's < 500k => revert.
    vm.expectRevert(NotEnoughGasForCall.selector);

    _limitedGasCallWithExactGas(
      500_000, // subcall has ~500k gas
      address(s_target),
      100_000, // _maxTotalGas is only 100k => triggers NotEnoughGasForCall
      abi.encodeWithSelector(TestTarget.returnData.selector),
      100
    );
  }

  function test__callWithExactGasSafeReturnData_RevertWhen_NotEnoughGasForPubdata() public {
    // We'll simulate pubdata usage:
    // s_mockSystemContext already starts with 0, let's set it to something
    s_mockSystemContext.setCurrentPubdataSpent(1000);

    // Then we mock the "after" usage to 5000 in the same call
    vm.mockCall(
      address(s_mockSystemContext),
      abi.encodeWithSelector(s_mockSystemContext.getCurrentPubdataSpent.selector),
      abi.encode(5000)
    );

    // This difference = 4000 pubdata * 10 gas/byte = 40,000 extra gas needed
    // We'll give the subcall 200,000 gas in total, and set _maxTotalGas=200k
    // But the overhead in the library might push it to revert for pubdata anyway.
    // We expect NotEnoughGasForPubdata now, not the NotEnoughGasForCall.

    vm.expectRevert(NotEnoughGasForPubdata.selector);

    _limitedGasCallWithExactGas(
      200_000, // subcall gas
      address(s_target),
      200_000, // library sees ~200k
      abi.encodeWithSelector(TestTarget.returnData.selector),
      100
    );
  }

  function test__callWithExactGasSafeReturnData_Success() public {
    // We'll provide the subcall with 1 million gas,
    // and pass _maxTotalGas=1e6 => no "NotEnoughGasForCall" revert
    (bool successCall, bytes memory retData) = _limitedGasCallWithExactGas(
      1_000_000_000_000,
      address(s_target),
      100_000_000_000_000,
      abi.encodeWithSelector(TestTarget.returnData.selector),
      10000
    );
    assertTrue(successCall, "Subcall itself must not revert");
    (bool success, bytes memory returnedData, uint256 pubdata) = _decodeResult(retData);

    assertTrue(success, "Target call must succeed");
    assertEq(pubdata, 0, "No extra pubdata usage expected");
    assertNotEq(returnedData.length, 0, "Should have returned some data");
    assertEq(abi.decode(returnedData, (string)), "Hello from TestTarget");
  }

  function test__callWithExactGasSafeReturnData_TruncatesData() public {
    (bool successCall, bytes memory retData) = _limitedGasCallWithExactGas(
      10_000_000,
      address(s_target),
      1_000_000_000,
      abi.encodeWithSelector(TestTarget.returnLargeData.selector),
      50 // only allow 50 bytes of return data
    );

    assertTrue(successCall, "Subcall must not revert");
    (bool success, bytes memory returnedData, ) = _decodeResult(retData);
    assertTrue(success, "Target call must succeed");
    assertEq(returnedData.length, 50, "Should have truncated the large data to 50 bytes");
  }

  function test__callWithExactGasSafeReturnData_RevertWhen_TargetRevertsWithReason() public {
    // We'll expect the revert reason "CustomRevertReason"
    vm.expectRevert(bytes("CustomRevertReason"));

    _limitedGasCallWithExactGas(
      1_000_000,
      address(s_target),
      1_000_000,
      abi.encodeWithSelector(TestTarget.revertWithReason.selector),
      100
    );
  }

  function test__callWithExactGasSafeReturnData_RevertWhen_TargetRevertsNoReason() public {
    vm.expectRevert(); // just expect some revert, no reason
    _limitedGasCallWithExactGas(
      1_000_000,
      address(s_target),
      1_000_000,
      abi.encodeWithSelector(TestTarget.revertNoReason.selector),
      100
    );
  }
}
