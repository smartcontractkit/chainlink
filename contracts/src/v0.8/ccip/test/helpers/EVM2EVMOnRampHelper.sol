// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import "../../onRamp/EVM2EVMOnRamp.sol";
import {IgnoreContractSize} from "./IgnoreContractSize.sol";

contract EVM2EVMOnRampHelper is EVM2EVMOnRamp, IgnoreContractSize {
  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig,
    RateLimiter.Config memory rateLimiterConfig,
    FeeTokenConfigArgs[] memory feeTokenConfigs,
    TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs,
    NopAndWeight[] memory nopsAndWeights
  )
    EVM2EVMOnRamp(
      staticConfig,
      dynamicConfig,
      rateLimiterConfig,
      feeTokenConfigs,
      tokenTransferFeeConfigArgs,
      nopsAndWeights
    )
  {}

  function getDataAvailabilityCost(
    uint112 dataAvailabilityGasPrice,
    uint256 messageDataLength,
    uint256 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) external view returns (uint256) {
    return
      _getDataAvailabilityCost(dataAvailabilityGasPrice, messageDataLength, numberOfTokens, tokenTransferBytesOverhead);
  }

  function getTokenTransferCost(
    address feeToken,
    uint224 feeTokenPrice,
    Client.EVMTokenAmount[] calldata tokenAmounts
  ) external view returns (uint256, uint32, uint32) {
    return _getTokenTransferCost(feeToken, feeTokenPrice, tokenAmounts);
  }

  function getSequenceNumber() external view returns (uint64) {
    return s_sequenceNumber;
  }
}
