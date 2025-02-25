// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

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
      s_destChainConfigs[destChainSelector].defaultTokenFeeUSDCents,
      s_destChainConfigs[destChainSelector].defaultTokenDestGasOverhead,
      destChainSelector,
      feeToken,
      feeTokenPrice,
      tokenAmounts
    );
  }

  function parseEVMExtraArgsFromBytes(
    bytes calldata extraArgs,
    uint64 destChainSelector
  ) external view returns (Client.EVMExtraArgsV2 memory) {
    return _parseEVMExtraArgsFromBytes(
      extraArgs,
      s_destChainConfigs[destChainSelector].defaultTxGasLimit,
      s_destChainConfigs[destChainSelector].maxPerMsgGasLimit,
      s_destChainConfigs[destChainSelector].enforceOutOfOrder
    );
  }

  function parseEVMExtraArgsFromBytes(
    bytes calldata extraArgs,
    uint64 destChainSelector,
    bool enforceOutOfOrder
  ) external view returns (Client.EVMExtraArgsV2 memory) {
    return _parseEVMExtraArgsFromBytes(
      extraArgs,
      s_destChainConfigs[destChainSelector].defaultTxGasLimit,
      s_destChainConfigs[destChainSelector].maxPerMsgGasLimit,
      enforceOutOfOrder
    );
  }

  function parseSVMExtraArgsFromBytes(
    bytes calldata extraArgs,
    DestChainConfig memory destChainConfig
  ) external pure returns (Client.SVMExtraArgsV1 memory) {
    return _parseSVMExtraArgsFromBytes(extraArgs, destChainConfig.maxPerMsgGasLimit, destChainConfig.enforceOutOfOrder);
  }

  function processChainFamilySelector(
    uint64 chainFamilySelector,
    bytes calldata messageReceiver,
    bytes calldata extraArgs
  ) external view returns (bytes memory, bool, bytes memory) {
    return _processChainFamilySelector(chainFamilySelector, messageReceiver, extraArgs);
  }

  function validateDestFamilyAddress(
    bytes4 chainFamilySelector,
    bytes memory destAddress,
    uint256 gasLimit
  ) external pure {
    _validateDestFamilyAddress(chainFamilySelector, destAddress, gasLimit);
  }

  function calculateRebasedValue(
    uint8 dataFeedDecimal,
    uint8 tokenDecimal,
    uint256 feedValue
  ) external pure returns (uint224) {
    return _calculateRebasedValue(dataFeedDecimal, tokenDecimal, feedValue);
  }
}
