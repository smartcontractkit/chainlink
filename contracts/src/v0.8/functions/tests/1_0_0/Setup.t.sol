// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {FunctionsRouter} from "../../dev/1_0_0/FunctionsRouter.sol";
import {FunctionsCoordinatorTestHelper} from "./testhelpers/FunctionsCoordinatorTestHelper.sol";
import {FunctionsBilling} from "../../dev/1_0_0/FunctionsBilling.sol";
import {FunctionsResponse} from "../../dev/1_0_0/libraries/FunctionsResponse.sol";
import {MockV3Aggregator} from "../../../tests/MockV3Aggregator.sol";
import {TermsOfServiceAllowList} from "../../dev/1_0_0/accessControl/TermsOfServiceAllowList.sol";
import {FunctionsClientUpgradeHelper} from "./testhelpers/FunctionsClientUpgradeHelper.sol";
import {MockLinkToken} from "../../../mocks/MockLinkToken.sol";

import "forge-std/Vm.sol";

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
  uint72 s_subscriptionDepositJuels = 11 * 1e18;

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
        gasOverheadAfterCallback: 44_615, // TODO: update
        gasOverheadBeforeCallback: 44_615, // TODO: update
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

  function setUp() public virtual override {
    FunctionsRouterSetup.setUp();

    address[] memory _signers = new address[](4);
    _signers[0] = NOP_SIGNER_ADDRESS_1;
    _signers[1] = NOP_SIGNER_ADDRESS_2;
    _signers[2] = NOP_SIGNER_ADDRESS_3;
    _signers[3] = NOP_SIGNER_ADDRESS_4;
    address[] memory _transmitters = new address[](4);
    _transmitters[0] = NOP_TRANSMITTER_ADDRESS_1;
    _transmitters[1] = NOP_TRANSMITTER_ADDRESS_2;
    _transmitters[2] = NOP_TRANSMITTER_ADDRESS_3;
    _transmitters[3] = NOP_TRANSMITTER_ADDRESS_4;
    uint8 _f = 1;
    bytes memory _onchainConfig = new bytes(0);
    uint64 _offchainConfigVersion = 1;
    bytes memory _offchainConfig = new bytes(0);
    // set OCR config
    s_functionsCoordinator.setConfig(
      _signers,
      _transmitters,
      _f,
      _onchainConfig,
      _offchainConfigVersion,
      _offchainConfig
    );
  }
}

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

contract FunctionsOwnerAcceptTermsOfServiceSetup is FunctionsRoutesSetup {
  function setUp() public virtual override {
    FunctionsRoutesSetup.setUp();

    bytes32 message = s_termsOfServiceAllowList.getMessage(OWNER_ADDRESS, OWNER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);
    s_termsOfServiceAllowList.acceptTermsOfService(OWNER_ADDRESS, OWNER_ADDRESS, r, s, v);
  }
}

contract FunctionsClientSetup is FunctionsOwnerAcceptTermsOfServiceSetup {
  FunctionsClientUpgradeHelper internal s_functionsClient;

  function setUp() public virtual override {
    FunctionsOwnerAcceptTermsOfServiceSetup.setUp();

    s_functionsClient = new FunctionsClientUpgradeHelper(address(s_functionsRouter));
  }
}

