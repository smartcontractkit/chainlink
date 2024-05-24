// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {FunctionsRouter} from "../../dev/v1_X/FunctionsRouter.sol";
import {FunctionsSubscriptions} from "../../dev/v1_X/FunctionsSubscriptions.sol";
import {FunctionsRequest} from "../../dev/v1_X/libraries/FunctionsRequest.sol";
import {FunctionsResponse} from "../../dev/v1_X/libraries/FunctionsResponse.sol";
import {FunctionsClientTestHelper} from "./testhelpers/FunctionsClientTestHelper.sol";

import {FunctionsRoutesSetup, FunctionsOwnerAcceptTermsOfServiceSetup, FunctionsSubscriptionSetup, FunctionsClientRequestSetup} from "./Setup.t.sol";

import "forge-std/Vm.sol";

/// @notice #acceptTermsOfService
contract Gas_AcceptTermsOfService is FunctionsRoutesSetup {
  bytes32 s_sigR;
  bytes32 s_sigS;
  uint8 s_sigV;

  function setUp() public virtual override {
    vm.pauseGasMetering();

    FunctionsRoutesSetup.setUp();

    bytes32 message = s_termsOfServiceAllowList.getMessage(OWNER_ADDRESS, OWNER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (s_sigV, s_sigR, s_sigS) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);
  }

  function test_AcceptTermsOfService_Gas() public {
    // Pull storage variables into memory
    address ownerAddress = OWNER_ADDRESS;
    bytes32 sigR = s_sigR;
    bytes32 sigS = s_sigS;
    uint8 sigV = s_sigV;
    vm.resumeGasMetering();

    s_termsOfServiceAllowList.acceptTermsOfService(ownerAddress, ownerAddress, sigR, sigS, sigV);
  }
}

/// @notice #createSubscription
contract Gas_CreateSubscription is FunctionsOwnerAcceptTermsOfServiceSetup {
  function test_CreateSubscription_Gas() public {
    s_functionsRouter.createSubscription();
  }
}

/// @notice #addConsumer
contract Gas_AddConsumer is FunctionsSubscriptionSetup {
  function setUp() public virtual override {
    vm.pauseGasMetering();

    FunctionsSubscriptionSetup.setUp();
  }

  function test_AddConsumer_Gas() public {
    // Keep input data in memory
    uint64 subscriptionId = s_subscriptionId;
    address consumerAddress = address(s_functionsCoordinator); // use garbage address
    vm.resumeGasMetering();

    s_functionsRouter.addConsumer(subscriptionId, consumerAddress);
  }
}

/// @notice #fundSubscription
contract Gas_FundSubscription is FunctionsSubscriptionSetup {
  function setUp() public virtual override {
    vm.pauseGasMetering();

    FunctionsSubscriptionSetup.setUp();
  }

  function test_FundSubscription_Gas() public {
    // Keep input data in memory
    address routerAddress = address(s_functionsRouter);
    uint96 s_subscriptionFunding = 10 * JUELS_PER_LINK; // 10 LINK
    bytes memory data = abi.encode(s_subscriptionId);
    vm.resumeGasMetering();

    s_linkToken.transferAndCall(routerAddress, s_subscriptionFunding, data);
  }
}

/// @notice #sendRequest
contract Gas_SendRequest is FunctionsSubscriptionSetup {
  bytes s_minimalRequestData;
  bytes s_maximalRequestData;

  function _makeStringOfBytesSize(uint16 bytesSize) internal pure returns (string memory) {
    return vm.toString(new bytes((bytesSize - 2) / 2));
  }

  function setUp() public virtual override {
    vm.pauseGasMetering();

    FunctionsSubscriptionSetup.setUp();

    {
      // Create minimum viable request data
      FunctionsRequest.Request memory minimalRequest;
      string memory minimalSourceCode = "return Functions.encodeString('hello world');";
      FunctionsRequest._initializeRequest(
        minimalRequest,
        FunctionsRequest.Location.Inline,
        FunctionsRequest.CodeLanguage.JavaScript,
        minimalSourceCode
      );
      s_minimalRequestData = FunctionsRequest._encodeCBOR(minimalRequest);
    }

    {
      // Create maximum viable request data - 30 KB encoded data
      FunctionsRequest.Request memory maxmimalRequest;

      // Create maximum viable request data - 30 KB encoded data
      string memory maximalSourceCode = _makeStringOfBytesSize(29_898); // CBOR size without source code is 102 bytes
      FunctionsRequest._initializeRequest(
        maxmimalRequest,
        FunctionsRequest.Location.Inline,
        FunctionsRequest.CodeLanguage.JavaScript,
        maximalSourceCode
      );
      s_maximalRequestData = FunctionsRequest._encodeCBOR(maxmimalRequest);
      assertEq(s_maximalRequestData.length, 30_000);
    }
  }

  /// @dev The order of these test cases matters as the first test will consume more gas by writing over default values
  function test_SendRequest_MaximumGas() public {
    // Pull storage variables into memory
    bytes memory maximalRequestData = s_maximalRequestData;
    uint64 subscriptionId = s_subscriptionId;
    uint32 callbackGasLimit = 300_000;
    bytes32 donId = s_donId;
    vm.resumeGasMetering();

    s_functionsClient.sendRequestBytes(maximalRequestData, subscriptionId, callbackGasLimit, donId);
  }

  function test_SendRequest_MinimumGas() public {
    // Pull storage variables into memory
    bytes memory minimalRequestData = s_minimalRequestData;
    uint64 subscriptionId = s_subscriptionId;
    uint32 callbackGasLimit = 5_500;
    bytes32 donId = s_donId;
    vm.resumeGasMetering();

    s_functionsClient.sendRequestBytes(minimalRequestData, subscriptionId, callbackGasLimit, donId);
  }
}

