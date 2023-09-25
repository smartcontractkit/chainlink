// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {FunctionsRouter} from "../../dev/v1_0_0/FunctionsRouter.sol";
import {FunctionsCoordinatorTestHelper} from "./testhelpers/FunctionsCoordinatorTestHelper.sol";
import {FunctionsBilling} from "../../dev/v1_0_0/FunctionsBilling.sol";
import {FunctionsResponse} from "../../dev/v1_0_0/libraries/FunctionsResponse.sol";
import {MockV3Aggregator} from "../../../tests/MockV3Aggregator.sol";
import {TermsOfServiceAllowList} from "../../dev/v1_0_0/accessControl/TermsOfServiceAllowList.sol";
import {FunctionsClientUpgradeHelper} from "./testhelpers/FunctionsClientUpgradeHelper.sol";
import {MockLinkToken} from "../../../mocks/MockLinkToken.sol";

import "forge-std/Vm.sol";

/// @notice Set up to deploy the following contracts: FunctionsRouter, FunctionsCoordinator, LINK/ETH Feed, ToS Allow List, and LINK token
contract FunctionsRouterSetup is BaseTest {
  FunctionsRouter internal s_functionsRouter;
  FunctionsCoordinatorTestHelper internal s_functionsCoordinator; // TODO: use actual FunctionsCoordinator instead of helper
  MockV3Aggregator internal s_linkEthFeed;
  TermsOfServiceAllowList internal s_termsOfServiceAllowList;
  MockLinkToken internal s_linkToken;

  uint16 internal s_maxConsumersPerSubscription = 3;
  uint72 internal s_adminFee = 100;
  uint72 internal s_donFee = 100;
  bytes4 internal s_handleOracleFulfillmentSelector = 0x0ca76175;
  uint16 s_subscriptionDepositMinimumRequests = 1;
  uint72 s_subscriptionDepositJuels = 11 * JUELS_PER_LINK;

  int256 internal LINK_ETH_RATE = 6000000000000000;

  uint256 internal TOS_SIGNER_PRIVATE_KEY = 0x3;
  address internal TOS_SIGNER = vm.addr(TOS_SIGNER_PRIVATE_KEY);

  function setUp() public virtual override {
    BaseTest.setUp();
    s_linkToken = new MockLinkToken();
    s_functionsRouter = new FunctionsRouter(address(s_linkToken), getRouterConfig());
    s_linkEthFeed = new MockV3Aggregator(0, LINK_ETH_RATE);
    s_functionsCoordinator = new FunctionsCoordinatorTestHelper(
      address(s_functionsRouter),
      getCoordinatorConfig(),
      address(s_linkEthFeed)
    );
    s_termsOfServiceAllowList = new TermsOfServiceAllowList(getTermsOfServiceConfig());
  }

  function getRouterConfig() public view returns (FunctionsRouter.Config memory) {
    uint32[] memory maxCallbackGasLimits = new uint32[](3);
    maxCallbackGasLimits[0] = 300_000;
    maxCallbackGasLimits[1] = 500_000;
    maxCallbackGasLimits[2] = 1_000_000;

    return
      FunctionsRouter.Config({
        maxConsumersPerSubscription: s_maxConsumersPerSubscription,
        adminFee: s_adminFee,
        handleOracleFulfillmentSelector: s_handleOracleFulfillmentSelector,
        maxCallbackGasLimits: maxCallbackGasLimits,
        gasForCallExactCheck: 5000,
        subscriptionDepositMinimumRequests: s_subscriptionDepositMinimumRequests,
        subscriptionDepositJuels: s_subscriptionDepositJuels
      });
  }

  function getCoordinatorConfig() public view returns (FunctionsBilling.Config memory) {
    return
      FunctionsBilling.Config({
        feedStalenessSeconds: 24 * 60 * 60, // 1 day
        gasOverheadAfterCallback: 50_000, // TODO: update
        gasOverheadBeforeCallback: 100_00, // TODO: update
        requestTimeoutSeconds: 60 * 5, // 5 minutes
        donFee: s_donFee,
        maxSupportedRequestDataVersion: 1,
        fulfillmentGasPriceOverEstimationBP: 5000,
        fallbackNativePerUnitLink: 5000000000000000
      });
  }

  function getTermsOfServiceConfig() public view returns (TermsOfServiceAllowList.Config memory) {
    return TermsOfServiceAllowList.Config({enabled: true, signerPublicKey: TOS_SIGNER});
  }
}

