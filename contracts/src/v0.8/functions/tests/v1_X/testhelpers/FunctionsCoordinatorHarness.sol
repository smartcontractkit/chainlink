// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsCoordinator} from "../../../dev/v1_X/FunctionsCoordinator.sol";
import {FunctionsBilling} from "../../../dev/v1_X/FunctionsBilling.sol";
import {FunctionsResponse} from "../../../dev/v1_X/libraries/FunctionsResponse.sol";
import {FunctionsBillingConfig} from "../../../dev/v1_X/interfaces/IFunctionsBilling.sol";

/// @title Functions Coordinator Test Harness
/// @notice Contract to expose internal functions for testing purposes
contract FunctionsCoordinatorHarness is FunctionsCoordinator {
  address s_linkToNativeFeed_HARNESS;
  address s_linkToUsdFeed_HARNESS;
  address s_router_HARNESS;

  constructor(
    address router,
    FunctionsBillingConfig memory config,
    address linkToNativeFeed,
    address linkToUsdFeed
  ) FunctionsCoordinator(router, config, linkToNativeFeed, linkToUsdFeed) {
    s_linkToNativeFeed_HARNESS = linkToNativeFeed;
    s_linkToUsdFeed_HARNESS = linkToUsdFeed;
    s_router_HARNESS = router;
  }

  function isTransmitter_HARNESS(address node) external view returns (bool) {
    return super._isTransmitter(node);
  }

  function beforeSetConfig_HARNESS(uint8 _f, bytes memory _onchainConfig) external {
    return super._beforeSetConfig(_f, _onchainConfig);
  }

  /// @dev Used by FunctionsBilling.sol
  function getTransmitters_HARNESS() external view returns (address[] memory) {
    return super._getTransmitters();
  }

  function report_HARNESS(DecodedReport memory decodedReport) external {
    return super._report(decodedReport);
  }

  function onlyOwner_HARNESS() external view {
    return super._onlyOwner();
  }

  // ================================================================
  // |                        Functions Billing                     |
  // ================================================================

  function getLinkToNativeFeed_HARNESS() external view returns (address) {
    return s_linkToNativeFeed_HARNESS;
  }

  function getLinkToUsdFeed_HARNESS() external view returns (address) {
    return s_linkToUsdFeed_HARNESS;
  }

  function getRouter_HARNESS() external view returns (address) {
    return s_router_HARNESS;
  }

  function calculateCostEstimate_HARNESS(
    uint32 callbackGasLimit,
    uint256 gasPriceWei,
    uint72 donFee,
    uint72 adminFee,
    uint72 operationFee
  ) external view returns (uint96) {
    return super._calculateCostEstimate(callbackGasLimit, gasPriceWei, donFee, adminFee, operationFee);
  }

  function startBilling_HARNESS(
    FunctionsResponse.RequestMeta memory request
  ) external returns (FunctionsResponse.Commitment memory commitment, uint72 operationFee) {
    return super._startBilling(request);
  }

  function fulfillAndBill_HARNESS(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    bytes memory onchainMetadata,
    bytes memory offchainMetadata,
    uint8 reportBatchSize
  ) external returns (FunctionsResponse.FulfillResult) {
    return super._fulfillAndBill(requestId, response, err, onchainMetadata, offchainMetadata, reportBatchSize);
  }

  function disperseFeePool_HARNESS() external {
    return super._disperseFeePool();
  }

  function owner_HARNESS() external view returns (address owner) {
    return super._owner();
  }

  // ================================================================
  // |                              OCR2                            |
  // ================================================================

  function configDigestFromConfigData_HARNESS(
    uint256 _chainId,
    address _contractAddress,
    uint64 _configCount,
    address[] memory _signers,
    address[] memory _transmitters,
    uint8 _f,
    bytes memory _onchainConfig,
    uint64 _encodedConfigVersion,
    bytes memory _encodedConfig
  ) internal pure returns (bytes32) {
    return
      super._configDigestFromConfigData(
        _chainId,
        _contractAddress,
        _configCount,
        _signers,
        _transmitters,
        _f,
        _onchainConfig,
        _encodedConfigVersion,
        _encodedConfig
      );
  }
}
