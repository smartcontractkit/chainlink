// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {HasRouter} from "./HasRouter.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IFunctionsSubscriptions} from "./interfaces/IFunctionsSubscriptions.sol";
import {IFunctionsRequest} from "./interfaces/IFunctionsRequest.sol";
import {AggregatorV3Interface} from "../../../interfaces/AggregatorV3Interface.sol";
import {IFunctionsBilling} from "./interfaces/IFunctionsBilling.sol";
import {FulfillResult} from "./interfaces/FulfillResultCodes.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";

/**
 * @title Functions Billing contract
 * @notice Contract that calculates payment from users to the nodes of the Decentralized Oracle Network (DON).
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 */
abstract contract FunctionsBilling is HasRouter, IFunctionsBilling {
  // ================================================================
  // |                  Request Commitment state                    |
  // ================================================================

  mapping(bytes32 requestId => bytes32 commitmentHash) private s_requestCommitments;

  event CommitmentDeleted(bytes32 requestId);

  // ================================================================
  // |                     Configuration state                      |
  // ================================================================

  struct Config {
    // Maximum amount of gas that can be given to a request's client callback
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
    // Max value is 2^80 - 1 == 1.2m LINK.
    uint80 donFee;
    // The highest support request data version supported by the node
    // All lower versions should also be supported
    uint16 maxSupportedRequestDataVersion;
    // Percentage of gas price overestimation to account for changes in gas price between request and response
    // Held as basis points (one hundredth of 1 percentage point)
    uint256 fulfillmentGasPriceOverEstimationBP;
    // fallback NATIVE CURRENCY / LINK conversion rate if the data feed is stale
    int256 fallbackNativePerUnitLink;
  }

  Config private s_config;
  event ConfigChanged(
    uint32 maxCallbackGasLimit,
    uint32 feedStalenessSeconds,
    uint32 gasOverheadBeforeCallback,
    uint32 gasOverheadAfterCallback,
    uint32 requestTimeoutSeconds,
    uint80 donFee,
    uint16 maxSupportedRequestDataVersion,
    uint256 fulfillmentGasPriceOverEstimationBP,
    int256 fallbackNativePerUnitLink
  );

  error UnsupportedRequestDataVersion();
  error InsufficientBalance();
  error InvalidSubscription();
  error UnauthorizedSender();
  error MustBeSubOwner(address owner);
  error InvalidLinkWeiPrice(int256 linkWei);
  error PaymentTooLarge();
  error NoTransmittersSet();
  error InvalidCalldata();

  // ================================================================
  // |                        Balance state                         |
  // ================================================================

  mapping(address transmitter => uint96 balanceJuelsLink) private s_withdrawableTokens;
  // Pool together DON fees and disperse them on withdrawal
  uint96 internal s_feePool;

  AggregatorV3Interface private s_linkToNativeFeed;

  // ================================================================
  // |                       Initialization                         |
  // ================================================================
  constructor(address router, bytes memory config, address linkToNativeFeed) HasRouter(router, config) {
    s_linkToNativeFeed = AggregatorV3Interface(linkToNativeFeed);
  }

  // ================================================================
  // |                    Configuration Methods                     |
  // ================================================================
  /**
   * @notice Sets the configuration of the Chainlink Functions billing registry
   * @param config bytes of abi.encoded config data to set the following:
   *  See the content of the Config struct above
   */
  function _updateConfig(bytes memory config) internal override {
    (
      uint32 maxCallbackGasLimit,
      uint32 feedStalenessSeconds,
      uint32 gasOverheadBeforeCallback,
      uint32 gasOverheadAfterCallback,
      uint32 requestTimeoutSeconds,
      uint80 donFee,
      uint16 maxSupportedRequestDataVersion,
      uint256 fulfillmentGasPriceOverEstimationBP,
      int256 fallbackNativePerUnitLink
    ) = abi.decode(config, (uint32, uint32, uint32, uint32, uint32, uint80, uint16, uint256, int256));

    if (fallbackNativePerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(fallbackNativePerUnitLink);
    }
    s_config = Config({
      maxCallbackGasLimit: maxCallbackGasLimit,
      feedStalenessSeconds: feedStalenessSeconds,
      gasOverheadBeforeCallback: gasOverheadBeforeCallback,
      gasOverheadAfterCallback: gasOverheadAfterCallback,
      requestTimeoutSeconds: requestTimeoutSeconds,
      donFee: donFee,
      maxSupportedRequestDataVersion: maxSupportedRequestDataVersion,
      fulfillmentGasPriceOverEstimationBP: fulfillmentGasPriceOverEstimationBP,
      fallbackNativePerUnitLink: fallbackNativePerUnitLink
    });
    emit ConfigChanged(
      maxCallbackGasLimit,
      feedStalenessSeconds,
      gasOverheadBeforeCallback,
      gasOverheadAfterCallback,
      requestTimeoutSeconds,
      donFee,
      maxSupportedRequestDataVersion,
      fulfillmentGasPriceOverEstimationBP,
      fallbackNativePerUnitLink
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
      uint32 gasOverheadBeforeCallback,
      uint32 gasOverheadAfterCallback,
      uint32 requestTimeoutSeconds,
      uint80 donFee,
      uint16 maxSupportedRequestDataVersion,
      uint256 fulfillmentGasPriceOverEstimationBP,
      int256 fallbackNativePerUnitLink,
      address linkPriceFeed
    )
  {
    return (
      s_config.maxCallbackGasLimit,
      s_config.feedStalenessSeconds,
      s_config.gasOverheadBeforeCallback,
      s_config.gasOverheadAfterCallback,
      s_config.requestTimeoutSeconds,
      s_config.donFee,
      s_config.maxSupportedRequestDataVersion,
      s_config.fulfillmentGasPriceOverEstimationBP,
      s_config.fallbackNativePerUnitLink,
      address(s_linkToNativeFeed)
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
  ) public view override returns (uint80) {
    // NOTE: Optionally, compute additional fee here
    return s_config.donFee;
  }

  /**
   * @inheritdoc IFunctionsBilling
   */
  function getAdminFee() public view override returns (uint96) {
    (, uint96 adminFee, , ) = _getRouter().getConfig();
    return adminFee;
  }

  function getFeedData() public view returns (int256) {
    uint32 feedStalenessSeconds = s_config.feedStalenessSeconds;
    bool staleFallback = feedStalenessSeconds > 0;
    (, int256 weiPerUnitLink, , uint256 timestamp, ) = s_linkToNativeFeed.latestRoundData();
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
    // Reasonable ceilings to prevent integer overflows
    _getRouter().isValidCallbackGasLimit(subscriptionId, callbackGasLimit);
    if (gasPrice > 1_000_000) {
      revert InvalidCalldata();
    }
    uint96 adminFee = getAdminFee();
    RequestBilling memory billing = RequestBilling(subscriptionId, msg.sender, callbackGasLimit, gasPrice, adminFee);
    uint96 donFee = getDONFee(data, billing);
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
    uint256 gasPriceWithOverestimation = gasPrice +
      ((gasPrice * s_config.fulfillmentGasPriceOverEstimationBP) / 10_000);
    // (1e18 juels/link) (wei/gas * gas) / (wei/link) = juels
    uint256 paymentNoFee = (1e18 * gasPriceWithOverestimation * executionGas) / uint256(weiPerUnitLink);
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
   * @return commitment - The parameters of the request that must be held consistent at response time
   */
  function _startBilling(
    bytes memory data,
    uint16 requestDataVersion,
    RequestBilling memory billing
  ) internal returns (IFunctionsRequest.Commitment memory commitment) {
    // Nodes should support all past versions of the structure
    if (requestDataVersion > s_config.maxSupportedRequestDataVersion) {
      revert UnsupportedRequestDataVersion();
    }

    // Check that subscription can afford the estimated cost
    uint80 donFee = getDONFee(data, billing);
    uint96 estimatedCost = _calculateCostEstimate(
      billing.callbackGasLimit,
      billing.expectedGasPrice,
      donFee,
      billing.adminFee
    );
    IFunctionsSubscriptions subscriptions = IFunctionsSubscriptions(address(_getRouter()));
    IFunctionsSubscriptions.Subscription memory subscription = subscriptions.getSubscription(billing.subscriptionId);
    (, uint64 initiatedRequests, ) = subscriptions.getConsumer(billing.client, billing.subscriptionId);

    if (subscription.balance - subscription.blockedBalance < estimatedCost) {
      revert InsufficientBalance();
    }

    bytes32 requestId = computeRequestId(address(this), billing.client, billing.subscriptionId, initiatedRequests + 1);

    commitment = IFunctionsRequest.Commitment({
      adminFee: billing.adminFee,
      coordinator: address(this),
      client: billing.client,
      subscriptionId: billing.subscriptionId,
      callbackGasLimit: billing.callbackGasLimit,
      estimatedTotalCostJuels: estimatedCost,
      timeoutTimestamp: uint40(block.timestamp + s_config.requestTimeoutSeconds),
      requestId: requestId,
      donFee: donFee,
      gasOverheadBeforeCallback: s_config.gasOverheadBeforeCallback,
      gasOverheadAfterCallback: s_config.gasOverheadAfterCallback
    });

    s_requestCommitments[requestId] = keccak256(abi.encode(commitment));
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
    bytes memory err,
    bytes memory onchainMetadata,
    bytes memory /* offchainMetadata TODO: use in getDonFee() for dynamic billing */
  ) internal returns (FulfillResult) {
    IFunctionsRequest.Commitment memory commitment = abi.decode(onchainMetadata, (IFunctionsRequest.Commitment));
    if (s_requestCommitments[requestId] == bytes32(0)) {
      return FulfillResult.INVALID_REQUEST_ID;
    }

    if (s_requestCommitments[requestId] != keccak256(abi.encode(commitment))) {
      return FulfillResult.INVALID_COMMITMENT;
    }

    int256 weiPerUnitLink;
    weiPerUnitLink = getFeedData();
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    // (1e18 juels/link) * (wei/gas) / (wei/link) = juels per gas
    uint256 juelsPerGas = (1e18 * tx.gasprice) / uint256(weiPerUnitLink);
    // Gas overhead without callback
    uint96 gasOverheadJuels = uint96(
      juelsPerGas * (commitment.gasOverheadBeforeCallback + commitment.gasOverheadAfterCallback)
    );

    // The Functions Router will perform the callback to the client contract
    (FulfillResult resultCode, uint96 callbackCostJuels) = _getRouter().fulfill(
      response,
      err,
      uint96(juelsPerGas),
      gasOverheadJuels + commitment.donFee, // costWithoutFulfillment
      msg.sender,
      commitment
    );

    if (resultCode == FulfillResult.USER_SUCCESS || resultCode == FulfillResult.USER_ERROR) {
      delete s_requestCommitments[requestId];
      // Reimburse the transmitter for the fulfillment gas cost
      s_withdrawableTokens[msg.sender] = gasOverheadJuels + callbackCostJuels;
      // Put donFee into the pool of fees, to be split later
      // Saves on storage writes that would otherwise be charged to the user
      s_feePool += commitment.donFee;
    }

    return resultCode;
  }

  // ================================================================
  // |                  Request Timeout Methods                     |
  // ================================================================
  /**
   * @inheritdoc IFunctionsBilling
   */
  function deleteCommitment(bytes32 requestId) external override onlyRouter returns (bool) {
    // Ensure that commitment exists
    if (s_requestCommitments[requestId] == bytes32(0)) {
      return false;
    }
    // Delete commitment
    delete s_requestCommitments[requestId];
    emit CommitmentDeleted(requestId);
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
    IFunctionsSubscriptions router = IFunctionsSubscriptions(address(_getRouter()));
    router.oracleWithdraw(recipient, amount);
  }

  function _disperseFeePool() internal {
    // All transmitters are assumed to also be observers
    // Pay out the DON fee to all transmitters
    address[] memory transmitters = _getTransmitters();
    if (transmitters.length == 0) {
      revert NoTransmittersSet();
    }
    uint96 feePoolShare = s_feePool / uint96(transmitters.length);
    // Bounded by "maxNumOracles" on OCR2Abstract.sol
    for (uint256 i = 0; i < transmitters.length; ++i) {
      s_withdrawableTokens[transmitters[i]] += feePoolShare;
    }
    s_feePool -= feePoolShare * uint96(transmitters.length);
  }

  function _getTransmitters() internal view virtual returns (address[] memory);
}
