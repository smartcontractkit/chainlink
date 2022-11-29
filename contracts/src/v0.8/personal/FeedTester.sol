// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../interfaces/AggregatorV3Interface.sol";
import "../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";
import "../vendor/@eth-optimism/contracts/0.8.6/contracts/L2/predeploys/OVM_GasPriceOracle.sol";

contract FeedTester {
  AggregatorV3Interface public fastGasFeedOpt = AggregatorV3Interface(0x5ad5CAdeBc6908b3aeb378a56659d08391C4C043);
  AggregatorV3Interface public linkNativeFeedOpt = AggregatorV3Interface(0x9FF1c5b77fCe72f9AA291BbF1b53A03B478d8Cf2);

  AggregatorV3Interface public fastGasFeedArb = AggregatorV3Interface(0x116542f62410Ac122C73ED3bC478937e781c5333);
  AggregatorV3Interface public linkNativeFeedArb = AggregatorV3Interface(0xE07eb28DcE1EAC2e6ea30379320Db88ED4b8a871);

  OVM_GasPriceOracle public OPTIMISM_ORACLE = OVM_GasPriceOracle(0x420000000000000000000000000000000000000F);
  ArbGasInfo public ARB_NITRO_ORACLE = ArbGasInfo(0x000000000000000000000000000000000000006C);

  function getOptOracleData() external view returns (uint256) {
    return OPTIMISM_ORACLE.getL1Fee(bytes("abc"));
  }

  function getArbOracleData() external view returns (uint256) {
    return ARB_NITRO_ORACLE.getCurrentTxL1GasFees();
  }

  function getOptFastGas() external view returns (int256 optFastFeedValue, uint256 optFastTimestamp) {
    (, optFastFeedValue, , optFastTimestamp, ) = fastGasFeedOpt.latestRoundData();
    return (optFastFeedValue, optFastTimestamp);
  }

  function getArbFastGas() external view returns (int256 arbFastFeedValue, uint256 arbFastTimestamp) {
    (, arbFastFeedValue, , arbFastTimestamp, ) = fastGasFeedArb.latestRoundData();
    return (arbFastFeedValue, arbFastTimestamp);
  }

  function getOptLinkNative() external view returns (int256 optLinkNativeFeedValue, uint256 optLinkNativeTimestamp) {
    (, optLinkNativeFeedValue, , optLinkNativeTimestamp, ) = linkNativeFeedOpt.latestRoundData();
    return (optLinkNativeFeedValue, optLinkNativeTimestamp);
  }

  function getArbLinkNative() external view returns (int256 arbLinkNativeFeedValue, uint256 arbLinkNativeTimestamp) {
    (, arbLinkNativeFeedValue, , arbLinkNativeTimestamp, ) = linkNativeFeedArb.latestRoundData();
    return (arbLinkNativeFeedValue, arbLinkNativeTimestamp);
  }
}
