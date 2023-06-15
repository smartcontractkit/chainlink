// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {Route, ITypeAndVersion} from "./Route.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IFunctionsSubscriptions} from "./interfaces/IFunctionsSubscriptions.sol";
import {LinkTokenInterface} from "../../../interfaces/LinkTokenInterface.sol";
import {AggregatorV3Interface} from "../../../interfaces/AggregatorV3Interface.sol";
import {IFunctionsBilling} from "./interfaces/IFunctionsBilling.sol";
import {IFunctionsClient} from "./interfaces/IFunctionsClient.sol";
import {ERC677ReceiverInterface} from "../../../interfaces/ERC677ReceiverInterface.sol";
import {IAuthorizedOriginReceiver} from "./accessControl/interfaces/IAuthorizedOriginReceiver.sol";
import {SafeCast} from "../../../shared/vendor/openzeppelin-solidity/v.4.8.0/contracts/utils/SafeCast.sol";

/**
 * @title Functions Billing contract
 * @notice Contract that calculates payment from users to the nodes of the Decentralized Oracle Network (DON).
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 */
abstract contract FunctionsBilling is Route, IFunctionsBilling {
  AggregatorV3Interface private LINK_TO_NATIVE_FEED;

  // ================================================================
  // |                  Request Commitment state                    |
  // ================================================================
  struct Commitment {
    uint64 subscriptionId;
    address client;
    uint32 gasLimit;
    uint256 gasPrice;
    address don;
    uint96 donFee;
    uint96 adminFee;
    uint96 estimatedCost;
    uint256 timestamp;
  }
  mapping(bytes32 => Commitment) /* requestID */ /* Commitment */
    private s_requestCommitments;
  event BillingStart(bytes32 indexed requestId, Commitment commitment);
  struct ItemizedBill {
    uint96 signerPayment;
    uint96 transmitterPayment;
    uint96 totalCost;
  }
  struct Payments {
    address[] to;
    uint96[] amount;
  }
  event BillingEnd(
    bytes32 indexed requestId,
    uint64 subscriptionId,
    uint96 signerPayment,
    uint96 transmitterPayment,
    uint96 totalCost,
    bool success
  );
  event RequestTimedOut(bytes32 indexed requestId);

  // ================================================================
  // |                     Configuration state                      |
  // ================================================================
  struct Config {
    // Maxiumum amount of gas that can be given to a request's client callback
    uint32 maxGasLimit;
    // stalenessSeconds is how long before we consider the feed price to be stale
    // and fallback to fallbackWeiPerUnitLink.
    uint32 stalenessSeconds;
    // Gas to cover transmitter oracle payment after we calculate the payment.
    // We make it configurable in case those operations are repriced.
    uint256 gasAfterPaymentCalculation;
    // Represents the average gas execution cost. Used in estimating cost beforehand.
    uint32 gasOverhead;
    // how many seconds it takes before we consider a request to be timed out
    uint32 requestTimeoutSeconds;
    // additional flat fee (in Juels of LINK) that will be split between node operators
    uint96 donFee;
  }
  int256 private s_fallbackWeiPerUnitLink;
  Config private s_config;
  event ConfigSet(
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint256 gasAfterPaymentCalculation,
    int256 fallbackWeiPerUnitLink,
    uint32 gasOverhead,
    uint96 fee
  );

  error InsufficientBalance();
  error InsufficientGasProvided();
  error InvalidSubscription();
  error UnauthorizedSender();
  error MustBeSubOwner(address owner);
  error GasLimitTooBig(uint32 have, uint32 want);
  error InvalidLinkWeiPrice(int256 linkWei);
  error PaymentTooLarge();

  // ================================================================
  // |                       Initialization                         |
  // ================================================================
  constructor(
    address router,
    bytes memory config,
    address linkToNativeFeed
  ) Route(router, config) {
    LINK_TO_NATIVE_FEED = AggregatorV3Interface(linkToNativeFeed);
  }

  // ================================================================
  // |                    Configuration Methods                     |
  // ================================================================
  /**
   * @notice Sets the configuration of the Chainlink Functions billing registry
   * @param config bytes of config data to set the following:
   *  - maxGasLimit: global max for request gas limit
   *  - stalenessSeconds: if the eth/link feed is more stale then this, use the fallback price
   *  - gasAfterPaymentCalculation: gas used in doing accounting after completing the gas measurement
   *  - fallbackWeiPerUnitLink: fallback eth/link price in the case of a stale feed
   *  - gasOverhead: average gas execution cost used in estimating total cost
   *  - requestTimeoutSeconds: e2e timeout after which user won't be charged
   */
  function _setConfig(bytes memory config) internal override {
    (
      uint32 maxGasLimit,
      uint32 stalenessSeconds,
      uint256 gasAfterPaymentCalculation,
      int256 fallbackWeiPerUnitLink,
      uint32 gasOverhead,
      uint32 requestTimeoutSeconds,
      uint96 donFee
    ) = abi.decode(config, (uint32, uint32, uint256, int256, uint32, uint32, uint96));

    if (fallbackWeiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(fallbackWeiPerUnitLink);
    }
    s_config = Config({
      maxGasLimit: maxGasLimit,
      stalenessSeconds: stalenessSeconds,
      gasAfterPaymentCalculation: gasAfterPaymentCalculation,
      gasOverhead: gasOverhead,
      requestTimeoutSeconds: requestTimeoutSeconds,
      donFee: donFee
    });
    s_fallbackWeiPerUnitLink = fallbackWeiPerUnitLink;
    emit ConfigSet(
      maxGasLimit,
      stalenessSeconds,
      gasAfterPaymentCalculation,
      fallbackWeiPerUnitLink,
      gasOverhead,
      donFee
    );
  }

  /**
   * @inheritdoc IFunctionsBilling
   */
  function getConfig()
    external
    view
    returns (
      uint32 maxGasLimit,
      uint32 stalenessSeconds,
      uint256 gasAfterPaymentCalculation,
      int256 fallbackWeiPerUnitLink,
      uint32 gasOverhead,
      address linkPriceFeed
    )
  {
    return (
      s_config.maxGasLimit,
      s_config.stalenessSeconds,
      s_config.gasAfterPaymentCalculation,
      s_fallbackWeiPerUnitLink,
      s_config.gasOverhead,
      address(LINK_TO_NATIVE_FEED)
    );
  }

  // ================================================================
  // |                 Cost Calculation Methods                     |
  // ================================================================
  /**
   * @inheritdoc IFunctionsBilling
   */
  function getDONFee(
    bytes memory, /* data */
    RequestBilling memory /* billing */
  ) public view override returns (uint96) {
    // NOTE: Optionally, compute additional fee here
    return s_config.donFee;
  }

  /**
   * @inheritdoc IFunctionsBilling
   */
  function getAdminFee(
    bytes memory, /* data */
    RequestBilling memory /* billing */
  ) public view override returns (uint96) {
    // NOTE: Optionally, compute additional fee here
    return IFunctionsRouter(address(s_router)).getAdminFee();
  }

  function getFeedData() private view returns (int256) {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    (, int256 weiPerUnitLink, , uint256 timestamp, ) = LINK_TO_NATIVE_FEED.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
      weiPerUnitLink = s_fallbackWeiPerUnitLink;
    }
    return weiPerUnitLink;
  }

  // ================================================================
  // |                  Cost Estimation Methods                     |
  // ================================================================
  /**
   * @inheritdoc IFunctionsBilling
   */
  function estimateCost(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 gasLimit,
    uint256 gasPrice
  ) external view override returns (uint96) {
    RequestBilling memory billing = RequestBilling(subscriptionId, msg.sender, gasLimit, gasPrice);
    uint96 donFee = getDONFee(data, billing);
    uint96 adminFee = getAdminFee(data, billing);
    return calculateCostEstimate(gasLimit, gasPrice, donFee, adminFee);
  }

  /**
   * @notice Uses current price feed data to estimate a cost
   */
  function calculateCostEstimate(
    uint32 gasLimit,
    uint256 gasPrice,
    uint96 donFee,
    uint96 adminFee
  ) internal view returns (uint96) {
    int256 weiPerUnitLink;
    weiPerUnitLink = getFeedData();
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    uint256 executionGas = s_config.gasOverhead + s_config.gasAfterPaymentCalculation + gasLimit;
    // (1e18 juels/link) (wei/gas * gas) / (wei/link) = juels
    uint256 paymentNoFee = (1e18 * gasPrice * executionGas) / uint256(weiPerUnitLink);
    uint256 fee = uint256(donFee) + uint256(adminFee);
    if (paymentNoFee > (1e27 - fee)) {
      revert PaymentTooLarge(); // Payment + fee cannot be more than all of the link in existence.
    }
    return uint96(paymentNoFee + fee);
  }

  // ================================================================
  // |                       Billing Methods                        |
  // ================================================================
  /**
   * @notice Initiate the billing process for an Functions request
   * @param data Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
   * @param billing Billing configuration for the request
   * @return requestId - A unique identifier of the request. Can be used to match a request to a response in fulfillRequest.
   * @dev Only callable by the Functions Router
   */
  function startBilling(bytes memory data, RequestBilling memory billing)
    internal
    returns (
      bytes32,
      uint96,
      uint256
    )
  {
    // No lower bound on the requested gas limit. A user could request 0
    // and they would simply be billed for the gas and computation.
    if (billing.gasLimit > s_config.maxGasLimit) {
      revert GasLimitTooBig(billing.gasLimit, s_config.maxGasLimit);
    }

    // Check that subscription can afford the estimated cost
    uint96 donFee = getDONFee(data, billing);
    uint96 adminFee = getAdminFee(data, billing);
    uint96 estimatedCost = calculateCostEstimate(billing.gasLimit, billing.gasPrice, donFee, adminFee);
    IFunctionsSubscriptions subscriptions = IFunctionsSubscriptions(address(s_router));
    (uint96 balance, uint96 blockedBalance, , , ) = subscriptions.getSubscription(billing.subscriptionId);
    (, uint64 initiatedRequests, ) = subscriptions.getConsumer(billing.client, billing.subscriptionId);

    uint96 effectiveBalance = balance - blockedBalance;
    if (effectiveBalance < estimatedCost) {
      revert InsufficientBalance();
    }

    bytes32 requestId = computeRequestId(address(this), billing.client, billing.subscriptionId, initiatedRequests + 1);

    Commitment memory commitment = Commitment(
      billing.subscriptionId,
      billing.client,
      billing.gasLimit,
      billing.gasPrice,
      msg.sender,
      donFee,
      adminFee,
      estimatedCost,
      block.timestamp
    );
    s_requestCommitments[requestId] = commitment;

    emit BillingStart(requestId, commitment);
    return (requestId, estimatedCost, s_config.requestTimeoutSeconds);
  }

  /**
   * @notice Generate a keccak hash request ID
   */
  function computeRequestId(
    address don,
    address client,
    uint64 subscriptionId,
    uint64 nonce
  ) private pure returns (bytes32) {
    return keccak256(abi.encode(don, client, subscriptionId, nonce));
  }

  /**
   * @notice Finalize billing process for an Functions request by sending a callback to the Client contract and then charging the subscription
   * @param requestId identifier for the request that was generated by the Registry in the beginBilling commitment
   * @param response response data from DON consensus
   * @param err error from DON consensus
   * @param transmitter the Oracle who sent the report
   * @param signers the Oracles who had a part in generating the report
   * @param signerCount the number of signers on the report
   * @param initialGas the initial amount of gas that should be used as a baseline to charge the single fulfillment for execution cost
   * @return result fulfillment result
   * @dev Only callable by a node that has been approved on the Registry
   * @dev simulated offchain to determine if sufficient balance is present to fulfill the request
   */
  function fulfillAndBill(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    /* bytes calldata metadata, */
    address transmitter,
    address[31] memory signers,
    uint8 signerCount,
    uint256 initialGas
  ) internal returns (FulfillResult) {
    Commitment memory commitment = s_requestCommitments[requestId];
    if (commitment.don == address(0)) {
      return FulfillResult.INVALID_REQUEST_ID;
    }

    if (gasleft() < commiment.gasLimit + s_config.gasAfterPaymentCalculation + s_config.gasOverhead) {
      revert InsufficientGasProvided();
    }

    delete s_requestCommitments[requestId];

    IFunctionsRouter router = IFunctionsRouter(address(s_router));
    bool success = router.callback(requestId, response, err);

    // We want to charge users exactly for how much gas they use in their callback.
    // The gasAfterPaymentCalculation is meant to cover these additional operations where we
    // decrement the subscription balance and increment the oracle's withdrawable balance.
    ItemizedBill memory bill = calculatePaymentAmount(
      initialGas,
      s_config.gasAfterPaymentCalculation,
      commitment.donFee,
      signerCount,
      commitment.adminFee,
      tx.gasprice
    );

    Payments memory payments;
    // Pay out signers their portion of the DON fee
    for (uint8 i = 0; i < signerCount; i++) {
      payments.to[i] = signers[i];
      payments.amount[i] = bill.signerPayment;
    }
    // Pay out the administration fee
    // address routerOwner = router.owner();
    // payments.to[signerCount + 1] = routerOwner;
    payments.amount[signerCount + 1] = commitment.adminFee;
    // Reimburse the transmitter for the execution gas cost + pay them their portion of the DON fee
    payments.to[signerCount + 2] = transmitter;
    payments.amount[signerCount + 2] = bill.transmitterPayment;
    // Remove blocked balance and mark the request as complete
    // IFunctionsSubscriptions subscriptions = IFunctionsSubscriptions(address(s_router));
    IFunctionsSubscriptions subscriptions = IFunctionsSubscriptions(address(s_router));
    subscriptions.pay(requestId, payments.to, payments.amount);

    // Include payment in the event for tracking costs.
    emit BillingEnd(
      requestId,
      commitment.subscriptionId,
      bill.signerPayment,
      bill.transmitterPayment,
      bill.totalCost,
      success
    );

    return success ? FulfillResult.USER_SUCCESS : FulfillResult.USER_ERROR;
  }

  /**
   * @notice Determine the cost breakdown for payment
   */
  function calculatePaymentAmount(
    uint256 startGas,
    uint256 gasAfterPaymentCalculation,
    uint96 donFee,
    uint8 signerCount,
    uint96 adminFee,
    uint256 weiPerUnitGas
  ) private view returns (ItemizedBill memory) {
    int256 weiPerUnitLink;
    weiPerUnitLink = getFeedData();
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    // (1e18 juels/link) (wei/gas * gas) / (wei/link) = juels
    uint256 paymentNoFee = (1e18 * weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft())) /
      uint256(weiPerUnitLink);
    uint256 fee = uint256(donFee) + uint256(adminFee);
    if (paymentNoFee > (1e27 - fee)) {
      revert PaymentTooLarge(); // Payment + fee cannot be more than all of the link in existence.
    }
    uint96 signerPayment = donFee / uint96(signerCount);
    uint96 transmitterPayment = uint96(paymentNoFee);
    uint96 totalCost = SafeCast.toUint96(paymentNoFee + fee);
    return ItemizedBill(signerPayment, transmitterPayment, totalCost);
  }

  // ================================================================
  // |                  Request Timeout Methods                     |
  // ================================================================
  /**
   * @inheritdoc IFunctionsBilling
   */
  function timeoutRequest(bytes32 requestId) external override onlyRouter returns (bool) {
    Commitment memory commitment = s_requestCommitments[requestId];
    // Ensure that commitment exists
    if (commitment.don == address(0)) {
      return false;
    }
    // Delete commitment
    delete s_requestCommitments[requestId];
    emit RequestTimedOut(requestId);
    return true;
  }
}
