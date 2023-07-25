// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Routable} from "./Routable.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IFunctionsSubscriptions} from "./interfaces/IFunctionsSubscriptions.sol";
import {AggregatorV3Interface} from "../../../interfaces/AggregatorV3Interface.sol";
import {IFunctionsBilling} from "./interfaces/IFunctionsBilling.sol";
import {FulfillResult} from "./FulfillResultCodes.sol";
import {SafeCast} from "../../../shared/vendor/openzeppelin-solidity/v.4.8.0/contracts/utils/SafeCast.sol";

/**
 * @title Functions Billing contract
 * @notice Contract that calculates payment from users to the nodes of the Decentralized Oracle Network (DON).
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 */
abstract contract FunctionsBilling is Routable, IFunctionsBilling {
  AggregatorV3Interface private LINK_TO_NATIVE_FEED;

  // ================================================================
  // |                  Request Commitment state                    |
  // ================================================================
  struct Commitment {
    uint64 subscriptionId;
    address client;
    uint32 callbackGasLimit;
    uint256 expectedGasPrice;
    address don;
    uint96 donFee;
    uint96 adminFee;
    uint96 estimatedTotalCostJuels;
    uint256 gasOverhead;
    uint256 timestamp;
  }
  mapping(bytes32 requestId => Commitment) private s_requestCommitments;

  event RequestTimedOut(bytes32 indexed requestId);

  // ================================================================
  // |                     Configuration state                      |
  // ================================================================
  struct Config {
    // Maxiumum amount of gas that can be given to a request's client callback
    uint32 maxCallbackGasLimit;
    // feedStalenessSeconds is how long before we consider the feed price to be stale
    // and fallback to fallbackNativePerUnitLink.
    uint32 feedStalenessSeconds;
    // Represents the average gas execution cost. Used in estimating cost beforehand.
    uint32 gasOverheadBeforeCallback;
    // Gas to cover transmitter oracle payment after we calculate the payment.
    // We make it configurable in case those operations are repriced.
    uint32 gasOverheadAfterCallback;
    // how many seconds it takes before we consider a request to be timed out
    uint32 requestTimeoutSeconds;
    // additional flat fee (in Juels of LINK) that will be split between Node Operators
    uint96 donFee;
    // fallback NATIVE CURRENCY / LINK conversion rate if the data feed is stale
    int256 fallbackNativePerUnitLink;
    // The highest support request data version supported by the node
    // All lower versions should also be supported
    uint16 maxSupportedRequestDataVersion;
  }
  Config private s_config;
  event ConfigSet(
    uint32 maxCallbackGasLimit,
    uint32 feedStalenessSeconds,
    uint32 gasOverheadBeforeCallback,
    uint32 gasOverheadAfterCallback,
    int256 fallbackNativePerUnitLink,
    uint96 fee,
    uint16 maxSupportedRequestDataVersion
  );

  error UnsupportedRequestDataVersion();
  error InsufficientBalance();
  error InvalidSubscription();
  error UnauthorizedSender();
  error MustBeSubOwner(address owner);
  error GasLimitTooBig(uint32 have, uint32 want);
  error InvalidLinkWeiPrice(int256 linkWei);
  error PaymentTooLarge();
  error NoTransmittersSet();

  // ================================================================
  // |                        Balance state                         |
  // ================================================================
  mapping(address transmitter => uint96 balanceJuelsLink) private s_withdrawableTokens;
  // Pool together DON fees and disperse them on withdrawal
  uint96 s_feePool;

  // ================================================================
  // |                         Cost Events                          |
  // ================================================================
  event BillingStart(bytes32 indexed requestId, Commitment commitment);
  event BillingEnd(
    bytes32 indexed requestId,
    uint64 subscriptionId,
    uint96 signerPayment,
    uint96 transmitterPayment,
    uint96 totalCost,
    FulfillResult result
  );

  // ================================================================
  // |                       Initialization                         |
  // ================================================================
  constructor(address router, bytes memory config, address linkToNativeFeed) Routable(router, config) {
    LINK_TO_NATIVE_FEED = AggregatorV3Interface(linkToNativeFeed);
  }

  // ================================================================
  // |                    Configuration Methods                     |
  // ================================================================
  /**
   * @notice Sets the configuration of the Chainlink Functions billing registry
   * @param config bytes of config data to set the following:
   *  - maxCallbackGasLimit: global max for request gas limit
   *  - feedStalenessSeconds: if the eth/link feed is more stale then this, use the fallback price
   *  - gasOverheadAfterCallback: gas used in doing accounting after completing the gas measurement
   *  - fallbackNativePerUnitLink: fallback eth/link price in the case of a stale feed
   *  - gasOverheadBeforeCallback: average gas execution cost used in estimating total cost
   *  - requestTimeoutSeconds: e2e timeout after which user won't be charged
   */
  function _setConfig(bytes memory config) internal override {
    (
      uint32 maxCallbackGasLimit,
      uint32 feedStalenessSeconds,
      uint32 gasOverheadAfterCallback,
      uint32 gasOverheadBeforeCallback,
      int256 fallbackNativePerUnitLink,
      uint32 requestTimeoutSeconds,
      uint96 donFee,
      uint16 maxSupportedRequestDataVersion
    ) = abi.decode(config, (uint32, uint32, uint32, uint32, int256, uint32, uint96, uint16));

    if (fallbackNativePerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(fallbackNativePerUnitLink);
    }
    s_config = Config({
      maxCallbackGasLimit: maxCallbackGasLimit,
      feedStalenessSeconds: feedStalenessSeconds,
      gasOverheadAfterCallback: gasOverheadAfterCallback,
      gasOverheadBeforeCallback: gasOverheadBeforeCallback,
      requestTimeoutSeconds: requestTimeoutSeconds,
      donFee: donFee,
      fallbackNativePerUnitLink: fallbackNativePerUnitLink,
      maxSupportedRequestDataVersion: maxSupportedRequestDataVersion
    });
    emit ConfigSet(
      maxCallbackGasLimit,
      feedStalenessSeconds,
      gasOverheadBeforeCallback,
      gasOverheadAfterCallback,
      fallbackNativePerUnitLink,
      donFee,
      maxSupportedRequestDataVersion
    );
  }

  /**
   * @inheritdoc IFunctionsBilling
   */
  function getConfig()
    external
    view
    override
    returns (
      uint32 maxCallbackGasLimit,
      uint32 feedStalenessSeconds,
      uint256 gasOverheadAfterCallback,
      int256 fallbackNativePerUnitLink,
      uint32 gasOverheadBeforeCallback,
      address linkPriceFeed,
      uint16 maxSupportedRequestDataVersion
    )
  {
    return (
      s_config.maxCallbackGasLimit,
      s_config.feedStalenessSeconds,
      s_config.gasOverheadAfterCallback,
      s_config.fallbackNativePerUnitLink,
      s_config.gasOverheadBeforeCallback,
      address(LINK_TO_NATIVE_FEED),
      s_config.maxSupportedRequestDataVersion
    );
  }

  // ================================================================
  // |                 Cost Calculation Methods                     |
  // ================================================================
  /**
   * @inheritdoc IFunctionsBilling
   */
  function getDONFee(
    bytes memory /* requestData */,
    RequestBilling memory /* billing */
  ) public view override returns (uint96) {
    // NOTE: Optionally, compute additional fee here
    return s_config.donFee;
  }

  /**
   * @inheritdoc IFunctionsBilling
   */
  function getAdminFee(
    bytes memory /* requestData */,
    RequestBilling memory /* billing */
  ) public view override returns (uint96) {
    // NOTE: Optionally, compute additional fee here
    return IFunctionsRouter(address(s_router)).getAdminFee();
  }

  function getFeedData() public view returns (int256) {
    uint32 feedStalenessSeconds = s_config.feedStalenessSeconds;
    bool staleFallback = feedStalenessSeconds > 0;
    (, int256 weiPerUnitLink, , uint256 timestamp, ) = LINK_TO_NATIVE_FEED.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    if (staleFallback && feedStalenessSeconds < block.timestamp - timestamp) {
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
    uint32 callbackGasLimit,
    uint256 gasPrice
  ) external view override returns (uint96) {
    RequestBilling memory billing = RequestBilling(subscriptionId, msg.sender, callbackGasLimit, gasPrice);
    uint96 donFee = getDONFee(data, billing);
    uint96 adminFee = getAdminFee(data, billing);
    return _calculateCostEstimate(callbackGasLimit, gasPrice, donFee, adminFee);
  }

  /**
   * @notice Uses current price feed data to estimate a cost
   */
  function _calculateCostEstimate(
    uint32 callbackGasLimit,
    uint256 gasPrice,
    uint96 donFee,
    uint96 adminFee
  ) internal view returns (uint96) {
    int256 weiPerUnitLink;
    weiPerUnitLink = getFeedData();
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    uint256 executionGas = s_config.gasOverheadBeforeCallback + s_config.gasOverheadAfterCallback + callbackGasLimit;
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
   * @dev Only callable by the Functions Router
   * @param data - Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
   * @param requestDataVersion - Version number of the structure of the request data
   * @param billing - Billing configuration for the request
   * @return requestId - A unique identifier of the request. Can be used to match a request to a response in fulfillRequest.
   * @return estimatedCost - The estimated cost in Juels of LINK that will be charged to the subscription if all callback gas is used
   * @return gasOverheadAfterCallback - The amount of gas that will be used after the user's callback
   * @return requestTimeoutSeconds - The number of seconds that this request can remain unfilled before being considered stale
   */
  function _startBilling(
    bytes memory data,
    uint16 requestDataVersion,
    RequestBilling memory billing
  ) internal returns (bytes32, uint96, uint256, uint256) {
    // Nodes should support all past versions of the structure
    if (requestDataVersion > s_config.maxSupportedRequestDataVersion) {
      revert UnsupportedRequestDataVersion();
    }

    // No lower bound on the requested gas limit. A user could request 0
    // and they would simply be billed for the gas and computation.
    if (billing.callbackGasLimit > s_config.maxCallbackGasLimit) {
      revert GasLimitTooBig(billing.callbackGasLimit, s_config.maxCallbackGasLimit);
    }

    // Check that subscription can afford the estimated cost
    uint96 donFee = getDONFee(data, billing);
    uint96 adminFee = getAdminFee(data, billing);
    uint96 estimatedCost = _calculateCostEstimate(billing.callbackGasLimit, billing.expectedGasPrice, donFee, adminFee);
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
      billing.callbackGasLimit,
      billing.expectedGasPrice,
      address(this),
      donFee,
      adminFee,
      estimatedCost,
      s_config.gasOverheadBeforeCallback + s_config.gasOverheadAfterCallback,
      block.timestamp
    );
    s_requestCommitments[requestId] = commitment;

    emit BillingStart(requestId, commitment);

    return (requestId, estimatedCost, s_config.gasOverheadAfterCallback, s_config.requestTimeoutSeconds);
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
   * @return result fulfillment result
   * @dev Only callable by a node that has been approved on the Coordinator
   * @dev simulated offchain to determine if sufficient balance is present to fulfill the request
   */
  function _fulfillAndBill(
    bytes32 requestId,
    bytes memory response,
    bytes memory err
  )
    internal
    returns (
      /* bytes calldata metadata, */
      FulfillResult
    )
  {
    Commitment memory commitment = s_requestCommitments[requestId];
    if (commitment.don == address(0)) {
      return FulfillResult.INVALID_REQUEST_ID;
    }
    delete s_requestCommitments[requestId];

    int256 weiPerUnitLink;
    weiPerUnitLink = getFeedData();
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    // (1e18 juels/link) * (gas/wei) / (wei/link) = juels per wei
    uint256 juelsPerGas = (1e18 * tx.gasprice) / uint256(weiPerUnitLink);
    // Gas overhead without callback
    uint96 gasOverheadJuels = uint96(juelsPerGas * commitment.gasOverhead);
    uint96 costWithoutFulfillment = gasOverheadJuels + commitment.donFee;

    // The Functions Router will perform the callback to the client contract
    IFunctionsRouter router = IFunctionsRouter(address(s_router));
    (uint8 result, uint96 callbackCostJuels) = router.fulfill(
      requestId,
      response,
      err,
      uint96(juelsPerGas),
      costWithoutFulfillment,
      msg.sender
    );

    // Reimburse the transmitter for the fulfillment gas cost
    s_withdrawableTokens[msg.sender] = gasOverheadJuels + callbackCostJuels;
    // Put donFee into the pool of fees, to be split later
    // Saves on storage writes that would otherwise be charged to the user
    s_feePool += commitment.donFee;

    emit BillingEnd(
      requestId,
      commitment.subscriptionId,
      commitment.donFee,
      gasOverheadJuels + callbackCostJuels,
      gasOverheadJuels + callbackCostJuels + commitment.donFee + commitment.adminFee,
      FulfillResult(result)
    );

    return FulfillResult(result);
  }

  // ================================================================
  // |                  Request Timeout Methods                     |
  // ================================================================
  /**
   * @inheritdoc IFunctionsBilling
   */
  function deleteCommitment(bytes32 requestId) external override onlyRouter returns (bool) {
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

  // ================================================================
  // |                   Node Operator methods                      |
  // ================================================================
  /*
   * @notice Oracle withdraw LINK earned through fulfilling requests
   * @notice If amount is 0 the full balance will be withdrawn
   * @notice Both signing and transmitting wallets will have a balance to withdraw
   * @param recipient where to send the funds
   * @param amount amount to withdraw
   */
  function oracleWithdraw(address recipient, uint96 amount) external {
    _disperseFeePool();

    if (amount == 0) {
      amount = s_withdrawableTokens[msg.sender];
    }
    if (s_withdrawableTokens[msg.sender] < amount) {
      revert InsufficientBalance();
    }
    s_withdrawableTokens[msg.sender] -= amount;
    IFunctionsSubscriptions router = IFunctionsSubscriptions(address(s_router));
    router.oracleWithdraw(recipient, amount);
  }

  function _disperseFeePool() internal {
    // All transmitters are assumed to also be observers
    // Pay out the DON fee to all transmitters
    // Bounded by "maxNumOracles" on OCR2Abstract.sol
    address[] memory transmitters = _getTransmitters();
    if (transmitters.length == 0) {
      revert NoTransmittersSet();
    }
    uint96 feePoolShare = s_feePool / uint96(transmitters.length);
    for (uint8 i = 0; i < transmitters.length; i++) {
      s_withdrawableTokens[transmitters[i]] += feePoolShare;
    }
    s_feePool -= feePoolShare * uint96(transmitters.length);
  }

  function _getTransmitters() internal view virtual returns (address[] memory);
}
