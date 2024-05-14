// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {FunctionsClientHarness} from "./testhelpers/FunctionsClientHarness.sol";
import {FunctionsRouterHarness, FunctionsRouter} from "./testhelpers/FunctionsRouterHarness.sol";
import {FunctionsCoordinatorHarness} from "./testhelpers/FunctionsCoordinatorHarness.sol";
import {FunctionsBilling} from "../../dev/v1_X/FunctionsBilling.sol";
import {FunctionsResponse} from "../../dev/v1_X/libraries/FunctionsResponse.sol";
import {MockV3Aggregator} from "../../../tests/MockV3Aggregator.sol";
import {TermsOfServiceAllowList} from "../../dev/v1_X/accessControl/TermsOfServiceAllowList.sol";
import {TermsOfServiceAllowListConfig} from "../../dev/v1_X/accessControl/interfaces/ITermsOfServiceAllowList.sol";
import {MockLinkToken} from "../../../mocks/MockLinkToken.sol";
import {FunctionsBillingConfig} from "../../dev/v1_X/interfaces/IFunctionsBilling.sol";

import "forge-std/Vm.sol";

/// @notice Set up to deploy the following contracts: FunctionsRouter, FunctionsCoordinator, LINK/ETH Feed, ToS Allow List, and LINK token
contract FunctionsRouterSetup is BaseTest {
  FunctionsRouterHarness internal s_functionsRouter;
  FunctionsCoordinatorHarness internal s_functionsCoordinator;
  MockV3Aggregator internal s_linkEthFeed;
  MockV3Aggregator internal s_linkUsdFeed;
  TermsOfServiceAllowList internal s_termsOfServiceAllowList;
  MockLinkToken internal s_linkToken;

  uint16 internal s_maxConsumersPerSubscription = 3;
  uint72 internal s_adminFee = 0; // Keep as 0. Setting this to anything else will cause fulfillments to fail with INVALID_COMMITMENT
  uint16 internal s_donFee = 100; // $1
  uint16 internal s_operationFee = 100; // $1
  bytes4 internal s_handleOracleFulfillmentSelector = 0x0ca76175;
  uint16 s_subscriptionDepositMinimumRequests = 1;
  uint72 s_subscriptionDepositJuels = 11 * JUELS_PER_LINK;

  int256 internal LINK_ETH_RATE = 6_000_000_000_000_000;
  uint8 internal LINK_ETH_DECIMALS = 18;
  int256 internal LINK_USD_RATE = 1_500_000_000;
  uint8 internal LINK_USD_DECIMALS = 8;

  uint256 internal TOS_SIGNER_PRIVATE_KEY = 0x3;
  address internal TOS_SIGNER = vm.addr(TOS_SIGNER_PRIVATE_KEY);

  function setUp() public virtual override {
    BaseTest.setUp();
    s_linkToken = new MockLinkToken();
    s_functionsRouter = new FunctionsRouterHarness(address(s_linkToken), getRouterConfig());
    s_linkEthFeed = new MockV3Aggregator(LINK_ETH_DECIMALS, LINK_ETH_RATE);
    s_linkUsdFeed = new MockV3Aggregator(LINK_USD_DECIMALS, LINK_USD_RATE);
    s_functionsCoordinator = new FunctionsCoordinatorHarness(
      address(s_functionsRouter),
      getCoordinatorConfig(),
      address(s_linkEthFeed),
      address(s_linkUsdFeed)
    );
    address[] memory initialAllowedSenders;
    address[] memory initialBlockedSenders;
    s_termsOfServiceAllowList = new TermsOfServiceAllowList(
      getTermsOfServiceConfig(),
      initialAllowedSenders,
      initialBlockedSenders
    );
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

  function getCoordinatorConfig() public view returns (FunctionsBillingConfig memory) {
    return
      FunctionsBillingConfig({
        feedStalenessSeconds: 24 * 60 * 60, // 1 day
        gasOverheadAfterCallback: 93_942,
        gasOverheadBeforeCallback: 105_000,
        requestTimeoutSeconds: 60 * 5, // 5 minutes
        donFeeCentsUsd: s_donFee,
        operationFeeCentsUsd: s_operationFee,
        maxSupportedRequestDataVersion: 1,
        fulfillmentGasPriceOverEstimationBP: 5000,
        fallbackNativePerUnitLink: 5000000000000000,
        fallbackUsdPerUnitLink: 1400000000,
        fallbackUsdPerUnitLinkDecimals: 8,
        minimumEstimateGasPriceWei: 1000000000 // 1 gwei
      });
  }

  function getTermsOfServiceConfig() public view returns (TermsOfServiceAllowListConfig memory) {
    return TermsOfServiceAllowListConfig({enabled: true, signerPublicKey: TOS_SIGNER});
  }
}

/// @notice Set up to set the OCR configuration of the Coordinator contract
contract FunctionsDONSetup is FunctionsRouterSetup {
  uint256 internal NOP_SIGNER_PRIVATE_KEY_1 = 0x400;
  address internal NOP_SIGNER_ADDRESS_1 = vm.addr(NOP_SIGNER_PRIVATE_KEY_1);
  uint256 internal NOP_SIGNER_PRIVATE_KEY_2 = 0x401;
  address internal NOP_SIGNER_ADDRESS_2 = vm.addr(NOP_SIGNER_PRIVATE_KEY_2);
  uint256 internal NOP_SIGNER_PRIVATE_KEY_3 = 0x402;
  address internal NOP_SIGNER_ADDRESS_3 = vm.addr(NOP_SIGNER_PRIVATE_KEY_3);
  uint256 internal NOP_SIGNER_PRIVATE_KEY_4 = 0x403;
  address internal NOP_SIGNER_ADDRESS_4 = vm.addr(NOP_SIGNER_PRIVATE_KEY_4);

  uint256 internal NOP_TRANSMITTER_PRIVATE_KEY_1 = 0x404;
  address internal NOP_TRANSMITTER_ADDRESS_1 = vm.addr(NOP_TRANSMITTER_PRIVATE_KEY_1);
  uint256 internal NOP_TRANSMITTER_PRIVATE_KEY_2 = 0x405;
  address internal NOP_TRANSMITTER_ADDRESS_2 = vm.addr(NOP_TRANSMITTER_PRIVATE_KEY_2);
  uint256 internal NOP_TRANSMITTER_PRIVATE_KEY_3 = 0x406;
  address internal NOP_TRANSMITTER_ADDRESS_3 = vm.addr(NOP_TRANSMITTER_PRIVATE_KEY_3);
  uint256 internal NOP_TRANSMITTER_PRIVATE_KEY_4 = 0x407;
  address internal NOP_TRANSMITTER_ADDRESS_4 = vm.addr(NOP_TRANSMITTER_PRIVATE_KEY_4);

  address[] internal s_signers;
  address[] internal s_transmitters;
  uint8 s_f = 1;
  bytes internal s_onchainConfig = new bytes(0);
  uint64 internal s_offchainConfigVersion = 1;
  bytes internal s_offchainConfig = new bytes(0);

  bytes s_thresholdKey =
    vm.parseBytes(
      "0x7b2247726f7570223a2250323536222c22475f626172223a22424f2f344358424575792f64547a436a612b614e774d666c2b645a77346d325036533246536b4966472f6633527547327337392b494e79642b4639326a346f586e67433657427561556a752b4a637a32377834484251343d222c2248223a224250532f72485065377941467232416c447a79395549466258776d46384666756632596d514177666e3342373844336f474845643247474536466e616f34552b4c6a4d4d5756792b464f7075686e77554f6a75427a64773d222c22484172726179223a5b22424d75546862414473337768316e67764e56792f6e3841316d42674b5a4b4c475259385937796a39695769337242502f316a32347571695869534531437554384c6f51446a386248466d384345477667517158494e62383d222c224248687974716d6e34314373322f4658416f43737548687151486236382f597930524b2b41354c6647654f645a78466f4e386c442b45656e4b587a544943784f6d3231636d535447364864484a6e336342645663714c673d222c22424d794e7a4534616e596258474d72694f52664c52634e7239766c347878654279316432452f4464335a744630546372386267567435582b2b42355967552b4b7875726e512f4d656b6857335845782b79506e4e4f584d3d222c22424d6a753272375a657a4a45545539413938746a6b6d547966796a79493735345742555835505174724a6578346d6766366130787373426d50325a7472412b55576d504e592b6d4664526b46674f7944694c53614e59453d225d7d"
    );
  bytes s_donKey =
    vm.parseBytes(
      "0xf2f9c47363202d89aa9fa70baf783d70006fe493471ac8cfa82f1426fd09f16a5f6b32b7c4b5d5165cd147a6e513ba4c0efd39d969d6b20a8a21126f0411b9c6"
    );

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

  function _getTransmitterBalances() internal view returns (uint256[4] memory balances) {
    return [
      s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_1),
      s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_2),
      s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_3),
      s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_4)
    ];
  }

  function _assertTransmittersAllHaveBalance(uint256[4] memory balances, uint256 expectedBalance) internal {
    assertEq(balances[0], expectedBalance);
    assertEq(balances[1], expectedBalance);
    assertEq(balances[2], expectedBalance);
    assertEq(balances[3], expectedBalance);
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
  FunctionsClientHarness internal s_functionsClient;

  function setUp() public virtual override {
    FunctionsOwnerAcceptTermsOfServiceSetup.setUp();

    s_functionsClient = new FunctionsClientHarness(address(s_functionsRouter));
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
  struct Report {
    bytes32[] rs;
    bytes32[] ss;
    bytes32 vs;
    bytes report;
    bytes32[3] reportContext;
  }

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
    FunctionsResponse.Commitment commitment; // Offchain commitment that contains operation fee in the place of admin fee
    FunctionsResponse.Commitment commitmentOnchain; // Commitment that is persisted as a hash in the Router
  }

  mapping(uint256 requestNumber => Request) s_requests;

  struct Response {
    uint96 totalCostJuels;
  }

  mapping(uint256 requestNumber => Response) s_responses;

  uint96 s_fulfillmentRouterOwnerBalance = 0;
  uint96 s_fulfillmentCoordinatorBalance = 0;
  uint8 s_requestsSent = 0;
  uint8 s_requestsFulfilled = 0;

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

  /// @notice Predicts the estimated cost (maximum cost) of a request
  /// @dev Meant only for Ethereum, does not add L2 chains' L1 fee
  function _getExpectedCostEstimate(uint256 callbackGas) internal view returns (uint96) {
    uint256 gasPrice = TX_GASPRICE_START < getCoordinatorConfig().minimumEstimateGasPriceWei
      ? getCoordinatorConfig().minimumEstimateGasPriceWei
      : TX_GASPRICE_START;
    uint256 gasPriceWithOverestimation = gasPrice +
      ((gasPrice * getCoordinatorConfig().fulfillmentGasPriceOverEstimationBP) / 10_000);
    uint96 juelsPerGas = uint96((1e18 * gasPriceWithOverestimation) / uint256(LINK_ETH_RATE));
    uint96 gasOverheadJuels = juelsPerGas *
      ((getCoordinatorConfig().gasOverheadBeforeCallback + getCoordinatorConfig().gasOverheadAfterCallback));
    uint96 callbackGasCostJuels = uint96(juelsPerGas * callbackGas);
    bytes memory emptyData = new bytes(0);
    return
      gasOverheadJuels +
      s_functionsCoordinator.getDONFeeJuels(emptyData) +
      s_adminFee +
      s_functionsCoordinator.getOperationFeeJuels() +
      callbackGasCostJuels;
  }

  /// @notice Predicts the actual cost of a request
  /// @dev Meant only for Ethereum, does not add L2 chains' L1 fee
  function _getExpectedCost(uint256 gasUsed) internal view returns (uint96) {
    uint96 juelsPerGas = uint96((1e18 * TX_GASPRICE_START) / uint256(LINK_ETH_RATE));
    uint96 gasOverheadJuels = juelsPerGas *
      (getCoordinatorConfig().gasOverheadBeforeCallback + getCoordinatorConfig().gasOverheadAfterCallback);
    uint96 callbackGasCostJuels = uint96(juelsPerGas * gasUsed);
    bytes memory emptyData = new bytes(0);
    return
      gasOverheadJuels +
      s_functionsCoordinator.getDONFeeJuels(emptyData) +
      s_adminFee +
      s_functionsCoordinator.getOperationFeeJuels() +
      callbackGasCostJuels;
  }

  /// @notice Send a request and store information about it in s_requests
  /// @param requestNumberKey - the key that the request will be stored in `s_requests` in
  /// @param sourceCode - Raw source code for Request.codeLocation of Location.Inline, URL for Request.codeLocation of Location.Remote, or slot decimal number for Request.codeLocation of Location.DONHosted
  /// @param secrets - Encrypted URLs for Request.secretsLocation of Location.Remote (use addSecretsReference()), or CBOR encoded slotid+version for Request.secretsLocation of Location.DONHosted (use addDONHostedSecrets())
  /// @param args - String arguments that will be passed into the source code
  /// @param bytesArgs - Bytes arguments that will be passed into the source code
  /// @param callbackGasLimit - Gas limit for the fulfillment callback
  /// @param client - The consumer contract to send the request from
  function _sendAndStoreRequest(
    uint256 requestNumberKey,
    string memory sourceCode,
    bytes memory secrets,
    string[] memory args,
    bytes[] memory bytesArgs,
    uint32 callbackGasLimit,
    address client
  ) internal {
    if (s_requests[requestNumberKey].requestId != bytes32(0)) {
      revert("Request already written");
    }

    vm.recordLogs();

    bytes32 requestId = FunctionsClientHarness(client).sendRequest(
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
      commitment: commitment, // Has operationFee in place of adminFee
      commitmentOnchain: FunctionsResponse.Commitment({
        coordinator: commitment.coordinator,
        client: commitment.client,
        subscriptionId: commitment.subscriptionId,
        callbackGasLimit: commitment.callbackGasLimit,
        estimatedTotalCostJuels: commitment.estimatedTotalCostJuels,
        timeoutTimestamp: commitment.timeoutTimestamp,
        requestId: commitment.requestId,
        donFee: commitment.donFee,
        gasOverheadBeforeCallback: commitment.gasOverheadBeforeCallback,
        gasOverheadAfterCallback: commitment.gasOverheadAfterCallback,
        adminFee: s_adminFee
      })
    });
    s_requestsSent += 1;
  }

  /// @notice Send a request and store information about it in s_requests
  /// @param requestNumberKey - the key that the request will be stored in `s_requests` in
  /// @param sourceCode - Raw source code for Request.codeLocation of Location.Inline, URL for Request.codeLocation of Location.Remote, or slot decimal number for Request.codeLocation of Location.DONHosted
  /// @param secrets - Encrypted URLs for Request.secretsLocation of Location.Remote (use addSecretsReference()), or CBOR encoded slotid+version for Request.secretsLocation of Location.DONHosted (use addDONHostedSecrets())
  /// @param args - String arguments that will be passed into the source code
  /// @param bytesArgs - Bytes arguments that will be passed into the source code
  /// @param callbackGasLimit - Gas limit for the fulfillment callback
  /// @dev @param client - The consumer contract to send the request from (overloaded to fill client with s_functionsClient)
  function _sendAndStoreRequest(
    uint256 requestNumberKey,
    string memory sourceCode,
    bytes memory secrets,
    string[] memory args,
    bytes[] memory bytesArgs,
    uint32 callbackGasLimit
  ) internal {
    _sendAndStoreRequest(
      requestNumberKey,
      sourceCode,
      secrets,
      args,
      bytesArgs,
      callbackGasLimit,
      address(s_functionsClient)
    );
  }

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
  ) internal pure returns (bytes32[] memory, bytes32[] memory, bytes32) {
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

  function _buildAndSignReport(
    uint256[] memory requestNumberKeys,
    string[] memory results,
    bytes[] memory errors
  ) internal view returns (Report memory) {
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

    return Report({report: report, reportContext: reportContext, rs: rawRs, ss: rawSs, vs: rawVs});
  }

  /// @notice Provide a response from the DON to fulfill one or more requests and store the updated balances of the DON & Admin
  /// @param requestNumberKeys - One or more requestNumberKeys that were used to store the request in `s_requests` of the requests, that will be added to the report
  /// @param results - The result that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @param errors - The error that will be sent to the consumer contract's callback. For each index, e.g. result[index] or errors[index], only one of should be filled.
  /// @param transmitter - The address that will send the `.report` transaction
  /// @param expectedToSucceed - Boolean representing if the report transmission is expected to produce a RequestProcessed event for every fulfillment. If not, we ignore retrieving the event log.
  /// @param requestProcessedStartIndex - On a successful fulfillment the Router will emit a RequestProcessed event. To grab that event we must know the order at which this event was thrown in the report transmission lifecycle. This can change depending on the test setup (e.g. the Client contract gives an extra event during its callback)
  /// @param transmitterGasToUse - Override the default amount of gas that the transmitter sends the `.report` transaction with
  function _reportAndStore(
    uint256[] memory requestNumberKeys,
    string[] memory results,
    bytes[] memory errors,
    address transmitter,
    bool expectedToSucceed,
    uint8 requestProcessedStartIndex,
    uint256 transmitterGasToUse
  ) internal {
    {
      if (requestNumberKeys.length != results.length || requestNumberKeys.length != errors.length) {
        revert("_reportAndStore arguments length mismatch");
      }
    }

    Report memory r = _buildAndSignReport(requestNumberKeys, results, errors);

    // Send as transmitter
    vm.stopPrank();
    vm.startPrank(transmitter, transmitter);

    // Send report
    vm.recordLogs();
    if (transmitterGasToUse > 0) {
      s_functionsCoordinator.transmit{gas: transmitterGasToUse}(r.reportContext, r.report, r.rs, r.ss, r.vs);
    } else {
      s_functionsCoordinator.transmit(r.reportContext, r.report, r.rs, r.ss, r.vs);
    }

    if (expectedToSucceed) {
      // Get actual cost from RequestProcessed event log
      (uint96 totalCostJuels, , , , , ) = abi.decode(
        vm.getRecordedLogs()[requestProcessedStartIndex].data,
        (uint96, address, FunctionsResponse.FulfillResult, bytes, bytes, bytes)
      );
      // Store response of first request
      // TODO: handle multiple requests
      s_responses[requestNumberKeys[0]] = Response({totalCostJuels: totalCostJuels});
      // Store profit amounts
      s_fulfillmentRouterOwnerBalance += s_adminFee * uint96(requestNumberKeys.length);
      // totalCostJuels = costWithoutCallbackJuels + adminFee + callbackGasCostJuels
      // TODO: handle multiple requests
      s_fulfillmentCoordinatorBalance += totalCostJuels - s_adminFee;
    }
    s_requestsFulfilled += 1;

    // Return prank to Owner
    vm.stopPrank();
    vm.startPrank(OWNER_ADDRESS, OWNER_ADDRESS);
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

    // Fast forward time by 30 seconds to simulate the DON executing the computation
    vm.warp(block.timestamp + 30);

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

    // Make 3 additional requests (1 already complete)

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

    //  *** Request #4 ***
    // Send
    _sendAndStoreRequest(4, sourceCode, secrets, args, bytesArgs, callbackGasLimit);
    // Fulfill as transmitter #1
    uint256[] memory requestNumberKeys3 = new uint256[](1);
    requestNumberKeys3[0] = 4;
    string[] memory results3 = new string[](1);
    results3[0] = "hello world!";
    bytes[] memory errors3 = new bytes[](1);
    errors3[0] = new bytes(0);
    _reportAndStore(requestNumberKeys3, results3, errors3, NOP_TRANSMITTER_ADDRESS_1, true);
  }
}
