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
      _to,
      _maxTotalGas,
      _data,
      _maxReturnBytes
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
  /// @notice Reverts if target has no code => "NoContract()"
  function test__callWithExactGasSafeReturnData_RevertWhen_NoContract() public {
    (bool successCall, bytes memory retData) = _limitedGasCallWithExactGas(
      2_000_000,
      address(12345), // no code
      1_000_000,
      abi.encodeWithSelector(TestTarget.returnData.selector),
      100
    );
    assertFalse(successCall, "Subcall itself must revert");

    if (retData.length >= 4) {
      bytes4 errSig;
      assembly {
        errSig := mload(add(retData, 32))
      }
      require(errSig == NoContract.selector, "Unexpected revert error");
    }
  }

  /// @notice Reverts if _maxTotalGas is greater than gasleft()
  function test__callWithExactGasSafeReturnData_FailsWhen_NotEnoughGasForCall() public {
    (bool successCall, bytes memory retData) = _limitedGasCallWithExactGas(
      500_000, // subcall has ~500k gas available
      address(s_target),
      600_000, // _maxTotalGas exceeds available gas => triggers NotEnoughGasForCall
      abi.encodeWithSelector(TestTarget.returnData.selector),
      100
    );
    assertTrue(successCall, "Subcall itself must not revert");
    (bool success, bytes memory returnedData, uint256 pubdata) = _decodeResult(retData);

    assertFalse(success, "Target call must fail");
    assertEq(pubdata, 0, "No extra pubdata usage expected");
    assertEq(returnedData.length, 0, "Should not return any data");
  }

  /// @notice Reverts if pubdata usage is too high => "NotEnoughGasForPubdata()"
  function test__callWithExactGasSafeReturnData_RevertWhen_NotEnoughGasForPubdata() public {
    // Simulate pubdata usage:
    // Set the initial pubdata value to 1000, then (via a mock) the after-call value to 5000.
    s_mockSystemContext.setCurrentPubdataSpent(1000);
    vm.mockCall(
      address(s_mockSystemContext),
      abi.encodeWithSelector(s_mockSystemContext.getCurrentPubdataSpent.selector),
      abi.encode(5000)
    );

    // This difference = 4000 pubdata * 10 gas/byte = 40,000 extra gas needed.
    // We'll provide allowed gas = 200,000, and _maxTotalGas = 200k.
    // With the overhead, the check should fail, triggering NotEnoughGasForPubdata.
    vm.expectRevert(NotEnoughGasForPubdata.selector);

    (bool successCall, ) = _limitedGasCallWithExactGas(
      400_000, // subcall gas
      address(s_target),
      200_000, // _maxTotalGas exactly 200k
      abi.encodeWithSelector(TestTarget.returnData.selector),
      100
    );
    assertFalse(successCall, "Subcall itself must revert");
  }

  /// @notice Succeeds under normal conditions, returning data.
  function test__callWithExactGasSafeReturnData_Success() public {
    (bool successCall, bytes memory retData) = _limitedGasCallWithExactGas(
      5_000_000,
      address(s_target),
      4_000_000,
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

  /// @notice Truncates return data if it exceeds _maxReturnBytes.
  function test__callWithExactGasSafeReturnData_TruncatesData() public {
    (bool successCall, bytes memory retData) = _limitedGasCallWithExactGas(
      500_000,
      address(s_target),
      300_000,
      abi.encodeWithSelector(TestTarget.returnLargeData.selector),
      50 // only allow 50 bytes of return data
    );

    assertTrue(successCall, "Subcall must not revert");
    (bool success, bytes memory returnedData, ) = _decodeResult(retData);
    assertTrue(success, "Target call must succeed");
    assertEq(returnedData.length, 50, "Should have truncated the large data to 50 bytes");
  }

  /// @notice Reverts with a revert reason when the target reverts with reason.
  function test__callWithExactGasSafeReturnData_RevertWhen_TargetRevertsWithReason() public {
    // Expect the revert reason "CustomRevertReason"
    vm.expectRevert(bytes("CustomRevertReason"));

    (bool successCall, ) = _limitedGasCallWithExactGas(
      1_000_000,
      address(s_target),
      1_000_000,
      abi.encodeWithSelector(TestTarget.revertWithReason.selector),
      100
    );
    assertFalse(successCall, "Subcall itself must revert");
  }

  /// @notice Reverts if the target reverts without a reason.
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
