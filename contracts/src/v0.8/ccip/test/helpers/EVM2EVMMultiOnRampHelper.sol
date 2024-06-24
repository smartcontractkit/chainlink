// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import "../../onRamp/EVM2EVMMultiOnRamp.sol";
import {IgnoreContractSize} from "./IgnoreContractSize.sol";

contract EVM2EVMMultiOnRampHelper is EVM2EVMMultiOnRamp, IgnoreContractSize {
  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig,
    DestChainConfigArgs[] memory destChainConfigs,
    PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs,
    TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs
  )
    EVM2EVMMultiOnRamp(
      staticConfig,
      dynamicConfig,
      destChainConfigs,
      premiumMultiplierWeiPerEthArgs,
      tokenTransferFeeConfigArgs
    )
  {}

  function getDataAvailabilityCost(
    uint64 destChainSelector,
    uint112 dataAvailabilityGasPrice,
    uint256 messageDataLength,
    uint256 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) external view returns (uint256) {
    return _getDataAvailabilityCost(
      destChainSelector, dataAvailabilityGasPrice, messageDataLength, numberOfTokens, tokenTransferBytesOverhead
    );
  }

  function getTokenTransferCost(
    uint64 destChainSelector,
    address feeToken,
    uint224 feeTokenPrice,
    Client.EVMTokenAmount[] calldata tokenAmounts
  ) external view returns (uint256, uint32, uint32) {
    return _getTokenTransferCost(destChainSelector, feeToken, feeTokenPrice, tokenAmounts);
  }
}
