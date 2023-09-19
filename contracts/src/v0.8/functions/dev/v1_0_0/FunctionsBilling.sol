// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IFunctionsSubscriptions} from "./interfaces/IFunctionsSubscriptions.sol";
import {AggregatorV3Interface} from "../../../interfaces/AggregatorV3Interface.sol";
import {IFunctionsBilling} from "./interfaces/IFunctionsBilling.sol";

import {Routable} from "./Routable.sol";
import {FunctionsResponse} from "./libraries/FunctionsResponse.sol";

import {SafeCast} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/math/SafeCast.sol";

/// @title Functions Billing contract
/// @notice Contract that calculates payment from users to the nodes of the Decentralized Oracle Network (DON).
/// @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
abstract contract FunctionsBilling is Routable, IFunctionsBilling {
  using FunctionsResponse for FunctionsResponse.RequestMeta;
  using FunctionsResponse for FunctionsResponse.Commitment;
  using FunctionsResponse for FunctionsResponse.FulfillResult;

  uint256 private constant REASONABLE_GAS_PRICE_CEILING = 1_000_000_000_000_000; // 1 million gwei
  // ================================================================
  // |                  Request Commitment state                    |
  // ================================================================

  mapping(bytes32 requestId => bytes32 commitmentHash) private s_requestCommitments;

  event CommitmentDeleted(bytes32 requestId);

  // ================================================================
  // |                     Configuration state                      |
  // ================================================================

  struct Config {
    uint32 fulfillmentGasPriceOverEstimationBP; // ══╗ Percentage of gas price overestimation to account for changes in gas price between request and response. Held as basis points (one hundredth of 1 percentage point)
    uint32 feedStalenessSeconds; //                  ║ How long before we consider the feed price to be stale and fallback to fallbackNativePerUnitLink.
    uint32 gasOverheadBeforeCallback; //             ║ Represents the average gas execution cost before the fulfillment callback. This amount is always billed for every request.
    uint32 gasOverheadAfterCallback; //              ║ Represents the average gas execution cost after the fulfillment callback. This amount is always billed for every request.
    uint32 requestTimeoutSeconds; //                 ║ How many seconds it takes before we consider a request to be timed out
    uint72 donFee; //                                ║ Additional flat fee (in Juels of LINK) that will be split between Node Operators. Max value is 2^80 - 1 == 1.2m LINK.
    uint16 maxSupportedRequestDataVersion; // ═══════╝ The highest support request data version supported by the node. All lower versions should also be supported.
    uint224 fallbackNativePerUnitLink; // ═══════════╸ fallback NATIVE CURRENCY / LINK conversion rate if the data feed is stale
  }

  Config private s_config;

  event ConfigUpdated(Config config);

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
  // Pool together collected DON fees
  // Disperse them on withdrawal or change in OCR configuration
  uint96 internal s_feePool;

  AggregatorV3Interface private s_linkToNativeFeed;

  // ================================================================
  // |                       Initialization                         |
  // ================================================================
  constructor(address router, Config memory config, address linkToNativeFeed) Routable(router) {
    s_linkToNativeFeed = AggregatorV3Interface(linkToNativeFeed);

    updateConfig(config);
  }

  // ================================================================
  // |                        Configuration                         |
  // ================================================================

  /// @notice Gets the Chainlink Coordinator's billing configuration
  /// @return config
  function getConfig() external view returns (Config memory) {
    return s_config;
  }

  /// @notice Sets the Chainlink Coordinator's billing configuration
  /// @param config - See the contents of the Config struct in IFunctionsBilling.Config for more information
  function updateConfig(Config memory config) public {
    _onlyOwner();

    s_config = config;
    emit ConfigUpdated(config);
  }

  // ================================================================
  // |                       Fee Calculation                        |
  // ================================================================

  /// @inheritdoc IFunctionsBilling
  function getDONFee(bytes memory /* requestData */) public view override returns (uint72) {
    return s_config.donFee;
  }

  /// @inheritdoc IFunctionsBilling
  function getAdminFee() public view override returns (uint72) {
    return _getRouter().getAdminFee();
  }

  /// @inheritdoc IFunctionsBilling
  function getWeiPerUnitLink() public view returns (uint256) {
    Config memory config = s_config;
    (, int256 weiPerUnitLink, , uint256 timestamp, ) = s_linkToNativeFeed.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    if (config.feedStalenessSeconds < block.timestamp - timestamp && config.feedStalenessSeconds > 0) {
      return config.fallbackNativePerUnitLink;
    }
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    return uint256(weiPerUnitLink);
  }

  function _getJuelsPerGas(uint256 gasPriceWei) private view returns (uint96) {
    // (1e18 juels/link) * (wei/gas) / (wei/link) = juels per gas
    // There are only 1e9*1e18 = 1e27 juels in existence, should not exceed uint96 (2^96 ~ 7e28)
    return SafeCast.toUint96((1e18 * gasPriceWei) / getWeiPerUnitLink());
  }

  // ================================================================
  // |                       Cost Estimation                        |
  // ================================================================

  /// @inheritdoc IFunctionsBilling
  function estimateCost(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 callbackGasLimit,
    uint256 gasPriceWei
  ) external view override returns (uint96) {
    _getRouter().isValidCallbackGasLimit(subscriptionId, callbackGasLimit);
    // Reasonable ceilings to prevent integer overflows
    if (gasPriceWei > REASONABLE_GAS_PRICE_CEILING) {
      revert InvalidCalldata();
    }
    uint72 adminFee = getAdminFee();
    uint72 donFee = getDONFee(data);
    return _calculateCostEstimate(callbackGasLimit, gasPriceWei, donFee, adminFee);
  }

  /// @notice Estimate the cost in Juels of LINK
  // that will be charged to a subscription to fulfill a Functions request
  // Gas Price can be overestimated to account for flucuations between request and response time
  function _calculateCostEstimate(
    uint32 callbackGasLimit,
    uint256 gasPriceWei,
    uint72 donFee,
    uint72 adminFee
  ) internal view returns (uint96) {
    uint256 executionGas = s_config.gasOverheadBeforeCallback + s_config.gasOverheadAfterCallback + callbackGasLimit;

    uint256 gasPriceWithOverestimation = gasPriceWei +
      ((gasPriceWei * s_config.fulfillmentGasPriceOverEstimationBP) / 10_000);
    /// @NOTE: Basis Points are 1/100th of 1%, divide by 10_000 to bring back to original units

    uint96 juelsPerGas = _getJuelsPerGas(gasPriceWithOverestimation);
    uint256 estimatedGasReimbursement = juelsPerGas * executionGas;
    uint96 fees = uint96(donFee) + uint96(adminFee);

    return SafeCast.toUint96(estimatedGasReimbursement + fees);
  }

  // ================================================================
  // |                           Billing                            |
  // ================================================================

  /// @notice Initiate the billing process for an Functions request
  /// @dev Only callable by the Functions Router
  /// @param request - Chainlink Functions request data, see FunctionsResponse.RequestMeta for the structure
  /// @return commitment - The parameters of the request that must be held consistent at response time
  function _startBilling(
    FunctionsResponse.RequestMeta memory request
  ) internal returns (FunctionsResponse.Commitment memory commitment) {
    Config memory config = s_config;

    // Nodes should support all past versions of the structure
    if (request.dataVersion > config.maxSupportedRequestDataVersion) {
      revert UnsupportedRequestDataVersion();
    }

    uint72 donFee = getDONFee(request.data);
    uint96 estimatedTotalCostJuels = _calculateCostEstimate(
      request.callbackGasLimit,
      tx.gasprice,
      donFee,
      request.adminFee
    );

    // Check that subscription can afford the estimated cost
    if ((request.availableBalance) < estimatedTotalCostJuels) {
      revert InsufficientBalance();
    }

    bytes32 requestId = _computeRequestId(
      address(this),
      request.requestingContract,
      request.subscriptionId,
      request.initiatedRequests + 1
    );

    commitment = FunctionsResponse.Commitment({
      adminFee: request.adminFee,
      coordinator: address(this),
      client: request.requestingContract,
      subscriptionId: request.subscriptionId,
      callbackGasLimit: request.callbackGasLimit,
      estimatedTotalCostJuels: estimatedTotalCostJuels,
      timeoutTimestamp: uint32(block.timestamp + config.requestTimeoutSeconds),
      requestId: requestId,
      donFee: donFee,
      gasOverheadBeforeCallback: config.gasOverheadBeforeCallback,
      gasOverheadAfterCallback: config.gasOverheadAfterCallback
    });

    s_requestCommitments[requestId] = keccak256(abi.encode(commitment));

    return commitment;
  }

  /// @notice Generate a keccak hash request ID
  /// @dev uses the number of requests that the consumer of a subscription has sent as a nonce
  function _computeRequestId(
    address don,
    address client,
    uint64 subscriptionId,
    uint64 nonce
  ) private pure returns (bytes32) {
    return keccak256(abi.encode(don, client, subscriptionId, nonce));
  }

  /// @notice Finalize billing process for an Functions request by sending a callback to the Client contract and then charging the subscription
  /// @param requestId identifier for the request that was generated by the Registry in the beginBilling commitment
  /// @param response response data from DON consensus
  /// @param err error from DON consensus
  /// @return result fulfillment result
  /// @dev Only callable by a node that has been approved on the Coordinator
  /// @dev simulated offchain to determine if sufficient balance is present to fulfill the request
  function _fulfillAndBill(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    bytes memory onchainMetadata,
    bytes memory /* offchainMetadata TODO: use in getDonFee() for dynamic billing */
  ) internal returns (FunctionsResponse.FulfillResult) {
    FunctionsResponse.Commitment memory commitment = abi.decode(onchainMetadata, (FunctionsResponse.Commitment));

    if (s_requestCommitments[requestId] == bytes32(0)) {
      return FunctionsResponse.FulfillResult.INVALID_REQUEST_ID;
    }

    if (s_requestCommitments[requestId] != keccak256(abi.encode(commitment))) {
      return FunctionsResponse.FulfillResult.INVALID_COMMITMENT;
    }

    uint96 juelsPerGas = _getJuelsPerGas(tx.gasprice);
    // Gas overhead without callback
    uint96 gasOverheadJuels = juelsPerGas *
      (commitment.gasOverheadBeforeCallback + commitment.gasOverheadAfterCallback);

    // The Functions Router will perform the callback to the client contract
    (FunctionsResponse.FulfillResult resultCode, uint96 callbackCostJuels) = _getRouter().fulfill(
      response,
      err,
      juelsPerGas,
      gasOverheadJuels + commitment.donFee, // costWithoutFulfillment
      msg.sender,
      commitment
    );

    // The router will only pay the DON on successfully processing the fulfillment
    // In these two fulfillment results the user has been charged
    // Otherwise, the Coordinator should hold on to the request commitment
    if (
      resultCode == FunctionsResponse.FulfillResult.FULFILLED ||
      resultCode == FunctionsResponse.FulfillResult.USER_CALLBACK_ERROR
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

  /// @inheritdoc IFunctionsBilling
  /// @dev Only callable by the Router
  /// @dev Used by FunctionsRouter.sol during timeout of a request
  function deleteCommitment(bytes32 requestId) external override onlyRouter {
    // Delete commitment
    delete s_requestCommitments[requestId];
    emit CommitmentDeleted(requestId);
  }

  // ================================================================
  // |                    Fund withdrawal                           |
  // ================================================================

  /// @inheritdoc IFunctionsBilling
  function oracleWithdraw(address recipient, uint96 amount) external {
    _disperseFeePool();

    if (amount == 0) {
      amount = s_withdrawableTokens[msg.sender];
    } else if (s_withdrawableTokens[msg.sender] < amount) {
      revert InsufficientBalance();
    }
    s_withdrawableTokens[msg.sender] -= amount;
    IFunctionsSubscriptions(address(_getRouter())).oracleWithdraw(recipient, amount);
  }

  /// @inheritdoc IFunctionsBilling
  /// @dev Only callable by the Coordinator owner
  function oracleWithdrawAll() external {
    _onlyOwner();
    _disperseFeePool();

    address[] memory transmitters = _getTransmitters();

    // Bounded by "maxNumOracles" on OCR2Abstract.sol
    for (uint256 i = 0; i < transmitters.length; ++i) {
      uint96 balance = s_withdrawableTokens[transmitters[i]];
      if (balance > 0) {
        s_withdrawableTokens[transmitters[i]] = 0;
        IFunctionsSubscriptions(address(_getRouter())).oracleWithdraw(transmitters[i], balance);
      }
    }
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

  // Overriden in FunctionsCoordinator.sol
  function _onlyOwner() internal view virtual;
}