contract FunctionsSubscriptionSetup is FunctionsClientSetup {
  uint96 constant JUELS_PER_LINK = 1e18;
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

contract FunctionsClientRequestSetup is FunctionsSubscriptionSetup {
  bytes32 s_requestId;
  FunctionsResponse.Commitment s_requestCommitment;

  function setUp() public virtual override {
    FunctionsSubscriptionSetup.setUp();

    // Send a minimal request
    string memory sourceCode = "return 'hello world';";
    bytes memory secrets;
    string[] memory args = new string[](0);
    bytes[] memory bytesArgs = new bytes[](0);
    uint32 callbackGasLimit = 5500;

    vm.recordLogs();
    s_requestId = s_functionsClient.sendRequest(
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
    s_requestCommitment = commitment;
  }
}

contract FunctionsFulfillmentSetup is FunctionsClientRequestSetup {
  uint96 s_fulfillmentRouterOwnerBalance = s_adminFee;
  uint96 s_fulfillmentCoordinatorBalance;

  function setUp() public virtual override {
    FunctionsClientRequestSetup.setUp();

    // Send as transmitter 1
    vm.stopPrank();
    vm.startPrank(NOP_TRANSMITTER_ADDRESS_1);

    // Build report
    bytes32[] memory requestIds = new bytes32[](1);
    requestIds[0] = s_requestId;
    bytes[] memory results = new bytes[](1);
    results[0] = bytes("hello world!");
    bytes[] memory errors = new bytes[](1);
    // No error
    bytes[] memory onchainMetadata = new bytes[](1);
    onchainMetadata[0] = abi.encode(s_requestCommitment);
    bytes[] memory offchainMetadata = new bytes[](1);
    // No offchain metadata
    bytes memory report = abi.encode(requestIds, results, errors, onchainMetadata, offchainMetadata);

    // Build signers
    address[31] memory signers;
    signers[0] = NOP_SIGNER_ADDRESS_1;

    // Send report
    vm.recordLogs();
    s_functionsCoordinator.callReportWithSigners(report, signers);

    // Get actual cost from RequestProcessed event log
    Vm.Log[] memory entries = vm.getRecordedLogs();
    (uint96 totalCostJuels, , , , , ) = abi.decode(
      entries[2].data,
      (uint96, address, FunctionsResponse.FulfillResult, bytes, bytes, bytes)
    );
    // totalCostJuels = costWithoutCallbackJuels + adminFee + callbackGasCostJuels
    s_fulfillmentCoordinatorBalance = totalCostJuels - s_adminFee;

    // Return prank to Owner
    vm.stopPrank();
    vm.startPrank(OWNER_ADDRESS);
  }
}

contract FunctionsMultipleFulfillmentsSetup is FunctionsFulfillmentSetup {
  bytes32 s_requestId2;
  FunctionsResponse.Commitment s_requestCommitment2;
  bytes32 s_requestId3;
  FunctionsResponse.Commitment s_requestCommitment3;

  function setUp() public virtual override {
    FunctionsFulfillmentSetup.setUp();

    // Make 2 additional requests (1 already complete)

    //  *** Request #2 ***
    vm.recordLogs();
    s_requestId2 = s_functionsClient.sendRequest(
      s_donId,
      "return 'hello world';",
      new bytes(0),
      new string[](0),
      new bytes[](0),
      s_subscriptionId,
      5500
    );

    // Get commitment data from OracleRequest event log
    Vm.Log[] memory entriesAfterRequest2 = vm.getRecordedLogs();
    (, , , , , , , FunctionsResponse.Commitment memory commitment2) = abi.decode(
      entriesAfterRequest2[0].data,
      (address, uint64, address, bytes, uint16, bytes32, uint64, FunctionsResponse.Commitment)
    );
    s_requestCommitment2 = commitment2;

    // Transmit as transmitter 2
    vm.stopPrank();
    vm.startPrank(NOP_TRANSMITTER_ADDRESS_2);

    // Build report
    bytes32[] memory requestIds2 = new bytes32[](1);
    requestIds2[0] = s_requestId2;
    bytes[] memory results2 = new bytes[](1);
    results2[0] = bytes("hello world!");
    bytes[] memory errors2 = new bytes[](1);
    // No error
    bytes[] memory onchainMetadata2 = new bytes[](1);
    onchainMetadata2[0] = abi.encode(s_requestCommitment2);
    bytes[] memory offchainMetadata2 = new bytes[](1);
    // No offchain metadata
    bytes memory report2 = abi.encode(requestIds2, results2, errors2, onchainMetadata2, offchainMetadata2);

    // Build signers
    address[31] memory signers2;
    signers2[0] = NOP_SIGNER_ADDRESS_2;

    // Send report
    vm.recordLogs();
    s_functionsCoordinator.callReportWithSigners(report2, signers2);

    // Get actual cost from RequestProcessed event log
    Vm.Log[] memory entriesAfterFulfill2 = vm.getRecordedLogs();
    (uint96 totalCostJuels2, , , , , ) = abi.decode(
      entriesAfterFulfill2[2].data,
      (uint96, address, FunctionsResponse.FulfillResult, bytes, bytes, bytes)
    );
    // totalCostJuels = costWithoutCallbackJuels + adminFee + callbackGasCostJuels
    s_fulfillmentCoordinatorBalance += totalCostJuels2 - s_adminFee;
    s_fulfillmentRouterOwnerBalance += s_adminFee;

    // Return prank to Owner
    vm.stopPrank();
    vm.startPrank(OWNER_ADDRESS);

    //  *** Request #3 ***
    vm.recordLogs();
    s_requestId3 = s_functionsClient.sendRequest(
      s_donId,
      "return 'hello world';",
      new bytes(0),
      new string[](0),
      new bytes[](0),
      s_subscriptionId,
      5500
    );

    // Get commitment data from OracleRequest event log
    Vm.Log[] memory entriesAfterRequest3 = vm.getRecordedLogs();
    (, , , , , , , FunctionsResponse.Commitment memory commitment3) = abi.decode(
      entriesAfterRequest3[0].data,
      (address, uint64, address, bytes, uint16, bytes32, uint64, FunctionsResponse.Commitment)
    );
    s_requestCommitment3 = commitment3;

    // Transmit as transmitter 3
    vm.stopPrank();
    vm.startPrank(NOP_TRANSMITTER_ADDRESS_3);

    // Build report
    bytes32[] memory requestIds3 = new bytes32[](1);
    requestIds3[0] = s_requestId3;
    bytes[] memory results3 = new bytes[](1);
    results3[0] = bytes("hello world!");
    bytes[] memory errors3 = new bytes[](1);
    // No error
    bytes[] memory onchainMetadata3 = new bytes[](1);
    onchainMetadata3[0] = abi.encode(s_requestCommitment3);
    bytes[] memory offchainMetadata3 = new bytes[](1);
    // No offchain metadata
    bytes memory report3 = abi.encode(requestIds3, results3, errors3, onchainMetadata3, offchainMetadata3);

    // Build signers
    address[31] memory signers3;
    signers3[0] = NOP_SIGNER_ADDRESS_3;

    // Send report
    vm.recordLogs();
    s_functionsCoordinator.callReportWithSigners(report3, signers3);

    // Get actual cost from RequestProcessed event log
    Vm.Log[] memory entriesAfterFulfill3 = vm.getRecordedLogs();
    (uint96 totalCostJuels3, , , , , ) = abi.decode(
      entriesAfterFulfill3[2].data,
      (uint96, address, FunctionsResponse.FulfillResult, bytes, bytes, bytes)
    );

    // totalCostJuels = costWithoutCallbackJuels + adminFee + callbackGasCostJuels
    s_fulfillmentCoordinatorBalance += totalCostJuels3 - s_adminFee;

    // Return prank to Owner
    vm.stopPrank();
    vm.startPrank(OWNER_ADDRESS);
  }
}