/// @notice Set up to set the OCR configuration of the Coordinator contract
contract FunctionsDONSetup is FunctionsRouterSetup {
  uint256 internal NOP_SIGNER_PRIVATE_KEY_1 = 0x100;
  address internal NOP_SIGNER_ADDRESS_1 = vm.addr(NOP_SIGNER_PRIVATE_KEY_1);
  uint256 internal NOP_SIGNER_PRIVATE_KEY_2 = 0x101;
  address internal NOP_SIGNER_ADDRESS_2 = vm.addr(NOP_SIGNER_PRIVATE_KEY_2);
  uint256 internal NOP_SIGNER_PRIVATE_KEY_3 = 0x102;
  address internal NOP_SIGNER_ADDRESS_3 = vm.addr(NOP_SIGNER_PRIVATE_KEY_3);
  uint256 internal NOP_SIGNER_PRIVATE_KEY_4 = 0x103;
  address internal NOP_SIGNER_ADDRESS_4 = vm.addr(NOP_SIGNER_PRIVATE_KEY_4);

  uint256 internal NOP_TRANSMITTER_PRIVATE_KEY_1 = 0x104;
  address internal NOP_TRANSMITTER_ADDRESS_1 = vm.addr(NOP_TRANSMITTER_PRIVATE_KEY_1);
  uint256 internal NOP_TRANSMITTER_PRIVATE_KEY_2 = 0x105;
  address internal NOP_TRANSMITTER_ADDRESS_2 = vm.addr(NOP_TRANSMITTER_PRIVATE_KEY_2);
  uint256 internal NOP_TRANSMITTER_PRIVATE_KEY_3 = 0x106;
  address internal NOP_TRANSMITTER_ADDRESS_3 = vm.addr(NOP_TRANSMITTER_PRIVATE_KEY_3);
  uint256 internal NOP_TRANSMITTER_PRIVATE_KEY_4 = 0x107;
  address internal NOP_TRANSMITTER_ADDRESS_4 = vm.addr(NOP_TRANSMITTER_PRIVATE_KEY_4);

  address[] internal s_signers;
  address[] internal s_transmitters;
  uint8 s_f = 1;
  bytes internal s_onchainConfig = new bytes(0);
  uint64 internal s_offchainConfigVersion = 1;
  bytes internal s_offchainConfig = new bytes(0);

  function setUp() public virtual override {
    FunctionsRouterSetup.setUp();

    s_signers = new address[](4);
    s_signers[0] = NOP_SIGNER_ADDRESS_1;
    s_signers[1] = NOP_SIGNER_ADDRESS_2;
    s_signers[2] = NOP_SIGNER_ADDRESS_3;
    s_signers[3] = NOP_SIGNER_ADDRESS_4;

    s_transmitters = new address[](4);
    s_transmitters[0] = NOP_TRANSMITTER_ADDRESS_1;
    s_transmitters[1] = NOP_TRANSMITTER_ADDRESS_2;
    s_transmitters[2] = NOP_TRANSMITTER_ADDRESS_3;
    s_transmitters[3] = NOP_TRANSMITTER_ADDRESS_4;

    // set OCR config
    s_functionsCoordinator.setConfig(
      s_signers,
      s_transmitters,
      s_f,
      s_onchainConfig,
      s_offchainConfigVersion,
      s_offchainConfig
    );
  }
}