// Setup Fulfill Gas tests
contract Gas_FulfillRequest_Setup is FunctionsClientRequestSetup {
  mapping(uint256 reportNumber => Report) s_reports;

  FunctionsClientTestHelper s_functionsClientWithMaximumReturnData;

  function _makeStringOfBytesSize(uint16 bytesSize) internal pure returns (string memory) {
    return vm.toString(new bytes((bytesSize - 2) / 2));
  }

  function setUp() public virtual override {
    FunctionsSubscriptionSetup.setUp();

    {
      // Deploy consumer that has large revert return data
      s_functionsClientWithMaximumReturnData = new FunctionsClientTestHelper(address(s_functionsRouter));
      s_functionsClientWithMaximumReturnData.setRevertFulfillRequest(true);
      string memory revertMessage = _makeStringOfBytesSize(30_000); // 30kb - FunctionsRouter cuts off response at MAX_CALLBACK_RETURN_BYTES = 4 + 4 * 32 = 132bytes, go well above that
      s_functionsClientWithMaximumReturnData.setRevertFulfillRequestMessage(revertMessage);
      s_functionsRouter.addConsumer(s_subscriptionId, address(s_functionsClientWithMaximumReturnData));
    }

    // Set up maximum gas test
    {
      // Send request #2 for maximum gas test
      uint8 requestNumber = 2;

      bytes memory secrets = new bytes(0);
      string[] memory args = new string[](0);
      bytes[] memory bytesArgs = new bytes[](0);
      uint32 callbackGasLimit = 300_000;

      // Create maximum viable request data - 30 KB encoded data
      string memory maximalSourceCode = _makeStringOfBytesSize(29_898); // CBOR size without source code is 102 bytes

      _sendAndStoreRequest(
        requestNumber,
        maximalSourceCode,
        secrets,
        args,
        bytesArgs,
        callbackGasLimit,
        address(s_functionsClientWithMaximumReturnData)
      );

      // Build the report transmission data
      uint256[] memory requestNumberKeys = new uint256[](1);
      requestNumberKeys[0] = requestNumber;
      string[] memory results = new string[](1);
      // Build a 256 byte response size
      results[0] = _makeStringOfBytesSize(256);
      bytes[] memory errors = new bytes[](1);
      errors[0] = new bytes(0); // No error

      (bytes memory report, bytes32[3] memory reportContext) = _buildReport(requestNumberKeys, results, errors);

      uint256[] memory signerPrivateKeys = new uint256[](3);
      signerPrivateKeys[0] = NOP_SIGNER_PRIVATE_KEY_1;
      signerPrivateKeys[1] = NOP_SIGNER_PRIVATE_KEY_2;
      signerPrivateKeys[2] = NOP_SIGNER_PRIVATE_KEY_3;

      (bytes32[] memory rawRs, bytes32[] memory rawSs, bytes32 rawVs) = _signReport(
        report,
        reportContext,
        signerPrivateKeys
      );

      // Store the report data
      s_reports[1] = Report({rs: rawRs, ss: rawSs, vs: rawVs, report: report, reportContext: reportContext});
    }

    // Set up minimum gas test
    {
      // Send requests minimum gas test
      uint8 requestsToSend = 1;
      uint8 requestNumberOffset = 3; // the setup already has request #1 sent, and the previous test case uses request #2, start from request #3

      string memory sourceCode = "return Functions.encodeString('hello world');";
      bytes memory secrets = new bytes(0);
      string[] memory args = new string[](0);
      bytes[] memory bytesArgs = new bytes[](0);
      uint32 callbackGasLimit = 5_500;

      for (uint256 i = 0; i < requestsToSend; ++i) {
        _sendAndStoreRequest(i + requestNumberOffset, sourceCode, secrets, args, bytesArgs, callbackGasLimit);
      }

      // Build the report transmission data
      uint256[] memory requestNumberKeys = new uint256[](requestsToSend);
      string[] memory results = new string[](requestsToSend);
      bytes[] memory errors = new bytes[](requestsToSend);
      for (uint256 i = 0; i < requestsToSend; ++i) {
        requestNumberKeys[i] = i + requestNumberOffset;
        results[i] = "hello world";
        errors[i] = new bytes(0); // no error
      }

      (bytes memory report, bytes32[3] memory reportContext) = _buildReport(requestNumberKeys, results, errors);

      uint256[] memory signerPrivateKeys = new uint256[](3);
      signerPrivateKeys[0] = NOP_SIGNER_PRIVATE_KEY_1;
      signerPrivateKeys[1] = NOP_SIGNER_PRIVATE_KEY_2;
      signerPrivateKeys[2] = NOP_SIGNER_PRIVATE_KEY_3;

      (bytes32[] memory rawRs, bytes32[] memory rawSs, bytes32 rawVs) = _signReport(
        report,
        reportContext,
        signerPrivateKeys
      );

      // Store the report data
      s_reports[2] = Report({rs: rawRs, ss: rawSs, vs: rawVs, report: report, reportContext: reportContext});
    }

    vm.stopPrank();
    vm.startPrank(NOP_TRANSMITTER_ADDRESS_1);
  }
}

