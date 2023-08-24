// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

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

  uint16 internal s_maxConsumersPerSubscription = 100;
  uint72 internal s_adminFee = 100;
  bytes4 internal s_handleOracleFulfillmentSelector = 0x0ca76175;

  address internal s_linkToken = 0x01BE23585060835E02B77ef475b0Cc51aA1e0709;

  int256 internal LINK_ETH_RATE = 6000000000000000;

  uint256 internal TOS_SIGNER_PRIVATE_KEY = 0x3;
  address internal TOS_SIGNER = vm.addr(TOS_SIGNER_PRIVATE_KEY);

  function setUp() public virtual override {
    BaseTest.setUp();
    s_functionsRouter = new FunctionsRouter(s_linkToken, getRouterConfig());
    s_linkEthFeed = new MockV3Aggregator(0, LINK_ETH_RATE);
    s_functionsCoordinator = new FunctionsCoordinator(
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
        gasForCallExactCheck: 5000
      });
  }

  function getCoordinatorConfig() public pure returns (FunctionsBilling.Config memory) {
    return
      FunctionsBilling.Config({
        maxCallbackGasLimit: 0, // NOTE: unused , TODO: remove
        feedStalenessSeconds: 24 * 60 * 60, // 1 day
        gasOverheadAfterCallback: 44_615, // TODO: update
        gasOverheadBeforeCallback: 44_615, // TODO: update
        requestTimeoutSeconds: 60 * 5, // 5 minutes
        donFee: 100,
        maxSupportedRequestDataVersion: 1,
        fulfillmentGasPriceOverEstimationBP: 5000,
        fallbackNativePerUnitLink: 5000000000000000
      });
  }

  function getTermsOfServiceConfig() public view returns (TermsOfServiceAllowList.Config memory) {
    return TermsOfServiceAllowList.Config({enabled: true, signerPublicKey: TOS_SIGNER});
  }
}

contract FunctionsSetupRoutes is FunctionsRouterSetup {
  function setUp() public virtual override {
    FunctionsRouterSetup.setUp();

    bytes32 allowListId = s_functionsRouter.getAllowListId();
    bytes32[] memory proposedContractSetIds = new bytes32[](2);
    proposedContractSetIds[0] = bytes32("1");
    proposedContractSetIds[1] = allowListId;
    address[] memory proposedContractSetAddresses = new address[](2);
    proposedContractSetAddresses[0] = address(s_functionsCoordinator);
    proposedContractSetAddresses[1] = address(s_termsOfServiceAllowList);

    s_functionsRouter.proposeContractsUpdate(proposedContractSetIds, proposedContractSetAddresses);
    s_functionsRouter.updateContracts();
  }
}

contract FunctionsOwnerAcceptTermsOfService is FunctionsSetupRoutes {
  function setUp() public virtual override {
    FunctionsSetupRoutes.setUp();

    bytes32 message = s_termsOfServiceAllowList.getMessage(OWNER_ADDRESS, OWNER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);
    s_termsOfServiceAllowList.acceptTermsOfService(OWNER_ADDRESS, OWNER_ADDRESS, r, s, v);
  }
}
