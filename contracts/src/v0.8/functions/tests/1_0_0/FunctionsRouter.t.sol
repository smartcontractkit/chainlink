// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IFunctionsRouter} from "../../dev/1_0_0/interfaces/IFunctionsRouter.sol";
import {IFunctionsBilling} from "../../dev/1_0_0/interfaces/IFunctionsBilling.sol";

import {BaseTest} from "./BaseTest.t.sol";
import {FunctionsRouter} from "../../dev/1_0_0/FunctionsRouter.sol";
import {FunctionsCoordinator} from "../../dev/1_0_0/FunctionsCoordinator.sol";
import {FunctionsBilling} from "../../dev/1_0_0/FunctionsBilling.sol";
import {MockV3Aggregator} from "../../../tests/MockV3Aggregator.sol";
import {TermsOfServiceAllowList} from "../../dev/1_0_0/accessControl/TermsOfServiceAllowList.sol";

contract FunctionsRouterSetup is BaseTest {
  FunctionsRouter internal s_functionsRouter;
  FunctionsCoordinator internal s_functionsCoordinator;
  MockV3Aggregator internal s_linkEthFeed;
  TermsOfServiceAllowList internal s_termsOfServiceAllowList;

  uint16 internal s_timelockBlocks = 0;
  uint16 internal s_maximumTimelockBlocks = 20;
  uint16 internal s_maxConsumersPerSubscription = 100;
  uint96 internal s_adminFee = 561724823;
  bytes4 internal s_handleOracleFulfillmentSelector = 0x0ca76175;

  address internal s_linkToken = 0x01BE23585060835E02B77ef475b0Cc51aA1e0709;

  int256 internal LINK_ETH_RATE = 5021530000000000;

  function setUp() public virtual override {
    BaseTest.setUp();
    s_functionsRouter = new FunctionsRouter(s_timelockBlocks, s_maximumTimelockBlocks, s_linkToken, getRouterConfig());
    s_linkEthFeed = new MockV3Aggregator(0, LINK_ETH_RATE);

    s_termsOfServiceAllowList = new TermsOfServiceAllowList(address(s_functionsRouter), getTermsOfServiceConfig());
  }

  function getRouterConfig() public view returns (IFunctionsRouter.Config memory) {
    uint32[] memory maxCallbackGasLimits = new uint32[](1);
    maxCallbackGasLimits[0] = type(uint32).max;

    return
      IFunctionsRouter.Config({
        maxConsumersPerSubscription: s_maxConsumersPerSubscription,
        adminFee: s_adminFee,
        handleOracleFulfillmentSelector: s_handleOracleFulfillmentSelector,
        maxCallbackGasLimits: maxCallbackGasLimits
      });
  }

  function getCoordinatorConfig() public pure returns (bytes memory) {
    IFunctionsBilling.Config memory billingConfig = IFunctionsBilling.Config({
      maxCallbackGasLimit: 5,
      feedStalenessSeconds: 5,
      gasOverheadAfterCallback: 5,
      gasOverheadBeforeCallback: 5,
      requestTimeoutSeconds: 1,
      donFee: 5,
      maxSupportedRequestDataVersion: 5,
      fulfillmentGasPriceOverEstimationBP: 5,
      fallbackNativePerUnitLink: 2874
    });

    return
      abi.encode(
        billingConfig.maxCallbackGasLimit,
        billingConfig.feedStalenessSeconds,
        billingConfig.gasOverheadBeforeCallback,
        billingConfig.gasOverheadAfterCallback,
        billingConfig.requestTimeoutSeconds,
        billingConfig.donFee,
        billingConfig.maxSupportedRequestDataVersion,
        billingConfig.fulfillmentGasPriceOverEstimationBP,
        billingConfig.fallbackNativePerUnitLink
      );
  }

  function getTermsOfServiceConfig() public pure returns (bytes memory) {
    bool enabled = false;
    address proofSignerPublicKey = address(132);
    return abi.encode(enabled, proofSignerPublicKey);
  }
}

contract FunctionsRouter_createSubscription is FunctionsRouterSetup {
  event SubscriptionCreated(uint64 indexed subscriptionId, address owner);

  function testCreateSubscriptionSuccess() public {
    s_functionsCoordinator = new FunctionsCoordinator(
      address(s_functionsRouter),
      getCoordinatorConfig(),
      address(s_linkEthFeed)
    );

    bytes32 allowListId = s_functionsRouter.getAllowListId();
    bytes32[] memory proposedContractSetIds = new bytes32[](2);
    proposedContractSetIds[0] = bytes32("1");
    proposedContractSetIds[1] = allowListId;
    address[] memory proposedContractSetAddresses = new address[](2);
    proposedContractSetAddresses[0] = address(s_functionsCoordinator);
    proposedContractSetAddresses[1] = address(s_termsOfServiceAllowList);

    s_functionsRouter.proposeContractsUpdate(proposedContractSetIds, proposedContractSetAddresses);
    s_functionsRouter.updateContracts();

    vm.expectEmit();
    emit SubscriptionCreated(1, OWNER);

    s_functionsRouter.createSubscription();

    vm.expectEmit();
    emit SubscriptionCreated(2, OWNER);

    s_functionsRouter.createSubscription();
  }
}
