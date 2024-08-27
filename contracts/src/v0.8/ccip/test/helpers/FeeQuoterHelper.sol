// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";

contract FeeQuoterHelper is FeeQuoter {
  constructor(
    StaticConfig memory staticConfig,
    address[] memory priceUpdaters,
    address[] memory feeTokens,
    TokenPriceFeedUpdate[] memory tokenPriceFeeds,
    TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs,
    PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs,
    DestChainConfigArgs[] memory destChainConfigArgs
  )
    FeeQuoter(
      staticConfig,
      priceUpdaters,
      feeTokens,
      tokenPriceFeeds,
      tokenTransferFeeConfigArgs,
      premiumMultiplierWeiPerEthArgs,
      destChainConfigArgs
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
      s_destChainConfigs[destChainSelector],
      dataAvailabilityGasPrice,
      messageDataLength,
      numberOfTokens,
      tokenTransferBytesOverhead
    );
  }

  function getTokenTransferCost(
    uint64 destChainSelector,
    address feeToken,
    uint224 feeTokenPrice,
    Client.EVMTokenAmount[] calldata tokenAmounts
  ) external view returns (uint256, uint32, uint32) {
    return _getTokenTransferCost(
      s_destChainConfigs[destChainSelector], destChainSelector, feeToken, feeTokenPrice, tokenAmounts
    );
  }

  function parseEVMExtraArgsFromBytes(
    bytes calldata extraArgs,
    uint64 destChainSelector
  ) external view returns (Client.EVMExtraArgsV2 memory) {
    return _parseEVMExtraArgsFromBytes(extraArgs, s_destChainConfigs[destChainSelector]);
  }

  function parseEVMExtraArgsFromBytes(
    bytes calldata extraArgs,
    DestChainConfig memory destChainConfig
  ) external pure returns (Client.EVMExtraArgsV2 memory) {
    return _parseEVMExtraArgsFromBytes(extraArgs, destChainConfig);
  }

  function validateDestFamilyAddress(bytes4 chainFamilySelector, bytes memory destAddress) external pure {
    _validateDestFamilyAddress(chainFamilySelector, destAddress);
  }

  function calculateRebasedValue(
    uint8 dataFeedDecimal,
    uint8 tokenDecimal,
    uint256 feedValue
  ) external pure returns (uint224) {
    return _calculateRebasedValue(dataFeedDecimal, tokenDecimal, feedValue);
  }
}