/// @notice Set up to add the Coordinator and ToS Allow Contract as routes on the Router contract
contract FunctionsRoutesSetup is FunctionsDONSetup {
  bytes32 s_donId = bytes32("1");

  function setUp() public virtual override {
    FunctionsDONSetup.setUp();

    bytes32 allowListId = s_functionsRouter.getAllowListId();
    bytes32[] memory proposedContractSetIds = new bytes32[](2);
    proposedContractSetIds[0] = s_donId;
    proposedContractSetIds[1] = allowListId;
    address[] memory proposedContractSetAddresses = new address[](2);
    proposedContractSetAddresses[0] = address(s_functionsCoordinator);
    proposedContractSetAddresses[1] = address(s_termsOfServiceAllowList);

    s_functionsRouter.proposeContractsUpdate(proposedContractSetIds, proposedContractSetAddresses);
    s_functionsRouter.updateContracts();
  }
}

/// @notice Set up for the OWNER_ADDRESS to accept the Terms of Service
contract FunctionsOwnerAcceptTermsOfServiceSetup is FunctionsRoutesSetup {
  function setUp() public virtual override {
    FunctionsRoutesSetup.setUp();

    bytes32 message = s_termsOfServiceAllowList.getMessage(OWNER_ADDRESS, OWNER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);
    s_termsOfServiceAllowList.acceptTermsOfService(OWNER_ADDRESS, OWNER_ADDRESS, r, s, v);
  }
}

/// @notice Set up to deploy a consumer contract
contract FunctionsClientSetup is FunctionsOwnerAcceptTermsOfServiceSetup {
  FunctionsClientUpgradeHelper internal s_functionsClient;

  function setUp() public virtual override {
    FunctionsOwnerAcceptTermsOfServiceSetup.setUp();

    s_functionsClient = new FunctionsClientUpgradeHelper(address(s_functionsRouter));
  }
}

/// @notice Set up to create a subscription, add the consumer contract as a consumer of the subscription, and fund the subscription with 's_subscriptionInitialFunding'
contract FunctionsSubscriptionSetup is FunctionsClientSetup {
  uint64 s_subscriptionId;
  uint96 s_subscriptionInitialFunding = 10 * JUELS_PER_LINK; // 10 LINK

  function setUp() public virtual override {
    FunctionsClientSetup.setUp();

    // Create subscription
    s_subscriptionId = s_functionsRouter.createSubscription();
    s_functionsRouter.addConsumer(s_subscriptionId, address(s_functionsClient));

    // Fund subscription
    s_linkToken.transferAndCall(address(s_functionsRouter), s_subscriptionInitialFunding, abi.encode(s_subscriptionId));
  }
}

