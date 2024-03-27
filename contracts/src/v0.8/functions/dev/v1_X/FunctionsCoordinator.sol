// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IFunctionsCoordinator} from "./interfaces/IFunctionsCoordinator.sol";
import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";

import {FunctionsBilling, FunctionsBillingConfig} from "./FunctionsBilling.sol";
import {OCR2Base} from "./ocr/OCR2Base.sol";
import {FunctionsResponse} from "./libraries/FunctionsResponse.sol";

/// @title Functions Coordinator contract
/// @notice Contract that nodes of a Decentralized Oracle Network (DON) interact with
contract FunctionsCoordinator is OCR2Base, IFunctionsCoordinator, FunctionsBilling {
  using FunctionsResponse for FunctionsResponse.RequestMeta;
  using FunctionsResponse for FunctionsResponse.Commitment;
  using FunctionsResponse for FunctionsResponse.FulfillResult;

  /// @inheritdoc ITypeAndVersion
  string public constant override typeAndVersion = "Functions Coordinator v1.3.0";

  event OracleRequest(
    bytes32 indexed requestId,
    address indexed requestingContract,
    address requestInitiator,
    uint64 subscriptionId,
    address subscriptionOwner,
    bytes data,
    uint16 dataVersion,
    bytes32 flags,
    uint64 callbackGasLimit,
    FunctionsResponse.Commitment commitment
  );
  event OracleResponse(bytes32 indexed requestId, address transmitter);

  error InconsistentReportData();
  error EmptyPublicKey();
  error UnauthorizedPublicKeyChange();

  bytes private s_donPublicKey;
  bytes private s_thresholdPublicKey;

  constructor(
    address router,
    FunctionsBillingConfig memory config,
    address linkToNativeFeed,
    address linkToUsdFeed
  ) OCR2Base() FunctionsBilling(router, config, linkToNativeFeed, linkToUsdFeed) {}

  /// @inheritdoc IFunctionsCoordinator
  function getThresholdPublicKey() external view override returns (bytes memory) {
    if (s_thresholdPublicKey.length == 0) {
      revert EmptyPublicKey();
    }
    return s_thresholdPublicKey;
  }

  /// @inheritdoc IFunctionsCoordinator
  function setThresholdPublicKey(bytes calldata thresholdPublicKey) external override onlyOwner {
    if (thresholdPublicKey.length == 0) {
      revert EmptyPublicKey();
    }
    s_thresholdPublicKey = thresholdPublicKey;
  }

  /// @inheritdoc IFunctionsCoordinator
  function getDONPublicKey() external view override returns (bytes memory) {
    if (s_donPublicKey.length == 0) {
      revert EmptyPublicKey();
    }
    return s_donPublicKey;
  }

  /// @inheritdoc IFunctionsCoordinator
  function setDONPublicKey(bytes calldata donPublicKey) external override onlyOwner {
    if (donPublicKey.length == 0) {
      revert EmptyPublicKey();
    }
    s_donPublicKey = donPublicKey;
  }

  /// @dev check if node is in current transmitter list
  function _isTransmitter(address node) internal view returns (bool) {
    // Bounded by "maxNumOracles" on OCR2Abstract.sol
    for (uint256 i = 0; i < s_transmitters.length; ++i) {
      if (s_transmitters[i] == node) {
        return true;
      }
    }
    return false;
  }

  /// @inheritdoc IFunctionsCoordinator
  function startRequest(
    FunctionsResponse.RequestMeta calldata request
  ) external override onlyRouter returns (FunctionsResponse.Commitment memory commitment) {
    uint72 operationFee;
    (commitment, operationFee) = _startBilling(request);

    emit OracleRequest(
      commitment.requestId,
      request.requestingContract,
      // solhint-disable-next-line avoid-tx-origin
      tx.origin,
      request.subscriptionId,
      request.subscriptionOwner,
      request.data,
      request.dataVersion,
      request.flags,
      request.callbackGasLimit,
      FunctionsResponse.Commitment({
        coordinator: commitment.coordinator,
        client: commitment.client,
        subscriptionId: commitment.subscriptionId,
        callbackGasLimit: commitment.callbackGasLimit,
        estimatedTotalCostJuels: commitment.estimatedTotalCostJuels,
        timeoutTimestamp: commitment.timeoutTimestamp,
        requestId: commitment.requestId,
        donFee: commitment.donFee,
        gasOverheadBeforeCallback: commitment.gasOverheadBeforeCallback,
        gasOverheadAfterCallback: commitment.gasOverheadAfterCallback,
        // The following line is done to use the Coordinator's operationFee in place of the Router's operation fee
        // With this in place the Router.adminFee must be set to 0 in the Router.
        adminFee: operationFee
      })
    );

    return commitment;
  }

  /// @dev DON fees are pooled together. If the OCR configuration is going to change, these need to be distributed.
  function _beforeSetConfig(uint8 /* _f */, bytes memory /* _onchainConfig */) internal override {
    if (_getTransmitters().length > 0) {
      _disperseFeePool();
    }
  }

  /// @dev Used by FunctionsBilling.sol
  function _getTransmitters() internal view override returns (address[] memory) {
    return s_transmitters;
  }

  function _beforeTransmit(
    bytes calldata report
  ) internal view override returns (bool shouldStop, DecodedReport memory decodedReport) {
    (
      bytes32[] memory requestIds,
      bytes[] memory results,
      bytes[] memory errors,
      bytes[] memory onchainMetadata,
      bytes[] memory offchainMetadata
    ) = abi.decode(report, (bytes32[], bytes[], bytes[], bytes[], bytes[]));
    uint256 numberOfFulfillments = uint8(requestIds.length);

    if (
      numberOfFulfillments == 0 ||
      numberOfFulfillments != results.length ||
      numberOfFulfillments != errors.length ||
      numberOfFulfillments != onchainMetadata.length ||
      numberOfFulfillments != offchainMetadata.length
    ) {
      revert ReportInvalid("Fields must be equal length");
    }

    for (uint256 i = 0; i < numberOfFulfillments; ++i) {
      if (_isExistingRequest(requestIds[i])) {
        // If there is an existing request, validate report
        // Leave shouldStop to default, false
        break;
      }
      if (i == numberOfFulfillments - 1) {
        // If the last fulfillment on the report does not exist, then all are duplicates
        // Indicate that it's safe to stop to save on the gas of validating the report
        shouldStop = true;
      }
    }

    return (
      shouldStop,
      DecodedReport({
        requestIds: requestIds,
        results: results,
        errors: errors,
        onchainMetadata: onchainMetadata,
        offchainMetadata: offchainMetadata
      })
    );
  }

  /// @dev Report hook called within OCR2Base.sol
  function _report(DecodedReport memory decodedReport) internal override {
    uint256 numberOfFulfillments = uint8(decodedReport.requestIds.length);

    // Bounded by "MaxRequestBatchSize" on the Job's ReportingPluginConfig
    for (uint256 i = 0; i < numberOfFulfillments; ++i) {
      FunctionsResponse.FulfillResult result = FunctionsResponse.FulfillResult(
        _fulfillAndBill(
          decodedReport.requestIds[i],
          decodedReport.results[i],
          decodedReport.errors[i],
          decodedReport.onchainMetadata[i],
          decodedReport.offchainMetadata[i],
          uint8(numberOfFulfillments) // will not exceed "MaxRequestBatchSize" on the Job's ReportingPluginConfig
        )
      );

      // Emit on successfully processing the fulfillment
      // In these two fulfillment results the user has been charged
      // Otherwise, the DON will re-try
      if (
        result == FunctionsResponse.FulfillResult.FULFILLED ||
        result == FunctionsResponse.FulfillResult.USER_CALLBACK_ERROR
      ) {
        emit OracleResponse(decodedReport.requestIds[i], msg.sender);
      }
    }
  }

  /// @dev Used in FunctionsBilling.sol
  function _onlyOwner() internal view override {
    _validateOwnership();
  }

  /// @dev Used in FunctionsBilling.sol
  function _owner() internal view override returns (address owner) {
    return this.owner();
  }
}
