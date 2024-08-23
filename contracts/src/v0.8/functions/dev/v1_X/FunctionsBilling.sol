// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IFunctionsSubscriptions} from "./interfaces/IFunctionsSubscriptions.sol";
import {AggregatorV3Interface} from "../../../shared/interfaces/AggregatorV3Interface.sol";
import {IFunctionsBilling, FunctionsBillingConfig} from "./interfaces/IFunctionsBilling.sol";

import {Routable} from "./Routable.sol";
import {FunctionsResponse} from "./libraries/FunctionsResponse.sol";

import {SafeCast} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/math/SafeCast.sol";

import {ChainSpecificUtil} from "./libraries/ChainSpecificUtil.sol";

/// @title Functions Billing contract
/// @notice Contract that calculates payment from users to the nodes of the Decentralized Oracle Network (DON).
abstract contract FunctionsBilling is Routable, IFunctionsBilling {
  using FunctionsResponse for FunctionsResponse.RequestMeta;
  using FunctionsResponse for FunctionsResponse.Commitment;
  using FunctionsResponse for FunctionsResponse.FulfillResult;

  uint256 private constant REASONABLE_GAS_PRICE_CEILING = 1_000_000_000_000_000; // 1 million gwei

  event RequestBilled(
    bytes32 indexed requestId,
    uint96 juelsPerGas,
    uint256 l1FeeShareWei,
    uint96 callbackCostJuels,
    uint72 donFeeJuels,
    uint72 adminFeeJuels,
    uint72 operationFeeJuels
  );

  // ================================================================
  // |                  Request Commitment state                    |
  // ================================================================

  mapping(bytes32 requestId => bytes32 commitmentHash) private s_requestCommitments;

  event CommitmentDeleted(bytes32 requestId);

  FunctionsBillingConfig private s_config;

  event ConfigUpdated(FunctionsBillingConfig config);

  error UnsupportedRequestDataVersion();
  error InsufficientBalance();
  error InvalidSubscription();
  error UnauthorizedSender();
  error MustBeSubOwner(address owner);
  error InvalidLinkWeiPrice(int256 linkWei);
  error InvalidUsdLinkPrice(int256 usdLink);
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
  AggregatorV3Interface private s_linkToUsdFeed;

  // ================================================================
  // |                       Initialization                         |
  // ================================================================
  constructor(
    address router,
    FunctionsBillingConfig memory config,
    address linkToNativeFeed,
    address linkToUsdFeed
  ) Routable(router) {
    s_linkToNativeFeed = AggregatorV3Interface(linkToNativeFeed);
    s_linkToUsdFeed = AggregatorV3Interface(linkToUsdFeed);

    updateConfig(config);
  }

  // ================================================================
  // |                        Configuration                         |
  // ================================================================

  /// @notice Gets the Chainlink Coordinator's billing configuration
  /// @return config
  function getConfig() external view returns (FunctionsBillingConfig memory) {
    return s_config;
  }

  /// @notice Sets the Chainlink Coordinator's billing configuration
  /// @param config - See the contents of the FunctionsBillingConfig struct in IFunctionsBilling.sol for more information
  function updateConfig(FunctionsBillingConfig memory config) public {
    _onlyOwner();

    s_config = config;
    emit ConfigUpdated(config);
  }

  // ================================================================
  // |                       Fee Calculation                        |
  // ================================================================

  /// @inheritdoc IFunctionsBilling
  function getDONFeeJuels(bytes memory /* requestData */) public view override returns (uint72) {
    // s_config.donFee is in cents of USD. Convert to dollars amount then get amount of Juels.
    return SafeCast.toUint72(_getJuelsFromUsd(s_config.donFeeCentsUsd) / 100);
  }

  /// @inheritdoc IFunctionsBilling
  function getOperationFeeJuels() public view override returns (uint72) {
    // s_config.donFee is in cents of USD. Convert to dollars then get amount of Juels.
    return SafeCast.toUint72(_getJuelsFromUsd(s_config.operationFeeCentsUsd) / 100);
  }

  /// @inheritdoc IFunctionsBilling
  function getAdminFeeJuels() public view override returns (uint72) {
    return _getRouter().getAdminFee();
  }

  /// @inheritdoc IFunctionsBilling
  function getWeiPerUnitLink() public view returns (uint256) {
    (, int256 weiPerUnitLink, , uint256 timestamp, ) = s_linkToNativeFeed.latestRoundData();
    // Only fallback if feedStalenessSeconds is set
    // solhint-disable-next-line not-rely-on-time
    if (s_config.feedStalenessSeconds < block.timestamp - timestamp && s_config.feedStalenessSeconds > 0) {
      return s_config.fallbackNativePerUnitLink;
    }
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    return uint256(weiPerUnitLink);
  }

  function _getJuelsFromWei(uint256 amountWei) private view returns (uint96) {
    // (1e18 juels/link) * wei / (wei/link) = juels
    // There are only 1e9*1e18 = 1e27 juels in existence, should not exceed uint96 (2^96 ~ 7e28)
    return SafeCast.toUint96((1e18 * amountWei) / getWeiPerUnitLink());
  }

  /// @inheritdoc IFunctionsBilling
  function getUsdPerUnitLink() public view returns (uint256, uint8) {
    (, int256 usdPerUnitLink, , uint256 timestamp, ) = s_linkToUsdFeed.latestRoundData();
    // Only fallback if feedStalenessSeconds is set
    // solhint-disable-next-line not-rely-on-time
    if (s_config.feedStalenessSeconds < block.timestamp - timestamp && s_config.feedStalenessSeconds > 0) {
      return (s_config.fallbackUsdPerUnitLink, s_config.fallbackUsdPerUnitLinkDecimals);
    }
    if (usdPerUnitLink <= 0) {
      revert InvalidUsdLinkPrice(usdPerUnitLink);
    }
    return (uint256(usdPerUnitLink), s_linkToUsdFeed.decimals());
  }

  function _getJuelsFromUsd(uint256 amountUsd) private view returns (uint96) {
    (uint256 usdPerLink, uint8 decimals) = getUsdPerUnitLink();
    // (usd) * (10**18 juels/link) * (10**decimals) / (link / usd) = juels
    // There are only 1e9*1e18 = 1e27 juels in existence, should not exceed uint96 (2^96 ~ 7e28)
    return SafeCast.toUint96((amountUsd * 10 ** (18 + decimals)) / usdPerLink);
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
    uint72 adminFee = getAdminFeeJuels();
    uint72 donFee = getDONFeeJuels(data);
    uint72 operationFee = getOperationFeeJuels();
    return _calculateCostEstimate(callbackGasLimit, gasPriceWei, donFee, adminFee, operationFee);
  }

  /// @notice Estimate the cost in Juels of LINK
  // that will be charged to a subscription to fulfill a Functions request
  // Gas Price can be overestimated to account for flucuations between request and response time
  function _calculateCostEstimate(
    uint32 callbackGasLimit,
    uint256 gasPriceWei,
    uint72 donFeeJuels,
    uint72 adminFeeJuels,
    uint72 operationFeeJuels
  ) internal view returns (uint96) {
    // If gas price is less than the minimum fulfillment gas price, override to using the minimum
    if (gasPriceWei < s_config.minimumEstimateGasPriceWei) {
      gasPriceWei = s_config.minimumEstimateGasPriceWei;
    }

    uint256 executionGas = s_config.gasOverheadBeforeCallback + s_config.gasOverheadAfterCallback + callbackGasLimit;
    uint256 l1FeeWei = ChainSpecificUtil._getL1FeeUpperLimit(s_config.transmitTxSizeBytes);
    uint256 totalFeeWei = (gasPriceWei * executionGas) + l1FeeWei;

    // Basis Points are 1/100th of 1%, divide by 10_000 to bring back to original units
    uint256 totalFeeWeiWithOverestimate = totalFeeWei +
      ((totalFeeWei * s_config.fulfillmentGasPriceOverEstimationBP) / 10_000);

    uint96 estimatedGasReimbursementJuels = _getJuelsFromWei(totalFeeWeiWithOverestimate);

    uint96 feesJuels = uint96(donFeeJuels) + uint96(adminFeeJuels) + uint96(operationFeeJuels);

    return estimatedGasReimbursementJuels + feesJuels;
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
  ) internal returns (FunctionsResponse.Commitment memory commitment, uint72 operationFee) {
    // Nodes should support all past versions of the structure
    if (request.dataVersion > s_config.maxSupportedRequestDataVersion) {
      revert UnsupportedRequestDataVersion();
    }

    uint72 donFee = getDONFeeJuels(request.data);
    operationFee = getOperationFeeJuels();
    uint96 estimatedTotalCostJuels = _calculateCostEstimate(
      request.callbackGasLimit,
      tx.gasprice,
      donFee,
      request.adminFee,
      operationFee
    );

    // Check that subscription can afford the estimated cost
    if ((request.availableBalance) < estimatedTotalCostJuels) {
      revert InsufficientBalance();
    }

    uint32 timeoutTimestamp = uint32(block.timestamp + s_config.requestTimeoutSeconds);
    bytes32 requestId = keccak256(
      abi.encode(
        address(this),
        request.requestingContract,
        request.subscriptionId,
        request.initiatedRequests + 1,
        keccak256(request.data),
        request.dataVersion,
        request.callbackGasLimit,
        estimatedTotalCostJuels,
        timeoutTimestamp,
        // solhint-disable-next-line avoid-tx-origin
        tx.origin
      )
    );

    commitment = FunctionsResponse.Commitment({
      adminFee: request.adminFee,
      coordinator: address(this),
      client: request.requestingContract,
      subscriptionId: request.subscriptionId,
      callbackGasLimit: request.callbackGasLimit,
      estimatedTotalCostJuels: estimatedTotalCostJuels,
      timeoutTimestamp: timeoutTimestamp,
      requestId: requestId,
      donFee: donFee,
      gasOverheadBeforeCallback: s_config.gasOverheadBeforeCallback,
      gasOverheadAfterCallback: s_config.gasOverheadAfterCallback
    });

    s_requestCommitments[requestId] = keccak256(abi.encode(commitment));

    return (commitment, operationFee);
  }

  /// @notice Finalize billing process for an Functions request by sending a callback to the Client contract and then charging the subscription
  /// @param requestId identifier for the request that was generated by the Registry in the beginBilling commitment
  /// @param response response data from DON consensus
  /// @param err error from DON consensus
  /// @param reportBatchSize the number of fulfillments in the transmitter's report
  /// @return result fulfillment result
  /// @dev Only callable by a node that has been approved on the Coordinator
  /// @dev simulated offchain to determine if sufficient balance is present to fulfill the request
  function _fulfillAndBill(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    bytes memory onchainMetadata,
    bytes memory /* offchainMetadata TODO: use in getDonFee() for dynamic billing */,
    uint8 reportBatchSize
  ) internal returns (FunctionsResponse.FulfillResult) {
    FunctionsResponse.Commitment memory commitment = abi.decode(onchainMetadata, (FunctionsResponse.Commitment));

    uint256 gasOverheadWei = (commitment.gasOverheadBeforeCallback + commitment.gasOverheadAfterCallback) * tx.gasprice;
    uint256 l1FeeShareWei = ChainSpecificUtil._getL1FeeUpperLimit(msg.data.length) / reportBatchSize;
    // Gas overhead without callback
    uint96 gasOverheadJuels = _getJuelsFromWei(gasOverheadWei + l1FeeShareWei);
    uint96 juelsPerGas = _getJuelsFromWei(tx.gasprice);

    // The Functions Router will perform the callback to the client contract
    (FunctionsResponse.FulfillResult resultCode, uint96 callbackCostJuels) = _getRouter().fulfill(
      response,
      err,
      juelsPerGas,
      // The following line represents: "cost without callback or admin fee, those will be added by the Router"
      // But because the _offchain_ Commitment is using operation fee in the place of the admin fee, this now adds admin fee (actually operation fee)
      // Admin fee is configured to 0 in the Router
      gasOverheadJuels + commitment.donFee + commitment.adminFee,
      msg.sender,
      FunctionsResponse.Commitment({
        adminFee: 0, // The Router should have adminFee set to 0. If it does not this will cause fulfillments to fail with INVALID_COMMITMENT instead of carrying out incorrect bookkeeping.
        coordinator: commitment.coordinator,
        client: commitment.client,
        subscriptionId: commitment.subscriptionId,
        callbackGasLimit: commitment.callbackGasLimit,
        estimatedTotalCostJuels: commitment.estimatedTotalCostJuels,
        timeoutTimestamp: commitment.timeoutTimestamp,
        requestId: commitment.requestId,
        donFee: commitment.donFee,
        gasOverheadBeforeCallback: commitment.gasOverheadBeforeCallback,
        gasOverheadAfterCallback: commitment.gasOverheadAfterCallback
      })
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
      s_withdrawableTokens[msg.sender] += gasOverheadJuels + callbackCostJuels;
      // Put donFee into the pool of fees, to be split later
      // Saves on storage writes that would otherwise be charged to the user
      s_feePool += commitment.donFee;
      // Pay the operation fee to the Coordinator owner
      s_withdrawableTokens[_owner()] += commitment.adminFee; // OperationFee is used in the slot for Admin Fee in the Offchain Commitment. Admin Fee is set to 0 in the Router (enforced by line 316 in FunctionsBilling.sol).
      emit RequestBilled({
        requestId: requestId,
        juelsPerGas: juelsPerGas,
        l1FeeShareWei: l1FeeShareWei,
        callbackCostJuels: callbackCostJuels,
        donFeeJuels: commitment.donFee,
        // The following two lines are because of OperationFee being used in the Offchain Commitment
        adminFeeJuels: 0,
        operationFeeJuels: commitment.adminFee
      });
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
    uint256 numberOfTransmitters = transmitters.length;
    if (numberOfTransmitters == 0) {
      revert NoTransmittersSet();
    }
    uint96 feePoolShare = s_feePool / uint96(numberOfTransmitters);
    if (feePoolShare == 0) {
      // Dust cannot be evenly distributed to all transmitters
      return;
    }
    // Bounded by "maxNumOracles" on OCR2Abstract.sol
    for (uint256 i = 0; i < numberOfTransmitters; ++i) {
      s_withdrawableTokens[transmitters[i]] += feePoolShare;
    }
    s_feePool -= feePoolShare * uint96(numberOfTransmitters);
  }

  // Overriden in FunctionsCoordinator.sol
  function _onlyOwner() internal view virtual;

  // Used in FunctionsCoordinator.sol
  function _isExistingRequest(bytes32 requestId) internal view returns (bool) {
    return s_requestCommitments[requestId] != bytes32(0);
  }

  // Overriden in FunctionsCoordinator.sol
  function _owner() internal view virtual returns (address owner);
}
