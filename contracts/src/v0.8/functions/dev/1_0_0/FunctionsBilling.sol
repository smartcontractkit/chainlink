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
    uint96 estimatedTotalCostJuels;
    uint256 fulfillmentGas;
    uint256 timestamp;
  }
  mapping(bytes32 => Commitment) /* requestID */ /* Commitment */
    private s_requestCommitments;
  struct ItemizedBill {
    uint96 juelsPerGas;
    uint96 signerPayment;
    uint96 transmitterPayment;
  }
  struct Payments {
    address[] to;
    uint96[] amount;
  }
  event RequestTimedOut(bytes32 indexed requestId);

  // ================================================================
  // |                     Configuration state                      |
  // ================================================================
  struct Config {
    // Maxiumum amount of gas that can be given to a request's client callback
    uint32 maxGasLimit;
    // stalenessSeconds is how long before we consider the feed price to be stale
    // and fallback to fallbackNativePerUnitLink.
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
    // fallback NATIVE CURRENCY / LINK conversion rate if the data feed is stale
    int256 fallbackNativePerUnitLink;
  }
  Config private s_config;
  event ConfigSet(
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint256 gasAfterPaymentCalculation,
    int256 fallbackNativePerUnitLink,
    uint32 gasOverhead,
    uint96 fee
  );

  error InsufficientBalance();
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
   *  - fallbackNativePerUnitLink: fallback eth/link price in the case of a stale feed
   *  - gasOverhead: average gas execution cost used in estimating total cost
   *  - requestTimeoutSeconds: e2e timeout after which user won't be charged
   */
  function _setConfig(bytes memory config) internal override {
    (
      uint32 maxGasLimit,
      uint32 stalenessSeconds,
      uint256 gasAfterPaymentCalculation,
      int256 fallbackNativePerUnitLink,
      uint32 gasOverhead,
      uint32 requestTimeoutSeconds,
      uint96 donFee
    ) = abi.decode(config, (uint32, uint32, uint256, int256, uint32, uint32, uint96));

    if (fallbackNativePerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(fallbackNativePerUnitLink);
    }
    s_config = Config({
      maxGasLimit: maxGasLimit,
      stalenessSeconds: stalenessSeconds,
      gasAfterPaymentCalculation: gasAfterPaymentCalculation,
      gasOverhead: gasOverhead,
      requestTimeoutSeconds: requestTimeoutSeconds,
      donFee: donFee,
      fallbackNativePerUnitLink: fallbackNativePerUnitLink
    });
    emit ConfigSet(
      maxGasLimit,
      stalenessSeconds,
      gasAfterPaymentCalculation,
      fallbackNativePerUnitLink,
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
      int256 fallbackNativePerUnitLink,
      uint32 gasOverhead,
      address linkPriceFeed
    )
  {
    return (
      s_config.maxGasLimit,
      s_config.stalenessSeconds,
      s_config.gasAfterPaymentCalculation,
      s_config.fallbackNativePerUnitLink,
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

  function getFeedData() public view returns (int256) {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    (, int256 weiPerUnitLink, , uint256 timestamp, ) = LINK_TO_NATIVE_FEED.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
      weiPerUnitLink = s_config.fallbackNativePerUnitLink;
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
    return _calculateCostEstimate(gasLimit, gasPrice, donFee, adminFee);
  }

  /**
   * @notice Uses current price feed data to estimate a cost
   */
  function _calculateCostEstimate(
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
   * @return estimatedCost
   * @return gasAfterPaymentCalculation
   * @return requestTimeoutSeconds
   * @dev Only callable by the Functions Router
   */
  function _startBilling(bytes memory data, RequestBilling memory billing)
    internal
    returns (
      bytes32,
      uint96,
      uint256,
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
    uint96 estimatedCost = _calculateCostEstimate(billing.gasLimit, billing.gasPrice, donFee, adminFee);
    IFunctionsSubscriptions subscriptions = IFunctionsSubscriptions(address(s_router));
    (uint96 balance, uint96 blockedBalance, , , ) = subscriptions.getSubscription(billing.subscriptionId);
    (, uint64 initiatedRequests, ) = subscriptions.getConsumer(billing.client, billing.subscriptionId);

    if (balance - blockedBalance < estimatedCost) {
      revert InsufficientBalance();
    }

    bytes32 requestId = computeRequestId(address(this), billing.client, billing.subscriptionId, initiatedRequests + 1);

    Commitment memory commitment = Commitment(
      billing.subscriptionId,
      billing.client,
      billing.gasLimit,
      billing.gasPrice,
      address(this),
      donFee,
      adminFee,
      estimatedCost,
      s_config.gasOverhead + s_config.gasAfterPaymentCalculation,
      block.timestamp
    );
    s_requestCommitments[requestId] = commitment;

    return (requestId, estimatedCost, s_config.gasAfterPaymentCalculation, s_config.requestTimeoutSeconds);
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
   * @return result fulfillment result
   * @dev Only callable by a node that has been approved on the Coordinator
   * @dev simulated offchain to determine if sufficient balance is present to fulfill the request
   */
  function _fulfillAndBill(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    /* bytes calldata metadata, */
    address transmitter,
    address[31] memory signers,
    uint8 signerCount
  ) internal returns (IFunctionsRouter.FulfillResult) {
    Commitment memory commitment = s_requestCommitments[requestId];
    if (commitment.don == address(0)) {
      return IFunctionsRouter.FulfillResult.INVALID_REQUEST_ID;
    }
    delete s_requestCommitments[requestId];

    ItemizedBill memory bill = calculatePaymentAmount(
      s_config.gasOverhead,
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
    // Reimburse the transmitter for the execution gas cost + pay them their portion of the DON fee
    payments.to[signerCount + 1] = transmitter;
    payments.amount[signerCount + 1] = bill.transmitterPayment;
    // Remove blocked balance and mark the request as complete
    IFunctionsRouter router = IFunctionsRouter(address(s_router));
    IFunctionsRouter.FulfillResult result = router.fulfill(
      requestId,
      response,
      err,
      bill.juelsPerGas,
      transmitter,
      payments.to,
      payments.amount
    );
    return result;
  }

  /**
   * @notice Determine the cost breakdown for payment
   */
  function calculatePaymentAmount(
    uint256 gasOverhead,
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
    // (1e18 juels/link) * (gas/wei) / (wei/link) = juels per wei
    uint256 juelsPerGas = (1e18 * weiPerUnitGas) / uint256(weiPerUnitLink);
    // Gas overhead without callback
    uint256 paymentNoFee = juelsPerGas * (gasOverhead + gasAfterPaymentCalculation);
    uint256 fee = uint256(donFee) + uint256(adminFee);
    if (paymentNoFee > (1e27 - fee)) {
      revert PaymentTooLarge(); // Payment + fee cannot be more than all of the link in existence.
    }
    uint96 signerPayment = donFee / uint96(signerCount);
    uint96 transmitterPayment = uint96(paymentNoFee);
    return ItemizedBill(uint96(juelsPerGas), signerPayment, transmitterPayment);
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