/// @notice Set up to initate a minimal request and store it in s_requests[1]
contract FunctionsClientRequestSetup is FunctionsSubscriptionSetup {
  struct RequestData {
    string sourceCode;
    bytes secrets;
    string[] args;
    bytes[] bytesArgs;
    uint32 callbackGasLimit;
  }
  struct Request {
    RequestData requestData;
    bytes32 requestId;
    FunctionsResponse.Commitment commitment;
  }

  mapping(uint256 => Request) s_requests;

  uint96 s_fulfillmentRouterOwnerBalance = 0;
  uint96 s_fulfillmentCoordinatorBalance = 0;

  function setUp() public virtual override {
    FunctionsSubscriptionSetup.setUp();

    // Send request #1
    string memory sourceCode = "return 'hello world';";
    bytes memory secrets = new bytes(0);
    string[] memory args = new string[](0);
    bytes[] memory bytesArgs = new bytes[](0);
    uint32 callbackGasLimit = 5500;
    _sendAndStoreRequest(1, sourceCode, secrets, args, bytesArgs, callbackGasLimit);
  }

  function _getExpectedCost(uint256 gasUsed) internal view returns (uint96 totalCostJuels) {
    uint96 juelsPerGas = uint96((1e18 * TX_GASPRICE_START) / uint256(LINK_ETH_RATE));
    uint96 gasOverheadJuels = juelsPerGas *
      (getCoordinatorConfig().gasOverheadBeforeCallback + getCoordinatorConfig().gasOverheadAfterCallback);
    uint96 callbackGasCostJuels = uint96(juelsPerGas * gasUsed);
    return gasOverheadJuels + s_donFee + s_adminFee + callbackGasCostJuels;
  }

  /// @notice Send a request and store information about it in s_requests
  /// @param requestNumberKey - the key that the request will be stored in `s_requests` in
  /// @param sourceCode - Raw source code for Request.codeLocation of Location.Inline, URL for Request.codeLocation of Location.Remote, or slot decimal number for Request.codeLocation of Location.DONHosted
  /// @param secrets - Encrypted URLs for Request.secretsLocation of Location.Remote (use addSecretsReference()), or CBOR encoded slotid+version for Request.secretsLocation of Location.DONHosted (use addDONHostedSecrets())
  /// @param args - String arguments that will be passed into the source code
  /// @param bytesArgs - Bytes arguments that will be passed into the source code
  /// @param callbackGasLimit - Gas limit for the fulfillment callback
  function _sendAndStoreRequest(
    uint256 requestNumberKey,
    string memory sourceCode,
    bytes memory secrets,
    string[] memory args,
    bytes[] memory bytesArgs,
    uint32 callbackGasLimit
  ) internal {
    if (s_requests[requestNumberKey].requestId != bytes32(0)) {
      revert("Request already written");
    }

    vm.recordLogs();

    bytes32 requestId = s_functionsClient.sendRequest(
      s_donId,
      sourceCode,
      secrets,
      args,
      bytesArgs,
      s_subscriptionId,
      callbackGasLimit
    );

    // Get commitment data from OracleRequest event log
    Vm.Log[] memory entries = vm.getRecordedLogs();
    (, , , , , , , FunctionsResponse.Commitment memory commitment) = abi.decode(
      entries[0].data,
      (address, uint64, address, bytes, uint16, bytes32, uint64, FunctionsResponse.Commitment)
    );
    s_requests[requestNumberKey] = Request({
      requestData: RequestData({
        sourceCode: sourceCode,
        secrets: secrets,
        args: args,
        bytesArgs: bytesArgs,
        callbackGasLimit: callbackGasLimit
      }),
      requestId: requestId,
      commitment: commitment
    });
  }

  /// @notice Send a request and store information about it in s_requests
  /// @param requestNumberKeys - One or more requestNumberKeys that were used to store the request in `s_requests` of the requests, that will be added to the report
  /// @param results - The result that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @param errors - The error that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @return report - Report bytes data
  /// @return reportContext - Report context bytes32 data
  function _buildReport(
    uint256[] memory requestNumberKeys,
    string[] memory results,
    bytes[] memory errors
  ) internal view returns (bytes memory report, bytes32[3] memory reportContext) {
    // Build report
    bytes32[] memory _requestIds = new bytes32[](requestNumberKeys.length);
    bytes[] memory _results = new bytes[](requestNumberKeys.length);
    bytes[] memory _errors = new bytes[](requestNumberKeys.length);
    bytes[] memory _onchainMetadata = new bytes[](requestNumberKeys.length);
    bytes[] memory _offchainMetadata = new bytes[](requestNumberKeys.length);
    for (uint256 i = 0; i < requestNumberKeys.length; ++i) {
      if (keccak256(bytes(results[i])) != keccak256(new bytes(0)) && keccak256(errors[i]) != keccak256(new bytes(0))) {
        revert("Report can only contain a result OR an error, one must remain empty.");
      }
      _requestIds[i] = s_requests[requestNumberKeys[i]].requestId;
      _results[i] = bytes(results[i]);
      _errors[i] = errors[i];
      _onchainMetadata[i] = abi.encode(s_requests[requestNumberKeys[i]].commitment);
      _offchainMetadata[i] = new bytes(0); // No off-chain metadata
    }
    report = abi.encode(_requestIds, _results, _errors, _onchainMetadata, _offchainMetadata);

    // Build report context
    uint256 h = uint256(
      keccak256(
        abi.encode(
          block.chainid,
          address(s_functionsCoordinator),
          1,
          s_signers,
          s_transmitters,
          s_f,
          s_onchainConfig,
          s_offchainConfigVersion,
          s_offchainConfig
        )
      )
    );
    uint256 prefixMask = type(uint256).max << (256 - 16); // 0xFFFF00..00
    uint256 prefix = 0x0001 << (256 - 16); // 0x000100..00
    bytes32 configDigest = bytes32((prefix & prefixMask) | (h & ~prefixMask));
    reportContext = [configDigest, configDigest, configDigest];

    return (report, reportContext);
  }

  /// @notice Gather signatures on report data
  /// @param report - Report bytes generated from `_buildReport`
  /// @param reportContext - Report context bytes32 generated from `_buildReport`
  /// @param signerPrivateKeys - One or more addresses that will sign the report data
  /// @return rawRs - Signature rs
  /// @return rawSs - Signature ss
  /// @return rawVs - Signature vs
  function _signReport(
    bytes memory report,
    bytes32[3] memory reportContext,
    uint256[] memory signerPrivateKeys
  ) internal pure returns (bytes32[] memory rawRs, bytes32[] memory rawSs, bytes32 rawVs) {
    bytes32[] memory rs = new bytes32[](signerPrivateKeys.length);
    bytes32[] memory ss = new bytes32[](signerPrivateKeys.length);
    bytes memory vs = new bytes(signerPrivateKeys.length);

    bytes32 reportDigest = keccak256(abi.encodePacked(keccak256(report), reportContext));

    for (uint256 i = 0; i < signerPrivateKeys.length; i++) {
      (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPrivateKeys[i], reportDigest);
      rs[i] = r;
      ss[i] = s;
      vs[i] = bytes1(v - 27);
    }

    return (rs, ss, bytes32(vs));
  }

  /// @notice Provide a response from the DON to fulfill one or more requests and store the updated balances of the DON & Admin
  /// @param requestNumberKeys - One or more requestNumberKeys that were used to store the request in `s_requests` of the requests, that will be added to the report
  /// @param results - The result that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @param errors - The error that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @param transmitter - The address that will send the `.report` transaction
  /// @param expectedToSucceed - Boolean representing if the report transmission is expected to produce a RequestProcessed event for every fulfillment. If not, we ignore retrieving the event log.
  /// @param requestProcessedIndex - On a successful fulfillment the Router will emit a RequestProcessed event. To grab that event we must know the order at which this event was thrown in the report transmission lifecycle. This can change depending on the test setup (e.g. the Client contract gives an extra event during its callback)
  /// @param transmitterGasToUse - Override the default amount of gas that the transmitter sends the `.report` transaction with
  function _reportAndStore(
    uint256[] memory requestNumberKeys,
    string[] memory results,
    bytes[] memory errors,
    address transmitter,
    bool expectedToSucceed,
    uint8 requestProcessedIndex,
    uint256 transmitterGasToUse
  ) internal {
    {
      if (requestNumberKeys.length != results.length || requestNumberKeys.length != errors.length) {
        revert("_reportAndStore arguments length mismatch");
      }
    }

    (bytes memory report, bytes32[3] memory reportContext) = _buildReport(requestNumberKeys, results, errors);

    // Sign the report
    // Need at least 3 signers to fulfill minimum number of: (configInfo.n + configInfo.f) / 2 + 1
    uint256[] memory signerPrivateKeys = new uint256[](3);
    signerPrivateKeys[0] = NOP_SIGNER_PRIVATE_KEY_1;
    signerPrivateKeys[1] = NOP_SIGNER_PRIVATE_KEY_2;
    signerPrivateKeys[2] = NOP_SIGNER_PRIVATE_KEY_3;
    (bytes32[] memory rawRs, bytes32[] memory rawSs, bytes32 rawVs) = _signReport(
      report,
      reportContext,
      signerPrivateKeys
    );

    // Send as transmitter
    vm.stopPrank();
    vm.startPrank(transmitter);

    // Send report
    vm.recordLogs();
    if (transmitterGasToUse > 0) {
      s_functionsCoordinator.transmit{gas: transmitterGasToUse}(reportContext, report, rawRs, rawSs, rawVs);
    } else {
      s_functionsCoordinator.transmit(reportContext, report, rawRs, rawSs, rawVs);
    }

    if (expectedToSucceed) {
      // Get actual cost from RequestProcessed event log
      (uint96 totalCostJuels, , , , , ) = abi.decode(
        vm.getRecordedLogs()[requestProcessedIndex].data,
        (uint96, address, FunctionsResponse.FulfillResult, bytes, bytes, bytes)
      );
      // Store profit amounts
      s_fulfillmentRouterOwnerBalance += s_adminFee;
      // totalCostJuels = costWithoutCallbackJuels + adminFee + callbackGasCostJuels
      s_fulfillmentCoordinatorBalance += totalCostJuels - s_adminFee;
    }

    // Return prank to Owner
    vm.stopPrank();
    vm.startPrank(OWNER_ADDRESS);
  }

  /// @notice Provide a response from the DON to fulfill one or more requests and store the updated balances of the DON & Admin
  /// @param requestNumberKeys - One or more requestNumberKeys that were used to store the request in `s_requests` of the requests, that will be added to the report
  /// @param results - The result that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @param errors - The error that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @param transmitter - The address that will send the `.report` transaction
  /// @param expectedToSucceed - Boolean representing if the report transmission is expected to produce a RequestProcessed event for every fulfillment. If not, we ignore retrieving the event log.
  /// @param requestProcessedIndex - On a successful fulfillment the Router will emit a RequestProcessed event. To grab that event we must know the order at which this event was thrown in the report transmission lifecycle. This can change depending on the test setup (e.g. the Client contract gives an extra event during its callback)
  /// @dev @param transmitterGasToUse is overloaded to give transmitterGasToUse as 0] - Sends the `.report` transaction with the default amount of gas
  function _reportAndStore(
    uint256[] memory requestNumberKeys,
    string[] memory results,
    bytes[] memory errors,
    address transmitter,
    bool expectedToSucceed,
    uint8 requestProcessedIndex
  ) internal {
    _reportAndStore(requestNumberKeys, results, errors, transmitter, expectedToSucceed, requestProcessedIndex, 0);
  }

  /// @notice Provide a response from the DON to fulfill one or more requests and store the updated balances of the DON & Admin
  /// @param requestNumberKeys - One or more requestNumberKeys that were used to store the request in `s_requests` of the requests, that will be added to the report
  /// @param results - The result that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @param errors - The error that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @param transmitter - The address that will send the `.report` transaction
  /// @param expectedToSucceed - Boolean representing if the report transmission is expected to produce a RequestProcessed event for every fulfillment. If not, we ignore retrieving the event log.
  /// @dev @param requestProcessedIndex is overloaded to give requestProcessedIndex as 3 (happy path value)] - On a successful fulfillment the Router will emit a RequestProcessed event. To grab that event we must know the order at which this event was thrown in the report transmission lifecycle. This can change depending on the test setup (e.g. the Client contract gives an extra event during its callback)
  /// @dev @param transmitterGasToUse is overloaded to give transmitterGasToUse as 0] - Sends the `.report` transaction with the default amount of gas
  function _reportAndStore(
    uint256[] memory requestNumberKeys,
    string[] memory results,
    bytes[] memory errors,
    address transmitter,
    bool expectedToSucceed
  ) internal {
    _reportAndStore(requestNumberKeys, results, errors, transmitter, expectedToSucceed, 3);
  }

  /// @notice Provide a response from the DON to fulfill one or more requests and store the updated balances of the DON & Admin
  /// @param requestNumberKeys - One or more requestNumberKeys that were used to store the request in `s_requests` of the requests, that will be added to the report
  /// @param results - The result that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @param errors - The error that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @param transmitter - The address that will send the `.report` transaction
  /// @dev @param expectedToSucceed is overloaded to give the value as true - The report transmission is expected to produce a RequestProcessed event for every fulfillment
  /// @dev @param requestProcessedIndex is overloaded to give requestProcessedIndex as 3 (happy path value)] - On a successful fulfillment the Router will emit a RequestProcessed event. To grab that event we must know the order at which this event was thrown in the report transmission lifecycle. This can change depending on the test setup (e.g. the Client contract gives an extra event during its callback)
  /// @dev @param transmitterGasToUse is overloaded to give transmitterGasToUse as 0] - Sends the `.report` transaction with the default amount of gas
  function _reportAndStore(
    uint256[] memory requestNumberKeys,
    string[] memory results,
    bytes[] memory errors,
    address transmitter
  ) internal {
    _reportAndStore(requestNumberKeys, results, errors, transmitter, true);
  }

  /// @notice Provide a response from the DON to fulfill one or more requests and store the updated balances of the DON & Admin
  /// @param requestNumberKeys - One or more requestNumberKeys that were used to store the request in `s_requests` of the requests, that will be added to the report
  /// @param results - The result that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @param errors - The error that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @dev @param transmitter is overloaded to give the value of transmitter #1 - The address that will send the `.report` transaction
  /// @dev @param expectedToSucceed is overloaded to give the value as true - The report transmission is expected to produce a RequestProcessed event for every fulfillment
  /// @dev @param requestProcessedIndex is overloaded to give requestProcessedIndex as 3 (happy path value)] - On a successful fulfillment the Router will emit a RequestProcessed event. To grab that event we must know the order at which this event was thrown in the report transmission lifecycle. This can change depending on the test setup (e.g. the Client contract gives an extra event during its callback)
  /// @dev @param transmitterGasToUse is overloaded to give transmitterGasToUse as 0] - Sends the `.report` transaction with the default amount of gas
  function _reportAndStore(
    uint256[] memory requestNumberKeys,
    string[] memory results,
    bytes[] memory errors
  ) internal {
    _reportAndStore(requestNumberKeys, results, errors, NOP_TRANSMITTER_ADDRESS_1);
  }
}

