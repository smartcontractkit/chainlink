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

  /// @dev Overload to give transmitterGasToUse as 0, which gives default tx gas
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

  /// @dev Overload to give requestProcessedIndex as 3 (happy path value)
  function _reportAndStore(
    uint256[] memory requestNumberKeys,
    string[] memory results,
    bytes[] memory errors,
    address transmitter,
    bool expectedToSucceed
  ) internal {
    _reportAndStore(requestNumberKeys, results, errors, transmitter, expectedToSucceed, 3);
  }

  /// @dev Overload to give expectedToSucceed as true
  function _reportAndStore(
    uint256[] memory requestNumberKeys,
    string[] memory results,
    bytes[] memory errors,
    address transmitter
  ) internal {
    _reportAndStore(requestNumberKeys, results, errors, transmitter, true);
  }

  /// @dev Overload to give transmitter as transmitter #1
  function _reportAndStore(
    uint256[] memory requestNumberKeys,
    string[] memory results,
    bytes[] memory errors
  ) internal {
    _reportAndStore(requestNumberKeys, results, errors, NOP_TRANSMITTER_ADDRESS_1);
  }
}

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
