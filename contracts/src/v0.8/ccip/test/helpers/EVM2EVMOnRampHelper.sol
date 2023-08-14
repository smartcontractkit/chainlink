// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import "../../onRamp/EVM2EVMOnRamp.sol";
import {IgnoreContractSize} from "./IgnoreContractSize.sol";

contract EVM2EVMOnRampHelper is EVM2EVMOnRamp, IgnoreContractSize {
  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig,
    Internal.PoolUpdate[] memory tokensAndPools,
    address[] memory allowlist,
    RateLimiter.Config memory rateLimiterConfig,
    FeeTokenConfigArgs[] memory feeTokenConfigs,
    TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs,
    NopAndWeight[] memory nopsAndWeights
  )
    EVM2EVMOnRamp(
      staticConfig,
      dynamicConfig,
      tokensAndPools,
      allowlist,
      rateLimiterConfig,
      feeTokenConfigs,
      tokenTransferFeeConfigArgs,
      nopsAndWeights
    )
  {}

  function getTokenTransferFeeUSD(
    address feeToken,
    uint192 feeTokenPrice,
    Client.EVMTokenAmount[] calldata tokenAmounts,
    FeeTokenConfig memory feeTokenConfig
  ) external view returns (uint256) {
    return _getTokenTransferFeeUSD(feeToken, feeTokenPrice, tokenAmounts, feeTokenConfig);
  }
}
