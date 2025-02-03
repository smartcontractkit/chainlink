// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {ZKSyncFunctionsRouter} from "../../v1_3_0_zksync/ZKSyncFunctionsRouter.sol";
import {FunctionsRouter} from "../../v1_0_0/FunctionsRouter.sol";
import {ZKSyncFunctionsRouterHarness} from "./testhelpers/ZKSyncFunctionsRouterHarness.sol";
import {ZKSyncFunctionsRouterSetup} from "./Setup.t.sol";
import {MockSystemContext} from "../../../shared/test/mocks/MockSystemContext.sol";

contract ZKSyncFunctionsRouter__Callback is ZKSyncFunctionsRouterSetup {
  MockClientSuccess internal s_mockClientSuccess;
  MockClientRevert internal s_mockClientRevert;
  MockSystemContext internal s_mockSystemContext;

  struct CallbackResult {
    bool success;
    uint256 gasUsed;
    bytes returnData;
  }

  function setUp() public virtual override {
    super.setUp();
    s_mockClientSuccess = new MockClientSuccess();
    s_mockClientRevert = new MockClientRevert();

    s_mockSystemContext = new MockSystemContext();
    // Write mock's code to 0x800b so library calls see it
    vm.etch(address(0x800b), address(s_mockSystemContext).code);
  }

  function test__callback_RevertWhen_NoClientCode() public {
    bytes32 reqId = bytes32("reqIdNoCode");
    bytes memory resp = bytes("responseData");
    bytes memory err = bytes("errData");
    uint32 callbackGasLimit = 5_000_000;
    // Create an address with no code
    address noCodeAddress = address(12345);

    ZKSyncFunctionsRouter.CallbackResult memory result = _callback(reqId, resp, err, callbackGasLimit, noCodeAddress);

    assertEq(result.success, false, "Should skip => success=false");
    assertEq(result.gasUsed, 0, "gasUsed=0 for skip");
    assertEq(result.returnData.length, 0, "no return data");
  }

  function test__callback_Success() public {
    bytes32 reqId = bytes32("reqSuccess");
    bytes memory resp = bytes("responseData");
    bytes memory err = bytes("errData");
    uint32 callbackGasLimit = 5_000_000;
    address client = address(s_mockClientSuccess);

    ZKSyncFunctionsRouter.CallbackResult memory result = _callback(reqId, resp, err, callbackGasLimit, client);

    assertEq(result.success, true, "callback should succeed");
    assertGt(result.gasUsed, 0, "some gas used");
    assertTrue(result.returnData.length > 0, "client returns a bool => should have data");
  }

  function test__callback_RevertWhen_ClientReverts() public {
    bytes32 reqId = bytes32("reqIdRevert");
    bytes memory resp = bytes("someResponse");
    bytes memory err = bytes("someErr");
    uint32 callbackGasLimit = 500_000;
    address client = address(s_mockClientRevert);
    ZKSyncFunctionsRouter.CallbackResult memory result = _callback(reqId, resp, err, callbackGasLimit, client);

    assertEq(result.success, false, "client revert => success=false");
    assertGt(result.gasUsed, 0, "some gas is consumed");
    // returnData should contain the revert reason "MockClientRevert"
    assertTrue(result.returnData.length > 0, "contains revert reason data");
  }

  /// @notice Example test verifying pubdata usage is zero and how to measure actual gas usage externally
  function test__callback_PubdataUsage_IsZero() public {
    bytes32 reqId = bytes32("reqPubdata");
    bytes memory resp = bytes("someResponse");
    bytes memory err = bytes("someErr");
    uint32 callbackGasLimit = 500_000;
    address client = address(s_mockClientSuccess);
    uint256 startGas = gasleft();
    ZKSyncFunctionsRouter.CallbackResult memory result = _callback(reqId, resp, err, callbackGasLimit, client);
    uint256 endGas = gasleft();
    uint256 actualUsed = startGas - endGas;
    assertTrue(result.success, "callback success");
    assertGt(result.gasUsed, 0, "callback claims >0 gas used");
    // Typically, `actualUsed` > `result.gasUsed` because:
    //   - There's overhead using `_callback` helper, function calls, memory expansions, etc.
    // Or in the opposite scenario, you may see them close but not identical. That is expected.
    // The important part is that if pubdataSpent=0, then 'result.gasUsed' won't have any extra overhead added.
    // We just ensure the router's number isn't bigger than the actual external usage by too much:
    assertLe(
      result.gasUsed,
      actualUsed + 5000,
      "Router's gasUsed should not exceed actual external usage by large margin"
    );
  }

  /// @notice Test pubdata usage _callback call by re-entrant logic in MockSystemContext
  function test__callback_PubdataUsage() public {
    // We'll rely on the re-entrant logic so the "before" read returns 1000, the "after" read returns 3000.
    // difference = 2000 pubdata bytes
    // total cost = 2000 * 5 (GasPerPubdataByte) = 10,000
    s_mockSystemContext.setGasPerPubdataByte(5);

    // Now do a single _callback call
    bytes32 reqId = bytes32("singleCallPubdata");
    bytes memory resp = bytes("response data");
    bytes memory err = bytes("error data");
    uint32 callbackGasLimit = 500_000;
    address client = address(s_mockClientSuccess);

    // The library will read getCurrentPubdataSpent() "before" => 1000
    // Then the mock sets storedPubdataSpent=3000 for the "after" read => difference=2000
    // => 2000 * 5 = 10,000 additional cost.

    FunctionsRouter.CallbackResult memory result = _callback(reqId, resp, err, callbackGasLimit, client);
    assertTrue(result.success, "Callback success expected");

    // We expect at least 10k more than an identical call if pubdata difference=0
    assertGt(result.gasUsed, 10_000, "Should have included ~10k pubdata cost");
  }

  /// @notice Confirm large return data gets truncated by `_maxReturnBytes`.
  function test__callback_ReturnDataTruncation() public {
    // Deploy large-return client
    MockClientLargeReturn bigClient = new MockClientLargeReturn();

    // typical callback
    bytes32 reqId = bytes32("reqLargeReturn");
    bytes memory resp = bytes("someResponse");
    bytes memory err = bytes("someErr");
    uint32 callbackGasLimit = 1_000_000;
    address client = address(bigClient);

    ZKSyncFunctionsRouter.CallbackResult memory result = _callback(reqId, resp, err, callbackGasLimit, client);
    assertTrue(result.success, "Should succeed");

    // If the library sets MAX_CALLBACK_RETURN_BYTES, e.g. 256 or 100, check truncation
    // E.g. if it's 50 => we expect only the first 50 bytes.
    uint256 expectedMax = s_functionsRouter.MAX_CALLBACK_RETURN_BYTES();
    // The result's length cannot exceed expectedMax
    assertEq(result.returnData.length, expectedMax, "Should truncate data to MAX_CALLBACK_RETURN_BYTES");
  }

  function _callback(
    bytes32 reqId,
    bytes memory resp,
    bytes memory err,
    uint32 callbackGasLimit,
    address client
  ) internal returns (FunctionsRouter.CallbackResult memory) {
    bytes memory payload = abi.encodeWithSelector(
      s_functionsRouter.exposed_callback.selector,
      reqId,
      resp,
      err,
      callbackGasLimit,
      client
    );

    (bool ok, bytes memory retData) = address(s_functionsRouter).call{gas: callbackGasLimit}(payload);
    assertTrue(ok, "callback should succeed");
    return abi.decode(retData, (FunctionsRouter.CallbackResult));
  }
}

contract MockClientSuccess {
  function handleOracleFulfillment(bytes32, bytes memory, bytes memory) external pure returns (bool) {
    return true;
  }
}

contract MockClientLargeReturn {
  function handleOracleFulfillment(bytes32, bytes memory, bytes memory) external pure returns (bytes memory) {
    // Return ~1,000 bytes
    bytes memory largeData = new bytes(1000);
    for (uint i = 0; i < 1000; i++) {
      largeData[i] = bytes1(uint8(65 + (i % 26))); // fill with A..Z
    }
    return largeData;
  }
}

contract MockClientRevert {
  function handleOracleFulfillment(bytes32, bytes memory, bytes memory) external pure returns (bool) {
    revert("MockClientRevert");
  }
}
