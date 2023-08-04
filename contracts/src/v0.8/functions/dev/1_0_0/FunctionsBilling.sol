// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IFunctionsSubscriptions} from "./interfaces/IFunctionsSubscriptions.sol";
import {AggregatorV3Interface} from "../../../interfaces/AggregatorV3Interface.sol";
import {IFunctionsBilling} from "./interfaces/IFunctionsBilling.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";

import {Routable} from "./Routable.sol";
import {FunctionsResponse} from "./libraries/FunctionsResponse.sol";

/**
 * @title Functions Billing contract
 * @notice Contract that calculates payment from users to the nodes of the Decentralized Oracle Network (DON).
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 */
abstract contract FunctionsBilling is Routable, IFunctionsBilling {
  using FunctionsResponse for FunctionsResponse.Commitment;
  using FunctionsResponse for FunctionsResponse.FulfillResult;

  uint32 private constant REASONABLE_GAS_PRICE_CEILING = 1_000_000;
  // ================================================================
  // |                  Request Commitment state                    |
  // ================================================================

  mapping(bytes32 requestId => bytes32 commitmentHash) private s_requestCommitments;

  event CommitmentDeleted(bytes32 requestId);

  // ================================================================
  // |                     Configuration state                      |
  // ================================================================

  Config private s_config;
  event ConfigChanged(Config config);

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
  constructor(address router, bytes memory config, address linkToNativeFeed) Routable(router, config) {
    s_linkToNativeFeed = AggregatorV3Interface(linkToNativeFeed);
  }

  // ================================================================
  // |                        Configuration                         |
  // ================================================================

  // @notice Sets the configuration of the Chainlink Functions billing registry
  // @param config bytes of abi.encoded config data to set the following:
  //  See the content of the Config struct above
  function _updateConfig(bytes memory config) internal override {
    Config memory _config = abi.decode(config, (Config));
    if (_config.fallbackNativePerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(_config.fallbackNativePerUnitLink);
    }
    s_config = _config;
    emit ConfigChanged(_config);
  }

  // @inheritdoc IFunctionsBilling
  function getConfig() external view override returns (Config memory) {
    return s_config;
  }

  // ================================================================
  // |                       Fee Calculation                        |
  // ================================================================

  // @inheritdoc IFunctionsBilling
  function getDONFee(
    bytes memory /* requestData */,
    RequestBilling memory /* billing */
  ) public view override returns (uint80) {
    // NOTE: Optionally, compute additional fee here
    return s_config.donFee;
  }

  // @inheritdoc IFunctionsBilling
  function getAdminFee() public view override returns (uint96) {
    return _getRouter().getConfig().adminFee;
  }

  // @inheritdoc IFunctionsBilling
  function getWeiPerUnitLink() public view returns (uint256) {
    uint32 feedStalenessSeconds = s_config.feedStalenessSeconds;
    bool staleFallback = feedStalenessSeconds > 0;
    (, int256 weiPerUnitLink, , uint256 timestamp, ) = s_linkToNativeFeed.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    if (staleFallback && feedStalenessSeconds < block.timestamp - timestamp) {
      weiPerUnitLink = s_config.fallbackNativePerUnitLink;
    }
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    return uint256(weiPerUnitLink);
  }

  // ================================================================
  // |                       Cost Estimation                        |
  // ================================================================

  // @inheritdoc IFunctionsBilling
  function estimateCost(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 callbackGasLimit,
    uint256 gasPrice
  ) external view override returns (uint96) {
    // Reasonable ceilings to prevent integer overflows
    _getRouter().isValidCallbackGasLimit(subscriptionId, callbackGasLimit);
    if (gasPrice > REASONABLE_GAS_PRICE_CEILING) {
      revert InvalidCalldata();
    }
    uint96 adminFee = getAdminFee();
    uint96 donFee = getDONFee(
      data,
      RequestBilling({
        subscriptionId: subscriptionId,
        client: msg.sender,
        callbackGasLimit: callbackGasLimit,
        expectedGasPrice: gasPrice,
        adminFee: adminFee
      })
    );
    return _calculateCostEstimate(callbackGasLimit, gasPrice, donFee, adminFee);
  }

  // @notice Estimate the cost in Juels of LINK
  // that will be charged to a subscription to fulfill a Functions request
  // Gas Price can be overestimated to account for flucuations between request and response time
  function _calculateCostEstimate(
    uint32 callbackGasLimit,
    uint256 gasPrice,
    uint96 donFee,
    uint96 adminFee
  ) internal view returns (uint96) {
    uint256 executionGas = s_config.gasOverheadBeforeCallback + s_config.gasOverheadAfterCallback + callbackGasLimit;

    uint256 gasPriceWithOverestimation = gasPrice +
      ((gasPrice * s_config.fulfillmentGasPriceOverEstimationBP) / 10_000);
    // @NOTE: Basis Points are 1/100th of 1%, divide by 10_000 to bring back to original units

    // (1e18 juels/link) (wei/gas * gas) / (wei/link) = juels
    uint256 estimatedGasReimbursement = (1e18 * gasPriceWithOverestimation * executionGas) / getWeiPerUnitLink();

    uint256 fees = uint256(donFee) + uint256(adminFee);

    return uint96(estimatedGasReimbursement + fees);
  }

  // ================================================================
  // |                           Billing                            |
  // ================================================================

  // @notice Initiate the billing process for an Functions request
  // @dev Only callable by the Functions Router
  // @param data - Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
  // @param requestDataVersion - Version number of the structure of the request data
  // @param billing - Billing configuration for the request
  // @return commitment - The parameters of the request that must be held consistent at response time
  function _startBilling(
    bytes memory data,
    uint16 requestDataVersion,
    RequestBilling memory billing
  ) internal returns (FunctionsResponse.Commitment memory commitment) {
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
    IFunctionsSubscriptions router = IFunctionsSubscriptions(address(_getRouter()));
    IFunctionsSubscriptions.Subscription memory subscription = router.getSubscription(billing.subscriptionId);
    if ((subscription.balance - subscription.blockedBalance) < estimatedCost) {
      revert InsufficientBalance();
    }

    (, uint64 initiatedRequests, ) = router.getConsumer(billing.client, billing.subscriptionId);

    bytes32 requestId = computeRequestId(address(this), billing.client, billing.subscriptionId, initiatedRequests + 1);

    commitment = FunctionsResponse.Commitment({
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

    return commitment;
  }

  // @notice Generate a keccak hash request ID
  function computeRequestId(
    address don,
    address client,
    uint64 subscriptionId,
    uint64 nonce
  ) private pure returns (bytes32) {
    return keccak256(abi.encode(don, client, subscriptionId, nonce));
  }

  // @notice Finalize billing process for an Functions request by sending a callback to the Client contract and then charging the subscription
  // @param requestId identifier for the request that was generated by the Registry in the beginBilling commitment
  // @param response response data from DON consensus
  // @param err error from DON consensus
  // @return result fulfillment result
  // @dev Only callable by a node that has been approved on the Coordinator
  // @dev simulated offchain to determine if sufficient balance is present to fulfill the request
  function _fulfillAndBill(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    bytes memory onchainMetadata,
    bytes memory /* offchainMetadata TODO: use in getDonFee() for dynamic billing */
  ) internal returns (FunctionsResponse.FulfillResult) {
    FunctionsResponse.Commitment memory commitment = abi.decode(onchainMetadata, (FunctionsResponse.Commitment));

    if (s_requestCommitments[requestId] != keccak256(abi.encode(commitment))) {
      return FunctionsResponse.FulfillResult.INVALID_COMMITMENT;
    }

    if (s_requestCommitments[requestId] == bytes32(0)) {
      return FunctionsResponse.FulfillResult.INVALID_REQUEST_ID;
    }

    // (1e18 juels/link) * (wei/gas) / (wei/link) = juels per gas
    uint256 juelsPerGas = (1e18 * tx.gasprice) / getWeiPerUnitLink();
    // Gas overhead without callback
    uint96 gasOverheadJuels = uint96(
      juelsPerGas * (commitment.gasOverheadBeforeCallback + commitment.gasOverheadAfterCallback)
    );

    // The Functions Router will perform the callback to the client contract
    (FunctionsResponse.FulfillResult resultCode, uint96 callbackCostJuels) = _getRouter().fulfill(
      response,
      err,
      uint96(juelsPerGas),
      gasOverheadJuels + commitment.donFee, // costWithoutFulfillment
      msg.sender,
      commitment
    );

    // The router will only pay the DON on successfully processing the fulfillment
    // In these two fulfillment results the user has been charged
    // Otherwise, the Coordinator should hold on to the request commitment
    if (
      resultCode == FunctionsResponse.FulfillResult.USER_SUCCESS ||
      resultCode == FunctionsResponse.FulfillResult.USER_ERROR
    ) {
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
  // |                       Request Timeout                        |
  // ================================================================

  // @inheritdoc IFunctionsBilling
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
  // |                    Fund withdrawal                           |
  // ================================================================

  // @notice Oracle withdraw LINK earned through fulfilling requests
  // @notice If amount is 0 the full balance will be withdrawn
  // @notice Both signing and transmitting wallets will have a balance to withdraw
  // @param recipient where to send the funds
  // @param amount amount to withdraw
  function oracleWithdraw(address recipient, uint96 amount) external {
    _disperseFeePool();

    if (amount == 0) {
      amount = s_withdrawableTokens[msg.sender];
    } else if (s_withdrawableTokens[msg.sender] < amount) {
      revert InsufficientBalance();
    }
    s_withdrawableTokens[msg.sender] -= amount;
    IFunctionsSubscriptions router = IFunctionsSubscriptions(address(_getRouter()));
    router.oracleWithdraw(recipient, amount);
  }

  // Overriden in FunctionsCoordinator, which has visibility into transmitters
  function _getTransmitters() internal view virtual returns (address[] memory);

  // DON fees are collected into a pool s_feePool
  // When OCR configuration changes, or any oracle withdraws, this must be dispersed
  function _disperseFeePool() internal {
    if (s_feePool == 0) {
      return;
    }
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
}
