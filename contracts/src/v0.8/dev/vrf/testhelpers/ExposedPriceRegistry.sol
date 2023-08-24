// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../VRFV2PlusPriceRegistry.sol";

contract ExposedPriceRegistry is VRFV2PlusPriceRegistry {
  constructor(address linkEthFeed, address linkUSDFeed, address ethUSDFeed) VRFV2PlusPriceRegistry(linkEthFeed, linkUSDFeed, ethUSDFeed){

  }

  function calculateFlatFeeFromUSDExternal(uint40 fulfillmentFlatFeeUSD, AggregatorV3Interface feed) external view returns (uint256) {
    return calculateFlatFeeFromUSD(fulfillmentFlatFeeUSD, feed);
  }

  function getUSDFeedDataExternal(AggregatorV3Interface feed) external view returns (int256 answer, uint8 decimals) {
    return getUSDFeedData(feed);
  }

  function getLINKEthFeedDataExternal() external view returns (int256 answer) {
    return getLINKEthFeedData();
  }
}