/// @notice Set up to have transmitter #1 send a report that fulfills request #1
contract FunctionsFulfillmentSetup is FunctionsClientRequestSetup {
  function setUp() public virtual override {
    FunctionsClientRequestSetup.setUp();

    // Fulfill request 1
    uint256[] memory requestNumberKeys = new uint256[](1);
    requestNumberKeys[0] = 1;
    string[] memory results = new string[](1);
    results[0] = "hello world!";
    bytes[] memory errors = new bytes[](1);
    errors[0] = new bytes(0);

    _reportAndStore(requestNumberKeys, results, errors, NOP_TRANSMITTER_ADDRESS_1, true);
  }
}

/// @notice Set up to send and fulfill two more requests, s_request[2] reported by transmitter #2 and s_request[3] reported by transmitter #3
contract FunctionsMultipleFulfillmentsSetup is FunctionsFulfillmentSetup {
  function setUp() public virtual override {
    FunctionsFulfillmentSetup.setUp();

    // Make 2 additional requests (1 already complete)

    //  *** Request #2 ***
    // Send
    string memory sourceCode = "return 'hello world';";
    bytes memory secrets = new bytes(0);
    string[] memory args = new string[](0);
    bytes[] memory bytesArgs = new bytes[](0);
    uint32 callbackGasLimit = 5500;
    _sendAndStoreRequest(2, sourceCode, secrets, args, bytesArgs, callbackGasLimit);
    // Fulfill as transmitter #2
    uint256[] memory requestNumberKeys1 = new uint256[](1);
    requestNumberKeys1[0] = 2;
    string[] memory results1 = new string[](1);
    results1[0] = "hello world!";
    bytes[] memory errors1 = new bytes[](1);
    errors1[0] = new bytes(0);
    _reportAndStore(requestNumberKeys1, results1, errors1, NOP_TRANSMITTER_ADDRESS_2, true);

    //  *** Request #3 ***
    // Send
    _sendAndStoreRequest(3, sourceCode, secrets, args, bytesArgs, callbackGasLimit);
    // Fulfill as transmitter #3
    uint256[] memory requestNumberKeys2 = new uint256[](1);
    requestNumberKeys2[0] = 3;
    string[] memory results2 = new string[](1);
    results2[0] = "hello world!";
    bytes[] memory errors2 = new bytes[](1);
    errors2[0] = new bytes(0);
    _reportAndStore(requestNumberKeys2, results2, errors2, NOP_TRANSMITTER_ADDRESS_3, true);
  }
}