/// @notice #fulfillRequest
contract Gas_FulfillRequest_Success is Gas_FulfillRequest_Setup {
  function setUp() public virtual override {
    vm.pauseGasMetering();

    Gas_FulfillRequest_Setup.setUp();
  }

  /// @dev The order of these test cases matters as the first test will consume more gas by writing over default values
  function test_FulfillRequest_Success_MaximumGas() public {
    // Pull storage variables into memory
    uint8 reportNumber = 1;
    bytes32[] memory rs = s_reports[reportNumber].rs;
    bytes32[] memory ss = s_reports[reportNumber].ss;
    bytes32 vs = s_reports[reportNumber].vs;
    bytes memory report = s_reports[reportNumber].report;
    bytes32[3] memory reportContext = s_reports[reportNumber].reportContext;
    vm.resumeGasMetering();

    // 1 fulfillment in the report, single request takes on all report validation cost
    // maximum request
    // maximum NOPs
    // maximum return data
    // first storage write to change default values
    s_functionsCoordinator.transmit(reportContext, report, rs, ss, vs);
  }

  function test_FulfillRequest_Success_MinimumGas() public {
    // Pull storage variables into memory
    uint8 reportNumber = 2;
    bytes32[] memory rs = s_reports[reportNumber].rs;
    bytes32[] memory ss = s_reports[reportNumber].ss;
    bytes32 vs = s_reports[reportNumber].vs;
    bytes memory report = s_reports[reportNumber].report;
    bytes32[3] memory reportContext = s_reports[reportNumber].reportContext;
    vm.resumeGasMetering();

    // max fulfillments in the report, cost of validation split between all
    // minimal request
    // minimum NOPs
    // no return data
    // not storage writing default values
    s_functionsCoordinator.transmit(reportContext, report, rs, ss, vs);
  }
}

/// @notice #fulfillRequest
contract Gas_FulfillRequest_DuplicateRequestID is Gas_FulfillRequest_Setup {
  function setUp() public virtual override {
    vm.pauseGasMetering();

    // Send requests
    Gas_FulfillRequest_Setup.setUp();
    // Fulfill request #1 & #2
    for (uint256 i = 1; i < 3; i++) {
      uint256 reportNumber = i;
      bytes32[] memory rs = s_reports[reportNumber].rs;
      bytes32[] memory ss = s_reports[reportNumber].ss;
      bytes32 vs = s_reports[reportNumber].vs;
      bytes memory report = s_reports[reportNumber].report;
      bytes32[3] memory reportContext = s_reports[reportNumber].reportContext;
      s_functionsCoordinator.transmit(reportContext, report, rs, ss, vs);
    }

    // Now tests will attempt to transmit reports with respones to requests that have already been fulfilled
  }

  /// @dev The order of these test cases matters as the first test will consume more gas by writing over default values
  function test_FulfillRequest_DuplicateRequestID_MaximumGas() public {
    // Pull storage variables into memory
    uint8 reportNumber = 1;
    bytes32[] memory rs = s_reports[reportNumber].rs;
    bytes32[] memory ss = s_reports[reportNumber].ss;
    bytes32 vs = s_reports[reportNumber].vs;
    bytes memory report = s_reports[reportNumber].report;
    bytes32[3] memory reportContext = s_reports[reportNumber].reportContext;
    vm.resumeGasMetering();

    // 1 fulfillment in the report, single request takes on all report validation cost
    // maximum request
    // maximum NOPs
    // maximum return data
    // first storage write to change default values
    s_functionsCoordinator.transmit(reportContext, report, rs, ss, vs);
  }

  function test_FulfillRequest_DuplicateRequestID_MinimumGas() public {
    // Pull storage variables into memory
    uint8 reportNumber = 2;
    bytes32[] memory rs = s_reports[reportNumber].rs;
    bytes32[] memory ss = s_reports[reportNumber].ss;
    bytes32 vs = s_reports[reportNumber].vs;
    bytes memory report = s_reports[reportNumber].report;
    bytes32[3] memory reportContext = s_reports[reportNumber].reportContext;
    vm.resumeGasMetering();

    // max fulfillments in the report, cost of validation split between all
    // minimal request
    // minimum NOPs
    // no return data
    // not storage writing default values
    s_functionsCoordinator.transmit(reportContext, report, rs, ss, vs);
  }
}
