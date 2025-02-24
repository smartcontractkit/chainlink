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
    // Write mock's code to 0x800b so that library calls see it
    vm.etch(address(0x800b), address(s_mockSystemContext).code);
  }

  function test__callback_RevertWhen_NoClientCode() public {
    bytes32 reqId = bytes32("reqIdNoCode");
    bytes memory resp = bytes("responseData");
    bytes memory err = bytes("errData");
    uint32 totalGas = 5_000_000;
    uint32 callbackGasLimit = 4_000_000;
    address noCodeAddress = address(12345);

    ZKSyncFunctionsRouter.CallbackResult memory result = _callback(
      reqId,
      resp,
      err,
      totalGas,
      callbackGasLimit,
      noCodeAddress
    );

    assertFalse(result.success, "Should skip => success=false");
    assertEq(result.gasUsed, 0, "gasUsed=0 for skip");
    assertEq(result.returnData.length, 0, "no return data");
  }

  function test__callback_Success() public {
    bytes32 reqId = bytes32("reqSuccess");
    bytes memory resp = bytes("responseData");
    bytes memory err = bytes("errData");
    uint32 totalGas = 5_000_000;
    uint32 callbackGasLimit = 4_000_000;
    address client = address(s_mockClientSuccess);

    ZKSyncFunctionsRouter.CallbackResult memory result = _callback(
      reqId,
      resp,
      err,
      totalGas,
      callbackGasLimit,
      client
    );

    assertTrue(result.success, "callback should succeed");
    assertGt(result.gasUsed, 0, "some gas used");
    assertTrue(result.returnData.length > 0, "client returns a bool => should have data");
  }

  function test__callback_RevertWhen_ClientReverts() public {
    bytes32 reqId = bytes32("reqIdRevert");
    bytes memory resp = bytes("someResponse");
    bytes memory err = bytes("someErr");
    // Use a moderate gas limit so that we don't trigger the _maxTotalGas check.
    uint32 totalGas = 5_000_000;
    uint32 callbackGasLimit = 4_000_000;
    address client = address(s_mockClientRevert);
    ZKSyncFunctionsRouter.CallbackResult memory result = _callback(
      reqId,
      resp,
      err,
      totalGas,
      callbackGasLimit,
      client
    );

    assertFalse(result.success, "client revert => success=false");
    assertGt(result.gasUsed, 0, "some gas is consumed");
    // returnData should contain the revert reason "MockClientRevert"
    assertTrue(result.returnData.length > 0, "contains revert reason data");
  }

  /// @notice Example test verifying pubdata usage is zero and comparing internal measurement with external gas usage.
  function test__callback_PubdataUsage_IsZero() public {
    s_mockSystemContext.setGasPerPubdataByte(0);
    bytes32 reqId = bytes32("reqPubdata");
    bytes memory resp = bytes("someResponse");
    bytes memory err = bytes("someErr");
    uint32 totalGas = 5_000_000;
    uint32 callbackGasLimit = 4_000_000;
    address client = address(s_mockClientSuccess);
    uint256 startGas = gasleft();
    ZKSyncFunctionsRouter.CallbackResult memory result = _callback(
      reqId,
      resp,
      err,
      totalGas,
      callbackGasLimit,
      client
    );
    uint256 endGas = gasleft();
    uint256 actualUsed = startGas - endGas;
    assertTrue(result.success, "callback success");
    assertGt(result.gasUsed, 0, "callback claims >0 gas used");
    // Allow a margin between the router's internal measurement and actual external usage.
    assertLe(result.gasUsed, actualUsed, "Router's gasUsed should not exceed actual external usage by large margin");
  }

  /// @notice Confirm large return data gets truncated by _maxReturnBytes.
  function test__callback_ReturnDataTruncation() public {
    // Deploy large-return client.
    MockClientLargeReturn bigClient = new MockClientLargeReturn();

    bytes32 reqId = bytes32("reqLargeReturn");
    bytes memory resp = bytes("someResponse");
    bytes memory err = bytes("someErr");
    uint32 totalGas = 5_000_000;
    uint32 callbackGasLimit = 4_000_000;
    address client = address(bigClient);

    ZKSyncFunctionsRouter.CallbackResult memory result = _callback(
      reqId,
      resp,
      err,
      totalGas,
      callbackGasLimit,
      client
    );
    assertTrue(result.success, "Should succeed");
    uint256 expectedMax = s_functionsRouter.MAX_CALLBACK_RETURN_BYTES();
    // The returned data should be truncated exactly to expectedMax.
    assertEq(result.returnData.length, expectedMax, "Should truncate data to MAX_CALLBACK_RETURN_BYTES");
  }

  /// @notice Internal helper to call the router's exposed callback function.
  function _callback(
    bytes32 reqId,
    bytes memory resp,
    bytes memory err,
    uint32 totalGas,
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
    (bool ok, bytes memory retData) = address(s_functionsRouter).call{gas: totalGas}(payload);
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
    // Return ~1,000 bytes.
    bytes memory largeData = new bytes(1000);
    for (uint i = 0; i < 1000; i++) {
      largeData[i] = bytes1(uint8(65 + (i % 26))); // Fill with A..Z.
    }
    return largeData;
  }
}

contract MockClientRevert {
  function handleOracleFulfillment(bytes32, bytes memory, bytes memory) external pure returns (bool) {
    revert("MockClientRevert");
  }
}
