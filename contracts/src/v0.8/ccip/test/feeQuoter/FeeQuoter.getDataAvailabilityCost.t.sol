// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_getDataAvailabilityCost is FeeQuoterSetup {
  function test_EmptyMessageCalculatesDataAvailabilityCost() public {
    uint256 dataAvailabilityCostUSD =
      s_feeQuoter.getDataAvailabilityCost(DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, 0, 0, 0);

    FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);

    uint256 dataAvailabilityGas = destChainConfig.destDataAvailabilityOverheadGas
      + destChainConfig.destGasPerDataAvailabilityByte * Internal.MESSAGE_FIXED_BYTES;
    uint256 expectedDataAvailabilityCostUSD =
      USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas * destChainConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);

    // Test that the cost is destination chain specific
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR + 1;
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityOverheadGas =
      destChainConfig.destDataAvailabilityOverheadGas * 2;
    destChainConfigArgs[0].destChainConfig.destGasPerDataAvailabilityByte =
      destChainConfig.destGasPerDataAvailabilityByte * 2;
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityMultiplierBps =
      destChainConfig.destDataAvailabilityMultiplierBps * 2;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR + 1);
    uint256 dataAvailabilityCostUSD2 =
      s_feeQuoter.getDataAvailabilityCost(DEST_CHAIN_SELECTOR + 1, USD_PER_DATA_AVAILABILITY_GAS, 0, 0, 0);
    dataAvailabilityGas = destChainConfig.destDataAvailabilityOverheadGas
      + destChainConfig.destGasPerDataAvailabilityByte * Internal.MESSAGE_FIXED_BYTES;
    expectedDataAvailabilityCostUSD =
      USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas * destChainConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD2);
    assertFalse(dataAvailabilityCostUSD == dataAvailabilityCostUSD2);
  }

  function test_SimpleMessageCalculatesDataAvailabilityCost() public view {
    uint256 dataAvailabilityCostUSD =
      s_feeQuoter.getDataAvailabilityCost(DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, 100, 5, 50);

    FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);

    uint256 dataAvailabilityLengthBytes =
      Internal.MESSAGE_FIXED_BYTES + 100 + (5 * Internal.MESSAGE_FIXED_BYTES_PER_TOKEN) + 50;
    uint256 dataAvailabilityGas = destChainConfig.destDataAvailabilityOverheadGas
      + destChainConfig.destGasPerDataAvailabilityByte * dataAvailabilityLengthBytes;
    uint256 expectedDataAvailabilityCostUSD =
      USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas * destChainConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);
  }

  function test_SimpleMessageCalculatesDataAvailabilityCostUnsupportedDestChainSelector() public view {
    uint256 dataAvailabilityCostUSD = s_feeQuoter.getDataAvailabilityCost(0, USD_PER_DATA_AVAILABILITY_GAS, 100, 5, 50);

    assertEq(dataAvailabilityCostUSD, 0);
  }

  function testFuzz_ZeroDataAvailabilityGasPriceAlwaysCalculatesZeroDataAvailabilityCost_Success(
    uint64 messageDataLength,
    uint32 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) public view {
    uint256 dataAvailabilityCostUSD = s_feeQuoter.getDataAvailabilityCost(
      DEST_CHAIN_SELECTOR, 0, messageDataLength, numberOfTokens, tokenTransferBytesOverhead
    );

    assertEq(0, dataAvailabilityCostUSD);
  }

  function testFuzz_CalculateDataAvailabilityCost_Success(
    uint64 destChainSelector,
    uint32 destDataAvailabilityOverheadGas,
    uint16 destGasPerDataAvailabilityByte,
    uint16 destDataAvailabilityMultiplierBps,
    uint112 dataAvailabilityGasPrice,
    uint64 messageDataLength,
    uint32 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) public {
    vm.assume(destChainSelector != 0);
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](1);
    FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(destChainSelector);
    destChainConfigArgs[0] =
      FeeQuoter.DestChainConfigArgs({destChainSelector: destChainSelector, destChainConfig: destChainConfig});
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityOverheadGas = destDataAvailabilityOverheadGas;
    destChainConfigArgs[0].destChainConfig.destGasPerDataAvailabilityByte = destGasPerDataAvailabilityByte;
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityMultiplierBps = destDataAvailabilityMultiplierBps;
    destChainConfigArgs[0].destChainConfig.defaultTxGasLimit = GAS_LIMIT;
    destChainConfigArgs[0].destChainConfig.maxPerMsgGasLimit = GAS_LIMIT;
    destChainConfigArgs[0].destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_EVM;

    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    uint256 dataAvailabilityCostUSD = s_feeQuoter.getDataAvailabilityCost(
      destChainConfigArgs[0].destChainSelector,
      dataAvailabilityGasPrice,
      messageDataLength,
      numberOfTokens,
      tokenTransferBytesOverhead
    );

    uint256 dataAvailabilityLengthBytes = Internal.MESSAGE_FIXED_BYTES + messageDataLength
      + (numberOfTokens * Internal.MESSAGE_FIXED_BYTES_PER_TOKEN) + tokenTransferBytesOverhead;

    uint256 dataAvailabilityGas =
      destDataAvailabilityOverheadGas + destGasPerDataAvailabilityByte * dataAvailabilityLengthBytes;
    uint256 expectedDataAvailabilityCostUSD =
      dataAvailabilityGasPrice * dataAvailabilityGas * destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);
  }
}
